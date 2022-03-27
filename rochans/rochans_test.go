package rochans_test

import (
	"fmt"
	"testing"

	"github.com/phelmkamp/immut/rochans"
)

func Example() {
	ch := make(chan int)
	roch := rochans.Freeze(ch)
	go func() {
		ch <- 42
	}()
	fmt.Println(roch.Recv())
	// Output: 42 true
}

func TestChan_Cap(t *testing.T) {
	type fields struct {
		ch chan int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "nil",
			fields: fields{
				ch: nil,
			},
			want: 0,
		},
		{
			name: "unbuffered",
			fields: fields{
				ch: make(chan int),
			},
			want: 0,
		},
		{
			name: "one",
			fields: fields{
				ch: make(chan int, 1),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := rochans.Freeze(tt.fields.ch)
			if got := ch.Cap(); got != tt.want {
				t.Errorf("Cap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChan_IsNil(t *testing.T) {
	type fields struct {
		ch chan string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "nil",
			fields: fields{
				ch: nil,
			},
			want: true,
		},
		{
			name: "non-nil",
			fields: fields{
				ch: make(chan string),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := rochans.Freeze(tt.fields.ch)
			if got := ch.IsNil(); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChan_Len(t *testing.T) {
	chone := make(chan bool, 1)
	chone <- true

	type fields struct {
		ch chan bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "nil",
			fields: fields{
				ch: nil,
			},
			want: 0,
		},
		{
			name: "unbuffered",
			fields: fields{
				ch: make(chan bool),
			},
			want: 0,
		},
		{
			name: "zero",
			fields: fields{
				ch: make(chan bool, 1),
			},
			want: 0,
		},
		{
			name: "one",
			fields: fields{
				ch: chone,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := rochans.Freeze(tt.fields.ch)
			if got := ch.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChan_Recv(t *testing.T) {
	ch := make(chan int)
	go func() {
		ch <- 1
		close(ch)
	}()

	roch := rochans.Freeze(ch)
	v, ok := roch.Recv()
	if v != 1 || !ok {
		t.Errorf("Recv() = %v %v, want %v %v", v, ok, 1, true)
	}

	v, ok = roch.Recv()
	if v != 0 || ok {
		t.Errorf("Recv() = %v %v, want %v %v", v, ok, 0, false)
	}
}

func TestChan_String(t *testing.T) {
	ch := make(chan int)

	type fields struct {
		ch chan int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil",
			fields: fields{
				ch: nil,
			},
			want: "<nil>",
		},
		{
			name: "non-nil",
			fields: fields{
				ch: ch,
			},
			want: fmt.Sprint(ch),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := rochans.Freeze(tt.fields.ch)
			if got := ch.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
