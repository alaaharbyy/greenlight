package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Declare a custom Runtime type, which has the underlying type int32 (the same as our // Movie struct field).
type Runtime int32

// SENDING RESPONSEEE

// Implement a MarshalJSON() method on the Runtime type so that it satisfies the
// json.Marshaler interface. This should return the JSON-encoded value for the movie // runtime (in our case, it will return a string in the format "<runtime> mins").
func (r *Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)
	// Use the strconv.Quote() function on the string to wrap it in double quotes. It // needs to be surrounded by double quotes in order to be a valid *JSON string*.
	quotedJSONValue := strconv.Quote(jsonValue)
	// Convert the quoted string value to a byte slice and return it.
	return []byte(quotedJSONValue), nil
}

// READING REQUESTTT

// Define an error that our UnmarshalJSON() method can return if we're unable to parse // or convert the JSON string successfully.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	value, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(value)
	// *r is calling the UnmarshalJSON method
	// overwrite that with just the parsed int32 part but first
	// cast it as Runtime

	return nil
}
