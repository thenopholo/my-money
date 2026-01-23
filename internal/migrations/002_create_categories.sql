-- 002_create_categories.sql

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    category_type VARCHAR(20) NOT NULL CHECK (category_type IN ('income', 'expense')),
    color VARCHAR(20),
    icon VARCHAR(50)
);

-- Índice para buscar categorias do usuário
CREATE INDEX idx_categories_user_id ON categories(user_id);

-- Índice para buscar por tipo
CREATE INDEX idx_categories_type ON categories(user_id, category_type);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_categories_type;
DROP INDEX IF EXISTS idx_categories_user_id;
DROP TABLE IF EXISTS categories;
