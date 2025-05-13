package client

import (
	"github.com/google/uuid"
)

type Client struct {
	Chips map[uint8]uint64
}

func NewClient(id uuid.UUID) *Client {
	return &Client{
		Chips: nil,
	}
}
