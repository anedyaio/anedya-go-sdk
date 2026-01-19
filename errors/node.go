// Package errors defines validation and API errors
// used by the Node Management APIs in the Anedya Go SDK.
package errors

import "errors"

// ----------------------------------------------------
// CreateNode validation errors
// ----------------------------------------------------

var (
	// ErrNodeNameRequired is returned when the node_name
	// field is missing or empty while creating a node.
	ErrNodeNameRequired = errors.New("node_name is required")
)

// ----------------------------------------------------
// GetNodeList validation errors
// ----------------------------------------------------

var (
	// ErrNodeListRequestNil is returned when the
	// GetNodeList request is nil.
	ErrNodeListRequestNil = errors.New("GetNodeListRequest cannot be nil")

	// ErrNodeListInvalidLimit is returned when limit
	// is not between 1 and 1000.
	ErrNodeListInvalidLimit = errors.New("limit must be between 1 and 1000")

	// ErrNodeListInvalidOrder is returned when order
	// is neither 'asc' nor 'desc'.
	ErrNodeListInvalidOrder = errors.New("order must be either 'asc' or 'desc'")
)

// ----------------------------------------------------
// GetNodeDetails validation errors
// ----------------------------------------------------

var (
	// ErrNodeDetailsRequestNil is returned when the nodes
	// list is missing or empty.
	ErrNodeDetailsRequestNil = errors.New("nodes list is required")
)

// ----------------------------------------------------
// AddChildNode validation errors
// ----------------------------------------------------

var (
	// ErrAddChildNodeRequestNil is returned when the
	// AddChildNode request is nil.
	ErrAddChildNodeRequestNil = errors.New("AddChildNodeRequest cannot be nil")

	// ErrAddChildNodeParentIDRequired is returned when
	// parentId is missing.
	ErrAddChildNodeParentIDRequired = errors.New("parentId is required")

	// ErrAddChildNodeEmptyChildren is returned when
	// no child nodes are provided.
	ErrAddChildNodeEmptyChildren = errors.New("at least one child node must be provided")

	// ErrAddChildNodeInvalidChild is returned when a child node
	// entry is missing nodeId or alias.
	ErrAddChildNodeInvalidChild = errors.New("each child node requires both nodeId and alias")
)

// ----------------------------------------------------
// RemoveChildNode validation errors
// ----------------------------------------------------

var (
	// ErrRemoveChildNodeRequestNil is returned when the
	// RemoveChildNode request is nil.
	ErrRemoveChildNodeRequestNil = errors.New("RemoveChildNodeRequest cannot be nil")

	// ErrRemoveChildNodeParentIDRequired is returned when
	// parentId is missing.
	ErrRemoveChildNodeParentIDRequired = errors.New("parentId is required")

	// ErrRemoveChildNodeChildIDRequired is returned when
	// childNode is missing.
	ErrRemoveChildNodeChildIDRequired = errors.New("childNode is required")
)

// ----------------------------------------------------
// ClearChildNodes validation errors
// ----------------------------------------------------

var (
	// ErrClearChildNodesRequestNil is returned when the
	// ClearChildNodes request is nil.
	ErrClearChildNodesRequestNil = errors.New("ClearChildNodesRequest cannot be nil")

	// ErrClearChildNodesParentIDRequired is returned when
	// parentId is missing.
	ErrClearChildNodesParentIDRequired = errors.New("parentId is required")
)

// ----------------------------------------------------
// ListChildNodes validation errors
// ----------------------------------------------------

var (
	// ErrListChildNodesRequestNil is returned when the
	// ListChildNodes request is nil.
	ErrListChildNodesRequestNil = errors.New("ListChildNodesRequest cannot be nil")

	// ErrListChildNodesParentIDRequired is returned when
	// parentId is missing.
	ErrListChildNodesParentIDRequired = errors.New("parentId is required")
)

// ----------------------------------------------------
// GetConnectionKey validation errors
// ----------------------------------------------------

var (
	// ErrGetConnectionKeyRequestNil is returned when the
	// GetConnectionKey request is nil.
	ErrGetConnectionKeyRequestNil = errors.New("GetConnectionKeyRequest cannot be nil")

	// ErrGetConnectionKeyNodeIDRequired is returned when
	// nodeid is missing.
	ErrGetConnectionKeyNodeIDRequired = errors.New("nodeid is required")
)

// ----------------------------------------------------
// UpdateNode validation errors
// ----------------------------------------------------

var (
	// ErrUpdateNodeRequestNil is returned when the
	// UpdateNode request is nil.
	ErrUpdateNodeRequestNil = errors.New("UpdateNodeRequest cannot be nil")

	// ErrUpdateNodeIDRequired is returned when
	// nodeid is missing.
	ErrUpdateNodeIDRequired = errors.New("nodeid is required")

	// ErrUpdateNodeEmptyUpdates is returned when
	// no updates are provided.
	ErrUpdateNodeEmptyUpdates = errors.New("updates must contain at least one update operation")
)

// ----------------------------------------------------
// AuthorizeDevice validation errors
// ----------------------------------------------------

var (
	// ErrAuthorizeDeviceRequestNil is returned when the
	// AuthorizeDevice request is nil.
	ErrAuthorizeDeviceRequestNil = errors.New("AuthorizeDeviceRequest cannot be nil")

	// ErrAuthorizeDeviceNodeIDRequired is returned when
	// nodeid is missing.
	ErrAuthorizeDeviceNodeIDRequired = errors.New("nodeid is required")

	// ErrAuthorizeDeviceDeviceIDRequired is returned when
	// deviceid is missing.
	ErrAuthorizeDeviceDeviceIDRequired = errors.New("deviceid is required")
)

// ----------------------------------------------------
// DeleteNode validation errors
// ----------------------------------------------------

var (
	// ErrDeleteNodeRequestNil is returned when the
	// DeleteNode request is nil.
	ErrDeleteNodeRequestNil = errors.New("DeleteNodeRequest cannot be nil")

	// ErrDeleteNodeIDRequired is returned when
	// nodeid is missing.
	ErrDeleteNodeIDRequired = errors.New("nodeid is required")
)

// ErrNodeManagementNotInitialized is returned when NodeManagement client is nil
var ErrNodeManagementNotInitialized = errors.New("node management client is not initialized")

// ErrNodeNotFound is returned when node details are not found
var ErrNodeNotFound = errors.New("node not found")
