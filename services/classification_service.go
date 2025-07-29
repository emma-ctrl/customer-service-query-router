package services

import (
    "context"
    "fmt"
    "strings"

    openai "github.com/sashabaranov/go-openai"
)

type ClassificationService struct {
    client *openai.Client
    intents []Intent
}

type Intent struct {
    Name  string
    Agent string
}

func NewClassificationService(apiKey string) *ClassificationService {
    intents := []Intent{
        {"account_access_issues", "account-support"},
        {"billing_discrepancies", "billing-team"},
        {"delivery_problems", "logistics-team"},
        {"installation_support_requests", "technical-support"},
        {"order_cancellation_requests", "order-management"},
        {"order_status_uncertainty", "order-tracking"},
        {"product_availability_inquiries", "inventory-team"},
        {"refund_processing_issues", "finance-team"},
        {"return_process_inquiries", "returns-team"},
        {"warranty_terms_inquiries", "warranty-team"},
        {"general", "general-agent"}, // Fallback
    }

    return &ClassificationService{
        client:  openai.NewClient(apiKey),
        intents: intents,
    }
}

func (cs *ClassificationService) ClassifyQuery(customerMessage string) (string, string, error) {
    prompt := cs.buildClassificationPrompt()
    
    resp, err := cs.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleSystem,
                    Content: prompt,
                },
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: customerMessage,
                },
            },
            MaxTokens:   50,
            Temperature: 0.1, // Low temperature for consistent classification
        },
    )

    if err != nil {
        return "", "", fmt.Errorf("OpenAI API error: %w", err)
    }

    intent := strings.TrimSpace(resp.Choices[0].Message.Content)
    agent := cs.getAgentForIntent(intent)
    
    return intent, agent, nil
}

func (cs *ClassificationService) buildClassificationPrompt() string {
    intentList := make([]string, len(cs.intents)-1) // Exclude "general" from the list
    for i, intent := range cs.intents[:len(cs.intents)-1] {
        intentList[i] = intent.Name
    }
    
    return fmt.Sprintf(`You are a customer service query classifier. 
Classify the following customer message into exactly ONE of these specific intents: %s

If the message doesn't clearly fit into any of these specific categories, respond with "general".

Respond with only the intent name, nothing else.`, strings.Join(intentList, ", "))
}

func (cs *ClassificationService) getAgentForIntent(intent string) string {
    for _, i := range cs.intents {
        if i.Name == intent {
            return i.Agent
        }
    }
    return "general-agent" // Default fallback
}

func (cs *ClassificationService) GetAllIntents() []Intent {
    return cs.intents
}