# Delve

High-performance nested data navigation library for Go with zero allocations and type-safe access

## Features

- **Zero Allocation** - All operations performed without heap allocations
- **Type-Safe Access** - 20+ type-specific methods (Int(), String(), Bool(), etc)
- **Path Precompilation** - Optimized access with `CQ` for static paths
- **Deep Navigation** - Unlimited depth traversal of maps/slices
- **Escape Handling** - Backslash escaping for special characters in keys
- **Custom Delimiters** - Flexible path separator configuration
- **Dual Qualifiers** - Choose between compiled (`CQ`) or string-based (`Q`) paths
- **Default Values** - Safe return patterns with configurable fallbacks

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
	
	// Compiled path access
	userID := nav.Int(delve.CQ("user.id"))          // 123
	
	// Dynamic path handling
	roleIndex := 0
	firstRole := nav.String(delve.Q(fmt.Sprintf("user.roles.%d", roleIndex)))  // "admin"
	
	// Nested navigation
	subNav := nav.Navigator(delve.CQ("user"))
	userName := subNav.String(delve.CQ("name"))  // "John Doe"
}
```

## Path Handling

### Special Characters
```go
data := map[string]any{
	"a.b": map[string]any{
		"c.d": 42,
	},
}

nav.Int(delve.CQ(`a\.b.c\.d`))  // 42
```

### Custom Delimiters
```go
// Unix-style paths
value := nav.Int(delve.CQ("a/b/c/d", '/'))

// Mixed delimiters
nav.Int(delve.CQ("first:second:third", ':'))
```

### Slice Navigation
```go
// Negative indexes for reverse access
lastItem := nav.String(delve.CQ("results.-1"))

// Multi-dimensional arrays
matrixValue := nav.Float64(delve.CQ("matrix.3.14"))
```

## Core API

### Value Access
```go
// Safe access with defaults
timeout := nav.Int(delve.CQ("config.timeout"), 30) 

// Type-specific methods
isActive := nav.Bool(delve.CQ("user.active"), false)
ratio := nav.Float32(delve.CQ("metrics.ratio"))

// Mandatory values
criticalID := nav.MustGetByQual(delve.CQ("system.id"))
```

### Navigation Control
```go
// Direct nested access
subNav := nav.Navigator(delve.CQ("user.profile"))
email := subNav.String(delve.CQ("contact.email"))

// Map value check
if val, exists := nav.GetByQual(delve.CQ("optional.path")); exists {
	// Handle value
}
```

## Qualifiers Explained

### Compiled Qualifiers (CQ)
```go
// Precompile static paths
userQual := delve.CQ("user.meta")
name := nav.String(userQual.Copy().String()) 

// Cache frequently used paths
var (
	qualCartTotal = delve.CQ("cart.totals.grand")
	qualInventory = delve.CQ("warehouse.stock.-1")
)

total := nav.Float64(qualCartTotal)
lastStock := nav.Int(qualInventory)
```

### String Qualifiers (Q)
```go
// Dynamic path construction
buildPath := func(region string) delve.IQual {
	return delve.Q(fmt.Sprintf("metrics.regions.%s.active", region))
}

// Runtime configuration
delimiter := getRuntimeDelimiter()
value := nav.Int(delve.Q("path/with/custom/delimiter", delimiter))
```

# Benchmark Results

Key insights from performance tests:

- **Direct access is ~5x faster** than Delve navigation (4.27ns vs 20.49ns)

- **Precompiled qualifiers (CQ)** show significant advantages:

  - 2-5x faster access than string-based qualifiers

  - Performance remains stable with increasing key length

- **Qualifier creation overhead** grows with complexity:

  - 81ns for simple paths vs 1.9μs for 506-character paths

  - Linear allocation growth with path depth

## Performance Comparison

### Basic Access

| Method         | ns/op  | Allocs/op |
|----------------|--------|-----------|
| Direct map     | 4.27   | 0         |
| Delve (CQ)     | 20.49  | 0         |

### Variable Key Length (ns/op)

| Key Length | CompiledQ | StringQ  |
|------------|-----------|----------|
| 2          | 21.84     | 31.95    |
| 14         | 22.47     | 42.82    |
| 506        | 30.38     | 372.4    |

### Nested Depth (ns/op)

| Depth | CompiledQ | StringQ  |
|-------|-----------|----------|
| 1     | 14.08     | 16.14    |
| 5     | 43.42     | 78.74    |
| 10    | 85.05     | 164.2    |

## Qualifier Creation Costs

### By Key Length

| Key Length | ns/op  | Allocs/op |
|------------|--------|-----------|
| 2          | 117.8  | 4         |
| 239        | 1032   | 8         |
| 506        | 1926   | 9         |

### By Nesting Depth

| Depth | ns/op  | Allocs/op |
|-------|--------|-----------|
| 1     | 81.88  | 3         |
| 5     | 299.7  | 7         |
| 10    | 565.5  | 12        |

## Optimization Guide

1. **Prefer CQ** for static paths in hot code paths
2. **Reuse qualifiers** with `Copy()` for similar paths
3. **Batch operations** using cached navigators
4. **Avoid qualifier creation** in tight loops
5. **Use type-specific methods** for direct access

```go
// Optimal pattern
var (
	qualUser = delve.CQ("user")
	qualName = delve.CQ("name")
)

func getUserName(nav *delve.Navigator) string {
	return nav.Navigator(qualUser).String(qualName)
}
```

## License

MIT License - See [LICENSE](LICENSE) for full text.