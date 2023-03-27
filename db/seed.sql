--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3 (Ubuntu 16.3-1.pgdg22.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: device_type_ids; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.device_type_ids (id, device_type, device_type_id) FROM stdin;
1	ambient_wx_station	1
2	avtech_sensor	2
3	pi_pico_sensor	3
\.


--
-- Data for Name: sensors; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sensors (id, sensor_name, sensor_location, device_type_id) FROM stdin;
1	ambient_wx_station	backyard	1
2	avtech_basement_rack	basement_rack	2
3	bedroom_sensor	bedroom	3
4	basement_sensor	basement	3
5	office_sensor	office	3
\.


--
-- Name: device_type_ids_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.device_type_ids_id_seq', 3, true);


--
-- Name: sensors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sensors_id_seq', 5, true);


--
-- PostgreSQL database dump complete
--

