package sess

// SessionStore is the interface that should be impimented for any new session stores
type SessionStore interface {

	/*
		Creates a new session
		Shoud overwrite an old session if one exists
	*/
	Create(id string, data *map[string]interface{}) error

	/*
		Updates an existing session
		Update should completely overwrite the old session value
		A session should be created if it does not exist
	*/
	Update(id string, data *map[string]interface{}) error

	/*
		Returns an existing session
		Returns a ErrSessionNonexistant error if the session does not exist
	*/
	Get(id string) (map[string]interface{}, error)

	/*
		Destroys an existing session
	*/
	Destroy(id string) error
}
