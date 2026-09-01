// Package migration is a stub representing a state migration package. All Optional+Computed
// schema fields here are intentionally uncommented — the analyzer must not report them.
package migration

import (
	"migration/pluginsdk"
	"migration/schema"
)

// ResourceV0ToV1Schema mirrors the real pattern: a historical schema snapshot returned by
// a StateUpgrader's Schema() method. O+C fields appear without a reason comment because
// migration schemas are stripped of all non-essential metadata.
func ResourceV0ToV1Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		// O+C without comment — must NOT be flagged in a migration package
		"disk_size_gb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		// O+C without comment, C+O field order — must NOT be flagged
		"disk_iops_read_write": {
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},

		// O+C without comment, pluginsdk alias — must NOT be flagged
		"tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
