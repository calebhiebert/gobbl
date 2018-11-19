package sess

type memoryStore struct {
	sessions map[string]map[string]interface{}
}

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
	return m.Create(id, data)
}

func (m memoryStore) Destroy(id string) error {
	delete(m.sessions, id)
	return nil
}
