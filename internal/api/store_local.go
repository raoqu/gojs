package api

type ScriptLocalStore struct {
	ScriptStore
}

func NewScriptLocalStore() *ScriptLocalStore {
	return &ScriptLocalStore{}
}

func (s *ScriptLocalStore) Load(callback ScriptLoadCallback) {
}

func (s *ScriptLocalStore) Save(scriptName string, scriptCode string) error {
	return nil
}

func (s *ScriptLocalStore) Get(scriptName string) (string, error) {
	return "", nil
}

func (s *ScriptLocalStore) Delete(scriptName string) error {
	return nil
}

func (s *ScriptLocalStore) List() ([]string, error) {
	return []string{}, nil
}

// ScriptExists checks if a script exists in local
func (s *ScriptLocalStore) Exists(scriptName string) (bool, error) {
	return true, nil
}
