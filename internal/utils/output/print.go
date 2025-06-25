package output

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

func SprintPrettyJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, " ", "  ")
	if err != nil {
		return fmt.Sprintf("output error: %v", err)
	}
	return string(data)
}

func PrintlnPrettyJSON(v interface{}) (n int, err error) {
	text := SprintPrettyJSON(v)
	return fmt.Println(text)
}

func SprintPrettyYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Sprintf("output error: %v", err)
	}
	return string(data)
}

func PrintlnPrettyYAML(v interface{}) (n int, err error) {
	text := SprintPrettyYAML(v)
	return fmt.Println(text)
}
