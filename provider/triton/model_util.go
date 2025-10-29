package triton

import (
	"fmt"
	"strings"
)

// splitModelID parses a model identifier of the form "name:version" and validates
// that both parts are present and non-empty. The component string is used only
// to contextualize error messages (e.g., "clinia/embed").
func splitModelID(component, id string) (name, version string, err error) {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return "", "", fmt.Errorf("clinia/%s: model id is required (expected 'name:version')", component)
	}

	parts := strings.SplitN(trimmed, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("clinia/%s: model version is required in id (expected 'name:version')", component)
	}
	name = strings.TrimSpace(parts[0])
	version = strings.TrimSpace(parts[1])
	if name == "" {
		return "", "", fmt.Errorf("clinia/%s: model name is required", component)
	}
	if version == "" {
		return "", "", fmt.Errorf("clinia/%s: model version is required", component)
	}
	return name, version, nil
}

// validateNameVersion trims and validates name and version provided separately.
// The component string is used only for clear error messages.
func validateNameVersion(component, name, version string) (string, string, error) {
	n := strings.TrimSpace(name)
	if n == "" {
		return "", "", fmt.Errorf("clinia/%s: model name is required", component)
	}
	v := strings.TrimSpace(version)
	if v == "" {
		return "", "", fmt.Errorf("clinia/%s: model version is required", component)
	}
	return n, v, nil
}
