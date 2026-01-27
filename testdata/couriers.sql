BEGIN;

-- ускорение на время наполнения
SET LOCAL synchronous_commit = OFF;
SET LOCAL work_mem = '256MB';

-- нужно один раз (можно оставить, если уже включено)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- очистка (если надо)
TRUNCATE TABLE delivery, couriers RESTART IDENTITY CASCADE;

-- 1) couriers: 100 000
INSERT INTO couriers (name, phone, status, created_at, updated_at, transport_type)
SELECT
  'Courier ' || gs::text AS name,
  -- уникальный телефон: +79XXXXXXXXX
  '+79' || lpad(gs::text, 9, '0') AS phone,
  (ARRAY['available','busy','inactive'])[1 + (random()*2)::int] AS status,
  now() - (random() * interval '365 days') AS created_at,
  now() - (random() * interval '30 days')  AS updated_at,
  (ARRAY['on_foot','bike','car','scooter'])[1 + (random()*3)::int] AS transport_type
FROM generate_series(1, 100000) gs;

-- 2) delivery: 1 000 000
INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
SELECT
  1 + (random() * 99999)::int AS courier_id,
  gen_random_uuid()::text     AS order_id,    -- UUID строкой
  ts                          AS assigned_at,
  ts + (15 + (random() * 180)) * interval '1 minute' AS deadline
FROM (
  SELECT
    gs,
    now() - (random() * interval '90 days') AS ts
  FROM generate_series(1, 1000000) gs
) s;

COMMIT;

ANALYZE couriers;
ANALYZE delivery;