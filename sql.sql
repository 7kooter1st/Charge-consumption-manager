-- =============================================
-- 1. Заполнение таблицы пользователей (users)
-- =============================================
INSERT INTO users (id, login, password, is_moderator) VALUES
(1, 'admin', '$2a$10$rN6J.lIsL5N6J.lIsL5N6J.lIsL5N6J.lIsL5N6J.lIsL5N6J.lIs', true),
(2, 'user1', '$2a$10$abc1.abc1abc1abc1abc1abc1abc1abc1abc1abc1abc1abc1abc', false),
(3, 'user2', '$2a$10$def2.def2def2def2def2def2def2def2def2def2def2def', false),
(4, 'moderator1', '$2a$10$ghi3.ghi3ghi3ghi3ghi3ghi3ghi3ghi3ghi3ghi3ghi3ghi', true),
(5, 'user3', '$2a$10$jkl4.jkl4jkl4jkl4jkl4jkl4jkl4jkl4jkl4jkl4jkl4jkl', false);

-- =============================================
-- 2. Заполнение таблицы сценариев (use_cases)
-- =============================================
INSERT INTO use_cases (id, name, url, description, consumption, is_delete) VALUES
(1, 'Освещение офиса', 'https://example.com/usecases/lighting', 'LED освещение рабочих зон', 150, false),
(2, 'Кондиционирование', 'https://example.com/usecases/ac', 'Охлаждение серверной', 800, false),
(3, 'Работа ПК', 'https://example.com/usecases/pc', '50 рабочих станций', 120, false),
(4, 'Серверное оборудование', 'https://example.com/usecases/servers', 'Основные серверы', 500, false),
(5, 'Вентиляция', 'https://example.com/usecases/ventilation', 'Приточная вентиляция', 200, false);

-- =============================================
-- 3. Заполнение таблицы потребления (consumptions)
-- =============================================
INSERT INTO consumptions (id, user_id, status, total_power, created_at, updated_at, moderated_at, moderator_id) VALUES
(1, 2, 'черновик', 0, EXTRACT(EPOCH FROM NOW() - INTERVAL '5 days')::bigint, EXTRACT(EPOCH FROM NOW() - INTERVAL '5 days')::bigint, NULL, NULL),
(2, 2, 'на модерации', 470, EXTRACT(EPOCH FROM NOW() - INTERVAL '3 days')::bigint, EXTRACT(EPOCH FROM NOW() - INTERVAL '1 day')::bigint, EXTRACT(EPOCH FROM NOW() - INTERVAL '1 day')::bigint, 4),
(3, 3, 'одобрено', 1050, EXTRACT(EPOCH FROM NOW() - INTERVAL '7 days')::bigint, EXTRACT(EPOCH FROM NOW() - INTERVAL '2 days')::bigint, EXTRACT(EPOCH FROM NOW() - INTERVAL '2 days')::bigint, 4),
(4, 4, 'черновик', 0, EXTRACT(EPOCH FROM NOW() - INTERVAL '1 day')::bigint, EXTRACT(EPOCH FROM NOW() - INTERVAL '1 day')::bigint, NULL, NULL),
(5, 5, 'отклонено', 320, EXTRACT(EPOCH FROM NOW() - INTERVAL '10 days')::bigint, EXTRACT(EPOCH FROM NOW() - INTERVAL '4 days')::bigint, EXTRACT(EPOCH FROM NOW() - INTERVAL '4 days')::bigint, 1);

-- =============================================
-- 4. Заполнение таблицы связей (usecase_consumptions)
-- =============================================
INSERT INTO usecase_consumptions (id, consumption_id, use_case_id, duration) VALUES
(1, 1, 1, 8),   
(2, 1, 3, 10),  
(3, 2, 1, 8),   
(4, 2, 2, 6),   
(5, 2, 4, 24),  
(6, 3, 2, 12),  
(7, 3, 4, 24),  
(8, 3, 5, 10),  
(9, 5, 1, 4),   
(10, 5, 3, 6);  


-- После вставки с явными id нужно синхронизировать sequence, иначе новая регистрация даст duplicate key (users_pkey)
SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT COALESCE(MAX(id), 0) + 1 FROM users), false);

-- Синхронизируем счётчик для таблицы use_cases
SELECT setval(pg_get_serial_sequence('use_cases', 'id'), (SELECT COALESCE(MAX(id), 0) + 1 FROM use_cases), false);

-- Синхронизируем счётчик для таблицы consumptions
SELECT setval(pg_get_serial_sequence('consumptions', 'id'), (SELECT COALESCE(MAX(id), 0) + 1 FROM consumptions), false);

-- Синхронизируем счётчик для таблицы usecase_consumptions
SELECT setval(pg_get_serial_sequence('usecase_consumptions', 'id'), (SELECT COALESCE(MAX(id), 0) + 1 FROM usecase_consumptions), false);

