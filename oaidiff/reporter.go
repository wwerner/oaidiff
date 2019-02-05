// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package oaidiff

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/wwerner/oaidiff/internal/value"
	"reflect"
	"strings"

)

type myReporter struct {
	cmp.Option
	diffs  []string // List of differences, possibly truncated
	ndiffs int      // Total number of differences
	nbytes int      // Number of bytes in diffs
	nlines int      // Number of lines in diffs
}

var _ = (*myReporter)(nil)

func (r *myReporter) Report(x, y reflect.Value, eq bool, p cmp.Path) {
	if eq {
		return // Ignore equal results
	}
	const maxBytes = 4096
	const maxLines = 256
	r.ndiffs++
	if r.nbytes < maxBytes && r.nlines < maxLines {
		sx := value.Format(x, value.FormatConfig{UseStringer: true})
		sy := value.Format(y, value.FormatConfig{UseStringer: true})
		if sx == sy {
			// Unhelpful output, so use more exact formatting.
			sx = value.Format(x, value.FormatConfig{PrintPrimitiveType: true})
			sy = value.Format(y, value.FormatConfig{PrintPrimitiveType: true})
		}
		s := fmt.Sprintf("%#v:\n\tXXX: %s\n\tYYY: %s\n", p, sx, sy)
		r.diffs = append(r.diffs, s)
		r.nbytes += len(s)
		r.nlines += strings.Count(s, "\n")
	}
}

func (r *myReporter) String() string {
	s := strings.Join(r.diffs, "")
	if r.ndiffs == len(r.diffs) {
		return s
	}
	return fmt.Sprintf("%s... %d more differences ...", s, r.ndiffs-len(r.diffs))
}
