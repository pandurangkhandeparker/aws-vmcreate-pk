package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var client *ec2.Client

// EC2CreateInstanceAPI defines the interface for the RunInstances and CreateTags functions.
// We use this interface to test the functions using a mocked service.
type EC2CreateInstanceAPI interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	CreateTags(ctx context.Context,
		params *ec2.CreateTagsInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
}

// EC2TerminateInstanceAPI defines the interface for the TerminateInstances functions.
// We use this interface to test the functions using a mocked service.
type EC2TerminateInstanceAPI interface {
	TerminateInstances(ctx context.Context,
		params *ec2.TerminateInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

// MakeInstance creates an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a RunInstancesOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to RunInstances.
func MakeInstance(c context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}

// MakeTags creates tags for an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a CreateTagsOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to CreateTags.
func MakeTags(c context.Context, api EC2CreateInstanceAPI, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return api.CreateTags(c, input)
}

// DeleteInstance terminates an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a TerminateInstancesOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to TerminateInstances.
func DeleteInstance(c context.Context, api EC2TerminateInstanceAPI, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return api.TerminateInstances(c, input)
}

// func DeleteInstancesCmd(instanceIds []string) {
// 	// fmt.Println("TO DO")

// 	input := &ec2.TerminateInstancesInput{
// 		InstanceIds: instanceIds,
// 		DryRun:      new(bool),
// 	}

// 	result, err := DeleteInstance(context.TODO(), client, input)
// 	if err != nil {
// 		fmt.Println("Got an error terminating the instance:")
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println("Terminated instance with id: ", *result.TerminatingInstances[0].InstanceId)
// }

func DeleteInstancesCmd(name *string, value *string) {
	// fmt.Println("TO DO")

	// tagInput := &ec2.DeleteTagsInput{
	// 	Resources: []string{*result.Instances[0].InstanceId},
	// 	Tags: []types.Tag{
	// 		{
	// 			Key:   name,
	// 			Value: value,
	// 		},
	// 	},
	// }

	// _, err = RemoveTags(context.TODO(), client, tagInput)
	// if err != nil {
	// 	fmt.Println("Got an error tagging the instance:")
	// 	fmt.Println(err)
	// 	return
	// }

	var instanceIds = make([]string, 0)

	val := strings.Split(*value, ",")
	tag := "tag:" + *name

	input1 := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(tag),
				Values: val,
			},
		},
	}
	result, err := client.DescribeInstances(context.TODO(), input1)
	if err != nil {
		fmt.Println("Got an error fetching the status of the instance")
		fmt.Println(err)
	} else {
		for _, r := range result.Reservations {
			fmt.Println("Instance IDs:")
			for _, i := range r.Instances {
				//value := *i.InstanceId
				instanceIds = append(instanceIds, *i.InstanceId)
			}
			fmt.Println(instanceIds)
		}

		input := &ec2.TerminateInstancesInput{
			InstanceIds: instanceIds,
			DryRun:      new(bool),
		}

		result, err := DeleteInstance(context.TODO(), client, input)
		if err != nil {
			fmt.Println("Got an error terminating the instance:")
			fmt.Println(err)
			return
		}

		fmt.Println("Terminated instance with id: ", *result.TerminatingInstances[0].InstanceId)
	}
}

func CreateInstancesCmd(name *string, value *string, imageId *string, instanceType *string) {
	// Create separate values if required.
	minMaxCount := int32(1)

	// insType := (types.InstanceType)(*instanceType)

	input := &ec2.RunInstancesInput{
		// ImageId:      aws.String("ami-0d0ca2066b861631c"),
		// InstanceType: types.InstanceTypeT2Micro, OR "t2.micro"
		ImageId:      imageId,
		InstanceType: (types.InstanceType)(*instanceType),
		MinCount:     &minMaxCount,
		MaxCount:     &minMaxCount,
	}

	result, err := MakeInstance(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error creating an instance:")
		fmt.Println(err)
		return
	}

	tagInput := &ec2.CreateTagsInput{
		Resources: []string{*result.Instances[0].InstanceId},
		Tags: []types.Tag{
			{
				Key:   name,
				Value: value,
			},
			{
				Key:   aws.String("Name"),
				Value: aws.String("pandurang-ec2"),
			},
		},
	}

	_, err = MakeTags(context.TODO(), client, tagInput)
	if err != nil {
		fmt.Println("Got an error tagging the instance:")
		fmt.Println(err)
		return
	}

	fmt.Println("Created tagged instance with ID " + *result.Instances[0].InstanceId)
}
func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client = ec2.NewFromConfig(cfg)

}
func main() {
	fmt.Println("Provisioning/De-provisioning EC2 in progress")
	command := flag.String("c", "", "command  create or delete")
	name := flag.String("n", "", "The name of the tag to attach to the instance")
	value := flag.String("v", "", "The value of the tag to attach to the instance")
	imageId := flag.String("i", "", "The image id of the instance")
	instanceType := flag.String("t", "", "The instance type of the new instance")
	// instanceId := flag.String("i", "", "The IDs of the instance to terminate")

	flag.Parse()

	if *command == "" {
		fmt.Println("You must supply an command  start or stop (-c start)")
		return
	}

	if *name == "" || *value == "" || *imageId == "" || *instanceType == "" {
		fmt.Println("You must supply a name and value for the tag (-n NAME -v VALUE -i IMAGEID -t INSTANCETYPE)")
		return
	}

	// if *instanceId == "" {
	// 	fmt.Println("You must supply an instance ID (-i INSTANCE-ID or comma separated list of ids")
	// 	return
	// }

	// instances := strings.Split(*instanceId, ",")

	if *command == "create" {
		CreateInstancesCmd(name, value, imageId, instanceType)
	}

	if *command == "delete" {
		DeleteInstancesCmd(name, value)
	}
}
