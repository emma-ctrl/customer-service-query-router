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
    
    // Initialize handlers
    routerHandler := handlers.NewRouterHandler(agentService)
    
    // Set up routes
    http.HandleFunc("/route", routerHandler.RouteQuery)
    http.HandleFunc("/agents", routerHandler.GetAgents)          // New!
    http.HandleFunc("/agents/stats", routerHandler.GetAgentStats) // New!
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Endpoints:")
    fmt.Println("  POST /route - Route a customer query")
    fmt.Println("  GET  /agents - View all agents")
    fmt.Println("  GET  /agents/stats - View system stats")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}