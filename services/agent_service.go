package services

import (
    "fmt"
    "customer-query-router/models"
)

type AgentService struct {
    agents map[string]*models.Agent
}

func NewAgentService() *AgentService {
    return &AgentService{
        agents: initializeAgents(),
    }
}

func (as *AgentService) FindAvailableAgent(intent string) (*models.Agent, error) {
    for _, agent := range as.agents {
        canHandle := false
        for _, specialty := range agent.Specialties {
            if specialty == intent {
                canHandle = true
                break
            }
        }
        
        if canHandle && agent.IsOnline && agent.CurrentLoad < agent.MaxCapacity {
            return agent, nil
        }
    }
    
    return nil, fmt.Errorf("no available agent for intent: %s", intent)
}

func (as *AgentService) AssignQuery(agentID string) {
    if agent, exists := as.agents[agentID]; exists {
        agent.CurrentLoad++
    }
}

func initializeAgents() map[string]*models.Agent {
    return map[string]*models.Agent{
        "billing-specialist": {
            ID:          "billing-specialist",
            Name:        "Sarah - Billing Expert",
            Specialties: []string{"billing_discrepancies", "refund_processing_issues"},
            MaxCapacity: 5,
            CurrentLoad: 2,
            IsOnline:    true,
        },
        // ... other agents
    }
}