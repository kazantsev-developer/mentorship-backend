arkdown

# Платформа менторства по Go (Backend)

## Запуск через Docker Compose

1. Склонируйте репозиторий и перейдите в папку проекта.
2. Создайте файл `.env` на основе `.env.example` (или используйте переменные окружения в docker-compose.yml).
3. Выполните:

   docker-compose up -d
   Бэкенд будет доступен по адресу http://localhost:8080.

### Проверка работы

/ping – health check

Регистрация: POST /api/auth/register

Логин: POST /api/auth/login

Полный список эндпоинтов в handlers

Демо-аккаунты (создаются через регистрацию)

Студент: {"login":"student1","password":"123","roles":["student"]}

Бадди: {"login":"buddy1","password":"123","roles":["buddy"]}

Админ: {"login":"admin1","password":"admin123","roles":["admin"]}

### Тестирование

Запустите тестовый скрипт:

chmod +x test_backend.sh
./test_backend.sh
Он проверит все основные сценарии.

### Остановка

docker-compose down
Требования
Docker и Docker Compose

### Краткое описание ключевых файлов и их роли

`cmd/api/main.go` - Инициализация зависимостей, запуск Gin-сервера, регистрация маршрутов
`internal/config/config.go` - Загрузка `.env`, структура конфигурации (БД, JWT, порт)
`internal/models/` - Все GORM-модели (пользователи, блоки, материалы, прогресс, бонусы и т.д.)
`internal/repositories/` - CRUD-операции с БД (каждая таблица – свой репозиторий)
`internal/services/` - Бизнес-логика (расчёт прогресса, начисление бонусов, выдача достижений)
`internal/handlers/` - Обработка HTTP-запросов, вызов сервисов, формирование JSON-ответов
`internal/middleware/auth.go` - Проверка JWT, извлечение `userID` и ролей из токена
`pkg/db/postgres.go` - Подключение к PostgreSQL, автоматические миграции (GORM)
`docker-compose.yml` - Описание сервисов (PostgreSQL + бэкенд)
`Dockerfile` - Многостадийная сборка Go-приложения
