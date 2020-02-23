package main

import (
	"fmt"
	"github.com/cipepser/go-failure-sample/db"
	"github.com/morikuni/failure"
)

func init() {
	_ = db.NewCustomer("Alice", "alice@example.com")
	_ = db.NewCustomer("Bob", "bob@example.com")
}

func main() {
	c := db.NewClient("user")
	userId := 0
	_, err := c.GetName(userId)
	if failure.Is(err, db.NotFound) {
		fmt.Println("error occurred: NotFound")
	}

	if err := db.ShowCustomers(); err != nil {
		panic(err)
	}
}
