package shortener

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	var buf bytes.Buffer
	printHelp(&buf)

	out := buf.String()
	if !strings.Contains(out, "-e") {
		t.Errorf("Help string lack of instruction: got=%s", out)
	}

	if !strings.HasPrefix(out, "USAGE") {
		t.Errorf("Doesn't start with USAGE: got=%s", out)
	}

	if !strings.HasSuffix(out, "\n") {
		t.Errorf("Doesn't end with new line: got=%s", out)
	}

}
