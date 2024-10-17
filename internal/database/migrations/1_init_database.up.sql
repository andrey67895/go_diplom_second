CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS auth(
                      id uuid not null primary key default uuid_generate_v4(),
                      "login" text unique not null,
                      "hash_pass" text not null,
                      "hash_pass_master" text not null);

CREATE TABLE IF NOT EXISTS secret_type (
                                    id uuid not null primary key default uuid_generate_v4(),
                                    name varchar not null
);

CREATE TABLE IF NOT EXISTS secret (
                                    id uuid not null primary key default uuid_generate_v4(),
                                    encoded bytea,
                                    type uuid not null references secret_type(id),
                                    metadata text
);

CREATE TABLE IF NOT EXISTS auth_secret_ref (
                                                    auth_id uuid not null references auth(id),
                                                    secret_id uuid not null unique references secret(id),
                                                    PRIMARY KEY (auth_id, secret_id)
);

create index auth_secret_ref_auth_id_index
    on auth_secret_ref(auth_id);

INSERT INTO secret_type (id, name)
VALUES (DEFAULT, 'FILE');

INSERT INTO secret_type (id, name)
VALUES (DEFAULT, 'WORD');

INSERT INTO secret_type (id, name)
VALUES (DEFAULT, 'CARD');

INSERT INTO secret_type (id, name)
VALUES (DEFAULT, 'LOG_PASS');