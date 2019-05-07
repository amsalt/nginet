package test

import (
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestWaitGroup(t *testing.T) {
	wg.Wait()
}
