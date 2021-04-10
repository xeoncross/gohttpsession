package sessiontoken_test

import (
	"github.com/xeoncross/gohttpsession/sessiontoken"
	"testing"
)

func TestNew(t *testing.T) {

	id := sessiontoken.New(32)

	if id == nil {
		t.Fatal("failed to generate sessiontoken")
	}

	encoded := sessiontoken.Encode(id)
	decoded := sessiontoken.Decode(encoded)

	if decoded == nil {
		t.Fatal("failed to decode sessiontoken")
	}

	for i := range id {
		if decoded[i] != id[i] {
			t.Errorf("%b != %b", id, decoded)
			break
		}
	}

}

func TestFailures(t *testing.T) {

	encoded := sessiontoken.Encode([]byte{})
	if encoded != "" {
		t.Fatal("should be empty string")
	}

	for _, s := range []string{" ", "\t", "0", "", "a"} {
		decoded := sessiontoken.Decode(s)
		if decoded != nil {
			t.Fatalf("should be empty slice: %v", decoded)
		}
	}

}
