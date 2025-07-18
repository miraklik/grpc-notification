# grpc-notification

`grpc-notification` — это сервис для обработки и отправки уведомлений с использованием **Go**, **gRPC**, **PostgreSQL**, **NATS**, а теперь также с поддержкой **Docker** и **Kubernetes**.

Проект реализует масштабируемую и отказоустойчивую архитектуру для доставки сообщений в реальном времени.

---

## ⚙️ Основные технологии

- **Go** — язык программирования для реализации серверной логики
- **gRPC** — высокопроизводительный фреймворк для удалённого вызова процедур
- **PostgreSQL** — реляционная база данных
- **NATS** — брокер сообщений для асинхронной передачи событий
- **Docker / Docker Compose** — контейнеризация и удобный запуск окружения
- **Kubernetes** — оркестрация и масштабирование микросервиса

---

## 🚀 Возможности

- Приём уведомлений через gRPC API
- Асинхронная обработка сообщений через NATS
- Хранение данных в PostgreSQL
- Запуск через Docker и масштабирование в Kubernetes

---

## 🐳 Быстрый запуск через Docker Compose

```bash
git clone https://github.com/miraklik/grpc-notification.git
cd grpc-notification
docker-compose up --build
```

## ☸️ Деплой в Kubernetes

1. Создай секрет с переменными окружения:
```bash
kubectl apply -f k8s/secret.yaml
```

2. Разверни сервис:
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

3. Проверь статус:
```bash 
kubectl get pods
kubectl get svc
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