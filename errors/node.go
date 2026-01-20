// Package errors defines validation and API errors
// used by the Node Management APIs in the Anedya Go SDK.
package errors

import "errors"

// ----------------------------------------------------
// General errors
// ----------------------------------------------------

var (
	// ErrNodeDeviceIDExists is returned when a device ID
	// already exists for a node.
	ErrNodeDeviceIDExists = errors.New("node device id already exists")

	// ErrNodeManagementNotInitialized is returned when
	// the NodeManagement client is nil.
	ErrNodeManagementNotInitialized = errors.New("client not initialized")

	// ErrNodeNotFound is returned when node details
	// are not found.
	ErrNodeNotFound = errors.New("node not found")

	// ErrNodeInvalidUUID is returned when node ID
	// is not a valid UUID.
	ErrNodeInvalidUUID = errors.New("node invalid uuid")
)

// ----------------------------------------------------
// Authentication & Authorization errors
// ----------------------------------------------------

var (
	ErrAccessDenied  = errors.New("access denied")
	ErrTokenNotFound = errors.New("authorization token not found")
)

// ----------------------------------------------------
// GetNodeList validation errors
// ----------------------------------------------------

var (
	// ErrNodeListRequestNil is returned when the
	// GetNodeList request is nil.
	ErrNodeListRequestNil = errors.New("request is nil")

	// ErrNodeListInvalidLimit is returned when limit
	// is not between 1 and 1000.
	ErrNodeListInvalidLimit = errors.New("invalid limit")

	// ErrNodeListInvalidOrder is returned when order
	// is neither 'asc' nor 'desc'.
	ErrNodeListInvalidOrder = errors.New("invalid order")
)

// ----------------------------------------------------
// GetNodeDetails validation errors
// ----------------------------------------------------

var (
	// ErrNodeDetailsRequestNil is returned when the
	// GetNodeDetails request is nil.
	ErrNodeDetailsRequestNil = errors.New("request is nil")
)

// ----------------------------------------------------
// AddChildNode validation errors
// ----------------------------------------------------

var (
	// ErrAddChildNodeRequestNil is returned when the
	// AddChildNode request is nil.
	ErrAddChildNodeRequestNil = errors.New("request is nil")

	// ErrAddChildNodeParentIDRequired is returned when
	// parentId is missing.
	ErrAddChildNodeParentIDRequired = errors.New("parent id required")

	// ErrAddChildNodeEmptyChildren is returned when
	// no child nodes are provided.
	ErrAddChildNodeEmptyChildren = errors.New("children required")

	// ErrAddChildNodeInvalidChild is returned when a child node
	// entry is missing nodeId or alias.
	ErrAddChildNodeInvalidChild = errors.New("invalid child")

	// Hierarchy constraint violations
	ErrNodeChildExists          = errors.New("child node already exists")
	ErrNodeUniqueAliasViolation = errors.New("duplicate alias not allowed for the same parent node")
	ErrNodeUniqueChildViolation = errors.New("duplicate child nodes not allowed for the same parent node")
)

// ----------------------------------------------------
// RemoveChildNode validation errors
// ----------------------------------------------------

var (
	// ErrRemoveChildNodeRequestNil is returned when the
	// RemoveChildNode request is nil.
	ErrRemoveChildNodeRequestNil = errors.New("request is nil")

	// ErrRemoveChildNodeParentIDRequired is returned when
	// parentId is missing.
	ErrRemoveChildNodeParentIDRequired = errors.New("parent id required")

	// ErrRemoveChildNodeChildIDRequired is returned when
	// childNode is missing.
	ErrRemoveChildNodeChildIDRequired = errors.New("child id required")

	ErrNodeChildNotFound   = errors.New("no such child node associated with the parent node")
	ErrNodeInvalidParentID = errors.New("invalid parent node ID")
	ErrNodeInvalidChildID  = errors.New("invalid child node ID")
)

// ----------------------------------------------------
// ClearChildNodes validation errors
// ----------------------------------------------------

var (
	// ErrClearChildNodesRequestNil is returned when the
	// ClearChildNodes request is nil.
	ErrClearChildNodesRequestNil = errors.New("request is nil")

	// ErrClearChildNodesParentIDRequired is returned when
	// parentId is missing.
	ErrClearChildNodesParentIDRequired = errors.New("parent id required")
)

// ----------------------------------------------------
// ListChildNodes validation errors
// ----------------------------------------------------

var (
	// ErrListChildNodesRequestNil is returned when the
	// ListChildNodes request is nil.
	ErrListChildNodesRequestNil = errors.New("request is nil")

	// ErrListChildNodesParentIDRequired is returned when
	// parentId is missing.
	ErrListChildNodesParentIDRequired = errors.New("parent id required")
)

// ----------------------------------------------------
// GetConnectionKey validation errors
// ----------------------------------------------------

var (
	// ErrGetConnectionKeyRequestNil is returned when the
	// GetConnectionKey request is nil.
	ErrGetConnectionKeyRequestNil = errors.New("request is nil")

	// ErrGetConnectionKeyNodeIDRequired is returned when
	// nodeId is missing.
	ErrGetConnectionKeyNodeIDRequired = errors.New("node id required")
)

// ----------------------------------------------------
// UpdateNode validation errors
// ----------------------------------------------------

var (
	// ErrUpdateNodeRequestNil is returned when the
	// UpdateNode request is nil.
	ErrUpdateNodeRequestNil = errors.New("request is nil")

	// ErrUpdateNodeIDRequired is returned when
	// nodeId is missing.
	ErrUpdateNodeIDRequired = errors.New("node id required")

	// ErrUpdateNodeEmptyUpdates is returned when
	// no updates are provided.
	ErrUpdateNodeEmptyUpdates = errors.New("no updates provided")
)

// ----------------------------------------------------
// AuthorizeDevice validation errors
// ----------------------------------------------------

var (
	// ErrAuthorizeDeviceRequestNil is returned when the
	// AuthorizeDevice request is nil.
	ErrAuthorizeDeviceRequestNil = errors.New("request is nil")

	// ErrAuthorizeDeviceNodeIDRequired is returned when
	// nodeId is missing.
	ErrAuthorizeDeviceNodeIDRequired = errors.New("node id required")

	// ErrAuthorizeDeviceDeviceIDRequired is returned when
	// deviceId is missing.
	ErrAuthorizeDeviceDeviceIDRequired = errors.New("device id required")

	ErrNodeDeviceNotFound = errors.New("node device not found")
)

// ----------------------------------------------------
// DeleteNode validation errors
// ----------------------------------------------------

var (
	// ErrDeleteNodeRequestNil is returned when the
	// DeleteNode request is nil.
	ErrDeleteNodeRequestNil = errors.New("request is nil")

	// ErrDeleteNodeIDRequired is returned when
	// nodeId is missing.
	ErrDeleteNodeIDRequired = errors.New("node id required")
)
