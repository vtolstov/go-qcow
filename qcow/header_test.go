package qcow

import (
	"fmt"
	"os"
	"testing"
)

func TestHeader(t *testing.T) {
	f, err := os.Open("debian-7-x86_64")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer f.Close()

	r, err := NewReader(f)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("%#v\n", r)
}
