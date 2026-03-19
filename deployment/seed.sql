--
-- PostgreSQL database dump (Taskify Sample Data)
--

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

-- 1. Disable constraints and triggers for clean restore
SET session_replication_role = 'replica';

-- 2. Clean existing data (except schema_migrations)
TRUNCATE public.user_roles, public.role_permissions, public.tasks, public.users, public.roles, public.permissions, public.casbin_rule RESTART IDENTITY CASCADE;

--
-- Data for Name: casbin_rule; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (1, 'p', 'admin', '/api/v1/*', '*', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (2, 'p', 'employee', '/api/v1/auth/signout', 'POST', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (3, 'p', 'employee', '/api/v1/auth/refresh', 'POST', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (4, 'p', 'employee', '/api/v1/users/permissions', 'GET', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (5, 'p', 'employee', '/api/v1/users/profile', 'GET', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (6, 'p', 'employee', '/api/v1/users/profile', 'PATCH', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (7, 'p', 'employee', '/api/v1/users/profile', 'DELETE', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (8, 'p', 'employee', '/api/v1/users/change-password', 'PATCH', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (9, 'p', 'employee', '/api/v1/users/avatar', 'POST', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (10, 'p', 'employee', '/api/v1/tasks*', '*', NULL, NULL, NULL);
INSERT INTO public.casbin_rule (id, ptype, v0, v1, v2, v3, v4, v5) VALUES (11, 'p', 'employee', '/api/v1/tasks/*', '*', NULL, NULL, NULL);


--
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.permissions (id, resource, action, description, created_at, updated_at) VALUES ('6d7d855d-73a1-4652-8a17-48bdd9a8e9f0', '*', '*', 'All actions on all resources', '2026-03-18 04:08:28.178588+00', '2026-03-18 04:08:28.178588+00');
INSERT INTO public.permissions (id, resource, action, description, created_at, updated_at) VALUES ('5fd21e12-0610-4b1d-900c-6b302facc34c', '/api/v1/tasks*', '*', 'All actions on tasks', '2026-03-18 04:08:28.178588+00', '2026-03-18 04:08:28.178588+00');
INSERT INTO public.permissions (id, resource, action, description, created_at, updated_at) VALUES ('29cc5d0f-055d-4ae5-888e-605799533f5a', '/api/v1/users/profile', 'GET', 'View own profile', '2026-03-18 04:08:28.178588+00', '2026-03-18 04:08:28.178588+00');
INSERT INTO public.permissions (id, resource, action, description, created_at, updated_at) VALUES ('cfcca3ee-4c36-4331-873f-1a95145d59ed', '/api/v1/users/profile', 'PATCH', 'Update own profile', '2026-03-18 04:08:28.178588+00', '2026-03-18 04:08:28.178588+00');
INSERT INTO public.permissions (id, resource, action, description, created_at, updated_at) VALUES ('30fd4f6d-9860-4383-8b85-6c6b57dc0799', '/api/v1/users/change-password', 'PATCH', 'Change own password', '2026-03-18 04:08:28.178588+00', '2026-03-18 04:08:28.178588+00');
INSERT INTO public.permissions (id, resource, action, description, created_at, updated_at) VALUES ('db9e8cc2-e2bc-4df5-969a-72e0626b5e40', '/api/v1/users/avatar', 'POST', 'Update own avatar', '2026-03-18 04:08:28.178588+00', '2026-03-18 04:08:28.178588+00');


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.roles (id, name, description, created_at, updated_at) VALUES ('3f371147-e089-4886-b003-fdb33d3a5a0c', 'admin', 'Full system access', '2026-03-18 04:08:28.178588+00', '2026-03-18 04:08:28.178588+00');
INSERT INTO public.roles (id, name, description, created_at, updated_at) VALUES ('84f58710-547e-4c06-866e-17990fc606fe', 'employee', 'Standard user access', '2026-03-18 04:08:28.178588+00', '2026-03-18 04:08:28.178588+00');


--
-- Data for Name: role_permissions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.role_permissions (role_id, permission_id) VALUES ('3f371147-e089-4886-b003-fdb33d3a5a0c', '6d7d855d-73a1-4652-8a17-48bdd9a8e9f0');
INSERT INTO public.role_permissions (role_id, permission_id) VALUES ('84f58710-547e-4c06-866e-17990fc606fe', '5fd21e12-0610-4b1d-900c-6b302facc34c');
INSERT INTO public.role_permissions (role_id, permission_id) VALUES ('84f58710-547e-4c06-866e-17990fc606fe', '29cc5d0f-055d-4ae5-888e-605799533f5a');
INSERT INTO public.role_permissions (role_id, permission_id) VALUES ('84f58710-547e-4c06-866e-17990fc606fe', 'cfcca3ee-4c36-4331-873f-1a95145d59ed');
INSERT INTO public.role_permissions (role_id, permission_id) VALUES ('84f58710-547e-4c06-866e-17990fc606fe', '30fd4f6d-9860-4383-8b85-6c6b57dc0799');
INSERT INTO public.role_permissions (role_id, permission_id) VALUES ('84f58710-547e-4c06-866e-17990fc606fe', 'db9e8cc2-e2bc-4df5-969a-72e0626b5e40');


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('29a1732a-1fe2-4b32-825b-717d3fd9e9ed', 'cr7@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$VuYeKJKIXgvdkrpn32W+fQ$OUwOrD1f2311lZbelR5oa2+mb8iZc6tY7gqb/UoTU80', 'Cristiano', 'Ronaldo', true, NULL, '2026-03-19 06:03:29.854741+00', '2026-03-19 06:03:29.854741+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('5e5727f9-d219-47e1-a5af-85991da0a1db', 'messi@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$3Z6OZ16n3gaqA9T031O2LA$i3GwyotIbjo9gUDGxJlO0OHbMOyiZrirtNuEv+ZnTuI', 'Lionel', 'Messi', true, NULL, '2026-03-19 06:05:25.864658+00', '2026-03-19 06:05:25.864658+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('d4798812-6ca5-4191-962e-ca1adb04eb1e', 'ney@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$xCdXlkth0koPMTZgjSDJvw$ITtCEfsl81xGdk+WeAQiyvdjSeXlelGGuQJ+CndqMF4', 'Neymar', 'Junior', true, NULL, '2026-03-19 06:06:24.311895+00', '2026-03-19 06:06:24.311895+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('6e3c24e1-0e19-4b2a-b30a-135035ea2c10', 'mbappe@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$RU8eoIqczXJOrkwzz3eDJA$BYvAoHkIPd/CcBzVU21GJ7knY9O4Fk01huqraXUYwDw', 'Kylian', 'Mbappé', true, NULL, '2026-03-19 06:07:05.997523+00', '2026-03-19 06:07:05.997523+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('d3694666-aee6-4a38-86c0-c52d25b4353f', 'haaland@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$rfZhyTzdtAA9p1pnX/JhbA$NBOHfwOQRAHaEGA3UOgYY1DiWKP4SW46mkSWNBfA6Nc', 'Erling', 'Haaland', true, NULL, '2026-03-19 06:07:32.244144+00', '2026-03-19 06:07:32.244144+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('b4c0aafa-85fd-45cf-bc7f-8750336fa469', 'salah@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$ZjwhXQicUpLqQLH4pxf3sQ$TrDcdaunAmEX73CPp0/QTcVrTeToDMllK7Upjm9WVeI', 'Mohamed', 'Salah', true, NULL, '2026-03-19 06:08:33.619233+00', '2026-03-19 06:08:33.619233+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('edebf19d-479a-447b-b643-32dfa22a9a97', 'vini@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$YKuq0OcYL0PjfDb4kPoUrQ$FOxyRGBtbGfwm1xBVnEn1eonHvxYTR0MODk+lDaM/fg', 'Vinícius', 'Junior', true, NULL, '2026-03-19 06:09:29.841727+00', '2026-03-19 06:09:29.841727+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('2b659f54-9db6-4ab0-8cf4-7a8d9c962e0a', 'saka@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$PvaX+DFa931eA4ICUXgT5A$4ZMuXAEaQAwWa51eJSLn5DZ0kfgTtft+i76woXTUO50', 'Bukayo', 'Saka', true, NULL, '2026-03-19 06:10:05.268145+00', '2026-03-19 06:10:05.268145+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('357d8a9b-a417-4673-aa5d-f62c9869fae2', 'mane@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$WxFXSTSFdC7QGjYsfiGG3A$YQlSgxPtXWuB5UgOkKbOxiIg9ulHpyMxvDJ4Ge+NNaw', 'Sadio', 'Mané', true, NULL, '2026-03-19 06:10:44.846908+00', '2026-03-19 06:10:44.846908+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('09df5ec5-9acd-4256-9cd2-3929a59bc423', 'rc3@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$aPjUhD0MtPuoOlR1W4sA0A$vCCVsiRYW8c8waQFpM4fr99+1qX8eWDF6qf4Nfrsu48', 'Roberto', 'Carlos', true, NULL, '2026-03-19 06:11:52.139117+00', '2026-03-19 06:11:52.139117+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('7e481672-8755-4af8-b09d-7ba15fe8f165', 'bellingham@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$Jj/olJ78ABBEQcoCWzq+/Q$yBH6BOzEQ3xCeDQwD/eNO2ojmmyXPZPi+aGDfRHTK+M', 'Jude', 'Bellingham', true, NULL, '2026-03-19 06:12:32.046672+00', '2026-03-19 06:12:32.046672+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('5133727e-47d4-4022-ac11-023cd7e349f3', 'obruxo@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$Zwx7cgSuf/ekp8d/gJZnMQ$kGvXA7ECuE00tRY3VI6vYOUPPe+yzrF5S/xkFN5BGeg', 'Ronaldinho', 'Gaucho', true, NULL, '2026-03-19 06:13:10.488187+00', '2026-03-19 06:13:10.488187+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('b55aa8e3-d217-44ae-9883-bc59a61102ba', 'lamine@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$2KbQOD+3Pdm3vRhve+kO3w$MhcarMbM/0cQ+Jga5aNfu6/0ol5qoK/LqVALlojQuP0', 'Lamine', 'Yamal', true, NULL, '2026-03-19 06:15:12.577258+00', '2026-03-19 06:15:12.577258+00');
INSERT INTO public.users (id, email, password_hash, first_name, last_name, is_active, avatar_url, created_at, updated_at) VALUES ('163613ae-9713-4318-881b-d8c329b73d51', 'quaresma@taskify.com', '$argon2id$v=19$m=65536,t=1,p=4$JB+UZjCSvVe/nyCBY9g0Sg$DxTSd0dZO/d/BjPXX2X5emneTdY8S4RSjEKfg1JaK8Y', 'Ricardo', 'Quaresma', true, NULL, '2026-03-19 06:13:53.477939+00', '2026-03-19 06:36:47.264262+00');


--
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('7ba19da9-a92e-4420-8c65-17227419e9b5', 'Implement SSO with OAuth2', 'Configure Google and GitHub login for the internal team.', 'in_progress', 'high', false, 'b55aa8e3-d217-44ae-9883-bc59a61102ba', 'b55aa8e3-d217-44ae-9883-bc59a61102ba', '2026-03-23 00:00:00+00', NULL, '2026-03-19 06:19:21.842102+00', '2026-03-19 06:20:20.317231+00', 8.00, NULL, false);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('d6f752e4-51f5-4580-90a4-b39507c9234b', 'Migrate to Svelte 5', 'Update all frontend components to use Svelte 5 Runes.', 'pending', 'high', false, 'b55aa8e3-d217-44ae-9883-bc59a61102ba', 'b55aa8e3-d217-44ae-9883-bc59a61102ba', '2026-03-31 00:00:00+00', NULL, '2026-03-19 06:21:16.398225+00', '2026-03-19 06:21:16.398225+00', 20.00, NULL, false);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('9500deba-7595-400d-894e-4e9bc76b9e13', 'Refactor Migrations', 'Clean up migration history and ensure schema consistency.', 'pending', 'medium', true, 'b55aa8e3-d217-44ae-9883-bc59a61102ba', 'b55aa8e3-d217-44ae-9883-bc59a61102ba', '2026-03-24 00:00:00+00', NULL, '2026-03-19 06:20:13.791733+00', '2026-03-19 06:21:30.979652+00', 8.00, NULL, false);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('350f6741-720f-40f9-9cb9-1d45e392cb83', 'Configure CI/CD', 'Create GitHub Actions pipelines for linter and testing.', 'pending', 'high', false, '29a1732a-1fe2-4b32-825b-717d3fd9e9ed', '29a1732a-1fe2-4b32-825b-717d3fd9e9ed', '2026-03-25 00:00:00+00', NULL, '2026-03-19 06:23:05.443673+00', '2026-03-19 06:23:05.443673+00', 8.00, NULL, false);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('10c8ec42-9e6a-4530-8d3e-0e1e1c5b1f5c', 'Redis Caching Layer', 'Implement caching for the most accessed task data.', 'cancelled', 'medium', false, '29a1732a-1fe2-4b32-825b-717d3fd9e9ed', '29a1732a-1fe2-4b32-825b-717d3fd9e9ed', '2026-03-22 00:00:00+00', NULL, '2026-03-19 06:24:13.619133+00', '2026-03-19 06:24:51.901051+00', 6.00, NULL, false);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('8244fd73-7e40-42ea-8f9c-0ee6a6b9019d', 'Fix Worker Memory Leak', 'Investigate excessive memory usage in the notification worker.', 'completed', 'critical', false, 'b55aa8e3-d217-44ae-9883-bc59a61102ba', 'b55aa8e3-d217-44ae-9883-bc59a61102ba', '2026-03-19 00:00:00+00', '2026-03-19 06:26:16.602804+00', '2026-03-19 06:26:04.418265+00', '2026-03-19 06:26:16.603076+00', 4.00, NULL, false);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('a148827c-129f-4739-ab68-71728abeb00e', 'Modularize API Routes', 'Split Auth, Task, and User routes into separate modules.', 'in_progress', 'medium', false, 'd4798812-6ca5-4191-962e-ca1adb04eb1e', 'd4798812-6ca5-4191-962e-ca1adb04eb1e', '2026-03-26 00:00:00+00', NULL, '2026-03-19 06:29:18.359724+00', '2026-03-19 06:29:21.67597+00', 20.00, NULL, false);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('cb471912-4f0c-4e02-b7f7-0144d989f708', 'Real-time Notifications', 'Integrate WebSockets for live Kanban board updates.', 'pending', 'medium', false, 'd4798812-6ca5-4191-962e-ca1adb04eb1e', 'd4798812-6ca5-4191-962e-ca1adb04eb1e', '2026-03-28 00:00:00+00', NULL, '2026-03-19 06:30:18.288707+00', '2026-03-19 06:30:25.963978+00', 20.00, NULL, true);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('147718e6-167f-418b-9750-b33a9cefbab3', 'Security Compliance Audit', 'Run OWASP ZAP security scan on the staging environment.', 'pending', 'critical', false, 'd4798812-6ca5-4191-962e-ca1adb04eb1e', 'd4798812-6ca5-4191-962e-ca1adb04eb1e', '2026-03-23 00:00:00+00', NULL, '2026-03-19 06:31:43.903135+00', '2026-03-19 06:31:43.903135+00', 10.00, NULL, false);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('195d2bf1-3dab-40e9-86e0-3f099c335ded', 'Password Reset Flow', 'Create email sending flow and link expiration logic.', 'in_progress', 'critical', false, '5e5727f9-d219-47e1-a5af-85991da0a1db', '5e5727f9-d219-47e1-a5af-85991da0a1db', '2026-03-21 00:00:00+00', NULL, '2026-03-19 06:33:33.66935+00', '2026-03-19 06:33:48.140823+00', 8.00, NULL, false);
INSERT INTO public.tasks (id, title, description, status, priority, is_blocked, created_by, assigned_to, due_date, completed_at, created_at, updated_at, estimated_hours, actual_hours, is_archived) VALUES ('c17f9edb-4319-46b8-a4fe-4cc63e6c0509', 'RBAC Integration Tests', 'Document new RBAC and attachment endpoints in Swagger.', 'cancelled', 'medium', false, 'd4798812-6ca5-4191-962e-ca1adb04eb1e', 'd4798812-6ca5-4191-962e-ca1adb04eb1e', '2026-03-20 00:00:00+00', NULL, '2026-03-19 06:34:40.47366+00', '2026-03-19 06:34:51.909582+00', 2.00, NULL, false);


--
-- Data for Name: task_attachments; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: task_notes; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: user_roles; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.user_roles (user_id, role_id) VALUES ('29a1732a-1fe2-4b32-825b-717d3fd9e9ed', '3f371147-e089-4886-b003-fdb33d3a5a0c');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('5e5727f9-d219-47e1-a5af-85991da0a1db', '3f371147-e089-4886-b003-fdb33d3a5a0c');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('d4798812-6ca5-4191-962e-ca1adb04eb1e', '84f58710-547e-4c06-866e-17990fc606fe');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('6e3c24e1-0e19-4b2a-b30a-135035ea2c10', '84f58710-547e-4c06-866e-17990fc606fe');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('d3694666-aee6-4a38-86c0-c52d25b4353f', '84f58710-547e-4c06-866e-17990fc606fe');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('b4c0aafa-85fd-45cf-bc7f-8750336fa469', '84f58710-547e-4c06-866e-17990fc606fe');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('edebf19d-479a-447b-b643-32dfa22a9a97', '84f58710-547e-4c06-866e-17990fc606fe');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('2b659f54-9db6-4ab0-8cf4-7a8d9c962e0a', '84f58710-547e-4c06-866e-17990fc606fe');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('357d8a9b-a417-4673-aa5d-f62c9869fae2', '3f371147-e089-4886-b003-fdb33d3a5a0c');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('09df5ec5-9acd-4256-9cd2-3929a59bc423', '3f371147-e089-4886-b003-fdb33d3a5a0c');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('7e481672-8755-4af8-b09d-7ba15fe8f165', '84f58710-547e-4c06-866e-17990fc606fe');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('5133727e-47d4-4022-ac11-023cd7e349f3', '3f371147-e089-4886-b003-fdb33d3a5a0c');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('b55aa8e3-d217-44ae-9883-bc59a61102ba', '84f58710-547e-4c06-866e-17990fc606fe');
INSERT INTO public.user_roles (user_id, role_id) VALUES ('163613ae-9713-4318-881b-d8c329b73d51', '3f371147-e089-4886-b003-fdb33d3a5a0c');


--
-- Name: casbin_rule_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.casbin_rule_id_seq', 11, true);

-- 3. Re-enable constraints and triggers
SET session_replication_role = 'origin';

--
-- PostgreSQL database dump complete
--
