package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMissing(t *testing.T) {
	ready := map[string]bool{"ao/slack-mcp": true}
	if got := missing([]string{"ao/slack-mcp"}, ready); len(got) != 0 {
		t.Fatalf("all ready: %v", got)
	}
	got := missing([]string{"ao/hubspot-mcp", "ao/slack-mcp", "ao/stripe-mcp"}, ready)
	if want := []string{"ao/hubspot-mcp", "ao/stripe-mcp"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// Every sealed starter is listed in catalog.json and every package it requires
// is a sealed connector, so the publisher never meets a starter it cannot place.
func TestStartersMatchCatalog(t *testing.T) {
	raw, err := os.ReadFile("../../catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	var cat struct {
		Namespace string `json:"namespace"`
		Starters  []struct {
			ID       string   `json:"id"`
			Requires []string `json:"requires"`
		} `json:"starters"`
	}
	if err := json.Unmarshal(raw, &cat); err != nil {
		t.Fatal(err)
	}
	inCatalog := map[string][]string{}
	for _, s := range cat.Starters {
		inCatalog[cat.Namespace+"/"+s.ID] = s.Requires
	}
	connectors := map[string]bool{}
	files, _ := filepath.Glob("../../index/*.json")
	for _, f := range files {
		ref, err := packageRef(f)
		if err != nil {
			t.Fatal(err)
		}
		connectors[ref] = true
	}
	files, _ = filepath.Glob("../../starters/*.json")
	if len(files) != len(inCatalog) {
		t.Fatalf("%d sealed starters, %d in catalog.json", len(files), len(inCatalog))
	}
	for _, f := range files {
		ref, err := packageRef(f)
		if err != nil {
			t.Fatal(err)
		}
		req, ok := inCatalog[ref]
		if !ok {
			t.Fatalf("%s is not in catalog.json starters", ref)
		}
		for _, r := range req {
			if !connectors[r] {
				t.Fatalf("%s requires %s, which is not in index/", ref, r)
			}
		}
	}
}
