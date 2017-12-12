package main

import (
	"errors"
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

func init() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(64)
	}
}

// privateip accepts an array of ec2.Reservation pointers and returns first private IP found
func privateip(r []*ec2.Reservation) (ip string, err error) {
	if len(r) > 0 && len(r[0].Instances) > 0 {
		ip = *r[0].Instances[0].PrivateIpAddress
	} else {
		err = errors.New("PrivateIpAddress not found")
	}
	return
}

// findinst searches for instances tagged with name
func findinst(cli *ec2.EC2, n string) (resp *ec2.DescribeInstancesOutput, err error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(n)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	resp, err = cli.DescribeInstances(params)
	return
}

func main() {
	inst := os.Args[1]

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ec2client := ec2.New(sess)

	resp, err := findinst(ec2client, inst)
	if err != nil {
		log.Fatal("Error finding instances: ", err)
	}

	ip, err := privateip(resp.Reservations)
	if err != nil {
		log.Fatal("Error getting IP: ", err)
	}
	fmt.Println(ip)
}
