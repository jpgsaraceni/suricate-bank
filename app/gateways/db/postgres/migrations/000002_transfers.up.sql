BEGIN;

CREATE TABLE IF NOT EXISTS transfers
(
    id UUID PRIMARY KEY,
    account_origin_id UUID NOT NULL REFERENCES accounts (id),
    account_destination_id UUID NOT NULL REFERENCES accounts (id) CHECK (account_destination_id != account_origin_id),
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT(now())
);

COMMIT;