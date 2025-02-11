// Package delve provides a flexible navigation and manipulation API for structured data (maps, slices).
// It enables type-safe traversal and modification of nested data structures using qualified path references.
//
// Key features:
// - Fluent interface for chaining operations
// - Support for both maps and slices through generic interfaces
// - Type-safe value retrieval with automatic wrapping
// - Nested structure navigation through sub-navigators
// - Optional panic-style error handling for mandatory operations
//
// The package operates through Navigator instances that wrap underlying data sources,
// providing consistent access patterns regardless of the source data format.
package delve

import (
	"fmt"

	"github.com/vloldik/delve/v2/internal/quals"
	"github.com/vloldik/delve/v2/internal/sources"
	"github.com/vloldik/delve/v2/internal/value"
	"github.com/vloldik/delve/v2/pkg/interfaces"
)

// Navigator represents a navigation interface into structured data.
// It wraps an underlying data source and provides methods for accessing
// and modifying the data structure. The zero value is not usable - create
// through New or From constructors.
type Navigator = *navigator

// sourceType constraints define valid base types for the generic New constructor.
// Supports either maps (string-keyed) or slices as root data structures.
type sourceType interface{ map[string]any | []any }

// New creates a Navigator from a compatible source structure.
// The generic parameter T must be either map[string]any or []any.
// For complex existing implementations, use From with a custom ISource.
func New[T sourceType](source T) Navigator {
	return &navigator{source: sources.GetSource(source)}
}

// From creates a Navigator from an existing ISource implementation.
// Allows integration with custom data sources that implement the ISource interface.
func From(source interfaces.ISource) Navigator {
	return &navigator{source: source}
}

// navigator implements the core data navigation and manipulation logic.
// Use the exported Navigator type alias instead of direct references.
type navigator struct {
	source interfaces.ISource
}

// Source returns the underlying ISource implementation.
// Useful for accessing low-level data source features not exposed by the Navigator interface.
func (fm *navigator) Source() interfaces.ISource {
	return fm.source
}

// QGetRaw retrieves a raw value from the data source using a qualified path.
// Returns the value and an existence flag. Prefer QGet for type-wrapped values.
func (fm *navigator) QGetRaw(qual interfaces.IQual) (any, bool) {
	return fm.qualGet(qual)
}

// QSet updates the data source at the specified qualified path with the given value.
// Returns true if the operation succeeded. Fails if the path doesn't exist or is read-only.
func (fm *navigator) QSet(qual interfaces.IQual, value any) bool {
	return fm.qualSet(qual, value)
}

// Get retrieves a value using a string-qualified path with optional delimiter customization.
// Default path delimiter is '.'. Returns a value.Value wrapper for type-safe operations.
func (fm *navigator) Get(qual string, _delimiter ...rune) *value.Value {
	return fm.QGet(quals.Q(qual, _delimiter...))
}

// QGet retrieves a qualified path value wrapped in a value.Value container.
// Returns nil-value container if path doesn't exist.
func (fm *navigator) QGet(qual interfaces.IQual) *value.Value {
	v, _ := fm.qualGet(qual)
	return value.New(v)
}

// QGetNavigator retrieves a sub-navigator for a qualified path.
// Useful for chaining operations on nested structures. Returns nil for nonexistent paths.
func (fm *navigator) QGetNavigator(qual interfaces.IQual) Navigator {
	v, ok := fm.qualGet(qual)
	if !ok {
		return nil
	}
	if source := sources.GetSource(v); source != nil {
		return &navigator{source: source}
	} else {
		return nil
	}
}

// GetNavigator retrieves a sub-navigator using string-qualified path.
// Returns nil if path doesn't exist or points to non-navigable data.
func (fm *navigator) GetNavigator(qual string, _delimiter ...rune) Navigator {
	return fm.QGetNavigator(quals.Q(qual, _delimiter...))
}

// QMust retrieves a raw value with panic on missing path.
// Use for mandatory value retrieval. Prefer QGet with existence checks for safer access.
func (fm *navigator) QMust(qual interfaces.IQual) any {
	if val, ok := fm.QGetRaw(qual); ok {
		return val
	}
	panic(fmt.Sprintf("could not get by qual %v", qual))
}

// SetMapSource replaces the underlying data source with a new map.
// Resets all navigation state. Prefer structural updates via QSet when possible.
func (fm *navigator) SetMapSource(source map[string]any) {
	fm.SetSource(sources.MapSource(source))
}

// SetListSource replaces the underlying data source with a new slice.
// Resets all navigation state. Useful for fundamentally changing data structure type.
func (fm *navigator) SetListSource(source []any) {
	fm.SetSource(sources.NewList(source))
}

// SetSource replaces the underlying ISource implementation.
// Allows switching between different data source types while preserving navigation logic.
func (fm *navigator) SetSource(source interfaces.ISource) {
	fm.source = source
}
