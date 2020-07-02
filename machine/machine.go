package machine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"github.com/kirsle/configdir"
)

type Machine struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

func (m *Machine) Read(file os.FileInfo) error {
	configPath := configdir.LocalConfig("dm")
	configFile := filepath.Join(configPath, "machines", file.Name(), "config.json")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	if m.Name == "" {
		m.Name = file.Name()
	}
	return nil
}
