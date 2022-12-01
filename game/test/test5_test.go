package main

import "testing"

func TestF1Run(t *testing.T) {
	for i := 0; i < 100000; i++ {
		testf1(i, i)
	}
}
