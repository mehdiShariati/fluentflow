package auth

import "testing"

func TestHashPassword_CheckPassword(t *testing.T) {
	h, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(h, "correct-horse-battery-staple") {
		t.Fatal("expected password to match")
	}
	if CheckPassword(h, "wrong") {
		t.Fatal("expected wrong password to fail")
	}
}
