package handlers

import (
    "encoding/json"
    "net/http"
    "customer-query-router/models"
    "customer-query-router/services"
)

type RouterHandler struct {
    agentService *services.AgentService
}

func NewRouterHandler(agentService *services.AgentService) *RouterHandler {
    return &RouterHandler{
        agentService: agentService,
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