package services

import (
	"encoding/json"
	"net/http"
)

type IPInfo struct {
	City string `json:"city"`
}

type IPService struct {
	client *http.Client
}

func NewIPService(client *http.Client) *IPService {
	return &IPService{client: client}
}

func (s *IPService) GetIPInfo(r *http.Request) (*IPInfo, error) {
	resp, err := s.client.Get("http://ip-api.com/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
