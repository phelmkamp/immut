// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package romaps_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/phelmkamp/immut/romaps"
)

func Example() {
	m1 := romaps.Freeze(map[string]int{"foo": 42, "bar": 7})
	m2 := map[string]int{"fiz": 3}
	fmt.Println(len(romaps.Keys(m1)))
	romaps.Copy(m2, m1)
	fmt.Println(m2)
	// Output: 2
	// map[bar:7 fiz:3 foo:42]
}

// Example_context demonstrates that a read-only Map can be used
// as a value in context.Context.
func Example_context() {
	type ctxKey string
	cfgKey := ctxKey("cfg")
	ctx := context.WithValue(context.Background(), cfgKey, romaps.Freeze(map[string]string{
		"foo": "42",
		"bar": "baz",
	}))

	cfg, _ := ctx.Value(cfgKey).(romaps.Map[string, string])
	if !cfg.IsNil() {
		v, _ := cfg.Index("foo")
		fmt.Println(v)
	}
	// Output: 42
}

func TestMap_Index(t *testing.T) {
	type fields struct {
		m map[string]int
	}
	type args struct {
		k string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantV  int
		wantOk bool
	}{
		{
			name: "nil",
			fields: fields{
				m: nil,
			},
			args: args{
				k: "foo",
			},
			wantV:  0,
			wantOk: false,
		},
		{
			name: "empty",
			fields: fields{
				m: make(map[string]int),
			},
			args: args{
				k: "foo",
			},
			wantV:  0,
			wantOk: false,
		},
		{
			name: "one",
			fields: fields{
				m: map[string]int{"foo": 42},
			},
			args: args{
				k: "foo",
			},
			wantV:  42,
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := romaps.Freeze(tt.fields.m)
			gotV, gotOk := m.Index(tt.args.k)
			if !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("Index() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Index() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMap_IsNil(t *testing.T) {
	type fields struct {
		m map[int]int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "nil",
			fields: fields{
				m: nil,
			},
			want: true,
		},
		{
			name: "empty",
			fields: fields{
				m: make(map[int]int),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := romaps.Freeze(tt.fields.m)
			if got := m.IsNil(); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_Len(t *testing.T) {
	type fields struct {
		m map[float32]string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "nil",
			fields: fields{
				m: nil,
			},
			want: 0,
		},
		{
			name: "empty",
			fields: fields{
				m: make(map[float32]string),
			},
			want: 0,
		},
		{
			name: "one",
			fields: fields{
				m: map[float32]string{42.0: "foo"},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := romaps.Freeze(tt.fields.m)
			if got := m.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_String(t *testing.T) {
	type fields struct {
		m map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil",
			fields: fields{
				m: nil,
			},
			want: "map[]",
		},
		{
			name: "empty",
			fields: fields{
				m: make(map[string]int),
			},
			want: "map[]",
		},
		{
			name: "one",
			fields: fields{
				m: map[string]int{"foo": 42},
			},
			want: "map[foo:42]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := romaps.Freeze(tt.fields.m)
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_Do(t *testing.T) {
	m := romaps.Freeze(map[int]uint8{0: 1, 1: 2, 3: 4})
	var n uint8
	count := func(int, uint8) { n++ }
	m.Do(doN(count, 2))
	if n != 2 {
		t.Errorf("count after Do() = %v, want %v", n, 3)
	}
	n = 0
	m.Do(doN(count, 4))
	if n != 3 {
		t.Errorf("count after Do() = %v, want %v", n, 7)
	}
}

func doN(f func(int, uint8), n int) func(int, uint8) bool {
	return func(k int, v uint8) bool { f(k, v); n--; return n > 0 }
}
