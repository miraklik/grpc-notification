# grpc-notification

grpc-notification — это сервис для обработки и отправки уведомлений с использованием Go, gRPC, PostgreSQL и NATS. Проект реализует масштабируемую и отказоустойчивую архитектуру для доставки сообщений в реальном времени.

---

## Основные технологии

- **Go** — язык программирования для реализации серверной логики
- **gRPC** — высокопроизводительный фреймворк для удалённого вызова процедур
- **PostgreSQL** — реляционная база данных для хранения данных уведомлений и состояния
- **NATS** — легковесный брокер сообщений для асинхронной коммуникации и передачи событий

---

## Возможности

- Приём уведомлений через gRPC API
- Асинхронная обработка сообщений через NATS
- Надёжное хранение данных в PostgreSQL
- Масштабируемая архитектура с поддержкой нескольких подписчиков и обработчиков

---

## Установка и запуск

### Предварительные требования

- Go 1.18+
- PostgreSQL 12+
- NATS Server

### Шаги

1. Клонируйте репозиторий

   ```bash
   git clone https://github.com/miraklik/grpc-notification.git
   cd grpc-notification
   make build
   ./grpc-notification
   ```

## Использование

### Сервис предоставляет gRPC API для отправки уведомлений. Пример запроса:

``` proto
    service NotificationService {
        rpc SendNotification(SendNotificationRequest) returns (SendNotificationResponse);
        rpc GetStatus(GetStatusRequest) returns (GetStatusResponse);
    }

    message SendNotificationRequest  {
        int64 user_id = 1;
        string type = 2;
        string message = 3;
        int32 priority = 4;
        int64 scheduled_at = 5;
    }

    message SendNotificationResponse {
        int64 notification_id = 1;
    }

    message GetStatusRequest {
        int64 notification_id = 1;
    }

    message GetStatusResponse {
        string status = 1; 
        int32 attempts = 2;
        string last_error = 3;
        int64 delivered_at = 4;
    }
```
