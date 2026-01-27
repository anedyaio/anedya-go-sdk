// Package commands provides types and structures used
// for managing and tracking commands issued to nodes
// within the Anedya platform.
package commands

// CommandStatus represents the current lifecycle status
// of a command issued to a node.
type CommandStatus string

const (
	// StatusPending indicates the command is queued
	// and has not yet been received by the device.
	StatusPending CommandStatus = "pending"

	// StatusReceived indicates the command has been
	// received by the device.
	StatusReceived CommandStatus = "received"

	// StatusProcessing indicates the command is currently
	// being processed by the device.
	StatusProcessing CommandStatus = "processing"

	// StatusSuccess indicates the command was executed
	// successfully by the device.
	StatusSuccess CommandStatus = "success"

	// StatusFailure indicates the command execution
	// failed on the device.
	StatusFailure CommandStatus = "failure"

	// StatusInvalidated indicates the command was
	// invalidated or cancelled before execution.
	StatusInvalidated CommandStatus = "invalidated"
)

// CommandDataType represents the format of command
// payload or acknowledgment data.
type CommandDataType string

const (
	// DataTypeString represents textual or JSON-based data.
	DataTypeString CommandDataType = "string"

	// DataTypeBinary represents binary command data.
	DataTypeBinary CommandDataType = "binary"
)

// CommandFilter defines filtering criteria used
// when listing or querying commands.
type CommandFilter struct {
	// IssuedAfter filters commands issued after
	// the specified timestamp.
	IssuedAfter string `json:"issuedAfter,omitempty"`

	// IssuedBefore filters commands issued before
	// the specified timestamp.
	IssuedBefore string `json:"issuedBefore,omitempty"`

	// NodeId filters commands associated with
	// a specific node.
	NodeId string `json:"nodeId,omitempty"`

	// Status filters commands by their current
	// lifecycle status.
	Status []CommandStatus `json:"status,omitempty"`

	// Identifier filters commands by a user-defined
	// command identifier.
	Identifier string `json:"identifier,omitempty"`
}

// CommandInfo represents the basic metadata
// associated with a command.
type CommandInfo struct {
	// Id is the unique system-generated identifier
	// for the command.
	Id string `json:"id"`

	// Identifier is a user-defined identifier
	// assigned to the command.
	Identifier string `json:"identifier"`

	// Status represents the current lifecycle
	// status of the command.
	Status CommandStatus `json:"status"`

	// UpdatedOn is the Unix timestamp (milliseconds)
	// when the command was last updated.
	UpdatedOn int64 `json:"updatedOn"`

	// Expired indicates whether the command
	// has expired.
	Expired bool `json:"expired"`

	// Expiry is the Unix timestamp (milliseconds)
	// when the command will expire.
	Expiry int64 `json:"expiry"`

	// IssuedAt is the Unix timestamp (milliseconds)
	// when the command was issued.
	IssuedAt int64 `json:"issuedAt"`
}

// CommandDetails represents detailed information
// about a command, including payload and acknowledgment data.
type CommandDetails struct {
	CommandInfo

	// AckData contains acknowledgment data returned
	// by the device after command execution.
	AckData string `json:"ackdata,omitempty"`

	// AckDataType specifies the format of the
	// acknowledgment data.
	AckDataType CommandDataType `json:"ackDataType,omitempty"`

	// Data contains the command payload
	// sent to the device.
	Data string `json:"data"`

	// DataType specifies the format of the
	// command payload.
	DataType CommandDataType `json:"dataType"`
}
