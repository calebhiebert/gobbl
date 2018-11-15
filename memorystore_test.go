package gbl

import "testing"

func TestGet(t *testing.T) {
	var sess SessionStore = CreateMemoryStore()

	_, err := sess.Get("dummy-id")
	if err != ErrSessionNonexistant {
		t.Error("Get is not returning an error for empty sessions")
	}
}

func TestCreate(t *testing.T) {
	var sess SessionStore = CreateMemoryStore()

	testData := map[string]interface{}{
		"test-data": "Wow",
	}

	err := sess.Create("test-id", &testData)
	if err != nil {
		t.Error("Received error on session creation")
	}

	session, err := sess.Get("test-id")
	if err != nil {
		t.Error("Received error on session retrieval")
	}

	s := *session

	if s["test-data"] != "Wow" {
		t.Errorf("Session Data incorrect, got: %s, want: %s", s["test-data"], "Wow")
	}
}

func TestUpdate(t *testing.T) {
	var sess SessionStore = CreateMemoryStore()

	testData := map[string]interface{}{
		"persistent_data":   "pickles",
		"to_be_overwritten": "not pickles",
	}

	err := sess.Create("test-id", &testData)
	if err != nil {
		t.Error("Received error on session creation")
	}

	updatedTestData := map[string]interface{}{
		"to_be_overwritten": "definitely pickles",
	}

	err = sess.Update("test-id", &updatedTestData)
	if err != nil {
		t.Error("Received error on session updating")
	}

	session, err := sess.Get("test-id")
	if err != nil {
		t.Error("Received error on session retrieval")
	}

	s := *session

	if s["to_be_overwritten"] != "definitely pickles" {
		t.Errorf("Improper update, expected: %s, got: %s", "definitely pickles", s["to_be_overwritten"])
	}

	if s["persistent_data"] != "pickles" {
		t.Errorf("Update not leaving old values intact, wanted: %s, got: %s", "pickles", s["persistent_data"])
	}
}

func TestUpdateCreate(t *testing.T) {
	var sess SessionStore = CreateMemoryStore()

	testData := map[string]interface{}{
		"test-data": "Wow",
	}

	err := sess.Update("test-id", &testData)
	if err != nil {
		t.Error("Received error on session creation")
	}

	session, err := sess.Get("test-id")
	if err != nil {
		if err != ErrSessionNonexistant {
			t.Error("Received error on session retrieval")
		}

		return
	}

	s := *session

	if s["test-data"] != "Wow" {
		t.Errorf("Session Data incorrect, got: %s, want: %s", s["test-data"], "Wow")
	}
}

func TestDestroy(t *testing.T) {
	var sess SessionStore = CreateMemoryStore()

	testData := map[string]interface{}{
		"test-data": "Wow",
	}

	err := sess.Create("test-id", &testData)
	if err != nil {
		t.Error("Received error on session creation")
	}

	err = sess.Destroy("test-id")
	if err != nil {
		t.Error("Received error on session destruction")
	}

	_, err = sess.Get("test-id")
	if err != ErrSessionNonexistant {
		t.Errorf("Session get should have returned ErrSessionNonexistant, instead got %+v", err)
	}
}
