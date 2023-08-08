package tools

import (
    "github.com/prometheus/client_golang/prometheus"
)

var SuccessfullyThrownCards = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "successfully_thrown_cards",
        Help: "The number of cards thrown on the table",
    },
)
var UnsuccessfullyThrownCards = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "unsuccessfully_thrown_cards",
        Help: "The number of unsuccessful attempts to throw a card",
    },
)

var SuccessfullyTakenCards = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "successfully_taken_cards",
        Help: "The number of cards taken from the table",
    },
)
var UnsuccessfullyTakenCards = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "unsuccessfully_taken_cards",
        Help: "The number of unsuccessful attempts to take a cards",
    },
)

var SuccessfullyStartedGame = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "successfully_started_game",
        Help: "The number of successfully started games",
    },
)
var UnsuccessfullyStartedGame = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "unsuccessfully_started_game",
        Help: "The number of unsuccessful attempts to start games",
    },
)

var SuccessfullyShowedCards = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "successfully_showed_cards",
        Help: "The number of successful show card requests",
    },
)
var UnsuccessfullyShowedCards = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "unsuccessfully_showed_cards",
        Help: "The number of unsuccessful show card requests",
    },
)

var ParsingErrorCounter = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "parsing_response_error_count",
        Help: "The number of parsing errors",
    },
)
var DatabaseErrorCounter = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "database_connection_error_count",
        Help: "The number of database connection errors",
    },
)

var ServiceStatus = prometheus.NewGauge(
    prometheus.GaugeOpts{
        Name: "service_status",
        Help: "If the value is 1, the service is active, if the value is 0, the service is inactive",
    },
)
