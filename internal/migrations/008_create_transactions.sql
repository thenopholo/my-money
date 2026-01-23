-- 008_create_transactions.sql

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    invoice_id UUID REFERENCES invoices(id) ON DELETE SET NULL,
    planned_income_id UUID REFERENCES planned_incomes(id) ON DELETE SET NULL,
    planned_expense_id UUID REFERENCES planned_expenses(id) ON DELETE SET NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('income', 'expense')),
    description VARCHAR(255) NOT NULL,
    transaction_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índice para buscar por conta
CREATE INDEX idx_transactions_account_id ON transactions(account_id);

-- Índice para buscar por data
CREATE INDEX idx_transactions_date ON transactions(account_id, transaction_date);

-- Índice para buscar por categoria
CREATE INDEX idx_transactions_category ON transactions(category_id);

-- Índice para buscar por tipo
CREATE INDEX idx_transactions_type ON transactions(account_id, transaction_type);

-- Índice para buscar transações de planejamento
CREATE INDEX idx_transactions_planned_income ON transactions(planned_income_id) WHERE planned_income_id IS NOT NULL;
CREATE INDEX idx_transactions_planned_expense ON transactions(planned_expense_id) WHERE planned_expense_id IS NOT NULL;

---- create above / drop below ----

DROP INDEX IF EXISTS idx_transactions_planned_expense;
DROP INDEX IF EXISTS idx_transactions_planned_income;
DROP INDEX IF EXISTS idx_transactions_type;
DROP INDEX IF EXISTS idx_transactions_category;
DROP INDEX IF EXISTS idx_transactions_date;
DROP INDEX IF EXISTS idx_transactions_account_id;
DROP TABLE IF EXISTS transactions;
