# Organizational Structure API

REST API для управления организационной структурой компании (подразделения и сотрудники).

## Технологии

- Go 1.21+
- PostgreSQL 15
- GORM
- Gorilla Mux
- Docker / Docker Toolbox      (Использовал это потому что обычный докер не могу использовать из за wsl2 и гипервизора)
- Goose (миграции)

## Запуск проекта

### Предварительные требования

- Docker Desktop или Docker Toolbox
- Go 1.21+

### 1. Запуск PostgreSQL

```bash
docker-compose up -d