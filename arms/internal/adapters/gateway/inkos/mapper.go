package inkos

import (
	"encoding/json"
	"fmt"
	"strings"
)

const maxContextRunes = 12000

// BuildWriteNextArgs returns argv tail for `inkos write next … --json` (excluding the binary).
func BuildWriteNextArgs(bookID, context string) []string {
	args := []string{"write", "next", bookID, "--count", "1", "--json"}
	ctx := truncateContext(strings.TrimSpace(context), maxContextRunes)
	if ctx != "" {
		args = append(args, "--context", ctx)
	}
	return args
}

func truncateContext(s string, maxRunes int) string {
	if maxRunes <= 0 || s == "" {
		return ""
	}
	r := []rune(s)
	if len(r) <= maxRunes {
		return s
	}
	return string(r[:maxRunes])
}

// ExternalRefFromInkOSStdout picks a stable segment from InkOS --json stdout, or returns empty.
func ExternalRefFromInkOSStdout(stdout []byte) string {
	s := strings.TrimSpace(string(stdout))
	if s == "" {
		return ""
	}
	var top map[string]interface{}
	if err := json.Unmarshal([]byte(s), &top); err != nil {
		return ""
	}
	if seg := pickRefSegment(top); seg != "" {
		return seg
	}
	if raw, ok := top["data"]; ok {
		if nested, ok := raw.(map[string]interface{}); ok {
			return pickRefSegment(nested)
		}
	}
	return ""
}

func pickRefSegment(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	order := []string{"requestId", "jobId", "traceId", "id", "taskId", "runId"}
	for _, k := range order {
		if v, ok := m[k]; ok {
			out := strings.TrimSpace(fmt.Sprint(v))
			if out != "" {
				return out
			}
		}
	}
	if v, ok := m["chapter"]; ok {
		return fmt.Sprintf("ch-%v", v)
	}
	if v, ok := m["bookId"]; ok {
		out := strings.TrimSpace(fmt.Sprint(v))
		if out != "" {
			return "book-" + out
		}
	}
	return ""
}
