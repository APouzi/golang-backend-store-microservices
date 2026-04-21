package helpers

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"net/http"
)

const maxBodyBytes = 1 << 20 

//  This constant is set using a bitwise left shift operation: 1 << 20. In binary terms, shifting the number 1 left by 20 places results in the value 1,048,576. This is equivalent to 2 raised to the 20th power.In practical terms, maxBodyBytes is often used to specify a size limit for data, such as the maximum number of bytes allowed in an HTTP request body. Setting this limit helps prevent excessive memory usage or potential denial-of-service attacks caused by very large payloads. In this case, the value represents 1 megabyte (MB), which is a common default for request size limits in web applications.

// Using a bitwise shift for powers of two is a concise and efficient way to express such values in code, and it makes the intent clear to developers familiar with binary operations.

// Strict semantic options for v2's Unmarshal.
var unmarshalOpts = json.JoinOptions(
	json.RejectUnknownMembers(true),       // fail on unknown fields
	json.MatchCaseInsensitiveNames(false), // keep strict, case-sensitive names
)

// ReadJSON reads exactly one JSON value into data with a hard body limit.
// Uses v2's streaming decoder variant; no []byte staging.
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	// v2 defaults already reject duplicate names & invalid UTF-8;
	// we pass them explicitly for clarity (safe defaults).
	if err := json.UnmarshalRead(
		r.Body,
		data,
		unmarshalOpts,
		jsontext.AllowDuplicateNames(false),
		jsontext.AllowInvalidUTF8(false),
	); err != nil {
		return err
	}
	return nil
}

// WriteJSON streams JSON directly to the ResponseWriter.
// This avoids an intermediate []byte and is the fastest path.
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	for _, h := range headers {
		for k, v := range h {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	// Stream out; choose options as needed.
	return json.MarshalWrite(
		w,
		data,
		json.OmitZeroStructFields(true), // drop zero-value struct fields
		// json.Deterministic(true),      // enable only if you need stable map order (slower)
	)
}

// ErrorJSON writes a standard error envelope.
func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	code := http.StatusBadRequest
	if len(status) > 0 {
		code = status[0]
	}

	payload := ErrorJSONResponse{
		Error:   true,
		Message: err.Error(),
	}
	return WriteJSON(w, code, payload,
		http.Header{"Cache-Control": []string{"no-store"}},
	)
}
