package services

import (
    "context"
    "fmt"
    "log"
    "strings"
    "time"

    openai "github.com/sashabaranov/go-openai"
)

type ClassificationService struct {
    client *openai.Client
    intents []Intent
    requestCount int64
    totalProcessingTime time.Duration
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

    log.Printf("[CLASSIFICATION SERVICE] Initialized with %d intent categories", len(intents)-1)
    
    // Create the service instance first
    service := &ClassificationService{
        client:  openai.NewClient(apiKey),
        intents: intents,
        requestCount: 0,
        totalProcessingTime: 0,
    }
    
    log.Printf("[CLASSIFICATION SERVICE] Available intents: %s", service.getIntentNames(intents))
    
    return service
    
    return &ClassificationService{
        client:  openai.NewClient(apiKey),
        intents: intents,
        requestCount: 0,
        totalProcessingTime: 0,
    }
}

func (cs *ClassificationService) ClassifyQuery(customerMessage string) (string, string, error) {
    startTime := time.Now()
    cs.requestCount++
    requestID := cs.requestCount
    
    log.Printf("[REQUEST %d] Starting classification process", requestID)
    log.Printf("[REQUEST %d] STEP 1 - Message received: \"%s\"", requestID, truncateMessage(customerMessage, 100))
    log.Printf("[REQUEST %d] Message length: %d characters", requestID, len(customerMessage))
    
    // Build the classification prompt
    prompt := cs.buildClassificationPrompt()
    log.Printf("[REQUEST %d] STEP 2 - Built classification prompt (%d characters)", requestID, len(prompt))
    
    // Log the actual prompt being sent (truncated for readability)
    log.Printf("[REQUEST %d] Prompt preview: \"%s\"", requestID, truncateMessage(prompt, 150))
    
    log.Printf("[REQUEST %d] STEP 3 - Sending request to OpenAI GPT-3.5-turbo", requestID)
    
    // Create the OpenAI request
    openaiRequest := openai.ChatCompletionRequest{
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
        Temperature: 0.1,
    }
    
    log.Printf("[REQUEST %d] OpenAI request configured - Model: %s, MaxTokens: %d, Temperature: %.1f", 
        requestID, openaiRequest.Model, openaiRequest.MaxTokens, openaiRequest.Temperature)
    
    // Make the API call
    apiStartTime := time.Now()
    resp, err := cs.client.CreateChatCompletion(context.Background(), openaiRequest)
    apiDuration := time.Since(apiStartTime)
    
    if err != nil {
        log.Printf("[REQUEST %d] ERROR - OpenAI API call failed: %v", requestID, err)
        log.Printf("[REQUEST %d] API call duration: %v", requestID, apiDuration)
        return "", "", fmt.Errorf("OpenAI API error: %w", err)
    }
    
    log.Printf("[REQUEST %d] STEP 4 - Received response from OpenAI in %v", requestID, apiDuration)
    log.Printf("[REQUEST %d] Response tokens used: %d", requestID, resp.Usage.TotalTokens)
    
    // Extract and process the classification result
    rawIntent := resp.Choices[0].Message.Content
    intent := strings.TrimSpace(rawIntent)
    
    log.Printf("[REQUEST %d] Raw OpenAI response: \"%s\"", requestID, rawIntent)
    log.Printf("[REQUEST %d] STEP 5 - Processed intent: \"%s\"", requestID, intent)
    
    // Validate the intent
    isValidIntent := cs.isValidIntent(intent)
    if !isValidIntent {
        log.Printf("[REQUEST %d] WARNING - Unrecognized intent \"%s\", falling back to \"general\"", requestID, intent)
        intent = "general"
    }
    
    // Get the assigned agent
    agent := cs.getAgentForIntent(intent)
    log.Printf("[REQUEST %d] STEP 6 - Agent assignment: \"%s\" -> \"%s\"", requestID, intent, agent)
    
    // Calculate final metrics
    totalDuration := time.Since(startTime)
    cs.totalProcessingTime += totalDuration
    avgProcessingTime := cs.totalProcessingTime / time.Duration(cs.requestCount)
    
    log.Printf("[REQUEST %d] STEP 7 - Classification complete", requestID)
    log.Printf("[REQUEST %d] METRICS - Total time: %v, API time: %v, Processing time: %v", 
        requestID, totalDuration, apiDuration, totalDuration-apiDuration)
    log.Printf("[REQUEST %d] METRICS - Request #%d, Average processing time: %v", 
        requestID, cs.requestCount, avgProcessingTime)
    
    log.Printf("[REQUEST %d] FINAL RESULT - Intent: \"%s\", Agent: \"%s\"", requestID, intent, agent)
    log.Printf("================================================================================")
    
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

func (cs *ClassificationService) isValidIntent(intent string) bool {
    for _, i := range cs.intents {
        if i.Name == intent {
            return true
        }
    }
    return false
}

func (cs *ClassificationService) GetAllIntents() []Intent {
    return cs.intents
}

func (cs *ClassificationService) GetStats() map[string]interface{} {
    avgProcessingTime := time.Duration(0)
    if cs.requestCount > 0 {
        avgProcessingTime = cs.totalProcessingTime / time.Duration(cs.requestCount)
    }
    
    return map[string]interface{}{
        "total_requests": cs.requestCount,
        "total_processing_time": cs.totalProcessingTime.String(),
        "average_processing_time": avgProcessingTime.String(),
    }
}

// Helper function to get all intent names for logging
func (cs *ClassificationService) getIntentNames(intents []Intent) string {
    names := make([]string, len(intents))
    for i, intent := range intents {
        names[i] = intent.Name
    }
    return strings.Join(names, ", ")
}

// Helper function to truncate long messages for logging
func truncateMessage(message string, maxLength int) string {
    if len(message) <= maxLength {
        return message
    }
    return message[:maxLength] + "..."
}