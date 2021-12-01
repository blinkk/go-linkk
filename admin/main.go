package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	// "net/url"
	"os"
	"regexp"
	"strings"

	"cloud.google.com/go/datastore"
)

func main() {
	http.HandleFunc("/_/api/edit", apiCreateHandler)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}

func apiCreateHandler(w http.ResponseWriter, r *http.Request) {
 	ctx := context.Background()

	err := authUserDomain(r)
	if err != nil {
		writeJSONError(ctx, w, err, "", http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" {
		dsClient, err := datastore.NewClient(ctx, "blinkk-linkk")
		if err != nil {
			log.Fatal(err)
		}

		linkk := &Linkk{
			Path:    r.FormValue("path"),
			URL:     r.FormValue("url"),
			Comment: r.FormValue("comment"),
		}
		linkk.Clean()
		err = linkk.Validate()
		if err != nil {
			writeJSONError(ctx, w, err, "", http.StatusInternalServerError)
			return
		}

		// Test for existing object to overwrite.
		key, _, err := getLinkkByPath(ctx, linkk.Path)
		if err != nil {
			writeJSONError(ctx, w, err, "Unable to search for linkk", http.StatusInternalServerError)
			return
		}
		if key == nil {
			key = datastore.NameKey("Linkk", "stringID", nil)
		}
		if _, err := dsClient.Put(ctx, key, linkk); err != nil {
			writeJSONError(ctx, w, err, "Unable to store new linkk", http.StatusInternalServerError)
			return
		}

		writeJSONResponse(ctx, w, EntityResponse{Key: key, Entity: linkk})
	}
}

func authUserDomain(r *http.Request) error {
	authInfo := r.Header.Get("X-Goog-Authenticated-User-Email")
	authEmailParts := strings.Split(authInfo, ":")
	authEmail := authEmailParts[1]
	domains := getAuthDomains()

	// TODO: Handle local dev.
	// if appengine.IsDevAppServer() {
	// 	domains = append(domains, "example.com")
	// }

	if len(domains) == 0 {
		return errors.New("No auth domains configured for authorization")
	}

	domainRegex, _ := regexp.Compile("@(.*)$")
	domain := domainRegex.FindStringSubmatch(authEmail)

	if !stringInSlice(domain[1], domains) {
		return fmt.Errorf("Invalid authorization domain: %s", domain[1])
	}

	return nil
}

func getAuthDomains() []string {
	return strings.Split(os.Getenv("AUTH_DOMAINS"), "|")
}
