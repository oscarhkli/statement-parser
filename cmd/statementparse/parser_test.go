package statementparse

import "testing"

func TestParseEmptyString(t *testing.T) {
	res := Parse("")
	if len(res) > 0 {
		t.Errorf("Parse(\"\") = %v; want []", res)
	}
}
