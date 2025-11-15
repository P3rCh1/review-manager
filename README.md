# Сервис назначения ревьюеров для Pull Request’ов
- Реализован в рамках тестового задания для стажировки backend avito.tech. 

## Допущения, комментарии к заданию
- В спецификации запроса `/pullRequest/reassign` не совпадают названия в свойствах и примере (я использовал old_user_id):
```
properties:
    pull_request_id: { type: string }
    old_user_id: { type: string }

example:
    pull_request_id: pr-1001
    old_reviewer_id: u2
```
- Возвращается createdAt несмотря на отсутствие в примерах (в спецификации компонента `PullRequest` поле присутствует)
- Pull request может открыть в том числе неактивный пользователь (требования запретить не было)
- При смене команды или выставлении флага `is_active` = `false` продолжает быть ревьюером там, где назначен.
(Предполагаю, что при требовании к перераспределению необходимо было как-то вернуть новых ревьюеров)
- При повторяющихся id пользователей в запросе на создание команды возвращается ошибка `repeatable ID`
- Расширено количество возможных ошибок в методах (например, запрос `/users/getReview` может вернуть ошибку `USER_NOT_FOUND`)

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

| Параметр            | YAML ключ | ENV переменная | Допустимые значения              |
| -------------------- | ---------- | --------------- | -------------------------------- |
| Уровень логирования | level      | LOG_LEVEL       | debug, info, warn, error         |
| Формат логов        | format     | LOG_FORMAT      | json, text                       |

  

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
| `make migrate-force VERSION=N` | **Принудительно** устанавливает версию миграции `N` (в том числе при флаге dirty) |
