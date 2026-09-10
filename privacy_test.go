package coas

import "testing"

func TestSealDraftRoundTrip(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	sealed, err := SealDraft([]byte("private human meaning"), key)
	if err != nil {
		t.Fatal(err)
	}
	plain, err := OpenDraft(sealed, key)
	if err != nil {
		t.Fatal(err)
	}
	if string(plain) != "private human meaning" {
		t.Fatalf("plaintext = %q", plain)
	}
}

func TestSealDraftRejectsBadKey(t *testing.T) {
	if _, err := SealDraft([]byte("x"), []byte("short")); err == nil {
		t.Fatal("expected short key to fail")
	}
}

func TestOpenDraftRejectsWrongKey(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	sealed, err := SealDraft([]byte("private"), key)
	if err != nil {
		t.Fatal(err)
	}
	wrong := []byte("abcdef0123456789abcdef0123456789")
	if _, err := OpenDraft(sealed, wrong); err == nil {
		t.Fatal("expected authentication failure")
	}
}
