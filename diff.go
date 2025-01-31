// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package diff

import (
	"fmt"
	"reflect"
)

// WeightedDiff takes in two sequences a and b, and tries to create a diff that
// maximizes the sum of w(s) where s is marked identical in the diff.
// The algorithm used is the basic O(nm) LCS algorithm.
// The returned diff has the property that a and b can be reconstructed:
// a = diff parts that are identical or removed
// b = diff parts that are identical or added
func WeightedDiff(a, b []string, w func(string) int) (diff []DiffPart) {
	// If a and b are identical there's no need to do anything
	if reflect.DeepEqual(a, b) {
		for _, part := range a {
			diff = append(diff, DiffPart{
				Action: DiffIdentical,
				Value:  part,
			})
		}
		return
	}

	n := len(a)
	m := len(b)

	// dp[i][j] = Number of matches possible for a[i..end] and b[j..end]
	// The last row and column are for matching nothing against nothing (a[n..end] and b[n..end]) for sake of code simplicity
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, m+1)
	}

	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			matching := 0
			if a[i] == b[j] {
				matching = w(a[i])
			}

			opt := dp[i+1][j+1] + matching
			if x := dp[i+1][j]; x > opt {
				opt = x
			}
			if x := dp[i][j+1]; x > opt {
				opt = x
			}
			dp[i][j] = opt
		}
	}

	// Traceback
	i := 0
	j := 0
	for i < n && j < m {
		if a[i] == b[j] {
			diff = append(diff, DiffPart{
				Action: DiffIdentical,
				Value:  a[i],
			})
			i++
			j++
		} else if dp[i+1][j] < dp[i][j+1] { // In case of tie, prefer to put the removal first
			diff = append(diff, DiffPart{
				Action: DiffAdded,
				Value:  b[j],
			})
			j++
		} else {
			diff = append(diff, DiffPart{
				Action: DiffRemoved,
				Value:  a[i],
			})
			i++
		}
	}

	// Process remaining removals
	for ; i < n; i++ {
		diff = append(diff, DiffPart{
			Action: DiffRemoved,
			Value:  a[i],
		})
	}
	// Process remaining additions
	for ; j < m; j++ {
		diff = append(diff, DiffPart{
			Action: DiffAdded,
			Value:  b[j],
		})
	}

	return
}

// Diff returns
func Diff(a, b []string) []DiffPart {
	return WeightedDiff(a, b, func(_ string) int { return 1 })
}

// MinCharDiff minimizes the characters added/deleted.
func MinCharDiff(a, b []string) []DiffPart {
	return WeightedDiff(a, b, func(s string) int { return len(s) })
}

type DiffAction int

const (
	DiffAdded DiffAction = iota
	DiffRemoved
	DiffIdentical
)

func (d DiffAction) String() string {
	return []string{
		DiffAdded:     "Added",
		DiffRemoved:   "Removed",
		DiffIdentical: "Identical",
	}[d]
}

type DiffPart struct {
	Action DiffAction
	Value  string
}

func (d DiffPart) String() string {
	return fmt.Sprintf("{%v: %v}", d.Action, d.Value)
}
