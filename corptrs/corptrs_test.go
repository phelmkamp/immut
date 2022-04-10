// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package corptrs_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/phelmkamp/immut/corptrs"
)

type big struct {
	a, b, c, d, e int
}

func fn1(b corptrs.Ptr[big]) {
	fn2(b)
}

func fn2(b corptrs.Ptr[big]) {
	fn3(b)
}

func fn3(b corptrs.Ptr[big]) {
	fn4(b)
}

func fn4(b corptrs.Ptr[big]) {
	fn5(b)
}

func fn5(b corptrs.Ptr[big]) {
	b2 := b.Clone()
	b2.e++
	fmt.Println(b2)
}

func Example() {
	b := big{1, 2, 3, 4, 5}
	p := corptrs.Freeze(&b)
	fn1(p)
	fmt.Println(p)
	// Output: &{1 2 3 4 6}
	// &{1 2 3 4 5}
}

func TestPtr_Clone(t *testing.T) {
	type fields struct {
		p *big
	}
	tests := []struct {
		name   string
		fields fields
		want   *big
	}{
		{
			name: "nil",
			fields: fields{
				p: nil,
			},
			want: nil,
		},
		{
			name: "non-nil",
			fields: fields{
				p: &big{1, 2, 3, 4, 5},
			},
			want: &big{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := corptrs.Freeze(tt.fields.p)
			if got := p.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPtr_IsNil(t *testing.T) {
	type fields struct {
		p *big
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "nil",
			fields: fields{
				p: nil,
			},
			want: true,
		},
		{
			name: "non-nil",
			fields: fields{
				p: &big{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := corptrs.Freeze(tt.fields.p)
			if got := p.IsNil(); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPtr_String(t *testing.T) {
	type fields struct {
		p *big
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil",
			fields: fields{
				p: nil,
			},
			want: "<nil>",
		},
		{
			name: "non-nil",
			fields: fields{
				p: &big{1, 2, 3, 4, 5},
			},
			want: "&{1 2 3 4 5}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := corptrs.Freeze(tt.fields.p)
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
