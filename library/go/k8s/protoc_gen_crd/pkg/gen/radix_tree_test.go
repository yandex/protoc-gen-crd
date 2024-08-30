package gen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	truePtr = func() *bool { val := true; return &val }()
)

func TestRadixTree(t *testing.T) {
	type testCase struct {
		name      string
		paths     map[string]bool
		childPath string
		wantVal   *bool
		wantOk    bool
	}
	tests := []testCase{
		{
			"nil map",
			nil,
			"a",
			nil,
			false,
		},
		{
			"empty map",
			map[string]bool{},
			"a",
			nil,
			false,
		},
		{
			"one path",
			map[string]bool{"a": true},
			"a",
			truePtr,
			true,
		},
		{
			"one wrong path",
			map[string]bool{"a": true},
			"b",
			nil,
			false,
		},
		{
			"multiple paths",
			map[string]bool{"a/b": true, "a/c": true},
			"a/b",
			truePtr,
			true,
		},
		{
			"wrong child",
			map[string]bool{"a/b": true, "a/c": true},
			"a/d",
			nil,
			false,
		},
		{
			"multiple paths with wrong path",
			map[string]bool{"a/b": true, "a/c": true},
			"d/e",
			nil,
			false,
		},
		{
			"nonexistent subchild",
			map[string]bool{"a/b": true, "a/c": true},
			"a/c/d",
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := NewRadixTree[bool]()
			for path, val := range tt.paths {
				pathList := strings.Split(path, "/")
				tree.Add(pathList, val)
			}
			childPathList := strings.Split(tt.childPath, "/")
			for _, path := range childPathList {
				tree = tree.Child(path)
			}
			gotVal, gotOk := tree.Value()
			assert.Equal(t, tt.wantVal, gotVal)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}
