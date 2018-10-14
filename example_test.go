package batcher_test

import (
	"fmt"
	"time"

	"github.com/aperdana/batcher"
)

func Example() {
	// Create a new Batcher
	b := batcher.New(func(batched []interface{}) {
		if len(batched) > 0 {
			fmt.Printf("Here's the batched items: %v\n", batched)
		}
	}, 100*time.Millisecond)

	// Starts listening for items
	b.Listen()

	// Batch item
	b.Batch("I")
	b.Batch("Am")
	b.Batch("A")
	b.Batch("Gopher")

	time.Sleep(110 * time.Millisecond)

	b.Batch("2nd batch")

	time.Sleep(100 * time.Millisecond)

	// Output:
	// Here's the batched items: [I Am A Gopher]
	// Here's the batched items: [2nd batch]
}
