package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ExpandedDate struct {
	Year   int `json:"year"`
	Month  int `json:"month"`
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}

func DateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	queryValues := r.URL.Query()
	offset := 0
	unit := 1 * time.Hour

	rawOffset := queryValues.Get("offset")
	if rawOffset != "" {
		offset, err = strconv.Atoi(rawOffset)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("offset must be an integer"))
			return
		}
	}
	rawUnit := queryValues.Get("unit")
	if rawUnit != "" {
		switch rawUnit {
		case "d":
			unit = 24 * time.Hour
		case "h":
			unit = time.Hour
		case "m":
			unit = time.Minute
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid unit"))
			return
		}
	}

	d := time.Now().Add(time.Duration(offset) * unit)
	w.Header().Set("Content-Type", "application/json")
	response, err := formatDate(d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "failed to format date"}`)))
	} else {
		w.Write(response)
	}
}

func formatDate(d time.Time) ([]byte, error) {
	expandedDate := ExpandedDate{
		Year:   d.Year(),
		Month:  int(d.Month()),
		Day:    d.Day(),
		Hour:   d.Hour(),
		Minute: d.Minute(),
		Second: d.Second(),
	}
	return json.Marshal(expandedDate)
}
