package realnamex

import (
	"context"
	"testing"
)

func TestIdCardOCRVerification(t *testing.T) {
	client, err := NewClient("", "", "")
	if err != nil {
		t.Fatal(err)
	}

	b, err := client.IdCardOCRVerification(context.Background(), "xxx", "33333")
	if err != nil {
		t.Fatal(err)
	}

	t.Error(b)
	t.Log(b)
}
