CREATE TABLE IF NOT EXISTS vacancy_skills (
    vacancy_id INT REFERENCES vacancies(id) ON DELETE CASCADE,
    skill_id   INT REFERENCES skills(id),
    PRIMARY KEY (vacancy_id, skill_id)
);

CREATE INDEX idx_vs_skill ON vacancy_skills(skill_id);