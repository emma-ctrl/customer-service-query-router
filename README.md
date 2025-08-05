# Customer Service Query Router

An intelligent customer service routing system built in Go that uses OpenAI to classify customer queries and automatically route them to specialized AI agents based on intent and agent availability.

## Features

- **AI-Powered Classification**: Uses OpenAI GPT-3.5-turbo to classify customer queries into specific intents
- **Smart Agent Routing**: Routes queries to specialized agents based on skills, availability, and capacity
- **Real-time Monitoring**: Web UI for monitoring agent status and system performance
- **RESTful API**: Complete API for integration with existing customer service platforms
- **Load Balancing**: Automatic load balancing across available agents

## Go Techniques Demonstrated

This project showcases several important Go programming concepts:

### Package Organization & Structure
- **Modular design** with separate packages (`models`, `services`, `handlers`)
- **Clean package boundaries** with clear responsibilities
- **Go modules** for dependency management (`go.mod`)

### HTTP Server & REST API
- **Native `net/http` package** for building web servers
- **HTTP middleware** pattern with CORS handling
- **Route multiplexing** with `http.HandleFunc`
- **File serving** with `http.FileServer` for static assets
- **JSON marshaling/unmarshaling** for API responses

### Struct Design & Methods
- **Struct composition** for modeling domain entities (Agent, Query, etc.)
- **Method receivers** for encapsulating behavior
- **JSON struct tags** for API serialization
- **Constructor functions** (`NewAgentService`, `NewClassificationService`)

### Error Handling
- **Idiomatic error handling** with explicit error returns
- **Error wrapping** using `fmt.Errorf` with `%w` verb
- **Graceful error propagation** through service layers

### Concurrency & Performance
- **Context package** for request lifecycle management
- **HTTP request timeouts** and cancellation
- **Goroutine-safe operations** (though not explicitly concurrent in this version)

### Third-Party Integration
- **External API integration** with OpenAI using `github.com/sashabaranov/go-openai`
- **Environment variable handling** with `os.Getenv`
- **HTTP client usage** for external service calls

### Data Management
- **In-memory data storage** with Go maps
- **File I/O operations** for loading conversation data
- **String manipulation** and text processing

### Logging & Observability
- **Structured logging** with detailed request tracking
- **Performance metrics** collection and reporting
- **Request/response logging** for debugging

## Quick Start

### Prerequisites

- Go 1.24.5 or later
- OpenAI API key

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd customer-service-query-router
   ```

2. Set your OpenAI API key:
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

4. Access the web interface at `http://localhost:8080`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/classify` | Classify customer queries with OpenAI |
| `POST` | `/api/route` | Route customer queries to appropriate agents |
| `GET` | `/api/agents` | Get all agents and their status |
| `GET` | `/api/agents/stats` | Get agent statistics |
| `POST` | `/api/test-conversations` | Test routing with sample conversations |
| `POST` | `/api/test-classification` | Test OpenAI classification on loaded conversations |

## Supported Query Types

The system handles the following customer service intents:

- **Account Support**: Account access issues
- **Billing**: Billing discrepancies, refund processing
- **Delivery**: Delivery problems, order status inquiries
- **Technical Support**: Installation support requests
- **Order Management**: Order cancellations, returns processing
- **Product Support**: Product availability, warranty inquiries
- **General**: Fallback for unclassified queries

## Architecture

- **Models**: Data structures for queries, agents, and routing responses
- **Services**: Core business logic for agent management, classification, and conversation handling
- **Handlers**: HTTP handlers for API endpoints and web UI
- **Static Files**: Web interface for monitoring and testing

## Agent Management

The system includes specialized agents with different capabilities:
- Billing specialists
- Account support representatives  
- Delivery tracking experts
- Product specialists
- Returns processors
- Warranty advisors
- Technical support engineers

Each agent has configurable capacity limits and availability status for intelligent load distribution.

## License

Licensed under the terms specified in the LICENSE file.
