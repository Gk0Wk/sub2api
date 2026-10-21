package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func TestNormalizeOpenAIResponsesCompactRequest_PreservesNativeCompactionUntilRouting(t *testing.T) {
	h := &OpenAIGatewayHandler{}
	body := []byte(`{
		"model":"gpt-6.1-sol",
		"stream":true,
		"input":[
			{"type":"message","role":"user","content":"hello"},
			{"type":"reasoning","id":"rs_1","encrypted_content":"account-a-ciphertext"},
			{"type":"compaction","encrypted_content":"account-a-compaction"},
			{"type":"message","role":"user","content":"continue"},
			{"type":"compaction_trigger"}
		]
	}`)
	c := newCompactBodySignalTestContext(t, "/v1/responses", body)

	normalized, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
	require.True(t, ok)
	require.True(t, isOpenAIRemoteCompactionV2Request(normalized))
	require.Equal(t, "account-a-ciphertext", gjson.GetBytes(normalized, `input.#(type=="reasoning").encrypted_content`).String())
	require.Equal(t, "account-a-compaction", gjson.GetBytes(normalized, `input.#(type=="compaction").encrypted_content`).String())
	require.Equal(t, "compaction_trigger", gjson.GetBytes(normalized, "input.4.type").String())
	require.Equal(t, "continue", gjson.GetBytes(normalized, "input.3.content").String())
}
