# Auth service

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)

Сервис аутентификации(часть)

Используемые технологии:
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Echo (веб фреймворк)
- golang-migrate/migrate (для миграций БД)
- pgx (драйвер для работы с PostgreSQL)


Сервис был написан с Clean Architecture, что позволяет легко расширять функционал сервиса и тестировать его.
Также был реализован Graceful Shutdown для корректного завершения работы сервиса


# Usage
Для запуска сервиса достаточно заполнить .env файл и прописать в терминале:
 `docker-compose up `