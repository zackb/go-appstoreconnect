package main

import (
	"fmt"
	"github.com/zackb/go-appstoreconnect/appstoreconnect"
)

func main() {
	client, err := appstoreconnect.NewClientFromCredentialsFile("credentials.yml")
	if err != nil {
		panic(err)
	}

	sr := appstoreconnect.NewSalesReport("MONTHLY")
	b, err := client.Get(sr)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
