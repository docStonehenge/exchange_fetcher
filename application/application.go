package application

import (
	"flag"
	"fmt"
	"github.com/docStonehenge/exchange_fetcher/connector"
	"github.com/docStonehenge/exchange_fetcher/exchange"
	"github.com/docStonehenge/exchange_fetcher/indices"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var symbols string
var onQueue bool

func Run() {
	parseCommandFlags()

	if onQueue {
		runProcessOnMQ()
	} else {
		logIndicesRequest()
	}
}

func parseCommandFlags() {
	flag.BoolVar(
		&onQueue, "mq", false,
		"Set to true if you want to run on RabbitMQ connection with a client application",
	)

	flag.StringVar(
		&symbols, "indices", "{}",
		"JSON representation with indices key pointing to an array of symbols",
	)

	flag.Parse()
}

func runProcessOnMQ() {
	loadEnvironment()

	fmt.Println("Connecting to AMQP server...")
	connection, err := connector.OpenConnection()
	defer connection.Close()
	logFailureAndCrash(err)
	fmt.Printf("Connected successfully to port %s\n", os.Getenv("AMQP_DEFAULT_PORT"))

	channel, err := connector.OpenChannel(connection)
	defer channel.Close()
	logFailureAndCrash(err)
	fmt.Println("Channel is now opened...")

	subQueueName := "exchange_fetcher.indices.requests"
	queueForSubscription, err := connector.DefineQueue(channel, subQueueName)
	logFailureAndCrash(err)
	fmt.Printf("Receiving indices on queue '%s'\n", queueForSubscription.Name)

	pubQueueName := "exchange_fetcher.indices.results"
	queueForPublishing, err := connector.DefineQueue(channel, pubQueueName)
	logFailureAndCrash(err)
	fmt.Printf("Publishing results on queue '%s'\n", queueForPublishing.Name)

	subscriber, err := connector.OpenSubscriber(channel, queueForSubscription.Name)
	logFailureAndCrash(err)

	fmt.Printf("\n\nWaiting for indices. Press Crtl+C to exit.\n\n")

	openConnectionToApplication := make(chan bool)
	indicesReceived := make(chan []string)

	go connector.HandleReceivedIndices(subscriber, indicesReceived)

	for indices := range indicesReceived {
		result := requestIndices(indices)
		err = connector.PublishIndices(channel, queueForPublishing.Name, result)
		logOperationResult(err, "Published results to subscribers.")
	}

	<-openConnectionToApplication
}

func logIndicesRequest() {
	result, _ := indices.Join(
		requestIndices(indices.Split([]byte(symbols))).Exchanges,
	)

	logOperationResult(nil, fmt.Sprintf("%s", result))
}

func requestIndices(indices []string) *exchange.ExchangesResult {
	fmt.Printf("Indices received are: %v\n", indices)

	url := exchange.BuildURL(indices)

	result, err := exchange.Fetch(url)
	logOperationResult(
		err, "Setting connection and fetching results from Yahoo! API...",
	)

	err = result.Parse()
	logOperationResult(
		err, "Successfully received results from Yahoo! API.",
	)

	return result
}

func logOperationResult(err error, message string) {
	if err == nil {
		log.Println(message)
	} else {
		log.Println(err)
	}
}

func loadEnvironment() {
	if envError := godotenv.Load(); envError != nil {
		log.Fatalf("An error occurred while loading environment: %v", envError)
	}
}

func logFailureAndCrash(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
