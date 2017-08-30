package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// DeleteResponse Structure of returning a delete request.
type DeleteResponse struct {
	Key *datastore.Key `json:"key"`
}

// EntityResponse Structure of returning an entity from a request.
type EntityResponse struct {
	Key    *datastore.Key `json:"key"`
	Entity interface{}    `json:"entity"`
}

// ListResponse Structure of returning multiple entities from a request.
type ListResponse struct {
	Items []EntityResponse `json:"items"`
}

// Determine if a string is in a list of strings.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Determine if a string has one of the prefixes in the list of strings.
func stringPrefixInSlice(a string, list []string) string {
	for _, b := range list {
		if strings.HasPrefix(a, b) {
			return b
		}
	}
	return ""
}

// Utility for writing a response in a json format.
func writeJSONResponse(ctx context.Context, w http.ResponseWriter, s interface{}) {
	js, err := json.Marshal(s)
	if err != nil {
		log.Errorf(ctx, "error marshalling: %v", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Utility for writing an error response in a json format.
func writeJSONError(ctx context.Context, w http.ResponseWriter, e error, msg string, code int) {
	var js []byte
	var err error

	log.Errorf(ctx, "error: %v", e)

	if reflect.TypeOf(e).Name() == "ValidationError" {
		js, err = json.Marshal(e)
	} else {
		var message string
		if msg != "" {
			message = msg
		} else {
			message = e.Error()
		}
		js, err = json.Marshal(map[string]interface{}{"Error": message})
	}

	if err != nil {
		log.Errorf(ctx, "error marshalling: %v", err)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}
