package parser

import "fmt"

// PatchOperation represents a single key-level change to apply.
type PatchOperation struct {
	Action string // "set", "delete", "rename"
	Key    string
	Value  string
	NewKey string // used for rename
}

// PatchOptions controls Patch behaviour.
type PatchOptions struct {
	IgnoreMissing bool // skip delete/rename if key not found
}

func DefaultPatchOptions() PatchOptions {
	return PatchOptions{IgnoreMissing: false}
}

// Patch applies a list of PatchOperations to entries and returns the result.
func Patch(entries []EnvEntry, ops []PatchOperation, opts PatchOptions) ([]EnvEntry, error) {
	result := make([]EnvEntry, len(entries))
	copy(result, entries)

	for _, op := range ops {
		switch op.Action {
		case "set":
			result = patchSet(result, op.Key, op.Value)
		case "delete":
			var err error
			result, err = patchDelete(result, op.Key, opts.IgnoreMissing)
			if err != nil {
				return nil, err
			}
		case "rename":
			var err error
			result, err = patchRename(result, op.Key, op.NewKey, opts.IgnoreMissing)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown patch action: %s", op.Action)
		}
	}
	return result, nil
}

func patchSet(entries []EnvEntry, key, value string) []EnvEntry {
	for i, e := range entries {
		if e.Key == key {
			entries[i].Value = value
			return entries
		}
	}
	return append(entries, EnvEntry{Key: key, Value: value})
}

func patchDelete(entries []EnvEntry, key string, ignoreMissing bool) ([]EnvEntry, error) {
	for i, e := range entries {
		if e.Key == key {
			return append(entries[:i], entries[i+1:]...), nil
		}
	}
	if !ignoreMissing {
		return nil, fmt.Errorf("patch delete: key %q not found", key)
	}
	return entries, nil
}

func patchRename(entries []EnvEntry, key, newKey string, ignoreMissing bool) ([]EnvEntry, error) {
	for i, e := range entries {
		if e.Key == key {
			entries[i].Key = newKey
			return entries, nil
		}
	}
	if !ignoreMissing {
		return nil, fmt.Errorf("patch rename: key %q not found", key)
	}
	return entries, nil
}
