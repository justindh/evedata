package apicache

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Custom transport to chain into the HTTPClient to gather statistics.
type transport struct {
	next *http.Transport
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to provide stats and handle error rates.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Loop until success
	tries := 0
	for {
		esiRateLimiter := true

		// Time our response
		start := time.Now().Nanosecond()

		// Tickup retry counter
		tries++

		// Run the request.
		res, err := t.next.RoundTrip(req)

		// We got a response
		if res != nil {
			// Get the ESI error information
			resetS := res.Header.Get("x-esi-error-limit-reset")
			tokensS := res.Header.Get("x-esi-error-limit-remain")

			if res.StatusCode >= 300 {
				metricAPIErrors.Inc()
				log.Printf("St: %d Res: %s Tok: %s - %s\n", res.StatusCode, resetS, tokensS, req.URL)
			}

			// If we cannot decode this is likely from another source.
			reset, err := strconv.Atoi(resetS)
			if err != nil {
				esiRateLimiter = false
			}
			tokens, err := strconv.Atoi(tokensS)
			if err != nil {
				esiRateLimiter = false
			}

			duration := (time.Now().Nanosecond() - start) / 1000
			metricAPICalls.With(
				prometheus.Labels{"host": req.Host},
			).Observe(float64(duration))

			// Sleep to prevent hammering CCP ESI if there are excessive errors
			if esiRateLimiter {
				time.Sleep(time.Second * time.Duration(float64(reset*2)*(1-(float64(tokens)/100))))
			}

			if res.StatusCode == 420 {
				time.Sleep(time.Second * time.Duration(reset+rand.Intn(30)))
			}

			if res.StatusCode == 420 || res.StatusCode >= 500 || res.StatusCode == 0 {
				// break out after 10 tries
				if tries > 10 {
					return res, err
				}
				if !esiRateLimiter {
					time.Sleep(time.Second * time.Duration(tries))
				}
				continue
			}
		}
		return res, err
	}
}

var (
	metricAPICalls = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "api",
		Name:      "calls",
		Help:      "API call statistics.",
		Buckets:   prometheus.ExponentialBuckets(10, 1.45, 20),
	},
		[]string{"host"},
	)

	metricAPIErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "evedata",
		Subsystem: "api",
		Name:      "errors",
		Help:      "Count of API errors.",
	})
)

func init() {
	prometheus.MustRegister(
		metricAPICalls,
		metricAPIErrors,
	)
}
