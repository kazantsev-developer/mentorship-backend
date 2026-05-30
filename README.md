# Платформа менторства по Go (Backend)

Backend-часть платформы для студентов, наставников и администраторов. Реализована на Go с использованием Gin, GORM и PostgreSQL.

## Запуск через Docker Compose

## Склонируйте репозиторий и перейдите в папку проекта:

git clone https://github.com/kazantsev-developer/mentorship-backend.git
cd mentorship-backend
Скопируйте файл с переменными окружения:

cp .env.example .env
При необходимости отредактируйте .env (пароли, JWT-секрет, настройки БД).

## Запустите контейнеры:

docker-compose up -d
Бэкенд будет доступен по адресу http://localhost:8080.

## Остановка

docker-compose down

Health check: GET /ping {"message":"pong"}

## Регистрация и логин

### Регистрация

curl -X POST http://localhost:8080/api/auth/register \
 -H "Content-Type: application/json" \
 -d '{"login":"test","password":"123","display_name":"Test","roles":["student"]}'

### Логин

curl -X POST http://localhost:8080/api/auth/login \
 -H "Content-Type: application/json" \
 -d '{"login":"test","password":"123"}'

## Демо-аккаунты уже созданы в базе

Роль Логин Пароль
Студент test_student 123
Бадди test_buddy 123
Админ admin 123

## Технологии

Go 1.25 (Alpine)
Gin – HTTP-роутер
GORM – ORM для PostgreSQL
JWT – аутентификация (golang-jwt)
Bcrypt – хэширование паролей
Docker Compose – контейнеризация

## Основные эндпоинты

## Публичные

POST /api/auth/register Регистрация
POST /api/auth/login Вход (JWT-токен)

## Защищённые (требуют токен)

## Профиль

GET /api/user/profile – текущий пользователь
PUT /api/user/profile – обновление профиля
GET /api/user/:id/profile – публичный профиль

## Roadmap и прогресс

GET /api/roadmap – блоки и материалы с прогрессом студента
POST /api/materials/view – отметка материала просмотренным

## Бонусы и достижения

GET /api/bonus/balance – баланс бонусов
GET /api/bonus/history – история операций
POST /api/bonus/convert – конвертация 100 бонусов → 1% скидки
GET /api/achievements – список достижений с прогрессом
Buddy (наставник)
GET /api/my-students – список закреплённых студентов
POST /api/blocks/approve – подтверждение блока
GET /api/buddy/students/:id – данные студента
GET /api/buddy/students/:id/roadmap – прогресс по блокам студента
GET /api/buddy/students/:id/activity – история активности студента

## Собеседования

POST /api/interviews/real – добавить real‑собеседование
POST /api/interviews/mock – добавить mock‑собеседование
GET /api/interviews/my – список собеседований студента
GET /api/interviews/real – общий каталог real‑собеседований

## Календарь

POST /api/calendar/events – создать событие (buddy)
GET /api/calendar/events – события для текущего пользователя
GET /api/calendar/upcoming – ближайшие 7 дней

## Заявки 1x1 (для студента)

POST /api/one-on-one – создать заявку
GET /api/one-on-one – список заявок студента
POST /api/one-on-one/approve – одобрить (только админ)
POST /api/one-on-one/reject – отклонить (только админ)

## Финальные проверки

POST /api/final-checks/schedule – назначить (buddy)
POST /api/final-checks/complete – завершить с результатом
GET /api/final-checks/student/:student_id – список проверок студента

## Администрирование (требуют роль admin)

## Пользователи (расширенное управление)

GET /api/admin/users – список пользователей
POST /api/admin/users – создание пользователя
GET /api/admin/users/:user_id – детальная информация
PUT /api/admin/users/:user_id – редактирование
DELETE /api/admin/users/:user_id – soft delete
POST /api/admin/users/:user_id/change-password – смена пароля
GET /api/admin/users/:user_id/progress – прогресс студента по блокам
POST /api/admin/users/:user_id/approve-block/:block_id – подтвердить блок (админ)

### Назначение Buddy

POST /api/admin/assign-buddy – привязать бадди к студенту

### Roadmap блоки

GET /api/admin/blocks – список блоков
POST /api/admin/blocks – создать блок
PUT /api/admin/blocks/:id – обновить
DELETE /api/admin/blocks/:id – удалить (soft delete)

### Материалы

GET /api/admin/materials – список материалов (фильтр по block_id)
POST /api/admin/materials – создать материал
PUT /api/admin/materials/:id – обновить
DELETE /api/admin/materials/:id – удалить (soft delete)
PATCH /api/admin/materials/:id/status – включить/отключить (is_active)

### Достижения

GET /api/admin/achievements – все достижения
POST /api/admin/achievements – создать
PUT /api/admin/achievements/:id – обновить
DELETE /api/admin/achievements/:id – удалить
PATCH /api/admin/achievements/:id/status – включить/отключить
GET /api/admin/achievements/:id/users – список пользователей, получивших достижение

### Заявки 1x1 (админская панель)

GET /api/admin/one-on-one – все заявки (с именем студента и балансом бонусов)
POST /api/admin/one-on-one/:id/approve – одобрить (списывает 1000 бонусов)
POST /api/admin/one-on-one/:id/reject – отклонить
POST /api/admin/one-on-one/:id/complete – отметить завершённой

### Статистика

GET /api/admin/stats – общие метрики: количество пользователей, студентов, бадди, активных заявок 1x1, выданных достижений

## Модели данных (основные)

users – логин, пароль, отображаемое имя, telegram, дата начала обучения
user_roles – связь пользователя с ролями (student, buddy, admin)
student_buddy_assignments – назначение бадди студенту
blocks – блоки roadmap (название, описание, порядок)
materials – материалы (тип, content_type, URL, обязательность)
material_progresses – отметка материалов студентами
block_progresses – статус блока для студента
achievements – достижения (название, бонус, условие)
user_achievements – выданные достижения
bonus_transactions – история бонусов
interviews – собеседования (mock/real, компания, позиция, фидбэк)
calendar_events – события календаря
one_on_one_requests – заявки на 1x1
final_checks – финальные проверки (техничка, прожарка)
activity_events – лог активности пользователя

## CORS

Разрешённые источники задаются переменной ALLOWED_ORIGINS(например, http://localhost:3000,http://185.75.189.130:3000).
