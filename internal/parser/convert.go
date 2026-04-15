package parser

import (
	"fmt"
	"strings"
)

// OutputFormat represents the target format for env conversion.
type OutputFormat string

const (
	FormatDotEnv OutputFormat = "dotenv"
	FormatExport OutputFormat = "export"
	FormatJSON    OutputFormat = "json"
	FormatYAML    OutputFormat = "yaml"
)

// ConvertOptions controls how env entries are serialized.
type ConvertOptions struct {
	Format      OutputFormat
	MaskSecrets bool
	MaskOpts    MaskOptions
}

// DefaultConvertOptions returns sensible defaults.
func DefaultConvertOptions() ConvertOptions {
	return ConvertOptions{
		Format:      FormatDotEnv,
		MaskSecrets: false,
		MaskOpts:    DefaultMaskOptions(),
	}
}

// Convert serializes a slice of EnvEntry into the requested format string.
func Convert(entries []EnvEntry, opts ConvertOptions) (string, error) {
	switch opts.Format {
	case FormatDotEnv:
		return toDotEnv(entries, opts), nil
	case FormatExport:
		return toExport(entries, opts), nil
	case FormatJSON:
		return toJSON(entries, opts), nil
	case FormatYAML:
		return toYAML(entries, opts), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", opts.Format)
	}
}

func resolveValue(e EnvEntry, opts ConvertOptions) string {
	if opts.MaskSecrets && IsSensitive(e.Key, opts.MaskOpts) {
		return MaskValue(e.Value, opts.MaskOpts)
	}
	return e.Value
}

func toDotEnv(entries []EnvEntry, opts ConvertOptions) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, resolveValue(e, opts))
	}
	return sb.String()
}

func toExport(entries []EnvEntry, opts ConvertOptions) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "export %s=%q\n", e.Key, resolveValue(e, opts))
	}
	return sb.String()
}

func toJSON(entries []EnvEntry, opts ConvertOptions) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", e.Key, resolveValue(e, opts), comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}

func toYAML(entries []EnvEntry, opts ConvertOptions) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s: %q\n", e.Key, resolveValue(e, opts))
	}
	return sb.String()
}
