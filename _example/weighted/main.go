package main

import (
	"fmt"
	"net/url"

	"github.com/appleboy/loadbalancer-algorithms/weighted"
)

type server struct {
	url    *url.URL
	weight int
}

func main() {
	servers := []*server{
		{
			url: &url.URL{
				Host: "192.168.1.10",
			},
			weight: 4,
		},
		{
			url: &url.URL{
				Host: "192.168.1.11",
			},
			weight: 3,
		},
		{
			url: &url.URL{
				Host: "192.168.1.12",
			},
			weight: 2,
		},
	}

	rb, err := weighted.New()
	if err != nil {
		panic(err)
	}

	for _, v := range servers {
		_ = rb.AddServer(v.url, v.weight)
	}

	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
	fmt.Println(rb.NextServer().Host)
}
