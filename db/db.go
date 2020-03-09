package db

import (
	"bufio"
	"fmt"
	"github.com/morikuni/failure"
	"os"
)

const (
	WHITELIST = "whitelist.txt"
)

var (
	NotFound  failure.StringCode = "NotFound"
	FORBIDDEN failure.StringCode = "Forbidden"
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

// customers represent GLOBAL STATE
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

// CheckPermitted ONLY returns `nil` in case of an input address is permitted.
// Other cases, it returns an `FORBIDDEN` error.
// To wrap an error, we open WHITELIST file every time checking the list,
// not using init function.
func (c *Client) CheckPermitted(address string) error {
	f, err := os.Open(WHITELIST)
	if err != nil {
		return failure.Translate(err, FORBIDDEN,
			failure.Context{"package": "os"},
			failure.Messagef("failed to open %s", WHITELIST),
		)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		permittedAddress := sc.Text()
		if permittedAddress == address {
			return nil
		}
	}

	return failure.New(FORBIDDEN)
}
