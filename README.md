# exchange_fetcher

`exchange_fetcher` is a SAMPLE application that requests stock exchanges results from Yahoo! finance API.
It allows two types of usage: directly, asking for specific stock symbols as command-line arguments; externally, via message queueing, which allows a client application running on <a href="https://www.rabbitmq.com/">RabbitMQ</a> to retrieve the results.
The sole purpose of this repository is a study of implementations on several tasks using Go, like JSON parsing, MQ communication, HTTP request handling.

## To whom this might be interesting
To anyone who is learning Go, like me, is an enthusiast of Go, also like me...to anyone who might be interesting on getting stock results fast enough...
Finally, I think, to anyone who likes programming, since you can contribute to this repo at any time!

## Requirements
This application was initially developed using Go 1.7.3.
For contributing, because of packages used on dependency management (for now), it requires at least 1.8.x versions.
To operate on message queueing, it's necessary to install <a href="https://www.rabbitmq.com/">RabbitMQ</a>.

## Installation
Get this repo via Go CLI. Run `go install`.

## Usage
`exchange_fetcher` requires symbols from stock exchanges, that are recognized by Yahoo! finance API.

(WIP) As an improvement, a list of common stock symbols will be inserted here.

The command-line can be used in two ways:

Directly:
```
$> exchange_fetcher --indices AAPL
// This will make an HTTP request directly to Yahoo! finance API and return a JSON response for AAPL stock result.
```
```
$> exchange_fetcher --indices 'AAPL, GOOGL'
// You can send a comma-separated list of symbols for multiple stock results !!
```

Externally:
```
$> exchange_fetcher -mq
// This will start connection with RabbitMQ, using proper configuration environment variables.
```

### External usage provides two queues:
`exchange_fetcher.indices.requests`, which a client can send requests with content-type of "application/json". JSON requests must have an `indices` key to an array of stock symbols as strings.
`exchage_fetcher.indices.results`, where the application will post all results parsed as a JSON representation. Example:

```
{"indices":["AAPL"]}
//JSON request sent via RabbitMQ driver; content-type application/json. It is possible to send multiple symbols in the indices array, on the same request.
```

```
{"Apple Inc.":{"Name":"Apple Inc.","Symbol":"AAPL","Price":155.39,"PreviousClose":153.48,"OpenPrice":154.18,"PercentChange":"+1.24%","ChangeInPoints":"+1.91","LastTradeDate":"10/5/2017","LastTradeTime":"4:00pm"}}
// result as a JSON representation
```

`exchange_fetcher` logs every process since the connection to MQ. At each request, the application displays which indices (symbols) were received and also the status of request/response for the stocks.

<a href="https://www.rabbitmq.com/">RabbitMQ</a> connection on the application requires environment variables set on `.env` file. It is necessary to define a `.env` file, based on `.env.example` file present on root of this repo.

## Contributing
Feel free to open a pull request, point an issue. I am on the search of learning Go the best way possible, so every opinion and any line of code are welcome!

  * First of all, fork this repository.
  * Run `go get github.com/docStonehenge/exchange_fetcher`; then
  * `cd $GOPATH/src/github.com/docStonehenge/exchange_fetcher`; then
  * `git remote set-url origin your-fork-url`; then
  * `go get -u github.com/golang/dep/cmd/dep`; then
  * `dep ensure`...whoa, then finally
  * Make your changes, following the lovely TDD approach :-P
## License
MIT License. Please, read LICENSE file.

### Godspeed!!
