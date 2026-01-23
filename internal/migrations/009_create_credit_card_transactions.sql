-- 009_create_credit_card_transactions.sql

CREATE TABLE IF NOT EXISTS credit_card_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES credit_cards(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    invoice_id UUID REFERENCES invoices(id) ON DELETE SET NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    description VARCHAR(255) NOT NULL,
    installments INTEGER NOT NULL DEFAULT 1 CHECK (installments BETWEEN 1 AND 48),
    current_installment INTEGER NOT NULL DEFAULT 1 CHECK (current_installment >= 1),
    installment_value NUMERIC(15, 2),
    transaction_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraint: current_installment não pode ser maior que installments
    CONSTRAINT chk_installment_range CHECK (current_installment <= installments)
);

-- Índice para buscar por cartão
CREATE INDEX idx_cc_transactions_card_id ON credit_card_transactions(card_id);

-- Índice para buscar por fatura
CREATE INDEX idx_cc_transactions_invoice ON credit_card_transactions(invoice_id);

-- Índice para buscar por data
CREATE INDEX idx_cc_transactions_date ON credit_card_transactions(card_id, transaction_date);

-- Índice para buscar por categoria
CREATE INDEX idx_cc_transactions_category ON credit_card_transactions(category_id);

-- Índice para buscar parcelas pendentes (sem fatura atribuída)
CREATE INDEX idx_cc_transactions_pending ON credit_card_transactions(card_id) WHERE invoice_id IS NULL;

---- create above / drop below ----

DROP INDEX IF EXISTS idx_cc_transactions_pending;
DROP INDEX IF EXISTS idx_cc_transactions_category;
DROP INDEX IF EXISTS idx_cc_transactions_date;
DROP INDEX IF EXISTS idx_cc_transactions_invoice;
DROP INDEX IF EXISTS idx_cc_transactions_card_id;
DROP TABLE IF EXISTS credit_card_transactions;
