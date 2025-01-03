# user-email-verification

test task

## Описание проекта

Данный проект представляет собой бэкенд-сервис для регистрации пользователей по e-mail с последующим подтверждением адреса e-mail. Сервис демонстрирует:

## Структура проекта:

код разбит на слои:

- handler — слой, отвечающий за приём HTTP-запросов, парсинг входных данных, вызов бизнес-логики и формирование ответов.
- service — слой бизнес-логики. Здесь реализованы правила регистрации, верификации, проверки токенов и обновления статусов.
- repository — слой доступа к данным в PostgreSQL. Операции Create/Read/Update для пользователей и токенов.
- models — структуры данных, описывающие модели в БД.
- db — инициализация и подключение к базе PostgreSQL.

Это упрощает сопровождение и расширение кода, поскольку позволяет разделить ответственность между слоями.

## Архитектура и расширяемость:

Отдельная таблица для токенов верификации (email_verification_tokens). Это сделано для того, чтобы:

- Было легко добавлять новые типы токенов (например, для сброса пароля или подтверждения изменения e-mail).
- Возможна гибкая логика истечения срока действия токена, отметки о его использовании.
- Хранить историю верификаций, если потребуется.

В более простой реализации можно было бы хранить токен прямо в таблице users, но отдельная таблица повышает масштабируемость и гибкость.

## Подтверждение e-mail:

При регистрации создаётся пользователь с email_verified = false.
Генерируется токен в отдельной таблице и возвращается в ответе (а в реальной ситуации отправляется по e-mail).
При переходе по ссылке с токеном пользователь считается подтверждённым (email_verified = true), а токен помечается как использованный.

#### Использован фреймворк Gin для удобного роутинга и обработки JSON.

## Скрипты миграций:

В директории scripts находятся SQL-скрипты, создающие таблицы и индексы. Это пример — его можно улучшать, использовать внешние пакеты для миграций (например, golang-migrate).
Применить миграции:

psql "postgres://test-user:password@localhost:5432/test_db?sslmode=disable" -f scripts/init_db.sql

## Тестирование:

Сервис можно протестировать с помощью Postman или curl.

**Пример запросов к /users/register и /users/verify:**

POST http://localhost:8080/users/register
Content-Type: application/json
{
"nick": "johndoe",
"name": "John Doe",
"email": "john@example.com"
}
Ответ будет содержать verification_link. Откройте его в браузере или сделайте GET запрос.

Подтверждение email:
GET http://localhost:8080/users/verify?token=<полученный токен>
Ожидаемый ответ:
{
"status": "verified",
"email": "john@example.com"
}

### Пароль и аутентификация:

В задании не упомянут пароль. Возможно, целью тестового могла быть именно проверка процесса подтверждения e-mail, а не аутентификация.
При необходимости, поле для пароля и логика хранения, и проверки пароля (с использованием bcrypt или других алгоритмов хэширования) могут быть добавлены. Так же можно добавить полноценную аутентификацию по паролю или по безпарольным токенам.
В рамках данного задания мы указали только минимально необходимые поля (nick, name, email) и логику верификации e-mail. В реальном проекте добавление пароля — стандартная задача и делается просто расширением модели users и логики аутентификации.

### Другие детали:

**_Пул подключений:_**
\*sql.DB — это не прямое соединение, а пул соединений. Мы не пишем отдельный пул, так как он встроен в стандартную библиотеку. При необходимости можно настроить параметры пула (SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime) в NewPostgresDB, если нужны оптимизации.

**_Генерация UUID:_**
Мы используем github.com/google/uuid для генерации уникальных идентификаторов. Это удобно для первичных ключей в БД и токенов подтверждения.

**_Ошибка и валидация:_**
В рамках тестового задания мы не делаем сложную валидацию полей (email, name). В реальном проекте можно добавить проверки, валидацию формы данных, нормализацию email и т.д.

**_Логи и конфигурация:_**
В этом примере логирование минимальное, конфигурация берётся из кода или переменных окружения. В реальной ситуации мы бы добавили более продвинутую систему логирования, конфигурационные файлы, поддержку разных окружений (dev, test, prod).

### Таким образом, это решение:

- Демонстрирует чистую архитектуру и разбиение на логические слои.
- Показывает понимание реляционной БД и проектирование таблиц.
- Демонстрирует умение работать с Gin, JSON, REST API.
- Учитывает масштабируемость и гибкость.
- Готово к расширению функционала (введение пароля, другие токены, дополнительная бизнес-логика).
