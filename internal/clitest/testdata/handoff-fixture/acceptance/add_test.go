package acceptance

import "testing"

func TestAdd(t *testing.T) {
	if got := Add(2, 2); got != 4 {
		t.Fatalf("Add(2, 2) = %d, want 4", got)
	}
	if got := Add(-1, 1); got != 0 {
		t.Fatalf("Add(-1, 1) = %d, want 0", got)
	}
	if got := Add(0, 0); got != 0 {
		t.Fatalf("Add(0, 0) = %d, want 0", got)
	}
}
