package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/simar7/gokv/types"

	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws/endpoints"

	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/simar7/gokv"
	"github.com/simar7/gokv-poc/benchmarks"
	"github.com/simar7/gokv/dynamodb"
)

type foo struct {
	Bar string
}

var customDynamoDBEndpoint = "http://localhost:8000"
var optionsDynamoDB = dynamodb.Options{
	Region:         endpoints.UsWest2RegionID,
	AWSAccessKeyID: "fookey",
	AWSSecretAccessKey: "barsecretkey",
	CustomEndpoint: customDynamoDBEndpoint,
	TableName: "gokv",
}

func main() {
	modeOp := flag.String("op", "demo", "mode of operation: default is demo mode")
	flag.Parse()

	switch *modeOp {
	case "demo":
		checkConnections()

		clientDynamoDB := setupDynamoClient()
		defer clientDynamoDB.Close()

		// TODO: create a function that sets up the table for interaction
		interactWithStore(clientDynamoDB)

	case "bench":
		log.Println("running boltdb benchmarks...")
		benchmarks.BoltUpdate()
		benchmarks.BoltBatch()
	default:
		log.Fatalf("invalid mode specified: %s", *modeOp)
	}
}

func checkConnections() {
	if !checkConnectionDynamoDB() {
		log.Fatal("couldn't establish connection with dynamodb")
	}
	log.Println("dynamodb session established.")
}

func setupDynamoClient() dynamodb.Store {
	log.Println("commencing ops with dynamodb...")
	clientDynamoDB, err := dynamodb.NewStore(optionsDynamoDB)
	if err != nil {
		log.Fatal(err)
	}
	return clientDynamoDB
}

// checkConnectionDynamoDB returns true if a connection could be made, false otherwise.
func checkConnectionDynamoDB() bool {
	sess, err := session.NewSession(aws.NewConfig().WithRegion(endpoints.UsWest2RegionID).WithEndpoint(customDynamoDBEndpoint))
	if err != nil {
		log.Printf("An error occurred during testing the connection to the server: %v\n", err)
		return false
	}
	svc := awsdynamodb.New(sess)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	limit := int64(1)
	listTablesInput := awsdynamodb.ListTablesInput{
		Limit: &limit,
	}
	_, err = svc.ListTablesWithContext(timeoutCtx, &listTablesInput)
	if err != nil {
		log.Printf("An error occurred during testing the connection to the server: %v\n", err)
		return false
	}
	return true
}

// interactWithStore stores, retrieves, prints and deletes a value.
// It's completely independent of the store implementation.
func interactWithStore(store gokv.Store) {
	// Store value
	val := foo{
		Bar: "bar",
	}
	val2 := foo{
		Bar: "baz",
	}

	err := store.BatchSet(types.BatchSetItemInput{
		Keys:   []string{"foo", "faz"},
		Values: []foo{val, val2},
	})
	if err != nil {
		panic(err)
	}

	// Retrieve value
	var retrievedVal1, retrievedVal2 foo
	found1, err := store.Get(types.GetItemInput{Key:"foo", Value: &retrievedVal1})
	if err != nil {
		panic(err)
	}
	if !found1 {
		panic("Value not found")
	}

	found2, err := store.Get(types.GetItemInput{Key:"faz", Value: &retrievedVal2})
	if err != nil {
		panic(err)
	}
	if !found2 {
		panic("Value not found")
	}

	fmt.Printf("foo: %+v\n", retrievedVal1)
	fmt.Printf("faz: %+v\n", retrievedVal2)

	// Delete value
	err = store.Delete(types.DeleteItemInput{Key:"foo"})
	if err != nil {
		panic(err)
	}

	// Delete value
	err = store.Delete(types.DeleteItemInput{Key:"faz"})
	if err != nil {
		panic(err)
	}
}
