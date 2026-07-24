package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var EventsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "swarmlens_events_total",
		Help: "Total number of agent events processed",
	},
	[]string{"swarm_id", "agent_id", "event_type"},
)
