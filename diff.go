/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package diff

import "fmt"

// WeightedDiff takes in two sequences a and b, and tries to create a diff that
// maximizes the sum of w(s) where s is marked identical in the diff.
// The algorithm used is the basic O(nm) LCS algorithm.
// The returned diff has the property that a and b can be reconstructed:
// a = diff parts that are identical or removed
// b = diff parts that are identical or added
func WeightedDiff(a, b []string, w func(string) int) (diff []DiffPart) {

	n := len(a)
	m := len(b)

	// dp[i][j] = Number of matches possible for a[i..end] and b[j..end]
	// The last row and column are for matching nothing against nothing (a[n..end] and b[n..end]) for sake of code simplicity
	dp := make([][]int, n+1)
	for i := range n + 1 {
		dp[i] = make([]int, m+1)
	}

	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			matching := 0
			if a[i] == b[j] {
				matching = w(a[i])
			}

			dp[i][j] = max(dp[i+1][j+1]+matching, dp[i+1][j], dp[i][j+1])
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
		} else if dp[i+1][j] < dp[i][j+1] {
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

	// Process remaining additions
	for ; j < m; j++ {
		diff = append(diff, DiffPart{
			Action: DiffAdded,
			Value:  b[j],
		})
	}

	// Process remaining removals
	for ; i < n; i++ {
		diff = append(diff, DiffPart{
			Action: DiffRemoved,
			Value:  a[i],
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
	switch d {
	case DiffAdded:
		return "Added"
	case DiffRemoved:
		return "Removed"
	case DiffIdentical:
		return "Identical"
	}
	panic("unreachable")
}

type DiffPart struct {
	Action DiffAction
	Value  string
}

func (d DiffPart) String() string {
	return fmt.Sprintf("{%v: %v}", d.Action, d.Value)
}
