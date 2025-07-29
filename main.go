package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "customer-query-router/handlers"
    "customer-query-router/services"
)

func main() {
    // Get OpenAI API key
    openaiKey := os.Getenv("OPENAI_API_KEY")
    if openaiKey == "" {
        log.Fatal("OPENAI_API_KEY environment variable is required")
    }

    // Initialize services
    agentService := services.NewAgentService()
    conversationService := services.NewConversationService()
    classificationService := services.NewClassificationService(openaiKey)
    
    // Load conversations
    err := conversationService.LoadConversations("data/conversations.txt")
    if err != nil {
        log.Fatal("Failed to load conversations:", err)
    }
    
    // Initialize handlers
    routerHandler := handlers.NewRouterHandler(agentService, conversationService, classificationService)
    uiHandler := handlers.NewUIHandler()
    
    // Set up UI routes
    http.HandleFunc("/", uiHandler.ServeHome)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
    
    // Set up API routes with CORS
    http.HandleFunc("/api/route", handlers.EnableCORS(routerHandler.RouteQuery))
    http.HandleFunc("/api/agents", handlers.EnableCORS(routerHandler.GetAgents))
    http.HandleFunc("/api/agents/stats", handlers.EnableCORS(routerHandler.GetAgentStats))
    http.HandleFunc("/api/test-conversations", handlers.EnableCORS(routerHandler.TestConversations))
    http.HandleFunc("/api/classify", handlers.EnableCORS(routerHandler.ClassifyQuery))
    http.HandleFunc("/api/test-classification", handlers.EnableCORS(routerHandler.TestClassificationOnConversations))
    
    fmt.Println("Server starting on :8080")
    fmt.Println("🌐 Web UI: http://localhost:8080")
    fmt.Println("\nAPI Endpoints:")
    fmt.Println("POST /api/classify - Classify customer queries with OpenAI")
    fmt.Println("POST /api/route - Route customer queries")  
    fmt.Println("GET  /api/agents - Get all agents")
    fmt.Println("GET  /api/agents/stats - Get agent statistics")
    fmt.Println("POST /api/test-conversations - Test conversations")
    fmt.Println("POST /api/test-classification - Test OpenAI classification on loaded conversations")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}