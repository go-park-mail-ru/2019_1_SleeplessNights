CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

-- table user_manager
CREATE TABLE users
(
  id          BIGSERIAL NOT NULL,
  email       CITEXT    NOT NULL,
  password    BYTEA     NOT NULL,
  salt        BYTEA     NOT NULL,
  nickname    TEXT      NOT NULL,
  avatar_path TEXT      NOT NULL,
  win_rate    FLOAT     NOT NULL,
  matches     INTEGER   NOT NULL,
  wins        INTEGER   NOT NULL,
  rating      INTEGER   NOT NULL
);

--TODO create move sql to functions
--TODO create trigger for win_rate auto recount
--TODO create index by rating
--TODO create index by win_rate
--TODO create index by matches

CREATE UNIQUE INDEX users_email_ui
  ON public.users (email);

ALTER TABLE ONLY public.users
  ADD CONSTRAINT users_pk PRIMARY KEY (id);
