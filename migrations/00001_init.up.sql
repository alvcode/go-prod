BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE public.currency
(
    id SERIAL PRIMARY KEY,
    name TEXT,
    symbol TEXT
);

CREATE TABLE public.category
(
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE public.product
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    image_id UUID,
    price BIGINT,
    currency_id INT,
    rating INT,
    category_id INT NOT NULL,
    specification JSONB,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

INSERT INTO public.currency (name, symbol)
VALUES ('рубль', 'Р');

INSERT INTO public.currency (name, symbol)
VALUES ('dollar', '$');

COMMIT;