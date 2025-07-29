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
        "account-helper": {
            ID:          "account-helper",
            Name:        "Mike - Account Support",
            Specialties: []string{"account_access_issues"},
            MaxCapacity: 3,
            CurrentLoad: 3, // At capacity!
            IsOnline:    true,
        },
        "delivery-tracker": {
            ID:          "delivery-tracker",
            Name:        "Emma - Delivery Support",
            Specialties: []string{"delivery_problems", "order_status_uncertainty"},
            MaxCapacity: 4,
            CurrentLoad: 1,
            IsOnline:    true,
        },
        "product-expert": {
            ID:          "product-expert",
            Name:        "Alex - Product Specialist",
            Specialties: []string{"product_quality_concerns", "product_availability_inquiries"},
            MaxCapacity: 6,
            CurrentLoad: 0,
            IsOnline:    true,
        },
        "returns-processor": {
            ID:          "returns-processor",
            Name:        "Jordan - Returns & Exchanges",
            Specialties: []string{"return_process_inquiries", "order_cancellation_requests"},
            MaxCapacity: 4,
            CurrentLoad: 2,
            IsOnline:    true,
        },
        "warranty-advisor": {
            ID:          "warranty-advisor",
            Name:        "Taylor - Warranty Support",
            Specialties: []string{"warranty_terms_inquiries"},
            MaxCapacity: 3,
            CurrentLoad: 1,
            IsOnline:    true,
        },
        "tech-support": {
            ID:          "tech-support",
            Name:        "Casey - Technical Support",
            Specialties: []string{"installation_support_requests"},
            MaxCapacity: 5,
            CurrentLoad: 0,
            IsOnline:    false, // Offline for maintenance
        },
    }
}

// GetAllAgents returns all agents and their current status
func (as *AgentService) GetAllAgents() map[string]*models.Agent {
    return as.agents
}

// GetAgent returns a specific agent by ID
func (as *AgentService) GetAgent(agentID string) (*models.Agent, bool) {
    agent, exists := as.agents[agentID]
    return agent, exists
}

// GetAgentStats returns summary statistics
func (as *AgentService) GetAgentStats() map[string]interface{} {
    totalAgents := len(as.agents)
    onlineAgents := 0
    totalCapacity := 0
    totalLoad := 0
    
    for _, agent := range as.agents {
        if agent.IsOnline {
            onlineAgents++
        }
        totalCapacity += agent.MaxCapacity
        totalLoad += agent.CurrentLoad
    }
    
    return map[string]interface{}{
        "total_agents":    totalAgents,
        "online_agents":   onlineAgents,
        "offline_agents":  totalAgents - onlineAgents,
        "total_capacity":  totalCapacity,
        "current_load":    totalLoad,
        "utilization":     float64(totalLoad) / float64(totalCapacity),
    }
}