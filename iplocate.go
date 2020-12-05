package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

var (
	region  *string
	project *string
	ip      *string
)

//Initializes the flags using the standard flag package
func init() {
	region = flag.String("region", "us-central1", "region to search in")
	project = flag.String("project", "", "project name")
	ip = flag.String("ip", "", "the IP address to search")
	flag.Parse()
	if *project == "" || *ip == "" {
		panic("project and ip are required flags")
	}
}

func main() {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		panic(err)
	}

	computeService, err := compute.New(client)
	if err != nil {
		panic(err)
	}
	//	getIPs(project, ctx, computeService)
	checkInstances(*project, ctx, computeService, *ip)
	//regions := []string{"us-central1"}
	//getRegionalAddresses(regions, project, ctx, computeService)
}

func checkRegionalAddresses(regions []string, project string, ctx context.Context, computeService *compute.Service) {
	for _, region := range regions {
		req := computeService.Addresses.List(project, region)
		err := req.Pages(ctx, func(page *compute.AddressList) error {
			for _, address := range page.Items {
				fmt.Printf("%#v\n", address.Address)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}

//Checks the supplied IP against running instances only, both external and internal IPs. This is
//not comprehensive, of course, but suffices for the most common use case
func checkInstances(project string, ctx context.Context, computeService *compute.Service, ip string) {
	req := computeService.Instances.AggregatedList(project)
	err := req.Pages(ctx, func(page *compute.InstanceAggregatedList) (err error) {
		for _, instancesList := range page.Items {
			instances := instancesList.Instances
			for _, instance := range instances {
				if instance.Status == "RUNNING" {
					//fmt.Println("Instance: " + instance.Name)
					//fmt.Printf("%#v\n", instance)
					ifaces := instance.NetworkInterfaces
					for _, iface := range ifaces {
						//fmt.Printf("%s\n", iface.NetworkIP)
						if ip == iface.NetworkIP {
							fmt.Printf("%s\n", instance.Name)
							return
						}
						acs := iface.AccessConfigs
						for _, ac := range acs {
							if ip == ac.NatIP {
								fmt.Printf("%s\n", instance.Name)
								return
							}
						}
					}
				}
			}
		}
		fmt.Printf("Not found\n")
		return
	})
	if err != nil {
		panic(err)
	}
}

func checkIPs(project string, ctx context.Context, computeService *compute.Service) {
	req := computeService.Addresses.AggregatedList(project)
	err := req.Pages(ctx, func(page *compute.AddressAggregatedList) error {
		for name, addr := range page.Items {
			addresses := addr.Addresses
			for _, ip := range addresses {
				fmt.Printf("%v : %#v\n", name, ip.Address)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
