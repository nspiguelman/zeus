CREATE TABLE kahoot (
    id SERIAL PRIMARY KEY,
    pin CHAR(6) NOT NULL,
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

CREATE TABLE kahootUser (
    id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(255),
    score INTEGER,
    kahoot_id BIGINT NOT NULL,
    CONSTRAINT fk_user_kahoot
        FOREIGN KEY (kahoot_id)
            REFERENCES kahoot(id)
)