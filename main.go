package main

import (
    "fmt"
    "log"
    "net/http"
    "customer-query-router/handlers"
    "customer-query-router/services"
)

func main() {
    // Initialize services
    agentService := services.NewAgentService()
    conversationService := services.NewConversationService()
    
    // Load conversations
    err := conversationService.LoadConversations("data/conversations.txt")
    if err != nil {
        log.Fatal("Failed to load conversations:", err)
    }
    
    // Initialize handlers
    routerHandler := handlers.NewRouterHandler(agentService, conversationService)
    
    // Set up routes
    http.HandleFunc("/route", routerHandler.RouteQuery)
    http.HandleFunc("/agents", routerHandler.GetAgents)
    http.HandleFunc("/agents/stats", routerHandler.GetAgentStats)
    http.HandleFunc("/test-conversations", routerHandler.TestConversations)
    
    fmt.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}