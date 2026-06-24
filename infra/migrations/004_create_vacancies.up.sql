CREATE TABLE IF NOT EXISTS vacancies (
    id            SERIAL PRIMARY KEY,
    external_id   TEXT NOT NULL,
    source_id     INT  REFERENCES sources(id),
    company_id    INT  REFERENCES companies(id),
    title         TEXT NOT NULL,
    description   TEXT,
    location      TEXT,
    salary_min    BIGINT,
    salary_max    BIGINT,
    currency      TEXT DEFAULT 'RUB',
    remote        BOOLEAN DEFAULT false,
    url           TEXT,
    published_at  TIMESTAMPTZ,
    collected_at  TIMESTAMPTZ DEFAULT NOW(),
    search_vector TSVECTOR,
    UNIQUE(source_id, external_id)
);

CREATE INDEX idx_vac_search   ON vacancies USING GIN(search_vector);
CREATE INDEX idx_vac_location ON vacancies(location);
CREATE INDEX idx_vac_salary   ON vacancies(salary_min, salary_max);
CREATE INDEX idx_vac_date     ON vacancies(published_at DESC);
CREATE INDEX idx_vac_remote   ON vacancies(remote);

CREATE OR REPLACE FUNCTION update_search_vector()
RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector(
        'russian',
        coalesce(NEW.title, '') || ' ' || coalesce(NEW.description, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trig_vacancies_search
BEFORE INSERT OR UPDATE ON vacancies
FOR EACH ROW EXECUTE FUNCTION update_search_vector();