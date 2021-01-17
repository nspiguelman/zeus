CREATE TABLE kahoot (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE question (
    id BIGINT NOT NULL PRIMARY KEY,
    question VARCHAR(255),
    description VARCHAR(255),
    kahoot_id BIGINT NOT NULL,
    CONSTRAINT fk_question_kahoot
        foreign key (kahoot_id)
            REFERENCES kahoot(id)
);

CREATE TABLE answer (
    id BIGINT NOT NULL,
    answer_id BIGINT NOT NULL,
    description VARCHAR(255),
    id_question BIGINT NOT NULL,
    is_true BOOLEAN,
    PRIMARY KEY(id, answer_id),
    CONSTRAINT fk_answer_question
        foreign key (id_question)
            REFERENCES question(id)
)