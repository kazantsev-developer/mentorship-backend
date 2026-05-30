# Платформа менторства по Go (Backend)

Backend-часть платформы для студентов, наставников и администраторов. Реализована на Go с использованием Gin, GORM и PostgreSQL.

## Запуск через Docker Compose

1. Склонируйте репозиторий и перейдите в папку проекта:

```bash
git clone https://github.com/kazantsev-developer/mentorship-backend.git
cd mentorship-backend
Скопируйте файл с переменными окружения:

cp .env.example .env
При необходимости отредактируйте .env (пароли, JWT-секрет, настройки БД).

Запустите контейнеры:
docker-compose up -d
Бэкенд будет доступен по адресу http://localhost:8080.

Остановка
docker-compose down

Health check: GET /ping {"message":"pong"}

Регистрация и логин

# Регистрация
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"test","password":"123","display_name":"Test","roles":["student"]}'

# Логин
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"test","password":"123"}'

Демо-аккаунты уже созданы в базе
Роль	Логин	Пароль
Студент	test_student	123
Бадди	test_buddy	123
Админ	admin	123


Технологии
Go 1.25 (Alpine)
Gin – HTTP-роутер
GORM – ORM для PostgreSQL
JWT – аутентификация (golang-jwt)
Bcrypt – хэширование паролей
Docker Compose – контейнеризация

Основные эндпоинты
Публичные
POST	/api/auth/register	Регистрация
POST	/api/auth/login	Вход (JWT-токен)
Защищённые (требуют токен)

Профиль

GET	/api/user/profile	Текущий пользователь
PUT	/api/user/profile	Обновление профиля
GET	/api/user/:id/profile	Публичный профиль
Roadmap и прогресс

GET	/api/roadmap	Блоки и материалы с прогрессом
POST	/api/materials/view	Отметка материала пройденным
Бонусы и достижения

GET	/api/bonus/balance	Баланс бонусов
GET	/api/bonus/history	История операций
POST	/api/bonus/convert	Конвертация 100 бонусов → 1%
GET	/api/achievements	Список достижений с прогрессом
Buddy (наставник)

GET	/api/my-students	Студенты текущего бадди
POST	/api/blocks/approve	Подтверждение блока
GET	/api/buddy/students/:id	Данные студента
GET	/api/buddy/students/:id/roadmap	Прогресс по блокам студента
GET	/api/buddy/students/:id/activity	История активности студента

Собеседования

POST	/api/interviews/real	Добавление real-собеседования
POST	/api/interviews/mock	Добавление mock-собеседования
GET	/api/interviews/my	Список собеседований студента
GET	/api/interviews/real	Общий каталог real

Администрирование (требуют роль admin)

GET	/api/admin/users	Список пользователей
POST	/api/admin/users	Создание пользователя
PUT	/api/admin/users/:id	Редактирование пользователя
DELETE	/api/admin/users/:id	Удаление (soft delete)
POST	/api/admin/assign-buddy	Назначение бадди студенту
GET	/api/admin/blocks	Список блоков
POST	/api/admin/blocks	Создание блока
PUT	/api/admin/blocks/:id	Обновление блока
DELETE	/api/admin/blocks/:id	Удаление блока

Модели данных (основные)
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

CORS
Разрешённые источники задаются переменной ALLOWED_ORIGINS(например, http://localhost:3000,http://185.75.189.130:3000).

```
