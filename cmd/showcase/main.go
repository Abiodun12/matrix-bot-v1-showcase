package main

import (
	"fmt"
	"time"

	"github.com/Abiodun12/matrix-bot-v1-showcase/pkg/latency"
	"github.com/Abiodun12/matrix-bot-v1-showcase/pkg/oms"
)

func main() {
	actions := oms.Reconcile(
		[]oms.DesiredOrder{
			{Key: "btc-5m:no:buy", Side: "buy", Price: 716, Size: 500},
			{Key: "btc-5m:yes:sell", Side: "sell", Price: 284, Size: 500, Reduce: true},
		},
		[]oms.OpenOrder{
			{ID: "ord-1", Key: "btc-5m:no:buy", Side: "buy", Price: 715, Size: 500},
			{ID: "ord-2", Key: "old-market:yes:buy", Side: "buy", Price: 500, Size: 100},
		},
	)

	fmt.Println("OMS reconcile plan")
	for _, action := range actions {
		fmt.Printf("- %-6s reason=%s order=%s desired=%s\n", action.Type, action.Reason, action.OrderID, action.Desired.Key)
	}

	summary, _ := latency.Summarize([]time.Duration{
		17 * time.Millisecond,
		28 * time.Millisecond,
		31 * time.Millisecond,
		41 * time.Millisecond,
		46 * time.Millisecond,
	})
	fmt.Printf("\nLatency sample: count=%d p50=%s p95=%s max=%s\n", summary.Count, summary.P50, summary.P95, summary.Max)
}
