package main

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func main() {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		panic(err)
	}
	compute, err := compute.New(client)
	if err != nil {
		panic(err)
	}
	fmt.Println("Got client %s", compute)

	project := "pivotal-canto-171605"

	req := compute.GlobalAddresses.List(project)
	err = req.Pages(ctx, func(page *compute.AddressList) error {
		for _, addr := range page.Items {
			fmt.Println(addr)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}
