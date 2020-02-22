package db

import (
	"github.com/morikuni/failure"
)

var NotFound failure.StringCode = "NotFound"

type Client struct {
	user string
}

func NewClient(user string) *Client {
	return &Client{
		user: user,
	}
}

func (c *Client) GetName(id int) (string, error) {
	return "", failure.New(NotFound)
}
