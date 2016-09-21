// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pathos

import (
	"runtime"
	"testing"
)

func TestTrimCommonSuffix(t *testing.T) {
	list := []struct {
		slash          rune
		base, suffix   string
		result, common string
	}{
		{slash: '/', base: "/a/b/c", suffix: "/x/y/b/c", result: "/a", common: "b/c"},
		{slash: '/', base: "/arg/borish/client", suffix: "fooish/client", result: "/arg/borish", common: "client"},
		{slash: '/', base: "/tmp/vendor_272718190/src/co2/go/pk1/", suffix: "co2/go/pk1", result: "/tmp/vendor_272718190/src", common: "co2/go/pk1"},
		{slash: '/', base: "/home/daniel/code/go/src/.cache/govendor/github.com/raphael/goa", suffix: "github.com/raphael/goa", result: "/home/daniel/code/go/src/.cache/govendor", common: "github.com/raphael/goa"},
		{slash: '\\', base: `d:\bob\alice\noob`, suffix: `c:\tmp\foo\alice\noob`, result: `d:\bob`, common: `alice\noob`},
	}

	for _, item := range list {
		slashSep = item.slash
		got, common := TrimCommonSuffix(item.base, item.suffix)
		if got != item.result || common != item.common {
			t.Errorf("For %#v got %q, common: %q", item, got, common)
		}
	}
}

func TestGoEnv(t *testing.T) {
	list := []struct {
		line   string
		name   string
		result string
		ok     bool
	}{
		{`set GOROOT=C:\Foo\Bar`, "GOROOT", `C:\Foo\Bar`, true},
		{`set GOPATH=C:\Foo\Bar`, "GOROOT", ``, false},
		{`set GOROOT=`, "GOROOT", ``, true},
		{`GOROOT="/foo/bar"`, "GOROOT", `/foo/bar`, true},
		{`GOPATH="/foo/bar"`, "GOROOT", ``, false},
		{`GOROOT=""`, "GOROOT", ``, true},
	}

	for index, item := range list {
		result, ok := GoEnv(item.name, item.line)
		if ok != item.ok {
			t.Errorf("index %d line %#v expected ok %t but got %t", index, item, item.ok, ok)
		}
		if result != item.result {
			t.Errorf("index %d line %#v expected result %q but got %q", index, item, item.result, result)
		}
	}
}

func TestEscapeImport(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skipf("not implemented for %s", runtime.GOOS)
	}

	list := []struct {
		path   string
		result string
	}{
		{`C:\Foo\Bar`, `C:\Foo\Bar`},
		{`D:Foo\Bar`, `D:Foo\Bar`},
		{`C:\Foo\Bar:Baz?Buzz`, `C:\Foo\Bar_Baz_Buzz`},
		{`Foo\Bar_z`, `Foo\Bar_z`},
		{`127.0.0.1:9000\herp\derp`, `127.0.0.1_9000\herp\derp`},
		{`foo\bar::bar\baz`, `foo\bar__bar\baz`},
	}
	for index, item := range list {
		result := EscapeImport(item.path)
		if result != item.result {
			t.Errorf("index %d line %#v expected result %q but got %q", index, item, item.result, result)
		}
	}
}
