package storage

import (
	"context"
	"database/sql"
	"time"
)

type DeviceRepo struct{ db *sql.DB }

func NewDeviceRepo(db *sql.DB) *DeviceRepo { return &DeviceRepo{db: db} }

func (r *DeviceRepo) Create(ctx context.Context, d Device) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO devices(id, name, type, mqtt_device_id, capabilities_json, created_at)
		VALUES(?, ?, ?, ?, ?, ?)
	`, d.ID, d.Name, d.Type, d.MQTTDeviceID, d.Capabilities, d.CreatedAt.UTC().Format(time.RFC3339Nano))
	return err
}

func (r *DeviceRepo) Get(ctx context.Context, id string) (Device, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, type, mqtt_device_id, capabilities_json, created_at
		FROM devices WHERE id = ?
	`, id)

	var d Device
	var created string
	if err := row.Scan(&d.ID, &d.Name, &d.Type, &d.MQTTDeviceID, &d.Capabilities, &created); err != nil {
		return Device{}, err
	}
	t, _ := time.Parse(time.RFC3339Nano, created)
	d.CreatedAt = t
	return d, nil
}

func (r *DeviceRepo) GetByMQTTDeviceID(ctx context.Context, mqttID string) (Device, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, type, mqtt_device_id, capabilities_json, created_at
		FROM devices WHERE mqtt_device_id = ?
	`, mqttID)

	var d Device
	var created string
	if err := row.Scan(&d.ID, &d.Name, &d.Type, &d.MQTTDeviceID, &d.Capabilities, &created); err != nil {
		return Device{}, err
	}
	t, _ := time.Parse(time.RFC3339Nano, created)
	d.CreatedAt = t
	return d, nil
}

func (r *DeviceRepo) List(ctx context.Context) ([]Device, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, type, mqtt_device_id, capabilities_json, created_at
		FROM devices ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Device
	for rows.Next() {
		var d Device
		var created string
		if err := rows.Scan(&d.ID, &d.Name, &d.Type, &d.MQTTDeviceID, &d.Capabilities, &created); err != nil {
			return nil, err
		}
		t, _ := time.Parse(time.RFC3339Nano, created)
		d.CreatedAt = t
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *DeviceRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM devices WHERE id = ?`, id)
	return err
}
