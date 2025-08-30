package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/huangrao121/CommunicationApp/BackendService/internal/gateway/websocket"
)

type GatewayService struct {
	messageServiceURL string
	httpClient        *http.Client
}

func NewGatewayService(messageServiceURL string) *GatewayService {
	return &GatewayService{
		messageServiceURL: messageServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// 实现MessageServiceClient接口
func (s *GatewayService) SendP2PMessage(ctx context.Context, req *websocket.SendP2PRequest) (*websocket.MessageResponse, error) {
	url := fmt.Sprintf("%s/api/v1/messages/p2p", s.messageServiceURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var messageResp websocket.MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&messageResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		messageResp.Success = false
		return &messageResp, fmt.Errorf("message service returned status: %d", resp.StatusCode)
	}

	messageResp.Success = true
	return &messageResp, nil
}

func (s *GatewayService) SendGroupMessage(ctx context.Context, req *websocket.SendGroupRequest) (*websocket.MessageResponse, error) {
	url := fmt.Sprintf("%s/api/v1/messages/group", s.messageServiceURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var messageResp websocket.MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&messageResp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		messageResp.Success = false
		return &messageResp, fmt.Errorf("message service returned status: %d", resp.StatusCode)
	}

	messageResp.Success = true
	return &messageResp, nil
}
