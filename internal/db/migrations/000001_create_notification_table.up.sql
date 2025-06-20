CREATE TYPE notification_type AS ENUM ('email', 'push', 'sms');
CREATE TYPE notification_status AS ENUM ('PENDING', 'IN_PROGRESS', 'RETRYING', 'SENT', 'FAILED');

CREATE TABLE IF NOT EXISTS notification (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    type notification_type NOT NULL,
    message TEXT NOT NULL,
    priority INT NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    status notification_status NOT NULL,
    attempts INT NOT NULL,
    last_error TEXT,
    delivered_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS notification_user_id_idx ON notification (id);