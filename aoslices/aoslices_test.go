// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package aoslices

import (
	"reflect"
	"testing"
)

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
				s: make([]int, 0),
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
			s := Slice[int]{s: tt.fields.s}
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
			s := Slice[string]{s: tt.fields.s}
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
			s := Slice[float32]{s: tt.fields.s}
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
			s := Slice[bool]{s: tt.fields.s}
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
			s := Slice[int]{s: tt.fields.s}
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	ints := []int{0, 1}
	wantS := make([]int, len(ints))
	wantN := copy(wantS, ints)
	gotS := make([]int, len(ints))
	if gotN := Copy(gotS, Slice[int]{s: ints}); gotN != wantN || !reflect.DeepEqual(gotS, wantS) {
		t.Errorf("Copy() = %v %v, want %v %v", gotS, gotN, wantS, wantN)
	}
}

func TestSlice_Slice(t *testing.T) {
	ints := []int{0, 1, 2}

	type fields struct {
		s []int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Slice[int]
	}{
		{
			name: "nil",
			fields: fields{
				s: nil,
			},
			args: args{
				i: 0,
			},
			want: Slice[int]{},
		},
		{
			name: "empty",
			fields: fields{
				s: make([]int, 0),
			},
			args: args{
				i: 0,
			},
			want: Slice[int]{s: make([]int, 0)},
		},
		{
			name: "1:",
			fields: fields{
				s: []int{0, 1, 2},
			},
			args: args{
				i: 1,
			},
			want: Slice[int]{s: ints[1:3]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Slice[int]{s: tt.fields.s}
			if got := s.Slice(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMake(t *testing.T) {
	type args struct {
		size []int
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
		wantCap int
	}{
		{
			name:    "empty",
			args:    args{},
			wantLen: 0,
			wantCap: 0,
		},
		{
			name: "len",
			args: args{
				size: []int{1},
			},
			wantLen: 1,
			wantCap: 1,
		},
		{
			name: "cap",
			args: args{
				size: []int{1, 2},
			},
			wantLen: 1,
			wantCap: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Make[struct{}](tt.args.size...)
			if got.Len() != tt.wantLen {
				t.Errorf("len after Make() = %v, want %v", got.Len(), tt.wantLen)
			}
			if got.Cap() != tt.wantCap {
				t.Errorf("cap after Make() = %v, want %v", got.Cap(), tt.wantCap)
			}
		})
	}
}

func TestAppend(t *testing.T) {
	s := Slice[int]{}
	if s = Append(s, 1); !reflect.DeepEqual(s, Slice[int]{s: []int{1}}) {
		t.Errorf("Append() = %v, want %v", s, Slice[int]{s: []int{1}})
	}
	if s = Append(s, 2, 3); !reflect.DeepEqual(s, Slice[int]{s: []int{1, 2, 3}}) {
		t.Errorf("Append() = %v, want %v", s, Slice[int]{s: []int{1, 2, 3}})
	}
}
