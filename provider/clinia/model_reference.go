package clinia

import "fmt"

func buildModelID(name, version string) string {
	return fmt.Sprintf("%s:%s", name, version)
}
