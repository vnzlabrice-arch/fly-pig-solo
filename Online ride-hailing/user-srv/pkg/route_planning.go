package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type DrivingDirectionResp struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Route    struct {
		Origin      string `json:"origin"`
		Destination string `json:"destination"`
		TaxiCost    string `json:"taxi_cost"`
		Paths       []struct {
			Distance    string `json:"distance"`
			Duration    string `json:"duration"`
			Restriction string `json:"restriction"`
			Steps       []struct {
				Instruction  string `json:"instruction"`
				Orientation  string `json:"orientation"`
				RoadName     string `json:"road_name,omitempty"`
				StepDistance string `json:"step_distance"`
			} `json:"steps"`
		} `json:"paths"`
	} `json:"route"`
}

func GetRouteByAddress(ctx context.Context, originAddr, destAddr string) (*DrivingDirectionResp, error) {
	// 1. 地址解析为经纬度
	originLng, originLat, err := AddressToLngLat(originAddr)
	if err != nil {
		return nil, fmt.Errorf("起点地址解析失败: %w", err)
	}
	originLoc := fmt.Sprintf("%f,%f", originLng, originLat)

	destLng, destLat, err := AddressToLngLat(destAddr)
	if err != nil {
		return nil, fmt.Errorf("终点地址解析失败: %w", err)
	}
	destLoc := fmt.Sprintf("%f,%f", destLng, destLat)

	// 2. 构造高德请求参数
	params := url.Values{}
	params.Set("origin", originLoc)
	params.Set("destination", destLoc)
	params.Set("key", apiKey)
	params.Set("strategy", "10")

	reqURL := fmt.Sprintf(
		"https://restapi.amap.com/v5/direction/driving?origin=%s&destination=%s&key=%s&strategy=%s",
		originLoc, // 111.553816,32.190534
		destLoc,   // 111.653077,32.263390
		apiKey,    // 你的 Key
		"0",       // 改成 0（速度优先），去掉 "10"
	)
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建路线请求失败: %w", err)
	}

	// 3. 发送请求并解析响应
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求高德路线接口失败: %w", err)
	}
	defer resp.Body.Close()

	var result DrivingDirectionResp
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析路线结果失败: %w", err)
	}

	// 4. 校验高德返回状态
	if result.Status != "1" {
		return nil, fmt.Errorf("路线规划失败，高德响应: status=%s, info=%s, infocode=%s",
			result.Status, result.Info, result.Infocode)
	}

	return &result, nil
}
