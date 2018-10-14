package batcher

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Example() {
	// Create a new Batcher
	b := New(func(batched []interface{}) {
		fmt.Printf("Here's the batched items: %v\n", batched)
	}, time.Second)

	// Starts listening for items
	b.Listen()

	// Batch item
	b.Batch("I")
	b.Batch("Am")
	b.Batch("The")
	b.Batch("Great")

	time.Sleep(time.Second)

	// Should print:
	// Here's the batched items: I Am The Great
}

func TestBatcher_receive(t *testing.T) {
	type fields struct {
		size      int
		timeout   time.Duration
		buffer    chan interface{}
		do        func([]interface{})
		listening bool
	}
	tests := []struct {
		name   string
		fields fields
		seeder func(chan interface{})
		want   []interface{}
	}{
		{
			"simple time-based",
			fields{
				timeout: 100 * time.Millisecond,
			},
			func(input chan interface{}) {
				input <- "My"
				input <- "Name"
				input <- "Is"
				input <- "Wawan"
			},
			[]interface{}{"My", "Name", "Is", "Wawan"},
		},
		{
			"time-based",
			fields{
				timeout: 100 * time.Millisecond,
			},
			func(input chan interface{}) {
				input <- "My"
				input <- "Name"
				input <- "Is"
				input <- "Wawan"

				time.Sleep(1 * time.Second)

				input <- "This should not be in this batch"
			},
			[]interface{}{"My", "Name", "Is", "Wawan"},
		},
		{
			"time-size-based",
			fields{
				timeout: 100 * time.Millisecond,
				size:    2,
			},
			func(input chan interface{}) {
				input <- "My"
				input <- "Name"
				input <- "Is"
				input <- "Wawan"
			},
			[]interface{}{"My", "Name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Batcher{
				size:    tt.fields.size,
				timeout: tt.fields.timeout,
				buffer:  make(chan interface{}, DefaultBufferSize),
			}
			go tt.seeder(b.buffer)
			if got := b.receive(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Batcher.receive() = %v, want %v", got, tt.want)
			}
		})
	}
}
