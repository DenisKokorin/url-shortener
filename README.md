## url-shortener

### 1. Создайте конфигурационный файл

Создайте файл `config.yaml` в папке config проекта:

```yaml
env: "local" # local | dev | prod
storage: "postgres" # postgres | memory
alias_length: 10
http_server:
  address: "0.0.0.0:8080"
  timeout: 4s
  idle_timeout: 30s
```

### 2. Создайте файл с переменными окружения (если используете Postgres в качестве хранилища)

Создайте файл `.env` в папке env проекта:

```env
DB_USER=user
DB_PASSWORD=password
DB_HOST=postgres
DB_PORT=5432
DB_NAME=db
```

### 3. Запустите сервер

если используете Postgres:

```bash
docker compose --env-file="env/.env" up
```
если используете внутреннюю память:

```bash
docker compose up
```
