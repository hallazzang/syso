package rsrc

import (
	"testing"
)

func TestIdentifier(t *testing.T) {
	for _, tc := range []struct {
		inp        interface{}
		isID       bool
		id         int
		name       string
		shouldFail bool
	}{
		{inp: 0, isID: true, id: 0},
		{inp: "foo", isID: false, name: "foo"},
		{inp: 1.2, shouldFail: true},
	} {
		t.Logf("input: %v", tc.inp)
		id, name, err := identifier(tc.inp)
		if tc.shouldFail {
			if err == nil {
				t.Fatal("should fail, but not failed")
			}
			continue
		} else if !tc.shouldFail && err != nil {
			t.Fatal(err)
		}
		if tc.isID {
			if id == nil {
				t.Fatal("result doesn't contain valid id address")
			} else if *id != tc.id {
				t.Fatalf("wrong id; got %d, expected %d", *id, tc.id)
			}
		} else {
			if name == nil {
				t.Fatal("result doesn't contain valid name address")
			} else if *name != tc.name {
				t.Fatalf("wrong name; got %s, expected %s", *name, tc.name)
			}
		}
	}
}
