-- +migrate Up
CREATE TABLE "users" (
id VARCHAR(255) PRIMARY KEY,
first_name VARCHAR(255) NOT NULL,
last_name VARCHAR(255)  NOT NULL,
user_name VARCHAR(255) UNIQUE NOT NULL,
is_admin BOOLEAN NOT NULL,
password VARCHAR(255),
wage int NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
deleted_at TIMESTAMP);

CREATE TABLE "workplaces" (
id VARCHAR(255) PRIMARY KEY,
name VARCHAR(255) NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
deleted_at TIMESTAMP);

CREATE TABLE "benefits" (
id VARCHAR(255) PRIMARY KEY,
date TIMESTAMP NOT NULL,
reason VARCHAR(255)  NOT NULL,
amount int NOT NULL,
user_id VARCHAR(255) NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
deleted_at TIMESTAMP
);
ALTER TABLE benefits ADD CONSTRAINT benefit_users_fkey FOREIGN KEY (user_id) REFERENCES users(id);

CREATE TABLE "todos" (
id VARCHAR(255) PRIMARY KEY,
name VARCHAR(255) NOT NULL,
done BOOLEAN NOT NULL,
date TIMESTAMP NOT NULL,
benefit INT NOT NULL,
workplace_id VARCHAR(255) NOT NULL,
user_id VARCHAR(255) NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
deleted_at TIMESTAMP
);
ALTER TABLE todos ADD CONSTRAINT todo_users_fkey FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE todos ADD CONSTRAINT todo_workplaces_fkey FOREIGN KEY (workplace_id) REFERENCES workplaces(id);

CREATE TABLE "users_workplaces" (
workplace_id VARCHAR(255) NOT NULL,
user_id VARCHAR(255) NOT NULL
); 
ALTER TABLE users_workplaces ADD CONSTRAINT user_workplace_pkey PRIMARY KEY (user_id, workplace_id);
ALTER TABLE users_workplaces ADD CONSTRAINT user_workplace_workplaces_fkey FOREIGN KEY (workplace_id) REFERENCES workplaces(id);
ALTER TABLE users_workplaces ADD CONSTRAINT user_workplace_users_fkey FOREIGN KEY (user_id) REFERENCES users(id);

CREATE TABLE "tokens" (
id VARCHAR(255) PRIMARY KEY,
user_id VARCHAR(255) NOT NULL,
token VARCHAR(255)  NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
deleted_at TIMESTAMP
);
ALTER TABLE tokens ADD CONSTRAINT token_users_fkey FOREIGN KEY (user_id) REFERENCES users(id);

INSERT INTO users(id, first_name, last_name, user_name, is_admin, wage, created_at, updated_at)
VALUES
 ('59701caf-5b69-47f2-a1a6-76cebce23497', 'Admin', 'Admin', 'admin', true, 999, NOW(), NOW());
-- +migrate Down
DROP ALL TABLES;
