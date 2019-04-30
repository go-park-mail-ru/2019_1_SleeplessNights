CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

-- table user
CREATE TABLE users
(
  id         BIGSERIAL        NOT NULL,
  email      CITEXT           NOT NULL,
  password   BYTEA            NOT NULL,
  salt       BYTEA            NOT NULL,
  won        BIGINT DEFAULT 0 NOT NULL,
  lost       BIGINT DEFAULT 0 NOT NULL,
  playtime   BIGINT DEFAULT 0 NOT NULL,
  nickname   TEXT             NOT NULL,
  avatarpath TEXT             NOT NULL
);

CREATE UNIQUE INDEX users_email_ui
  ON public.users (email);

ALTER TABLE ONLY public.users
  ADD CONSTRAINT users_pk PRIMARY KEY (id);
