CREATE TABLE kahoot (
    id SERIAL PRIMARY KEY,
    pin CHAR(6) NOT NULL,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE question (
    id SERIAL PRIMARY KEY,
    question VARCHAR(255),
    description VARCHAR(255),
    kahoot_id BIGINT NOT NULL,
    CONSTRAINT fk_question_kahoot
        foreign key (kahoot_id)
            REFERENCES kahoot(id)
);

CREATE TABLE answer (
    id SERIAL NOT NULL PRIMARY KEY,
    description VARCHAR(255),
    question_id BIGINT NOT NULL,
    is_true BOOLEAN,
    CONSTRAINT fk_answer_question
        foreign key (question_id)
            REFERENCES question(id)
);

CREATE TABLE kahootUser (
    id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(255),
    token VARCHAR (200),
    kahoot_id BIGINT NOT NULL,
    CONSTRAINT fk_user_kahoot
        FOREIGN KEY (kahoot_id)
            REFERENCES kahoot(id)
);