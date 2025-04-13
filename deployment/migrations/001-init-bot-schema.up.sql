CREATE SCHEMA IF NOT EXISTS digikeeper;

CREATE TABLE IF NOT EXISTS digikeeper.user (
    user_uid UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE SCHEMA IF NOT EXISTS digikeeper_tg;

CREATE TABLE IF NOT EXISTS digikeeper_tg.tg_user_sessions (
    user_uid UUID NOT NULL references digikeeper.user (user_uid),
    tg_user_id BIGINT PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    data JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_updated_at ON digikeeper_tg.user_sessions (updated_at);

CREATE TABLE IF NOT EXISTS digikeeper.message_logs (
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    message_text TEXT,
    message_type TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_message_logs_user_id ON digikeeper.message_logs (user_id);

CREATE INDEX IF NOT EXISTS idx_message_logs_chat_id ON digikeeper.message_logs (chat_id);

SELECT
    create_hypertable (
        'digikeeper.message_logs',
        by_range ('timestamp', INTERVAL '1 day'),
        if_not_exists = > TRUE
    );

SELECT
    add_retention_policy (
        'digikeeper.message_logs',
        INTERVAL '30 days',
        if_not_exists = > TRUE
    );
