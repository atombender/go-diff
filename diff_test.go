package diff_test

import (
	"github.com/atombender/go-diff"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPruneContext(t *testing.T) {
	hunks := diff.Diff(
		[]string{"aaa", "bbb", "ccc", "ddd", "eee", "fff"},
		[]string{"aaa", "bbb", "XXX", "ddd", "eee", "fff"},
	)
	assert.Equal(t, []diff.Hunk{
		{LineNum: 1, Operation: diff.OpUnchanged, Line: "bbb"},
		{LineNum: 2, Operation: diff.OpDelete, Line: "ccc"},
		{LineNum: 2, Operation: diff.OpInsert, Line: "XXX"},
		{LineNum: 3, Operation: diff.OpUnchanged, Line: "ddd"},
	}, diff.PruneContext(hunks, 1))
}

func TestDiff(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		a, b   []string
		expect []diff.Hunk
	}{
		{
			desc:   "empty",
			a:      []string{},
			b:      []string{},
			expect: []diff.Hunk{},
		},
		{
			desc: "unchanged",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"aaa", "bbb", "ccc"},
			expect: []diff.Hunk{
				{Operation: diff.OpUnchanged, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpUnchanged, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpUnchanged, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "insert from empty",
			a:    []string{},
			b:    []string{"aaa", "bbb", "ccc"},
			expect: []diff.Hunk{
				{Operation: diff.OpInsert, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpInsert, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpInsert, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "insert middle",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"aaa", "bbb", "ZZZ", "ccc"},
			expect: []diff.Hunk{
				{Operation: diff.OpUnchanged, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpUnchanged, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpInsert, LineNum: 2, Line: "ZZZ"},
				{Operation: diff.OpUnchanged, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "insert start",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"ZZZ", "aaa", "bbb", "ccc"},
			expect: []diff.Hunk{
				{Operation: diff.OpInsert, LineNum: 0, Line: "ZZZ"},
				{Operation: diff.OpUnchanged, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpUnchanged, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpUnchanged, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "insert end",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"aaa", "bbb", "ccc", "ZZZ"},
			expect: []diff.Hunk{
				{Operation: diff.OpUnchanged, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpUnchanged, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpUnchanged, LineNum: 2, Line: "ccc"},
				{Operation: diff.OpInsert, LineNum: 3, Line: "ZZZ"},
			},
		},
		{
			desc: "delete all",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{},
			expect: []diff.Hunk{
				{Operation: diff.OpDelete, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpDelete, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpDelete, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "delete middle",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"aaa", "ccc"},
			expect: []diff.Hunk{
				{Operation: diff.OpUnchanged, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpDelete, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpUnchanged, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "delete start",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"bbb", "ccc"},
			expect: []diff.Hunk{
				{Operation: diff.OpDelete, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpUnchanged, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpUnchanged, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "delete end",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"aaa", "bbb"},
			expect: []diff.Hunk{
				{Operation: diff.OpUnchanged, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpUnchanged, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpDelete, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "replace all",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"xxx", "yyy", "zzz"},
			expect: []diff.Hunk{
				{Operation: diff.OpDelete, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpInsert, LineNum: 0, Line: "xxx"},
				{Operation: diff.OpDelete, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpInsert, LineNum: 1, Line: "yyy"},
				{Operation: diff.OpDelete, LineNum: 2, Line: "ccc"},
				{Operation: diff.OpInsert, LineNum: 2, Line: "zzz"},
			},
		},
		{
			desc: "change middle",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"aaa", "ZZZ", "ccc"},
			expect: []diff.Hunk{
				{Operation: diff.OpUnchanged, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpDelete, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpInsert, LineNum: 1, Line: "ZZZ"},
				{Operation: diff.OpUnchanged, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "change start",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"ZZZ", "bbb", "ccc"},
			expect: []diff.Hunk{
				{Operation: diff.OpDelete, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpInsert, LineNum: 0, Line: "ZZZ"},
				{Operation: diff.OpUnchanged, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpUnchanged, LineNum: 2, Line: "ccc"},
			},
		},
		{
			desc: "change end",
			a:    []string{"aaa", "bbb", "ccc"},
			b:    []string{"aaa", "bbb", "ZZZ"},
			expect: []diff.Hunk{
				{Operation: diff.OpUnchanged, LineNum: 0, Line: "aaa"},
				{Operation: diff.OpUnchanged, LineNum: 1, Line: "bbb"},
				{Operation: diff.OpDelete, LineNum: 2, Line: "ccc"},
				{Operation: diff.OpInsert, LineNum: 2, Line: "ZZZ"},
			},
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			actual := diff.Diff(tc.a, tc.b)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
