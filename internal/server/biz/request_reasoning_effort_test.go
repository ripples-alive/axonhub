package biz

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/looplj/axonhub/internal/contexts"
	"github.com/looplj/axonhub/internal/ent/project"
	"github.com/looplj/axonhub/llm"
	"github.com/looplj/axonhub/llm/httpclient"
)

func TestRequestServiceCreateRequestPersistsReasoningEffort(t *testing.T) {
	svc, client, ctx := setupTestRequestService(t)
	defer client.Close()

	proj, err := client.Project.Create().
		SetName("test-project").
		SetStatus(project.StatusActive).
		Save(ctx)
	require.NoError(t, err)

	ctx = contexts.WithProjectID(ctx, proj.ID)
	stream := true
	req, err := svc.CreateRequest(ctx, &llm.Request{
		Model:           "gpt-5.5",
		ReasoningEffort: "high",
		Stream:          &stream,
	}, &httpclient.Request{
		JSONBody: []byte(`{"model":"gpt-5.5","reasoning_effort":"high","stream":true}`),
	}, llm.APIFormatOpenAIChatCompletion)
	require.NoError(t, err)

	saved, err := client.Request.Get(ctx, req.ID)
	require.NoError(t, err)
	require.Equal(t, "high", saved.ReasoningEffort)
}
