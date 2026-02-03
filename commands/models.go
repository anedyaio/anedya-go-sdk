package commands

// CommandStatus represents the current status of a command
type CommandStatus string

const (
	StatusPending     CommandStatus = "pending"     // Command is queued, not yet received by device
	StatusReceived    CommandStatus = "received"    // Command received by device
	StatusProcessing  CommandStatus = "processing"  // Command is being processed by device
	StatusSuccess     CommandStatus = "success"     // Command executed successfully
	StatusFailure     CommandStatus = "failure"     // Command execution failed
	StatusInvalidated CommandStatus = "invalidated" // Command was invalidated/cancelled
)

// CommandDataType represents the type of command data
type CommandDataType string

const (
	DataTypeString CommandDataType = "string" // Text/JSON data
	DataTypeBinary CommandDataType = "binary" // Binary data
)

// CommandFilter contains filter criteria for listing commands
type CommandFilter struct {
	// IssuedAfter filters commands issued after this timestamp
	IssuedAfter string `json:"issuedAfter,omitempty"`

	// IssuedBefore filters commands issued before this timestamp
	IssuedBefore string `json:"issuedBefore,omitempty"`

	// NodeId filters commands for a specific node
	NodeId string `json:"nodeId,omitempty"`

	// Status filters commands by their current status
	Status []CommandStatus `json:"status,omitempty"`

	// Identifier filters commands by their identifier
	Identifier string `json:"identifier,omitempty"`
}

// CommandInfo represents basic information about a command
type CommandInfo struct {
	// Id is the unique identifier for the command
	Id string `json:"id"`

	// Identifier is a user-defined identifier for the command
	Identifier string `json:"identifier"`

	// Status is the current status of the command
	Status CommandStatus `json:"status"`

	// UpdatedOn is the Unix timestamp (milliseconds) when the command was last updated
	UpdatedOn int64 `json:"updatedOn"`

	// Expired indicates whether the command has expired
	Expired bool `json:"expired"`

	// Expiry is the Unix timestamp (milliseconds) when the command will expire
	Expiry int64 `json:"expiry"`

	// IssuedAt is the Unix timestamp (milliseconds) when the command was issued
	IssuedAt int64 `json:"issuedAt"`
}

// CommandDetails represents detailed information about a command
type CommandDetails struct {
	CommandInfo

	// AckData is the acknowledgment data returned by the device
	AckData string `json:"ackdata,omitempty"`

	// AckDataType is the type of acknowledgment data
	AckDataType CommandDataType `json:"ackDataType,omitempty"`

	// Data is the command payload sent to the device
	Data string `json:"data"`

	// DataType is the type of command data
	DataType CommandDataType `json:"dataType"`
}
