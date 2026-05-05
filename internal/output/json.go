package output

import (
	"encoding/json"
	"io"

	"config-audit/internal/finding"
)

func WriteJSON(writer io.Writer, findings []finding.Finding) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(findings)
}
