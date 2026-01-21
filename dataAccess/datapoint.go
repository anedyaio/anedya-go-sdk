package dataAccess

import "encoding/json"

// GeoValue represents a geographical coordinate.
//
// It contains latitude and longitude values and is typically
// used when a variable stores location-based data.
type GeoValue struct {
	Lat  float64 `json:"lat"`  // Latitude
	Long float64 `json:"long"` // Longitude
}

// DataPoint represents a single data value recorded at a specific timestamp.
//
// The Value field is stored as raw JSON to support multiple
// data types, such as:
//   - numeric values (float, int)
//   - structured objects (e.g., GeoValue)
//
// Helper methods like AsFloat and AsGeo can be used to
// safely decode the value into the expected type.
type DataPoint struct {
	Timestamp int64           `json:"timestamp"` // Unix timestamp in milliseconds
	Value     json.RawMessage `json:"value"`     // Raw JSON-encoded value
}
