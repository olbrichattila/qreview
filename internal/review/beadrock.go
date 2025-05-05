package review

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/olbrichattila/qreview/internal/report"
	"github.com/olbrichattila/qreview/internal/retriever"
)

// newBedrock creates a new AWS bedrock reviewer using the AWS SDK
func newBedrock(
	retr retriever.Retriever,
	prompt string,
	reporters []report.Reporter,
	commentOnPR bool,
) Reviewer {
	return &bedrock{
		reporters:   reporters,
		retr:        retr,
		prompt:      prompt,
		commentOnPR: commentOnPR,
	}
}

type bedrock struct {
	reporters   []report.Reporter
	retr        retriever.Retriever
	prompt      string
	commentOnPR bool
}

// Claude message structure
type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Claude request structure
type claudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version"`
	MaxTokens        int             `json:"max_tokens"`
	Messages         []claudeMessage `json:"messages"`
}

// Claude response structure
type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
}

// AnalyzeCode returns the result of the analyzed code using Amazon Q via AWS SDK
func (a *bedrock) AnalyzeCode(fileName string) error {
	content, err := a.retr.Get(fileName)
	if err != nil {
		return fmt.Errorf("Analyze code %w", err)
	}

	// Load AWS configuration from environment variables or shared credentials file
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create Amazon Bedrock Runtime client (which is used for Amazon Q)
	bedrockClient := bedrockruntime.NewFromConfig(cfg)

	// Prepare the message with the prompt and file content
	message := a.prompt + content.FileContent

	// Create the Claude request (Amazon Q uses Claude under the hood)
	claudeReq := claudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        4096,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	// Convert the request to JSON
	reqBody, err := json.Marshal(claudeReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Set a timeout for the API call
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Call the Amazon Bedrock API with Claude model (which powers Amazon Q)
	// You can use different model IDs based on your needs
	// modelID := "anthropic.claude-3-sonnet-20240229-v1:0"
	modelID := "anthropic.claude-v2"

	invokeResp, err := bedrockClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
		Body:        reqBody,
	})

	if err != nil {
		return fmt.Errorf("failed to get response from Amazon Bedrock: %w", err)
	}

	// Parse the response
	var claudeResp claudeResponse
	err = json.Unmarshal(invokeResp.Body, &claudeResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the response text
	aiResponse := ""
	for _, content := range claudeResp.Content {
		if content.Type == "text" {
			aiResponse += content.Text
		}
	}

	// Handle PR comments if needed
	if a.commentOnPR {
		err = commentOnPRIfNecessary(fileName, aiResponse, content.DiffContent)
		if err != nil {
			return err
		}
	}

	return generateReports(a.reporters, fileName, aiResponse)
}

// Summary implements Reviewer.
func (a *bedrock) Summary() error {
	return summary(a.reporters)
}
