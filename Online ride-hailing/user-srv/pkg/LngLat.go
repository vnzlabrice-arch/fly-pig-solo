package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type GeoCodeResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Geocodes []struct {
		FormattedAddress string   `json:"formatted_address"`
		Province         string   `json:"province"`
		City             string   `json:"city"`
		Citycode         string   `json:"citycode"`
		District         string   `json:"district"`
		Township         []string `json:"township"`
		Neighborhood     []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"neighborhood"`
		Building []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"building"`
		Street   []string `json:"street"`
		Number   []string `json:"number"`
		Location string   `json:"location"`
		Level    string   `json:"level"`
	} `json:"geocodes"`
}

func AddressToLngLat(address string) (float64, float64, error) {
	// 高德地图API key，需要替换为实际的key
	apiKey := "907ccf1bfdf101822c8c65b6eb2f6a8d"

	// 构建请求URL
	baseURL := "https://restapi.amap.com/v3/geocode/geo"
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("address", address)
	params.Add("output", "json")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 发送请求
	resp, err := http.Get(fullURL)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	// 解析响应
	var geoResp GeoCodeResponse
	err = json.Unmarshal(body, &geoResp)
	if err != nil {
		return 0, 0, err
	}

	// 检查响应状态
	if geoResp.Status != "1" {
		return 0, 0, fmt.Errorf("geocode failed: %s", geoResp.Info)
	}

	// 检查是否有结果
	if len(geoResp.Geocodes) == 0 {
		return 0, 0, fmt.Errorf("no geocode results")
	}

	// 解析经纬度
	var lng, lat float64
	_, err = fmt.Sscanf(geoResp.Geocodes[0].Location, "%f,%f", &lng, &lat)
	if err != nil {
		return 0, 0, err
	}

	return lng, lat, nil
}
