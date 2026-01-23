-- 004_create_credit_cards.sql

CREATE TABLE IF NOT EXISTS credit_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    credit_limit NUMERIC(15, 2) NOT NULL,
    close_day INTEGER NOT NULL CHECK (close_day BETWEEN 1 AND 28),
    due_day INTEGER NOT NULL CHECK (due_day BETWEEN 1 AND 31),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índice para buscar cartões do usuário
CREATE INDEX idx_credit_cards_user_id ON credit_cards(user_id);

-- Índice para buscar cartões ativos
CREATE INDEX idx_credit_cards_active ON credit_cards(user_id, is_active);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_credit_cards_active;
DROP INDEX IF EXISTS idx_credit_cards_user_id;
DROP TABLE IF EXISTS credit_cards;
