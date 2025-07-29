package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

// 11 hardcoded categories (for now)
var intentCategories = []string{
    "account_access_issues",
    "product_quality_concerns", 
    "order_status_uncertainty",
    "product_availability_inquiries",
    "billing_discrepancies",
    "delivery_problems",
    "return_process_inquiries",
    "warranty_terms_inquiries", 
    "order_cancellation_requests",
    "refund_processing_issues",
    "installation_support_requests",
}

//5. this Go struct
type Query struct {
    Content string `json:"content"`
    Intent  string `json:"intent"`
}

type RoutingResponse struct {
    AgentID string `json:"agent_id"`
    Intent  string `json:"intent"`
}

//3. Go matched the /route url and then called here routeQuery

func routeQuery(w http.ResponseWriter, r *http.Request) {
    var query Query
    json.NewDecoder(r.Body).Decode(&query) // 4. reads the json from curl command + converts content into Go struct
    
    // 6. processed the request using simple mapping
    agentID := "agent-" + query.Intent
    
    response := RoutingResponse{
        AgentID: agentID,
        Intent:  query.Intent,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// this runs when server starts up
// route query function is registered to handle requests to /route
// 1. server is ready and waiting for HTTP requests
func main() {
    http.HandleFunc("/route", routeQuery)
    
    fmt.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

//2. send a curl request which contains URL, header, body
/*

curl -X POST http://localhost:8080/route \
  -H "Content-Type: application/json" \
  -d '{"content": "I cant log into my account", "intent": "account_access_issues"}'

  7. output

  {"agent_id":"agent-account_access_issues","intent":"account_access_issues"}

*/