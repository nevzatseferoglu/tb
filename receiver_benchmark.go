package main

import (
	"testing"
)

type MockType struct {
	i8   int8
	ui64 uint64
}

func (t *MockType) PointerReceiver() (int8, uint64) {
	return t.i8, t.ui64
}

func (t MockType) ValueReceiver() (int8, uint64) {
	return t.i8, t.ui64
}

func BenchmarkChangePointerReceiver(b *testing.B) {}

func BenchmarkChangeItValueReceiver(b *testing.B) {}
