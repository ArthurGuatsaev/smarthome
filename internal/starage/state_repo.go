package storage

import (
	"context"
	"database/sql"
	"time"
)

type StateRepo struct{ db *sql.DB }

func NewStateRepo(db *sql.DB) *StateRepo { return &StateRepo{db: db} }

func (r *StateRepo) Upsert(ctx context.Context, s DeviceState) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO device_state(device_id, state_json, updated_at)
		VALUES(?, ?, ?)
		ON CONFLICT(device_id) DO UPDATE SET
		  state_json = excluded.state_json,
		  updated_at = excluded.updated_at
	`, s.DeviceID, s.StateJSON, s.UpdatedAt.UTC().Format(time.RFC3339Nano))
	return err
}

func (r *StateRepo) Get(ctx context.Context, deviceID string) (DeviceState, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT device_id, state_json, updated_at
		FROM device_state
		WHERE device_id = ?
	`, deviceID)

	var s DeviceState
	var updated string
	if err := row.Scan(&s.DeviceID, &s.StateJSON, &updated); err != nil {
		return DeviceState{}, err
	}
	t, _ := time.Parse(time.RFC3339Nano, updated)
	s.UpdatedAt = t
	return s, nil
}
