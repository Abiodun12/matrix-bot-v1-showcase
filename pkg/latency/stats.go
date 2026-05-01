package latency

import (
	"errors"
	"sort"
	"time"
)

// Summary is a compact latency distribution suitable for README tables,
// dashboards, and post-run scorecards.
type Summary struct {
	Count int
	Min   time.Duration
	P50   time.Duration
	P90   time.Duration
	P95   time.Duration
	P99   time.Duration
	Max   time.Duration
	Mean  time.Duration
}

// Summarize returns common latency percentiles for a slice of samples.
func Summarize(samples []time.Duration) (Summary, error) {
	if len(samples) == 0 {
		return Summary{}, errors.New("empty latency sample")
	}

	vals := append([]time.Duration(nil), samples...)
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })

	var total time.Duration
	for _, v := range vals {
		total += v
	}

	return Summary{
		Count: len(vals),
		Min:   vals[0],
		P50:   percentile(vals, 50),
		P90:   percentile(vals, 90),
		P95:   percentile(vals, 95),
		P99:   percentile(vals, 99),
		Max:   vals[len(vals)-1],
		Mean:  total / time.Duration(len(vals)),
	}, nil
}

func percentile(vals []time.Duration, pct int) time.Duration {
	if len(vals) == 1 {
		return vals[0]
	}

	rank := float64(len(vals)-1) * float64(pct) / 100.0
	low := int(rank)
	high := low + 1
	if high >= len(vals) {
		return vals[len(vals)-1]
	}

	frac := rank - float64(low)
	lo := float64(vals[low])
	hi := float64(vals[high])
	return time.Duration(lo + (hi-lo)*frac)
}
