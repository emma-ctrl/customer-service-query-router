package handlers

import (
    "encoding/json"
    "net/http"
    "customer-query-router/models"
    "customer-query-router/services"
)

type RouterHandler struct {
    agentService *services.AgentService
	conversationService *services.ConversationService 
}

func NewRouterHandler(agentService *services.AgentService, conversationService *services.ConversationService) *RouterHandler {
    return &RouterHandler{
        agentService: agentService,
		conversationService: conversationService,
    }
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