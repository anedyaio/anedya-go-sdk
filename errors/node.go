package errors

import "fmt"

// -------- CreateNode Errors --------

// ErrNodeNameRequired is thrown when someone tries to create a node
// without providing the `node_name`. This is mandatory for creating any node.
var (
	ErrNodeNameRequired = &AnedyaError{
		Message: "node_name is required",
		Err:     ErrVariableNameRequired,
	}
)

// -------- GetNodeList Errors --------

// ErrNodeListRequestNil is returned when the GetNodeList request itself is nil.
// Always make sure you pass a valid request object.
var (
	ErrNodeListRequestNil = &AnedyaError{
		Message: "GetNodeListRequest cannot be nil",
		Err:     fmt.Errorf("request is nil"),
	}

	// ErrNodeListInvalidLimit happens when the `limit` provided is not in the valid range.
	// Limit should always be between 1 and 1000.
	ErrNodeListInvalidLimit = &AnedyaError{
		Message: "limit must be between 1 and 1000",
		Err:     fmt.Errorf("invalid limit"),
	}

	// ErrNodeListInvalidOrder occurs if `order` is neither 'asc' nor 'desc'.
	ErrNodeListInvalidOrder = &AnedyaError{
		Message: "order must be either 'asc' or 'desc'",
		Err:     fmt.Errorf("invalid order"),
	}
)

// -------- GetNodeDetails Errors --------

// ErrNodeDetailsRequestNil is returned when the nodes list is empty or missing
// while fetching node details. Always pass at least one node ID.
var (
	ErrNodeDetailsRequestNil = &AnedyaError{
		Message: "`nodes` list is required",
		Err:     fmt.Errorf("nodes list is empty"),
	}
)

// -------- AddChildNode Errors --------

// ErrAddChildNodeRequestNil is returned when the AddChildNodeRequest is nil.
// We cannot proceed without a proper request.
var (
	ErrAddChildNodeRequestNil = &AnedyaError{
		Message: "AddChildNodeRequest cannot be nil",
		Err:     fmt.Errorf("request is nil"),
	}

	// ErrAddChildNodeParentIdRequired is returned when `parentId` is missing.
	// Every child node must have a parent.
	ErrAddChildNodeParentIdRequired = &AnedyaError{
		Message: "`parentId` is required",
		Err:     fmt.Errorf("parentId missing"),
	}

	// ErrAddChildNodeEmptyChildren occurs when no child nodes are provided
	// in the request. We need at least one child node.
	ErrAddChildNodeEmptyChildren = &AnedyaError{
		Message: "at least one child node must be provided",
		Err:     fmt.Errorf("childNodes empty"),
	}
)

// -------- RemoveChildNode Errors --------

// ErrRemoveChildNodeRequestNil is returned when the RemoveChildNode request is nil.
var (
	ErrRemoveChildNodeRequestNil = &AnedyaError{
		Message: "RemoveChildNodeRequest cannot be nil",
		Err:     fmt.Errorf("request is nil"),
	}

	// ErrRemoveChildNodeParentIdRequired occurs if `parentId` is missing.
	ErrRemoveChildNodeParentIdRequired = &AnedyaError{
		Message: "`parentId` is required",
		Err:     fmt.Errorf("parentId missing"),
	}

	// ErrRemoveChildNodeChildIdRequired occurs if `childNode` is missing in the request.
	ErrRemoveChildNodeChildIdRequired = &AnedyaError{
		Message: "`childNode` is required",
		Err:     fmt.Errorf("childNode missing"),
	}
)

// -------- ClearChildNodes Errors --------

// ErrClearChildNodesRequestNil happens when the request object is nil.
var (
	ErrClearChildNodesRequestNil = &AnedyaError{
		Message: "ClearChildNodesRequest cannot be nil",
		Err:     fmt.Errorf("request is nil"),
	}

	// ErrClearChildNodesParentIdRequired happens if parentId is not provided.
	ErrClearChildNodesParentIdRequired = &AnedyaError{
		Message: "`parentId` is required",
		Err:     fmt.Errorf("parentId missing"),
	}
)

// -------- ListChildNodes Errors --------

// ErrListChildNodesRequestNil happens when the request object is nil.
var (
	ErrListChildNodesRequestNil = &AnedyaError{
		Message: "ListChildNodesRequest cannot be nil",
		Err:     fmt.Errorf("request is nil"),
	}

	// ErrListChildNodesParentIdRequired occurs when parentId is missing in the request.
	ErrListChildNodesParentIdRequired = &AnedyaError{
		Message: "`parentId` is required",
		Err:     fmt.Errorf("parentId missing"),
	}
)

// -------- GetConnectionKey Errors --------

// ErrGetConnectionKeyRequestNil is returned when request object is nil.
var (
	ErrGetConnectionKeyRequestNil = &AnedyaError{
		Message: "GetConnectionKeyRequest cannot be nil",
		Err:     fmt.Errorf("request is nil"),
	}

	// ErrGetConnectionKeyNodeIDRequired occurs when nodeid is missing.
	ErrGetConnectionKeyNodeIDRequired = &AnedyaError{
		Message: "`nodeid` is required",
		Err:     fmt.Errorf("nodeid missing"),
	}
)

// -------- UpdateNode Errors --------

// ErrUpdateNodeRequestNil happens when request object is nil.
var (
	ErrUpdateNodeRequestNil = &AnedyaError{
		Message: "UpdateNodeRequest cannot be nil",
		Err:     fmt.Errorf("request is nil"),
	}

	// ErrUpdateNodeIDRequired occurs when nodeid is missing in update request.
	ErrUpdateNodeIDRequired = &AnedyaError{
		Message: "`nodeid` is required",
		Err:     fmt.Errorf("nodeid missing"),
	}

	// ErrUpdateNodeEmptyUpdates occurs when no updates are provided in the request.
	ErrUpdateNodeEmptyUpdates = &AnedyaError{
		Message: "`updates` must contain at least one update operation",
		Err:     fmt.Errorf("no updates provided"),
	}
)

// -------- AuthorizeDevice Errors --------

// ErrAuthorizeDeviceRequestNil occurs when the request object is nil.
var (
	ErrAuthorizeDeviceRequestNil = &AnedyaError{
		Message: "AuthorizeDeviceRequest cannot be nil",
		Err:     fmt.Errorf("request is nil"),
	}

	// ErrAuthorizeDeviceNodeIDRequired occurs when nodeid is missing in request.
	ErrAuthorizeDeviceNodeIDRequired = &AnedyaError{
		Message: "`nodeid` is required",
		Err:     fmt.Errorf("nodeid missing"),
	}

	// ErrAuthorizeDeviceDeviceIDRequired occurs when deviceid is missing.
	ErrAuthorizeDeviceDeviceIDRequired = &AnedyaError{
		Message: "`deviceid` is required",
		Err:     fmt.Errorf("deviceid missing"),
	}
)

// -------- DeleteNode Errors --------

// ErrDeleteNodeRequestNil happens when request object is nil.
var (
	ErrDeleteNodeRequestNil = &AnedyaError{
		Message: "DeleteNodeRequest cannot be nil",
		Err:     fmt.Errorf("request is nil"),
	}

	// ErrDeleteNodeIDRequired occurs when nodeid is missing in delete request.
	ErrDeleteNodeIDRequired = &AnedyaError{
		Message: "`nodeid` is required",
		Err:     fmt.Errorf("nodeid missing"),
	}
)
