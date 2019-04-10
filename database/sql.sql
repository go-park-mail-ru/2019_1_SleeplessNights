CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

-- table user
CREATE TABLE users
(
  id         SERIAL           NOT NULL,
  email      CITEXT           NOT NULL,
  password   BYTEA            NOT NULL,
  salt       BYTEA            NOT NULL,
  won        BIGINT DEFAULT 0 NOT NULL,
  lost       BIGINT DEFAULT 0 NOT NULL,
  playtime   BIGINT DEFAULT 0 NOT NULL,
  nickname   CITEXT           NOT NULL,
  avatarpath CITEXT           NOT NULL
);

ALTER TABLE public.users
  ADD CONSTRAINT users_pk PRIMARY KEY (email);
