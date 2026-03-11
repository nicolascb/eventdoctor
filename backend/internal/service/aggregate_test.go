package service

import (
	"testing"

	"github.com/nicolascb/eventdoctor/internal/api/response"
)

// --- orderedMap tests ---

func TestOrderedMap_GetOrCreate_NewKey(t *testing.T) {
	m := newOrderedMap[string, int]()
	val := m.getOrCreate("a", func() int { return 42 })
	if *val != 42 {
		t.Fatalf("expected 42, got %d", *val)
	}
}

func TestOrderedMap_GetOrCreate_ExistingKey(t *testing.T) {
	m := newOrderedMap[string, int]()
	m.getOrCreate("a", func() int { return 1 })
	val := m.getOrCreate("a", func() int { return 99 })
	if *val != 1 {
		t.Fatalf("expected existing value 1, got %d", *val)
	}
}

func TestOrderedMap_Collect_PreservesOrder(t *testing.T) {
	m := newOrderedMap[string, int]()
	m.getOrCreate("c", func() int { return 3 })
	m.getOrCreate("a", func() int { return 1 })
	m.getOrCreate("b", func() int { return 2 })

	result := m.collect()
	if len(result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result))
	}
	if result[0] != 3 || result[1] != 1 || result[2] != 2 {
		t.Fatalf("unexpected order: %v", result)
	}
}

func TestOrderedMap_Collect_Empty(t *testing.T) {
	m := newOrderedMap[string, int]()
	result := m.collect()
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}
}

func TestOrderedMap_MutateViaPointer(t *testing.T) {
	m := newOrderedMap[string, int]()
	p := m.getOrCreate("x", func() int { return 10 })
	*p = 20
	result := m.collect()
	if result[0] != 20 {
		t.Fatalf("expected 20 after mutation, got %d", result[0])
	}
}

// --- findOrAppend tests ---

func TestFindOrAppend_AppendsWhenNotFound(t *testing.T) {
	s := []int{}
	p := findOrAppend(&s, func(v *int) bool { return *v == 5 }, func() int { return 5 })
	if *p != 5 {
		t.Fatalf("expected 5, got %d", *p)
	}
	if len(s) != 1 {
		t.Fatalf("expected slice length 1, got %d", len(s))
	}
}

func TestFindOrAppend_ReturnsExisting(t *testing.T) {
	s := []int{3, 7}
	p := findOrAppend(&s, func(v *int) bool { return *v == 7 }, func() int { return 99 })
	if *p != 7 {
		t.Fatalf("expected 7, got %d", *p)
	}
	if len(s) != 2 {
		t.Fatalf("slice should not grow, got length %d", len(s))
	}
}

func TestFindOrAppend_MutateExisting(t *testing.T) {
	s := []int{1, 2, 3}
	p := findOrAppend(&s, func(v *int) bool { return *v == 2 }, func() int { return 99 })
	*p = 42
	if s[1] != 42 {
		t.Fatalf("expected s[1]=42, got %d", s[1])
	}
}

// --- appendHeader tests ---

func TestAppendHeader_NilName(t *testing.T) {
	headers := []response.EventHeaderView{}
	appendHeader(&headers, nil, nil)
	if len(headers) != 0 {
		t.Fatal("expected no header appended when name is nil")
	}
}

func TestAppendHeader_WithName(t *testing.T) {
	headers := []response.EventHeaderView{}
	name := "X-Request-ID"
	desc := "unique request identifier"
	appendHeader(&headers, &name, &desc)
	if len(headers) != 1 {
		t.Fatalf("expected 1 header, got %d", len(headers))
	}
	if headers[0].Name != name || headers[0].Description != desc {
		t.Fatalf("unexpected header: %+v", headers[0])
	}
}

func TestAppendHeader_NilDescription(t *testing.T) {
	headers := []response.EventHeaderView{}
	name := "Authorization"
	appendHeader(&headers, &name, nil)
	if len(headers) != 1 {
		t.Fatalf("expected 1 header, got %d", len(headers))
	}
	if headers[0].Description != "" {
		t.Fatalf("expected empty description, got %q", headers[0].Description)
	}
}
