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
CREATE OR REPLACE FUNCTION public.func_clean_db()
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
