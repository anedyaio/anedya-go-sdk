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
