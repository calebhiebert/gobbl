package main

type MemoryStore struct {
	sessions map[string]map[string]interface{}
}

func (m MemoryStore) Create(id string, data *map[string]interface{}) error {
	m.sessions[id] = *data
	return nil
}

func (m MemoryStore) Get(id string) (*map[string]interface{}, error) {
	session := m.sessions[id]

	return &session, nil
}

func (m MemoryStore) Update(id string, data *map[string]interface{}) error {

	existingSession, err := m.Get(id)
	if err != nil {
		return err
	}

	if existingSession != nil {
		newSession := make(map[string]interface{})

		// Insert all the old session variables into the new session map
		for k, v := range *existingSession {
			newSession[k] = v
		}

		// Insert all the new session variables into the new session map
		for k, v := range *data {
			newSession[k] = v
		}

		m.sessions[id] = newSession

	} else {
		return m.Create(id, data)
	}

	return nil

}

func (m MemoryStore) Destroy(id string) error {
	m.sessions[id] = nil
	return nil
}
