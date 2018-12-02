package sess

import "sync"

type memoryStore struct {
	sessions map[string]map[string]interface{}
	mutex    *sync.Mutex
}

// MemoryStore creates a new memory session store
func MemoryStore() *memoryStore {
	ms := memoryStore{
		mutex: &sync.Mutex{},
	}
	ms.sessions = make(map[string]map[string]interface{})
	return &ms
}

func (m *memoryStore) Create(id string, data *map[string]interface{}) error {
	m.mutex.Lock()
	m.sessions[id] = *data
	m.mutex.Unlock()
	return nil
}

func (m *memoryStore) Get(id string) (map[string]interface{}, error) {
	m.mutex.Lock()
	session, ok := m.sessions[id]
	m.mutex.Unlock()
	if !ok {
		return nil, ErrSessionNonexistant
	}

	return session, nil
}

func (m *memoryStore) Update(id string, data *map[string]interface{}) error {
	err := m.Create(id, data)
	return err
}

func (m *memoryStore) Destroy(id string) error {
	m.mutex.Lock()
	delete(m.sessions, id)
	m.mutex.Unlock()
	return nil
}
