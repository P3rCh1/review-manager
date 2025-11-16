# Сервис назначения ревьюеров для Pull Request’ов
### Реализован в рамках тестового задания для стажировки backend avito.tech

## Допущения, комментарии к заданию
* В спецификации запроса `/pullRequest/reassign` не совпадают названия в свойствах и примере:
```
properties:
    pull_request_id: { type: string }
    old_user_id: { type: string }

example:
    pull_request_id: pr-1001
    old_reviewer_id: u2
```
**В `docs/openapi.yaml` исправлен пример*
* Возвращается createdAt несмотря на отсутствие в примерах (в спецификации компонента `PullRequest` поле присутствует)
* Pull request может открыть в том числе неактивный пользователь (требования запретить не было) 
* При смене команды или выставлении флага `is_active` = `false` продолжает быть ревьюером там, где назначен (Предполагаю, что при требовании к перераспределению необходимо было бы как-то вернуть новых ревьюеров)
* При повторяющихся id пользователей в запросе на создание команды возвращается ошибка `REPEATABLE_IDS`
* Расширено количество возможных ошибок в методах (например, запрос `/users/getReview` может вернуть ошибку `USER_NOT_FOUND`)
* В базе данных используется uuid вместо team_name в качестве первичного ключа (использование имени помешало бы, например, в будущем менять имя команды)

| Код ошибки | HTTP Status | Сообщение |
|------------|-------------|-----------|
| `TEAM_EXISTS` | 400 | team_name already exists |
| `PR_EXISTS` | 409 | PR id already exists |
| `PR_MERGED` | 409 | cannot reassign on merged PR |
| `NOT_ASSIGNED` | 409 | reviewer is not assigned to this PR |
| `NO_CANDIDATE` | 409 | no active replacement candidate in team |
| `PR_NOT_FOUND` | 404 | pull request not found |
| `TEAM_NOT_FOUND` | 404 | team not found |
| `USER_NOT_FOUND` | 404 | user not found |
| `INVALID_INPUT` | 400 | invalid request body |
| `REPEATABLE_IDS` | 400 | repeatable IDs |
| `INTERNAL_ERROR` | 500 | internal server error |

## Конфигурация
- Выполняется с помощью yaml-файла и переменных окружения
- Переменные окружения имеют приоритет
### Все параметры
| Раздел | Параметр | YAML | ENV | Значение по умолчанию |
|---------|-----------|------------|----------------|------------------------|
| **Logger** | Уровень логирования | `level` | `LOG_LEVEL` | `info` |
|  | Формат логов | `format` | `LOG_FORMAT` | `json` |
| **HTTP** | Хост | `host` | `HTTP_HOST` | `localhost` |
|  | Порт | `port` | `HTTP_PORT` | `8080` |
|  | Таймаут чтения | `read_timeout` | `HTTP_READ_TIMEOUT` | `10s` |
|  | Таймаут записи | `write_timeout` | `HTTP_WRITE_TIMEOUT` | `10s` |
|  | Idle timeout | `idle_timeout` | `HTTP_IDLE_TIMEOUT` | `60s` |
|  | Время на корректное завершение | `shutdown_timeout` | `HTTP_SHUTDOWN_TIMEOUT` | `20s` |
| **Postgres** | Хост | `host` | `POSTGRES_HOST` | `localhost` |
|  | Порт | `port` | `POSTGRES_PORT` | `5432` |
|  | Пользователь | — | `POSTGRES_USER` | — *(обязателен)* |
|  | Пароль | — | `POSTGRES_PASSWORD` | — *(обязателен)* |
|  | Имя базы данных | — | `POSTGRES_DB` | — *(обязателен)* |
|  | Режим SSL | `ssl_mode` | `POSTGRES_SSL_MODE` | `disable` |
| **Business** | Количество ревьюверов | `reviewers_count` | `REVIEWERS_COUNT` | `2` |


### Допустимые значения параметров логов  

| Параметр | YAML | ENV | Допустимые значения |
| -------------------- | ---------- | --------------- | -------------------------------- |
| Уровень логирования | level | LOG_LEVEL | debug, info, warn, error |
| Формат логов | format | LOG_FORMAT | json, text |

  

## Запуск
- Назначте обязательные переменные окружения  
Пример:
```
POSTGRES_USER=review-manager
POSTGRES_PASSWORD=you_need_to_change_this
POSTGRES_DB=review-manager
```
- Запустите с помощью docker compose
```
docker compose up -d --build
```
- Либо используйте make  

| Команда | Описание |
|---------|-------------------------------------------------------|
| `make или make up` | Собирает и запускает все сервисы в фоновом режиме|
| `make down` | Останавливает и удаляет все контейнеры|
| `make clean` | Останавливает и удаляет все контейнеры и тома |
| `make migrate-up` | Применяет миграции |
| `make migrate-down` | Выполняет откат последней миграции |
| `make migrate-version` | Выводит текущую версию миграции |
| `make migrate-force VERSION=N` | Принудительно устанавливает версию миграции `N` (в том числе при флаге dirty) |

## Дополнительные задания
### 1) Добавить простой эндпоинт статистики (например, количество назначений по пользователям и/или по PR)
| Раздел  | Параметр | Тип     | Описание |
|---------|-----------------------------------|---------|----------|
| **service** | total_users | integer | Общее количество пользователей в системе |
| | active_users | integer | Количество активных пользователей |
| | total_teams | integer | Общее количество команд |
| | merged_prs | integer | Количество замёрдженных PR |
| | open_prs | integer | Количество открытых PR |
| **user** | avg_reviews_per_active_user | float   | Среднее количество назначенных ревью на одного активного пользователя |
| | max_reviews_on_user | integer | Максимальное количество назначенных ревью на одного пользователя |
| | active_users_with_zero_reviews | integer | Количество активных пользователей, у которых нет ревью |
| **pr** | open_prs_with_0_reviewers | integer | Количество открытых PR без ревьюверов |
| | open_prs_with_1_reviewer | integer | Количество открытых PR с одним ревьювером |
| | open_prs_with_2_reviewers | integer | Количество открытых PR с двумя ревьюверами |

**Спецификация добавлена в `docs/openapi.yaml`*
  
### 2) Описать конфигурацию линтера
* Конфигурация находится в `.golangci.yaml`
* Настроены GitHub Actions для линтера

### 3) Провести нагрузочное тестирование полученного решения и приложить краткие результаты тестирования к решению  
- Инструмент: `k6`
- База данных: `~5к команд`, `~11к пользователей`, `~6к PR`  
- Виртуальных пользователей: `50`
- Задержка между отправками запроса: `0.1s` 
- Время каждого теста: `15s` 
- Результаты в формате *json* и краткие отчеты в формате *txt* находятся в `k6/results`
- Использованные скрипты в находятся в `k6/scripts`

| Эндпоинт | RPS | Ср. время ответа | p95 | Ошибки |
|----------|-----|------------------|-----|---------|
| `POST /team/add` | 371.1/s | 31.92ms | 107.77ms | 0% |
| `POST /pr/create` | 431.2/s | 14.21ms | 44.19ms | 0% |
| `POST /pullRequest/merge` | 431.7/s | 13.92ms | 44.13ms | 0% |
| `POST /users/setIsActive` | 400.0/s | 23.25ms | 109.25ms | 0% |
| `GET /team/get` | 446.7/s | 9.79ms | 38.45ms | 0% ||
| `POST /pullRequest/reassign` | 181.4/s | 35.40ms | 130.77ms | 0% |
| `GET /users/getReview` | 173.2/s | 287.70ms | 411.73ms | 0% |
| `GET /stats` | 123.5/s | 402.06ms | 722.13ms | 0% |
#### Пояснения
- Все эндпоинты прошли тесты без ошибок
- Задержка остаётся в разумных пределах, для тяжёлых эндпоинтов может достигать ~700ms
- Результаты превышающее ориентировочные SLA 300ms связаны с объемными выборками данных RPS