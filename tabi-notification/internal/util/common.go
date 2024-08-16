package util

import (
	"encoding/json"
	"os"
)

// ReadFile reads file and returns map[string]interface{}
func ReadFile(filename string) (map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var out map[string]interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// WriteFile writes content to file
func WriteFile(filename string, content map[string]interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(content)
	if err != nil {
		return err
	}

	return nil
}

// InterfaceToArrayString converts []interface{} to []string
func InterfaceToArrayString(in []interface{}) []string {
	out := []string{}
	for _, v := range in {
		out = append(out, v.(string))
	}
	return out
}

// TernaryOperator returns trueVal if condition is true, otherwise returns falseVal
func TernaryOperator(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}

	return falseVal
}
