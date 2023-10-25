package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReverseServer_Direct(t *testing.T) {
	input := "The quick brown fox jumped over the lazy dog"
	rev := Reverse(input)
	doubleRev := Reverse(rev)

	if input != doubleRev {
		t.Errorf("expected %q to be equal to %q", rev, doubleRev)
		t.Fail()
	}
}

func TestReverseServer_Handler(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"empty", ""},
		{"ascii", "The quick brown fox jumped over the lazy dog"},
		{"non-ascii", "日本語"},
		{"another", "�"},
		{"fails", "\xe6\x83"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("text: " + test.text)
			path := "/reverse"
			reverseRequest := ReverseRequest{Param: test.text}
			requestBytes, err := json.Marshal(reverseRequest)
			if err != nil {
				t.Fatal(err)
				return
			}
			req, err := http.NewRequest("POST", path, bytes.NewReader(requestBytes))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ReverseHandler)
			handler.ServeHTTP(rr, req)
			status := rr.Code
			if status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v", status)
				t.Fail()
				return
			}

			var reverseResponse ReverseResponse
			err = json.Unmarshal(rr.Body.Bytes(), &reverseResponse)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				t.Fail()
				return
			}

			if Reverse(test.text) != reverseResponse.Reversed {
				t.Errorf("expected reversed %v to be equal to %v", Reverse(test.text), reverseResponse.Reversed)
				t.Fail()
			}
		})
	}
}

func FuzzReverseServer(f *testing.F) {
	f.Add("")
	f.Add("The quick brown fox jumped over the lazy dog")
	f.Add("日本語")
	f.Add("�")
	f.Fuzz(func(t *testing.T, input string) {
		path := "/reverse"
		reverseRequest := ReverseRequest{Param: input}
		requestBytes, err := json.Marshal(reverseRequest)
		if err != nil {
			t.Fatal(err)
			return
		}
		req, err := http.NewRequest("POST", path, bytes.NewReader(requestBytes))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ReverseHandler)
		handler.ServeHTTP(rr, req)
		status := rr.Code
		if status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v", status)
			t.Fail()
			return
		}

		var reverseResponse ReverseResponse
		err = json.Unmarshal(rr.Body.Bytes(), &reverseResponse)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			t.Fail()
			return
		}

		if Reverse(input) != reverseResponse.Reversed {
			t.Errorf("expected reverse %q to be equal to %q", Reverse(string(input)), reverseResponse.Reversed)
			t.Fail()
		}
	})
}
