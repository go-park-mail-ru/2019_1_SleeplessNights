------------------------------------------------------------------------------------------------------------------------
-- TABLES --------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- table users
CREATE TABLE public.users
(
  id          BIGSERIAL NOT NULL,
  uid         BIGINT    NOT NULL,
  nickname    TEXT      NOT NULL,
  avatar_path TEXT      NOT NULL
);

ALTER TABLE ONLY public.users
  ADD CONSTRAINT users_pk PRIMARY KEY (id);

-- table rooms
CREATE TABLE public.rooms
(
  id    BIGSERIAL NOT NULL,
  users BIGINT[]
);

ALTER TABLE ONLY public.rooms
  ADD CONSTRAINT rooms_pk PRIMARY KEY (id);

-- table messages
CREATE TABLE public.messages
(
  id      BIGSERIAL NOT NULL,
  payload JSON      NOT NULL,
  user_id BIGINT    NOT NULL,
  room_id BIGINT    NOT NULL
);

ALTER TABLE ONLY public.messages
  ADD CONSTRAINT messages_pk PRIMARY KEY (id);

ALTER TABLE ONLY public.messages
  ADD CONSTRAINT "message_user_id_fk" FOREIGN KEY (user_id) REFERENCES public.users (id);

ALTER TABLE ONLY public.messages
  ADD CONSTRAINT "message_room_id_fk" FOREIGN KEY (room_id) REFERENCES public.rooms (id);


------------------------------------------------------------------------------------------------------------------------
-- INDEXES -------------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- index user_id_idx
CREATE INDEX user_id_idx ON public.messages USING btree (user_id);

-- index users_idx
CREATE INDEX users_idx ON public.rooms USING btree (users);

-- index room_id_idx
CREATE INDEX room_id_idx ON public.messages USING btree (room_id);


------------------------------------------------------------------------------------------------------------------------
-- FUNCTIONS -----------------------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------

-- func update_user
CREATE OR REPLACE FUNCTION public.func_update_user(arg_uid BIGINT,
                                                   arg_nickname TEXT,
                                                   arg_avatar_path TEXT)
  RETURNS BIGINT
AS
$BODY$
DECLARE
  result BIGINT;
BEGIN
  UPDATE public.users
  SET avatar_path = arg_avatar_path,
      nickname    = arg_nickname
  WHERE uid = arg_uid RETURNING id into result;
  IF NOT FOUND THEN
    INSERT INTO public.users (uid, nickname, avatar_path)
    VALUES (arg_uid, arg_nickname, arg_avatar_path) RETURNING id INTO result;
  END IF;
  RETURN result;
END;
$BODY$
  LANGUAGE plpgsql;

-- func add_message
CREATE OR REPLACE FUNCTION public.func_add_message(arg_user_id BIGINT,
                                                   arg_room_id BIGINT,
                                                   arg_payload JSON)
  RETURNS VOID
AS
$BODY$
BEGIN
  INSERT INTO public.messages (user_id, payload, room_id)
  VALUES (arg_user_id,
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
  INSERT INTO public.rooms (users)
  VALUES (arg_authors) RETURNING id INTO result;
  RETURN result;
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

-- func clean_db
CREATE OR REPLACE FUNCTION public.func_clean_db()
  RETURNS VOID
AS
$BODY$
BEGIN
  TRUNCATE TABLE public.rooms, public.messages, public.users RESTART IDENTITY;
END;
$BODY$
  LANGUAGE plpgsql;
