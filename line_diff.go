// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package diff

import (
	"github.com/wk-y/diff/internal/strutils"
)

func LineDiff(a, b string) []DiffPart {
	return Diff(strutils.SplitLines(a), strutils.SplitLines(b))
}
