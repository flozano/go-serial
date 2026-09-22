//
// Copyright 2014-2026 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Testing code idea and fix thanks to @angri
// https://github.com/bugst/go-serial/pull/42

package serial

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func startSocatAndWaitForPort(t *testing.T, ctx context.Context) *exec.Cmd {
	t.Helper()
	cmd := exec.CommandContext(ctx, "socat", "-D", "STDIO", "pty,link=/tmp/faketty")
	r, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Let our fake serial port node appear.
	// socat will write to stderr before starting transfer phase;
	// we don't really care what, just that it did, because then it's ready.
	buf := make([]byte, 1024)
	if _, err = r.Read(buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return cmd
}

func TestSerialReadAndCloseConcurrency(t *testing.T) {

	// Run this test with race detector to actually test that
	// the correct multitasking behaviour is happening.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := startSocatAndWaitForPort(t, ctx)
	go cmd.Wait()

	port, err := Open("/tmp/faketty", &Mode{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := make([]byte, 100)
	go port.Read(buf)
	// let port.Read to start
	time.Sleep(time.Millisecond * 1)
	port.Close()
}

func TestDoubleCloseIsNoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := startSocatAndWaitForPort(t, ctx)
	go cmd.Wait()

	port, err := Open("/tmp/faketty", &Mode{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := port.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := port.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReadIntervalTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := startSocatAndWaitForPort(t, ctx)
	go cmd.Wait()

	port, err := Open("/tmp/faketty", &Mode{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer port.Close()

	if err := port.SetReadIntervalTimeout(0, NoTimeout); err == nil {
		t.Error("a zero interval was accepted")
	}
	if err := port.SetReadIntervalTimeout(20*time.Millisecond, 50*time.Millisecond); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	start := time.Now()
	n, err := port.Read(make([]byte, 10))
	if n != 0 || err != nil {
		t.Errorf("read %d bytes, %v; want the total timeout: 0 and no error", n, err)
	}
	if waited := time.Since(start); waited < 50*time.Millisecond {
		t.Errorf("returned after %v, before the total timeout", waited)
	}
}
