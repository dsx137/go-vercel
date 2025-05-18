package vercelkit

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func CamelToSnake(s string) string {
	result := strings.Builder{}
	for i, r := range []rune(s) {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func ConvertBytesToHuman(bytes int) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	if bytes == 0 {
		return "0B"
	}

	floatBytes := float64(bytes)

	var unitIndex int
	for floatBytes >= 1024 && unitIndex < len(units)-1 {
		floatBytes = floatBytes / 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f%s", floatBytes, units[unitIndex])
}

func ReadFrom[T any](reader io.ReadCloser) (*T, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	obj := new(T)
	err = json.Unmarshal(body, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
