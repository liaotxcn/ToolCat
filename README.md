# Weave - A microkernel and layered architecture plugin/service development aggregation platform designed to provide high-performance, highly scalable, secure, and reliable plugin/service development

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Microkernel-Layered-6BA539?style=for-the-badge" alt="Architecture">
  <img src="https://img.shields.io/badge/AI--LLM-74AA9C?style=for-the-badge&logo=brain&logoColor=white" alt="AI-LLM">
  <img src="https://img.shields.io/badge/Plugin%20and Service-FF6F00?style=for-the-badge&logo=ai&logoColor=white" alt="Plugin and Service">
  <img src="https://img.shields.io/badge/Cloud_Native-3371E3?style=for-the-badge&logo=Docker&logoColor=white" alt="Cloud Native">

  **Language Options:** [中文](README.zh-CN.md) | [English](README.md)
</div>


## 📋 Project Introduction

Weave from a simple thread to a complex tapestry, weaving is the creative process from simplicity to complexity. Developers can use Weave to develop various plugins or services, and through connection and integration, plugins and services can be seamlessly integrated and deeply fused, thereby building efficient and stable application systems. Weave, crafting infinite possibilities.

<img width="2590" height="1200" alt="Weaving" src="https://github.com/user-attachments/assets/5dfaa7bd-9817-42f8-847c-446d2f03ce05" />

A high-performance, high-efficiency, pluggable, and easily extensible tool/service development platform based on Golang. It adopts a microkernel + layered architecture design, allowing developers to efficiently develop and easily integrate and manage various tools/services while maintaining system high performance and scalability.

Main application scenarios include:
- Tool development and integration
- Service development and aggregation
- Data/service flow middleware platform
- API gateway and service orchestration
- Efficient development and prototyping platform

---

## 🏗️ Overall Architecture

<img width="1640" height="626" alt="image" src="https://github.com/user-attachments/assets/ad521b7a-4aab-4cc4-8e73-33542a8d7f6c" />

Weave adopts a **microkernel + layered architecture** design pattern, fully leveraging the advantages of both architectures to ensure system availability and achieve high flexibility, scalability, and good performance.

### Fusion of Microkernel and Layered Architecture

Weave integrates the design philosophy of layered architecture on the basis of the microkernel architecture, forming a complete, efficient, and flexible architectural system:

1. **Microkernel Architecture (Plugin System)**: Provides plugin management, lifecycle control, and inter-plugin communication mechanisms.
2. **Layered Architecture (Core System)**: Separates core functions by concerns, forming a clear hierarchical structure.

### Microkernel Architecture Components

1. **Core Kernel**: Provides basic runtime environment, plugin management, configuration management, logging services, security mechanisms, and other basic functions.
2. **Plugin System**: The plugin manager is responsible for plugin registration, lifecycle management, dependency resolution, and conflict detection.
3. **Extension Plugins**: Integrated into the core system through plugin interfaces to implement various business functions.

### Layered Architecture Components

1. **Interface Layer**: Handles HTTP requests, including route management and controllers.
2. **Business Layer**: Contains core business logic and the plugin system.
3. **Data Layer**: Responsible for data storage and access.
4. **Infrastructure Layer**: Provides services such as logging, configuration, and security.

### Architectural Features

**Loose Coupling Design**: The core system and plugins communicate through well-defined interfaces, reducing inter-module dependencies.

**Hot-Plug Capability**: Plugins can be dynamically loaded and unloaded at runtime without restarting the system.

**Functional Isolation**: Each plugin independently encapsulates functionality, having its own namespace and route prefix.

**Dependency and Conflict Management**: Built-in dependency resolution and conflict detection mechanisms ensure harmonious coexistence among plugins.

**Unified Interface**: All plugins implement the same `Plugin` interface, standardizing the development process.

**Extensibility**: System functions can be extended on demand without modifying the kernel code.

**Clear Hierarchy**: The core system uses a layered design, with reasonable code organization, making it easy to maintain and extend.

**High Performance**: The layered design optimizes the request processing flow, improving system response speed.

The core of the system is an efficient and stable plugin mechanism and service aggregation, allowing functional modules to be independently developed and deployed as plugins/services, while interacting through unified interfaces. The overall architectural design emphasizes modularity, scalability, and high performance.

---

## 🌟 Project Features

### 🏗️ Microkernel + Layered Architecture
- **Stable Core & Clear Hierarchy**: The core system remains minimal, and the layered design makes code organization reasonable, easy to maintain, and extend.
- **Flexible Functional Extension**: Extend system functions on demand through the plugin mechanism without modifying the kernel code.
- **Low Coupling & High Cohesion**: Loose coupling between system components facilitates maintenance and upgrades.
- **Hot-Plug Capability**: Plugins can be dynamically loaded and unloaded at runtime without restarting the system.
- **Functional Isolation & Unified Management**: Each plugin independently encapsulates functionality, has its own namespace and route prefix, while core services are uniformly managed through the layered architecture.
- **Dependency and Conflict Management**: Built-in dependency resolution and conflict detection mechanisms ensure harmonious coexistence among plugins.
- **Unified Interface**: All plugins implement the same `Plugin` interface, standardizing the development process.
- **High Performance**: The layered design optimizes the request processing flow, improving system response speed.

### 🚀 High Performance/Efficiency
- Built based on the Gin framework, offering fast request processing and strong concurrency capabilities.
- Optimized database connection pool supporting high concurrent access.
- Modular architecture design with clear code structure, easy to maintain and extend.
- Supports environment variable overrides for easy configuration across different environments.
- Efficient route management supporting dynamic routing and parameter binding.
- The layered architecture optimizes the request processing flow, improving system response speed.

### 🔌 Pluggable & Easy to Extend
- Unified plugin interface design supporting hot-plugging.
- Plugin manager uniformly registers, manages, and executes plugins.
- Plugins can independently register routes and have independent namespaces.
- Plugin dependency and conflict detection mechanisms.
- Scaffolding tool for conveniently generating plugin framework code.
- Example plugins demonstrating the complete plugin development process.

### 🧠 Deep Service Integration
- Seamlessly integrated with the microkernel architecture, highly extensible, and can be used as service aggregation/plugins.
- For example, integrating services like LLM, RAG, etc., providing intelligent Q&A and document retrieval functions.
- Efficient text retrieval based on RedisSearch.
- Supports embedding, indexing, custom embedding models, and retrieval parameters for various document formats.
- Provides flexible retrieval APIs for easy plugin integration.

### 🔒 Secure & Reliable
- JWT-based authentication and authorization system.
- Comprehensive CSRF protection mechanism.
- Rate limiting middleware based on the token bucket algorithm.
- Password hashing storage and verification.
- Detailed login history records.
- Unified error handling middleware.
- Supports HTTPS (can be enabled in configuration).
- The layered architecture encapsulates security mechanisms uniformly in the infrastructure layer, facilitating unified management and maintenance.

### 📊 Observability
- Integrated structured logging system (zap).
- Health check interface for monitoring system status.
- Detailed request/response logging.
- Supports custom monitoring metrics.
- The layered architecture independently encapsulates monitoring functions, ensuring observability of the operational status of each system layer.
- Integrated Prometheus and Grafana monitoring system providing visual dashboards.
- Supports custom alert rule configuration.

### 🚀 Developer Friendly
- Complete plugin development documentation and examples.
- Plugin scaffolding tool for quickly generating plugin templates.
- Supports local development and Docker deployment.
- Clear project structure and code specifications.

---

## 📂 Project Structure

Weave adopts a microkernel + layered architecture, and the project structure clearly reflects this design philosophy. The core system is organized in layers, while functional extensions are achieved through the plugin mechanism/service aggregation.

```
├── .github/
├── .gitignore
├── Dockerfile           # Docker build file
├── Makefile             # Build scripts
├── README.md
├── config/              # Configuration management
├── controllers/         # API controllers [Interface Layer]
├── docker-compose.yaml  # Docker Compose configuration
├── docs/                # Project documentation
├── go.mod
├── go.sum
├── main.go
├── middleware/          # Middleware
├── models/              # Data models [Data Layer]
├── pkg/                 # Common packages [Infrastructure Layer]
├── plugins/             # Plugin system [Core of Microkernel Architecture]
│   ├── core/              # Core plugin functionality
│   ├── doc.go             # Plugin package documentation
│   ├── examples/          # Example plugins
│   ├── features/          # Feature plugins (extensible)
│   ├── init.go            # Plugin initialization
│   ├── loader/            # Plugin loader
│   ├── templates/         # Plugin templates
│   └── watcher/           # Plugin watcher
├── routers/             # Route definition and registration
├── services/            # Service aggregation
│   ├── llm/               # LLM Service
│   ├── rag/               # RAG Service  
│   └── extended/          # Extensible services
├── test/                # Unit/Integration tests
├── tools/               # Development tools
├── utils/               # Utility functions
└── web/                 # Frontend code
```

---

## 🧩 Core Components

### 🔌 Plugin System - Core Implementation of Microkernel Architecture
The plugin system is an important component of Weave, responsible for plugin registration, loading, unloading, and lifecycle management. It implements a complete plugin mechanism, enabling the system to extend functionality in the form of plugins. In the microkernel + layered architecture, the plugin system connects the core kernel with various business extensions.

Under the microkernel architecture, the plugin system has the following characteristics:
- **Complete Lifecycle Management**: Full lifecycle management from plugin initialization, registration, activation to shutdown.
- **Automatic Dependency Resolution**: Automatically resolves dependencies between plugins via the `GetDependencies()` method.
- **Conflict Detection Mechanism**: Avoids functional conflicts between plugins via the `GetConflicts()` method.
- **Automatic Route Registration**: Supports two methods of route registration; the recommended `GetRoutes()` method aligns better with the microkernel architecture design philosophy.
- **Namespace Isolation**: Each plugin has an independent namespace to avoid resource conflicts.
- **Unified Middleware Management**: Supports global and plugin-level middleware configuration.

```go
// Plugin interface definition
type Plugin interface {
    Name() string              // Plugin name
    Description() string       // Plugin description
    Version() string           // Plugin version
    Init() error               // Initialize plugin
    Shutdown() error           // Shutdown plugin
    
    // Route management (new way) - Recommended
    GetRoutes() []Route
    GetDefaultMiddlewares() []gin.HandlerFunc
    
    // Route management (old way) - Retained for compatibility
    RegisterRoutes(*gin.Engine) // Register routes
    
    Execute(map[string]interface{}) (interface{}, error) // Execute function
}

// The Route struct defines route metadata and handler functions.
// This is the core of the new route definition method.
type Route struct {
    Path         string                 // Route path
    Method       string                 // HTTP method (GET, POST, PUT, DELETE, etc.)
    Handler      gin.HandlerFunc        // Request handler function
    Middlewares  []gin.HandlerFunc      // Route-specific middlewares
    Description  string                 // Route description
    AuthRequired bool                   // Whether authentication is required
    Tags         []string               // Route tags, for documentation generation
    Params       map[string]string      // Parameter descriptions, for documentation generation
    Metadata     map[string]interface{} // Custom metadata
}
```

The plugin manager is responsible for the entire lifecycle management of plugins, including registration, deregistration, querying, and executing plugin functions.

### 🧩 Service Aggregation
Service aggregation is an important extension capability of Weave based on the microkernel + layered architecture, providing a mechanism for unified management and invocation of various services, data sources, and functions. For example, the LLM-RAG service aggregation.

- **LLM-RAG Retrieval-Augmented Generation Service**, as one of Weave's service aggregations, provides intelligent text retrieval and enhanced generation functions.
- **Efficient Vector Retrieval**: High-performance vector similarity search based on RedisSearch.
- **Multi-format Document Support**: Supports parsing, chunking, and vectorization of various document formats.
- **Flexible Retrieval API**: Provides rich retrieval interfaces supporting multiple retrieval strategies.
- **Configurable Embedding Models**: Supports switching different embedding models to adapt to different scenario requirements.
- **Integration with Plugin System**: Can be used as infrastructure called by various plugins, enhancing plugin intelligence capabilities.
- **Independent Deployment Option**: Supports running as an independent service or integrating into the main application.

The service aggregation design enhances system functional flexibility, allowing the system to extend and integrate various services and data sources, providing stronger underlying capability support.

### 🔐 Authentication System
The authentication system is located in the infrastructure layer of the layered architecture, providing comprehensive identity authentication and authorization mechanisms, supporting multiple authentication methods. The authentication system is closely integrated with the plugin system, ensuring secure access to plugins, while achieving unified management of security mechanisms through layered design.

- JWT-based token authentication.
- Supports access token and refresh token mechanisms.
- Password hashing storage enhances security.
- Login history records facilitate auditing and tracking.
- Role-based access control.

### 🔄 Middleware System
The middleware system is located between the interface layer and the business layer of the layered architecture, supporting global middleware and plugin-level middleware. It can be used for scenarios such as logging, request validation, and performance monitoring. The middleware system adopts a chain invocation pattern, flexibly combining various functions, reflecting the request processing optimization of the layered architecture.

- **Authentication Middleware**: Verifies user identity.
- **Rate Limiting Middleware**: Prevents API abuse.
- **CORS Middleware**: Handles cross-origin requests.
- **CSRF Protection Middleware**: Prevents cross-site request forgery.
- **Error Handling Middleware**: Uniformly handles and logs errors.

### 📈 Monitoring System
Weave integrates a complete Prometheus + Grafana monitoring system:

- Automatically collects application runtime metrics.
- Pre-configured with various visualization dashboards.
- Supports custom alert rules.
- Real-time monitoring of system health status and performance metrics.

### 🩺 Health Check
The health check function covers all layers of the layered architecture, periodically checking the operational status of various system components to ensure stable system operation. Supports custom health check items to meet the needs of different scenarios. Through the microkernel + layered architecture design, health checks can precisely target the operational status of each plugin and each layer.

- Database connection health check.
- Plugin system status check.
- Overall system health assessment.
- Returns appropriate HTTP status codes based on health status.

## Quick Start

### Environment Preparation
- **Go 1.21+** (for local development)
- **Git** (for cloning the repository)
- **Docker** and **Docker Compose** (for containerized deployment)
- **MySQL 8.0+**
- **PostgreSQL、Redis、Prometheus、Grafana** (Optional, for extension)

### Deployment Methods

#### 1. Docker Compose Deployment (Recommended)

1. Clone the repository
```bash
git clone https://github.com/liaotxcn/weave.git
cd weave
```

2. Create an environment variable file (Optional but recommended)
Create a .env file to set environment variables for enhanced security.

3. Start the services
Use Docker Compose to start the entire service stack with one command:
```bash
docker-compose up -d
```

   On the first startup, Docker Compose will automatically:
   - Build the Docker image for the Weave application.
   - Create the MySQL database container.
   - Create the RedisSearch vector database container.
   - Configure the Prometheus and Grafana monitoring system.
   - Configure networks and volumes.
   - Start all services.
   
   After the services start, you can access the following addresses:
   - Application Backend: http://localhost:8081
   - Prometheus Monitoring: http://localhost:9090
   - Grafana Dashboard: http://localhost:3000 (Default credentials: admin/admin)

4. Verify Service Status
Check if all services are running normally:
```bash
docker-compose ps
```
Under normal circumstances,`weave-app`、`weave-mysql`and`weave-redis` should all show Up status.

### Docker Compose Commands

```bash
docker-compose down    // Stop services
docker-compose logs -f weave-app   // View application logs
docker-compose logs -f weave-mysql // View database logs
docker-compose logs -f weave-redis    // View Redis logs
docker-compose exec weave-app /bin/sh             // Enter the application container
docker-compose exec weave-mysql mysql -u root -p  // Enter the database container
docker-compose exec weave-redis redis-cli    // Enter the Redis container
docker-compose up --build -d        // Rebuild and start services

// Clean up old containers and volume data
docker-compose down -v 
docker system prune -f
docker-compose build --no-cache     // Rebuild images
docker-compose up --force-recreate -d   // Start with --force-recreate option
```

#### 2. Local Development Environment Setup

1. Clone the repository and enter the project directory
```bash
git clone https://github.com/liaotxcn/weave.git
cd weave
```

2. Install dependencies
```bash
go mod download
```

3. Configure the database
Ensure the local MySQL service is started and create the database:
```sql
CREATE DATABASE weave;
```

4. Set environment variables or modify the default configuration in `config/config.go`

5. Run the application
```bash
go run main.go
```

6. Build the application
```bash
go build
```

#### Frontend Build
```bash
cd web
npm install
npm run dev
```

### Notes

1. **Data Persistence**:
   - MySQL data is stored in the `mysql-data` volume, ensuring data is not lost
   - RedisSearch data is stored in the `redis-data` volume, ensuring vector index data is not lost
2. **Health Check**: The system provides a `/health` interface to monitor service health status
3. **Resource Limits**: CPU and memory limits are configured by default and can be adjusted in `docker-compose.yaml` according to actual needs
4. **First Startup**: The first startup requires some time to build images and initialize services
5. **Port Mapping**:
   - By default, the container's port 8081 is mapped to the host's port 8081
   - By default, the container's port 6379 is mapped to the host's port 6379 (RedisSearch)

---

## Project Documentation

### Please read in detail
[API Documentation](./docs/API.md)
[Plugin Development Guide](./docs/PLUGIN_DEVELOPMENT_GUIDE.md)
[Plugin Scaffold Tool Usage](./docs/PLUGIN_SCAFFOLD_USAGE.md)
[Database Migration Guide](./docs/DATABASE_MIGRATION.md)
[Monitoring System Guide](./docs/GRAFANA_MONITORING_GUIDE.md)

### 🔧 Creating a New Plugin

In Weave's microkernel + layered architecture, creating a new plugin is one way to extend system functionality. A plugin is a Go struct that implements the `Plugin` interface. Through this interface, the plugin can interact with the core system. The microkernel architecture provides plugin flexibility, while the layered architecture provides good guidance for the internal code organization of the plugin.

Creating a new plugin is very efficient, just follow these steps:
1. Implement the `plugins.Plugin` interface, defining the plugin's basic information, lifecycle, and functionality
2. Register the plugin in the `registerPlugins` function in `main.go`

Advantages of plugin development in the microkernel + layered architecture:
- **Low Intrusiveness**: Extend system functionality without modifying core code
- **Independent Evolution**: Plugins can be developed, tested, and deployed independently
- **Standardized Interface**: Unified plugin interface simplifies the development process
- **Flexible Combination**: Users can combine different plugins according to their needs
- **Clear Structure**: The layered architecture philosophy guides the internal code organization of plugins, improving maintainability

### Plugin Example (Using the recommended GetRoutes method)
```go
// Example plugin structure
type MyPlugin struct{}

// Methods implementing the Plugin interface
func (p *MyPlugin) Name() string { return "myplugin" }
func (p *MyPlugin) Description() string { return "My custom plugin" }
func (p *MyPlugin) Version() string { return "1.0.0" }
func (p *MyPlugin) Init() error { /* Initialization logic */ return nil }
func (p *MyPlugin) Shutdown() error { /* Shutdown logic */ return nil }

// Register routes using the recommended GetRoutes method
func (p *MyPlugin) GetRoutes() []Route {
    return []Route{
        {
            Path:        "/",
            Method:      "GET",
            Handler:     p.handleIndex,
            Description: "Plugin homepage",
            AuthRequired: false,
            Tags:        []string{"home"},
        },
        {
            Path:        "/api/data",
            Method:      "GET",
            Handler:     p.handleGetData,
            Description: "Get data API",
            AuthRequired: true,
            Tags:        []string{"data", "api"},
            Params: map[string]string{
                "id": "Data ID",
            },
        },
    }
}

// Define the plugin's default middlewares
func (p *MyPlugin) GetDefaultMiddlewares() []gin.HandlerFunc {
    return []gin.HandlerFunc{
        p.logMiddleware,
    }
}

// Route handler functions
func (p *MyPlugin) handleIndex(c *gin.Context) {
    c.JSON(200, gin.H{
        "plugin": p.Name(),
        "version": p.Version(),
    })
}

func (p *MyPlugin) handleGetData(c *gin.Context) {
    id := c.Query("id")
    c.JSON(200, gin.H{
        "id": id,
        "data": "Example data",
    })
}

// Middleware example
func (p *MyPlugin) logMiddleware(c *gin.Context) {
    // Log request
    c.Next()
}

// RegisterRoutes method retained for compatibility
func (p *MyPlugin) RegisterRoutes(router *gin.Engine) {
    // Note: Using GetRoutes method is recommended, this method is only retained for compatibility
    // Can keep empty implementation or add log hint here
}

// Plugin execution logic
func (p *MyPlugin) Execute(params map[string]interface{}) (interface{}, error) {
    // Implement plugin functionality
    return map[string]interface{}{"result": "success"}, nil
}
```

### Plugin Example (RegisterRoutes method - Retained only for compatibility)

```go
// Register plugin routes
func (p *MyPlugin) RegisterRoutes(router *gin.Engine) {
    group := router.Group(fmt.Sprintf("/plugins/%s", p.Name()))
    {
        group.GET("/", func(c *gin.Context) {
            c.JSON(200, gin.H{"plugin": p.Name()})
        })
        // Add more routes...
    }
}
```

### Comparison of Two Route Registration Methods

| Feature | GetRoutes Method (Recommended) | RegisterRoutes Method (Compatibility) |
|---------|-------------------------------|--------------------------------------|
| Route Definition | Uses Route struct array | Directly operates gin.Engine object |
| Metadata Support | ✅ Full support | ❌ Not supported |
| Automatic Route Group | ✅ Automatically created | ❌ Requires manual creation |
| Middleware Management | ✅ Supports global and route level | ❌ Requires manual addition |
| Documentation Generation | ✅ Supports automatic API doc generation | ❌ Not supported |

### 📊 Database Migration Tool

Weave provides an efficient and powerful database migration tool located in the `pkg/migrate` directory, supporting version management of the database structure:

- Implemented based on the `golang-migrate` library
- Supports migration application, rollback, status query, etc.
- Automatically generates version numbers to avoid conflicts
- Supports migration status check and dirty state handling

---

## 🤝 Contribution Guide

Welcome contributions to the project! Thank you!

1. **Fork the repository** and clone it locally
2. **Create a branch** for development (`git checkout -b feature/your-feature`)
3. **Commit your code** and ensure tests pass
4. **Create a Pull Request** describing your changes
5. Wait for **code review** and make modifications based on feedback

---

### <div align="center"> <strong>✨ Continuously updating and improving... ✨</strong> </div>


