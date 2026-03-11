package service

import (
	"context"
	"testing"
)

// --- GetServiceView integration tests ---

func TestGetServiceView(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.GetServiceView(context.Background(), "service-a")
	if err != nil {
		t.Fatalf("GetServiceView() error: %v", err)
	}

	if result.Service != "service-a" {
		t.Errorf("expected service-a, got %s", result.Service)
	}
	if len(result.Produces) != 2 {
		t.Errorf("expected 2 produces entries, got %d", len(result.Produces))
	}
	if len(result.Consumes) != 1 {
		t.Errorf("expected 1 consumes entry, got %d", len(result.Consumes))
	}
}

func TestGetServiceView_ServiceB(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.GetServiceView(context.Background(), "service-b")
	if err != nil {
		t.Fatalf("GetServiceView() error: %v", err)
	}

	if result.Service != "service-b" {
		t.Errorf("expected service-b, got %s", result.Service)
	}
	if len(result.Produces) != 1 {
		t.Errorf("expected 1 produces entry, got %d", len(result.Produces))
	}
	if len(result.Consumes) != 1 {
		t.Errorf("expected 1 consumes entry, got %d", len(result.Consumes))
	}
}

func TestGetServiceView_NonExistent(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	result, err := svc.GetServiceView(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("GetServiceView() error: %v", err)
	}

	if len(result.Produces) != 0 {
		t.Errorf("expected 0 produces for non-existent service, got %d", len(result.Produces))
	}
	if len(result.Consumes) != 0 {
		t.Errorf("expected 0 consumes for non-existent service, got %d", len(result.Consumes))
	}
}
