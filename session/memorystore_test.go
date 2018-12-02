package sess

import "testing"

func TestGet(t *testing.T) {
	var sess = MemoryStore()

	_, err := sess.Get("dummy-id")
	if err != ErrSessionNonexistant {
		t.Error("Get is not returning an error for empty sessions")
	}
}

func TestCreate(t *testing.T) {
	var sess = MemoryStore()

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

	s := session

	if s["test-data"] != "Wow" {
		t.Errorf("Session Data incorrect, got: %s, want: %s", s["test-data"], "Wow")
	}
}

func TestUpdate(t *testing.T) {
	var sess SessionStore = MemoryStore()

	testData := map[string]interface{}{
		"to_be_deleted":     "pickles",
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

	s := session

	if s["to_be_overwritten"] != "definitely pickles" {
		t.Errorf("Improper update, expected: %s, got: %s", "definitely pickles", s["to_be_overwritten"])
	}

	if _, exists := s["persistent_data"]; exists == true {
		t.Error("Update is not removing old values")
	}
}

func TestUpdateCreate(t *testing.T) {
	var sess SessionStore = MemoryStore()

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

	s := session

	if s["test-data"] != "Wow" {
		t.Errorf("Session Data incorrect, got: %s, want: %s", s["test-data"], "Wow")
	}
}

func TestDestroy(t *testing.T) {
	var sess SessionStore = MemoryStore()

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
