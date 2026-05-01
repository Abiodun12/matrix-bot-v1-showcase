package book

import (
	"testing"
	"time"
)

func TestTopOfBook(t *testing.T) {
	b := New()
	now := time.Unix(100, 0)
	b.Apply(now,
		Delta{Side: Bid, Price: 715, Size: 100},
		Delta{Side: Bid, Price: 716, Size: 50},
		Delta{Side: Ask, Price: 718, Size: 70},
		Delta{Side: Ask, Price: 719, Size: 90},
	)

	top := b.Top()
	if !top.OK {
		t.Fatal("top not ok")
	}
	if top.BidPrice != 716 || top.AskPrice != 718 || top.Spread != 2 {
		t.Fatalf("top=%+v", top)
	}
}

func TestApplyRemovesLevel(t *testing.T) {
	b := New()
	now := time.Unix(100, 0)
	b.Apply(now,
		Delta{Side: Bid, Price: 715, Size: 100},
		Delta{Side: Ask, Price: 718, Size: 70},
	)
	b.Apply(now.Add(time.Second), Delta{Side: Bid, Price: 715, Size: 0})

	if b.Top().OK {
		t.Fatal("expected incomplete book after bid removal")
	}
}

func TestFreshness(t *testing.T) {
	b := New()
	now := time.Unix(100, 0)
	if b.Fresh(now, time.Second) {
		t.Fatal("empty book should not be fresh")
	}

	b.Apply(now, Delta{Side: Bid, Price: 715, Size: 100})
	if !b.Fresh(now.Add(500*time.Millisecond), time.Second) {
		t.Fatal("book should be fresh")
	}
	if b.Fresh(now.Add(2*time.Second), time.Second) {
		t.Fatal("book should be stale")
	}
}

func TestLevelsAreSorted(t *testing.T) {
	b := New()
	now := time.Unix(100, 0)
	b.Apply(now,
		Delta{Side: Bid, Price: 714, Size: 10},
		Delta{Side: Bid, Price: 716, Size: 10},
		Delta{Side: Ask, Price: 719, Size: 10},
		Delta{Side: Ask, Price: 718, Size: 10},
	)

	bids := b.Levels(Bid)
	asks := b.Levels(Ask)
	if bids[0].Price != 716 || bids[1].Price != 714 {
		t.Fatalf("bids=%+v", bids)
	}
	if asks[0].Price != 718 || asks[1].Price != 719 {
		t.Fatalf("asks=%+v", asks)
	}
}
