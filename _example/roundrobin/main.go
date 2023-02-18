package main

import (
	"fmt"
	"net/url"

	"github.com/appleboy/loadbalancer-algorithms/roundrobin"
)

func main() {
	servers := []*url.URL{
		{Host: "192.168.1.10"},
		{Host: "192.168.1.11"},
		{Host: "192.168.1.12"},
		{Host: "192.168.1.13"},
	}

	rb, err := roundrobin.New(servers...)
	if err != nil {
		panic(err)
	}

	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
}
