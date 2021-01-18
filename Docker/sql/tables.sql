CREATE TABLE kahoot (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE question (
    id SERIAL PRIMARY KEY,
    question VARCHAR(255),
    description VARCHAR(255),
    kahoot_id INT NOT NULL,
    CONSTRAINT fk_question_kahoot
        foreign key (kahoot_id)
            REFERENCES kahoot(id)
);

CREATE TABLE answer (
    id SERIAL PRIMARY KEY,
    answer_id INT NOT NULL,
    description VARCHAR(255),
    id_question INT NOT NULL,
    is_true BOOLEAN,
    CONSTRAINT fk_answer_question
        foreign key (id_question)
            REFERENCES question(id)
)