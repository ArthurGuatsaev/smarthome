package storage

import (
	"context"
	"database/sql"
	"time"
)

type CommandRepo struct{ db *sql.DB }

func NewCommandRepo(db *sql.DB) *CommandRepo { return &CommandRepo{db: db} }

func (r *CommandRepo) Create(ctx context.Context, c Command) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO commands(id, device_id, action, params_json, status, error, created_at, acked_at)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`, c.ID, c.DeviceID, c.Action, c.ParamsJSON, c.Status, c.Error,
		c.CreatedAt.UTC().Format(time.RFC3339Nano),
		nil,
	)
	return err
}

func (r *CommandRepo) Get(ctx context.Context, id string) (Command, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, device_id, action, params_json, status, error, created_at, acked_at
		FROM commands WHERE id = ?
	`, id)

	var c Command
	var created string
	var acked sql.NullString

	if err := row.Scan(&c.ID, &c.DeviceID, &c.Action, &c.ParamsJSON, &c.Status, &c.Error, &created, &acked); err != nil {
		return Command{}, err
	}

	ct, _ := time.Parse(time.RFC3339Nano, created)
	c.CreatedAt = ct

	if acked.Valid {
		at, _ := time.Parse(time.RFC3339Nano, acked.String)
		c.AckedAt = &at
	}

	return c, nil
}

func (r *CommandRepo) SetAck(ctx context.Context, id string, ok bool, errMsg string, at time.Time) error {
	status := "acked"
	if !ok {
		status = "failed"
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE commands
		SET status = ?, error = ?, acked_at = ?
		WHERE id = ?
	`, status, errMsg, at.UTC().Format(time.RFC3339Nano), id)
	return err
}

func (r *CommandRepo) SetTimeout(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE commands
		SET status = 'timeout'
		WHERE id = ? AND status = 'pending'
	`, id)
	return err
}
