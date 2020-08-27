package shortener

import (
	"bytes"
	"strings"
	"testing"
)

func help_test(t *testing.T) {
	var buf bytes.Buffer
	printHelp(&buf)

	out := buf.String()
	if strings.Contains(out, "-e") {
		t.Errorf("Help string lack of instruction: got=%s", out)
	}

	if !strings.HasSuffix(out, "\n") {
		t.Errorf("Doesn't end with new line: got=%s", out)
	}

}
