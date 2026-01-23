package valuestore

// Scope defines the namespace level at which a value is stored.
type Scope string

const (
	ScopeGlobal Scope = "global"
	ScopeNode   Scope = "node"
)

// ValueType represents the supported data type of a stored value.
type ValueType string

const (
	// ValueTypeString represents plain string data.
	ValueTypeString ValueType = "string"

	// ValueTypeBoolean represents boolean data (true/false).
	ValueTypeBoolean ValueType = "boolean"

	// ValueTypeFloat represents floating point numeric data.
	ValueTypeFloat ValueType = "float"

	// ValueTypeBinary represents base64 encoded binary data.
	ValueTypeBinary ValueType = "binary"
)

// isValidScope validates whether the provided scope is supported.
func isValidScope(s Scope) bool {
	switch s {
	case ScopeGlobal, ScopeNode:
		return true
	default:
		return false
	}
}

// isValidValueType validates whether the provided value type is supported.
func isValidValueType(t ValueType) bool {
	switch t {
	case ValueTypeString, ValueTypeBinary, ValueTypeFloat, ValueTypeBoolean:
		return true
	default:
		return false
	}
}

// NameSpace defines where the value is stored in the platform.
type NameSpace struct {

	// Scope determines whether the value is stored globally
	// or under a specific node.
	//
	// Allowed values:
	//   - "global"
	//   - "node"
	Scope Scope `json:"scope"`

	// Id specifies:
	//   - a project-wide unique ID when scope = global
	//   - a valid node ID when scope = node
	Id string `json:"id"`
}

// SortOrder defines the direction of the sort for list operations.
type SortOrder string

const (
	// SortOrderAsc sorts results in ascending order.
	SortOrderAsc SortOrder = "asc"

	// SortOrderDesc sorts results in descending order.
	SortOrderDesc SortOrder = "desc"
)

// ScanOrderBy specifies the field used to sort the results when scanning keys.
type ScanOrderBy string

const (
	// ScanOrderByNamespace sorts results based on the namespace identifier.
	ScanOrderByNamespace ScanOrderBy = "namespace"

	// ScanOrderByKey sorts results alphabetically by the key name.
	ScanOrderByKey ScanOrderBy = "key"

	// ScanOrderByCreated sorts results based on the creation timestamp of the key.
	ScanOrderByCreated ScanOrderBy = "created"
)

// isValidSortOrder validates whether the provided sort order is supported (asc or desc).
func isValidSortOrder(t SortOrder) bool {
	switch t {
	case SortOrderAsc, SortOrderDesc:
		return true
	default:
		return false
	}
}

// isValidOrderby validates whether the provided sort field is supported by the scan API.
func isValidOrderby(t ScanOrderBy) bool {
	switch t {
	case ScanOrderByNamespace, ScanOrderByKey, ScanOrderByCreated:
		return true
	default:
		return false
	}
}
