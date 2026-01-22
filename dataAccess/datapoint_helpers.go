package dataAccess

import "encoding/json"

// AsFloat attempts to decode the DataPoint value as a float64.
//
// This method is typically used when the variable represents
// a numeric metric such as temperature, voltage, count, etc.
//
// Returns:
//   - (float64, true) if the value can be successfully decoded.
//   - (0, false) if decoding fails or the underlying value
//     is not a valid float64.
func (dp DataPoint) AsFloat() (float64, bool) {
	var v float64
	if err := json.Unmarshal(dp.Value, &v); err != nil {
		return 0, false
	}
	return v, true
}

// AsGeo attempts to decode the DataPoint value as a GeoValue.
//
// This method is used for variables that represent geographical
// coordinates (latitude and longitude).
//
// A basic sanity check is performed to ensure the decoded
// coordinates are not both zero.
//
// Returns:
//   - (GeoValue, true) if decoding succeeds and values are valid.
//   - (GeoValue{}, false) if decoding fails or the coordinates
//     are invalid.
func (dp DataPoint) AsGeo() (GeoValue, bool) {
	var g GeoValue
	if err := json.Unmarshal(dp.Value, &g); err != nil {
		return GeoValue{}, false
	}

	// basic sanity check for invalid coordinates
	if g.Lat == 0 && g.Long == 0 {
		return GeoValue{}, false
	}

	return g, true
}
