package db

import (
	"fmt"
	"github.com/morikuni/failure"
)

var (
	NotFound failure.StringCode = "NotFound"
)

type Customer struct {
	id      int
	name    string
	address string
}

func NewCustomer(name, address string) *Customer {
	c := &Customer{
		id:      len(customers),
		name:    name,
		address: address,
	}
	customers = append(customers, c)
	return c
}

var customers []*Customer

func ShowCustomers() error {
	for _, c := range customers {
		fmt.Printf("%+v\n", *c)
	}
	return nil
}

type Client struct {
	user string
}

func NewClient(user string) *Client {
	return &Client{
		user: user,
	}
}

func (c *Client) GetName(id int) (string, error) {
	for _, c := range customers {
		if c.id == id {
			return c.name, nil
		}
 	}
	return "", failure.New(NotFound)
}
