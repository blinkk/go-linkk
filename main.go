package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func init() {
	http.HandleFunc("/_/api/create", apiCreateHandler)
	http.HandleFunc("/~/", infoHandler)
	http.HandleFunc("/", redirectHandler)
}

func apiCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	err := authUserDomain(u)
	if err != nil {
		writeJSONError(ctx, w, err, "", http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" {
		// TODO: Check for existing linkk with same path.
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

		key := datastore.NewIncompleteKey(ctx, "Linkk", nil)
		if _, err := datastore.Put(ctx, key, linkk); err != nil {
			writeJSONError(ctx, w, err, "Unable to store new linkk", http.StatusInternalServerError)
			return
		}
	}
}

func authUserDomain(u *user.User) error {
	if !appengine.IsDevAppServer() && u.AuthDomain != "gmail.com" {
		return fmt.Errorf("Invalid auth domain set for authorization: %s", u.AuthDomain)
	}

	domains := getAuthDomains()

	if appengine.IsDevAppServer() {
		domains = append(domains, "example.com")
	}

	if len(domains) == 0 {
		return errors.New("No auth domains configured for authorization")
	}

	r, _ := regexp.Compile("@(.*)$")
	domain := r.FindStringSubmatch(u.Email)

	if !stringInSlice(domain[1], domains) {
		return fmt.Errorf("Invalid authorization domain: %s", domain[1])
	}

	return nil
}

func getAuthDomains() []string {
	return strings.Split(os.Getenv("AUTH_DOMAINS"), "|")
}

func getLinkkByPath(ctx context.Context, path string) (key *datastore.Key, linkk *Linkk, err error) {
	path = strings.ToLower(path)
	q := datastore.NewQuery("Linkk").Filter("Path =", path)
	t := q.Run(ctx)
	for {
		var linkk Linkk
		key, err := t.Next(&linkk)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		return key, &linkk, nil
	}
	// No linkk found.
	return nil, nil, nil
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
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
	ctx := appengine.NewContext(r)

	// Root path, redirect to create ui.
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/_/ui/create/index.html", 301)
	}

	// Get the linkk from cache if available.
	if item, err := memcache.Get(ctx, r.URL.Path); err == memcache.ErrCacheMiss {
		// Not found, ignore.
	} else if err != nil {
		writeJSONError(ctx, w, err, "Unable to check cache for linkk", http.StatusInternalServerError)
		return
	} else {
		http.Redirect(w, r, string(item.Value), 301)
		return
	}

	// Search for the path in the existing linkks.
	_, linkk, err := getLinkkByPath(ctx, r.URL.Path)
	if err != nil {
		writeJSONError(ctx, w, err, "Unable to search for linkk", http.StatusInternalServerError)
		return
	}

	// Save to caceh and redirect.
	if linkk != nil {
		item := &memcache.Item{
			Key:   linkk.Path,
			Value: []byte(linkk.URL),
		}
		if err := memcache.Set(ctx, item); err != nil {
			writeJSONError(ctx, w, err, "Unable to cache linkk", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, linkk.URL, 301)
		return
	}

	// Not found, redirect to page to create the redirect.
	http.Redirect(w, r, fmt.Sprintf("/_/ui/create/index.html?path=%s", url.QueryEscape(r.URL.Path)), 301)
}
