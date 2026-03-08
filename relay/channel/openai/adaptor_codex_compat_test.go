package openai

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestConvertOpenAIResponsesRequest_AddsCodexCompatDefaults(t *testing.T) {
	adaptor := &Adaptor{}

	request := dto.OpenAIResponsesRequest{
		Model: "gpt-5.4",
	}
	convertedAny, err := adaptor.ConvertOpenAIResponsesRequest(nil, nil, request)
	require.NoError(t, err)

	converted, ok := convertedAny.(dto.OpenAIResponsesRequest)
	require.True(t, ok)
	require.Equal(t, `""`, string(converted.Instructions))
	require.Equal(t, "false", string(converted.Store))
}

func TestConvertOpenAIResponsesRequest_NoCompatDefaultsForNormalModel(t *testing.T) {
	adaptor := &Adaptor{}

	request := dto.OpenAIResponsesRequest{
		Model: "gpt-4o-mini",
	}
	convertedAny, err := adaptor.ConvertOpenAIResponsesRequest(nil, nil, request)
	require.NoError(t, err)

	converted, ok := convertedAny.(dto.OpenAIResponsesRequest)
	require.True(t, ok)
	require.Len(t, converted.Instructions, 0)
	require.Len(t, converted.Store, 0)
}

func TestSetupRequestHeader_AddsCodexCompatHeadersForResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	headers := make(http.Header)
	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType:       constant.ChannelTypeOpenAI,
			ApiKey:            "sk-test",
			UpstreamModelName: "gpt-5.4",
		},
		RelayMode: relayconstant.RelayModeResponses,
	}

	adaptor := &Adaptor{}
	err := adaptor.SetupRequestHeader(c, &headers, info)
	require.NoError(t, err)
	require.Equal(t, "Bearer sk-test", headers.Get("Authorization"))
	require.Equal(t, "responses=experimental", headers.Get("OpenAI-Beta"))
	require.Equal(t, "codex_cli_rs", headers.Get("originator"))
}

func TestSetupRequestHeader_RespectsExistingCodexCompatHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	headers := make(http.Header)
	headers.Set("OpenAI-Beta", "custom")
	headers.Set("originator", "custom-originator")

	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType:       constant.ChannelTypeOpenAI,
			ApiKey:            "sk-test",
			UpstreamModelName: "gpt-5.4",
		},
		RelayMode: relayconstant.RelayModeResponses,
	}

	adaptor := &Adaptor{}
	err := adaptor.SetupRequestHeader(c, &headers, info)
	require.NoError(t, err)
	require.Equal(t, "custom", headers.Get("OpenAI-Beta"))
	require.Equal(t, "custom-originator", headers.Get("originator"))
}

func TestShouldApplyCodexResponsesCompat(t *testing.T) {
	require.True(t, shouldApplyCodexResponsesCompat("gpt-5"))
	require.True(t, shouldApplyCodexResponsesCompat("gpt-5.4"))
	require.True(t, shouldApplyCodexResponsesCompat("gpt-5.1-codex-mini"))
	require.False(t, shouldApplyCodexResponsesCompat("gpt-4o-mini"))
}
