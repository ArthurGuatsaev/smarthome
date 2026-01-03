PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS devices (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  mqtt_device_id TEXT NOT NULL UNIQUE,
  capabilities_json TEXT NOT NULL,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS device_state (
  device_id TEXT PRIMARY KEY,
  state_json TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY(device_id) REFERENCES devices(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS commands (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL,
  action TEXT NOT NULL,
  params_json TEXT NOT NULL,
  status TEXT NOT NULL,         -- pending|acked|failed|timeout
  error TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  acked_at TEXT,
  FOREIGN KEY(device_id) REFERENCES devices(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_commands_device_created
ON commands(device_id, created_at DESC);
