package main_test

import (
	"testing"
)

func TestChanCap(t *testing.T) {
	for _, tt := range []struct {
		ch   chan string
		want int
	}{
		{ch: make(chan string), want: 0},
		{ch: make(chan string, 1), want: 1},
	} {
		got := cap(tt.ch)
		if got != tt.want {
			t.Errorf("got: %d, want: %d", got, tt.want)
		}
	}
}

func TestRcv(t *testing.T) {
	ch := make(chan bool)
	go func() {
		ch <- true
		ch <- false
		close(ch)
	}()

	got, ok := <-ch
	if got != true || ok != true {
		t.Errorf("got: %v, want: true", got)
	}

	got, ok = <-ch
	if got != false || ok != true {
		t.Errorf("got: %v, want: false", got)
	}

	got, ok = <-ch
	if got != false || ok != false {
		t.Errorf("got: %v, want: false", got)
	}
	got, ok = <-ch
	if got != false || ok != false {
		t.Errorf("got: %v, want: false", got)
	}
}
