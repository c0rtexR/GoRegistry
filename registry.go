package registry

import (
	"fmt"
	"sync"
)

// Registry is a generic, thread-safe registry for any type.
type Registry[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

// TypeRegistry adds type-safe registration capabilities.
type TypeRegistry[T any, K comparable] struct {
	registry *Registry[T]
	types    map[K]struct{}
}

// NewRegistry creates a new Registry.
func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		items: make(map[string]T),
	}
}

// NewTypeRegistry creates a new TypeRegistry with specified valid types.
// Example:
//
//	type Tool string
//	const (
//	    Parser Tool = "parser"
//	    Formatter Tool = "formatter"
//	)
//	registry := NewTypeRegistry[MyInterface, Tool](Parser, Formatter)
func NewTypeRegistry[T any, K comparable](validTypes ...K) *TypeRegistry[T, K] {
	types := make(map[K]struct{}, len(validTypes))
	for _, t := range validTypes {
		types[t] = struct{}{}
	}
	return &TypeRegistry[T, K]{
		registry: NewRegistry[T](),
		types:    types,
	}
}

// RegisterWithType adds an item with the given type and name.
// Returns an error if:
// - The type is not registered as valid
// - An item with the same name already exists
// Example:
//
//	registry.RegisterWithType(Parser, "myParser", parserImpl)
func (r *TypeRegistry[T, K]) RegisterWithType(itemType K, name string, item T) error {
	if _, exists := r.types[itemType]; !exists {
		return fmt.Errorf("invalid type: %v", itemType)
	}

	key := fmt.Sprintf("%v:%s", itemType, name)
	return r.registry.Register(key, item)
}

// GetByType retrieves an item by its type and name.
func (r *TypeRegistry[T, K]) GetByType(itemType K, name string) (T, bool) {
	key := fmt.Sprintf("%v:%s", itemType, name)
	return r.registry.Get(key)
}

// DeleteByType removes an item by its type and name.
func (r *TypeRegistry[T, K]) DeleteByType(itemType K, name string) bool {
	key := fmt.Sprintf("%v:%s", itemType, name)
	return r.registry.Delete(key)
}

// ItemsByType returns all items of a specific type.
func (r *TypeRegistry[T, K]) ItemsByType(itemType K) map[string]T {
	// Use write lock to ensure thread safety during map copy
	r.registry.mu.Lock()
	defer r.registry.mu.Unlock()

	prefix := fmt.Sprintf("%v:", itemType)
	result := make(map[string]T)

	// Create a complete copy while holding the lock
	for key, value := range r.registry.items {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			name := key[len(prefix):]
			// Make a deep copy if T is a pointer type
			result[name] = value
		}
	}

	return result
}

// Register adds an item with the given key.
// Returns an error if the key already exists.
func (r *Registry[T]) Register(key string, item T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[key]; exists {
		return fmt.Errorf("item with key '%s' already exists", key)
	}

	r.items[key] = item
	return nil
}

// Get retrieves an item by its key.
// Returns the item and true if found, or zero value and false if not found.
func (r *Registry[T]) Get(key string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[key]
	return item, exists
}

// Delete removes an item from the registry.
// Returns true if the item was found and deleted, false otherwise.
func (r *Registry[T]) Delete(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[key]; !exists {
		return false
	}

	delete(r.items, key)
	return true
}

// Items returns a copy of all registered items.
func (r *Registry[T]) Items() map[string]T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copy := make(map[string]T, len(r.items))
	for k, v := range r.items {
		copy[k] = v
	}
	return copy
}

// Len returns the number of items in the registry.
func (r *Registry[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items)
}
