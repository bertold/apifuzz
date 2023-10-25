package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type ReverseRequest struct {
	Param string
}

type ReverseResponse struct {
	Original string
	Reversed string
}

func ReverseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "failed to read request"}`))
	}

	var reverseRequest ReverseRequest
	err = json.Unmarshal(req, &reverseRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "failed to unmarshal request"}`))
	}
	response := ReverseResponse{
		Original: reverseRequest.Param,
		Reversed: Reverse(string(reverseRequest.Param)),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "failed to marshal response"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(responseBytes)
	}
}

func Reverse(s string) string {
	return reverse_Fixed(s)
}

func reverse_Buggy(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func reverse_Fixed(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
