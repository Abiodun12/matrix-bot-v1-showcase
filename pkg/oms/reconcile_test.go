package oms

import "testing"

func TestReconcileKeepsMatchingOrder(t *testing.T) {
	actions := Reconcile(
		[]DesiredOrder{{Key: "market:no:buy", Side: "buy", Price: 715, Size: 500}},
		[]OpenOrder{{ID: "ord-1", Key: "market:no:buy", Side: "buy", Price: 715, Size: 500}},
	)

	if len(actions) != 1 {
		t.Fatalf("actions len=%d", len(actions))
	}
	if actions[0].Type != Keep {
		t.Fatalf("action=%s, want keep", actions[0].Type)
	}
}

func TestReconcileCancelsStaleAndPostsMissing(t *testing.T) {
	actions := Reconcile(
		[]DesiredOrder{{Key: "market:no:buy", Side: "buy", Price: 716, Size: 500}},
		[]OpenOrder{{ID: "ord-1", Key: "market:no:buy", Side: "buy", Price: 715, Size: 500}},
	)

	if len(actions) != 2 {
		t.Fatalf("actions len=%d", len(actions))
	}
	if actions[0].Type != Cancel {
		t.Fatalf("first action=%s, want cancel", actions[0].Type)
	}
	if actions[1].Type != Post {
		t.Fatalf("second action=%s, want post", actions[1].Type)
	}
}

func TestReconcileCancelsUndesiredOpenOrders(t *testing.T) {
	actions := Reconcile(nil, []OpenOrder{{ID: "ord-1", Key: "market:no:buy", Side: "buy", Price: 715, Size: 500}})

	if len(actions) != 1 {
		t.Fatalf("actions len=%d", len(actions))
	}
	if actions[0].Type != Cancel || actions[0].Reason != "not_desired" {
		t.Fatalf("action=%+v, want not_desired cancel", actions[0])
	}
}
