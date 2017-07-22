# exchange_fetcher

`exchange_fetcher` is a SAMPLE application that requests stock exchanges results from Yahoo! finance API.
It allows two types of usage: directly, asking for specific stock symbols as command-line arguments; externally, via message queueing, which allows a client application running on `RabbitMQ` to retrieve the results.
The sole purpose of this repository is a study of implementations on several tasks using Go, like JSON parsing, MQ communication, HTTP request handling.

## To whom this might be interesting
To anyone who is learning Go, like me, is an enthusiast of Go, also like me...to anyone who might be interesting on getting stock results fast enough...
Finally, I think, to anyone who likes programming, since you can contribute to this repo at any time!

## Requirements
This application was developed with Go `1.7.3`.
To operate on message queueing, it's necessary to install <a href="https://www.rabbitmq.com/">RabbitMQ</a>.

## Installation
Get this repo via Go CLI. Run `go install`.

## Usage
`exchange_fetcher` requires symbols from stock exchanges, that are recognized by Yahoo! finance API.

(WIP) As an improvement, a list of common stock symbols will be inserted here.

The command-line can be used in two ways:

Directly:
```
$> exchange_fetcher --indices '{"indices":["GOOGL"]}'
// This will make an HTTP request directly to Yahoo! finance API and return a JSON response of the stocks sent.
```

Externally:
```
$> exchange_fetcher -mq
// This will start connection with RabbitMQ, using proper configuration environment variables.
```

<a href="https://www.rabbitmq.com/">RabbitMQ</a> connection on the application requires environment variables set on `.env` file. It is necessary to define a `.env` file, based on `.env.example` file present on root of this repo.

## Contributing
Feel free to open a pull request, point an issue. I am on the search of learning Go the best way possible, so every opinion and any line of code are welcome!
Fork this repository, make your changes and open a pull request.

## License
MIT License. Please, read LICENSE file.

### Godspeed!!
