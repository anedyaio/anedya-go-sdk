package commands

// CommandStatus represents the execution state of a command
// issued to a node within the Anedya platform.
type CommandStatus string

const (
	// StatusPending indicates the command has been created
	// but not yet delivered to the target node.
	StatusPending CommandStatus = "pending"

	// StatusReceived indicates the command has been received
	// by the target node but not yet processed.
	StatusReceived CommandStatus = "received"

	// StatusProcessing indicates the command is currently
	// being executed by the target node.
	StatusProcessing CommandStatus = "processing"

	// StatusSuccess indicates the command executed successfully.
	StatusSuccess CommandStatus = "success"

	// StatusFailure indicates the command execution failed.
	StatusFailure CommandStatus = "failure"

	// StatusInvalidated indicates the command was invalidated
	// before execution due to expiry or manual cancellation.
	StatusInvalidated CommandStatus = "invalidated"
)

// CommandDataType represents the format of command payload
// or acknowledgement data exchanged with nodes.
type CommandDataType string

const (
	// DataTypeString indicates plain text data.
	DataTypeString CommandDataType = "string"

	// DataTypeBinary indicates binary or encoded data.
	DataTypeBinary CommandDataType = "binary"
)

// CommandInfo contains the core metadata associated
// with a command issued to a node.
type CommandInfo struct {
	// CommandId is the unique identifier of the command.
	CommandId string `json:"commandId"`

	// Command is the name or type of the command issued.
	Command string `json:"command"`

	// Status represents the current execution status of the command.
	Status CommandStatus `json:"status"`

	// UpdatedOn indicates the last status update timestamp (Unix milliseconds).
	UpdatedOn int64 `json:"updatedOn"`

	// Expired indicates whether the command has expired.
	Expired bool `json:"expired"`

	// Expiry represents the expiry timestamp of the command (Unix milliseconds).
	Expiry int64 `json:"expiry"`

	// IssuedAt indicates when the command was originally issued (Unix milliseconds).
	IssuedAt int64 `json:"issuedAt"`
}

// CommandDetails extends CommandInfo by including
// payload and acknowledgement data associated with the command.
type CommandDetails struct {
	CommandInfo

	// AckData contains acknowledgement data returned by the node.
	AckData string `json:"ackdata,omitempty"`

	// AckDataType specifies the format of the acknowledgement data.
	AckDataType CommandDataType `json:"ackdatatype,omitempty"`

	// Data contains the original command payload.
	Data string `json:"data,omitempty"`

	// DataType specifies the format of the command payload.
	DataType CommandDataType `json:"datatype,omitempty"`
}
