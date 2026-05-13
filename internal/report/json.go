package report

import (
	"encoding/json"

	"siteintel/internal/model"
)

// JSON renders pretty, stable JSON suitable for MCP/API wrapping.
func JSON(result model.AuditResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}
