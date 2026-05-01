# Matrix Bot v1

Live-tested CLOB execution infrastructure for prediction markets.

Matrix Bot v1 is a Go-based execution stack built around Polymarket-style central limit order book markets. It was designed to solve the hard part of automated prediction-market trading: ingesting fast market data, managing live orders safely, enforcing runtime risk controls, and producing auditable post-trade records.

This public repository is a sanitized project overview. The live trading implementation, deployment runbooks, strategy logic, and private operational details are not published here.

## Public Code Included

This repo includes a small venue-neutral code sample:

- `pkg/oms`: deterministic desired-vs-open order reconciliation
- `pkg/latency`: percentile summaries for runtime latency samples
- `cmd/showcase`: demo CLI showing the public components
- `docs/ARCHITECTURE.md`: high-level system diagram

Run it locally:

```bash
go test ./...
go run ./cmd/showcase
```

## What It Demonstrates

Most trading bots stop at "fetch book, place order." Matrix Bot v1 focuses on everything that has to work after orders are live:

- WebSocket-first market-data ingestion
- Authenticated order routing
- OMS cancel/replace reconciliation
- Explicit paper/live arming controls
- Cancel-only safety mode
- Max-fill and max-runtime risk caps
- Reduce-only unwind behavior
- Position reconciliation
- User-fill tracking
- Realized PnL accounting
- Final REST-based trade audit
- Replayable logs for research and postmortems

## Observed Runtime Metrics

Measured from runtime/live logs:

| Metric | Observed |
|---|---:|
| Quote diagnostic samples | 16k+ |
| Health samples | 800+ |
| Market WS freshness at health, p50 | ~46 ms |
| User WS freshness at health, p50 | ~35 ms |
| Cancel-all REST latency, p50 | ~31 ms |
| Cancel-all REST latency, p95 | ~46 ms |
| Completed execution-session log entries | 100+ |
| Final audit records | 180+ |
| Final audit drift on clean sessions | 0.000000 |

These are operational measurements from the system's logs. They are not a claim of sub-millisecond end-to-end order execution latency.

## System Areas

### Market Data

- Market WebSocket ingestion
- Book and price-change handling
- Feed freshness checks
- Reconnect tracking
- Deduplication
- Raw capture for replay and analysis

### Execution

- Signed order construction
- Authenticated REST order flow
- Open-order reconciliation
- Cancel/replace loop
- Duplicate-order detection
- Post rejection handling
- Cancel-all safety path

### Risk Controls

- Explicit live arming
- Dry-run and cancel-only modes
- Max fill caps
- Max runtime caps
- No-fill session timeout
- No-new-risk mode after fills
- Reduce-only unwind behavior
- Dust/minimum-order handling
- Shutdown flatten window

### Accounting And Audit

- User WebSocket fill tracking
- Position truth sync
- Realized PnL ledger
- REST-based final audit
- Drift checks between runtime ledger and venue activity
- Session summaries for post-trade review

### Research Tooling

- Markout analysis
- Shadow mode
- Short-duration observer
- Queue-priority probes
- Evidence scorecards
- Replayable JSONL captures

## Why It Exists

Prediction-market execution is fragile. A bot needs to know when market data is stale, when an order is no longer desired, when inventory has moved, when to stop taking new risk, and whether final venue activity matches the local ledger.

Matrix Bot v1 was built as a live-tested execution platform for those problems.

## Access

The full implementation is private.

For commercial access, collaboration, or technical due diligence, contact the repository owner through GitHub.

## Status

Matrix Bot v1 is a historical private-system showcase. Current development focus has moved toward a cleaner multi-venue architecture and Kalshi research.

## Disclaimer

This repository is for technical overview only. It is not financial advice, not a trading recommendation, and not an offer to manage funds. Prediction-market trading involves risk. Users are responsible for complying with all applicable laws, venue rules, and terms of service.
