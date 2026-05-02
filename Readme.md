# Control Plane

## Запуск

```bash
# 1. minikube start --driver='docker'
# 2. minikube addons enable ingress
# 3. kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 8081:80
# 4. Запуск backend
# - cd app
# - Создать .env.prod .env.dev
# - make run-dev
# 5. Запуск frontend
# - cd control-plane-ui
# - npm run dev
```

## API Документация

### GET /clients

Получение списка API-клиентов.

**Query:**

| Параметр |  Тип   | Обязательный |    Ограничения    |
| -------- | ------ | ------------ | ----------------- |
| status   | string | нет          | enum ClientStatus |
| limit    | int    | нет          | 0–100             |
| offset   | int    | нет          | ≥0                |

**Response 200:**

**JSON:**
```json
{
  "items": [
    {
      "id": "uuid",
      "name": "Payment API",
      "api_service_id": "svc-1",
      "status": "running"
    }
  ],
  "total": 1
}
```

**Возможные ошибки:**
- 400 — некорректные query параметры
- 500 — внутренняя ошибка
___

### GET /clients/:client_id

Получение одного клиента.

**Path Params**
| Параметр  | Тип    | Обязательный |
| --------- | ------ | ------------ |
| client_id | string | да           |

**Response 200 OK**
**JSON**
```json
{
  "id": "uuid",
  "name": "Payment API",
  "api_service_id": "svc-1",
  "status": "running",
  "active_config_id": "cfg-1",
  "created_at": "2026-02-22T10:00:00Z"
}
```

**Возможные ошибки:**
- 400 — некорректный client_id
- 404 — client not found
- 500 — внутренняя ошибка
___

### POST /clients

Создание клиента.

**Body:**

| Параметр       | Тип    | Обязательное |
|       ---      |  ---   |      ---     |
| name           | string | да           |
| api_service_id | string | да           |
| description    | string | нет          |

**Response 201:**
**JSON:**
```json
{
  "id": "uuid",
  "name": "Payment API",
  "api_service_id": "svc-1",
  "status": "created",
  "created_at": "2026-02-22T10:00:00Z"
}
```

**Возможные ошибки:**
- 400 — некорректные данные
- 500 — внутренняя ошибка
___

### POST /clients/:client_id/restart

Асинхронный перезапуск клиента.

**Path Params**
| Параметр  | Тип    |
| --------- | ------ |
| client_id | string |


**Body:**
| Поле   | Тип    | Обязательное |
| ------ | ------ | ------------ |
| reason | string | нет          |

**Response 202:**

**JSON:**
```json
{
  "client_id": "uuid",
  "action": "restart",
  "status": "restarting"
}
```

**Возможные ошибки:**
- 404 - клиент не найден
- 409 — invalid state transition
- 500 — внутренняя ошибка
___

### POST /clients/:client_id/delete

Асинхронное удаление клиента.

**Path Params**
| Параметр  | Тип    |
| --------- | ------ |
| client_id | string |

**Response 202:**

**JSON:**
```json
{
  "client_id": "uuid",
  "action": "delete",
  "status": "deleting"
}
```

**Возможные ошибки:**
- 404 — некорректные данные
- 409 — invalid state transition
- 500 — внутренняя ошибка
___

### GET /clients/:client_id/configs

Получение конфигураций клиента.

| Параметр  | Тип    |
| --------- | ------ |
| client_id | string |

**Response 200 OK:**

**JSON:**
```json
[
  {
    "id": "cfg-1",
    "client_id": "uuid",
    "version": "v1",
    "auth_type": "api_key",
    "auth_ref": "secret_ref",
    "timeout_ms": 1000,
    "retry_count": 3,
    "retry_backoff": 500,
    "headers": {},
    "created_at": "2026-02-22T10:00:00Z",
    "created_by": "1"
  }
]
```

**Возможные ошибки:**
- 404 — client not found
- 500 — внутренняя ошибка
___

### POST /clients/:client_id/configs

Создание новой версии конфигурации.

**Path Params**
| Параметр  | Тип    |
| --------- | ------ |
| client_id | string |

**Body:**
| Поле          | Тип    | Ограничения          |
| ------------- | ------ | -------------------- |
| version       | string | required             |
| auth_type     | string | none/api_key/bearer  |
| auth_ref      | string | нет                  |
| timeout_ms    | int    | 0–60000              |
| retry_count   | int    | 0–10                 |
| retry_backoff | int    | 0–60000              |
| headers       | object | нет                  |

**Response 201 CREATE:**

**JSON:**
```json
{
  "id": "cfg-1",
  "client_id": "uuid",
  "version": "v2",
  "auth_type": "api_key",
  "timeout_ms": 1000,
  "retry_count": 3,
  "retry_backoff": 500,
  "headers": {},
  "created_at": "2026-02-22T10:00:00Z",
  "created_by": "1"
}
```

**Возможные ошибки:**
- 400 — invalid body
- 404 — клиент не найден
- 409 — config version exists
- 500 — внутренняя ошибка
___

### POST /clients/:client_id/configs/:config_id/deploy

Активация конфигурации (асинхронная операция).

**Path Params**

| Параметр  | Тип    |
| --------- | ------ |
| client_id | string |
| config_id | string |

**Response 202 Accepted:**

**JSON:**
```json
{
  "client_id": "uuid",
  "config_id": "cfg-1",
  "status": "deploying"
}
```

**Возможные ошибки:**
- 404 — client not found
- 404 — config not found
- 409 — invalid state transition
- 500 — внутренняя ошибка
___
