-- USER ----------------------------------------------------------------------------------------------------------------
-- TABLES --------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

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


-- USER ----------------------------------------------------------------------------------------------------------------
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


-- USER ----------------------------------------------------------------------------------------------------------------
-- INDEXES -------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- index rating_idx
CREATE INDEX rating_idx ON public.users (rating DESC);


-- USER ----------------------------------------------------------------------------------------------------------------
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


-- USER ----------------------------------------------------------------------------------------------------------------
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

-- func clean_user_db
CREATE OR REPLACE FUNCTION public.func_clean_user_db()
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
      RETURN NEXT result;
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

-- CHAT ----------------------------------------------------------------------------------------------------------------
-- TABLES --------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- table talkers
CREATE TABLE public.talkers
(
  id          BIGSERIAL NOT NULL,
  uid         BIGINT    NOT NULL,
  nickname    TEXT      NOT NULL,
  avatar_path TEXT      NOT NULL
);

ALTER TABLE ONLY public.talkers
  ADD CONSTRAINT talkers_pk PRIMARY KEY (id);

-- table rooms
CREATE TABLE public.rooms
(
  id      BIGSERIAL NOT NULL,
  talkers BIGINT[]  NOT NULL
);

ALTER TABLE ONLY public.rooms
  ADD CONSTRAINT rooms_pk PRIMARY KEY (id);

-- table messages
CREATE TABLE public.messages
(
  id        BIGSERIAL NOT NULL,
  payload   JSON      NOT NULL,
  talker_id BIGINT    NOT NULL,
  room_id   BIGINT    NOT NULL
);

ALTER TABLE ONLY public.messages
  ADD CONSTRAINT messages_pk PRIMARY KEY (id);

ALTER TABLE ONLY public.messages
  ADD CONSTRAINT "message_talker_id_fk" FOREIGN KEY (talker_id) REFERENCES public.talkers (id);

ALTER TABLE ONLY public.messages
  ADD CONSTRAINT "message_room_id_fk" FOREIGN KEY (room_id) REFERENCES public.rooms (id) ON DELETE CASCADE;


-- CHAT ----------------------------------------------------------------------------------------------------------------
-- INDEXES -------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- index talker_id_idx
CREATE INDEX talker_id_idx ON public.messages USING btree (talker_id);

-- index talkers_idx
CREATE INDEX talkers_idx ON public.rooms USING btree (talkers);

-- index room_id_idx
CREATE INDEX room_id_idx ON public.messages USING btree (room_id);


-- CHAT ----------------------------------------------------------------------------------------------------------------
-- TYPES ---------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- type room
CREATE TYPE public.type_room AS
  (
  id BIGINT,
  talkers BIGINT[]
  );


-- CHAT ----------------------------------------------------------------------------------------------------------------
-- FUNCTIONS -----------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- func update_talker
CREATE OR REPLACE FUNCTION public.func_update_talker(arg_uid BIGINT,
                                                     arg_nickname TEXT,
                                                     arg_avatar_path TEXT)
  RETURNS BIGINT
AS
$BODY$
DECLARE
  result BIGINT;
BEGIN
  IF arg_uid = 0::BIGINT THEN
    INSERT INTO public.talkers (uid, nickname, avatar_path)
    VALUES (arg_uid, arg_nickname, arg_avatar_path) RETURNING id INTO result;
  ELSE
    UPDATE public.talkers
    SET avatar_path = arg_avatar_path,
        nickname    = arg_nickname
    WHERE uid = arg_uid RETURNING id into result;
    IF NOT FOUND THEN
      INSERT INTO public.talkers (uid, nickname, avatar_path)
      VALUES (arg_uid, arg_nickname, arg_avatar_path) RETURNING id INTO result;
    END IF;
  END IF;
  RETURN result;
END;
$BODY$
  LANGUAGE plpgsql;

-- func add_message
CREATE OR REPLACE FUNCTION public.func_add_message(arg_talker_id BIGINT,
                                                   arg_room_id BIGINT,
                                                   arg_payload JSON)
  RETURNS VOID
AS
$BODY$
BEGIN
  INSERT INTO public.messages (talker_id, payload, room_id)
  VALUES (arg_talker_id,
          arg_payload,
          arg_room_id);
END;
$BODY$
  LANGUAGE plpgsql;

-- func add_room
CREATE OR REPLACE FUNCTION public.func_add_room(arg_authors BIGINT[])
  RETURNS BIGINT
AS
$BODY$
DECLARE
  result BIGINT;
BEGIN
  INSERT INTO public.rooms (talkers)
  VALUES (arg_authors) RETURNING id INTO result;
  RETURN result;
END;
$BODY$
  LANGUAGE plpgsql;

-- func get_rooms
CREATE OR REPLACE FUNCTION public.func_get_rooms()
  RETURNS SETOF public.type_room
AS
$BODY$
DECLARE
  result public.type_room;
  rec    RECORD;
BEGIN
  FOR rec in SELECT * FROM public.rooms
    LOOP
      result.id := rec.id;
      result.talkers := rec.talkers;
      RETURN NEXT result;
    END LOOP;
END;
$BODY$
  LANGUAGE plpgsql;

-- func delete_room
CREATE OR REPLACE FUNCTION public.func_delete_room(arg_room_id BIGINT)
  RETURNS VOID
AS
$BODY$
BEGIN
  DELETE FROM public.rooms WHERE id = arg_room_id ;
END;
$BODY$
  LANGUAGE plpgsql;

-- func get_messages
CREATE OR REPLACE FUNCTION public.func_get_messages(arg_room_id BIGINT,
                                                    arg_since BIGINT,
                                                    arg_limit BIGINT)
  RETURNS SETOF JSON
AS
$BODY$
DECLARE
  result JSON;
  rec    RECORD;
BEGIN
  FOR rec in SELECT payload
             FROM public.messages
             WHERE id < arg_since
               AND room_id = arg_room_id
             LIMIT arg_limit
    LOOP
      result := rec.payload;
      RETURN NEXT result;
    END LOOP;
END;
$BODY$
  LANGUAGE plpgsql;

-- func clean_chat_db
CREATE OR REPLACE FUNCTION public.func_clean_chat_db()
  RETURNS VOID
AS
$BODY$
BEGIN
  TRUNCATE TABLE public.rooms, public.messages, public.talkers RESTART IDENTITY;
END;
$BODY$
  LANGUAGE plpgsql;
CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


------------------------------------------------------------------------------------------------------------------------
-- TABLES --------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- table Question'sPack
CREATE TABLE public.question_pack
(
  id        BIGSERIAL    NOT NULL,
  icon_path VARCHAR(100) NOT NULL,
  theme     VARCHAR(100) NOT NULL
);

ALTER TABLE ONLY public.question_pack
  ADD CONSTRAINT question_pack_pk PRIMARY KEY (id);

-- table question
CREATE TABLE public.question
(
  id      BIGSERIAL      NOT NULL,
  answers VARCHAR(200)[] NOT NULL,
  correct INTEGER        NOT NULL,
  text    TEXT           NOT NULL,
  pack_id BIGINT         NOT NULL
);

ALTER TABLE ONLY public.question
  ADD CONSTRAINT question_pk PRIMARY KEY (id);

ALTER TABLE ONLY public.question
  ADD CONSTRAINT "question_pack_id_fk" FOREIGN KEY (pack_id) REFERENCES public.question_pack (id);


------------------------------------------------------------------------------------------------------------------------
-- TYPES ---------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- type question_pack
CREATE TYPE public.type_question_pack AS
  (
  id BIGINT,
  icon_path VARCHAR(100),
  theme VARCHAR(100)
  );

-- type question
CREATE TYPE public.type_question AS
  (
  id BIGINT,
  answers VARCHAR(200)[],
  correct INTEGER,
  text TEXT,
  pack_id BIGINT
  );


------------------------------------------------------------------------------------------------------------------------
-- FUNCTIONS -----------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- func add_pack
CREATE OR REPLACE FUNCTION public.func_add_pack(arg_theme CITEXT,
                                                arg_icon_path CITEXT)
  RETURNS VOID
AS
$BODY$
BEGIN
  INSERT INTO public.question_pack (theme, icon_path)
  VALUES (arg_theme,
          arg_icon_path);
END;
$BODY$
  LANGUAGE plpgsql;

-- func add_question
CREATE OR REPLACE FUNCTION public.func_add_question(arg_answers VARCHAR(60)[],
                                                    arg_correct INTEGER,
                                                    arg_text TEXT,
                                                    arg_pack_id BIGINT)
  RETURNS VOID
AS
$BODY$
BEGIN
  INSERT INTO public.question (answers,
                               correct,
                               text,
                               pack_id)
  VALUES (arg_answers,
          arg_correct,
          arg_text,
          arg_pack_id);
END;
$BODY$
  LANGUAGE plpgsql;

-- func clean_db
CREATE OR REPLACE FUNCTION public.func_clean_game_db()
  RETURNS VOID
AS
$BODY$
BEGIN
  TRUNCATE TABLE public.question, public.question_pack RESTART IDENTITY;
END;
$BODY$
  LANGUAGE plpgsql;

-- func get_packs
CREATE OR REPLACE FUNCTION public.func_get_packs(arg_number BIGINT)
  RETURNS SETOF public.question_pack
AS
$BODY$
DECLARE
  result public.question_pack;
  rec    RECORD;
BEGIN
  FOR rec IN SELECT *
             FROM (SELECT DISTINCT ON (theme) *
                   FROM public.question_pack
                   ORDER BY theme) AS qp
             ORDER BY random()
             LIMIT arg_number
    LOOP
      result.id := rec.id;
      result.icon_path := rec.icon_path;
      result.theme := rec.theme;
      RETURN next result;
    END LOOP;
END;
$BODY$
  LANGUAGE plpgsql;

-- func get_questions
CREATE OR REPLACE FUNCTION public.func_get_questions(arg_ids BIGINT[])
  RETURNS SETOF public.question
AS
$BODY$
DECLARE
  result public.question;
  rec    RECORD;
BEGIN
  FOR rec IN SELECT *
             FROM public.question
             WHERE pack_id = ANY (arg_ids)
    LOOP
      result.id := rec.id;
      result.answers := rec.answers;
      result.correct := rec.correct;
      result.text := rec.text;
      result.pack_id := rec.pack_id;
      RETURN next result;
    END LOOP;
END;
$BODY$
  LANGUAGE plpgsql;
