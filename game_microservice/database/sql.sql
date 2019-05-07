CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

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
  id      BIGSERIAL     NOT NULL,
  answers VARCHAR(60)[] NOT NULL,
  correct INTEGER       NOT NULL,
  text    TEXT          NOT NULL,
  pack_id BIGINT        NOT NULL
);

ALTER TABLE ONLY public.question
  ADD CONSTRAINT question_pk PRIMARY KEY (id);
