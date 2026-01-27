package db

import (
	"core/errors"
	"database/sql"
)

type Ref struct {
	Name         string
	SnapshotHash *string
	UpdatedAt    int64
}

func ListRefs(db DBTX) ([]Ref, error) {
	var refs []Ref

	rows, err := db.Query("SELECT name, snapshot_hash, updated_at FROM refs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			ref          Ref
			snapshotHash sql.NullString
		)

		if err := rows.Scan(&ref.Name, &snapshotHash, &ref.UpdatedAt); err != nil {
			return nil, err
		}

		if snapshotHash.Valid {
			ref.SnapshotHash = &snapshotHash.String
		}

		refs = append(refs, ref)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return refs, nil
}

func SetRef(db DBTX, name string, snapshotHash *string) error {
	_, err := db.Exec(
		"INSERT OR IGNORE INTO refs (name, snapshot_hash) VALUES (?, ?)",
		name,
		snapshotHash,
	)
	return err
}

func UpdateRef(db DBTX, name string, snapshotHash *string) error {
	res, err := db.Exec(
		"UPDATE refs SET snapshot_hash = ?, updated_at = (strftime('%s', 'now')) WHERE name = ?",
		snapshotHash,
		name,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.ErrRefNotFound
	}

	return nil
}

func GetRef(db DBTX, name string) (*Ref, error) {
	var (
		ref          Ref
		snapshotHash sql.NullString
	)

	err := db.QueryRow(
		"SELECT name, snapshot_hash, updated_at FROM refs WHERE name = ?",
		name,
	).Scan(&ref.Name, &snapshotHash, &ref.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.ErrRefNotFound
	}
	if err != nil {
		return nil, err
	}

	if snapshotHash.Valid {
		ref.SnapshotHash = &snapshotHash.String
	}

	return &ref, nil
}

func DeleteRef(db DBTX, name string) error {
	res, err := db.Exec(
		"DELETE FROM refs WHERE name = ?",
		name,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.ErrRefNotFound
	}

	return nil
}
