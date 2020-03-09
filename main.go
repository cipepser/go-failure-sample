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
	if err := c.CheckPermitted("mallory@example.com"); err != nil {
		fmt.Println("============ Error ============")
		fmt.Printf("Error = %v\n", err)

		code, _ := failure.CodeOf(err)
		fmt.Printf("Code = %v\n", code)

		msg, _ := failure.MessageOf(err)
		fmt.Printf("Message = %v\n", msg)

		cs, _ := failure.CallStackOf(err)
		fmt.Printf("CallStack = %v\n", cs)

		fmt.Printf("Cause = %v\n", failure.CauseOf(err))

		fmt.Println()
		fmt.Println("============ Detail ============")
		fmt.Printf("%+v\n", err)

		fmt.Println()
		fmt.Println("============ Unwrap ============")
		iter := failure.NewIterator(err)
		for iter.Next() {
			fmt.Println("***********")
			err := iter.Error()
			fmt.Printf("Error = %v\n", err)
			code, ok := failure.CodeOf(err)
			if !ok {
				fmt.Println("not ok")
				continue
			}
			switch code {
			case db.Forbidden:
				fmt.Println(db.Forbidden)
			case db.NotFound:
				fmt.Println(db.NotFound)
			default:
				fmt.Println("default")
			}
		}
	}
}
