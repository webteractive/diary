package hash

import "testing"

func TestContentIsDeterministicAndExcludesHash(t *testing.T) {
	fields := map[string]any{
		"id":      "record-1",
		"project": "diary",
		"hash":    "sha256:ignored",
	}

	first, err := Content(fields, "hello\r\nworld")
	if err != nil {
		t.Fatal(err)
	}

	fields["hash"] = "sha256:different"
	second, err := Content(fields, "hello\nworld")
	if err != nil {
		t.Fatal(err)
	}

	if first != second {
		t.Fatalf("expected stable hash, got %q and %q", first, second)
	}
}

func TestContentChangesWhenBodyChanges(t *testing.T) {
	fields := map[string]any{"id": "record-1"}

	first, err := Content(fields, "hello")
	if err != nil {
		t.Fatal(err)
	}
	second, err := Content(fields, "goodbye")
	if err != nil {
		t.Fatal(err)
	}

	if first == second {
		t.Fatal("expected body change to change hash")
	}
}
