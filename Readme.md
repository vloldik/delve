# Delve

High-performance nested data navigation library for Go with zero allocations, type-safe access, and mutability.

## Features

-   **Zero Allocation** -  Most operations, including reads and many writes, are performed without heap allocations.
-   **Type-Safe Access** - 20+ type-specific methods (Int(), String(), Bool(), etc.).
-   **Path Precompilation** - Optimized access with `CQ` for static paths.
-   **Deep Navigation** - Unlimited depth traversal of maps/slices.
-   **Escape Handling** - Backslash escaping for special characters in keys.
-   **Custom Delimiters** - Flexible path separator configuration.
-   **Dual Qualifiers** - Choose between compiled (`CQ`) or string-based (`Q`) paths.
-   **Default Values** - Safe return patterns with configurable fallbacks.
-   **Mutable Data Structures** -  Modify nested maps and slices in-place with `QualSet`.
-   **Append to Lists** - Add elements to the end of nested lists using the `"+"` index.
-   **Negative List Indexing** - Access and modify list elements from the end using negative indices.

## Installation

```bash
go get github.com/vloldik/delve
```

## Quick Start

```go
package main

import (
    "fmt"
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
    userID := nav.Int(delve.CQ("user.id")) // 123

    // Dynamic path handling
    roleIndex := 0
    firstRole := nav.String(delve.Q(fmt.Sprintf("user.roles.%d", roleIndex))) // "admin"

    // Nested navigation
    subNav := nav.Navigator(delve.CQ("user"))
    userName := subNav.String(delve.CQ("name")) // "John Doe"

    // Mutating the data:  Change the user's name.
    nav.QualSet(delve.CQ("user.name"), "Jane Smith")
    fmt.Println(nav.String(delve.CQ("user.name"))) // Output: Jane Smith

    // Append a role to the user's roles:
    nav.QualSet(delve.CQ("user.roles.+"), "viewer")
    fmt.Println(nav.String(delve.CQ("user.roles.2"))) // Output: viewer

	// Mutating with a list
	listData := []any{
		map[string]any{"id": 1},
		map[string]any{"id": 2},
	}
  	listNav := delve.FromList(listData)

	// Changing a value deep withing a nested structure
	listNav.QualSet(delve.CQ("0.id"), 5)
	fmt.Println(listNav.Int(delve.CQ("0.id"))) // Output: 5

	// Using negative indexing to update the last elements id to 10
	listNav.QualSet(delve.CQ("-1.id"), 10)
	fmt.Println(listNav.Int(delve.CQ("-1.id"))) // Output: 10
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

### Value Modification

```go
// Set a value deep within a nested structure
nav.QualSet(delve.CQ("user.profile.address.city"), "New York")

// Overwrite an existing value
nav.QualSet(delve.CQ("user.id"), 456)

// Append to a list
nav.QualSet(delve.CQ("items.+"), "new_item")

// Modify an element in a list using a negative index
nav.QualSet(delve.CQ("items.-1"), "updated_item")

// Create non existed maps in path
initialMap := map[string]any{"a": 1} // 'b' key is absent initially
delveNav := delve.FromMap(initialMap)
// Attempting to access a non-existent path returns nil
if delveNav.Navigator(delve.CQ("b")) == nil {
	fmt.Println("b is nil initially") // Output b is nil initially
}
// Set operation creates the necessary intermediate maps automatically.
delveNav.QualSet(delve.CQ("b.nestedValue"), 42)
// You can see that 'b' is new map created
fmt.Printf("b is %#v\n", delveNav.Get("b")) // Output b map[string]interface {}{"nestedValue":42}
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

## Benchmark Results

Key insights from performance tests:

-   **Direct access is ~5x faster** than Delve navigation (4.27ns vs 20.49ns).  This highlights the inherent overhead of any abstraction.
-   **Precompiled qualifiers (CQ)** show significant advantages:
    -   2-5x faster access than string-based qualifiers.
    -   Performance remains stable with increasing key length.
-   **Qualifier creation overhead** grows with complexity:
    -   81ns for simple paths vs 1.9μs for 506-character paths.
    -   Linear allocation growth with path depth.
- **Setting Values:** Modifying values with `QualSet` introduces some overhead compared to direct map/slice operations. `QualSet` is optimized for zero allocations when modifying *existing* values, but will allocate when creating new nested structures.

## Performance Comparison (Get)

### Basic Access

| Method       | ns/op  | Allocs/op |
| ------------ | ------ | --------- |
| Direct map   | 4.27   | 0         |
| Delve (CQ)   | 20.49  | 0         |

### Variable Key Length (ns/op)

| Key Length | CompiledQ | StringQ |
| ---------- | --------- | ------- |
| 2          | 21.84     | 31.95   |
| 14         | 22.47     | 42.82   |
| 506        | 30.38     | 372.4   |

### Nested Depth (ns/op)

| Depth | CompiledQ | StringQ |
| ----- | --------- | ------- |
| 1     | 14.08     | 16.14   |
| 5     | 43.42     | 78.74   |
| 10    | 85.05     | 164.2   |

## Qualifier Creation Costs

### By Key Length

| Key Length | ns/op  | Allocs/op |
| ---------- | ------ | --------- |
| 2          | 117.8  | 4         |
| 239        | 1032   | 8         |
| 506        | 1926   | 9         |

### By Nesting Depth

| Depth | ns/op | Allocs/op |
| ----- | ----- | --------- |
| 1     | 81.88 | 3         |
| 5     | 299.7 | 7         |
| 10    | 565.5 | 12        |

## Set Benchmarks
These benchmarks compare various scenarios of setting values within nested data structures using `delve.QualSet`.

| Benchmark                     | ns/op | Allocs/op |
|-------------------------------|-------|-----------|
| SetValueInMap                 | 23.62 |     0    |
| OverwriteValueInMap           | 26.05 |     0    |
| SetValueInList                | 52.49 |     1    |
| SetNestedValueInMap           | 140.0 |     2    |
| AppendToList                  | 53.33 |     1    |
| SetNegativeIndexInList        | 52.43 |     1    |

*   **SetValueInMap**:  Sets a new key-value pair in an initially empty map. This is very fast and allocation-free *after* the initial map allocation.
*   **OverwriteValueInMap**:  Modifies the value of an *existing* key. This demonstrates the zero-allocation optimization for existing paths.
*   **SetValueInList**:  Sets a value at a specific index in a pre-populated list. Involves an allocation for boxing the integer as an `any`.
*   **SetNestedValueInMap**:  Creates a new *nested* path (`a.b.c`) and sets a value. The allocations are for creating the intermediate `map[string]any` levels. On successive calls it will have 0 allocations and time will be about 48ns.
*   **AppendToList**:  Adds an element to the end of a list using the `"+"` index.
*  **SetNegativeIndexInList**: Example of allocations of setting list element by negative index.

## Optimization Guide

1.  **Prefer CQ** for static paths in hot code paths.
2.  **Reuse qualifiers** with `Copy()` for similar paths.
3.  **Batch operations** using cached navigators.
4.  **Avoid qualifier creation** in tight loops.
5.  **Use type-specific methods** for direct access.
6.  **Pre-allocate** maps and slices if you know their maximum size.
7.  When modifying deeply nested structures, consider whether direct access (if possible) would be more efficient *if* the structure of the maps and slices is well-defined and won't change.  `delve` is most valuable when the shape of your data isn't known at compile time.

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