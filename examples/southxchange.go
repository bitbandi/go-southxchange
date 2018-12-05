package main

import (
	"fmt"

	"github.com/bitbandi/go-southxchange"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func main() {
	// southxchange client
	southxchange := southxchange.New(API_KEY, API_SECRET, "user-agent 1.0")

	// GetBalances
	balances, _ := southxchange.GetBalances()
	fmt.Println(len(balances))

	for i, _ := range balances {
		if balances[i].Currency == "BTC" {
			fmt.Println(balances[i].Available)
		}
	}

}
