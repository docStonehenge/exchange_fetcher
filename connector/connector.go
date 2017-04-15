package connector

import (
	"encoding/json"
	"fmt"
	"github.com/docStonehenge/exchange_fetcher/exchange"
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

func ReceiveIndices(channel *amqp.Channel, queueName string) ([]string, error) {
	indices := []string{}

	msgs, err := channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for message := range msgs {
		var jsonIndices map[string]interface{}
		json.Unmarshal(message.Body, &jsonIndices)

		for _, idx := range jsonIndices["indices"].([]interface{}) {
			if i, ok := idx.(string); ok {
				indices = append(indices, i)
			}
		}
	}

	return indices, err
}

func PublishIndices(channel *amqp.Channel, queueName string, result exchange.ExchangesResult) error {
	response, err := json.Marshal(result.Exchanges)

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
