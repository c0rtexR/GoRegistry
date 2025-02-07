package registry

import (
	"sync"
	"testing"
)

func TestRegistry(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		registry := NewRegistry[string]()

		// Test registration
		err := registry.Register("key1", "value1")
		if err != nil {
			t.Errorf("Failed to register item: %v", err)
		}

		// Test duplicate registration
		err = registry.Register("key1", "value2")
		if err == nil {
			t.Error("Expected error when registering duplicate key")
		}

		// Test retrieval
		value, exists := registry.Get("key1")
		if !exists {
			t.Error("Item not found")
		}
		if value != "value1" {
			t.Errorf("Expected 'value1', got '%s'", value)
		}

		// Test non-existent key
		_, exists = registry.Get("nonexistent")
		if exists {
			t.Error("Found item that shouldn't exist")
		}

		// Test items count
		if registry.Len() != 1 {
			t.Errorf("Expected length 1, got %d", registry.Len())
		}

		// Test deletion
		if !registry.Delete("key1") {
			t.Error("Failed to delete existing item")
		}
		if registry.Delete("key1") {
			t.Error("Deleted non-existent item")
		}
	})

	t.Run("concurrent operations", func(t *testing.T) {
		registry := NewRegistry[int]()
		var wg sync.WaitGroup
		numGoroutines := 100

		// Concurrent registrations
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				_ = registry.Register(string(rune(n)), n)
			}(i)
		}
		wg.Wait()

		// Verify all items were registered
		items := registry.Items()
		if len(items) != numGoroutines {
			t.Errorf("Expected %d items, got %d", numGoroutines, len(items))
		}
	})
}

type ComponentType string

const (
	Tool  ComponentType = "tool"
	Stage ComponentType = "stage"
)

type TestComponent struct {
	Name string
}

func TestTypeRegistry(t *testing.T) {
	t.Run("type-safe operations", func(t *testing.T) {
		registry := NewTypeRegistry[TestComponent, ComponentType](Tool, Stage)

		// Test valid registration
		toolComp := TestComponent{Name: "parser"}
		err := registry.RegisterWithType(Tool, "parser", toolComp)
		if err != nil {
			t.Errorf("Failed to register tool: %v", err)
		}

		// Test invalid type
		invalidType := ComponentType("invalid")
		err = registry.RegisterWithType(invalidType, "test", TestComponent{})
		if err == nil {
			t.Error("Expected error when registering with invalid type")
		}

		// Test duplicate name within same type
		err = registry.RegisterWithType(Tool, "parser", TestComponent{})
		if err == nil {
			t.Error("Expected error when registering duplicate name within same type")
		}

		// Test same name different type
		stageComp := TestComponent{Name: "parser"}
		err = registry.RegisterWithType(Stage, "parser", stageComp)
		if err != nil {
			t.Errorf("Failed to register stage with same name: %v", err)
		}

		// Test retrieval by type
		tool, exists := registry.GetByType(Tool, "parser")
		if !exists {
			t.Error("Tool not found")
		}
		if tool.Name != "parser" {
			t.Errorf("Expected tool name 'parser', got '%s'", tool.Name)
		}

		// Test items by type
		tools := registry.ItemsByType(Tool)
		if len(tools) != 1 {
			t.Errorf("Expected 1 tool, got %d", len(tools))
		}
		if _, exists := tools["parser"]; !exists {
			t.Error("Parser tool not found in tools map")
		}

		stages := registry.ItemsByType(Stage)
		if len(stages) != 1 {
			t.Errorf("Expected 1 stage, got %d", len(stages))
		}
		if _, exists := stages["parser"]; !exists {
			t.Error("Parser stage not found in stages map")
		}

		// Test deletion by type
		if !registry.DeleteByType(Tool, "parser") {
			t.Error("Failed to delete tool")
		}
		if registry.DeleteByType(Tool, "parser") {
			t.Error("Deleted non-existent tool")
		}
	})
}
