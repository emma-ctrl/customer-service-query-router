package services

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type ConversationService struct {
    conversations []string
}

func NewConversationService() *ConversationService {
    return &ConversationService{}
}

func (cs *ConversationService) LoadConversations(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    var currentConversation strings.Builder
    var conversations []string
    
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        line := scanner.Text()
        
        // If line starts with quote and "Agent: Thank you", it's a new conversation
        if strings.HasPrefix(line, `"Agent: Thank you for calling`) {
            // Save previous conversation if exists
            if currentConversation.Len() > 0 {
                conversations = append(conversations, strings.TrimSpace(currentConversation.String()))
                currentConversation.Reset()
            }
        }
        
        currentConversation.WriteString(line)
        currentConversation.WriteString("\n")
    }
    
    // Don't forget the last conversation
    if currentConversation.Len() > 0 {
        conversations = append(conversations, strings.TrimSpace(currentConversation.String()))
    }
    
    cs.conversations = conversations
    fmt.Printf("Loaded %d conversations\n", len(conversations))
    
    return scanner.Err()
}

func (cs *ConversationService) GetConversations() []string {
    return cs.conversations
}

func (cs *ConversationService) GetFirstCustomerMessage(conversation string) string {
    lines := strings.Split(conversation, "\n")
    
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "Customer: ") {
            // Remove "Customer: " prefix and return the message
            return strings.TrimPrefix(line, "Customer: ")
        }
    }
    
    return "No customer message found"
}