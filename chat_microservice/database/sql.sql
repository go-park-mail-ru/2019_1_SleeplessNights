--
-- PostgreSQL database dump
--

-- Dumped from database version 10.7 (Ubuntu 10.7-0ubuntu0.18.04.1)
-- Dumped by pg_dump version 10.7 (Ubuntu 10.7-0ubuntu0.18.04.1)

CREATE SCHEMA public;

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


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: Authors; Type: TABLE; Schema: public; Owner: maxim
--

CREATE TABLE public."Authors"
(
  id          bigint                  NOT NULL,
  uid         bigint                  NOT NULL,
  nickname    character varying(2044) NOT NULL,
  avatar_path character varying(2044) NOT NULL
);


--
-- Name: Authors_id_seq; Type: SEQUENCE; Schema: public; Owner: maxim
--

CREATE SEQUENCE public."Authors_id_seq"
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;



--
-- Name: Authors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maxim
--

ALTER SEQUENCE public."Authors_id_seq" OWNED BY public."Authors".id;


--
-- Name: Messages; Type: TABLE; Schema: public; Owner: maxim
--

CREATE TABLE public."Messages"
(
  id        bigint NOT NULL,
  payload   json   NOT NULL,
  author_id bigint NOT NULL,
  room_id   bigint NOT NULL
);



--
-- Name: Message_id_seq; Type: SEQUENCE; Schema: public; Owner: maxim
--

CREATE SEQUENCE public."Message_id_seq"
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;



--
-- Name: Message_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maxim
--

ALTER SEQUENCE public."Message_id_seq" OWNED BY public."Messages".id;


--
-- Name: Rooms; Type: TABLE; Schema: public; Owner: maxim
--

CREATE TABLE public."Rooms"
(
  id      bigint NOT NULL,
  authors bigint[]
);



--
-- Name: Rooms_id_seq; Type: SEQUENCE; Schema: public; Owner: maxim
--

CREATE SEQUENCE public."Rooms_id_seq"
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


--
-- Name: Rooms_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maxim
--

ALTER SEQUENCE public."Rooms_id_seq" OWNED BY public."Rooms".id;


--
-- Name: Authors id; Type: DEFAULT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Authors"
  ALTER COLUMN id SET DEFAULT nextval('public."Authors_id_seq"'::regclass);


--
-- Name: Messages id; Type: DEFAULT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Messages"
  ALTER COLUMN id SET DEFAULT nextval('public."Message_id_seq"'::regclass);


--
-- Name: Rooms id; Type: DEFAULT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Rooms"
  ALTER COLUMN id SET DEFAULT nextval('public."Rooms_id_seq"'::regclass);

--
-- Name: Authors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maxim
--

SELECT pg_catalog.setval('public."Authors_id_seq"', 1, false);


--
-- Name: Message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maxim
--

SELECT pg_catalog.setval('public."Message_id_seq"', 1, false);


--
-- Name: Rooms_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maxim
--

SELECT pg_catalog.setval('public."Rooms_id_seq"', 1, false);


--
-- Name: Authors Authors_pkey; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Authors"
  ADD CONSTRAINT "Authors_pkey" PRIMARY KEY (id);


--
-- Name: Messages Message_pkey; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Messages"
  ADD CONSTRAINT "Message_pkey" PRIMARY KEY (id);


--
-- Name: Authors unique_Authors_id; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Authors"
  ADD CONSTRAINT "unique_Authors_id" UNIQUE (id);


--
-- Name: Authors unique_Authors_uid; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Authors"
  ADD CONSTRAINT "unique_Authors_uid" UNIQUE (uid);


--
-- Name: Messages unique_Message_id; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Messages"
  ADD CONSTRAINT "unique_Message_id" UNIQUE (id);


--
-- Name: Rooms unique_Rooms_id; Type: CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Rooms"
  ADD CONSTRAINT "unique_Rooms_id" PRIMARY KEY (id);


--
-- Name: index_author_id; Type: INDEX; Schema: public; Owner: maxim
--

CREATE INDEX index_author_id ON public."Messages" USING btree (author_id);


--
-- Name: index_authors; Type: INDEX; Schema: public; Owner: maxim
--

CREATE INDEX index_authors ON public."Rooms" USING btree (authors);


--
-- Name: index_room_id; Type: INDEX; Schema: public; Owner: maxim
--

CREATE INDEX index_room_id ON public."Messages" USING btree (room_id);


--
-- Name: Messages lnk_Authors_Messages; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Messages"
  ADD CONSTRAINT "lnk_Authors_Messages" FOREIGN KEY (author_id) REFERENCES public."Authors" (id) MATCH FULL ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Messages lnk_Rooms_Messages; Type: FK CONSTRAINT; Schema: public; Owner: maxim
--

ALTER TABLE ONLY public."Messages"
  ADD CONSTRAINT "lnk_Rooms_Messages" FOREIGN KEY (room_id) REFERENCES public."Rooms" (id) MATCH FULL ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

