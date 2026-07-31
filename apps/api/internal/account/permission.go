package account

import (
	"encoding/json"
	"errors"
	"regexp"
	"slices"
	"strings"
)

// ErrExpression is returned when a permission expression cannot be parsed.
var ErrExpression = errors.New("invalid permission expression")

// Expression shape mirrors the frozen R3/web contract: only $context.user.* /
// $context.features.* paths, == / != / contains, and string or boolean literals.
var expressionPattern = regexp.MustCompile(
	`^\$context\.(user|features)\.([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)\s+(==|!=|contains)\s+(.+)$`,
)

// Evaluate reports whether expr holds against the given user/features
// snapshot. Undeclared paths fail closed (false), never crash.
func Evaluate(expr string, user User, features map[string]bool) (bool, error) {
	match := expressionPattern.FindStringSubmatch(strings.TrimSpace(expr))
	if match == nil {
		return false, ErrExpression
	}
	rootName, path, operator, literal := match[1], match[2], match[3], strings.TrimSpace(match[4])

	var root json.RawMessage
	switch rootName {
	case "user":
		b, err := json.Marshal(user)
		if err != nil {
			return false, err
		}
		root = b
	case "features":
		b, err := json.Marshal(features)
		if err != nil {
			return false, err
		}
		root = b
	}

	actual, err := lookup(root, path)
	if err != nil {
		return false, err
	}

	switch operator {
	case "==":
		return equal(actual, literal)
	case "!=":
		eq, err := equal(actual, literal)
		return !eq, err
	case "contains":
		return contains(actual, literal)
	}
	return false, ErrExpression
}

// lookup walks a dotted path over JSON. Missing keys return nil with no error;
// walking into a non-object or non-array value returns ErrExpression.
func lookup(root json.RawMessage, path string) (json.RawMessage, error) {
	current := root
	for part := range strings.SplitSeq(path, ".") {
		if len(current) == 0 || string(current) == "null" {
			return nil, nil
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(current, &object); err != nil {
			return nil, ErrExpression
		}
		var ok bool
		current, ok = object[part]
		if !ok {
			return nil, nil
		}
	}
	return current, nil
}

// parseLiteral mirrors the R3 web contract: booleans and strings; strings are
// JSON literals (e.g. "A"), numbers are compared as strings.
func parseLiteral(literal string) (any, error) {
	if literal == "true" {
		return true, nil
	}
	if literal == "false" {
		return false, nil
	}
	var value string
	if err := json.Unmarshal([]byte(literal), &value); err != nil {
		return nil, ErrExpression
	}
	return value, nil
}

// equal compares a JSON value with a parsed literal.
func equal(value json.RawMessage, literal string) (bool, error) {
	if len(value) == 0 || string(value) == "null" {
		return false, nil
	}
	expected, err := parseLiteral(literal)
	if err != nil {
		return false, err
	}
	switch want := expected.(type) {
	case bool:
		var got bool
		if err := json.Unmarshal(value, &got); err != nil {
			return false, nil
		}
		return got == want, nil
	default:
		var got string
		if err := json.Unmarshal(value, &got); err != nil {
			return false, nil
		}
		return got == want, nil
	}
}

// contains reports whether a JSON array contains the parsed literal.
func contains(value json.RawMessage, literal string) (bool, error) {
	expected, err := parseLiteral(literal)
	if err != nil {
		return false, err
	}
	switch want := expected.(type) {
	case bool:
		var values []bool
		if err := json.Unmarshal(value, &values); err != nil {
			return false, nil
		}
		return slices.Contains(values, want), nil
	default:
		var values []string
		if err := json.Unmarshal(value, &values); err != nil {
			return false, nil
		}
		return slices.Contains(values, want.(string)), nil
	}
}

// Allow is the fail-closed authorization entry point for business routes.
// Unparsable expressions deny; undeclared $context paths deny.
func Allow(expr string, user User, features map[string]bool) bool {
	if strings.TrimSpace(expr) == "" {
		return true
	}
	allowed, err := Evaluate(expr, user, features)
	return err == nil && allowed
}
