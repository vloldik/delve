# Delve v3

Delve is a high-performance, zero-allocation, type-safe Go library for navigating and manipulating nested data structures (maps and slices). It uses path strings or precompiled qualifiers for efficient access and modification.

## Key Features

*   **Zero Allocation:** Minimizes heap allocations for reads and common write operations.
*   **Type-Safe:**  Retrieves values with type-specific methods (e.g., `.Int()`, `.String()`, `.Bool()`).
*   **Compiled Qualifiers (CQ):**  Optimized for repeated access to the *same* paths.
*   **Deep Navigation:**  Traverses nested maps and slices.
*   **In-Place Modification:**  Modifies nested structures with `QSet`.
*   **List Manipulation:** Appends to lists (`"+"`) and supports negative indices (access from the end).
*   **Custom Delimiters:** Configure the path separator (default is `.`).
*   **Safe Raw Access:** `SafeInterface()` retrieves `any` values with type-checking and default values.
* **Length Retrieval:** `Len()` gets the length of strings, slices, arrays, maps and channels.

## Installation

```bash
go get github.com/vloldik/delve/v3
```

## Core API

*   **`New(data any)` / `From(source idelve.ISource)`:** Creates a `Navigator` from `map[string]any`, `[]any`, or a custom `idelve.ISource`.
*   **`Get(path string, ...delimiter rune)`:** Gets a wrapped value (`*delve.Value`) using a string path.  Uses `Q` internally.
*   **`QGet(qualifier idelve.IQual)`:** Gets a wrapped value (`*delve.Value`).  Use methods on `Value` (e.g., `.Int()`, `.String()`).
*   **`QSet(qualifier idelve.IQual, value any)`:** Sets a value. Creates intermediate maps/slices as needed.
*   **`CQ(path string, ...delimiter rune)`:** Creates a *compiled* qualifier. Use for static, frequently accessed paths.
*   **`Q(path string, ...delimiter rune)`:** Creates a qualifier from a string. Use for dynamic paths.
*   **`GetNavigator(qualifier string, ...delimiter rune)` / `QGetNavigator(qualifier idelve.IQual)`:** Gets a `Navigator` for a nested section (or `nil`).
*   **`Value.SafeInterface(defaultValue any)`:** Gets the value as `any`.  Returns `defaultValue` if the path is invalid or type is not assignable.
*   **`Value.Len() int`:** Gets length (strings, slices, arrays, maps, channels). Returns -1 if not applicable.
*   **`IterList[V any](val *Value, callback func(int, V))`:** Iterates a slice in a `Value`.
*   **`IterMap[K comparable, V any](val *Value, callback func(K, V))`:** Iterates a map in a `Value`.

## Qualifiers: `CQ` vs. `Q`

*   **`CQ`:**  For *static* paths, compiled for faster, repeated use.
*   **`Q`:**  For *dynamic* paths built at runtime.

## Path Features

*   **Escaping:**  `\` escapes special characters like `.`:  `delve.CQ("data.with\\.dots.key")`.
*   **List Append:**  `"+"` appends: `nav.QSet(delve.CQ("list.+"), value)`.
*   **Negative Indices:**  Access from the end: `nav.QGet(delve.CQ("list.-1"))`.
*   **Custom Delimiters:** `delve.CQ("a/b/c", '/')`.

## Performance

`CQ` is substantially faster than `Q` for repeated access to the same path.  `QSet` has zero allocation cost if the path already exists, allocating only when maps need creating or slices are being appended to by boxing. `delve.Len()` is slower than native `len()` as type casting is required. Direct map/slice access (e.g., `data["key"]`) remains the fastest option when possible.

## License

MIT License