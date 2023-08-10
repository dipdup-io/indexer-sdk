package duplicator

import "fmt"

// GetInputName -
func GetInputName(idx int) string {
	return fmt.Sprintf("%s_%d", InputName, idx)
}

// GetOutputName -
func GetOutputName(idx int) string {
	return fmt.Sprintf("%s_%d", OutputName, idx)
}
