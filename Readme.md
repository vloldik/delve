# Delve

## TODO
Add placeholders to compiled

High-performance nested data navigation library for Go with zero allocations, type-safe access, and mutability.

## Features

-   **Zero Allocation** - Most operations, including reads and many writes, are performed without heap allocations.
-   **Type-Safe Access** - 20+ type-specific methods (Int(), String(), Bool(), etc.).
-   **Path Precompilation** - Optimized access with `CQ` for static paths.
-   **Deep Navigation** - Unlimited depth traversal of maps/slices.
-   **Escape Handling** - Backslash escaping for special characters in keys.
-   **Custom Delimiters** - Flexible path separator configuration.
-   **Dual Qualifiers** - Choose between compiled (`CQ`) or string-based (`Q`) paths.
-   **Default Values** - Safe return patterns with configurable fallbacks.
-   **Mutable Data Structures** - Modify nested maps and slices in-place with `QualSet`.
-   **Append to Lists** - Add elements to the end of nested lists using the `"+"` index.
-   **Negative List Indexing** - Access and modify list elements from the end using negative indices.
-   **Interface Access** -  Retrieve raw `any` values, with safe and unsafe options.
-   **Length Retrieval** - Determine the length of strings, slices, and maps.

## Installation

```bash
go get github.com/vloldik/delve/v2
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/vloldik/delve/v2"
)

func main() {
	data := map[string]any{
		"user": map[string]any{
			"id":    123,
			"name":  "John Doe",
			"roles": []any{"admin", "editor"},
		},
	}

	nav := delve.New(data)

	// Compiled path access
	userID := nav.Int(quals.CQ("user.id")) // 123

	// Dynamic path handling
	roleIndex := 0
	firstRole := nav.String(quals.Q(fmt.Sprintf("user.roles.%d", roleIndex))) // "admin"

	// Nested navigation
	subNav := nav.Navigator(quals.CQ("user"))
	userName := subNav.String(quals.CQ("name")) // "John Doe"

	// Mutating the data:  Change the user's name.
	nav.QualSet(quals.CQ("user.name"), "Jane Smith")
	fmt.Println(nav.String(quals.CQ("user.name"))) // Output: Jane Smith

	// Append a role to the user's roles:
	nav.QualSet(quals.CQ("user.roles.+"), "viewer")
	fmt.Println(nav.String(quals.CQ("user.roles.2"))) // Output: viewer

	// Mutating with a list
	listData := []any{
		map[string]any{"id": 1},
		map[string]any{"id": 2},
	}
	listNav := delve.New(listData)

	// Changing a value deep within a nested structure
	listNav.QualSet(quals.CQ("0.id"), 5)
	fmt.Println(listNav.Int(quals.CQ("0.id"))) // Output: 5

	// Using negative indexing to update the last elements id to 10
	listNav.QualSet(quals.CQ("-1.id"), 10)
	fmt.Println(listNav.Int(quals.CQ("-1.id"))) // Output: 10

    // Using Interface() to get the raw value of roles
    roles := nav.Interface(quals.CQ("user.roles"))
    fmt.Printf("Roles (raw): %v\n", roles)

    //Using SafeInterface() to get user object with default value as empty map
    safeUser := nav.SafeInterface(quals.CQ("user.dog"), map[string]any{})
    fmt.Printf("Safe User: %#v", safeUser)

    // Getting length of roles by Len() method
    rolesLen := nav.Len(quals.CQ("user.roles"))
    fmt.Printf("Roles length is %d\n", rolesLen)
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

nav.Int(quals.CQ(`a\.b.c\.d`))  // 42
```

### Custom Delimiters

```go
// Unix-style paths
value := nav.Int(quals.CQ("a/b/c/d", '/'))

// Mixed delimiters
nav.Int(quals.CQ("first:second:third", ':'))
```

### Slice Navigation

```go
// Negative indexes for reverse access
lastItem := nav.String(quals.CQ("results.-1"))

// Multi-dimensional arrays
matrixValue := nav.Float64(quals.CQ("matrix.3.14"))
```

## Core API

### Value Access

```go
// Safe access with defaults
timeout := nav.Int(quals.CQ("config.timeout"), 30)

// Type-specific methods
isActive := nav.Bool(quals.CQ("user.active"), false)
ratio := nav.Float32(quals.CQ("metrics.ratio"))

// Mandatory values
criticalID := nav.MustGetByQual(quals.CQ("system.id"))
```

### Navigation Control

```go
// Direct nested access
subNav := nav.Navigator(quals.CQ("user.profile"))
email := subNav.String(quals.CQ("contact.email"))

// Map value check
if val, exists := nav.GetByQual(quals.CQ("optional.path")); exists {
	// Handle value
}
```

### Value Modification

```go
// Set a value deep within a nested structure
nav.QualSet(quals.CQ("user.profile.address.city"), "New York")

// Overwrite an existing value
nav.QualSet(quals.CQ("user.id"), 456)

// Append to a list
nav.QualSet(quals.CQ("items.+"), "new_item")

// Modify an element in a list using a negative index
nav.QualSet(quals.CQ("items.-1"), "updated_item")

// Create non existed maps in path
initialMap := map[string]any{"a": 1} // 'b' key is absent initially
delveNav := delve.New(initialMap)
// Attempting to access a non-existent path returns nil
if delveNav.Navigator(quals.CQ("b")) == nil {
	fmt.Println("b is nil initially") // Output b is nil initially
}
// Set operation creates the necessary intermediate maps automatically.
delveNav.QualSet(quals.CQ("b.nestedValue"), 42)
// You can see that 'b' is new map created
fmt.Printf("b is %#v\n", delveNav.Get("b")) // Output b map[string]interface {}{"nestedValue":42}
```

###  `Interface()` and `SafeInterface()`

These methods provide access to the underlying `any` value at a given path.

- **`Interface(interfaces.IQual) any`**:  Returns the raw `any` value.  If the path does not exist or any intermediate node is of the wrong type, it returns `nil`.  This is the *unsafe* version because it's up to the caller to perform type assertions.

- **`SafeInterface(interfaces.IQual, defaultValue any) any`**:  Returns the `any` value, but *only* if it's the correct type (or nil). If the value at the path exists but is not of the expected type, `defaultValue` is returned.  If the path does not exist, `defaultValue` is returned.  This is the *safe* version, as it performs type checking before returning.

```go
// Unsafe access - requires type assertion (and possible panic)
rawValue := nav.Interface(quals.CQ("user.roles"))
rolesArray, ok := rawValue.([]any)
if ok {
    fmt.Println(rolesArray)
}

// Much Safer:  Provides []any{} as a default if the path/type is wrong.
safeValue := nav.SafeInterface(quals.CQ("user.roles"), []any{}).([]any{}) // Will not panic even if user.roles is not []any

safeMap := nav.SafeInterface(quals.CQ("wrong.path", map[string]any{"a": 1})) //will return a map[string]any always
```

### `Len(interfaces.IQual) int`

Retrieves the length of a value at a given path. The behavior depends on the underlying type:

-   **Strings**: Returns the number of *runes* (Unicode code points) in the string.
-   **Slices/Arrays**: Returns the number of elements.
-   **Maps**: Returns the number of key-value pairs.
-    **Other Types:** Returns -1

```go
stringLen := nav.Len(quals.CQ("user.name")) // Length of the name string
rolesLen := nav.Len(quals.CQ("user.roles")) // Number of roles
invalidLen := nav.Len(quals.CQ("user.id"))    // Returns -1 (int is not countable)
notExistLen := nav.Len(quals.CQ("not.exists.path")) // Returns -1
```

## Qualifiers Explained

### Compiled Qualifiers (CQ)

```go
// Precompile static paths
userQual := quals.CQ("user.meta")
name := nav.String(userQual.Copy().String())

// Cache frequently used paths
var (
	qualCartTotal = quals.CQ("cart.totals.grand")
	qualInventory = quals.CQ("warehouse.stock.-1")
)

total := nav.Float64(qualCartTotal)
lastStock := nav.Int(qualInventory)
```

### String Qualifiers (Q)

```go
// Dynamic path construction
buildPath := func(region string) delve.IQual {
	return quals.Q(fmt.Sprintf("metrics.regions.%s.active", region))
}

// Runtime configuration
delimiter := getRuntimeDelimiter()
value := nav.Int(quals.Q("path/with/custom/delimiter", delimiter))
```

## Benchmark Results

Key insights from performance tests:

-   **Direct access is ~5x faster** than Delve navigation (4.27ns vs ~20ns).  This highlights the inherent overhead of *any* abstraction.  Direct access is always ideal when possible.
-   **Precompiled qualifiers (CQ)** show significant advantages:
    -   2-5x faster access than string-based qualifiers (Q).
    -   Performance remains stable with increasing key length, unlike string-based qualifiers, which degrade.
-   **Qualifier creation overhead** grows with complexity:
    -   ~100ns for simple paths, increasing linearly with path depth and length.
    -   Avoid qualifier creation within tight loops.
- **Setting Values:** Modifying values with `QualSet` introduces some overhead compared to direct map/slice operations. `QualSet` is optimized for zero allocations when modifying *existing* values, but will allocate when creating new nested structures.
- **Getting Len from a string with delve is slower than getting the value as a string and calling `len()`, but the difference is minimal.**.

## Performance Comparison (Get)

### Basic Access

| Method                                  | ns/op  | Allocs/op |
| --------------------------------------- | ------ | --------- |
| Direct map                              | 4.27   | 0         |
| Delve (CQ)                              | 20.49  | 0         |
| Delve Getting Inner value               | 122.0  | 2       |
| Delve Getting array (unsafe)            | 39.69  | 0         |
| Delve Getting array (safe)              | 46.45  | 0         |
| Delve Getting string len with  len()    | 41.08  | 0         |
| Delve Getting string len with String()  | 38.97  | 0         |

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

These benchmarks compare various scenarios of setting values within nested data structures using `quals.QualSet`.

| Benchmark                     | ns/op | Allocs/op |
| ----------------------------- | ----- | --------- |
| SetValueInMap                 | 23.62 | 0         |
| OverwriteValueInMap           | 26.05 | 0         |
| SetValueInList                | 52.49 | 1         |
| SetNestedValueInMap           | 140.0 | 2         |
| AppendToList                  | 53.33 | 1         |
| SetNegativeIndexInList        | 52.43 | 1         |

*   **SetValueInMap**:  Sets a new key-value pair in an initially empty map. This is very fast and allocation-free *after* the initial map allocation.
*   **OverwriteValueInMap**:  Modifies the value of an *existing* key. This demonstrates the zero-allocation optimization for existing paths.
*   **SetValueInList**:  Sets a value at a specific index in a pre-populated list. Involves an allocation for boxing the integer as an `any`.
*   **SetNestedValueInMap**:  Creates a new *nested* path (`a.b.c`) and sets a value. The allocations are for creating the intermediate `map[string]any` levels.  On successive calls *to the same path*, it will have 0 allocations.
*   **AppendToList**:  Adds an element to the end of a list using the `"+"` index.
*   **SetNegativeIndexInList**: Example of allocations of setting list element by negative index.

## Optimization Guide

1.  **Prefer CQ over Q:**  Use compiled qualifiers (`quals.CQ`) for static paths, especially in frequently executed code.  String qualifiers (`quals.Q`) are better for dynamically constructed paths.

2.  **Cache Navigators:** If you need to access multiple values within the *same* nested section, create a `Navigator` once and reuse it:

    ```go
    userNav := nav.Navigator(quals.CQ("user"))
    name := userNav.String(quals.CQ("name"))
    id := userNav.Int(quals.CQ("id"))
    ```

3.  **Avoid Qualifier Creation in Loops:** Creating qualifiers (especially `CQ`) has a cost.  Do it *outside* of loops:

    ```go
    // BAD:  Recompiles the qualifier on every iteration
    for i := 0; i < len(myArray); i++ {
        value := nav.Int(quals.CQ(fmt.Sprintf("data.%d.value", i)))
    }

    // GOOD:  Pre-compile the base, and copy/extend inside the loop.
    baseQual := quals.CQ("data")
    for i := 0; i < len(myArray); i++ {
        value := nav.Int(baseQual.Copy().String() + fmt.Sprintf(".%d.value", i))
    }

     // Better:  If the structure is regular use a list
	dataNav := nav.Navigator(quals.CQ("data"))              // Get 'data' navigator
	qualValue := quals.CQ("value") // Pre-compile ".value"
	for i := 0; i < len(myArray); i++ {
	    value := dataNav.Navigator(quals.Q(strconv.Itoa(i))).Int(qualValue) // Access each element and its ".value"
	}
    // Ideal:
    // for _, elem := range myArray {
    //    value := elem.Value
    // }
    ```

4.  **Use Type-Specific Methods:** Use methods like `Int()`, `String()`, `Bool()`, etc., when you know the expected type. Avoid unnecessary type assertions.

5.  **Pre-allocate Maps and Slices:**  If you have a good estimate of the maximum size of maps or slices *before* using Delve, pre-allocate them using `make(map[string]any, size)` or `make([]any, size)`.  This can reduce reallocations if Delve needs to grow these structures.

6.  **Consider Direct Access (When Possible):** Delve is most helpful when you *don't* know the exact structure of your data at compile time. If you *do* have a fixed, well-defined data structure, direct access (e.g., `myData.User.Profile.Address.City`) will *always* be faster.  Use Delve for flexibility, not as a replacement for direct access when the structure is known.

7. **Use `SafeInterface()` for type-safe access to raw `any` values**

```go
// Optimal pattern
var (
	qualUser = quals.CQ("user")
	qualName = quals.CQ("name")
)

func getUserName(nav *delve.Navigator) string {
	return nav.Navigator(qualUser).String(qualName)
}
```

## License

MIT License - See [LICENSE](LICENSE) for full text.
