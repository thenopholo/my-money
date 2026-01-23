-- 005_create_invoices.sql

CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES credit_cards(id) ON DELETE CASCADE,
    reference_date DATE NOT NULL,
    due_date DATE NOT NULL,
    total_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed', 'paid', 'expired')),
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índice para buscar faturas do cartão
CREATE INDEX idx_invoices_card_id ON invoices(card_id);

-- Índice para buscar por status
CREATE INDEX idx_invoices_status ON invoices(card_id, status);

-- Índice para buscar por mês de referência
CREATE INDEX idx_invoices_reference ON invoices(card_id, reference_date);

-- Constraint: não pode ter duas faturas para o mesmo mês/cartão
CREATE UNIQUE INDEX idx_invoices_unique_month ON invoices(card_id, reference_date);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_invoices_unique_month;
DROP INDEX IF EXISTS idx_invoices_reference;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_card_id;
DROP TABLE IF EXISTS invoices;
