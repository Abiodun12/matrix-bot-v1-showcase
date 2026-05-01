package book

import (
	"sort"
	"time"
)

// Side identifies which side of the book is being updated.
type Side string

const (
	Bid Side = "bid"
	Ask Side = "ask"
)

// Level is a price level in integer ticks.
type Level struct {
	Price int64
	Size  int64
}

// Delta is a single level update. Size <= 0 removes the level.
type Delta struct {
	Side  Side
	Price int64
	Size  int64
}

// Top is the current best bid/ask snapshot.
type Top struct {
	BidPrice int64
	BidSize  int64
	AskPrice int64
	AskSize  int64
	Spread   int64
	OK       bool
}

// Book is a deterministic in-memory order book. It is intentionally small and
// venue-neutral for the public showcase.
type Book struct {
	bids       map[int64]int64
	asks       map[int64]int64
	lastUpdate time.Time
}

// New returns an empty book.
func New() *Book {
	return &Book{
		bids: make(map[int64]int64),
		asks: make(map[int64]int64),
	}
}

// Apply updates one or more price levels and records the receive timestamp.
func (b *Book) Apply(now time.Time, deltas ...Delta) {
	for _, d := range deltas {
		levels := b.bids
		if d.Side == Ask {
			levels = b.asks
		}

		if d.Size <= 0 {
			delete(levels, d.Price)
			continue
		}
		levels[d.Price] = d.Size
	}
	b.lastUpdate = now
}

// Top returns the best bid and ask. OK is false until both sides exist and the
// book is not crossed.
func (b *Book) Top() Top {
	bid, bidSize, bidOK := bestBid(b.bids)
	ask, askSize, askOK := bestAsk(b.asks)
	if !bidOK || !askOK || bid >= ask {
		return Top{}
	}
	return Top{
		BidPrice: bid,
		BidSize:  bidSize,
		AskPrice: ask,
		AskSize:  askSize,
		Spread:   ask - bid,
		OK:       true,
	}
}

// Fresh reports whether the book has been updated within maxAge.
func (b *Book) Fresh(now time.Time, maxAge time.Duration) bool {
	if b.lastUpdate.IsZero() {
		return false
	}
	return now.Sub(b.lastUpdate) <= maxAge
}

// Levels returns sorted levels for diagnostics and tests.
func (b *Book) Levels(side Side) []Level {
	src := b.bids
	if side == Ask {
		src = b.asks
	}

	out := make([]Level, 0, len(src))
	for price, size := range src {
		out = append(out, Level{Price: price, Size: size})
	}

	sort.Slice(out, func(i, j int) bool {
		if side == Bid {
			return out[i].Price > out[j].Price
		}
		return out[i].Price < out[j].Price
	})
	return out
}

func bestBid(levels map[int64]int64) (price int64, size int64, ok bool) {
	for p, s := range levels {
		if !ok || p > price {
			price, size, ok = p, s, true
		}
	}
	return price, size, ok
}

func bestAsk(levels map[int64]int64) (price int64, size int64, ok bool) {
	for p, s := range levels {
		if !ok || p < price {
			price, size, ok = p, s, true
		}
	}
	return price, size, ok
}
