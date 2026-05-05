package config

import (
	"fmt"
	"strconv"
	"strings"
)

type Document struct {
	Source string
	Format Format
	Root   any
}

type Node struct {
	Path  string
	Key   string
	Value any
}

func (d Document) Walk(fn func(Node)) {
	walkValue(d.Root, "", "", fn)
}

func ScalarString(value any) (string, bool) {
	switch typed := value.(type) {
	case string:
		return typed, true
	case fmt.Stringer:
		return typed.String(), true
	default:
		return "", false
	}
}

func ScalarBool(value any) (bool, bool) {
	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		switch strings.ToLower(strings.TrimSpace(typed)) {
		case "true", "yes", "y", "1", "on", "enabled":
			return true, true
		case "false", "no", "n", "0", "off", "disabled":
			return false, true
		default:
			return false, false
		}
	default:
		return false, false
	}
}

func walkValue(value any, path string, key string, fn func(Node)) {
	fn(Node{
		Path:  path,
		Key:   key,
		Value: value,
	})

	switch typed := value.(type) {
	case map[string]any:
		for childKey, childValue := range typed {
			walkValue(childValue, joinPath(path, childKey), childKey, fn)
		}
	case []any:
		for i, childValue := range typed {
			index := "[" + strconv.Itoa(i) + "]"
			walkValue(childValue, path+index, index, fn)
		}
	}
}

func joinPath(parent string, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}
