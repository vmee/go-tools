package realnamex

import (
	"context"
	"testing"
)

func TestTencentRealName(t *testing.T) {
	b, err := TencentRealName(context.Background(), "xxx", "xxx", "111", "2222")
	if err != nil {
		t.Errorf("TencentRealName failed: %v", err)
		return
	}
	if !b {
		t.Error("TencentRealName returned false")
	}
}
