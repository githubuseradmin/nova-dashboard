# nova

[![CI](https://github.com/githubuseradmin/nova-dashboard/actions/workflows/ci.yml/badge.svg)](https://github.com/githubuseradmin/nova-dashboard/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Go](https://img.shields.io/badge/Go-1.23%2B-00ADD8)
![Svelte](https://img.shields.io/badge/Svelte-5-FF3E00)

Чистый современный full-stack стартер — **Go · Svelte · PostgreSQL**. Лендинг, сессионная авторизация, пользовательский дашборд и админ-панель — со своим, вручную написанным UI и без «фреймворк-мусора».

> 🇬🇧 English version: [README.md](README.md)

## Почему такой стек

| Слой | Выбор | Почему |
|---|---|---|
| **Backend** | Go — stdlib `net/http` + [`chi`](https://github.com/go-chi/chi) | Нативные хендлеры `net/http`, без `fasthttp`. `chi` — только роутинг и middleware. |
| **База** | PostgreSQL + [`pgx`](https://github.com/jackc/pgx) | Явный SQL и миграции, без магии ORM. |
| **Frontend** | Svelte 5 + TypeScript + Vite | Бандл ~17 КБ gzip, компилируется, без virtual DOM. |
| **Стили** | Свой CSS на дизайн-токенах | Без Tailwind — своя тема light/dark на CSS-переменных. |
| **Auth** | Сессия в куке (`HttpOnly`/`Secure`/`SameSite`) + CSRF + argon2id | Безопаснее JWT в `localStorage`. |
| **Edge** | Nginx + Certbot | TLS, отдача статики и reverse-proxy на `/api`. |

## Архитектура

```
nova/
├── cmd/server/          # точка входа: конфиг, логгер, HTTP-сервер, graceful shutdown
├── internal/
│   ├── config/          # конфигурация из окружения
│   ├── http/            # транспорт: роутер (chi), middleware, хендлеры, ответы
│   ├── auth/            # сессии, хеширование паролей (argon2id), CSRF
│   ├── store/           # доступ к данным (pgx), запуск миграций
│   └── models/          # доменные типы
├── migrations/          # SQL-миграции (встроены через embed)
└── web/
    ├── site/            # статический лендинг (свой CSS)
    └── app/             # SPA на Svelte 5 + Vite: login, dashboard, admin
```

Зависимости направлены в одну сторону: `http → auth / store → models`. Транспорт не ходит в SQL напрямую — только через `store`.

## API

| Метод | Путь | Защита | Описание |
|---|---|---|---|
| `GET` | `/healthz` | — | Liveness-проба |
| `GET` | `/api/csrf` | — | Выдать CSRF-токен (кука + тело) |
| `POST` | `/api/auth/login` | CSRF | Проверить учётные данные, открыть сессию |
| `POST` | `/api/auth/logout` | CSRF | Закрыть сессию |
| `GET` | `/api/auth/me` | сессия | Текущий пользователь |
| `GET` | `/api/dashboard` | сессия | Данные дашборда |
| `GET` | `/api/admin/users` | admin | Список пользователей |
| `PUT` | `/api/admin/users/{id}/role` | admin + CSRF | Сменить роль пользователя |

## Быстрый старт (разработка)

Требования: **Go 1.23+**, **Node 20+**, **Docker**.

```sh
cp .env.example .env          # при необходимости отредактируй

make db-up                    # поднять Postgres (Docker)
make dev                      # API на http://localhost:8080 (миграции + сид админа)
make web-install              # один раз — установить зависимости фронта
make web-dev                  # dev-сервер Svelte на http://localhost:5173 (проксирует /api)
```

Открой http://localhost:5173 и войди под сид-админом
(`NOVA_ADMIN_EMAIL` / `NOVA_ADMIN_PASSWORD`, значения по умолчанию в `.env.example`).

## Тесты

```sh
make test                     # юнит-тесты идут всегда
# интеграционные тесты store запускаются только при заданной тестовой БД:
NOVA_TEST_DATABASE_URL=postgres://nova:nova@localhost:5432/nova?sslmode=disable make test
```

CI (GitHub Actions) на каждый push гоняет `go vet`, `go test -race` против сервиса
Postgres и `svelte-check` + прод-сборку фронта.

## Продакшн

```sh
make web-build                # собирает web/app/dist с base /app/
make build                    # собирает ./bin/server
```

Один Go-бинарь отдаёт лендинг по `/`, SPA по `/app/` и JSON API по `/api`.
Спереди — Nginx + Certbot для TLS.

## Лицензия

[MIT](LICENSE).
