// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package roslices_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/phelmkamp/immut/roslices"
	expslices "golang.org/x/exp/slices"
)

func Example() {
	ints1 := roslices.Freeze([]int{2, 1, 3})
	if !roslices.IsSorted(ints1) {
		// must clone to sort
		ints2 := roslices.Clone(ints1)
		expslices.Sort(ints2)
		fmt.Println(ints1, ints2)
	}
	// Output: [2 1 3] [1 2 3]
}

func TestSlice_Cap(t *testing.T) {
	type fields struct {
		s []int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "nil",
			fields: fields{
				s: nil,
			},
			want: 0,
		},
		{
			name: "empty",
			fields: fields{
				s: []int{},
			},
			want: 0,
		},
		{
			name: "zero",
			fields: fields{
				s: make([]int, 0, 0),
			},
			want: 0,
		},
		{
			name: "one",
			fields: fields{
				s: make([]int, 0, 1),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := roslices.Freeze(tt.fields.s)
			if got := s.Cap(); got != tt.want {
				t.Errorf("Cap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_Index(t *testing.T) {
	type fields struct {
		s []string
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "0",
			fields: fields{
				s: []string{"0", "1"},
			},
			args: args{
				i: 0,
			},
			want: "0",
		},
		{
			name: "1",
			fields: fields{
				s: []string{"0", "1"},
			},
			args: args{
				i: 1,
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := roslices.Freeze(tt.fields.s)
			if got := s.Index(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_IsNil(t *testing.T) {
	type fields struct {
		s []float32
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "nil",
			fields: fields{
				s: nil,
			},
			want: true,
		},
		{
			name: "non-nil",
			fields: fields{
				s: make([]float32, 0),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := roslices.Freeze(tt.fields.s)
			if got := s.IsNil(); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_Len(t *testing.T) {
	type fields struct {
		s []bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "nil",
			fields: fields{
				s: nil,
			},
			want: 0,
		},
		{
			name: "empty",
			fields: fields{
				s: []bool{},
			},
			want: 0,
		},
		{
			name: "one",
			fields: fields{
				s: make([]bool, 1),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := roslices.Freeze(tt.fields.s)
			if got := s.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_String(t *testing.T) {
	type fields struct {
		s []int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil",
			fields: fields{
				s: make([]int, 0),
			},
			want: "[]",
		},
		{
			name: "empty",
			fields: fields{
				s: make([]int, 0),
			},
			want: "[]",
		},
		{
			name: "one",
			fields: fields{
				s: []int{0},
			},
			want: "[0]",
		},
		{
			name: "two",
			fields: fields{
				s: []int{0, 1},
			},
			want: "[0 1]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := roslices.Freeze(tt.fields.s)
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_Slice(t *testing.T) {
	ints := []int{0, 1, 2}

	type fields struct {
		s []int
	}
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   roslices.Slice[int]
	}{
		{
			name: "nil",
			fields: fields{
				s: nil,
			},
			args: args{
				i: 0,
				j: 0,
			},
			want: roslices.Freeze[int](nil),
		},
		{
			name: "empty",
			fields: fields{
				s: make([]int, 0),
			},
			args: args{
				i: 0,
				j: 0,
			},
			want: roslices.Freeze([]int{}),
		},
		{
			name: "0:2",
			fields: fields{
				s: []int{0, 1, 2},
			},
			args: args{
				i: 0,
				j: 2,
			},
			want: roslices.Freeze(ints[0:2]),
		},
		{
			name: "1:3",
			fields: fields{
				s: []int{0, 1, 2},
			},
			args: args{
				i: 1,
				j: 3,
			},
			want: roslices.Freeze(ints[1:3]),
		},
		{
			name: "0:3",
			fields: fields{
				s: []int{0, 1, 2},
			},
			args: args{
				i: 0,
				j: 3,
			},
			want: roslices.Freeze(ints[0:3]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := roslices.Freeze(tt.fields.s)
			if got := s.Slice(tt.args.i, tt.args.j); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	ints := []int{0, 1}
	wantS := make([]int, len(ints))
	wantN := copy(wantS, ints)
	gotS := make([]int, len(ints))
	if gotN := roslices.Copy(gotS, roslices.Freeze(ints)); gotN != wantN || !reflect.DeepEqual(gotS, wantS) {
		t.Errorf("Copy() = %v %v, want %v %v", gotS, gotN, wantS, wantN)
	}
}
