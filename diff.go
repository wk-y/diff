// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package diff

import (
	"fmt"
)

// Diff takes two string slices and tries to find as many identical lines between them.
// The output is an array of added, removed, and identical parts such that:
// a = removed and identical lines
// b = added and identical lines
func Diff(a, b []string) []DiffPart {
	d := DiffAlgorithm(len(a), len(b), func(i, j int) bool {
		return a[i] == b[j]
	})
	result := make([]DiffPart, len(d))
	var i, j int
	for k, action := range d {
		result[k].Action = action
		switch action {
		case DiffIdentical:
			result[k].Value = a[i]
			i++
			j++
		case DiffAdded:
			result[k].Value = b[j]
			j++
		case DiffRemoved:
			result[k].Value = a[i]
			i++
		}
	}
	return result
	// return WeightedDiff(a, b, func(_ string) int { return 1 })
}

type DiffAction int

const (
	DiffIdentical DiffAction = iota
	DiffAdded
	DiffRemoved
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

func DiffAlgorithm(n, m int, equal func(i, j int) bool) []DiffAction {
	if n == m {
		same := true
		for i := 0; i < n; i++ {
			if !equal(i, i) {
				same = false
				break
			}
		}
		if same {
			return make([]DiffAction, n)
		}
	}

	// dp[i][j] = Number of matches possible for a[i..end] and b[j..end]
	// The last row and column are for matching nothing against nothing (a[n..end] and b[n..end]) for sake of code simplicity
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, m+1)
	}

	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			matching := 0
			if equal(i, j) {
				matching = 1
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
	diff := []DiffAction{}
	for i < n && j < m {
		if equal(i, j) {
			diff = append(diff, DiffIdentical)
			i++
			j++
		} else if dp[i+1][j] < dp[i][j+1] { // In case of tie, prefer to put the removal first
			diff = append(diff, DiffAdded)
			j++
		} else {
			diff = append(diff, DiffRemoved)
			i++
		}
	}

	// Process remaining removals
	for ; i < n; i++ {
		diff = append(diff, DiffRemoved)
	}
	// Process remaining additions
	for ; j < m; j++ {
		diff = append(diff, DiffAdded)
	}

	return diff
}
