package machine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Machine struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

func (m *Machine) Read(file os.FileInfo) error {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/config.json", 
		"/home/pwhelan/.config/dm/machines", file.Name()))
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
