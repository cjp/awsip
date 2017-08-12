package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Looks up current PrivateIpAddress of a named EC2 instance\n\n")
	fmt.Fprintf(os.Stderr, "usage: %s <instance_name>\n", filepath.Base(os.Args[0]))
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(64)
	}

	inst := os.Args[1]

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ec2client := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(inst)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	resp, err := ec2client.DescribeInstances(params)
	if err != nil {
		log.Fatal("Error listing instances", err)
	}

	var ip string = *resp.Reservations[0].Instances[0].PrivateIpAddress

	if len(ip) > 0 {
		fmt.Println(ip)
	} else {
		fmt.Fprintf(os.Stderr, "Not found\n")
		os.Exit(1)
	}
}
