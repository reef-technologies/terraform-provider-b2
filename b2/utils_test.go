//####################################################################
//
// File: b2/utils_test.go
//
// Copyright 2025 Backblaze Inc. All Rights Reserved.
//
// License https://www.backblaze.com/using_b2_code.html
//
//####################################################################

package b2

import (
	"testing"
)

func TestOrderLifecycleRules(t *testing.T) {
	cases := map[string]struct {
		rules   []LifecycleRule
		desired []interface{}
		want    []string // expected order of file name prefixes
	}{
		"empty": {
			rules:   nil,
			desired: nil,
			want:    []string{},
		},
		"single rule": {
			rules:   []LifecycleRule{{FileNamePrefix: "a/"}},
			desired: []interface{}{map[string]interface{}{"file_name_prefix": "a/"}},
			want:    []string{"a/"},
		},
		"already ordered": {
			rules: []LifecycleRule{
				{FileNamePrefix: "a/"},
				{FileNamePrefix: "b/"},
				{FileNamePrefix: "c/"},
			},
			desired: []interface{}{
				map[string]interface{}{"file_name_prefix": "a/"},
				map[string]interface{}{"file_name_prefix": "b/"},
				map[string]interface{}{"file_name_prefix": "c/"},
			},
			want: []string{"a/", "b/", "c/"},
		},
		"reordered by B2": {
			// B2 returns rules in an arbitrary order; the desired order must win
			rules: []LifecycleRule{
				{FileNamePrefix: "avatars/"},
				{FileNamePrefix: "thumbnails/"},
				{FileNamePrefix: "uploads/"},
				{FileNamePrefix: "files/"},
			},
			desired: []interface{}{
				map[string]interface{}{"file_name_prefix": "uploads/"},
				map[string]interface{}{"file_name_prefix": "files/"},
				map[string]interface{}{"file_name_prefix": "thumbnails/"},
				map[string]interface{}{"file_name_prefix": "avatars/"},
			},
			want: []string{"uploads/", "files/", "thumbnails/", "avatars/"},
		},
		"no desired rules": {
			// e.g. empty state (import); API order is kept
			rules: []LifecycleRule{
				{FileNamePrefix: "b/"},
				{FileNamePrefix: "a/"},
			},
			desired: nil,
			want:    []string{"b/", "a/"},
		},
		"rule added out-of-band": {
			// rules missing from desired keep API order at the end
			rules: []LifecycleRule{
				{FileNamePrefix: "a/"},
				{FileNamePrefix: "new/"},
				{FileNamePrefix: "b/"},
			},
			desired: []interface{}{
				map[string]interface{}{"file_name_prefix": "b/"},
				map[string]interface{}{"file_name_prefix": "a/"},
			},
			want: []string{"b/", "a/", "new/"},
		},
		"rule removed out-of-band": {
			// desired rules missing from B2 are skipped
			rules: []LifecycleRule{
				{FileNamePrefix: "a/"},
			},
			desired: []interface{}{
				map[string]interface{}{"file_name_prefix": "a/"},
				map[string]interface{}{"file_name_prefix": "removed/"},
			},
			want: []string{"a/"},
		},
		"duplicate prefixes": {
			rules: []LifecycleRule{
				{FileNamePrefix: "a/", DaysFromHidingToDeleting: 2},
				{FileNamePrefix: "a/", DaysFromHidingToDeleting: 1},
			},
			desired: []interface{}{
				map[string]interface{}{"file_name_prefix": "a/"},
				map[string]interface{}{"file_name_prefix": "a/"},
			},
			want: []string{"a/", "a/"},
		},
		"non-rule desired entry": {
			rules: []LifecycleRule{
				{FileNamePrefix: "a/"},
				{FileNamePrefix: "b/"},
			},
			desired: []interface{}{
				"not a rule",
				map[string]interface{}{"file_name_prefix": "b/"},
			},
			want: []string{"b/", "a/"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := orderLifecycleRules(tc.rules, tc.desired)

			prefixes := make([]string, len(got))
			for i, rule := range got {
				prefixes[i] = rule.FileNamePrefix
			}

			if len(prefixes) != len(tc.want) {
				t.Fatalf("expected prefixes %v, got %v", tc.want, prefixes)
			}
			for i := range tc.want {
				if prefixes[i] != tc.want[i] {
					t.Fatalf("expected prefixes %v, got %v", tc.want, prefixes)
				}
			}
		})
	}
}

func TestOrderLifecycleRulesKeepsRuleFields(t *testing.T) {
	rules := []LifecycleRule{
		{FileNamePrefix: "avatars/", DaysFromHidingToDeleting: 1},
		{FileNamePrefix: "uploads/", DaysFromUploadingToHiding: 1},
	}
	desired := []interface{}{
		map[string]interface{}{"file_name_prefix": "uploads/", "days_from_uploading_to_hiding": 1},
		map[string]interface{}{"file_name_prefix": "avatars/", "days_from_hiding_to_deleting": 1},
	}

	got := orderLifecycleRules(rules, desired)

	if got[0].FileNamePrefix != "uploads/" || got[0].DaysFromUploadingToHiding != 1 || got[0].DaysFromHidingToDeleting != 0 {
		t.Errorf("first rule does not match B2 data for 'uploads/': %+v", got[0])
	}
	if got[1].FileNamePrefix != "avatars/" || got[1].DaysFromHidingToDeleting != 1 || got[1].DaysFromUploadingToHiding != 0 {
		t.Errorf("second rule does not match B2 data for 'avatars/': %+v", got[1])
	}
}
