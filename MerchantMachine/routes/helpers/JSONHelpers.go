package helpers

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"net/http"
)

const maxBodyBytes = 1 << 20 

var unmarshalOpts = json.JoinOptions(
	json.RejectUnknownMembers(true),       // fail on unknown fields
	json.MatchCaseInsensitiveNames(false), // keep strict, case-sensitive names
)
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
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

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	for _, h := range headers {
		for k, v := range h {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil{
	return json.MarshalWrite(
		w,
		data,
		json.OmitZeroStructFields(true), // drop zero-value struct fields
		// json.Deterministic(true),      // enable only if you need stable map order (slower)
	)
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error{
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
