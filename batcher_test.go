package batcher

import (
	"reflect"
	"testing"
	"time"
)

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

func TestSetBufferSize(t *testing.T) {
	type args struct {
		s int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"s <= 0",
			args{-1},
			DefaultBufferSize,
		},
		{
			"s > 0",
			args{1},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(nil, time.Second, SetBufferSize(tt.args.s)); got.bufSize != tt.want {
				t.Errorf("SetBufferSize() result in Buffer.bufSize = %v, want %v", got.bufSize, tt.want)
			}
		})
	}
}

func TestSetMaxBatchSize(t *testing.T) {
	type args struct {
		s int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"s <= 0",
			args{-1},
			0,
		},
		{
			"s > 0",
			args{1},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(nil, time.Second, SetMaxBatchSize(tt.args.s)); got.size != tt.want {
				t.Errorf("SetMaxBatchSize() result in Buffer.size = %v, want %v", got.size, tt.want)
			}
		})
	}
}
