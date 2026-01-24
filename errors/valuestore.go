package errors

import "errors"

var (
	// ErrValueNamespaceIdRequired indicates that the namespace ID was not provided in the request.
	ErrValueNamespaceIdRequired = errors.New("namespace id is required")

	// ErrInvalidNamespaceScope indicates that the provided namespace scope is not supported.
	ErrInvalidNamespaceScope = errors.New("invalid namespace scope")

	// ErrValueKeyRequired indicates that the key field is missing in the request.
	ErrValueKeyRequired = errors.New("key is required")

	// ErrInvalidValueType indicates that the provided value type is not supported.
	ErrInvalidValueType = errors.New("invalid value type")

	// ErrTypeMismatch indicates that the returned value does not match the expected type.
	ErrTypeMismatch = errors.New("value type does not match the expected type")

	// ErrNamespaceScopeRequired indicates that the namespace scope field is missing from the request.
	ErrNamespaceScopeRequired = errors.New("namespace scope is required")

	// ErrOrderByRequired indicates that the 'orderby' field is missing, which is mandatory for scanning keys.
	ErrOrderByRequired = errors.New("orderby is required")

	// ErrInvalidSortOrder indicates that the provided sort order is not 'asc' or 'desc'.
	ErrInvalidSortOrder = errors.New("invalid sort order")

	// ErrInvalidOrderBy indicates that the provided 'orderby' field is not one of the allowed values (namespace, key, created).
	ErrInvalidOrderBy = errors.New("invalid orderby value")

	// ErrValueRequired indicates that the value field is missing, which is mandatory for set value.
	ErrValueRequired = errors.New("Value is required")

	// ErrKeyNotFound indicated that the enterd key is not available
	ErrKeyNotFound = errors.New("Key not found")
)
