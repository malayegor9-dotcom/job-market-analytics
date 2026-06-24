CREATE TABLE IF NOT EXISTS sources (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    url        TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO sources (name, url) VALUES
    ('hh',       'https://api.hh.ru'),
    ('remoteok', 'https://remoteok.com/api')
ON CONFLICT (name) DO NOTHING;

