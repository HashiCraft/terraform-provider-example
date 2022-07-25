package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MinecraftClient struct {
	baseURL string
}

type BlockRequest struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Z        int    `json:"z"`
	Material string `json:"material"`
	Facing   string `json:"facing"`
	Half     string `json:"half"`
}

func NewClient(url string) *MinecraftClient {
	return &MinecraftClient{url}
}

func (c *MinecraftClient) CreateBlock(block BlockRequest) error {
	url := fmt.Sprintf("%s/block", c.baseURL)

	// convert the object to json
	d, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("unable to marshal block to json: %s", err)
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(d))
	if err != nil {
		return fmt.Errorf("unable to create request: %s", err)
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got status %d", resp.StatusCode)
	}

	return nil
}

func (c *MinecraftClient) DeleteBlock(block BlockRequest) error {
	url := fmt.Sprintf("%s/block/%d/%d/%d", c.baseURL, block.X, block.Y, block.Z)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %s", err)
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got status %d", resp.StatusCode)
	}

	return nil
}

func (c *MinecraftClient) GetBlock(x, y, z int) (*BlockRequest, error) {
	url := fmt.Sprintf("%s/block/%d/%d/%d", c.baseURL, x, y, z)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %s", err)
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status 200, got status %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response from server %s", err)
	}

	br := &BlockRequest{}
	err = json.Unmarshal(d, br)

	return br, err
}
