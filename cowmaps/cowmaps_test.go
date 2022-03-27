package cowmaps_test

import (
	"fmt"
	"testing"

	"github.com/phelmkamp/immut/cowmaps"
	"github.com/phelmkamp/immut/romaps"
)

func Example() {
	m1 := map[string]int{"fiz": 3}
	m2 := cowmaps.CopyOnWrite(m1)
	cowmaps.Copy(&m2, romaps.Freeze(map[string]int{"foo": 42, "bar": 7}))
	fmt.Println(m2)
	fmt.Println(m1)
	// Output: map[bar:7 fiz:3 foo:42]
	// map[fiz:3]
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
			m := cowmaps.CopyOnWrite(tt.fields.m)
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
