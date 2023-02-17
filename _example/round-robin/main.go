package main

import (
	"fmt"
	"net/url"

	roundrobin "github.com/appleboy/loadbalancer-algorithms/round-robin"
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

	fmt.Println(rb.Next().Host)
	fmt.Println(rb.Next().Host)
	fmt.Println(rb.Next().Host)
	fmt.Println(rb.Next().Host)
	fmt.Println(rb.Next().Host)
}
