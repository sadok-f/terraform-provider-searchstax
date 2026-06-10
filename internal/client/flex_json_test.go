package client

import (
	"encoding/json"
	"testing"
)

func TestFlexStringListUnmarshal(t *testing.T) {
	t.Run("array", func(t *testing.T) {
		var list FlexStringList
		if err := json.Unmarshal([]byte(`["a","b"]`), &list); err != nil {
			t.Fatal(err)
		}
		if len(list) != 2 || list[0] != "a" {
			t.Fatalf("unexpected list: %#v", list)
		}
	})

	t.Run("space delimited string", func(t *testing.T) {
		var list FlexStringList
		if err := json.Unmarshal([]byte(`"a b c"`), &list); err != nil {
			t.Fatal(err)
		}
		if len(list) != 3 || list[2] != "c" {
			t.Fatalf("unexpected list: %#v", list)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		var list FlexStringList
		if err := json.Unmarshal([]byte(`""`), &list); err != nil {
			t.Fatal(err)
		}
		if list != nil {
			t.Fatalf("expected nil list, got %#v", list)
		}
	})
}

func TestPlanUnmarshalApplicationString(t *testing.T) {
	var plan Plan
	if err := json.Unmarshal([]byte(`{"plan":"NDC4-GCP-G","application":"Solr","plan_type":"DedicatedPlan"}`), &plan); err != nil {
		t.Fatal(err)
	}
	if plan.Application != "Solr" || plan.Name != "NDC4-GCP-G" {
		t.Fatalf("unexpected plan: %#v", plan)
	}
}
