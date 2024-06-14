# API для калькулятора валют

## Объем задач
- CRUD для операций с валютой
- Фоновый воркер для получения курсов с FastForex
- Хранение курсов в памяти приложения
- Хранение валют в PostgreSQL
- Написаны unit тесты с моками зависимостей через mockery
- Функциональные тесты для проверки БД
- Предусмотрена валидация входящий данных
- К проекту подключен golangci-lint
- Документация Swagger api/swagger/swagger.yaml

## Запуск
- скопировать `.env.example` в `.env`, добавить API ключ FAST_FOREX_API_KEY=
- `make first-init` инициализация зависимостей проекта (docker, миграции)
- `make run` запуск приложения на 8080 порту
- `make test` запуск unit тестов
- `make local-test` запуск тестов, включая функциональные
- `make lint` запуск линтера
- `make swag` перегенерерация документации swagger

## Примечания к решению
- :trident: чистая архитектура (handler->service->repository)
- :book: Стандартная схема проекта GO
- :cd: docker compose + Makefile
- :card_file_box: миграции для PostgreSQL
- :heavy_check_mark: коллекции для Postman (examples/postman)