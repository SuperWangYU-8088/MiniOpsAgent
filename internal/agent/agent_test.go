package agent

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/itwanger/paicli-go/internal/llm"
)

func TestFormatToolCallPrettyJSON(t *testing.T) {
	args := json.RawMessage(`{"path":"README.md","limit":5}`)
	got := formatToolCall(llm.ToolCall{
		ID: "call_1",
		Function: llm.FunctionCall{
			Name:      "read_file",
			Arguments: args,
		},
	})

	if !strings.Contains(got, `"path": "README.md"`) {
		t.Fatalf("formatted tool call missing path: %q", got)
	}
	if !strings.Contains(got, `"limit": 5`) {
		t.Fatalf("formatted tool call missing limit: %q", got)
	}
}

func TestRepeatedToolCall(t *testing.T) {
	call := llm.ToolCall{
		ID: "call_1",
		Function: llm.FunctionCall{
			Name:      "web_search",
			Arguments: json.RawMessage(`{"query":"沉默王二"}`),
		},
	}
	seen := map[string]int{}
	if repeatedToolCall([]llm.ToolCall{call}, seen) {
		t.Fatal("first tool call should not be repeated")
	}
	if !repeatedToolCall([]llm.ToolCall{call}, seen) {
		t.Fatal("second identical tool call should be repeated")
	}
}
