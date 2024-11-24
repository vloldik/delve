# FlexMap

FlexMap is a Go package that provides a flexible map implementation for handling nested data structures with easy access using qualified paths.

## Features

- Access nested values using dot-notation paths (e.g., "a.b.c")
- Custom delimiter support for qualified paths
- Type-safe value retrieval
- Handles both map and array/slice nested structures
- JSON unmarshaling support

## Installation

```bash
go get github.com/vloldik/flexmap
```

## Usage

### Basic Usage

```go
package main

import (
    "encoding/json"
    "github.com/vloldik/flexmap"
)

func main() {
    // Create and populate FlexMap from JSON
    jsonData := `{
        "user": {
            "profile": {
                "name": "John",
                "scores": [10, 20, 30]
            }
        }
    }`
    
    flexMap := make(flexmap.FlexMap)
    json.Unmarshal([]byte(jsonData), &flexMap)

    // Access nested values using dot notation
    if value, ok := flexMap.GetByQual("user.profile.name"); ok {
        name := value.(string) // "John"
    }

    // Access array elements
    if value, ok := flexMap.GetByQual("user.profile.scores.1"); ok {
        score := value.(float64) // 20
    }
}
```

### Type-Safe Value Retrieval

```go
// Using GetTypedByQual for type-safe value retrieval
name, ok := flexmap.GetTypedByQual[string]("user.profile.name", flexMap)
if ok {
    // name is already a string, no type assertion needed
}
```

### Custom Delimiter

```go
// Change the default delimiter
flexmap.QDelemiter = "/"

// Now use forward slashes for path separation
value, ok := flexMap.GetByQual("user/profile/name")
```

## Interface

### Types

- `FlexMap`: Main type that implements the flexible map functionality
- `FlexList`: Slice type for handling array/slice values
- `IAnyGetter`: Interface for getting values by key

### Methods

#### FlexMap
- `Get(key string) (any, bool)`: Get value by key
- `GetByQual(qual string) (any, bool)`: Get value using qualified path
- `GetTypedByQual[T any](qual string, from FlexMap, allowNil ...bool) (T, bool)`: Get typed value using qualified path

## License

[MIT License](LICENSE)