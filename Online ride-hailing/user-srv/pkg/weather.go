package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	weatherURL  = "https://restapi.amap.com/v3/weather/weatherInfo"
	districtURL = "https://restapi.amap.com/v3/config/district"
	extensions  = "base"
)

type WeatherResponse struct {
	Status  string `json:"status"`
	Info    string `json:"info"`
	Infocode string `json:"infocode"`
	Lives   []Live `json:"lives"`
}

type Live struct {
	Province   string `json:"province"`
	City       string `json:"city"`
	Weather    string `json:"weather"`
	Temperature string `json:"temperature"`
	Winddirection string `json:"winddirection"`
	Windpower  string `json:"windpower"`
	Humidity   string `json:"humidity"`
	Reporttime string `json:"reporttime"`
}

// 城市名 → adcode
func cityName2Adcode(cityName string) (string, error) {
	if cityName == "" {
		return "", fmt.Errorf("城市名不能为空")
	}

	params := url.Values{}
	params.Set("keywords", cityName)
	params.Set("subdistrict", "0")
	params.Set("key", apiKey)

	reqURL := fmt.Sprintf("%s?%s", districtURL, params.Encode())
	resp, err := HTTPClient.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("请求城市编码失败: %w", err)
	}
	defer resp.Body.Close()

	var res struct {
		Status    string `json:"status"`
		Info      string `json:"info"`
		Districts []struct {
			Adcode string `json:"adcode"`
		} `json:"districts"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("解析城市编码失败: %w", err)
	}

	if res.Status != "1" || len(res.Districts) == 0 {
		return "", fmt.Errorf("未找到城市: %s", cityName)
	}

	return res.Districts[0].Adcode, nil
}

// adcode → 天气
func queryWeatherByAdcode(adcode string) (*WeatherResponse, error) {
	params := url.Values{}
	params.Set("city", adcode)
	params.Set("key", apiKey)
	params.Set("extensions", extensions)

	reqURL := fmt.Sprintf("%s?%s", weatherURL, params.Encode())
	resp, err := HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("请求天气失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("天气接口状态码错误: %d", resp.StatusCode)
	}

	var result WeatherResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析天气失败: %w", err)
	}

	return &result, nil
}

func GetWeatherByCityName(cityName string) (*WeatherResponse, error) {
	adcode, err := cityName2Adcode(cityName)
	if err != nil {
		return nil, err
	}
	return queryWeatherByAdcode(adcode)
}

func GetWeatherByLocation(loc *LocationResult) (*WeatherResponse, error) {
	if loc == nil || loc.Adcode == "" {
		return nil, fmt.Errorf("IP 定位无 adcode")
	}
	return queryWeatherByAdcode(loc.Adcode)
}
