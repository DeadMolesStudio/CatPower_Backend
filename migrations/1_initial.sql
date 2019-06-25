-- +migrate Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS user_profile (
    user_id serial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    password varchar(64) NOT NULL,

    username citext UNIQUE NOT NULL
--     avatar text
);

CREATE TABLE IF NOT EXISTS money_category (
    category_id serial PRIMARY KEY,
    user_id integer REFERENCES user_profile NOT NULL,
    name citext NOT NULL,
    is_income boolean DEFAULT FALSE, -- FALSE – consumption
    pic text,

    sum integer DEFAULT 0 -- denormalization
);

CREATE TABLE IF NOT EXISTS money_action (
    action_uuid text PRIMARY KEY,
    user_id integer REFERENCES user_profile NOT NULL,
    delta integer NOT NULL,
    from_category integer REFERENCES money_category NOT NULL,
    to_category integer REFERENCES money_category NOT NULL,
    photo text,

    added timestamp with time zone DEFAULT now()
);

ALTER DATABASE catpower SET timezone TO 'UTC-3';

-- +migrate Down
ALTER DATABASE catpower SET timezone TO 'UTC';

DROP TABLE IF EXISTS money_action;
DROP TABLE IF EXISTS money_category;
DROP TABLE IF EXISTS user_profile;

DROP EXTENSION IF EXISTS citext;
