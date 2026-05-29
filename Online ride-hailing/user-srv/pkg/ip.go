package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
)

const apiKey = "4ea8a94eb3efbdc5aaa66a0b0d4d59f4"

type LocationResult struct {
	Ip        string `json:"ip"`
	Status    string `json:"status"`
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Adcode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

func GetLocationByIP(ctx context.Context, ip string) (*LocationResult, error) {
	url := fmt.Sprintf(
		"https://restapi.amap.com/v3/ip?ip=%s&key=%s&output=json",
		ip,
		apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.InternalServer("REQUEST_ERROR", "创建请求失败")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.InternalServer("AMAP_API_ERROR", "请求高德失败")
	}
	defer resp.Body.Close()

	var result LocationResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.InternalServer("PARSE_ERROR", "解析失败")
	}

	result.Ip = ip

	if result.Status != "1" {
		return nil, errors.BadRequest("AMAP_ERROR", result.Info)
	}

	return &result, nil
}
