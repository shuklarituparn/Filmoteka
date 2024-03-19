package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"testing"
)

func TestMetrics(t *testing.T) {

	Init()
	CreateMovieApiPingCounter.Inc()
	CreateActorApiPingCounter.Inc()
	prometheus.Unregister(CreateMovieApiPingCounter)
	prometheus.Unregister(CreateActorApiPingCounter)
	prometheus.Unregister(RegisterApiPingCounter)
	prometheus.Unregister(LoginApiPingCounter)
	prometheus.Unregister(UpdateMovieApiPingCounter)
	prometheus.Unregister(UpdateActorApiPingCounter)
	prometheus.Unregister(PatchActorApiPingCounter)
	prometheus.Unregister(PatchMovieApiPingCounter)
	prometheus.Unregister(DeleteMovieApiPingCounter)
	prometheus.Unregister(DeleteActorApiPingCounter)
	prometheus.Unregister(ReadAllActorApiPingCounter)
	prometheus.Unregister(ReadOneActorApiPingCounter)
	prometheus.Unregister(ReadOneMovieApiPingCounter)
	prometheus.Unregister(ReadAllMovieApiPingCounter)

	// Simulate some API pings for testing purposes

	// Add more increments for other metrics as needed

	// Test if the metrics have been incremented
	// For example, testing CreateMovieApiPingCounter
	if got := testutil.ToFloat64(CreateMovieApiPingCounter); got != 1 {
		t.Errorf("CreateMovieApiPingCounter was expected to be 1, got %v", got)
	}
}
