package main

import (
	"fmt"
	"github.com/docStonehenge/exchange_fetcher/connector"
	"github.com/docStonehenge/exchange_fetcher/exchange"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	loadEnvironmentOrCrash()

	fmt.Println("Connecting to AMQP server...")
	connection, err := connector.OpenConnection()
	defer connection.Close()
	logFailureAndCrash(err)
	fmt.Printf("Connected successfully to port %s\n", os.Getenv("AMQP_DEFAULT_PORT"))

	channel, err := connector.OpenChannel(connection)
	defer channel.Close()
	logFailureAndCrash(err)
	fmt.Println("Channel is now opened...")

	subQueueName := "exchange_fetcher.indices.subscription"
	queueForSubscription, err := connector.DefineQueue(channel, subQueueName)
	logFailureAndCrash(err)
	fmt.Printf("Receiving indices on queue '%s'\n", queueForSubscription.Name)

	pubQueueName := "exchange_fetcher.indices.publishing"
	queueForPublishing, err := connector.DefineQueue(channel, pubQueueName)
	logFailureAndCrash(err)
	fmt.Printf("Publishing results on queue '%s'\n", queueForPublishing.Name)

	subscriber, err := connector.OpenSubscriber(channel, queueForSubscription.Name)
	logFailureAndCrash(err)

	fmt.Printf("\n\nWaiting for indices. Press Crtl+C to exit.\n\n")

	openConnectionToApplication := make(chan bool)

	go connector.HandleReceivedIndices(
		subscriber, func(indices []string) {
			fmt.Printf("Indices received are: %v\n", indices)

			url := exchange.BuildURL(indices)

			result, err := exchange.Fetch(url)
			log.Println("Setting connection and fetching results from Yahoo! API...")

			if err != nil {
				log.Println(err)
			}

			err = result.Parse()

			if err != nil {
				log.Println(err)
			}

			log.Println("Successfully received results from Yahoo! API.")
			err = connector.PublishIndices(channel, queueForPublishing.Name, result)

			if err != nil {
				log.Printf("There was a problem on setting publisher: %v", err)
			}
		},
	)

	<-openConnectionToApplication
}

func loadEnvironmentOrCrash() {
	if envError := godotenv.Load(); envError != nil {
		log.Fatalf("An error occurred while loading environment: %v", envError)
	}
}

func logFailureAndCrash(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
