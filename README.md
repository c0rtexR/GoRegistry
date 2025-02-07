# GoRegistry

A thread-safe, generic registry implementation in Go for managing and orchestrating components.

## Overview

GoRegistry provides a robust, thread-safe implementation of a generic registry pattern in Go. It's designed to help manage and orchestrate different components in your application while maintaining type safety and concurrent access.

## Features

- Thread-safe operations
- Generic type support
- Concurrent access support
- Simple and clean API
- Type-safe implementations
- Support for type-based registration and grouping

## Installation

```bash
go get github.com/c0rtexR/GoRegistry
```

## Usage

### Basic Registry
For simple key-value registration:

```go
import "github.com/c0rtexR/GoRegistry"

// Create a new registry for string types
reg := registry.NewRegistry[string]()

// Register items
err := reg.Register("key1", "value1")
if err != nil {
    // Handle error
}

// Get an item
value, exists := reg.Get("key1")
if exists {
    // Use value
}

// Get all items
items := reg.Items()

// Delete an item
deleted := reg.Delete("key1")

// Get count of items
count := reg.Len()
```

### Type-Safe Registry
For registering components with type safety:

```go
import "github.com/c0rtexR/GoRegistry"

// Define your component types
type ComponentType string

const (
    Tool  ComponentType = "tool"
    Stage ComponentType = "stage"
)

// Define your component interface or struct
type Component interface {
    Execute() error
}

// Create a type-safe registry with allowed component types
reg := registry.NewTypeRegistry[Component, ComponentType](Tool, Stage)

// Register a tool
parser := NewParser()
err := reg.RegisterWithType(Tool, "parser", parser)
if err != nil {
    // Handle error
}

// Register a stage with the same name (allowed!)
parserStage := NewParserStage()
err = reg.RegisterWithType(Stage, "parser", parserStage)

// Get a specific tool
tool, exists := reg.GetByType(Tool, "parser")
if exists {
    // Use tool
}

// Get all tools
tools := reg.ItemsByType(Tool)
for name, tool := range tools {
    // Use tools
}

// Delete a specific tool
deleted := reg.DeleteByType(Tool, "parser")
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests. 