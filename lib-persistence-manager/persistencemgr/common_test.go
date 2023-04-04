package persistencemgr

import (
	"fmt"
	"reflect"
	"testing"
)

func TestByteSlicess(t *testing.T) {
	reply := []interface{}{[]byte("task:task7ce6d6d2-129c-412d-93bf-e11c7e1b7c74"), []byte("task:task88304bc2-74ba-48b6-97e0-f6e11bfddfb1")}
	result, err := ByteSlices(reply, nil)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	expected := [][]byte{[]byte("task:task7ce6d6d2-129c-412d-93bf-e11c7e1b7c74"), []byte("task:task88304bc2-74ba-48b6-97e0-f6e11bfddfb1")}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result. Expected %v, got %v", expected, result)
	}
}

func TestSliceHelper(t *testing.T) {
	reply := []interface{}{[]byte("task:task7ce6d6d2-129c-412d-93bf-e11c7e1b7c74"), []byte("task:task88304bc2-74ba-48b6-97e0-f6e11bfddfb1")}
	var result [][]byte
	err := sliceHelper(reply, nil, "SliceHelper", func(n int) { result = make([][]byte, n) }, func(i int, v interface{}) error {
		p, ok := v.([]byte)
		if !ok {
			return fmt.Errorf("Unexpected element type")
		}
		result[i] = p
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	expected := [][]byte{[]byte("task:task7ce6d6d2-129c-412d-93bf-e11c7e1b7c74"), []byte("task:task88304bc2-74ba-48b6-97e0-f6e11bfddfb1")}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result. Expected %v, got %v", expected, result)
	}
}
