package main

import (
	"fmt"

	"github.com/teris-io/shortid"
)

func main() {
	fmt.Println("starting scripts")

	//cl := client.New(&client.Config{
	//BaseURL: "http://localhost:9000",
	//Client:  http.DefaultClient,
	//})

	//token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Impkb2VAZ21haWwuY29tIiwiZXhwIjoxNzE1ODQxNDA5fQ.DWAF9i07IgdnhZoyXHUsFQMzJ-6Hzz_-rD6rprpbCeo"

	//res, err := cl.ValidToken(token)
	//check(err)

	//fmt.Printf("%+v\n", res)
	for i := 1; i < 11; i++ {
		sid, err := shortid.Generate()
		check(err)
		fmt.Printf("%d. %s\n", i, sid)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
