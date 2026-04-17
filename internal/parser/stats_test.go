package parser

import (
	"testing"
)

func makeStatsEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_SECRET", Value: ""},
		{Key: "PORT", Value: "8080"},
	}
}

func TestComputeStats_Total(t *testing.T) {
	entries := makeStatsEntries()
	s := ComputeStats(entries, DefaultStatsOptions())
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
}

func TestComputeStats_Empty(t *testing.T) {
	entries := makeStatsEntries()
	s := ComputeStats(entries, DefaultStatsOptions())
	if s.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", s.Empty)
	}
}

func TestComputeStats_Sensitive(t *testing.T) {
	entries := makeStatsEntries()
	s := ComputeStats(entries, DefaultStatsOptions())
	if s.Sensitive < 1 {
		t.Errorf("expected at least 1 sensitive key, got %d", s.Sensitive)
	}
}

func TestComputeStats_Prefixes(t *testing.T) {
	entries := makeStatsEntries()
	s := ComputeStats(entries, DefaultStatsOptions())
	if s.Prefixes["DB"] != 2 {
		t.Errorf("expected DB prefix count=2, got %d", s.Prefixes["DB"])
	}
	if s.Prefixes["APP"] != 2.Errorf("expected APP prefix count=2, got %d", s.Prefixes["APP"])
	}
}

func TestComputeStats_NoPrefixKey(t *testing.T) {
	entries := makeStatsEntries()
	s := ComputeStats(entries, DefaultStatsOptions())
	if _, ok := s.Prefixes["PORT"]; ok {
		t.Error("PORT has no underscore, should not appear in prefixes")
	}
}
