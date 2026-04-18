package parser

// DiffKeysOptions controls behavior of DiffKeys.
type DiffKeysOptions struct {
	// IgnoreCase treats keys as case-insensitive when comparing.
	IgnoreCase bool
}

// DefaultDiffKeysOptions returns sensible defaults.
func DefaultDiffKeysOptions() DiffKeysOptions {
	return DiffKeysOptions{
		IgnoreCase: false,
	}
}

// KeyDiff holds the result of comparing keys across two entry slices.
type KeyDiff struct {
	OnlyInBase   []string
	OnlyInOther  []string
	InBoth       []string
}

// DiffKeys compares the keys of base and other, returning which keys are
// exclusive to each side and which are shared.
func DiffKeys(base, other []EnvEntry, opts DiffKeysOptions) KeyDiff {
	normalize := func(k string) string {
		if opts.IgnoreCase {
			return strings.ToLower(k)
		}
		return k
	}

	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		baseMap[normalize(e.Key)] = e.Key
	}

	otherMap := make(map[string]string, len(other))
	for _, e := range other {
		otherMap[normalize(e.Key)] = e.Key
	}

	var result KeyDiff

	for norm, orig := range baseMap {
		if _, found := otherMap[norm]; found {
			result.InBoth = append(result.InBoth, orig)
		} else {
			result.OnlyInBase = append(result.OnlyInBase, orig)
		}
	}

	for norm, orig := range otherMap {
		if _, found := baseMap[norm]; !found {
			result.OnlyInOther = append(result.OnlyInOther, orig)
		}
	}

	sort.Strings(result.OnlyInBase)
	sort.Strings(result.OnlyInOther)
	sort.Strings(result.InBoth)

	return result
}
