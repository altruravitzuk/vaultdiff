package diff

import "testing"

func TestCompare_NoChanges(t *testing.T) {
	old := map[string]interface{}{"key1": "val1", "key2": "val2"}
	new := map[string]interface{}{"key1": "val1", "key2": "val2"}

	result := Compare(old, new)
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestCompare_Added(t *testing.T) {
	old := map[string]interface{}{}
	new := map[string]interface{}{"newkey": "newval"}

	result := Compare(old, new)
	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(result.Changes) != 1 || result.Changes[0].Type != Added {
		t.Errorf("expected Added change, got %+v", result.Changes)
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]interface{}{"gone": "value"}
	new := map[string]interface{}{}

	result := Compare(old, new)
	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	if result.Changes[0].Type != Removed {
		t.Errorf("expected Removed, got %s", result.Changes[0].Type)
	}
}

func TestCompare_Modified(t *testing.T) {
	old := map[string]interface{}{"key": "old"}
	new := map[string]interface{}{"key": "new"}

	result := Compare(old, new)
	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	if result.Changes[0].Type != Modified {
		t.Errorf("expected Modified, got %s", result.Changes[0].Type)
	}
}

func TestMask_HidesValue(t *testing.T) {
	if mask("supersecret") != "***" {
		t.Error("expected value to be masked")
	}
}

func TestMask_EmptyValue(t *testing.T) {
	if mask("") != "" {
		t.Error("expected empty string for empty input")
	}
}
