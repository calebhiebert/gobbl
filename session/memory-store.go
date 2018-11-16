package sess

import "errors"

type memoryStore struct {
	sessions map[string]map[string]interface{}
}

var ErrSessionNonexistant = errors.New("Session did not exist")

func MemoryStore() memoryStore {
	ms := memoryStore{}
	ms.sessions = make(map[string]map[string]interface{})
	return ms
}

func (m memoryStore) Create(id string, data *map[string]interface{}) error {
	m.sessions[id] = *data
	return nil
}

func (m memoryStore) Get(id string) (map[string]interface{}, error) {
	session, ok := m.sessions[id]
	if !ok {
		return nil, ErrSessionNonexistant
	}

	return session, nil
}

func (m memoryStore) Update(id string, data *map[string]interface{}) error {

	existingSession, err := m.Get(id)
	if err != nil {
		if err == ErrSessionNonexistant {
			return m.Create(id, data)
		} else {
			return err
		}
	}

	if existingSession != nil {
		newSession := make(map[string]interface{})

		// Insert all the old session variables into the new session map
		for k, v := range existingSession {
			newSession[k] = v
		}

		// Insert all the new session variables into the new session map
		for k, v := range *data {
			newSession[k] = v
		}

		m.sessions[id] = newSession

	}

	return nil

}

func (m memoryStore) Destroy(id string) error {
	delete(m.sessions, id)
	return nil
}
