// Package errors defines SDK and API error types used by the
// Anedya Go SDK.
package errors

import (
	"fmt"
)

// AnedyaError represents a structured SDK or API error.
//
// It wraps a sentinel error for programmatic checks and includes
// a human-readable message for debugging or logging.
type AnedyaError struct {
	// Message is a human-readable error description.
	Message string

	// Err is the underlying sentinel error.
	Err error
}

// Error implements the error interface.
func (e *AnedyaError) Error() string {
	return fmt.Sprintf("anedya api error: %s: %v", e.Message, e.Err)
}

// Unwrap allows errors.Is to work with AnedyaError.
func (e *AnedyaError) Unwrap() error {
	return e.Err
}

// codeMap maps API reason codes to SDK sentinel errors.
var codeMap = map[string]error{
	// generic errors
	"generic::malformedrequest": ErrInvalidInput,

	// node errors
	"node::devidexists":          ErrNodeDeviceIDExists,
	"node::childexists":          ErrNodeChildExists,
	"node::uniquealiasviolation": ErrNodeUniqueAliasViolation,
	"node::uniquechildviolation": ErrNodeUniqueChildViolation,
	"node::childnotfound":        ErrNodeChildNotFound,
	"node::invalidparentid":      ErrNodeInvalidParentID,
	"node::invalidchildid":       ErrNodeInvalidChildID,
	"node::nodenotfound":         ErrNodeNotFound,
	"node::devicenotfound":       ErrNodeDeviceNotFound,
	"node::invaliduuid":          ErrNodeInvalidUUID,
	"data::invalidorder":         ErrInvalidOrder,

	// Data API errors
	"data::variablenotfound": ErrVariableNotFound,
	"data::invalidnodeid":    ErrInvalidNodeID,

	// Device Logs API errors (verified from Postman)
	"logs::invalidnodeid":    ErrInvalidNodeID,
	"logs::invalidtimerange": ErrInvalidTimeRange,

	// variable errors
	"variable::namerequired":     ErrVariableNameRequired,
	"variable::variablerequired": ErrVariableRequired,
	"variable::typerequired":     ErrVariableTypeRequired,

	// Health API errors
	"health::limitexceeded": ErrHealthLimitExceeded,

	// Aggregation API errors
	"aggregates::invalidaggregationmethod": ErrInvalidAggregationMethod,
	"aggregates::missingaggregationmethod": ErrAggregationMethodRequired,
	"aggregates::invalidintervalmeasure":   ErrInvalidIntervalMeasure,
	"aggregates::invalidinterval":          ErrInvalidInterval,
	"aggregates::invalidtimezone":          ErrInvalidTimezone,
	"aggregates::invalidfiltertype":        ErrInvalidFilterType,
	"aggregates::invalidtimerange":         ErrInvalidTimeRange,

	// accesstoken errors
	"fa::invalidexpiry": ErrExpiryRequried,
	"fa::tokennofound":  ErrInvalidToken,

	// Commands API errors
	"cmd::invalidexpiry": ErrInvalidExpiry,
	"cmd::unknowntype":   ErrInvalidCommandType,
	"cmd::notfound":      ErrCommandNotFound,
	"cmd::invalidfilter": ErrInvalidFilter,
	"cmd::cmdidtoolong":  ErrInvalidCommandID,

	// valuestore errors
	"vs::invalidscope":   ErrInvalidNamespaceScope,
	"vs::invalidtype":    ErrInvalidValueType,
	"vs::invalidns":      ErrInvalidNamespaceScope,
	"vs::invalidrequest": ErrInvalidInput,
	"vs::keynotfound":    ErrKeyNotFound,
}

// GetError converts an API reason code and message into an AnedyaError.
func GetError(code, message string) error {
	sentinel, ok := codeMap[code]
	if !ok {
		sentinel = ErrUnknown
	}

	return &AnedyaError{
		Message: message,
		Err:     sentinel,
	}
}
