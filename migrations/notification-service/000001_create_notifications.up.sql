CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    recipient_email TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    idempotency_key TEXT NOT NULL UNIQUE,
    error_message TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    sent_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS notification_settings (
    user_id TEXT PRIMARY KEY,
    email_enabled BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS notification_templates (
    name TEXT PRIMARY KEY,
    subject TEXT NOT NULL,
    body TEXT NOT NULL
);
