package handlers

import (
    "net/http"
    "path/filepath"
)

type UIHandler struct{}

func NewUIHandler() *UIHandler {
    return &UIHandler{}
}

func (uh *UIHandler) ServeHome(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    
    http.ServeFile(w, r, filepath.Join("static", "index.html"))
}

// Enable CORS for the API endpoints
func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}