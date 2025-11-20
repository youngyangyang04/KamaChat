-- End-to-End Encryption Support Migration
-- Version: v2.0
-- Created: 2025-11-19

-- 1. Add encryption fields to user_info table
ALTER TABLE user_info ADD COLUMN IF NOT EXISTS identity_key_public TEXT;
ALTER TABLE user_info ADD COLUMN IF NOT EXISTS signed_pre_key_id INT DEFAULT 1;
ALTER TABLE user_info ADD COLUMN IF NOT EXISTS signed_pre_key_public TEXT;
ALTER TABLE user_info ADD COLUMN IF NOT EXISTS signed_pre_key_signature TEXT;
ALTER TABLE user_info ADD COLUMN IF NOT EXISTS signed_pre_key_timestamp TIMESTAMP;
ALTER TABLE user_info ADD COLUMN IF NOT EXISTS key_generation INT DEFAULT 1;

-- 2. Create one_time_pre_keys table
CREATE TABLE IF NOT EXISTS one_time_pre_keys (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    key_id INT NOT NULL,
    public_key TEXT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP,
    used_by_user_id VARCHAR(20),
    UNIQUE(user_id, key_id)
);

CREATE INDEX IF NOT EXISTS idx_otpk_user_unused ON one_time_pre_keys(user_id, used) WHERE used = FALSE;
CREATE INDEX IF NOT EXISTS idx_otpk_user_id ON one_time_pre_keys(user_id);

-- 3. Add encryption fields to message table
ALTER TABLE message ADD COLUMN IF NOT EXISTS is_encrypted BOOLEAN DEFAULT FALSE;
ALTER TABLE message ADD COLUMN IF NOT EXISTS encryption_version INT DEFAULT 1;
ALTER TABLE message ADD COLUMN IF NOT EXISTS message_type VARCHAR(20) DEFAULT 'SignalMessage';
ALTER TABLE message ADD COLUMN IF NOT EXISTS sender_identity_key TEXT;
ALTER TABLE message ADD COLUMN IF NOT EXISTS sender_ephemeral_key TEXT;
ALTER TABLE message ADD COLUMN IF NOT EXISTS used_one_time_pre_key_id INT;
ALTER TABLE message ADD COLUMN IF NOT EXISTS ratchet_key TEXT;
ALTER TABLE message ADD COLUMN IF NOT EXISTS counter INT;
ALTER TABLE message ADD COLUMN IF NOT EXISTS prev_counter INT;
ALTER TABLE message ADD COLUMN IF NOT EXISTS iv TEXT;
ALTER TABLE message ADD COLUMN IF NOT EXISTS auth_tag TEXT;

CREATE INDEX IF NOT EXISTS idx_message_encrypted ON message(is_encrypted);

-- 4. Create key_replenishment_log table
CREATE TABLE IF NOT EXISTS key_replenishment_log (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    replenished_at TIMESTAMP NOT NULL DEFAULT NOW(),
    keys_added INT NOT NULL,
    remaining_keys INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_krl_user_time ON key_replenishment_log(user_id, replenished_at);

-- Migration completed

