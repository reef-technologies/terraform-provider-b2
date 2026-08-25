// ####################################################################
//
// File: b2/utils.go
//
// Copyright 2024 Backblaze Inc. All Rights Reserved.
//
// License https://www.backblaze.com/using_b2_code.html
//
// ####################################################################

package b2

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// suppressMapKeyCaseDiff suppresses diffs of a map attribute caused only by B2 storing keys in lower case.
func suppressMapKeyCaseDiff(attribute string) schema.SchemaDiffSuppressFunc {
	prefix := attribute + "."

	return func(k, old, new string, d *schema.ResourceData) bool {
		key, found := strings.CutPrefix(k, prefix)
		if !found || key == "%" {
			return false
		}

		config, ok := configMapAttribute(d, attribute)
		if !ok {
			return false
		}
		state, _ := d.GetChange(attribute)

		key = strings.ToLower(key)
		stateValue, inState := lookupLowerCaseKey(stateMapAttribute(state), key)
		configValue, inConfig := lookupLowerCaseKey(config, key)

		return inState && inConfig && stateValue == configValue
	}
}

// configMapAttribute reads a map attribute from the raw config, so that values missing from the
// config are not read from the state instead. It reports false when the value cannot be compared.
func configMapAttribute(d *schema.ResourceData, attribute string) (map[string]string, bool) {
	config := d.GetRawConfig()
	if config.IsNull() || !config.Type().IsObjectType() || !config.Type().HasAttribute(attribute) {
		return nil, false
	}

	value := config.GetAttr(attribute)
	if !value.IsKnown() {
		return nil, false
	}
	if value.IsNull() {
		return map[string]string{}, true
	}

	result := make(map[string]string, value.LengthInt())
	for it := value.ElementIterator(); it.Next(); {
		elemKey, elemValue := it.Element()
		if elemValue.IsNull() || !elemValue.IsKnown() {
			return nil, false
		}
		result[elemKey.AsString()] = elemValue.AsString()
	}

	return result, true
}

// stateMapAttribute converts a map attribute read from the state. Values that are not strings
// are dropped, which makes their key look absent and leaves its diff in place.
func stateMapAttribute(state interface{}) map[string]string {
	m, _ := state.(map[string]interface{})

	result := make(map[string]string, len(m))
	for key, value := range m {
		if s, ok := value.(string); ok {
			result[key] = s
		}
	}

	return result
}

// lookupLowerCaseKey returns the value stored under the key whose lower case form is the given one.
func lookupLowerCaseKey(m map[string]string, lowerCaseKey string) (string, bool) {
	for key, value := range m {
		if strings.ToLower(key) == lowerCaseKey {
			return value, true
		}
	}

	return "", false
}
