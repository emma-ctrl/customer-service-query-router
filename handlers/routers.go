package handlers

import (
    "encoding/json"
    "net/http"
    "customer-query-router/models"
    "customer-query-router/services"
)

type RouterHandler struct {
    agentService          *services.AgentService
    conversationService   *services.ConversationService
    classificationService *services.ClassificationService
}

func NewRouterHandler(agentService *services.AgentService, conversationService *services.ConversationService, classificationService *services.ClassificationService) *RouterHandler {
    return &RouterHandler{
        agentService:          agentService,
        conversationService:   conversationService,
        classificationService: classificationService,
    }
}

// ClassifyQuery handles OpenAI-based query classification
func (rh *RouterHandler) ClassifyQuery(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        CustomerMessage string `json:"customer_message"`
    }

    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if request.CustomerMessage == "" {
        errorResponse := map[string]string{
            "error": "customer_message field is required",
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    intent, agent, err := rh.classificationService.ClassifyQuery(request.CustomerMessage)
    if err != nil {
        errorResponse := map[string]string{
            "error": "Classification failed: " + err.Error(),
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    response := map[string]interface{}{
        "intent":           intent,
        "recommended_agent": agent,
        "message":          request.CustomerMessage,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (rh *RouterHandler) RouteQuery(w http.ResponseWriter, r *http.Request) {
    var query models.Query
    err := json.NewDecoder(r.Body).Decode(&query)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    agent, err := rh.agentService.FindAvailableAgent(query.Intent)
    if err != nil {
        errorResponse := map[string]string{
            "error": err.Error(),
            "status": "no_agent_available",
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }
    
    rh.agentService.AssignQuery(agent.ID)
    
    response := models.RoutingResponse{
        AgentID: agent.ID,
        Intent:  query.Intent,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// GetAgents returns all agents and their status
func (rh *RouterHandler) GetAgents(w http.ResponseWriter, r *http.Request) {
    agents := rh.agentService.GetAllAgents()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(agents)
}

// GetAgentStats returns system statistics
func (rh *RouterHandler) GetAgentStats(w http.ResponseWriter, r *http.Request) {
    stats := rh.agentService.GetAgentStats()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

func (rh *RouterHandler) TestConversations(w http.ResponseWriter, r *http.Request) {
    // This is just to test our file reading
    response := map[string]interface{}{
        "message": "Conversation loading test",
        "total_conversations": len(rh.conversationService.GetConversations()),
    }
    
    // Show first few customer messages as examples
    conversations := rh.conversationService.GetConversations()
    examples := []string{}
    
    for i, conv := range conversations {
        if i >= 3 { // Just show first 3
            break
        }
        firstMessage := rh.conversationService.GetFirstCustomerMessage(conv)
        examples = append(examples, firstMessage)
    }
    
    response["examples"] = examples
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// TestClassificationOnConversations tests the OpenAI classification on your loaded conversations
func (rh *RouterHandler) TestClassificationOnConversations(w http.ResponseWriter, r *http.Request) {
    conversations := rh.conversationService.GetConversations()
    if len(conversations) == 0 {
        errorResponse := map[string]string{
            "error": "No conversations loaded",
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    results := []map[string]interface{}{}
    
    // Test classification on first 5 conversations
    limit := 5
    if len(conversations) < limit {
        limit = len(conversations)
    }

    for i := 0; i < limit; i++ {
        conv := conversations[i]
        firstMessage := rh.conversationService.GetFirstCustomerMessage(conv)
        
        intent, agent, err := rh.classificationService.ClassifyQuery(firstMessage)
        
        result := map[string]interface{}{
            "conversation_id": i + 1,
            "customer_message": firstMessage,
            "classified_intent": intent,
            "recommended_agent": agent,
        }
        
        if err != nil {
            result["error"] = err.Error()
        }
        
        results = append(results, result)
    }

    response := map[string]interface{}{
        "message": "Classification test on conversations",
        "total_tested": limit,
        "results": results,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}