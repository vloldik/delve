# Delve

High-performance nested data navigation library for Go with zero allocations and type-safe access

## Features

- **Zero Allocation** - All operations performed without heap allocations
- **Type-Safe Access** - 15+ type-specific methods (Int(), String(), Bool(), etc)
- **Deep Navigation** - Unlimited depth traversal of map/struct hierarchies
- **Path Escaping** - Handle keys containing delimiters via backslash escaping
- **Dual Source Support** - Native handling of both maps and slices
- **Compiled Paths** - Precompiled qualifiers for repeated access
- **Default Values** - Safe return patterns with configurable defaults

## Installation

```bash
go get github.com/vloldik/delve
```

## Quick Start

```go
package main

import (
	"github.com/vloldik/delve"
)

func main() {
	data := map[string]any{
		"user": map[string]any{
			"id":    123,
			"name":  "John Doe",
			"roles": []any{"admin", "editor"},
		},
	}

	nav := delve.FromMap(data)
	
	// Direct value access
	userID := nav.Int(delve.Qual("user.id"))          // 123
	profile := nav.Navigator(delve.Qual("user"))      // Get nested navigator
	
	// Slice navigation
	firstRole := nav.String(delve.Qual("user.roles.0"))  // "admin"
	
	// Special character handling
	nav.Int(delve.Qual(`payments.invoice\.id`))  // Escaped key access
}
```

## Path Navigation

### Special Characters

```go
// Access keys containing dots
data := map[string]any{
	"a.b": map[string]any{
		"c.d": 42,
	},
}

nav.Int(delve.Qual(`a\.b.c\.d`))  // 42
```

### Custom Delimiters

```go
// Use '/' as path separator
nav := delve.FromMap(data)
value := nav.Int(delve.Qual("a/b/c/d", '/'))
```

### Nested Structures

```go
// Direct sub-navigation
subNav := nav.Navigator(delve.Qual("user.profile"))
email := subNav.String(delve.Qual("contact.email"))
```

## Performance Benchmarks

Tested on 12th Gen Intel® CoreTM i5-12500H (Windows/amd64):

| Test Case                  | Operations/sec | Time/Op | Allocs/Op |
|----------------------------|----------------|---------|-----------|
| Direct Map Access          | 263,452,660    | 4.53ns  | 0         |
| Shallow Delve Access       | 65,676,771     | 17.74ns | 0         |
| Depth 10 Access            | 16,352,204     | 73.35ns | 0         |
| 506 Character Key Access   | 42,016,806     | 27.48ns | 0         |


**Performance Characteristics:**
- Shallow access (1-2 levels) adds ~13ns overhead vs direct access
- Each nesting level adds 6-7ns
- Key length variations (2-506 chars) add <10ns penalty
- All operations achieve 0 allocations
- Linear O(n) time complexity relative to path depth

## API Reference

### Core Methods

| Method               | Return Type     | Description                                |
|----------------------|-----------------|--------------------------------------------|
| `GetByQual()`        | `(any, bool)`   | Direct value access with success check     |
| `MustGetByQual()`    | `any`           | Panic-on-error value retrieve             |
| `Navigator()`        | `*Navigator`    | Sub-context for nested data                |
| `Int()/String()/etc` | `<type>`        | Type-converted values with default support|

### Supported Conversions

**Numeric Types**
- All int/uint variants (8/16/32/64)
- float32/float64
- complex64/complex128

**Special Types**
- `[]byte` and `[]rune` slices
- `map[string]string` maps
- Nested `[]any` slices
- Sub-`Navigator` contexts

### Usage Patterns

**Safe Value Access**
```go
// With default value
age := nav.Int(delve.Qual("user.age"), 25) 

// Default to search depth
depth := nav.Uint32(delve.Qual("config.pagination.limit"), 50)
```

**Mandatory Values**
```go
// Panics if value not found
criticalID := nav.MustGetByQual(delve.Qual("system.critical_id")) 
```

**Slice Operations**
```go
// Access slice elements
lastItem := nav.String(delve.Qual("results.-1"))  
 
// Multidimensional slices
matrixValue := nav.Float64(delve.Qual("matrix.3.14"))
```

## License

MIT License - See [LICENSE](LICENSE) for full text.