package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	// "cloud.google.com/go/memcache/apiv1beta2"
)

func main() {
	http.HandleFunc("/~/", infoHandler)
	http.HandleFunc("/", redirectHandler)

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

func infoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	path := r.URL.Path[2:]

	// Search for the path in the existing linkks.
	key, linkk, err := getLinkkByPath(ctx, path)
	if err != nil {
		writeJSONError(ctx, w, err, "Unable to search for linkk", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(ctx, w, EntityResponse{Key: key, Entity: linkk})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Root path, redirect to edit ui.
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/_/ui/edit/index.html", 302)
	}

	// TODO: Enable caching
	// // Get the linkk from cache if available.
	// if item, err := memcache.Get(ctx, r.URL.Path); err == memcache.ErrCacheMiss {
	// 	// Not found, ignore.
	// } else if err != nil {
	// 	writeJSONError(ctx, w, err, "Unable to check cache for linkk", http.StatusInternalServerError)
	// 	return
	// } else {
	// 	http.Redirect(w, r, string(item.Value), 302)
	// 	return
	// }

	// Search for the path in the existing linkks.
	_, linkk, err := getLinkkByPath(ctx, r.URL.Path)
	if err != nil {
		writeJSONError(ctx, w, err, "Unable to search for linkk", http.StatusInternalServerError)
		return
	}

	// Save to cache and redirect.
	if linkk != nil {
		// TODO: Enable caching
		// item := &memcache.Item{
		// 	Key:   linkk.Path,
		// 	Value: []byte(linkk.URL),
		// }
		// if err := memcache.Set(ctx, item); err != nil {
		// 	return
		// }
		http.Redirect(w, r, linkk.URL, 302)
		return
	}

	// Not found, redirect to page to edit the redirect.
	http.Redirect(w, r, fmt.Sprintf("/_/ui/edit/index.html?path=%s", url.QueryEscape(r.URL.Path)), 302)
}
