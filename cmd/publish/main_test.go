package main

import (
	"encoding/json"
	"os"
	"testing"
)

// "Not verified is not published": the rule is read from catalog.json, and a
// connector without the field counts as not verified.
func TestVerifiedSelection(t *testing.T) {
	if verified(obj{}) || verified(obj{"verified": "true"}) || verified(obj{"verified": false}) || !verified(obj{"verified": true}) {
		t.Fatal("verified must be true only for boolean true")
	}
	raw, err := os.ReadFile("../../catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	var cat struct{ Connectors []obj }
	if err := json.Unmarshal(raw, &cat); err != nil {
		t.Fatal(err)
	}
	var held []string
	for _, c := range cat.Connectors {
		if !verified(c) {
			held = append(held, str(c, "id"))
		}
	}
	if len(held) != 2 || held[0] != "ao/rdstation-mcp" || held[1] != "ao/nuvemshop-mcp" {
		t.Fatalf("unverified connectors = %v", held)
	}
}
