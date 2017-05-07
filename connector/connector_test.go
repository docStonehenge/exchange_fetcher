package connector

import (
	"fmt"
	"github.com/docStonehenge/exchange_fetcher/exchange"
	"github.com/streadway/amqp"
	"os"
	"regexp"
	"strings"
	"testing"
)

var connErrMatch = regexp.MustCompile("There was a problem when opening connection to AMQP")
var channelErrMatch = regexp.MustCompile("There was a problem when opening channel on connection")

func TestOpenConnectionWithSuccess(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		defer connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	if err != nil {
		t.Fatalf("OpenConnection() should return correct connection to amqp, but returned error: %v", err)
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestOpenConnectionWithError(t *testing.T) {
	_, err := OpenConnection()

	if err == nil {
		t.Fatalf("OpenConnection() should return an error, but error is null: %v", err)
	}

	msg := fmt.Sprintf("%v", err)

	if !connErrMatch.MatchString(msg) {
		t.Fatalf("Error message should be for connection error, but is %s", msg)
	}
}

func TestConnectionError(t *testing.T) {
	_, err := OpenConnection()

	connError := ConnectionError{err: err}
	msg := connError.Error()

	if !connErrMatch.MatchString(msg) {
		t.Fatalf("Error message should have prefix of 'There was a problem when opening connection to AMQP', but message is: %s", msg)
	}
}

func TestOpenChannelWithSuccess(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		defer connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	channel, err := OpenChannel(connection)
	defer channel.Close()

	if err != nil {
		t.Fatalf("OpenChannel() should return an opened channel, but returned error: %v", err)
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestOpenChannelWithError(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	_, err = OpenChannel(connection)

	if err == nil {
		t.Fatalf("OpenChannel() should return an error, but returned nothing: %v", err)
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestChannelError(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	_, err = OpenChannel(connection)
	channelError := &ChannelError{err: err}

	msg := channelError.Error()

	if !channelErrMatch.MatchString(msg) {
		t.Fatalf("Channel error message should start with 'There was a problem when opening channel on connection', but is %s", msg)
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestDefineQueue(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		defer connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	channel, err := OpenChannel(connection)

	if channel != nil {
		defer channel.Close()
	} else {
		t.Fatalf("There was a problem on opening channel: %v", err)
	}

	queue, err := DefineQueue(channel, "test_queue1")

	if err != nil {
		t.Fatalf("DefineQueue() should return a queue, but returned error: %v", err)
	}

	if queue.Name != "test_queue1" {
		t.Fatalf("Queue name should be %s, but is %s", "test_queue1", queue.Name)
	}

	if _, err := channel.QueueDelete("test_queue1", false, false, false); err != nil {
		t.Fatal("Queue 'test_queue1' should be deleted.")
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

// These are integration tests that must be run only on controlled environment, since amqp connection cannot be closed from within the test.
// To run these tests:
//   * uncomment the imports
//   * open RabbitMQ dashboard at http://localhost:15672 and force-close connection after test run, to exit it.

func TestOpenSubscriber(t *testing.T) {
	integrationEnvironmentForTest(
		t,
		func(channel *amqp.Channel, queueName string) {
			testBody := "{\"indices\": [\"^BVSP\", \"AAPL\"]}"
			channel.Publish(
				"",
				queueName,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(testBody),
				},
			)

			_, err := OpenSubscriber(channel, queueName)

			if err != nil {
				t.Fatalf("OpenSubscriber() should return a channel for messages, but returned error: %v", err)
			}
		},
	)
}

func TestOpenSubscriberRaisesErrorOnOpening(t *testing.T) {
	integrationEnvironmentForTest(
		t,
		func(channel *amqp.Channel, queueName string) {
			testBody := "{\"indices\": [\"^BVSP\", \"AAPL\"]}"
			channel.Publish(
				"",
				queueName,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(testBody),
				},
			)

			_, err := OpenSubscriber(channel, "foo")

			if err == nil {
				t.Fatal("OpenSubscriber() should raise an error, but nothing was raised")
			}
		},
	)
}

func TestHandleReceivedIndicesPutsIndicesOnChannel(t *testing.T) {
	integrationEnvironmentForTest(
		t,
		func(channel *amqp.Channel, queueName string) {
			testBody := "{\"indices\": [\"^BVSP\", \"AAPL\"]}"
			channel.Publish(
				"",
				queueName,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(testBody),
				},
			)

			subscriber, err := channel.Consume(
				queueName, "", true, false, false, false, nil,
			)

			if err != nil {
				t.Fatal()
			}

			indicesChannel := make(chan []string)

			go HandleReceivedIndices(subscriber, indicesChannel)
			receivedIndices := <-indicesChannel

			if strings.Join(receivedIndices, ",") != "^BVSP,AAPL" {
				t.Fatal("Should handle subscriber received indices with handler, but nothing happened.")
			}
		},
	)
}

func TestHandleReceivedIndicesSkipsEmptyJSON(t *testing.T) {
	integrationEnvironmentForTest(
		t,
		func(channel *amqp.Channel, queueName string) {
			testBody := "{}"
			channel.Publish(
				"",
				queueName,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(testBody),
				},
			)

			subscriber, err := channel.Consume(
				queueName, "", true, false, false, false, nil,
			)

			if err != nil {
				t.Fatal()
			}

			indicesChannel := make(chan []string)

			go HandleReceivedIndices(subscriber, indicesChannel)
			receivedIndices := <-indicesChannel

			if strings.Join(receivedIndices, ",") != "" {
				t.Fatal("Should return an empty collection without raising error.")
			}
		},
	)
}

func TestHandleReceivedIndicesSkipsEmptyStringBodyReceived(t *testing.T) {
	integrationEnvironmentForTest(
		t,
		func(channel *amqp.Channel, queueName string) {
			testBody := ""
			channel.Publish(
				"",
				queueName,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(testBody),
				},
			)

			subscriber, err := channel.Consume(
				queueName, "", true, false, false, false, nil,
			)

			if err != nil {
				t.Fatal()
			}

			indicesChannel := make(chan []string)

			go HandleReceivedIndices(subscriber, indicesChannel)
			receivedIndices := <-indicesChannel

			if strings.Join(receivedIndices, ",") != "" {
				t.Fatal("Should return an empty collection without raising error.")
			}
		},
	)
}

func TestPublishIndices(t *testing.T) {
	integrationEnvironmentForTest(
		t,
		func(channel *amqp.Channel, queueName string) {
			result := &exchange.ExchangesResult{
				Exchanges: map[string]exchange.Exchange{
					"Nikkei 225":    exchange.Exchange{Name: "Nikkei 225", Symbol: "^n225", PercentChange: "-0.91%", ChangeInPoints: "-172.98", LastTradeDate: "4/14/2017", LastTradeTime: "3:15pm"},
					"Alphabet Inc.": exchange.Exchange{Name: "Alphabet Inc.", Symbol: "GOOGL", PercentChange: "-0.09%", ChangeInPoints: "-0.76", LastTradeDate: "4/13/2017", LastTradeTime: "4:00pm"},
				},
			}

			PublishIndices(channel, queueName, result)

			msgs, _ := channel.Consume(
				queueName, "", true, false, false, false, nil,
			)

			expected := "{\"Alphabet Inc.\":{\"Name\":\"Alphabet Inc.\",\"Symbol\":\"GOOGL\",\"PercentChange\":\"-0.09%\",\"ChangeInPoints\":\"-0.76\",\"LastTradeDate\":\"4/13/2017\",\"LastTradeTime\":\"4:00pm\"},\"Nikkei 225\":{\"Name\":\"Nikkei 225\",\"Symbol\":\"^n225\",\"PercentChange\":\"-0.91%\",\"ChangeInPoints\":\"-172.98\",\"LastTradeDate\":\"4/14/2017\",\"LastTradeTime\":\"3:15pm\"}}"

			for msg := range msgs {
				parsedBody := string(msg.Body)

				if parsedBody != expected {
					t.Fatalf("Published body should be %s, but it is %s", expected, parsedBody)
				}
			}
		},
	)
}

func TestPublishIndicesRaisesErrorOnPublishProblem(t *testing.T) {
	integrationEnvironmentForTest(
		t,
		func(channel *amqp.Channel, queueName string) {
			result := &exchange.ExchangesResult{
				Exchanges: map[string]exchange.Exchange{
					"Nikkei 225":    exchange.Exchange{Name: "Nikkei 225", Symbol: "^n225", PercentChange: "-0.91%", ChangeInPoints: "-172.98", LastTradeDate: "4/14/2017", LastTradeTime: "3:15pm"},
					"Alphabet Inc.": exchange.Exchange{Name: "Alphabet Inc.", Symbol: "GOOGL", PercentChange: "-0.09%", ChangeInPoints: "-0.76", LastTradeDate: "4/13/2017", LastTradeTime: "4:00pm"},
				},
			}

			channel.Close()
			publishError := PublishIndices(channel, queueName, result)

			if publishError == nil {
				t.Fatal("PublishIndices should run with errors, but no error was raised")
			}
		},
	)
}

func integrationEnvironmentForTest(t *testing.T, handler func(channel *amqp.Channel, queueName string)) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()
	defer connection.Close()

	if err != nil {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	channel, err := OpenChannel(connection)
	defer channel.Close()

	if err != nil {
		t.Fatalf("There was a problem on opening channel: %v", err)
	}

	queue, err := DefineQueue(channel, "test_queue6")

	if err != nil {
		t.Fatalf("There was a problem on defining queue: %v", err)
	}

	handler(channel, queue.Name)

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}
