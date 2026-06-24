CREATE TABLE IF NOT EXISTS skills (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    category   TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(name)
);

INSERT INTO skills (name, category) VALUES
    ('go', 'language'), ('python', 'language'), ('java', 'language'),
    ('typescript', 'language'), ('rust', 'language'),
    ('gin', 'framework'), ('echo', 'framework'), ('fastapi', 'framework'),
    ('django', 'framework'), ('react', 'framework'),
    ('postgresql', 'database'), ('mysql', 'database'),
    ('mongodb', 'database'), ('redis', 'database'), ('clickhouse', 'database'),
    ('docker', 'tool'), ('kubernetes', 'tool'),
    ('kafka', 'tool'), ('nats', 'tool'), ('nginx', 'tool'),
    ('aws', 'cloud'), ('gcp', 'cloud'), ('azure', 'cloud')
ON CONFLICT (name) DO NOTHING;