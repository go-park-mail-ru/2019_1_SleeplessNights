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


-- table Question'sPack
CREATE TABLE question_pack
(
  id    BIGSERIAL    NOT NULL,
  theme VARCHAR(100) NOT NULL
);

ALTER TABLE ONLY public.question_pack
  ADD CONSTRAINT question_pack_pk PRIMARY KEY (id);


-- table question
CREATE TABLE question
(
  id         BIGSERIAL    NOT NULL,
  answers    VARCHAR(100)[]       NOT NULL,
  correct    INTEGER      NOT NULL,
  text       TEXT         NOT NULL,
  pack_id    BIGINT       NOT NULL
);

ALTER TABLE ONLY public.question
  ADD CONSTRAINT question_pk PRIMARY KEY (id);


