# Delve v3

Delve is a Go library that simplifies navigating and manipulating complex, nested data structures (maps and slices) with minimal overhead. It provides a type-safe, fluent interface for accessing and modifying data using path strings or precompiled qualifiers, optimizing for both ease of use and performance.  Think of it as a safer, more Go-idiomatic alternative to reflection for common data manipulation tasks.

## Why Delve?

*   **Reduce Boilerplate:** Avoid deeply nested `if x, ok := ...` checks when accessing nested data.
*   **Type Safety:**  No more `interface{}` casting everywhere!  Delve's `Value` wrapper provides type-specific methods.
*   **Performance:**  Minimize allocations, especially when reading data or modifying *existing* paths.
*   **Flexibility:**  Handle dynamic paths (`Q`) and optimize for frequently accessed, static paths (`CQ`).
*   **Expressive Paths:**  Append to lists, access elements from the end, and escape special characters.

## Key Features

*   **Zero Allocation (Reads):** Reading values using `Get`, `QGet`, `CQ` involves *zero* heap allocations.
*   **Minimal Allocation (Writes):**  `QSet` only allocates memory when creating *new* maps or appending to slices.  Modifying existing paths is allocation-free.
*   **Type-Safe Access:** The `Value` wrapper provides methods like `.Int()`, `.String()`, `.Bool()`, `.Float64()`, etc., ensuring type safety.  Fallback to `.SafeInterface()` for type-checked access to the underlying `any`.
*   **Compiled Qualifiers (CQ):**  Create optimized qualifiers for paths you access repeatedly.  This significantly boosts performance for static paths.
*   **Dynamic Qualifiers (Q):** Build paths at runtime.  Useful for situations where the path isn't known in advance.
*   **Deep Navigation:**  Easily traverse deeply nested data structures using dot-separated paths (or custom delimiters).
*   **In-Place Modification:** Modify nested data directly using `QSet`.  Delve handles creating intermediate maps and slices as needed.
*   **List Manipulation:**
    *   **Append:**  Use the `"+"` qualifier to append to a list.
    *   **Negative Indices:** Access list elements from the end using negative indices (e.g., `-1` for the last element).
*   **Custom Delimiters:**  Configure the path separator (default is `.`).  This allows you to work with paths that use slashes (`/`) or other separators.
*   **Safe Raw Access:**  The `Value.SafeInterface(defaultValue any)` method retrieves the underlying `any` value with type-checking.  If the path or type is invalid, it returns the provided default value.
* **Length Retrieval:** `Value.Len()` gets the length of strings, slices, arrays, maps or channels. Returns -1 if not applicable.
*   **Sub-Navigators:**  Use `GetNavigator` or `QGetNavigator` to obtain a new `Navigator` instance focused on a specific sub-section of your data.  This allows you to chain operations.
*   **Iteration:** `IterList` and `IterMap` provide type safe ways to iterate through slices and maps.

## Getting Started

```go
package main

import (
	"fmt"
	"github.com/vloldik/delve/v3"
)

func main() {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
			"age":  30,
			"address": map[string]any{
				"street": "123 Main St",
				"city":   "Anytown",
			},
			"roles": []string{"admin", "editor"},
		},
	}

	nav := delve.New(data)

	// Get a value using a string path (dynamic qualifier)
	age := nav.Get("user.age").Int()
	fmt.Println("Age:", age) // Output: Age: 30

	// Create a compiled qualifier (for frequently accessed paths)
	streetQual := delve.CQ("user.address.street")
	street := nav.QGet(streetQual).String()
	fmt.Println("Street:", street) // Output: Street: 123 Main St

    // Setting an existing value
    ok := nav.QSet(delve.CQ("user.age"), 31)
    fmt.Println("Set successful:", ok)
    fmt.Println("New Age:", nav.Get("user.age").Int())

	// Append to a list
	nav.QSet(delve.CQ("user.roles.+"), "viewer")
	fmt.Println("Roles:", nav.Get("user.roles").StringSlice()) // Output: [admin editor viewer]

	// Access from the end of a list
	lastRole := nav.Get("user.roles.-1").String() // or nav.QGet(delve.CQ("user.roles.-1"))
	fmt.Println("Last Role:", lastRole)       // Output: Last Role: viewer


	// Get a sub-navigator
	userNav := nav.GetNavigator("user")
	name := userNav.Get("name").String()
	fmt.Println("Name (via sub-navigator):", name) // Output: Name (via sub-navigator): Alice


    // SafeInterface with a default value (if the requested type is not assignable, you get the default)
    invalid := nav.Get("user.address.zipcode").SafeInterface(12345).(int)
    fmt.Println("Zipcode (default):", invalid) // Output: 12345

    // Attempting to get a Value as an incorrect type
    notAnInt := nav.Get("user.name").Int() // String "Alice" is not an int
    fmt.Println("notAnInt value:", notAnInt)

    // Using Len()
    fmt.Println("Length of roles", nav.Get("user.roles").Len())
}
```

## Core API

*   **`New(data any)` / `From(source idelve.ISource)`:** Creates a `Navigator` instance.
    *   `New` works directly with `map[string]any` and `[]any`.
    *   `From` allows you to use a custom data source that implements the `idelve.ISource` interface.

    ```go
    data := map[string]any{"a": 1}
    nav := delve.New(data)

    // Or, with a custom source:
    // mySource := myCustomSource{}
    // nav := delve.From(mySource)
    ```

*   **`Get(path string, ...delimiter rune)`:** Retrieves a value using a string path.  Uses `Q` internally to create a dynamic qualifier. Returns a `*delve.Value` wrapper.
    ```go
    value := nav.Get("path.to.value")
    intValue := value.Int() // Get as int (returns 0 if not an int)
    ```

*   **`QGet(qualifier idelve.IQual)`:** Retrieves a value using a pre-created qualifier (either from `Q` or `CQ`).  Returns a `*delve.Value` wrapper.

    ```go
    qual := delve.CQ("path.to.value")
    value := nav.QGet(qual)
    stringValue := value.String()
    ```

*   **`QSet(qualifier idelve.IQual, value any)`:** Sets a value at the specified path.  Creates any necessary intermediate maps or slices.  Returns `true` on success, `false` on failure (e.g., trying to set a value on a non-existent list index without appending).

    ```go
    qual := delve.CQ("path.to.new.value")
    success := nav.QSet(qual, 42) // Creates "path", "to", "new" maps if needed
    ```

*   **`CQ(path string, ...delimiter rune)`:** Creates a *compiled* qualifier.  Use this for paths that you access repeatedly.  The compilation step happens only once, leading to significant performance gains for frequent access.
    ```go
      var myQual = delve.CQ("user.profile.settings.theme")
      // use nav.QGet(myQual) many times
    ```

*   **`Q(path string, ...delimiter rune)`:** Creates a qualifier from a string path.  Use this for dynamic paths that are constructed at runtime.
    ```go
    key := "dynamicKey"
	qual := delve.Q("data." + key)
    value := nav.QGet(qual)
    ```

*   **`GetNavigator(path string, ...delimiter rune)` / `QGetNavigator(qualifier idelve.IQual)`:** Gets a `Navigator` for a nested section of the data.  Returns `nil` if the path doesn't exist or points to a non-navigable value (like a primitive type).  This is extremely useful for chaining operations.

    ```go
    userNav := nav.GetNavigator("user.profile") // Get a navigator for the "user.profile" section
    theme := userNav.Get("theme").String()     // Access "theme" relative to "user.profile"
    ```

*   **`Value.SafeInterface(defaultValue any)`:** Retrieves the value as `any`, but with type safety. It returns the `defaultValue` if the value at the path has a different type than `defaultValue` or does not exist.
    ```go
    intValue := nav.Get("maybe.an.int").SafeInterface(0).(int)
    stringValue := nav.Get("maybe.a.string").SafeInterface("default").(string)
    boolValue := nav.Get("maybe.a.bool").SafeInterface(false).(bool)
    ```

*   **`Value.Len() int`:**  Gets the length of the underlying value if it's a string, slice, array, map or channel. Returns -1 if the value is `nil` or does not have a length.

    ```go
    mySlice := []int{1, 2, 3}
    nav := delve.New(mySlice)
    length := nav.Get("").Len() // Get length of the root slice.  Returns 3.

    myString := "hello"
    strNav := delve.New(map[string]any{"str": myString})
    strLength := strNav.Get("str").Len() // Returns 5

    myMap := map[string]int{"a": 1, "b": 2}
    mapNav := delve.New(myMap)
    mapLength := mapNav.Get("").Len()
    ```

 *   **`IterList[V any](val *Value, callback func(int, V))`:** Iterates a slice in a `Value`.
   
    ```go
        mySlice := []int{1, 2, 3}
        nav := delve.New(map[string]any{"list_key":mySlice})

        delve.IterList(nav.Get("list_key"), func(i int, v int) bool{
          fmt.Printf("Index %d: %d\n", i, v)
           return false; // continue
        })
     ```

 *   **`IterMap[K comparable, V any](val *Value, callback func(K, V))`:** Iterates a map in a `Value`.
     ```go
        myMap := map[string]int{"a": 1, "b": 2}
        nav := delve.New(map[string]any{"map_key":myMap})
        delve.IterMap(nav.Get("map_key"), func(k string, v int) bool {
           fmt.Printf("Key %s: %d\n", k, v)
            return false; // continue
        })
     ```

## Qualifiers: `CQ` vs. `Q`

*   **`CQ (Compiled Qualifier)`:**
    *   Use `CQ` for *static* paths that are known at compile time and are accessed *multiple times*.
    *   `CQ` *compiles* the path into an internal representation that allows for much faster lookups.
    *   The compilation overhead is paid only *once*, so subsequent `QGet` calls using the same `CQ` are very efficient.

*   **`Q (Dynamic Qualifier)`:**
    *   Use `Q` for *dynamic* paths that are built at runtime or are only used once.
    *   `Q` also parses the path string, but it doesn't perform the same level of optimization as `CQ`.

**Recommendation:**  If you have a path that you know you'll be using repeatedly, *always* use `CQ`.  If you're building a path on the fly, use `Q`.

## Path Features

*   **Escaping Special Characters:** Use a backslash (`\`) to escape special characters within your path string. For example, if you have a key that contains a dot, you would escape it like this:

    ```go
    data := map[string]any{"user.name": "Alice"}
    nav := delve.New(data)
    name := nav.Get("user\\.name").String() // Access the "user.name" key
    fmt.Println(name) // Output: Alice
    ```

*   **List Append:** Use the `"+"` qualifier to append a value to the end of a list.

    ```go
    data := map[string]any{"numbers": []int{1, 2, 3}}
    nav := delve.New(data)
    nav.QSet(delve.CQ("numbers.+"), 4) // Append 4 to the "numbers" list
    fmt.Println(nav.Get("numbers").StringSlice()) // Output: [1 2 3 4]
    ```

*   **Negative List Indices:** Access list elements from the end using negative indices.  `-1` refers to the last element, `-2` to the second-to-last, and so on.

    ```go
    data := map[string]any{"numbers": []int{1, 2, 3}}
    nav := delve.New(data)
    last := nav.Get("numbers.-1").Int()   // Get the last element (3)
    secondToLast := nav.Get("numbers.-2").Int() // Get the second-to-last element (2)
    fmt.Println(last, secondToLast) // Output: 3 2
    ```

*   **Custom Delimiters:** You can specify a custom delimiter for your path strings.  The default delimiter is `.`.

    ```go
    data := map[string]any{"a/b": map[string]any{"c/d": 123}}
    nav := delve.New(data)

    value := nav.Get("a/b/c/d", '/').Int()  // Use '/' as the delimiter
    fmt.Println(value)
    // or
    value1 := nav.QGet(delve.Q("a/b/c/d", '/')).Int()  // Use '/' as the delimiter and Q
    fmt.Println(value1)
    // or

    qualifier := delve.CQ("a/b/c/d", '/') // Use '/' as the delimiter
	fmt.Println(nav.QGet(qualifier).Int())
    ```

## Performance

*   **`CQ` vs. `Q`:**  `CQ` is significantly faster than `Q` for repeated access to the same path. This is because `CQ` pre-compiles the path.  `Q` is suitable for one-off or dynamically generated paths.
*   **`QSet` Allocation:**  `QSet` *only* allocates memory when it needs to create new maps along the path or when appending to a list. If you're setting a value at a path that already exists, `QSet` will not allocate any new memory on the heap.  This makes Delve very efficient for modifying existing data structures. List appends cause an allocation due to boxing.
*   Prefer direct struct \ map access when possible.

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

## Error Handling
Delve prioritizes type safety and offers multiple ways to handle potential errors:
*   **`Value` Methods:** Type specific methods (e.g. `.Int()`,`.String()`) return the corresponding zero-value of the type, if the value cannot be converted.
*  **`SafeInterface`:** The `.SafeInterface()` method lets you provide a default, and ensures you always receive a value of the type you expect.
*   **`QSet` Return Value:** `QSet` returns `true` if successful, or `false` if not.

## Advanced Usage

### Custom Data Sources with `idelve.ISource`

You can use Delve with data sources other than `map[string]any` and `[]any` by implementing the `idelve.ISource` interface. This allows you to use Delve with custom data providers.
