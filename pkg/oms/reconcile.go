package oms

import (
	"sort"
)

// DesiredOrder is the minimal intent the strategy wants resting on the book.
// The public demo intentionally keeps this venue-neutral.
type DesiredOrder struct {
	Key    string
	Side   string
	Price  int64
	Size   int64
	Reduce bool
}

// OpenOrder is a currently resting order from the venue.
type OpenOrder struct {
	ID     string
	Key    string
	Side   string
	Price  int64
	Size   int64
	Reduce bool
}

// ActionType is the reconciler output.
type ActionType string

const (
	Keep   ActionType = "keep"
	Cancel ActionType = "cancel"
	Post   ActionType = "post"
)

// Action describes the smallest set of operations needed to move from open
// orders to desired orders.
type Action struct {
	Type    ActionType
	OrderID string
	Reason  string
	Desired DesiredOrder
	Open    OpenOrder
}

// Reconcile computes a deterministic cancel/keep/post plan.
//
// This mirrors the core shape of a live OMS loop without exposing venue
// signing, auth, strategy, or deployment details.
func Reconcile(desired []DesiredOrder, open []OpenOrder) []Action {
	desiredByKey := make(map[string]DesiredOrder, len(desired))
	for _, d := range desired {
		desiredByKey[d.Key] = d
	}

	openByKey := make(map[string][]OpenOrder, len(open))
	for _, o := range open {
		openByKey[o.Key] = append(openByKey[o.Key], o)
	}
	for key := range openByKey {
		sort.Slice(openByKey[key], func(i, j int) bool {
			return openByKey[key][i].ID < openByKey[key][j].ID
		})
	}

	var actions []Action
	for key, orders := range openByKey {
		d, wanted := desiredByKey[key]
		if !wanted {
			for _, o := range orders {
				actions = append(actions, Action{Type: Cancel, OrderID: o.ID, Reason: "not_desired", Open: o})
			}
			continue
		}

		kept := false
		for _, o := range orders {
			if !kept && matches(d, o) {
				actions = append(actions, Action{Type: Keep, OrderID: o.ID, Reason: "matches", Desired: d, Open: o})
				kept = true
				continue
			}
			actions = append(actions, Action{Type: Cancel, OrderID: o.ID, Reason: "duplicate_or_stale", Desired: d, Open: o})
		}
		if kept {
			delete(desiredByKey, key)
			continue
		}
		actions = append(actions, Action{Type: Post, Reason: "stale_replace", Desired: d})
		delete(desiredByKey, key)
	}

	var missing []DesiredOrder
	for key, d := range desiredByKey {
		if len(openByKey[key]) == 0 {
			missing = append(missing, d)
		}
	}
	sort.Slice(missing, func(i, j int) bool { return missing[i].Key < missing[j].Key })
	for _, d := range missing {
		actions = append(actions, Action{Type: Post, Reason: "missing", Desired: d})
	}

	return actions
}

func matches(d DesiredOrder, o OpenOrder) bool {
	return d.Key == o.Key &&
		d.Side == o.Side &&
		d.Price == o.Price &&
		d.Size == o.Size &&
		d.Reduce == o.Reduce
}
