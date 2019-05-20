CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


------------------------------------------------------------------------------------------------------------------------
-- TABLES --------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- table users
CREATE TABLE public.users
(
  id          BIGSERIAL                  NOT NULL,
  email       CITEXT                     NOT NULL,
  password    BYTEA                      NOT NULL,
  salt        BYTEA                      NOT NULL,
  nickname    TEXT                       NOT NULL,
  avatar_path TEXT                       NOT NULL,
  win_rate    DOUBLE PRECISION DEFAULT 0 NOT NULL,
  matches     DOUBLE PRECISION DEFAULT 0 NOT NULL,
  wins        DOUBLE PRECISION DEFAULT 0 NOT NULL,
  rating      INT              DEFAULT 0 NOT NULL
);

--TODO create index by win_rate
--TODO create index by matches

CREATE UNIQUE INDEX users_email_ui
  ON public.users (email);

ALTER TABLE ONLY public.users
  ADD CONSTRAINT users_pk PRIMARY KEY (id);


------------------------------------------------------------------------------------------------------------------------
-- TRIGGERS ------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- trigger update_rating
CREATE OR REPLACE FUNCTION public.update_rating()
  RETURNS trigger AS
$BODY$
BEGIN
  NEW.win_rate := NEW.wins / NEW.matches;
  --   NEW.rating := OLD.rating + NEW.rating;
  RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql;

CREATE TRIGGER update_users_rating
  BEFORE UPDATE OF matches, wins
  ON public.users
  FOR EACH ROW
EXECUTE PROCEDURE public.update_rating();


------------------------------------------------------------------------------------------------------------------------
-- INDEXES -------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- index rating_idx
CREATE INDEX rating_idx ON public.users (rating DESC);


------------------------------------------------------------------------------------------------------------------------
-- TYPES ---------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- type user
CREATE TYPE public.type_user AS
  (
  id BIGINT,
  email CITEXT,
  password BYTEA,
  salt BYTEA,
  nickname TEXT,
  avatar_path TEXT,
  win_rate FLOAT,
  matches INTEGER,
  wins INTEGER,
  rating INTEGER
  );


------------------------------------------------------------------------------------------------------------------------
-- FUNCTIONS -----------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- func add_user
CREATE OR REPLACE FUNCTION public.func_add_user(arg_email CITEXT,
                                                arg_password BYTEA,
                                                arg_salt BYTEA,
                                                arg_nickname TEXT,
                                                arg_avatar_path TEXT)
  RETURNS public.type_user
AS
$BODY$
DECLARE
  result public.type_user;
BEGIN
  INSERT INTO public.users (email,
                            password,
                            salt,
                            nickname,
                            avatar_path)
  VALUES (arg_email,
          arg_password,
          arg_salt,
          arg_nickname,
          arg_avatar_path) RETURNING *
           INTO
             result.id,
             result.email,
             result.password,
             result.salt,
             result.nickname,
             result.avatar_path,
             result.win_rate,
             result.matches,
             result.wins,
             result.rating;
  RETURN result;
EXCEPTION
  WHEN unique_violation THEN
    RAISE unique_violation;
  WHEN integrity_constraint_violation THEN
    RAISE integrity_constraint_violation;
END;
$BODY$
  LANGUAGE plpgsql;

-- func clean_db
CREATE OR REPLACE FUNCTION public.func_clean_db()
  RETURNS VOID
AS
$BODY$
BEGIN
  TRUNCATE TABLE public.users RESTART IDENTITY;
END;
$BODY$
  LANGUAGE plpgsql;

-- func get_user_by_email
CREATE OR REPLACE FUNCTION public.func_get_user_by_email(arg_email CITEXT)
  RETURNS public.type_user
AS
$BODY$
DECLARE
  result public.type_user;
BEGIN
  SELECT * INTO
    result.id,
    result.email,
    result.password,
    result.salt,
    result.nickname,
    result.avatar_path,
    result.win_rate,
    result.matches,
    result.wins,
    result.rating
  FROM public.users
  WHERE email = arg_email;
  IF NOT FOUND THEN
    RAISE no_data_found;
  END IF;
  RETURN result;
END;
$BODY$
  LANGUAGE plpgsql;

-- func get_user_by_id
CREATE OR REPLACE FUNCTION public.func_get_user_by_id(arg_id BIGINT)
  RETURNS public.type_user
AS
$BODY$
DECLARE
  result public.type_user;
BEGIN
  SELECT * INTO
    result.id,
    result.email,
    result.password,
    result.salt,
    result.nickname,
    result.avatar_path,
    result.win_rate,
    result.matches,
    result.wins,
    result.rating
  FROM public.users
  WHERE id = arg_id;
  IF NOT FOUND THEN
    RAISE no_data_found;
  END IF;
  RETURN result;
END;
$BODY$
  LANGUAGE plpgsql;

-- func get_users
CREATE OR REPLACE FUNCTION public.func_get_users(arg_since BIGINT,
                                                 arg_limit BIGINT)
  RETURNS SETOF public.type_user
AS
$BODY$
DECLARE
  result public.type_user;
  rec    RECORD;
BEGIN
  FOR rec IN SELECT *
             FROM users
             WHERE id > arg_since
             ORDER BY rating
             LIMIT arg_limit OFFSET arg_since
    LOOP
      result.id := rec.id;
      result.email := rec.email;
      result.password := rec.password;
      result.salt := rec.salt;
      result.nickname := rec.nickname;
      result.avatar_path := rec.avatar_path;
      result.win_rate := rec.win_rate;
      result.matches := rec.matches;
      result.wins := rec.wins;
      result.rating := rec.rating;
      RETURN next result;
    END LOOP;
END;
$BODY$
  LANGUAGE plpgsql;

-- func update_user
CREATE OR REPLACE FUNCTION public.func_update_user(arg_nickname TEXT,
                                                   arg_avatar_path TEXT,
                                                   arg_id BIGINT)
  RETURNS public.type_user
AS
$BODY$
DECLARE
  result public.type_user;
BEGIN
  UPDATE public.users
  SET nickname    = CASE
                      WHEN arg_nickname = '' THEN nickname
                      ELSE arg_nickname END,
      avatar_path = CASE
                      WHEN arg_avatar_path = '' THEN avatar_path
                      ELSE arg_avatar_path END
  WHERE id = arg_id;
  IF NOT FOUND THEN
    RAISE no_data_found;
  END IF;
  RETURN result;
END;
$BODY$
  LANGUAGE plpgsql;

