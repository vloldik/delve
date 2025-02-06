# delve

delve is a Go package that provides a flexible map implementation for handling nested data structures with easy access using qualified paths.

## Features

- Access nested values using dot-notation paths (e.g., "a.b.c")
- Custom delimiter support for qualified paths
- Type-safe value retrieval
- Handles both map and array/slice nested structures
- JSON unmarshaling support

## Installation

```bash
go get github.com/vloldik/delve
```

## Usage

### Basic Usage

```go
package main

import (
    "encoding/json"
    "github.com/vloldik/delve"
)

func main() {
    // Create and populate delve from JSON
    jsonData := `{
        "user": {
            "profile": {
                "name": "John",
                "scores": [10, 20, 30]
            }
        }
    }`
    
    delve := make(delve.delve)
    json.Unmarshal([]byte(jsonData), &delve)

    // Access nested values using dot notation
    if value, ok := delve.GetByQual("user.profile.name"); ok {
        name := value.(string) // "John"
    }

    // Access array elements
    if value, ok := delve.GetByQual("user.profile.scores.1"); ok {
        score := value.(float64) // 20
    }
}
```

### Type-Safe Value Retrieval

```go
// Using GetTypedByQual for type-safe value retrieval
name, ok := delve.GetTypedByQual[string]("user.profile.name", delve)
if ok {
    // name is already a string, no type assertion needed
}
```

### Custom Delimiter

```go
// Change the default delimiter
delve.QDelemiter = "/"

// Now use forward slashes for path separation
value, ok := delve.GetByQual("user/profile/name")
```

## Interface

### Types

- `delve`: Main type that implements the flexible map functionality
- `FlexList`: Slice type for handling array/slice values
- `IAnyGetter`: Interface for getting values by key

### Methods

#### delve
- `Get(key string) (any, bool)`: Get value by key
- `GetByQual(qual string) (any, bool)`: Get value using qualified path
- `GetTypedByQual[T any](qual string, from delve, allowNil ...bool) (T, bool)`: Get typed value using qualified path

## License

[MIT License](LICENSE)