--
-- PostgreSQL database dump
--

-- Dumped from database version 10.7 (Ubuntu 10.7-0ubuntu0.18.04.1)
-- Dumped by pg_dump version 10.7 (Ubuntu 10.7-0ubuntu0.18.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: Vote; Type: TYPE; Schema: public; Owner: maxim
--

CREATE TYPE public."Vote" AS ENUM (
    '-1',
    '1'
);


ALTER TYPE public."Vote" OWNER TO maxim;

--
-- Name: type_forum; Type: TYPE; Schema: public; Owner: maxim
--

CREATE TYPE public.type_forum AS (
	title character varying,
	"user" character varying,
	slug character varying,
	posts bigint,
	threads integer,
	is_new boolean
);


ALTER TYPE public.type_forum OWNER TO maxim;

--
-- Name: type_post; Type: TYPE; Schema: public; Owner: maxim
--

CREATE TYPE public.type_post AS (
	id bigint,
	parent bigint,
	author character varying,
	message character varying,
	"isEdited" boolean,
	forum character varying,
	thread bigint,
	created timestamp with time zone,
	is_new boolean
);


ALTER TYPE public.type_post OWNER TO maxim;

--
-- Name: type_post_data; Type: TYPE; Schema: public; Owner: maxim
--

CREATE TYPE public.type_post_data AS (
	parent bigint,
	author character varying,
	message character varying
);


ALTER TYPE public.type_post_data OWNER TO maxim;

--
-- Name: type_status; Type: TYPE; Schema: public; Owner: maxim
--

CREATE TYPE public.type_status AS (
	"user" integer,
	forum integer,
	thread integer,
	post integer
);


ALTER TYPE public.type_status OWNER TO maxim;

--
-- Name: type_thread; Type: TYPE; Schema: public; Owner: maxim
--

CREATE TYPE public.type_thread AS (
	is_new boolean,
	id bigint,
	title character varying(256),
	author character varying,
	forum character varying,
	message character varying,
	votes integer,
	slug character varying,
	created timestamp with time zone
);


ALTER TYPE public.type_thread OWNER TO maxim;

--
-- Name: type_user; Type: TYPE; Schema: public; Owner: maxim
--

CREATE TYPE public.type_user AS (
	is_new boolean,
	nickname character varying,
	fullname character varying,
	about character varying,
	email character varying
);


ALTER TYPE public.type_user OWNER TO maxim;

--
-- Name: do_not_use_func_thread_post_parent_tree(bigint, bigint, boolean, integer); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.do_not_use_func_thread_post_parent_tree(arg_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer DEFAULT 100) RETURNS SETOF public.type_post
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_post;
    rec RECORD;
BEGIN
    result.is_new := false;
    FOR rec IN (
        ----Рекурсивный запрос, делающий выборку одного дерева----
        WITH RECURSIVE tree(id, "parent-id", nickname, message, "is-edited", slug, "thread-id", created) AS (
            (SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            WHERE p."thread-id" = arg_id
            AND p."parent-id" IS NULL
            and CASE
                when arg_since_id is null then true
                WHEN arg_desc THEN p.id < arg_since_id
                ELSE p.id > arg_since_id
            END
            --ORDER BY
            --    (case WHEN arg_desc THEN p.created END) DESC,
            --    (CASE WHEN not arg_desc THEN p.created END) ASC
            LIMIT 1)
        UNION ALL
            SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            JOIn tree ON p."parent-id" = tree.id
            --ORDER BY
            --    (case WHEN arg_desc THEN p.created END) DESC,
            --    (CASE WHEN not arg_desc THEN p.created END) ASC
        )
        SELECT * FROM tree
    )
    LOOP
        result.id := rec.id;
        IF rec."parent-id" is null then
            result.parent := 0;
        ELSE
            result.parent := rec."parent-id";
        end if;
        result.author := rec.nickname;
        result.message := rec.message;
        result."isEdited" := rec."is-edited";
        result.forum := rec.slug;
        result.thread := rec."thread-id";
        result.created := rec.created;
        RETURN next result;
    END LOOP;
END;
$$;


ALTER FUNCTION public.do_not_use_func_thread_post_parent_tree(arg_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer) OWNER TO maxim;

--
-- Name: func_add_post_to_forum(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_add_post_to_forum() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
BEGIN
    UPDATE "Forums" SET posts = posts + 1 WHERE id = NEW."forum-id";
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.func_add_post_to_forum() OWNER TO maxim;

--
-- Name: func_add_thread_to_forum(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_add_thread_to_forum() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
BEGIN
    UPDATE "Forums" SET threads = threads + 1 WHERE id = NEW."forum-id";
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.func_add_thread_to_forum() OWNER TO maxim;

--
-- Name: func_add_vote_to_thread(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_add_vote_to_thread() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    voice int;
BEGIN
    voice := NEW.voice;
    UPDATE "Threads" SET votes = votes + voice WHERE id = NEW."thread-id";
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.func_add_vote_to_thread() OWNER TO maxim;

--
-- Name: func_check_post_before_adding(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_check_post_before_adding() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    parent RECORD;
    thread RECORD;
BEGIN
    IF  NEW."parent-id" IS NOT NULL
    AND NEW."parent-id" != 0 THEN
        SELECT * INTO parent
        from "Posts"
        WHERE id = NEW."parent-id";
    
        SELECT * into thread
        FROM "Threads"
        WHERE id = NEW."thread-id";
    
        if NEW."forum-id" != parent."forum-id"
        OR NEW."thread-id" != parent."thread-id"
        OR NEW."forum-id" != thread."forum-id"
        THEN
            RAISE integrity_constraint_violation;
        END IF;
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.func_check_post_before_adding() OWNER TO maxim;

--
-- Name: func_convert_post_parent_zero_into_null(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_convert_post_parent_zero_into_null() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
BEGIN
    IF NEW."parent-id" = 0 THEN
       NEW."parent-id" := NULL;
    END IF;
    RETURN NEW;
END; 
$$;


ALTER FUNCTION public.func_convert_post_parent_zero_into_null() OWNER TO maxim;

--
-- Name: func_convert_post_parent_zero_into_null$$(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public."func_convert_post_parent_zero_into_null$$"() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
BEGIN
    IF NEW."parent-id" = 0 THEN
       NEW."parent-id" := NULL;
    END IF;
    RETURN NEW;
END; 
$$;


ALTER FUNCTION public."func_convert_post_parent_zero_into_null$$"() OWNER TO maxim;

--
-- Name: func_delete_post_from_forum(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_delete_post_from_forum() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
BEGIN
  	UPDATE "Forums" SET posts = posts -1 WHERE id = OLD."forum-id";
  	RETURN NEW;
END;
$$;


ALTER FUNCTION public.func_delete_post_from_forum() OWNER TO maxim;

--
-- Name: func_delete_thread_from_forum(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_delete_thread_from_forum() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
BEGIN
    UPDATE "Forums" SET threads = threads - 1 WHERE id = OLD."forum-id";
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.func_delete_thread_from_forum() OWNER TO maxim;

--
-- Name: func_delete_vote_from_thread(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_delete_vote_from_thread() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    voice int;
BEGIN
    voice := OLD.voice;
    UPDATE "Threads" SET votes = votes - voice WHERE id = OLD."thread-id";
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.func_delete_vote_from_thread() OWNER TO maxim;

--
-- Name: func_edit_post(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_edit_post() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
BEGIN
    if OLD.message != NEW.message THEN
        NEW."is-edited" = TRUE;
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.func_edit_post() OWNER TO maxim;

--
-- Name: func_forum_create(character varying, character varying, character varying); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_forum_create(arg_title character varying, arg_nickname character varying, arg_slug character varying) RETURNS TABLE(res_is_new boolean, res_title character varying, res_user character varying, res_slug character varying, res_posts bigint, res_threads bigint)
    LANGUAGE plpgsql
    AS $$
DECLARE
    user_id BIGINT;
BEGIN
    SELECT id, nickname into user_id, res_user
    from "Users"
    WHERE lower(nickname) = lower(arg_nickname) FOR UPDATE;
    IF not found then
        RAISe no_data_found;
    end if;
    begin
        res_is_new := true;
        INSERT into "Forums"(title, "user-id", slug)
        VALUES (arg_title, user_id, arg_slug)
        RETURNING title, slug, posts, threads
        into res_title, res_slug, res_posts, res_threads;
        RETURN NEXT;
    EXCEPTION
        WHEN unique_violation THEN
            begin
                res_is_new := false;
                select f.title, f.slug, f.posts, f.threads
                into res_title, res_slug, res_posts, res_threads
                FROM "Forums" f
                where lower(f.slug) = lower(arg_slug);
                return NEXT;
            end;
    END;
END;
$$;


ALTER FUNCTION public.func_forum_create(arg_title character varying, arg_nickname character varying, arg_slug character varying) OWNER TO maxim;

--
-- Name: func_forum_create_thread(character varying, character varying, character varying, character varying, character varying, timestamp with time zone); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_forum_create_thread(arg_forum_slug character varying, arg_thread_slug character varying, arg_title character varying, arg_author character varying, arg_message character varying, arg_created timestamp with time zone) RETURNS public.type_thread
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_thread;
    user_id BIGINT;
    forum_id BIGINT;
BEGIN
    select id INTO user_id
    FROM "Users"
    WHERE lower(nickname) = lower(arg_author) for update;
    IF not found then
        RAISe no_data_found;
    end if;
    result.author := arg_author;
    
    SELECT f.id, f.slug INTO forum_id, result.forum
    FROM "Forums" f
    WHERE lower(f.slug) = lower(arg_forum_slug) FOR UPDATE;
    if not found then
        RAISe no_data_found;
    end if;
    
    begin
        result.is_new := true;
        INSERT INTO "Threads" ("author-id", created, "forum-id", message, slug, title)
        VALUES (user_id, arg_created, forum_id, arg_message,arg_thread_slug, arg_title)
        RETURNING id, title, message, votes, slug, created
        INTO result.id, result.title, result.message, result.votes, result.slug, result.created;
        RETURN result;
    EXCEPTION
        WHEN unique_violation THEN
            begin
                result.is_new := false;
                SELECT t.id, t.title, u.nickname, f.slug, t.message, t.votes, t.slug, t.created
                INTO result.id, result.title, result.author, result.forum, result.message, result.votes, result.slug, result.created
                FROM "Threads" t
                JOIN "Users" u ON u.id = t."author-id"
                JOIN "Forums" f ON f.id = t."forum-id"
                WHERE lower(t.slug) = lower(arg_thread_slug);
                return result;
            end;
    end;
END;
$$;


ALTER FUNCTION public.func_forum_create_thread(arg_forum_slug character varying, arg_thread_slug character varying, arg_title character varying, arg_author character varying, arg_message character varying, arg_created timestamp with time zone) OWNER TO maxim;

--
-- Name: func_forum_details(character varying); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_forum_details(arg_slug character varying) RETURNS public.type_forum
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_forum;
BEGIN
    result.is_new := FALSE;
    SELECT f.title, u.nickname, f.slug, f.posts, f.threads
    INTO result.title, result.user, result.slug, result.posts, result.threads
    FROM "Forums" f
    JOIN "Users" u ON u.id = f."user-id"
    WHERE lower(f.slug) = lower(arg_slug);
    if not found then
        RAISe no_data_found;
    end if;
    RETURN result;
END;
$$;


ALTER FUNCTION public.func_forum_details(arg_slug character varying) OWNER TO maxim;

--
-- Name: func_forum_threads(character varying, timestamp with time zone, boolean, integer); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_forum_threads(arg_slug character varying, arg_since timestamp with time zone, arg_desc boolean, arg_limit integer DEFAULT 100) RETURNS SETOF public.type_thread
    LANGUAGE plpgsql
    AS $$
DECLARE
    forum_id BIGINT;
    result type_thread;
    rec RECORD;
BEGIN
    result.is_new := false;
    
    SELECT id, slug into forum_id, result.forum
    from "Forums"
    where lower(slug) = lower(arg_slug) FOR UPDATE;
    if not found then
        RAISe no_data_found;
    end if;
    
    FOR rec in SELECT t.id, t.title, u.nickname, t.message, t.votes, t.slug, t.created
        FROM "Threads" t 
        JOIN "Users" u on u.id = t."author-id"
        WHERE t."forum-id" = forum_id
        and CASE
            when arg_since is null then true
            WHEN arg_desc THEN t.created <= arg_since
            ELSE t.created >= arg_since
        END
        ORDER BY
            (case WHEN arg_desc THEN t.created END) DESC,
            (CASE WHEN not arg_desc THEN t.created END) ASC
        LIMIT arg_limit
    LOOp
        result.id := rec.id;
        result.title := rec.title;
        result.author := rec.nickname;
        result.message := rec.message;
        result.votes := rec.votes;
        result.slug := rec.slug;
        result.created := rec.created;
        RETURN next result;
    end loop;
END;
$$;


ALTER FUNCTION public.func_forum_threads(arg_slug character varying, arg_since timestamp with time zone, arg_desc boolean, arg_limit integer) OWNER TO maxim;

--
-- Name: func_forum_users(character varying, character varying, boolean, integer); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_forum_users(arg_slug character varying, arg_since character varying, arg_desc boolean, arg_limit integer DEFAULT 100) RETURNS SETOF public.type_user
    LANGUAGE plpgsql
    AS $$
DECLARE
    forum_id BIGINT;
    result type_user;
    rec RECORD;
BEGIN
    result.is_new := false;
    
    SELECT id into forum_id
    from "Forums"
    where lower(slug) = lower(arg_slug) FOR UPDATE;
    if not found then
        RAISe no_data_found;
    end if;
    
    FOR rec in SELECT u.nickname, u.fullname, u.about, u.email
        FROM "Users" u
        join (
            SELECT DISTINCT "author-id" AS id
            FROM "Threads"
            where "forum-id" = forum_id
            UNION
            SELECT DISTINCT "author-id" AS id
            FROM "Posts"
            WHERE "forum-id" = forum_id
        ) forum_users on forum_users.id = u.id  
        WHERE CASE
            when arg_since is null then true
            WHEN arg_desc THEN lower(u.nickname)::bytea < lower(arg_since)::bytea
            ELSE lower(u.nickname)::bytea > lower(arg_since)::bytea
        END
        ORDER BY
            (case WHEN arg_desc THEN lower(u.nickname)::bytea END) DESC,
            (CASE WHEN not arg_desc THEN lower(u.nickname)::bytea END) ASC
        LIMIT arg_limit
    LOOp
        result.nickname := rec.nickname;
        result.fullname := rec.fullname;
        result.about := rec.about;
        result.email := rec.email;
        RETURN next result;
    end loop;
END;
$$;


ALTER FUNCTION public.func_forum_users(arg_slug character varying, arg_since character varying, arg_desc boolean, arg_limit integer) OWNER TO maxim;

--
-- Name: func_make_path_for_post(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_make_path_for_post() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    parent RECORD;
BEGIN
    IF  NEW."parent-id" IS NOT NULL
    AND NEW."parent-id" != 0
    THEN
        SELECT * INTO parent
        FROM "Posts"
        WHERE id = NEW."parent-id";
        NEW.path := parent.path || parent.id;
    END IF;
    RETURN NEW; 
END;
$$;


ALTER FUNCTION public.func_make_path_for_post() OWNER TO maxim;

--
-- Name: func_post_change(bigint, character varying); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_post_change(arg_id bigint, arg_post character varying) RETURNS public.type_post
    LANGUAGE plpgsql
    AS $$
DECLARE
    author_id BIGINT;
    forum_id BIGINT;
    result type_post;
BEGIN
    result.is_new := FALSE;
    UPDATE "Posts"
    SET message = CASE
        WHEN arg_post != '' THEN arg_post
        ELSE message END
    WHERE id = arg_id
    REturning id, "parent-id", "author-id", message, "is-edited", "forum-id", "thread-id", created
    INTO result.id, result.parent, author_id, result.message, result."isEdited", forum_id, result.thread, result.created;
    if not found then
        RAISe no_data_found;
    end if;
    
    SELECT nickname INTo result.author FROM "Users" where id = author_id FOR UPDATE;
    if not found then
        RAISe no_data_found;
    end if;
    
    SELECT slug InTO result.forum FROM "Forums" Where id = forum_id FOR UPDATE;
    if not found then
        RAISe no_data_found;
    end if;
    
    if result.parent is null then
        result.parent = 0;
    end if;
    
    return result;
END;
$$;


ALTER FUNCTION public.func_post_change(arg_id bigint, arg_post character varying) OWNER TO maxim;

--
-- Name: func_post_details(bigint); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_post_details(arg_id bigint) RETURNS public.type_post
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_post;
BEGIN
    result.is_new := false;
    SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created
    INTO result.id, result.parent, result.author, result.message, result."isEdited", result.forum, result.thread, result.created
    FROM "Posts" p
    JOIN "Users" u on u.id = p."author-id"
    JOIN "Forums" f ON f.id = p."forum-id"
    WHERE p.id = arg_id;
    if not found then
        RAISe no_data_found;
    end if;
    IF result.parent IS NULL THEN
        result.parent := 0;
    END IF;
    RETURN result;
END;
$$;


ALTER FUNCTION public.func_post_details(arg_id bigint) OWNER TO maxim;

--
-- Name: func_service_clear(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_service_clear() RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
BEGIN
    TRUNCATE TABLE "Users", "Forums", "Threads", "Posts", "Votes";
END;
$$;


ALTER FUNCTION public.func_service_clear() OWNER TO maxim;

--
-- Name: func_service_status(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_service_status() RETURNS public.type_status
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_status;
BEGIN
    SELECT count(*) INTO result.user FROM (
        SELECT * FROM "Users"
    ) u;
    SELECT count(*) INTO result.forum FROM (
        SELECT * FROM "Forums" f
    ) f;
    SELECT count(*) INTO result.thread FROM (
        SELECT * FROM "Threads"
    ) t;
    SELECT count(*) INTO result.post FROM (
        SELECT * FROM "Posts"
    ) p;
    RETURN result;
END;
$$;


ALTER FUNCTION public.func_service_status() OWNER TO maxim;

--
-- Name: func_thread_change(bigint, character varying, character varying); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_change(arg_id bigint, arg_title character varying, arg_message character varying) RETURNS public.type_thread
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_thread;
    user_id BIGINT;
    forum_id BIGINT;
BEGIN
    result.is_new := FALSE;
    UPDATE "Threads"
    SET title = CASE
            WHEN arg_title != '' THEN arg_title
            ELSE title END,
        message = CASE
            WHEN arg_message != '' THEN arg_message
            ELSE message END
    Where id = arg_id
    RETURNING id, title, "author-id", message, votes, slug, created, "forum-id"
    INTO result.id, result.title, user_id, result.message, result.votes, result.slug, result.created, forum_id;
    if not found then
        RAISe no_data_found;
    end if;
    SELECT nickname into result.author From "Users" WHERE id = user_id for update;
    if not found then
        RAISe no_data_found;
    end if;
    SELECT slug into result.forum From "Forums" WHERE id = forum_id for update;
    if not found then
        RAISe no_data_found;
    end if;
    return result;
END;
$$;


ALTER FUNCTION public.func_thread_change(arg_id bigint, arg_title character varying, arg_message character varying) OWNER TO maxim;

--
-- Name: func_thread_create_posts(bigint, bigint[], character varying[], character varying[]); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_create_posts(arg_id bigint, arg_parents bigint[], arg_authors character varying[], arg_messages character varying[]) RETURNS SETOF public.type_post
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_post;
    user_id BIGINT;
    forum_id BIGINT;
    array_length INTEGER;
    i Integer;
    post type_post_data;
BEGIN
    array_length := array_length(arg_parents, 1);
    IF array_length != array_length(arg_authors, 1)
    OR array_length != array_length(arg_messages, 1)
    THEN
        RAISE invalid_parameter_value;
    END IF;
    
    result.is_new := false;
    result.created := now(); 
    
    SELECT "forum-id" INTO forum_id
    FROM "Threads"
    WHERE id = arg_id;
    if not found then
        RAISe no_data_found;
    end if;
    result.thread := arg_id;
    
    SELECT slug INTO result.forum
    FROM "Forums"
    WHERE id = forum_id;
    if not found then
        RAISe no_data_found;
    end if;
    
    IF array_length is null then
        RETURN;
    END IF;
    
    i := 1;
    
    LOOP
        EXIT WHEN i > array_length;
        
        post.parent  := arg_parents[i];
        IF post.parent = 0 THEN
            post.parent := NULL;
        END IF;
        post.author  := arg_authors[i];
        post.message := arg_messages[i];
        
        if post.parent is null then
            result.parent = 0;
        else
            result.parent := post.parent;
        end if;   
    
        SELECT id into user_id
        from "Users"
        WHERE lower(nickname) = lower(post.author) FOR UPDATE;
        if not found then
            RAISe no_data_found;
        end if;
        result.author := post.author;
        
        INSERT into "Posts"("author-id", created, "forum-id", message, "parent-id", "thread-id")
        VALUES (user_id, result.created, forum_id, post.message, post.parent, arg_id)
        RETURNING  id, message, "is-edited"
        INTO result.id, result.message, result."isEdited";
        
        RETURN NEXT result;
        
        i := i + 1;
    END LOOP;
EXCEPTION
    WHEN unique_violation THEN
        RAISE unique_violation;
    WHEN integrity_constraint_violation THEN
        RAISE integrity_constraint_violation;
END;
$$;


ALTER FUNCTION public.func_thread_create_posts(arg_id bigint, arg_parents bigint[], arg_authors character varying[], arg_messages character varying[]) OWNER TO maxim;

--
-- Name: func_thread_details(bigint); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_details(arg_id bigint) RETURNS public.type_thread
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_thread;
BEGIN
    result.is_new := false;
    SELECT t.id, t.title, u.nickname, f.slug, t.message, t.votes, t.slug, t.created
    into result.id, result.title, result.author, result.forum, result.message, result.votes, result.slug, result.created
    from "Threads" t
    JOIN "Users" u ON u.id = t."author-id"
    JOIN "Forums" f ON f.id = t."forum-id"
    WHERE t.id = arg_id;
    if not found then
        RAISe no_data_found;
    end if;
    return result;
END;
$$;


ALTER FUNCTION public.func_thread_details(arg_id bigint) OWNER TO maxim;

--
-- Name: func_thread_get_id_by_slug(character varying); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_get_id_by_slug(arg_slug character varying) RETURNS bigint
    LANGUAGE plpgsql
    AS $$
DECLARE
    result BIGINT;
BEGIN
    SELECT id into result
    from "Threads"
    where lower(slug) = lower(arg_slug);
    if not found then
        RAISe no_data_found;
    end if;
    return result;
END;
$$;


ALTER FUNCTION public.func_thread_get_id_by_slug(arg_slug character varying) OWNER TO maxim;

--
-- Name: func_thread_get_post_layer(bigint, bigint, bigint, boolean, integer); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_get_post_layer(arg_thread_id bigint, arg_parent_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer DEFAULT NULL::integer) RETURNS SETOF public.type_post
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_post;
    rec RECORD;
BEGIN
    result.is_new := false;
    FOR rec in SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created
        FROM "Posts" p 
        JOIN "Users" u on u.id = p."author-id"
        JOIN "Forums" f ON f.id = p."forum-id"
        WHERE p."thread-id" = arg_thread_id
        and case
            WHEN arg_parent_id = 0 THEN p."parent-id" IS NULL
            ELSE p."parent-id" = arg_parent_id
        END
        and CASE
            when arg_since_id is null then true
            WHEN arg_desc THEN p.id < arg_since_id
            ELSE p.id > arg_since_id
        END
        ORDER BY
            (case WHEN arg_desc THEN p.created END) DESC,
            (CASE WHEN not arg_desc THEN p.created END) ASC
        LIMIT arg_limit
    LOOp
        result.id := rec.id;
         IF rec."parent-id" is null then
            result.parent := 0;
        ELSE
            result.parent := rec."parent-id";
        end if;
        result.author := rec.nickname;
        result.message := rec.message;
        result."isEdited" := rec."is-edited";
        result.forum := rec.slug;
        result.thread := rec."thread-id";
        result.created := rec.created;
        RETURN next result;
    end loop;
END;
$$;


ALTER FUNCTION public.func_thread_get_post_layer(arg_thread_id bigint, arg_parent_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer) OWNER TO maxim;

--
-- Name: func_thread_posts_flat(bigint, bigint, boolean, integer); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_posts_flat(arg_thread_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer DEFAULT 100) RETURNS SETOF public.type_post
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_post;
    rec RECORD;
BEGIN
    result.is_new := false;
    SELECT *
    Into rec
    from "Threads"
    where id = arg_thread_id FOR UPDATE;
    if not found then
        RAISe no_data_found;
    end if;
    
    IF arg_desc THEN
        FOR rec iN SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            WHERE p."thread-id" = arg_thread_id
            AND CASE
                when arg_since_id is null OR arg_since_id = 0 then true
                ELSE p.id < arg_since_id
            END
            ORDER BY p.id DESC
                --(case WHEN arg_desc THEN 1 END) DESC,
                --(CASE WHEN not arg_desc THEN 1 END) ASC
            LIMIT arg_limit
        LOOp
            result.id := rec.id;
            IF rec."parent-id" is null then
                result.parent := 0;
            ELSE
                result.parent := rec."parent-id";
            end if;
            result.author := rec.nickname;
            result.message := rec.message;
            result."isEdited" := rec."is-edited";
            result.forum := rec.slug;
            result.thread := rec."thread-id";
            result.created := rec.created;
            RETURN next result;
        end loop;
    ELSE
        FOR rec iN SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            WHERE p."thread-id" = arg_thread_id
            AND CASE
                when arg_since_id is null then true
                ELSE p.id > arg_since_id
            END
            ORDER BY p.id
                --(case WHEN arg_desc THEN 1 END) DESC,
                --(CASE WHEN not arg_desc THEN 1 END) ASC
            LIMIT arg_limit
        LOOp
            result.id := rec.id;
            IF rec."parent-id" is null then
                result.parent := 0;
            ELSE
                result.parent := rec."parent-id";
            end if;
            result.author := rec.nickname;
            result.message := rec.message;
            result."isEdited" := rec."is-edited";
            result.forum := rec.slug;
            result.thread := rec."thread-id";
            result.created := rec.created;
            RETURN next result;
        end loop;
    END IF;
END;
$$;


ALTER FUNCTION public.func_thread_posts_flat(arg_thread_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer) OWNER TO maxim;

--
-- Name: func_thread_posts_parent_tree(bigint, bigint, boolean, integer); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_posts_parent_tree(arg_thread_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer DEFAULT 100) RETURNS SETOF public.type_post
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_post;
    rec RECORD;
    root RECORD;
    since_root_id BIGINT;
    --max_branch_len INTEGER;
    depth INTEGER;
BEGIN
    SELECT *
    Into rec
    from "Threads"
    where id = arg_thread_id FOR UPDATE;
    if not found then
        RAISe no_data_found;
    end if;

    SELECT path[2]
    INTO since_root_id
    FROM "Posts"
    WHERE id = arg_since_id;
    
    result.is_new := false;
    IF arg_desc THEN
        --SELECT max(array_length(path,1))
        --INTO max_branch_len
        --FROM "Posts"
        --WHERE "thread-id" = arg_thread_id
        --AND CASE
        --    when since_root_id is null then true
        --    ELSE id < since_root_id
        --END;
    
        SELECT array_length(path,1)
        INTO depth
        FROM "Posts"
        WHERE id = arg_since_id;
        IF NOT FOUND THEN
            depth := 1;
        END IF;
        
        FOR root IN
            SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created, p.path || p.id as tree_path
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            --WHERE array_length(p.path,1) = max_branch_len
            WHERE p."parent-id" IS NULL
            AND p."thread-id" = arg_thread_id
            AND CASE
                when since_root_id is null then true
                ELSE p.id < since_root_id
            END
            --WHERE CASE
            --    when arg_since_id is null OR arg_since_id = 0 then true
            --    ELSE p.id < arg_since_id
            --END
            --AND array_length(p.path, 1) = depth
            ORDER BY p.id DESC
                --(case WHEN arg_desc THEN 1 END) DESC,
                --(CASE WHEN not arg_desc THEN 1 END) ASC
            LIMIT arg_limit
        LOOP
            --result.id := root.id;
            --IF root."parent-id" is null then
            --    result.parent := 0;
            --ELSE
            --    result.parent := root."parent-id";
            --end if;
            --result.author := root.nickname;
            --result.message := root.message;
            --result."isEdited" := root."is-edited";
            --result.forum := root.slug;
            --result.thread := root."thread-id";
            --result.created := root.created;
            --RETURN next result;
        
            FOR rec in
                SELECT * from func_thread_posts_tree_from_root(arg_thread_id, root.id, TRUE, NULL)
            LOOP
                RETURN next rec;
            END LOOP;
        END LOOP;
    ELSE
        FOR root IN
            SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created, p.path || p.id as tree_path
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            WHERE p."parent-id" IS NULL
            AND p."thread-id" = arg_thread_id
            AND CASE
                when since_root_id is null then true
                ELSE p.id > since_root_id
            END
            ORDER BY p.id
                --(case WHEN arg_desc THEN 1 END) DESC,
                --(CASE WHEN not arg_desc THEN 1 END) ASC
            LIMIT arg_limit
        LOOP
            --result.id := root.id;
            --IF root."parent-id" is null then
            --    result.parent := 0;
            --ELSE
            --    result.parent := root."parent-id";
            --end if;
            --result.author := root.nickname;
            --result.message := root.message;
            --result."isEdited" := root."is-edited";
            --result.forum := root.slug;
            --result.thread := root."thread-id";
            --result.created := root.created;
            --RETURN next result;
        
            FOR rec in
                SELECT * from func_thread_posts_tree_from_root(arg_thread_id, root.id, FALSE, NULL)
            LOOP
                RETURN next rec;
            END LOOP;
        END LOOP;
    END IF;
END;
$$;


ALTER FUNCTION public.func_thread_posts_parent_tree(arg_thread_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer) OWNER TO maxim;

--
-- Name: func_thread_posts_tree(bigint, bigint, boolean, integer); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_posts_tree(arg_thread_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer DEFAULT 100) RETURNS SETOF public.type_post
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_post;
    rec RECORD;
    since_path BIGINT[];
BEGIN
    result.is_new := false;
    
    SELECT *
    Into rec
    from "Threads"
    where id = arg_thread_id FOR UPDATE;
    if not found then
        RAISe no_data_found;
    end if;
    
    IF arg_since_id = 0 THEN
        arg_since_id = NULL;
    END IF;
    
    IF arg_since_id IS NOT NULL THEN
        SELECT (path || id) INTO since_path
        FROM "Posts"
        WHERE id = arg_since_id
        FOR UPDATE;
        IF not found THEN
            RAISE no_data_found;
        END IF;
    END IF;
    
    IF arg_desc THEN
        FOR rec IN
            SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created, p.path || p.id AS tree_path
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            WHERE p."thread-id" = arg_thread_id
            AND CASE
                when arg_since_id is null OR arg_since_id = 0 then true
                ELSE (p.path || p.id) < since_path
            END
            ORDER BY tree_path DESC, p.id
                --tree_path, (case WHEN arg_desc THEN 1 END) DESC,
                --tree_path, (CASE WHEN not arg_desc THEN 1 END) ASC
            LIMIT arg_limit
        LOOP
            result.id := rec.id;
            IF rec."parent-id" is null then
                result.parent := 0;
            ELSE
                result.parent := rec."parent-id";
            end if;
            result.author := rec.nickname;
            result.message := rec.message;
            result."isEdited" := rec."is-edited";
            result.forum := rec.slug;
            result.thread := rec."thread-id";
            result.created := rec.created;
            RETURN next result;
        END LOOP;
    ELSE
        FOR rec IN
            SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created, p.path || p.id AS tree_path
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            WHERE p."thread-id" = arg_thread_id
            AND CASE
                when arg_since_id is  null then true
                ELSE (p.path || p.id) > since_path
            END
            ORDER BY tree_path, p.id
                --tree_path, (case WHEN arg_desc THEN 1 END) DESC,
                --tree_path, (CASE WHEN not arg_desc THEN 1 END) ASC
            LIMIT arg_limit
        LOOP
            result.id := rec.id;
            IF rec."parent-id" is null then
                result.parent := 0;
            ELSE
                result.parent := rec."parent-id";
            end if;
            result.author := rec.nickname;
            result.message := rec.message;
            result."isEdited" := rec."is-edited";
            result.forum := rec.slug;
            result.thread := rec."thread-id";
            result.created := rec.created;
            RETURN next result;
        END LOOP;
    END IF;
END;
$$;


ALTER FUNCTION public.func_thread_posts_tree(arg_thread_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer) OWNER TO maxim;

--
-- Name: func_thread_posts_tree_from_root(bigint, bigint, boolean, integer); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_posts_tree_from_root(arg_thread_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer DEFAULT 100) RETURNS SETOF public.type_post
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_post;
    rec RECORD;
    depth INTEGER;
    since_parent BIGINT;
BEGIN
    result.is_new := false;
    
    IF arg_since_id = 0 THEN
        arg_since_id = NULL;
    END IF;
    
    IF arg_since_id IS NULL THEN
        depth := 1;
    ELSE
        SELECT array_length(path, 1) + 1 INTO depth
        FROM "Posts"
        WHERE id = arg_since_id
        FOR UPDATE;
        IF not found THEN
            RAISE no_data_found;
        END IF;
    END IF;
    
    IF arg_desc THEN
        FOR rec IN
            SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created, p.path || p.id AS tree_path
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            WHERE p."thread-id" = arg_thread_id
            --AND CASE
            --    WHEN arg_since_id IS NULL THEN TRUE
            --    ELSE p.id < arg_since_id
            --END
            AND CASE
                when arg_since_id is null then true
                ELSE (p.path || p.id)[depth] = arg_since_id
            END
            ORDER BY tree_path, p.id DESC
                --tree_path, (case WHEN arg_desc THEN 1 END) DESC,
                --tree_path, (CASE WHEN not arg_desc THEN 1 END) ASC
            LIMIT arg_limit
        LOOP
            result.id := rec.id;
            IF rec."parent-id" is null then
                result.parent := 0;
            ELSE
                result.parent := rec."parent-id";
            end if;
            result.author := rec.nickname;
            result.message := rec.message;
            result."isEdited" := rec."is-edited";
            result.forum := rec.slug;
            result.thread := rec."thread-id";
            result.created := rec.created;
            RETURN next result;
        END LOOP;
    ELSE
        FOR rec IN
            SELECT p.id, p."parent-id", u.nickname, p.message, p."is-edited", f.slug, p."thread-id", p.created, p.path || p.id AS tree_path
            FROM "Posts" p
            JOIN "Users" u on u.id = p."author-id"
            JOIN "Forums" f ON f.id = p."forum-id"
            WHERE p."thread-id" = arg_thread_id
            AND CASE
                when arg_since_id is  null then true
                ELSE (p.path || p.id)[depth] = arg_since_id
            END
            ORDER BY tree_path, p.id
                --tree_path, (case WHEN arg_desc THEN 1 END) DESC,
                --tree_path, (CASE WHEN not arg_desc THEN 1 END) ASC
            LIMIT arg_limit
        LOOP
            result.id := rec.id;
            IF rec."parent-id" is null then
                result.parent := 0;
            ELSE
                result.parent := rec."parent-id";
            end if;
            result.author := rec.nickname;
            result.message := rec.message;
            result."isEdited" := rec."is-edited";
            result.forum := rec.slug;
            result.thread := rec."thread-id";
            result.created := rec.created;
            RETURN next result;
        END LOOP;
    END IF;
END;
$$;


ALTER FUNCTION public.func_thread_posts_tree_from_root(arg_thread_id bigint, arg_since_id bigint, arg_desc boolean, arg_limit integer) OWNER TO maxim;

--
-- Name: func_thread_vote(bigint, character varying, boolean); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_thread_vote(arg_id bigint, arg_nickname character varying, arg_like boolean) RETURNS public.type_thread
    LANGUAGE plpgsql
    AS $$
DECLARE
    voice_val "public"."Vote";
    user_id BIGINT;
    result type_thread;
BEGIN
    result.is_new := false;
    if arg_like then
        voice_val := 1;
    else
        voice_val := -1;
    END IF;
    SELECT id into user_id
    from "Users"
    WHERE lower(nickname) = lower(arg_nickname) FOR UPDATE;
    IF NOT FOUND THEN
        RAISE no_data_found;
    END IF;
    INSERT into "Votes"("user-id", "thread-id", voice) VALUES (user_id, arg_id, voice_val);
    SELECT t.id, t.title, u.nickname, f.slug, t.message, t.votes, t.slug, t.created
    Into result.id, result.title, result.author, result.forum, result.message, result.votes, result.slug, result.created
    FROM "Threads" t 
    JOIN "Users" u on u.id = t."author-id"
    JOIN "Forums" f ON f.id = t."forum-id"
    WHERE t.id = arg_id;
    return result;
exception
    when unique_violation then
        UPDATE "Votes"
        SET voice = voice_val
        WHERE "user-id" = user_id
        AND "thread-id" = arg_id;
        SELECT t.id, t.title, u.nickname, f.slug, t.message, t.votes, t.slug, t.created
        Into result.id, result.title, result.author, result.forum, result.message, result.votes, result.slug, result.created
        FROM "Threads" t 
        JOIN "Users" u on u.id = t."author-id"
        JOIN "Forums" f ON f.id = t."forum-id"
        WHERE t.id = arg_id;
        return result;
    WHEN foreign_key_violation THEN
        RAISE no_data_found;
END;
$$;


ALTER FUNCTION public.func_thread_vote(arg_id bigint, arg_nickname character varying, arg_like boolean) OWNER TO maxim;

--
-- Name: func_update_vote(); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_update_vote() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    voice int;
BEGIN
    if (OLD.voice != NEW.voice)
    then
        voice := NEW.voice;
        voice := 2 * voice;
        UPDATE "Threads" SET votes = votes + voice WHERE id = NEW."thread-id";
    end if;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.func_update_vote() OWNER TO maxim;

--
-- Name: func_user_change_profile(character varying, character varying, character varying, character varying); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_user_change_profile(arg_nickname character varying, arg_fullname character varying, arg_about character varying, arg_email character varying) RETURNS public.type_user
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_user;
    user_id BIGINT;
BEGIN
    result.is_new := FALSE;
    SELECT id INTO user_id
    FROM "Users"
    Where lower(nickname) = lower(arg_nickname)
    FOR UPDATE;
    if not found then
        RAISe no_data_found;
    end if;
    
    UPDATE "Users"
    SET fullname = CASE
            WHEN arg_fullname != '' THEN arg_fullname
            ELSE fullname END,
        about = CASE
            WHEN arg_about != '' THEN arg_about
            ELSE about END,
        email = CASE
            WHEN arg_email != '' THEN arg_email
            ELSE email END
    Where id = user_id
    RETURNING nickname, fullname, about, email
    INTO result.nickname, result.fullname, result.about, result.email;
    return result;
exception
    when unique_violation THEN
        raise unique_violation;
END;
$$;


ALTER FUNCTION public.func_user_change_profile(arg_nickname character varying, arg_fullname character varying, arg_about character varying, arg_email character varying) OWNER TO maxim;

--
-- Name: func_user_create(character varying, character varying, character varying, character varying); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_user_create(arg_nickname character varying, arg_fullname character varying, arg_about character varying, arg_email character varying) RETURNS SETOF public.type_user
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_user;
    rec RECORD;
BEGIN
    begin
        result.is_new := true;
        INSERT INTO "Users" (nickname, fullname, about, email)
        VALUES (arg_nickname, arg_fullname, arg_about, arg_email)
        RETURNING nickname, fullname, about, email
        INTO result.nickname, result.fullname, result.about, result.email;
        RETURN next result;
    EXCEPTION
        WHEN unique_violation THEN
            begin
                result.is_new := false;
                FOR rec IN SELECT nickname, fullname, about, email
                    FROM "Users"
                    WHERE lower(nickname) = lower(arg_nickname)
                    OR lower(email) = lower(arg_email)
                LOOP
                    result.nickname := rec.nickname;
                    result.fullname := rec.fullname;
                    result.about := rec.about;
                    result.email := rec.email;
                    RETURN NEXT result;
                END LOOP;
            end;
    end;
END;
$$;


ALTER FUNCTION public.func_user_create(arg_nickname character varying, arg_fullname character varying, arg_about character varying, arg_email character varying) OWNER TO maxim;

--
-- Name: func_user_details(character varying); Type: FUNCTION; Schema: public; Owner: maxim
--

CREATE FUNCTION public.func_user_details(arg_nickname character varying) RETURNS public.type_user
    LANGUAGE plpgsql
    AS $$
DECLARE
    result type_user;
BEGIN
    result.is_new := FALSE;
    SELECT nickname, fullname, about, email
    INTO result.nickname, result.fullname, result.about, result.email
    FROM "Users"
    WHERE lower(nickname) = lower(arg_nickname);
    if not found then
        RAISe no_data_found;
    end if;
    RETURN result;
END;
$$;


ALTER FUNCTION public.func_user_details(arg_nickname character varying) OWNER TO maxim;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: Forums; Type: TABLE; Schema: public; Owner: maxim
--

CREATE TABLE public."Forums" (
    id bigint NOT NULL,
    posts bigint DEFAULT 0 NOT NULL,
    slug character varying(2044) NOT NULL,
    threads integer DEFAULT 0 NOT NULL,
    title character varying(256) NOT NULL,
    "user-id" bigint NOT NULL
);


ALTER TABLE public."Forums" OWNER TO maxim;

--
-- Name: Forum_id_seq; Type: SEQUENCE; Schema: public; Owner: maxim
--

CREATE SEQUENCE public."Forum_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Forum_id_seq" OWNER TO maxim;

--
-- Name: Forum_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maxim
--

ALTER SEQUENCE public."Forum_id_seq" OWNED BY public."Forums".id;


--
-- Name: Posts; Type: TABLE; Schema: public; Owner: maxim
--

CREATE TABLE public."Posts" (
    id bigint NOT NULL,
    "author-id" bigint NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    "forum-id" bigint NOT NULL,
    "is-edited" boolean DEFAULT false NOT NULL,
    message character varying(2044) NOT NULL,
    "parent-id" bigint,
    "thread-id" bigint NOT NULL,
    path bigint[] DEFAULT '{0}'::bigint[] NOT NULL
);


ALTER TABLE public."Posts" OWNER TO maxim;

--
-- Name: Post_id_seq; Type: SEQUENCE; Schema: public; Owner: maxim
--

CREATE SEQUENCE public."Post_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Post_id_seq" OWNER TO maxim;

--
-- Name: Post_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maxim
--

ALTER SEQUENCE public."Post_id_seq" OWNED BY public."Posts".id;


--
-- Name: Threads; Type: TABLE; Schema: public; Owner: maxim
--

CREATE TABLE public."Threads" (
    id bigint NOT NULL,
    "author-id" bigint NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    "forum-id" bigint NOT NULL,
    message character varying(2044) NOT NULL,
    slug character varying(2044) DEFAULT ''::character varying NOT NULL,
    title character varying(2044) NOT NULL,
    votes integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."Threads" OWNER TO maxim;

--
-- Name: Thread_id_seq; Type: SEQUENCE; Schema: public; Owner: maxim
--

CREATE SEQUENCE public."Thread_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Thread_id_seq" OWNER TO maxim;

--
-- Name: Thread_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maxim
--

ALTER SEQUENCE public."Thread_id_seq" OWNED BY public."Threads".id;


--
-- Name: Users; Type: TABLE; Schema: public; Owner: maxim
--

CREATE TABLE public."Users" (
    id bigint NOT NULL,
    about character varying(512) NOT NULL,
    email character varying(2044) NOT NULL,
    fullname character varying(128) NOT NULL,
    nickname character varying(2044) NOT NULL
);


ALTER TABLE public."Users" OWNER TO maxim;

--
-- Name: Users_id_seq; Type: SEQUENCE; Schema: public; Owner: maxim
--

CREATE SEQUENCE public."Users_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Users_id_seq" OWNER TO maxim;

--
-- Name: Users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maxim
--

ALTER SEQUENCE public."Users_id_seq" OWNED BY public."Users".id;


--
-- Name: Votes; Type: TABLE; Schema: public; Owner: maxim
--

CREATE TABLE public."Votes" (
    id bigint NOT NULL,
    voice public."Vote" NOT NULL,
    "thread-id" bigint NOT NULL,
    "user-id" bigint NOT NULL
);


ALTER TABLE public."Votes" OWNER TO maxim;

--
-- Name: Votes_id_seq; Type: SEQUENCE; Schema: public; Owner: maxim
--

CREATE SEQUENCE public."Votes_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Votes_id_seq" OWNER TO maxim;

--
-- Name: Votes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maxim
--

ALTER SEQUENCE public."Votes_id_seq" OWNED BY public."Votes".id;


--
-- Name: Forums id; Type: DEFAULT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Forums" ALTER COLUMN id SET DEFAULT nextval('public."Forum_id_seq"'::regclass);


--
-- Name: Posts id; Type: DEFAULT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Posts" ALTER COLUMN id SET DEFAULT nextval('public."Post_id_seq"'::regclass);


--
-- Name: Threads id; Type: DEFAULT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Threads" ALTER COLUMN id SET DEFAULT nextval('public."Thread_id_seq"'::regclass);


--
-- Name: Users id; Type: DEFAULT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Users" ALTER COLUMN id SET DEFAULT nextval('public."Users_id_seq"'::regclass);


--
-- Name: Votes id; Type: DEFAULT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Votes" ALTER COLUMN id SET DEFAULT nextval('public."Votes_id_seq"'::regclass);


--
-- Data for Name: Forums; Type: TABLE DATA; Schema: public; Owner: maxim
--

COPY public."Forums" (id, posts, slug, threads, title, "user-id") FROM stdin;
15702	1	uBk0kX92epI6k	1	Nec consonant suspirat eos vi, dicam lux hi mittere.	60628
15703	1	Sm6NS2WvX0-iR	1	Dura mendacio pro sit cur cotidie, sese verae recuperatae.	60631
15704	1	1tMN8X9eV06ik	1	Abs.	60634
15705	1	ZXy08E9VeN-oR	1	Quaerebant lustravi esau civium audio simul se praeteritorum talium.	60637
15706	1	w9T4R2WXv06-s	1	Vita didici sacramento tunc rupisti ebrietas filium dicens.	60640
15707	1	GsP4K2yX206oR	1	Umquam eum me dilabuntur.	60643
15708	1	VJkSeEU22PoiK	1	His pius me toleramus ob hac.	60646
15709	1	9HVRXVyxepIi8	1	Os prorsus cognosceremus confessio graeci displiceant dilabuntur voluptatibus.	60649
15710	1	SK3kx2WX2nOOre	1	Ad attamen oculos an elapsum.	60651
15711	1	LLC8exY2v0OIke	1	Ergo.	60653
15712	1	IZwRXVye2N-OKv	1	Si memoriam praeteritorum erro unus me virtus ex e.	60655
\.


--
-- Data for Name: Posts; Type: TABLE DATA; Schema: public; Owner: maxim
--

COPY public."Posts" (id, "author-id", created, "forum-id", "is-edited", message, "parent-id", "thread-id", path) FROM stdin;
124047	60627	2019-04-24 17:00:03.85542+03	15702	f	Si nota orare fama, ore. Sententiis coapta fabricasti pro neminem nares mortales, eundem passim discerem vel et. A opibus seu e quot veritatem, sinis. Avaritiam saucium cogit caritatis me ad si nos ibi multiplicitas abstinentia ab ascendam ac misertus ubi. Sensum ex heremo palpa liquida vi maerore meminissemus eam tria pede hic ne et voluisti viam adipiscendae. Haurimus. Pecora es an vere os, corporalis, habent ferre hic fui contineam magnus. In o senectute surditatem en paene vox retribuet o ne excipiens conprehendant interiore magni. Numquid ei vel pecora, peregrinorum requiro an e forte utrique remota.	\N	24227	{0}
124048	60630	2019-04-24 17:00:03.968071+03	15703	f	Ob vultu sui quidem tua via redarguentem plus pulsatori amplius, mel, locus immo praeciderim. Volo iustum. Idem transibo multa veris ingentibus lux soporem aliterque tui at grex, timeri porro praeterita, eos concordiam. Quibusve de at transire magis tamen, an ne pax esto an me.	\N	24228	{0}
124049	60633	2019-04-24 17:00:04.084339+03	15704	f	Vis salvi audiar suo vi omnino rapiatur.	\N	24229	{0}
124050	60636	2019-04-24 17:00:04.200541+03	15705	f	Vi meas non qui numeros o dicebam dona hac sui, cuiusque loco infirma sui cedendo porro eis.	\N	24230	{0}
124051	60639	2019-04-24 17:00:04.318146+03	15706	f	Habitas iubes da innumerabiles adduxerunt ac an est. Sum cur gustavi spe conspectu necessarium voluntate reficiatur possit mecum colligere, nulla de gradibus salvi. Conmemoravi sarcina credituri hi verus da in conceptaculum intra ibi quidem ad ago ago ascendant anima ac ac. Sinus olfactum innumerabilia seu genuit e suffragia difiniendo me gusta considero, adducor anhelo te una spes. Scis causa me certum erat diversitate tu repeterent eras artificiosas mors eam ait, his ubi. Careamus vocis vi. Pecco fit. Utcumque agebam perfecturum sane sui inlecebrosa facit angelos, vindicavit alium des. Ad alas es cum rem beatam en aliud delectamur volebant soni muta.	\N	24231	{0}
124052	60642	2019-04-24 17:00:04.432061+03	15707	f	Fuit rem sub et turpibus quaesitionum sui familiaritate, da, aperit sapores e solo infirmitati sonant fleo ea adhaesit. Dixi hos se passim vi ex o ea oblitumque alis cessas dixeris mundis, ore ob mutant esto has generis. Tu sine terrae vox potu es vicit vi flexu amari quanta id da cibo os potest seu. Auri israel ametur ei sitis ex da adversitas tu ut aer has quaesiveram ruga ducere putare gustatae via amor. Laude eo rei. Re sanctae eo. De ex ad domi qui laqueis sonant et veritas ea tuas at amaris sonat pendenda honoris in an. Cur ita inveniremus. Cantantur siderum memor temptationes si vis dum. Istae corda aestimanda an homines recordationis praebeo, sat cui. Casu vix mel sui pater, caveam sapiat solam latinique. Da qui duabus ab pollutum abundare e. Id hoc teque imago vestigio piae simile tua. Dari refrenare quantis mordeor mira, rem. Sanabis narro latis natura, die, dubitant.	\N	24232	{0}
124053	60645	2019-04-24 17:00:04.542775+03	15708	f	Veritatem tenebris nullum. Videndi ore dum videbat hominum inhaesero etiam corones malim pro si, sapientiorem ipsam hi e percipitur iacob e cum. Has eo ait manduco nutu, nam ex sanaturi, maerore. Propterea corones mutans e eam. Audiam tu dispulerim falsitate sui monuisti leporem venio abesset ne amore libeatque has id laetitiam.	\N	24233	{0}
124054	60648	2019-04-24 17:00:04.654783+03	15709	f	Amore die an consolatione die labor o idem a. Ea lux linguarum sensarum desideraris laudabunt beata coniunctam quod inhaeseram augebis variis tuae quam an. Dei e refero febris enim beatitudinis imagines in consuetudinem merito et. Tot ac istorum se calamitas, delectatio his subditi tuum ingentes amaremus supervacuanea peccatoris solitis an dextera. Saturari. Impium contemnere es satietas vos nimirum ita innumerabiles subinde has sui. Eloquia. Suo oblitum inlexit es eram, aves david genuit at, toleramus docentem hac capiendum sic sat vi. Vel. Hi sic quid gaudent suo populus auget, ei en. Aspectui nulla. Sim res ulla aer domine re duabus oculi artes a o habet cur aqua. Peccati te oblivionis pax. Meritis sane dicam stet huc opus utroque latinae aufer qua offendamus spargant o pulchra, iniquitatis eoque subditus. Omnibus bono intus tenacius, responsio interdum, iam nimis saeculum sub tria memoria attingi conmunem arrogantis en nascendo. Licet et vi commemoro se, qui beatitudinis sonet det voluisti ab lux faciens amat. Hac haereo incipio hos debeo os ex videri, me meliore deum, penetrale oportet potuero dicis vix usum das.	\N	24234	{0}
124055	60652	2019-04-24 17:00:04.74935+03	15710	t	PERFECTURUM EN CAREO SONUS AN VIX OBLECTAMENTA SI VIA NE VOLUIT SUA TAM EA, OBLIVISCAR. HIC HAC BENE AGITO. DE TUAM RESPICIENS SCIO TE QUANTIS EN VOX INVENIMUS ISTUC TANTO A SAT MEMORIAM. INFIRMIOR DULCIDINE. PER HI FIANT PATIENTER UT PERICULOSA CUI TU SCIUNT RE SCRIBENTUR, A BEATUM PRIVATIO TU RE ALIO MERITO. TALE RECONDIDI RIDICULUM IPSIS AMISI CUBILE EX DELECTOR ALTA DINOSCENS. EX MOVERI MAVULT QUAE PACEM NUTANTIBUS PLUS PULCHRIS GAUDEO. PER QUAM TENEBANT EN SI EO POTU EI MIRUM IBI HILARESCIT ALAS LAUDANDUM VISA AER TRISTITIA CUR, EN LENTICULAE.	\N	24235	{0}
124056	60654	2019-04-24 17:00:04.852487+03	15711	f	Esca facis latere confiteor et temptari liberasti sapiat mystice neglecta es, sciri sola carent cupimus quadam dividendo. Tenuissimas cui etiam solet psalterium inde sim melior ubi cotidianum, vicit potes ad fidei boni eum. Sit dicentium me cur animalia paratus. Melos ad mira ita cuncta das, vi ago doces curo. Si oboritur sic e malis, eos vix praesentia proximi sedem eo eum te habites nomen res rogo potens. A ubi via ut fit, teneam o fui eam. Ab memoriam nitidos commendavi invenimus vix tuus agenti nam in rogo perditum at. Ex. Eris de ut. Et. Potuero de exserentes.	\N	24236	{0}
124057	60656	2019-04-24 17:00:04.953941+03	15712	f	Accipiat habebunt istis flexu fac nos estis ad. Si vide fidei fuerunt tristitia vitam incideram huiuscemodi huc ad vigilanti licet da modestis eo desperarem mea. Admiratio da vestigio gaudent est mel extra securior tum. Iam ac o nos nec tui. Fixit re cito auram ac. Hos beatitudinis vivere hi videbam certe malo nonnumquam ita lata ex, dare orare. Facta cantandi una hi, non requiro nimii vis exhorreas, tam sparsis. Ferre haurimus ubi ad colligimus ne meae in pro. Sed faciei diu indidem o tu solus, narium languidus sapor pacem, vos piae gerit exitum metumve. Ipse ab inlecebras sublevas potuere caro passionis inlecebris ullo, reminiscimur bone se. Vel deum at licet actione fiat fundamenta tui an remisisti conmunem. Ut volebant das tam, verum voluit scit res in montes aditum ne hae.	\N	24237	{0}
\.


--
-- Data for Name: Threads; Type: TABLE DATA; Schema: public; Owner: maxim
--

COPY public."Threads" (id, "author-id", created, "forum-id", message, slug, title, votes) FROM stdin;
24227	60629	2019-04-07 21:31:33.238+03	15702	Haec tegitur re hi quaeram cui reminiscerer sua longum ore vera diversitate sub confessiones. Discendi. Periculo lucente typho deum si, vel, cordibus capio nominatur alius ac thesauri mole mutare quem es agro putem earum. Ne alteri esca cur modus cur similitudinem utriusque malus qui timeo tu opificiis peccatis silva contristamur recolenda. Credendum profero ea es sonare creditarum cibus tum voce bonam. Mel et et. Placere idoneus interiora quoquo nos, palpa ut eram vi non natura eosdem, mirificum sitis. Qua sui vocem liberamenta suavi iniuste miles vox nosti instat factus una praeciditur asperum des fac tertio suis.	k6xpkE92X0O6r	Sopitur se eras isto habet placent valde occurrat locutus.	0
24228	60632	2019-05-18 17:03:24.369+03	15703	Fueramus vae homo in si abigo o e vasis ut e pede soli. Bibo cibus eras re ex picturis mirantur se e fructus dura carnem sane facie ore vim aula exemplo ipso. Proceditur. Audito ob de quaererem de sinis aliam ex nos timere ferre hoc hi duxi da. Re hae alis addiderunt, variando, det adversus habemus aspectui modos significatur. Victoria a de vide enubiletur insidiis hi quoquo in solis a tuus turibulis eant ista tu. Os isto appetitu donec ex ad equus me re temptetur grave indicat concurrunt dicerem nigrum die me succurrat. Quot iubentis me. Filum filiis ab est hos iudex iam en christus en. Nutu solet quo laudatur fac piam auras sectantur suo, proruunt. En. Ac meos responderent mira dum ac numeramus, an capiamur placent.	CgIPRe9ExNOO8	Tui.	0
24229	60635	2018-06-08 19:36:57.968+03	15704	Typho itaque vi officiis. E vi gemitus nec iam prae. Nemo os se loquendo ea inpressit quantis mihi res istorum retranseo inruentes mutaveris es, obtentu abundabimus abigo. Re libidine quo. Emendicata talia eum et transcendi adsuefacta spe, malo perturbatione sit iussisti se ut et. Voluisti modulatione tui solem saluti, sensifico hae an artibus cantu adamavi cogeremur ei deo forsitan. Contristamur gero affectu os cubilia inde da at his nares, intervallis. Retenta per ago vana sim vel consolatione sit appareat visco iterum meo habet potest varia en donasti vitam. Salus ago modi quia vult sint praeditum loco flenda enervandam. Aqua ob aut tobis, ei propter, at nam unum disputando sancti molesta. Falsi distorta voluit. Bonis omnipotenti me suo de distinguere cepit invidentes fac es conmemoravi vae de locutum ne, vix die seu. Vicit tua hac ad mentem viva oleum ut dei reminiscenti qui. Cor mihi utrubique fit cessatione aliquam fames agit plenus via inplicentur, suaveolentiam sentio tuo cuiuslibet indicatae amo verba scis. O da odorum cor undique re mea eo proprios, pecco sint ab disseritur insinuat tota cupiditate, deviare tuo sui. Quia tanta dispersione nuda, ut terrae cor agnoscere pusilla promissio in doleam voluerit ab maneas te soporis ob.	sx5NK29ExP-6k	Accepisse.	0
24230	60638	2018-11-01 18:31:54.998+03	15705	Maeroribus laude aliterque me ei partes his beatum nesciebam quo amplexum. Mentem adhaesit magni tuo palleant manifestetur eo, diversitate, contraria coloratae. Verbum quadrupedibus hae nos. Quid duobus mecum dabis agit poterimus obliviscamur scit bibendi inhaesero eis ullo libeatque amarent moveri illuc ullo, agito.	wYwPSEux2NOoR	Nos comitatum.	0
24231	60641	2018-09-03 05:02:56.714+03	15706	Scierim mors innecto lege id sim, omnem hi sermo ei in aliquam scierim agro os si trium. Aliquantulum subeundam ne diei. Indicabo sic constrictione spes, factis fuit securior, improbat. Respondi imprimi quotiens te eris tam quibus es cubile, ego auget. Nuntiavit vae copiarum num cotidianas mare redigens quid quo pelluntur castrorum nostri hos pax interiore vana eis. Naturam qui offeretur melos locus insidiis. Amari si es artificosa, consensio. Os graece timore sint qua absurdissimum. Ubi ab. Iucunditas templi eras eas transfigurans norunt nescirem at spiritus eo laude. Agenti sentitur transibo sonos at sit sarcina. Et intentum. Manus vae de audiant res eventa o ex id catervas simile at vi. Nisi faciebat me cum visa a invocari cognovi verax voce laetamur exultans praesidens calidum quomodo iube sed deinde. Ait ore mors fac me tuo quia iesus.	lr1NkEU2Vp6iK	Narium quis.	0
24232	60644	2019-06-25 12:11:49.578+03	15707	An iam foris recognoscitur en responsio es affectent languor parva laudem, ac vanus corpore. Et contendunt cavis indicatae ubique oboritur fidelis stilo flatus. Ne viventis quidam surdis esca transeo vivere dignationem carentes noe filiorum ac indicavi. Si. Seu ore hae david ac at desiderant foris quaeratur paene imago visa es illum desiderant permanens at ad ob. Retrusa temptat hoc das, noe. Dona des. Copia nomino ei. Re sum interroget ubiubi aves avide eum vi invisibilia aliquid est dici cedentibus duobus solo fide. Re transit de o hic proprie. Similitudines me actiones novit familiaritate ut, nisi eo donec oculo os. E vocasti sua diu die has innotescunt haec. At vana ab a repetamus misericordiam peccator, ei pati num dicat contrectatae. Quod es egenus insania ut falsissime hanc necessaria die iam eo visione falsa eam incurrunt, simile at. Num amari via meo auri eis catervatim suam talia cordi eum temptari ac fateor unus deinde per vi.	K5P0Sxu2e4Iok	Item bonae fluctuo da hos qui consumma vis persuaserit.	0
24233	60647	2018-11-06 03:52:14.393+03	15708	Ne grave ipsam nec pluris nutu, pronuntianti mundi. Tui e nonnumquam corpore rem verbis noe tuo auras deo, usque diu rem te hos non alios. Rebus tum loqueretur operum et generis invenisse grave subinde beatus re evidentius, dei repositum factum, intentioni meas ridiculum. Sonuerunt audiar tu delicta libet iesus via o fit huc, dispulerim nuda, tum me. Sine horrendum illo adquiesco at at se isti possideri paulatim volo veris tu meo caecus agnoscimus sub. Viam non iactantia se eras detestor ea. Si pusillus audiat. Scio sui vox moderatum en coniugio, vocatur si. Sim aves potu sanum sanare. Fuisse adpropinquet eunt se, hebesco. Penitus quantum nitidos vos fastu meditor his alia ob cantu ambitum doce omnimodarum, sibimet in alii en. Invisibilia euge ita fructu laude temptationem suis miris ego. Corpori fugam vis inventa lene, satietate e stilo os referrem. Interpellat hos eam desideraris rationes solem defenditur sudoris castrorum. Ecce cuius tuo conectitur alimentorum tanto an atque, fui varia oculos cor utimur locutum odorum da recondidi fallax. Mare potuero tempore me. Lege animant eram re pulvis qua direxi pulchritudo omne exitum laetatum a futurae vivat vide nocte vae tui. Has exclamaverunt pacto. Contemnenda aer dulce etiam.	3_rrX2y2epi68	Ita tota alas.	0
24234	60650	2019-09-15 14:37:07.436+03	15709	Ea efficeret an recordabor cum fletus ex. Cor aer de mulier amplitudines parum si stet se sub abscondo dei abs si. Fui. Trahunt inperturbata ad eam excipiens os vellem hi. Conatur si tuam ipsi deserens iste parte toto verborum mala si te rem delet ad en unum pane specto.	pvosEXYVX46-8	Eum humilem exterminantes haeret seu.	0
24235	60652	2019-12-09 11:19:27.914+03	15710	Res surditatem te interrogare audivi gloriae ac vigilantes mirari duobus. Quo quadam mel cor eis da fit es de excusatio. In curiosa ac fui ut celeritate servientes beatus. Suo meminerim a animo vim gratiae an desperatione quo in cor dormies appareat copia at. Quaere serviunt ut sciunt, meruit cura inconcussus. Posside. Viae minuit olorem res labor per his fuerunt. Intus sua humilem disputandi eo, manifestus dignatus languidus certa. Potuere eius cui agro gustatae, sero vim in conduntur ita veritas amare sic eo odium vae praesentiam. Alas intellegentis intime ne unicus enim, occultum surditatem, sui ac cuiuslibet det spe factito thesauri solem mortales parvulus. Nosse iubens ac vim quis sumendi seu se amplius ulla eius at recedimus vos spe meo. Removeri transcendi a potui laudor sine fuit laudor ita splendorem ab omni ponamus o verax rei poenaliter pronuntianti. Fide. Omnesque ego.	8RlSv2W2E4-ir	Amo pulchras contristamur ait ex die miseria dixit abs.	0
24236	60654	2019-11-25 00:13:56.35+03	15711	Mei oleat sarcina serviant pede his numquam has fui ea aut da quaerens pater fui tota o fierem timere. Etiam tenacius at affectiones rogo lene varias ac, es cognoscendam si lenia hi de. Cibus ne pugno vivente, victor enumerat diu duabus o re admoneas an. Ac profunda vix an defrito ipsius. Tam ad hi potuero fecit mei tetigi oneri natura et minister notatum una trium immo cogenda. Via odore id. Quis amo redditur numeros. Fit. Ebriosus vi occurrit os vim, abs petitur an dicentem tetigi os iamque hi digna. Cupidatatium ex eloquia latine praeposita tali ubi sim modi confessio conmendat opus contemptu fletur ore. Nostra quo facit ne, interpellat res. Dei excellentiam prout muscas prodest, rapiunt sono auri veritatem resolvisti olet edunt esca ebriosus admitti, valeret ait vana.	lL5Rx2yxx0oik	Alias da ex putem, te cor plagas infligi.	0
24237	60656	2018-12-09 16:11:19.378+03	15712	Tuis autem habens aliae ne tribuere dum gradibus vix, numerorum aperit absentia homo absurdissimum victor meam alta te quaecumque. Das furens speciem amari tu una illico ingressae me fortius vi soni. Vocatur eo praegravatis coniunctione vis a agro fac sui nuntiavimus venio utendi ac possideas tria ad sermo. Proprie quippe hoc vis loco amo ex meo, magni resistere muta tenetur luminoso faciam ac partes apud. Melius pulvere varia et ab, hi de surgere. Sensu dolet ab sententiis redarguentem resistere neglecta restat ei alterum mea alias te se velut meus de clamat, officium. Apparens lingua rebus ne regem gressum a contristor diei volito de antris fructu. Amor alicui amat id caecus des. Vae reccido silva abscondo ne, modum id. E relaxari vis cotidianas, ut possit nova putant inventa es. Valde alia abs an. Fidei usque flatus aliquid sine sub concurrunt ubi assuescunt ministerium ad, pulchris mole diei ac e. Aliquid num oblitumque quo prodest ne hi o confessio iudicare imaginesque tu eo possim det hi perscrutanda quos. Habet ei gratiam actione conscribebat prout et intra tua ob deo tui rideat hi subiugaverant faciet subeundam, gaudebit me. Hae mei os meo tangendo nec animae en mala intellexisse nosse sequi vae aenigmate adparet vel facultas tum.	-Z98Ev9xE4-OR	Iubentis tu.	0
\.


--
-- Data for Name: Users; Type: TABLE DATA; Schema: public; Owner: maxim
--

COPY public."Users" (id, about, email, fullname, nickname) FROM stdin;
60627	Istum terram hos. Ad. Nos iamque audito. Redimas minutissimis de expertus ea tua amorem. Contrectatae tu os miris.	o.B8Jn71X11gZHj@earumeius.org	Elijah Thompson	venit.YT74pdtV126hpu
60628	Probet da viam evellas tam. Quaerimus doctrinae abiciam tuo, en. Plagas vi gloriatur parum. Invenit factumque dura possideri offensionem dari unde, oportet at. Absconditi ut. An infirmior audi eam mei habeas hymno. Abditis eo egenus meus. Hoc contristat sonos penuriam vos. Mirantur.	teneam.Tbj4pVS1d2h6J@exan.net	Addison Harris	e.TyR2JVsv14ZZJV
60629	Fundamenta illud abs an. Deteriore ne tui amor. Cavis quousque dormientis ab bonum sit. Pars esse hoc te illinc. Refero terrae es. Oblectamenta erat ecclesia suo vel drachmam ipse tuos.	animi.phUnRDXud4h6J@idoneustempli.com	Emma Robinson	plenis.7mVNPutVDnZhJ1
60630	Dici ob. Inter caelum tuis in interdum, inlusio hi. Piae pars inprobari praestat dicam, tecum siderum nota et. Iam theatra putare ea dixit suo es sane ea. Voce bonam malim die lux nisi lux.	me.3dhGr1X1d2mhJ@habitisum.com	Michael Harris	dixeris.luM2j1Sud2ZhJ1
60631	Oneri populi in. Aves os mea ex hos. Eras cogo dixi nominum, dona possint.	desperare.j3h4JD8Uv46Z7@alienotunc.com	Jacob Martin	nuntii.RL64Pv8vuGMmjV
60632	In filiorum una in. Fortasse conferunt cogitari. Canto vi. Os multiplicius. Iam nutu prohibuisti. Casu mei at.	sero.CQZGpUsU12zzR@itidemme.org	Mia Thompson	flammam.5Qh4Pu8uVg6HRu
60633	Es re suo. Vestigio nota offeratur. Catervatim talibus aut abs ei disputando e, suspirent. Verbum mel nos molestiam, vos famulatum. Aliam peccator celeritate ei hos stat. E lux. Hae boni. Quo fugam gemitus ad, eum peccator pius recordationem.	tetigi.7s9GPDxUd4ZH7@nimiasolis.org	Alexander Harris	gutture.PXl4r1X1unzzRD
60634	Bene. Una amo una iacitur en hic vi. Etsi fit. Interpellante tu. Peccatorum sacramento antiqua per grex et, reficiatur. Es es dulcidine o aer ob re. An fuit ab hi. Hac loqueremur mel copiosum falli os affectio, utcumque clamore.	retarder.aEK2PDSVdgMZ7@psalmisaluti.net	Aiden Martinez	adversus.aeLg71X1DNMmRv
60635	Turpis da sunt.	in.pUF2RdSvdN6Mr@detubi.org	Matthew Garcia	corporis.jvigJvsV1nhHju
60636	Vide e. Rem aliis o magnam esse dei bibo. Loco erro meretur sapores, retarder, vetare. Inquit. Eis abs dum eliqua vi corporalium verbum decus. Scire meae oboritur vituperare alius scio fac. Os ibi castam ventris cui. Sat liquida vi peccatoris spem.	si.PRsGPVxddGhhR@daturte.com	Zoey Taylor	e.7rsNj18uuGMZPd
60637	Fueramus concubitu si es. Potu tot mea nosse silva se ne. Vero hi bellum reminiscenti subiugaverant.	a.o18Gr1sVd2HHP@cavisbibendi.com	Ella Thomas	da.QuT2jvsV12Zzpu
60638	Intraverint at. Oportet vult fluctuo o fit sentio te ut hac. Totius est corporalis ei sic, re secum fueram. Tunc mella desuper at imprimi pars vae. Videndi antris desuper rem sequentes sitis transfigurans voluero. Lineas de peste malis melos. Carne sciunt sed des ubi tu multiplicius diiudicas mel.	placeam.stS4RdsdvGmh7@respicepro.org	Emily Jackson	nolo.XSS47UxVU4hZPd
60639	Aurium ac.	num.y3e2jUsvVGh6p@talilateat.net	Anthony Davis	sub.Y3qg7v8v14hh71
60640	Reminiscenti e ait. O delector hi. Te esset pecora iugo nos tot en donum temptat. Ad tempus errans. Vos ea at. Des quo eum qui.	a.8TE2p1xu14zHp@furenscupimus.com	Natalie Martinez	deo.T8w2pUs1UnzMrd
60641	Aeger. Rei mallem dei essem, si iaceat, tum. Hoc fudi opertis ei conmixta adtendi oculus affectus. Rei si digna. Fratribus.	viam.kJYGJuTuV46MP@quippemaris.org	Charlotte White	infirma.3JY27vtuu26Mp1
60642	Nihil gaudeam ei delet os bonis non venit, tua. Manet diei.	a.jn0n7usUug6H7@extrafalsum.org	Michael Wilson	e.J40nR1Tdu4hHPU
60643	Odorem hac adsurgat. Mei generis actionibus ea, magicas, sempiterna. O transcendi ipsas vix cui medice admiratio. Num foras quos detestetur amisi fidelis sum tenacius. Ebrietas.	metum.WPgGR1x1V4M6r@faciatlata.com	Lily White	genuit.wjn271T112ZZjd
60644	Falsissime. Tam eum vi vi amissum eum hi, est. Praesentiam os nunc reconditae sim aves. Creditarum intellexisse.	viribus.Pin4jDtvvnhh7@sudorismalint.org	David Wilson	o.rin4pDXUV2ZZr1
60645	Ex adsurgere usque quo melodias amor es. Rei oblectamenta vestigio vanias os infirmior. Ut dari species ego vera dicuntur mutare sim in. Conspectum hos ea sonuerunt. Benedicere. Speculum quaerit pax res tuam inconcussus fiant tu. Os et ut. Des verbo tot meas.	capiamur.Emj7U1sUv26M7@vaecomitum.org	Charlotte Jackson	a.e67pVusDvnmZrV
60646	Loquor fudi vivunt tuorum se dico. Ad aliis una o ne remota. Immo mea intus nonne, tobis latina, animarum.	re.HCRPuvXvUNHz7@fatemurcetera.net	Ethan Johnson	sacerdos.ZC771D8V12zzRd
60647	Has de via fuerunt. Illuc. Pleno habitaculum. Coniunctam moveri placuit. Habites dicuntur.	sui.KAjpvVXVuG6zp@diceremnaturae.org	Elizabeth Williams	eum.9Y7R1u81dgzz71
60648	Quaestio audio conantes avertit attigi.	ad.pwUpV1sV14Hzr@nutuipsa.org	Elijah Brown	ac.JOuJ1vSDu46HjV
60649	Vi docens quaesisse sim illam est flenda putant, cognosceremus.	e.TyUjuD81dGZhj@sapidaprobet.net	Abigail Davis	tactus.sbV7UvxvuGzzjU
60650	Sed doleat nonnullius in his eant doce huc. Passim. Percipitur eum ingerantur simus sane infligi. Unum.	inpressit.ndz7UUxV1nzMj@subirelata.org	Avery Miller	alas.gUm7uv8uDNMhjd
60651	Transeatur. Inmunditiam bona sum re multiplicitas, numquid nescirem. Et deserens eunt lux id cognitus, mare.	gustavi.JR3Pud8VdgmhJ@ianuasonare.com	Chloe Smith	a.PJKJdu8Vv46zRd
60652	Ingressae lascivos e soli mirum solum iumenti voce. Earum parit amet. Per ita eunt novi meo. Domi sive.	quidam.K33R1u8uVnh6j@totoimmensa.org	Mason Thomas	cognitus.l99711Tddghhp1
60653	Pulchras. Erigo infirmitas munerum ne qui huc. Respice maior spem sim faciant moderationi videtur noe. Abditis cogitando indica velit. Eligam temptationum habeatur et, prospera viae eripe rem piae. Hanc ex prodest. Oblivio raptae an dum quos.	re.93Fpd18uUg66R@cupiantcura.net	Liam Miller	penetrale.3kijdv8UUGZmJ1
60654	Mulier. Notus ab te sacerdos ecce nos.	desidero.5QFPuUXu1Nhhj@dain.org	Ethan Thomas	adlapsu.fOIJ11xdd4M6P1
60655	Est palpa.	aer.ZoS7VVTDU2ZzR@videturvivat.com	Emily Wilson	quae.zqtPvu8dU4MzrU
60656	Clauditur cohiberi alter equus, hi cito misericordiam laetatus. Omnis est vox sola per bene alii. Fuerunt. Proferatur rebus moveat quam tua.	formosa.2n8p1181D46z7@tuahomo.com	Mia Jones	haurimus.g2Sj1V8uVgmhPU
\.


--
-- Data for Name: Votes; Type: TABLE DATA; Schema: public; Owner: maxim
--

COPY public."Votes" (id, voice, "thread-id", "user-id") FROM stdin;
\.


--
-- Name: Forum_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maxim
--

SELECT pg_catalog.setval('public."Forum_id_seq"', 15712, true);


--
-- Name: Post_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maxim
--

SELECT pg_catalog.setval('public."Post_id_seq"', 124057, true);


--
-- Name: Thread_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maxim
--

SELECT pg_catalog.setval('public."Thread_id_seq"', 24237, true);


--
-- Name: Users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maxim
--

SELECT pg_catalog.setval('public."Users_id_seq"', 60656, true);


--
-- Name: Votes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maxim
--

SELECT pg_catalog.setval('public."Votes_id_seq"', 1147, true);


--
-- Name: Forums Forum_pkey; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Forums"
    ADD CONSTRAINT "Forum_pkey" PRIMARY KEY (id);


--
-- Name: Posts Post_pkey; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Posts"
    ADD CONSTRAINT "Post_pkey" PRIMARY KEY (id);


--
-- Name: Threads Thread_pkey; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Threads"
    ADD CONSTRAINT "Thread_pkey" PRIMARY KEY (id);


--
-- Name: Users Users_pkey; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "Users_pkey" PRIMARY KEY (id);


--
-- Name: Votes Votes_pkey; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Votes"
    ADD CONSTRAINT "Votes_pkey" PRIMARY KEY (id);


--
-- Name: Forums unique_Forum_id; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Forums"
    ADD CONSTRAINT "unique_Forum_id" UNIQUE (id);


--
-- Name: Posts unique_Post_id; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Posts"
    ADD CONSTRAINT "unique_Post_id" UNIQUE (id);


--
-- Name: Threads unique_Thread_id; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Threads"
    ADD CONSTRAINT "unique_Thread_id" UNIQUE (id);


--
-- Name: Users unique_Users_id; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "unique_Users_id" UNIQUE (id);


--
-- Name: Votes unique_Votes_id; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Votes"
    ADD CONSTRAINT "unique_Votes_id" UNIQUE (id);


--
-- Name: Votes unique_Votes_user_thread_pair; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Votes"
    ADD CONSTRAINT "unique_Votes_user_thread_pair" UNIQUE ("thread-id", "user-id");


--
-- Name: idx_forum_ci_slug; Type: INDEX; Schema: public; Owner: maxim
--

CREATE UNIQUE INDEX idx_forum_ci_slug ON public."Forums" USING btree (lower((slug)::text));


--
-- Name: idx_thread_ci_slug; Type: INDEX; Schema: public; Owner: maxim
--

CREATE UNIQUE INDEX idx_thread_ci_slug ON public."Threads" USING btree (lower((slug)::text)) WHERE ((slug)::text <> ''::text);


--
-- Name: idx_user_email; Type: INDEX; Schema: public; Owner: maxim
--

CREATE UNIQUE INDEX idx_user_email ON public."Users" USING btree (lower((email)::text));


--
-- Name: idx_user_nickname; Type: INDEX; Schema: public; Owner: maxim
--

CREATE UNIQUE INDEX idx_user_nickname ON public."Users" USING btree (lower((nickname)::text));


--
-- Name: Posts trg_add_post_to_forum; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_add_post_to_forum AFTER INSERT ON public."Posts" FOR EACH ROW EXECUTE PROCEDURE public.func_add_post_to_forum();


--
-- Name: Threads trg_add_thread_to_forum; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_add_thread_to_forum AFTER INSERT ON public."Threads" FOR EACH ROW EXECUTE PROCEDURE public.func_add_thread_to_forum();


--
-- Name: Votes trg_add_vote_to_thread; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_add_vote_to_thread AFTER INSERT ON public."Votes" FOR EACH ROW EXECUTE PROCEDURE public.func_add_vote_to_thread();


--
-- Name: Posts trg_check_post_before_adding; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_check_post_before_adding BEFORE INSERT OR UPDATE ON public."Posts" FOR EACH ROW EXECUTE PROCEDURE public.func_check_post_before_adding();


--
-- Name: Posts trg_convert_post_parent_zero_into_null; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_convert_post_parent_zero_into_null BEFORE INSERT OR UPDATE ON public."Posts" FOR EACH ROW EXECUTE PROCEDURE public.func_convert_post_parent_zero_into_null();


--
-- Name: Threads trg_delete thread_from_forum; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER "trg_delete thread_from_forum" AFTER DELETE ON public."Threads" FOR EACH ROW EXECUTE PROCEDURE public.func_delete_thread_from_forum();


--
-- Name: Posts trg_delete_post_from_forum; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_delete_post_from_forum BEFORE DELETE ON public."Posts" FOR EACH ROW EXECUTE PROCEDURE public.func_delete_post_from_forum();


--
-- Name: Votes trg_delete_vote_from_thread; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_delete_vote_from_thread AFTER DELETE ON public."Votes" FOR EACH ROW EXECUTE PROCEDURE public.func_delete_vote_from_thread();


--
-- Name: Posts trg_edit_post; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_edit_post BEFORE UPDATE ON public."Posts" FOR EACH ROW EXECUTE PROCEDURE public.func_edit_post();


--
-- Name: Posts trg_make_post_path; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_make_post_path BEFORE INSERT ON public."Posts" FOR EACH ROW EXECUTE PROCEDURE public.func_make_path_for_post();


--
-- Name: Votes trg_update_vote; Type: TRIGGER; Schema: public; Owner: maxim
--

CREATE TRIGGER trg_update_vote AFTER UPDATE ON public."Votes" FOR EACH ROW EXECUTE PROCEDURE public.func_update_vote();


--
-- Name: Posts lnk_Forums_Posts; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Posts"
    ADD CONSTRAINT "lnk_Forums_Posts" FOREIGN KEY ("forum-id") REFERENCES public."Forums"(id) MATCH FULL ON DELETE CASCADE;


--
-- Name: Threads lnk_Forums_Threads; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Threads"
    ADD CONSTRAINT "lnk_Forums_Threads" FOREIGN KEY ("forum-id") REFERENCES public."Forums"(id) MATCH FULL ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Posts lnk_Posts_Posts; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Posts"
    ADD CONSTRAINT "lnk_Posts_Posts" FOREIGN KEY ("parent-id") REFERENCES public."Posts"(id) MATCH FULL ON DELETE CASCADE;


--
-- Name: Posts lnk_Threads_Posts; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Posts"
    ADD CONSTRAINT "lnk_Threads_Posts" FOREIGN KEY ("thread-id") REFERENCES public."Threads"(id) MATCH FULL ON DELETE CASCADE;


--
-- Name: Votes lnk_Threads_Votes; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Votes"
    ADD CONSTRAINT "lnk_Threads_Votes" FOREIGN KEY ("thread-id") REFERENCES public."Threads"(id) MATCH FULL ON DELETE CASCADE;


--
-- Name: Forums lnk_Users_Forums; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Forums"
    ADD CONSTRAINT "lnk_Users_Forums" FOREIGN KEY ("user-id") REFERENCES public."Users"(id) MATCH FULL;


--
-- Name: Posts lnk_Users_Posts; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Posts"
    ADD CONSTRAINT "lnk_Users_Posts" FOREIGN KEY ("author-id") REFERENCES public."Users"(id) MATCH FULL;


--
-- Name: Threads lnk_Users_Threads; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Threads"
    ADD CONSTRAINT "lnk_Users_Threads" FOREIGN KEY ("author-id") REFERENCES public."Users"(id) MATCH FULL;


--
-- Name: Votes lnk_Users_Votes; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Votes"
    ADD CONSTRAINT "lnk_Users_Votes" FOREIGN KEY ("user-id") REFERENCES public."Users"(id) MATCH FULL;


--
-- PostgreSQL database dump complete
--

