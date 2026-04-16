// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package db

const schemaSQL = `
PRAGMA journal_mode=WAL;
PRAGMA foreign_keys=ON;


CREATE TABLE IF NOT EXISTS settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS connectors (
  name TEXT PRIMARY KEY,
  enabled INTEGER NOT NULL DEFAULT 0,
  auth_status TEXT NOT NULL DEFAULT 'not_configured',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS memories (
  id TEXT PRIMARY KEY,
  scope TEXT,
  content TEXT NOT NULL,
  tags TEXT,
  created_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS archive_index (
  id TEXT PRIMARY KEY,
  object_ref TEXT NOT NULL,
  location TEXT NOT NULL,
  checksum TEXT,
  state TEXT NOT NULL,
  created_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  surface TEXT NOT NULL,
  user_id TEXT NOT NULL,
  workspace TEXT,
  started_at DATETIME NOT NULL,
  ended_at DATETIME,
  summary TEXT
);
CREATE TABLE IF NOT EXISTS messages (
  id TEXT PRIMARY KEY,
  session_id TEXT NOT NULL,
  role TEXT NOT NULL,
  content TEXT NOT NULL,
  tokens_est INTEGER DEFAULT 0,
  created_at DATETIME NOT NULL,
  FOREIGN KEY(session_id) REFERENCES sessions(id)
);
CREATE TABLE IF NOT EXISTS memory_entries (
  id TEXT PRIMARY KEY,
  scope_user TEXT,
  scope_project TEXT,
  kind TEXT NOT NULL,
  content TEXT NOT NULL,
  tags TEXT,
  importance INTEGER DEFAULT 1,
  pinned INTEGER DEFAULT 0,
  last_access_at DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS memory_pins (
  id TEXT PRIMARY KEY,
  memory_id TEXT NOT NULL,
  reason TEXT,
  pinned_by TEXT,
  created_at DATETIME NOT NULL,
  FOREIGN KEY(memory_id) REFERENCES memory_entries(id)
);
CREATE TABLE IF NOT EXISTS tasks (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  schedule_expr TEXT,
  payload_json TEXT,
  enabled INTEGER NOT NULL DEFAULT 1,
  last_run_at DATETIME,
  next_run_at DATETIME,
  created_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS connector_auth (
  id TEXT PRIMARY KEY,
  connector TEXT NOT NULL,
  account_label TEXT,
  scopes TEXT,
  token_ref TEXT,
  status TEXT NOT NULL,
  expires_at DATETIME,
  updated_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS archive_objects (
  id TEXT PRIMARY KEY,
  class TEXT NOT NULL,
  local_path TEXT,
  drive_file_id TEXT,
  checksum_sha256 TEXT NOT NULL,
  size_bytes INTEGER NOT NULL,
  compressed INTEGER DEFAULT 0,
  state TEXT NOT NULL,
  restored_until DATETIME,
  last_access_at DATETIME,
  created_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS file_metadata (
  id TEXT PRIMARY KEY,
  source TEXT NOT NULL,
  source_message_id TEXT,
  mime_type TEXT,
  original_name TEXT,
  local_path TEXT,
  size_bytes INTEGER,
  checksum_sha256 TEXT,
  archive_object_id TEXT,
  created_at DATETIME NOT NULL,
  FOREIGN KEY(archive_object_id) REFERENCES archive_objects(id)
);
CREATE TABLE IF NOT EXISTS command_approvals (
  id TEXT PRIMARY KEY,
  command_text TEXT NOT NULL,
  risk_tier INTEGER NOT NULL,
  requested_by TEXT NOT NULL,
  scope TEXT,
  status TEXT NOT NULL,
  expires_at DATETIME,
  decided_by TEXT,
  decided_at DATETIME,
  created_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS audit_log (
  id TEXT PRIMARY KEY,
  actor TEXT NOT NULL,
  action TEXT NOT NULL,
  target TEXT,
  connector TEXT,
  outcome TEXT NOT NULL,
  details_json TEXT,
  created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_messages_session_created ON messages(session_id, created_at);
CREATE INDEX IF NOT EXISTS idx_memory_scope_kind ON memory_entries(scope_user, scope_project, kind);
CREATE INDEX IF NOT EXISTS idx_tasks_next_run ON tasks(next_run_at) WHERE enabled = 1;
CREATE INDEX IF NOT EXISTS idx_archive_state ON archive_objects(state, last_access_at);
CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_log(created_at);
`
