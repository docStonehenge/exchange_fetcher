package connector

import (
	"fmt"
	"github.com/docStonehenge/exchange_fetcher/exchange"
	"github.com/docStonehenge/exchange_fetcher/indices"
	"github.com/streadway/amqp"
	"os"
)

type ConnectionError struct {
	err error
}

type ChannelError struct {
	err error
}

func OpenConnection() (*amqp.Connection, error) {
	connection, err := amqp.Dial(formatAmqpURL())

	if err == nil {
		return connection, nil
	}

	return nil, &ConnectionError{err: err}
}

func OpenChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	channel, err := connection.Channel()

	if err == nil {
		return channel, nil
	}

	return nil, &ChannelError{err: err}
}

func DefineQueue(channel *amqp.Channel, queueName string) (amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		queueName,
		false,
		true,
		false,
		false,
		nil,
	)

	return queue, err
}

func OpenSubscriber(channel *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	subscriber, err := channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err == nil {
		return subscriber, nil
	}

	return nil, err
}

func HandleReceivedIndices(subscriber <-chan amqp.Delivery, indicesChannel chan []string) {
	for delivery := range subscriber {
		indicesChannel <- indices.Split(delivery.Body)
	}
}

func PublishIndices(channel *amqp.Channel, queueName string, result *exchange.ExchangesResult) error {
	response, err := indices.Join(result.Exchanges)

	if err != nil {
		return err
	}

	if publishingError := channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        response,
		},
	); publishingError != nil {
		return publishingError
	}

	return nil
}

func (connError *ConnectionError) Error() string {
	return fmt.Sprintf("There was a problem when opening connection to AMQP: %v", connError.err)
}

func (channelError *ChannelError) Error() string {
	return fmt.Sprintf(
		"There was a problem when opening channel on connection: %v",
		channelError.err,
	)
}

func formatAmqpURL() string {
	return fmt.Sprintf(
		"amqp://%s:%s@localhost:%s/",
		os.Getenv("AMQP_USERNAME"),
		os.Getenv("AMQP_PASSWORD"),
		os.Getenv("AMQP_DEFAULT_PORT"),
	)
}
