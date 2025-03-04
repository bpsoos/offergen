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

COPY public.items (id, owner_id, name, description, category, price) FROM stdin;
213c32b0-8cb2-4f49-ae17-7e94bc233db1	2d318789-1438-49ca-a51c-28927da2c68b	Muhammara	A rich and smoky red pepper dip blended with walnuts, pomegranate molasses, and aromatic spices. This Syrian specialty delivers a perfect balance of sweetness, spice, and nuttiness, served with warm pita.	\N	2900
46f73dd0-faba-4dd0-aaf7-f3c591b35c1c	2d318789-1438-49ca-a51c-28927da2c68b	Baba Ghanouj	\N	\N	1950
1d39681c-e9be-4282-a50b-8782176bbe38	2d318789-1438-49ca-a51c-28927da2c68b	Falafel Plate (Small)	Crispy, golden-brown falafel served with fresh vegetables, tahini sauce, and pita. A satisfying vegetarian favorite with bold Middle Eastern flavors.	\N	1290
b340cf31-e08d-4cf8-8ae3-b2e0fe8c3393	2d318789-1438-49ca-a51c-28927da2c68b	Quinoa Tabbouleh	\N	\N	2600
8dd53a2c-d8eb-463f-ac8e-dd9c76a2ef46	2d318789-1438-49ca-a51c-28927da2c68b	Falafel plate large	\N	\N	2450
e8869138-d238-433e-8147-f5f6c0f816e2	2d318789-1438-49ca-a51c-28927da2c68b	Avocado Wrap	A wholesome wrap filled with creamy avocado, crisp greens, and fresh veggies, all wrapped in a soft tortilla for a refreshing, nutritious bite.	\N	2400
8e2ead39-aacf-43df-9736-c162a416156b	2d318789-1438-49ca-a51c-28927da2c68b	TD Gold	\N	\N	2100
a4d65a6d-7dd8-492a-8ab1-6ed8a7179600	2d318789-1438-49ca-a51c-28927da2c68b	Hummus	A smooth and creamy blend of chickpeas, tahini, lemon, and garlic, drizzled with olive oil. This classic Mediterranean dip is perfect for pairing with warm pita or fresh veggies.	\N	2300
4c9a7a95-c5a7-4023-bc51-906210fd5dba	2d318789-1438-49ca-a51c-28927da2c68b	Warak Enab 	Tender grape leaves stuffed with a flavorful mix of rice, herbs, and spices. These bite-sized delights are a perfect balance of tangy and savory.	\N	1200
977b7cf2-5217-4f8e-b3e7-7c2c6d845893	2d318789-1438-49ca-a51c-28927da2c68b	Green Wrap	A fresh and vibrant wrap packed with crisp greens, cucumbers, herbs, and a light dressing, all wrapped in a soft tortilla for a refreshing and healthy meal.	\N	1000
0d262fd5-c460-4d9e-ae77-0b888ec2ba84	2d318789-1438-49ca-a51c-28927da2c68b	Black Cola	\N	\N	690
35d8fcc5-3c97-442d-b5ca-92741abcb79f	2d318789-1438-49ca-a51c-28927da2c68b	Red Cola	\N	\N	690
07429c02-1c26-4d43-962a-98dc65f2a0d4	2d318789-1438-49ca-a51c-28927da2c68b	Orange Frizzle	\N	\N	590
823b9b5a-5f46-44be-9ca5-caeacc5b21ec	2d318789-1438-49ca-a51c-28927da2c68b	Green Frizzle	\N	\N	490
4ecf8ccf-656e-4df2-b2d1-813b5a7c90a2	2d318789-1438-49ca-a51c-28927da2c68b	Plain Water	\N	\N	490
f3fe8a5d-cd4a-4370-a33c-b5ee6d8db5ec	2d318789-1438-49ca-a51c-28927da2c68b	Sparkly Water	\N	\N	390
16b0626b-99c1-4d87-a4e8-0c787feccb35	2d318789-1438-49ca-a51c-28927da2c68b	Craft beer	\N	\N	390
bf91b74a-1397-4265-8f12-151ae1ae3feb	2d318789-1438-49ca-a51c-28927da2c68b	Draught beer	\N	\N	990
\.


--
-- PostgreSQL database dump complete
--

