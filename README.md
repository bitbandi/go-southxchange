go-southxchange
===============

go-southxchange is an implementation of the SouthXchange API (public and private) in Golang.

## Import
	import "github.com/bitbandi/go-southxchange"
	
## Usage

In order to use the client with go's default http client settings you can do:

~~~ go
package main

import (
	"fmt"
	"github.com/bitbandi/go-southxchange"
)

const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
)

func main() {
	// southxchange client
	southxchange := southxchange.New(API_KEY, API_SECRET, "user-agent")

	// Get balances
	balances, err := southxchange.GetBalances()
	fmt.Println(err, balances)
}
~~~

In order to use custom settings for the http client do:

~~~ go
package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/bitbandi/go-southxchange"
)

const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
)

func main() {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	// southxchange client
	bc := southxchange.NewWithCustomHttpClient(API_KEY, API_SECRET, "user-agent", httpClient)

	// Get balances
	balances, err := southxchange.GetBalances()
	fmt.Println(err, balances)
}
~~~

See ["Examples" folder for more... examples](https://github.com/bitbandi/go-southxchange/blob/master/examples/southxchange.go)
