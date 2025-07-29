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
    
    fmt.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}