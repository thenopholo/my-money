-- 006_create_planned_incomes.sql

CREATE TABLE IF NOT EXISTS planned_incomes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    description VARCHAR(255) NOT NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    due_day INTEGER NOT NULL CHECK (due_day BETWEEN 1 AND 31),
    frequency VARCHAR(20) NOT NULL CHECK (frequency IN ('once', 'monthly', 'yearly')),
    start_date DATE,
    end_date DATE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraint: end_date deve ser >= start_date (se ambos definidos)
    CONSTRAINT chk_planned_incomes_dates CHECK (
        end_date IS NULL OR start_date IS NULL OR end_date >= start_date
    )
);

-- Índice para buscar por usuário
CREATE INDEX idx_planned_incomes_user_id ON planned_incomes(user_id);

-- Índice para buscar ativos
CREATE INDEX idx_planned_incomes_active ON planned_incomes(user_id, is_active);

-- Índice para buscar por conta
CREATE INDEX idx_planned_incomes_account ON planned_incomes(account_id);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_planned_incomes_account;
DROP INDEX IF EXISTS idx_planned_incomes_active;
DROP INDEX IF EXISTS idx_planned_incomes_user_id;
DROP TABLE IF EXISTS planned_incomes;
