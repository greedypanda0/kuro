package db

import (
	"database/sql"
	"testing"
)

func TestDefaultsCreated(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := ApplySchema(db); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	head, err := GetConfig(db, "head")
	if err != nil {
		t.Fatalf("get head: %v", err)
	}
	if head != "main" {
		t.Fatalf("expected head to be main, got %s", head)
	}

	ref, err := GetRef(db, "main")
	if err != nil {
		t.Fatalf("get main ref: %v", err)
	}
	if ref.Name != "main" {
		t.Fatalf("expected main ref, got %s", ref.Name)
	}
	if ref.SnapshotHash != nil {
		t.Fatalf("expected main ref snapshot to be nil")
	}
}

func TestStagingLifecycle(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := ApplySchema(db); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	if err := AddStageFile(db, "file-a.txt"); err != nil {
		t.Fatalf("add stage file: %v", err)
	}
	if err := AddStageFile(db, "dir/file-b.txt"); err != nil {
		t.Fatalf("add stage file: %v", err)
	}

	files, err := GetStageFiles(db)
	if err != nil {
		t.Fatalf("get stage files: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 staged files, got %d", len(files))
	}

	if err := RemoveStageFile(db, "file-a.txt"); err != nil {
		t.Fatalf("remove stage file: %v", err)
	}

	files, err = GetStageFiles(db)
	if err != nil {
		t.Fatalf("get stage files: %v", err)
	}
	if len(files) != 1 || files[0].Path != "dir/file-b.txt" {
		t.Fatalf("unexpected staged files after remove")
	}

	if err := ClearStage(db); err != nil {
		t.Fatalf("clear stage: %v", err)
	}

	files, err = GetStageFiles(db)
	if err != nil {
		t.Fatalf("get stage files: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected empty stage, got %d", len(files))
	}
}
