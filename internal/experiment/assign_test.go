package experiment

import "testing"

func TestVariantForUser_deterministic(t *testing.T) {
	variants := []string{"a", "b", "c"}
	v1 := VariantForUser("user-1", "exp1", variants)
	v2 := VariantForUser("user-1", "exp1", variants)
	if v1 != v2 {
		t.Fatalf("expected sticky assignment, got %q vs %q", v1, v2)
	}
}

func TestVariantForUser_differentUsers(t *testing.T) {
	variants := []string{"x", "y"}
	a := VariantForUser("u1", "e", variants)
	b := VariantForUser("u2", "e", variants)
	// Extremely unlikely both collide forever; just ensure non-empty
	if a == "" || b == "" {
		t.Fatal("empty variant")
	}
}
