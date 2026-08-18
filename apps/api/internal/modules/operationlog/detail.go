package operationlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// DetailSchemaVersion is the version of newly written structured operation
// details. Legacy detail strings remain readable but are not parsed as this
// envelope.
const DetailSchemaVersion = 1

// RedactedValue is used for values that must never be persisted in an audit
// detail. Keeping the marker stable makes downstream redaction explicit.
const RedactedValue = "[REDACTED]"

// DetailChange represents one field's before/after value. Omitted sides mean
// that the field did not exist on that side of the mutation.
type DetailChange struct {
	Before any `json:"before,omitempty"`
	After  any `json:"after,omitempty"`
}

// DetailEnvelope is the versioned structured detail contract for new writes.
// Before/after are optional snapshots; diff contains only changed fields.
type DetailEnvelope struct {
	SchemaVersion int                     `json:"schemaVersion"`
	Action        string                  `json:"action"`
	Before        map[string]any          `json:"before,omitempty"`
	After         map[string]any          `json:"after,omitempty"`
	Diff          map[string]DetailChange `json:"diff"`
}

// NewDetail builds a versioned, recursively redacted detail envelope.
func NewDetail(action string, before, after map[string]any) (string, error) {
	if strings.TrimSpace(action) == "" {
		return "", errors.New("operationlog: detail action is required")
	}
	cleanBefore, err := redactObject(before)
	if err != nil {
		return "", fmt.Errorf("operationlog: redact before detail: %w", err)
	}
	cleanAfter, err := redactObject(after)
	if err != nil {
		return "", fmt.Errorf("operationlog: redact after detail: %w", err)
	}

	diff := make(map[string]DetailChange)
	keys := make(map[string]struct{}, len(before)+len(after))
	for key := range before {
		keys[key] = struct{}{}
	}
	for key := range after {
		keys[key] = struct{}{}
	}
	for key := range keys {
		rawBefore, beforeOK := before[key]
		rawAfter, afterOK := after[key]
		if beforeOK && afterOK && reflect.DeepEqual(rawBefore, rawAfter) {
			continue
		}
		change := DetailChange{}
		if beforeOK {
			change.Before = cleanBefore[key]
		}
		if afterOK {
			change.After = cleanAfter[key]
		}
		diff[key] = change
	}

	envelope := DetailEnvelope{
		SchemaVersion: DetailSchemaVersion,
		Action:        strings.TrimSpace(action),
		Before:        cleanBefore,
		After:         cleanAfter,
		Diff:          diff,
	}
	raw, err := json.Marshal(envelope)
	if err != nil {
		return "", fmt.Errorf("operationlog: marshal detail: %w", err)
	}
	return string(raw), nil
}

// ParseDetail validates and decodes a structured detail envelope. Legacy
// detail strings intentionally return an error so callers can keep them on a
// compatibility path instead of guessing their schema.
func ParseDetail(raw string) (DetailEnvelope, error) {
	var envelope DetailEnvelope
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return DetailEnvelope{}, err
	}
	if envelope.SchemaVersion != DetailSchemaVersion {
		return DetailEnvelope{}, fmt.Errorf("operationlog: unsupported detail schema version %d", envelope.SchemaVersion)
	}
	if strings.TrimSpace(envelope.Action) == "" || envelope.Diff == nil {
		return DetailEnvelope{}, errors.New("operationlog: invalid detail envelope")
	}
	return envelope, nil
}

func redactObject(input map[string]any) (map[string]any, error) {
	if input == nil {
		return nil, nil
	}
	value, err := redactValue(input, "")
	if err != nil {
		return nil, err
	}
	result, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("operationlog: detail object is not a map")
	}
	return result, nil
}

func redactValue(value any, key string) (any, error) {
	if isSensitiveKey(key) {
		return RedactedValue, nil
	}
	switch typed := value.(type) {
	case nil, string, bool, float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return typed, nil
	case map[string]any:
		out := make(map[string]any, len(typed))
		for childKey, childValue := range typed {
			clean, err := redactValue(childValue, childKey)
			if err != nil {
				return nil, err
			}
			out[childKey] = clean
		}
		return out, nil
	case []any:
		out := make([]any, len(typed))
		for index, childValue := range typed {
			clean, err := redactValue(childValue, "")
			if err != nil {
				return nil, err
			}
			out[index] = clean
		}
		return out, nil
	case []string:
		out := make([]any, len(typed))
		for index, childValue := range typed {
			out[index] = childValue
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported detail value type %T", value)
	}
}

func isSensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.NewReplacer("_", "", "-", "", ".", "").Replace(key))
	if normalized == "" {
		return false
	}
	for _, fragment := range []string{
		"password", "accesstoken", "refreshtoken", "authorization", "cookie",
		"secret", "privatekey", "clientsecret", "apikey", "credential",
		"recoverycode", "recoverycodes", "otpauth", "one time password",
	} {
		if strings.Contains(normalized, strings.ReplaceAll(fragment, " ", "")) {
			return true
		}
	}
	if normalized == "token" || strings.HasSuffix(normalized, "token") || normalized == "code" || normalized == "otp" {
		return true
	}
	return normalized == "url" || strings.HasSuffix(normalized, "url")
}
