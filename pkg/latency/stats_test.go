package latency

import (
	"testing"
	"time"
)

func TestSummarize(t *testing.T) {
	got, err := Summarize([]time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got.Count != 5 {
		t.Fatalf("count=%d", got.Count)
	}
	if got.Min != 10*time.Millisecond || got.Max != 50*time.Millisecond {
		t.Fatalf("min/max=%s/%s", got.Min, got.Max)
	}
	if got.P50 != 30*time.Millisecond {
		t.Fatalf("p50=%s", got.P50)
	}
}

func TestSummarizeRejectsEmptySamples(t *testing.T) {
	if _, err := Summarize(nil); err == nil {
		t.Fatal("expected error")
	}
}
