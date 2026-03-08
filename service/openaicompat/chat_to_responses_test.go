package openaicompat

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/stretchr/testify/require"
)

func TestChatCompletionsRequestToResponsesRequest_CodexStyleStringContent(t *testing.T) {
	req := &dto.GeneralOpenAIRequest{
		Model: "gpt-5.4",
		Messages: []dto.Message{
			{
				Role:    "user",
				Content: "hello",
			},
		},
	}

	out, err := ChatCompletionsRequestToResponsesRequest(req)
	require.NoError(t, err)

	var inputItems []map[string]any
	require.NoError(t, common.Unmarshal(out.Input, &inputItems))
	require.Len(t, inputItems, 1)

	content, ok := inputItems[0]["content"].([]any)
	require.True(t, ok, "codex-compatible model should use array content")
	require.Len(t, content, 1)

	part, ok := content[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "input_text", part["type"])
	require.Equal(t, "hello", part["text"])
}

func TestChatCompletionsRequestToResponsesRequest_DefaultStringContent(t *testing.T) {
	req := &dto.GeneralOpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []dto.Message{
			{
				Role:    "user",
				Content: "hello",
			},
		},
	}

	out, err := ChatCompletionsRequestToResponsesRequest(req)
	require.NoError(t, err)

	var inputItems []map[string]any
	require.NoError(t, common.Unmarshal(out.Input, &inputItems))
	require.Len(t, inputItems, 1)
	require.Equal(t, "hello", inputItems[0]["content"])
}
