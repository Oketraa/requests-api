CREATE TABLE IF NOT EXISTS requests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'new',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
INSERT INTO requests (title, description, status, created_at, updated_at) VALUES
(
    'Не работает авторизация в CRM', 
    'При попытке войти под учетной записью Active Directory выдает ошибку: "500 Internal Server Error". У всей группы бухгалтерии.', 
    'new', 
    NOW() - INTERVAL '10 minutes', 
    NOW() - INTERVAL '10 minutes'
),
(
    'Запрос на установку Rider для нового стажера', 
    'Необходимо установить JetBrains Rider и Docker Desktop на рабочее место нового разработчика (кабинет 404).', 
    'new', 
    NOW() - INTERVAL '1 hour', 
    NOW() - INTERVAL '1 hour'
),
(
    'Закончился картридж в принтере', 
    'В приемной директора принтер HP LaserJet начал печатать с белыми полосами. Нужна замена картриджа.', 
    'new', 
    NOW() - INTERVAL '2 hours', 
    NOW() - INTERVAL '2 hours'
),
(
    'Падение производительности БД в пиковые часы', 
    'Каждый день с 11:00 до 12:30 резко возрастает время ответа API. Подозрение на долгие блокировки в таблице транзакций.', 
    'in_progress', 
    NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '2 hours'
),
(
    'Настройка VPN для удаленных сотрудников', 
    'Заявка от отдела маркетинга. Нужно выдать доступы к корпоративному OpenVPN трем новым сотрудникам на удаленке.', 
    'in_progress', 
    NOW() - INTERVAL '5 hours', 
    NOW() - INTERVAL '30 minutes'
),
(
    'Ошибка 403 при выгрузке отчетов в PDF', 
    'Аналитики не могли скачать отчет за прошлый месяц. Проблема исправлена изменением прав на директорию хранения на сервере.', 
    'resolved', 
    NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '1 day'
),
(
    'Замена клавиатуры', 
    'Пролила кофе на клавиатуру, залипают клавиши WASD. Оформили замену со склада.', 
    'closed', 
    NOW() - INTERVAL '7 days', 
    NOW() - INTERVAL '6 days'
),
(
    'Обновление SSL-сертификата на тестовом стенде', 
    'Истек срок действия сертификата Let''s Encrypt на staging. Сертификат успешно перевыпущен через certbot.', 
    'closed', 
    NOW() - INTERVAL '10 days', 
    NOW() - INTERVAL '10 days'
);