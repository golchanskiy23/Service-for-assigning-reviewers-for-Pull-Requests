# Примеры использования API

Данный документ содержит примеры запросов и ответов для всех эндпоинтов сервиса.

**Базовый URL**: `http://localhost:8080`

---

## 1. POST /team/add - Создание команды с участниками

Создает команду и обновляет/создает пользователей.

### Запрос:
```bash
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": true
      },
      {
        "user_id": "u2",
        "username": "Bob",
        "is_active": true
      },
      {
        "user_id": "u3",
        "username": "Charlie",
        "is_active": true
      }
    ]
  }'
```

### Успешный ответ (201):
```json
{
  "team": {
    "team_name": "backend",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": true
      },
      {
        "user_id": "u2",
        "username": "Bob",
        "is_active": true
      },
      {
        "user_id": "u3",
        "username": "Charlie",
        "is_active": true
      }
    ]
  }
}
```

### Ошибка - команда уже существует (400):
```json
{
  "error": {
    "code": "TEAM_EXISTS",
    "message": "team_name already exists"
  }
}
```

---

## 2. GET /team/get - Получение команды

Получает информацию о команде и её участниках.

### Запрос:
```bash
curl -X GET "http://localhost:8080/team/get?team_name=backend"
```

### Успешный ответ (200):
```json
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "u1",
      "username": "Alice",
      "is_active": true
    },
    {
      "user_id": "u2",
      "username": "Bob",
      "is_active": true
    },
    {
      "user_id": "u3",
      "username": "Charlie",
      "is_active": true
    }
  ]
}
```

### Ошибка - команда не найдена (404):
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "team not found"
  }
}
```

---

## 3. POST /users/setIsActive - Установка флага активности пользователя

Устанавливает флаг активности пользователя.

### Запрос:
```bash
curl -X POST http://localhost:8080/users/setIsActive \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "u2",
    "is_active": false
  }'
```

### Успешный ответ (200):
```json
{
  "user": {
    "user_id": "u2",
    "username": "Bob",
    "team_name": "backend",
    "is_active": false
  }
}
```

### Активация пользователя:
```bash
curl -X POST http://localhost:8080/users/setIsActive \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "u2",
    "is_active": true
  }'
```

### Ошибка - пользователь не найден (404):
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "user not found"
  }
}
```

---

## 4. GET /users/getReview - Получение PR'ов пользователя

Получает список PR'ов, где пользователь назначен ревьювером.

### Запрос:
```bash
curl -X GET "http://localhost:8080/users/getReview?user_id=u2"
```

### Успешный ответ (200):
```json
{
  "user_id": "u2",
  "pull_requests": [
    {
      "pull_request_id": "pr-1001",
      "pull_request_name": "Add search",
      "author_id": "u1",
      "status": "OPEN"
    },
    {
      "pull_request_id": "pr-1002",
      "pull_request_name": "Fix bug",
      "author_id": "u3",
      "status": "MERGED"
    }
  ]
}
```

### Пустой список (200):
```json
{
  "user_id": "u2",
  "pull_requests": []
}
```

### Ошибка - пользователь не найден (404):
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "user not found"
  }
}
```

---

## 5. POST /pullRequest/create - Создание PR

Создает PR и автоматически назначает до 2 ревьюверов из команды автора.

### Запрос:
```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search functionality",
    "author_id": "u1"
  }'
```

### Успешный ответ (201):
```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search functionality",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": ["u2", "u3"],
    "created_at": "2025-10-24T12:34:56Z",
    "merged_at": null
  }
}
```

### Пример с одним ревьювером (если в команде только 2 человека):
```json
{
  "pr": {
    "pull_request_id": "pr-1002",
    "pull_request_name": "Fix bug",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": ["u2"],
    "created_at": "2025-10-24T12:35:00Z",
    "merged_at": null
  }
}
```

### Пример без ревьюверов (если в команде только автор):
```json
{
  "pr": {
    "pull_request_id": "pr-1003",
    "pull_request_name": "Update docs",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": [],
    "created_at": "2025-10-24T12:36:00Z",
    "merged_at": null
  }
}
```

### Ошибка - PR уже существует (409):
```json
{
  "error": {
    "code": "PR_EXISTS",
    "message": "PR id already exists"
  }
}
```

### Ошибка - автор/команда не найдены (404):
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "author/team not found"
  }
}
```

---

## 6. POST /pullRequest/merge - Merge PR

Помечает PR как MERGED (идемпотентная операция).

### Запрос:
```bash
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001"
  }'
```

### Успешный ответ (200):
```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search functionality",
    "author_id": "u1",
    "status": "MERGED",
    "assigned_reviewers": ["u2", "u3"],
    "created_at": "2025-10-24T12:34:56Z",
    "merged_at": "2025-10-24T13:00:00Z"
  }
}
```

### Повторный вызов (идемпотентность) - тот же ответ (200):
```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search functionality",
    "author_id": "u1",
    "status": "MERGED",
    "assigned_reviewers": ["u2", "u3"],
    "created_at": "2025-10-24T12:34:56Z",
    "merged_at": "2025-10-24T13:00:00Z"
  }
}
```

### Ошибка - PR не найден (404):
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "PR not found"
  }
}
```

---

## 7. POST /pullRequest/reassign - Переназначение ревьювера

Переназначает конкретного ревьювера на другого из его команды.

### Запрос:
```bash
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "old_user_id": "u2"
  }'
```

### Успешный ответ (200):
```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search functionality",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": ["u3", "u4"],
    "created_at": "2025-10-24T12:34:56Z",
    "merged_at": null
  },
  "replaced_by": "u4"
}
```

### Ошибка - PR не найден (404):
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "PR or user not found"
  }
}
```

### Ошибка - PR уже MERGED (409):
```json
{
  "error": {
    "code": "PR_MERGED",
    "message": "cannot reassign on merged PR"
  }
}
```

### Ошибка - ревьювер не назначен (409):
```json
{
  "error": {
    "code": "NOT_ASSIGNED",
    "message": "reviewer is not assigned to this PR"
  }
}
```

### Ошибка - нет доступных кандидатов (409):
```json
{
  "error": {
    "code": "NO_CANDIDATE",
    "message": "no active replacement candidate in team"
  }
}
```

---

## Полный сценарий тестирования

### Шаг 1: Создать команду
```bash
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend",
    "members": [
      {"user_id": "u1", "username": "Alice", "is_active": true},
      {"user_id": "u2", "username": "Bob", "is_active": true},
      {"user_id": "u3", "username": "Charlie", "is_active": true},
      {"user_id": "u4", "username": "David", "is_active": true}
    ]
  }'
```

### Шаг 2: Получить команду
```bash
curl -X GET "http://localhost:8080/team/get?team_name=backend"
```

### Шаг 3: Создать PR (автоматически назначатся ревьюверы)
```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1"
  }'
```

### Шаг 4: Получить PR'ы пользователя u2
```bash
curl -X GET "http://localhost:8080/users/getReview?user_id=u2"
```

### Шаг 5: Переназначить ревьювера
```bash
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "old_user_id": "u2"
  }'
```

### Шаг 6: Деактивировать пользователя
```bash
curl -X POST http://localhost:8080/users/setIsActive \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "u3",
    "is_active": false
  }'
```

### Шаг 7: Merge PR
```bash
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001"
  }'
```

### Шаг 8: Попытка переназначить ревьювера после merge (должна вернуть ошибку)
```bash
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "old_user_id": "u4"
  }'
```

---

## Примечания

1. **Автоматическое назначение ревьюверов**: При создании PR автоматически назначаются до 2 активных ревьюверов из команды автора, исключая самого автора.

2. **Переназначение**: При переназначении новый ревьювер выбирается из команды **заменяемого** ревьювера, а не автора.

3. **Идемпотентность merge**: Повторный вызов merge для уже объединенного PR не приводит к ошибке и возвращает актуальное состояние.

4. **Неактивные пользователи**: Пользователи с `is_active: false` не назначаются на ревью и не могут быть выбраны при переназначении.

5. **После merge**: После того как PR помечен как MERGED, изменение списка ревьюверов запрещено.

