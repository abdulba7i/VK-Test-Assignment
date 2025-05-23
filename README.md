
# 🎬 Фильмотека

Тестовое задание от **VK** на позицию **Junior Golang-разработчика**.

## 📌 Задача

Разработать бэкенд приложения **"Фильмотека"**, предоставляющий **REST API** для управления базой данных фильмов и актёров.

---

## 🔧 Функциональность API

### 🎭 Работа с актёрами

* Добавление информации об актёре (имя, пол, дата рождения)
* Редактирование информации об актёре (полное и частичное)
* Удаление информации об актёре
  
### 🎬 Работа с фильмами

* Добавление информации о фильме:

  * Название (1–150 символов)
  * Описание (до 1000 символов)
  * Дата выпуска
  * Рейтинг (0–10)
  * Список актёров
* Редактирование информации о фильме (полное и частичное)
* Удаление фильма

### 🔍 Поиск и фильтрация

* Получение списка фильмов с сортировкой:

  * По названию
  * По рейтингу (по умолчанию — по убыванию)
  * По дате выпуска
* Поиск по фрагменту названия фильма или имение актёра
* Получение списка всех актёров с участием их фильмов

### 👤 Работа с пользователями и ролями

* Авторизация обязательна
* Поддержка двух ролей:

  * **Обычный пользователь** — только чтение и поиск
  * **Администратор** — полный доступ
* Распределение ролей задаётся вручную (например, напрямую через БД)

---

## 🛠 Технологии и требования

| Компонент             | Использовано                        |
| --------------------- | ----------------------------------- |
| Язык программирования | Go                                  |
| HTTP-сервер           | `net/http` (стандартная библиотека) |
| База данных           | PostgreSQL                          |
| Архитектура API       | OpenAPI 3.0                         |
| Авторизация           | JWT                  |
| Логирование           | Включает базовые запросы и ошибки   |
| Покрытие тестами      | ≥ 70%                               |
| Сборка                | Docker                              |
| Окружение             | Docker Compose                      |

---

## 🚀 Запуск проекта

### 🔧 Требования

* Docker
* Docker Compose
* Env
### Добавьте env файл
```
DB_HOST=vk-test-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres

SECRET_KEY=5a0f8e0a3f6c4d8992fbd3b25e509987
```
### ⚙️ Команды

```bash
# Сборка и запуск
docker compose up --build

# Приложение будет доступно на:
# http://localhost:8080
```
---
## 🧪 Тестирование

```
go test ./... -cover
```
---
## 📄 Документация API

Автоматически генерируется в формате OpenAPI 3.0
Доступна по адресу: http://localhost:8080/swagger/index.html