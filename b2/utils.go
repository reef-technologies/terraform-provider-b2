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

// orderLifecycleRules reorders lifecycle rules returned by B2 to match the order of
// rules in desired. B2 does not guarantee the order of lifecycle rules in bucket
// responses, which would otherwise cause a permanent diff (and a needless bucket
// revision bump on every apply) for list-typed lifecycle_rules. Rules are matched
// by file name prefix; rules missing from desired (e.g. added out-of-band) keep
// the API order at the end of the list.
func orderLifecycleRules(rules []LifecycleRule, desired []interface{}) []LifecycleRule {
	if len(rules) < 2 || len(desired) == 0 {
		return rules
	}

	consumed := make([]bool, len(rules))
	result := make([]LifecycleRule, 0, len(rules))

	for _, d := range desired {
		wanted, ok := d.(map[string]interface{})
		if !ok {
			continue
		}
		prefix, _ := wanted["file_name_prefix"].(string)

		for i, rule := range rules {
			if !consumed[i] && rule.FileNamePrefix == prefix {
				consumed[i] = true
				result = append(result, rule)
				break
			}
		}
	}

	for i, rule := range rules {
		if !consumed[i] {
			result = append(result, rule)
		}
	}

	return result
}

// suppressMapKeyCaseDiff suppresses diffs of a map attribute caused only by B2 storing keys in lower
// case. Keys listed in serverAddedKeys are the ones B2 adds on its own, so finding them in the state
// but not in the config is not a change. Attributes without such keys must not pass any, as then every
// key missing from the config has been removed from it.
func suppressMapKeyCaseDiff(attribute string, serverAddedKeys ...string) schema.SchemaDiffSuppressFunc {
	prefix := attribute + "."

	serverAdded := make(map[string]struct{}, len(serverAddedKeys))
	for _, key := range serverAddedKeys {
		serverAdded[strings.ToLower(key)] = struct{}{}
	}

	return func(k, old, new string, d *schema.ResourceData) bool {
		key, found := strings.CutPrefix(k, prefix)
		if !found {
			return false
		}

		config, ok := configMapAttribute(d, attribute)
		if !ok {
			return false
		}
		state, _ := d.GetChange(attribute)
		stateMap := stateMapAttribute(state)

		// The element count diff carries key removals that the element diffs do not, so it can only be
		// suppressed once the whole map is known to hold what the config asks for
		if key == "%" {
			return len(serverAdded) > 0 && mapStoresConfig(config, stateMap, serverAdded)
		}

		key = strings.ToLower(key)
		configValue, inConfig := lookupLowerCaseKey(config, key)
		if !inConfig {
			_, isServerAdded := serverAdded[key]
			return isServerAdded
		}

		stateValue, inState := lookupLowerCaseKey(stateMap, key)

		return inState && stateValue == configValue
	}
}

// mapStoresConfig reports whether every configured key is stored with the same value, and whether
// every stored key is either configured or one that B2 adds on its own.
func mapStoresConfig(config, state map[string]string, serverAdded map[string]struct{}) bool {
	for key, configValue := range config {
		stateValue, inState := lookupLowerCaseKey(state, strings.ToLower(key))
		if !inState || stateValue != configValue {
			return false
		}
	}

	for key := range state {
		key = strings.ToLower(key)
		if _, inConfig := lookupLowerCaseKey(config, key); inConfig {
			continue
		}
		if _, isServerAdded := serverAdded[key]; !isServerAdded {
			return false
		}
	}

	return true
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
