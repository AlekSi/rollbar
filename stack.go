// Code in this file is based on code from Go standard library.
// See package runtime/debug (file src/pkg/runtime/debug/stack.go).
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

package rollbar

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"strings"
)

var SkipPath = "github.com/AlekSi/rollbar"

type frame struct {
	Filename string `json:"filename"`
	Lineno   int    `json:"lineno"`
	Method   string `json:"method"`
	Code     string `json:"code"`
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// stack implements Stack, skipping 2 frames
func stack() (frames []frame) {
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := 3; ; i++ { // Caller we care about is the user, 3 frames up
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if len(SkipPath) > 0 && strings.Contains(file, SkipPath) {
			continue
		}

		f := frame{Filename: file, Lineno: line}

		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				frames = append(frames, f)
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		line-- // in stack trace, lines are 1-indexed but our array is 0-indexed
		f.Method = string(function(pc))
		f.Code = string(source(lines, line))
		frames = append(frames, f)
	}
	return
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.Trim(lines[n], " \t")
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
