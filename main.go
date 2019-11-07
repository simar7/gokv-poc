package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	goredis "github.com/go-redis/redis"

	"github.com/aws/aws-sdk-go/aws/endpoints"

	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/dynamodb"
	"github.com/philippgille/gokv/redis"
	"github.com/simar7/gokv-poc/benchmarks"
)

type foo struct {
	Bar string
}

var customDynamoDBEndpoint = "http://localhost:8000"
var testRedisDbNumber = 15
var optionsDynamoDB = dynamodb.Options{
	Region:         endpoints.UsWest2RegionID,
	CustomEndpoint: customDynamoDBEndpoint,
}

func main() {
	modeOp := flag.String("op", "demo", "mode of operation: default is demo mode")
	flag.Parse()

	switch *modeOp {
	case "demo":
		checkConnections()

		clientDynamoDB := setupDynamoClient()
		defer clientDynamoDB.Close()
		interactWithStore(clientDynamoDB)

		clientRedis := setupRedisClient()
		defer clientRedis.Close()
		interactWithStore(clientRedis)
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

	if !checkConnectionRedis(testRedisDbNumber) {
		log.Fatal("couldn't establish connection with redis")
	}
	log.Println("redis session established.")
}

func setupDynamoClient() dynamodb.Client {
	log.Println("commencing ops with dynamodb...")
	clientDynamoDB, err := dynamodb.NewClient(optionsDynamoDB)
	if err != nil {
		log.Fatal(err)
	}
	return clientDynamoDB
}

func setupRedisClient() redis.Client {
	log.Println("commencing ops with redis...")
	optionsRedis := redis.Options{
		DB: testRedisDbNumber,
	}
	clientRedis, err := redis.NewClient(optionsRedis)
	if err != nil {
		log.Fatal(err)
	}
	return clientRedis
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

// checkConnection returns true if a connection could be made, false otherwise.
func checkConnectionRedis(number int) bool {
	client := goredis.NewClient(&goredis.Options{
		Addr:     redis.DefaultOptions.Address,
		Password: redis.DefaultOptions.Password,
		DB:       number,
	})
	defer client.Close()
	err := client.Ping().Err()
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
		Bar: "baz",
	}
	err := store.Set("foo123", val)
	if err != nil {
		panic(err)
	}

	// Retrieve value
	retrievedVal := new(foo)
	found, err := store.Get("foo123", retrievedVal)
	if err != nil {
		panic(err)
	}
	if !found {
		panic("Value not found")
	}

	fmt.Printf("foo: %+v\n", *retrievedVal) // Prints `foo: {Bar:baz}`

	// Delete value
	err = store.Delete("foo123")
	if err != nil {
		panic(err)
	}
}
