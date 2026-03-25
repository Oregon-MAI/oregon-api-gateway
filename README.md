# Oregon API Gateway
Gateway поднимается на `8000` порту.

## Как запустить

1. Go run
```sh
go run cmd/api-gateway/main.go
```
2. Docker
```sh
docker build -t "api-gw:latest" .
docker run api-gw
```

## Эндпоинты ресурсов

### 1. Создание ресурса
**POST** `/api/v1/resources`

Переговорка:
```json
{
  "name": "Meeting Room Alpha",
  "type": "RESOURCE_TYPE_MEETING_ROOM",
  "location": "Office 1, Floor 3",
  "details": {
    "capacity": 12,
    "has_projector": true,
    "has_whiteboard": true
  }
}
```

Рабочее место:
```json
{
  "name": "Workspace 304",
  "type": "RESOURCE_TYPE_WORKSPACE",
  "location": "Office 1, Openspace 2",
  "details": {
    "has_monitor": true
  }
}
```

### 2. Получение доступных ресурсов
**GET** `/api/v1/resources`

Поддерживает параметры запроса для фильтрации:
```
GET /api/v1/resources?type=RESOURCE_TYPE_MEETING_ROOM&type=RESOURCE_TYPE_WORKSPACE&location=Office
```

### 3. Получение списка всех ресурсов
**GET** `/api/v1/resources/list`

Поддерживает параметры запроса для фильтрации по типу:
```
GET /api/v1/resources/list?type=RESOURCE_TYPE_DEVICE
```

### 4. Получение ресурса по ID
**GET** `/api/v1/resources/:id`

### 5. Обновление ресурса
**PUT** `/api/v1/resources/:id`

Массив `paths` выступает в роли маски полей (field mask), чтобы указать, какие именно поля нужно обновить.

Пример тела запроса:
```json
{
  "resource": {
    "name": "Meeting Room Omega",
    "location": "Office 2, Floor 1"
  },
  "paths": ["name", "location"]
}
```

### 6. Изменение статуса ресурса (Админ)
**PATCH** `/api/v1/resources/:id/status`

Доступные статусы: `RESOURCE_STATUS_AVAILABLE`, `RESOURCE_STATUS_OCCUPIED`, `RESOURCE_STATUS_MAINTENANCE`, `RESOURCE_STATUS_EMERGENCY`.

Пример тела запроса:
```json
{
  "status": "RESOURCE_STATUS_MAINTENANCE",
  "reason": "Projector maintenance"
}
```

### 7. Удаление ресурса
**DELETE** `/api/v1/resources/:id`

### 8. Проверка статуса ресурса (Бронирование)
**GET** `/api/v1/resources/:id/status`

### 9. Обновление статуса занятости ресурса (Бронирование)
**PATCH** `/api/v1/resources/:id/occupancy`

Используется системой бронирования для отметки ресурса как занятого или свободного.

Пример тела запроса:
```json
{
  "is_occupied": true
}
```

## Эндпоинты авторизации (SSO)

Все запросы из этой группы автоматически проксируются в микросервис SSO с сохранением заголовков и пробросом контекста распределенной трассировки.

### 1. Регистрация
**POST** `/api/v1/auth/register`

### 2. Авторизация (Логин)
**POST** `/api/v1/auth/login`

### 3. Обновление токена
**POST** `/api/v1/auth/refresh`
