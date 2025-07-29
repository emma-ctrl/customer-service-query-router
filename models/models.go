package models

type Query struct {
    Content string `json:"content"`
    Intent  string `json:"intent"`
}

type Agent struct {
    ID           string   `json:"id"`
    Name         string   `json:"name"`
    Specialties  []string `json:"specialties"`
    MaxCapacity  int      `json:"max_capacity"`
    CurrentLoad  int      `json:"current_load"`
    IsOnline     bool     `json:"is_online"`
}

type RoutingResponse struct {
    AgentID string `json:"agent_id"`
    Intent  string `json:"intent"`
}