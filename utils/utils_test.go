package utils

import "testing"

func TestParseInt(t *testing.T) {
	if got := ParseInt(" 42 ", -1); got != 42 {
		t.Fatalf("ParseInt() = %d, want 42", got)
	}
	if got := ParseInt("bad", -1); got != -1 {
		t.Fatalf("ParseInt() fallback = %d, want -1", got)
	}
}

func TestParseInt64(t *testing.T) {
	if got := ParseInt64(" 42 ", -1); got != 42 {
		t.Fatalf("ParseInt64() = %d, want 42", got)
	}
	if got := ParseInt64("bad", -1); got != -1 {
		t.Fatalf("ParseInt64() fallback = %d, want -1", got)
	}
}

func TestParseUint(t *testing.T) {
	if got := ParseUint(" 42 ", 7); got != 42 {
		t.Fatalf("ParseUint() = %d, want 42", got)
	}
	if got := ParseUint("-1", 7); got != 7 {
		t.Fatalf("ParseUint() fallback = %d, want 7", got)
	}
}

func TestParseUint64(t *testing.T) {
	if got := ParseUint64(" 42 ", 7); got != 42 {
		t.Fatalf("ParseUint64() = %d, want 42", got)
	}
	if got := ParseUint64("-1", 7); got != 7 {
		t.Fatalf("ParseUint64() fallback = %d, want 7", got)
	}
}

func TestParseFloat64(t *testing.T) {
	if got := ParseFloat64(" 3.14 ", -1); got != 3.14 {
		t.Fatalf("ParseFloat64() = %f, want 3.14", got)
	}
	if got := ParseFloat64("bad", -1); got != -1 {
		t.Fatalf("ParseFloat64() fallback = %f, want -1", got)
	}
}

func TestParseBool(t *testing.T) {
	if got := ParseBool(" true ", false); !got {
		t.Fatalf("ParseBool() = %t, want true", got)
	}
	if got := ParseBool("bad", true); !got {
		t.Fatalf("ParseBool() fallback = %t, want true", got)
	}
}
