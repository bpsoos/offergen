--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2 (Debian 17.2-1.pgdg120+1)
-- Dumped by pg_dump version 17.2 (Debian 17.2-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: testuser
--

COPY public.users (id, email) FROM stdin;
2d318789-1438-49ca-a51c-28927da2c68b	test@user.com
\.


--
-- Data for Name: inventory; Type: TABLE DATA; Schema: public; Owner: testuser
--

COPY public.inventories (owner_id, title, is_published) FROM stdin;
2d318789-1438-49ca-a51c-28927da2c68b	Offering	f
\.


--
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: testuser
--

COPY public.items (id, owner_id, name, price) FROM stdin;
57145c97-ab4d-4418-969b-98a2d0b12dd1	2d318789-1438-49ca-a51c-28927da2c68b	Hummus	1800
46f73dd0-faba-4dd0-aaf7-f3c591b35c1c	2d318789-1438-49ca-a51c-28927da2c68b	Baba Ghanouj	1950
fffc5fec-c598-4ee0-bb11-2d948f4f1dfa	2d318789-1438-49ca-a51c-28927da2c68b	Muhammara	1950
e6e3551a-b18e-4dc9-b32f-8a49570527ac	2d318789-1438-49ca-a51c-28927da2c68b	Warak Enab	2150
b340cf31-e08d-4cf8-8ae3-b2e0fe8c3393	2d318789-1438-49ca-a51c-28927da2c68b	Quinoa Tabbouleh	2600
0348a9a2-546a-4385-8a9b-ca8b64aadbb7	2d318789-1438-49ca-a51c-28927da2c68b	Falafel plate small	1950
8dd53a2c-d8eb-463f-ac8e-dd9c76a2ef46	2d318789-1438-49ca-a51c-28927da2c68b	Falafel plate large	2450
d9e00039-23e1-474e-8521-9d3691078fad	2d318789-1438-49ca-a51c-28927da2c68b	Green wrap	2100
d692ed5f-f29c-4546-8c0c-f39ad32ee8da	2d318789-1438-49ca-a51c-28927da2c68b	Avocado wrap	2100
8e2ead39-aacf-43df-9736-c162a416156b	2d318789-1438-49ca-a51c-28927da2c68b	TD Gold	2100
0d262fd5-c460-4d9e-ae77-0b888ec2ba84	2d318789-1438-49ca-a51c-28927da2c68b	Black Cola	690
35d8fcc5-3c97-442d-b5ca-92741abcb79f	2d318789-1438-49ca-a51c-28927da2c68b	Red Cola	690
07429c02-1c26-4d43-962a-98dc65f2a0d4	2d318789-1438-49ca-a51c-28927da2c68b	Orange Frizzle	590
823b9b5a-5f46-44be-9ca5-caeacc5b21ec	2d318789-1438-49ca-a51c-28927da2c68b	Green Frizzle	490
4ecf8ccf-656e-4df2-b2d1-813b5a7c90a2	2d318789-1438-49ca-a51c-28927da2c68b	Plain Water	490
f3fe8a5d-cd4a-4370-a33c-b5ee6d8db5ec	2d318789-1438-49ca-a51c-28927da2c68b	Sparkly Water	390
16b0626b-99c1-4d87-a4e8-0c787feccb35	2d318789-1438-49ca-a51c-28927da2c68b	Craft beer	390
bf91b74a-1397-4265-8f12-151ae1ae3feb	2d318789-1438-49ca-a51c-28927da2c68b	Draught beer	990
\.


--
-- PostgreSQL database dump complete
--

