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

INSERT INTO genders VALUES ('MALE');
INSERT INTO genders VALUES ('FEMALE');

INSERT INTO users VALUES ('Ve3be104c-9c8f-46da-b837-663a03c4d15d', 'Jane Smith' ,'janesmith@email.com', '2000-01-01', 'FEMALE', 'janeSmith');
INSERT INTO users VALUES ('9b9d25e0-fd5e-4432-b630-67668b00394f', 'Zara Ali' ,'zaraali@email.com', '1990-05-26', 'FEMALE', 'zaraAli');
INSERT INTO users VALUES ('c097ff70-eb72-4d13-a383-b6b7fb115aa8', 'Mia Rodriguez' ,'miarodriguez@email.com', '1996-10-14', 'FEMALE', 'miaRodriguez');

INSERT INTO users VALUES ('fb1a0b06-3dfb-41a9-b38c-3dc47f37f956', 'Zain Gill' ,'zaingill@email.com', '1997-03-14', 'MALE', 'zainGill');
INSERT INTO users VALUES ('3bbcd2eb-7599-4bb0-99dd-4bf03e591fff', 'John Smith' ,'johnsmith@email.com', '1998-12-03', 'MALE', 'johnSmith');