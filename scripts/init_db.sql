CREATE DATABASE dating_users;

USE dating_users;

CREATE TABLE genders(
    gender VARCHAR(128),
    PRIMARY KEY (gender)
);

CREATE TABLE users(
    id varchar(256),
    name varchar(256) NOT NULL,
    email varchar(256) NOT NULL,
    date_of_birth date NOT NULL,
    gender VARCHAR(128) NOT NULL,
    password  varchar(256) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (email),
    CONSTRAINT fk_userGender FOREIGN KEY (gender) REFERENCES genders(gender)
);

CREATE TABLE preferences(
    user_id varchar(256),
    preference_gender  VARCHAR(128) NOT NULL,
    max_age INTEGER NULL,
    min_age INTEGER NOT NULL,
    max_distance INTEGER NOT NULL,
    PRIMARY KEY (user_id),
    CONSTRAINT fk_preferenceGender FOREIGN KEY (preference_gender) REFERENCES genders(gender),
    CONSTRAINT fk_preferenceUser FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE locations(
    user_id varchar(256),
    lat FLOAT,
    `long` FLOAT,
    PRIMARY KEY (user_id),
    CONSTRAINT fk_locationUser FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO genders VALUES ('MALE');
INSERT INTO genders VALUES ('FEMALE');

INSERT INTO users VALUES ('Ve3be104c-9c8f-46da-b837-663a03c4d15d', 'Jane Smith' ,'janesmith@email.com', '2000-01-01', 'FEMALE', 'janeSmith');
INSERT INTO users VALUES ('9b9d25e0-fd5e-4432-b630-67668b00394f', 'Zara Ali' ,'zaraali@email.com', '1986-05-26', 'FEMALE', 'zaraAli');
INSERT INTO users VALUES ('c097ff70-eb72-4d13-a383-b6b7fb115aa8', 'Mia Rodriguez' ,'miarodriguez@email.com', '1996-10-14', 'FEMALE', 'miaRodriguez');

INSERT INTO users VALUES ('fb1a0b06-3dfb-41a9-b38c-3dc47f37f956', 'Zain Gill' ,'zaingill@email.com', '1997-03-14', 'MALE', 'zainGill');
INSERT INTO users VALUES ('3bbcd2eb-7599-4bb0-99dd-4bf03e591fff', 'John Smith' ,'johnsmith@email.com', '1984-12-03', 'MALE', 'johnSmith');

INSERT INTO preferences VALUES ('Ve3be104c-9c8f-46da-b837-663a03c4d15d', 'MALE', 50, 18, 200);
INSERT INTO preferences VALUES ('9b9d25e0-fd5e-4432-b630-67668b00394f', 'MALE', 45, 25, 50);
INSERT INTO preferences VALUES ('c097ff70-eb72-4d13-a383-b6b7fb115aa8', 'MALE', 27, 18, 50);

INSERT INTO preferences VALUES ('fb1a0b06-3dfb-41a9-b38c-3dc47f37f956', 'FEMALE', 30, 23, 75);
INSERT INTO preferences VALUES ('3bbcd2eb-7599-4bb0-99dd-4bf03e591fff', 'FEMALE', 55, 30, 75);

INSERT INTO locations VALUES ('Ve3be104c-9c8f-46da-b837-663a03c4d15d', 51.5055, 0.00754);
INSERT INTO locations VALUES ('9b9d25e0-fd5e-4432-b630-67668b00394f', 51.5007, -0.1246);
INSERT INTO locations VALUES ('c097ff70-eb72-4d13-a383-b6b7fb115aa8', 53.4808, -2.2426);

INSERT INTO locations VALUES ('fb1a0b06-3dfb-41a9-b38c-3dc47f37f956', 51.5072, -0.1247);
INSERT INTO locations VALUES ('3bbcd2eb-7599-4bb0-99dd-4bf03e591fff', 51.4934, 0.0098);