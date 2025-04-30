package tools

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type baseEnvVar struct {
	Error    error
	RawValue string
}

// IntEnvVar переменная окружения типа int.
type IntEnvVar struct {
	baseEnvVar

	Exists bool
	Value  int
}

// GetIntFromEnv достаёт переменную окружения типа int.
func GetIntFromEnv(name string) IntEnvVar {
	var (
		exists bool
		v      IntEnvVar
	)

	v.RawValue, exists = os.LookupEnv(name)
	if !exists {
		return v
	}

	v.Exists = true

	v.Value, v.Error = strconv.Atoi(v.RawValue)

	return v
}

// StrEnvVar переменная окружения типа string.
type StrEnvVar struct {
	baseEnvVar

	Exists bool
	Value  string
}

// GetStrFromEnv достаёт переменную окружения типа string.
func GetStrFromEnv(name string) StrEnvVar {
	var (
		exists bool
		v      StrEnvVar
	)

	v.RawValue, exists = os.LookupEnv(name)
	if !exists {
		return v
	}

	v.Value = v.RawValue

	return v
}

// BoolEnvVar переменная окружения типа bool.
type BoolEnvVar struct {
	baseEnvVar

	Exists bool
	Value  bool
}

// GetBoolFromEnv достаёт переменную окружения типа bool.
func GetBoolFromEnv(name string) BoolEnvVar {
	var (
		exists bool
		v      BoolEnvVar
	)

	v.RawValue, exists = os.LookupEnv(name)
	if !exists {
		return v
	}

	v.Exists = true

	switch strings.TrimSpace(strings.ToLower(v.RawValue)) {
	case "true":
		v.Value = true
	case "false":
		v.Value = false
	default:
		v.Error = fmt.Errorf("invalid value for boolean variable: %q", v.RawValue)
	}

	return v
}
