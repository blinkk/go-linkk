package main

import (
	"context"
	"log"
  "strings"

	"cloud.google.com/go/datastore"
)

func getLinkkByPath(ctx context.Context, path string) (key *datastore.Key, linkk *Linkk, err error) {
	dsClient, err := datastore.NewClient(ctx, "blinkk-linkk")
	if err != nil {
		log.Fatal(err)
	}

	path = strings.ToLower(path)

	var entities []Linkk
	q := datastore.NewQuery("Linkk").Filter("Path =", path).Limit(1)
	if _, err := dsClient.GetAll(ctx, q, &entities); err != nil {
		log.Fatal(err)
	}

	if len(entities) > 0 {
		return key, &entities[0], nil
	}

	// No linkk found.
	return nil, nil, nil
}
