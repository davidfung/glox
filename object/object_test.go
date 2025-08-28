package object

import (
	"fmt"
	"testing"
)

func TestHashString(t *testing.T) {
	s := "kilroy"
	h := hashString(s, len(s))
	r := uint32(788470611)
	if h != r {
		t.Errorf("hashString(\"%s\")=%d, expect %d\n", s, h, r)
	}
	fmt.Println(s, h)
}
