-- 003_create_bank_accounts.sql

CREATE TABLE IF NOT EXISTS bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    bank_name VARCHAR(100) NOT NULL,
    account_type VARCHAR(20) NOT NULL CHECK (account_type IN ('checking', 'savings')),
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índice para buscar contas do usuário
CREATE INDEX idx_bank_accounts_user_id ON bank_accounts(user_id);

-- Índice para buscar contas ativas
CREATE INDEX idx_bank_accounts_active ON bank_accounts(user_id, is_active);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_bank_accounts_active;
DROP INDEX IF EXISTS idx_bank_accounts_user_id;
DROP TABLE IF EXISTS bank_accounts;
