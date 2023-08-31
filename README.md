## Slug service

Сервис динамического сегментирования пользователей. Разработан в
качестве [тестового задания](https://github.com/avito-tech/backend-trainee-assignment-2023/blob/main/README.md).
Сервис предоставляет JSON HTTP API.

Реализованы методы:

1. Создание сегмента.
2. Удаление сегмента.
3. Добавление и удаление сегментов пользователю,
4. Получение списков активных сегментов пользователя.

Были реализованы опциональные задания:

1. Получение отчета (попадания/выбывания пользователя из сегмента) по пользователю за определенный период.
2. В методе добавления сегмента пользователю предусмотрено опциональное добавление
   ttl (время автоматического удаления пользователя из сегмента).
3. В методе создания сегмента предусмотрена опция для передачи процента пользователей, которых необходимо добавить
   автоматически в сегмент.

## Стек

1. Язык разработки: Golang.
2. Реляционная СУБД: PostgreSQL.
3. Веб фреймворк: Echo.
4. Объектное хранилище: MinIO.
5. Docker и docker-compose для развертывания сервиса.
6. OpenAPI (Swagger) для API.
7. golang/mock, testify (для тестирования).

## Допущения, сделанные при разработке

1. За регистрацию и авторизацию пользователя отвечает другой сервис. В данном сервисе подразумевается, что приходит
   валидный userID.
2. Имя сегмента уникально.
3. Метод удаления сегмента фактически не удаляет сегмент из базы, а помечает сегмент удаленным. Физически не удаляем
   запись по сегменту для истории.
4. Если сегмент был удален, то его нужно удалить у всех пользователей. Есть риск, что пользователей с удаленными
   сегментом может быть очень много, поэтому для удаления у всех пользователей потребуется много времени. Как исключаем
   риск: при запросе сегментов у пользователя проверяем, удален ли полученный сегмент. Если сегмент был ранее удален, то
   не включаем его в ответ, в фоне удаляем его у пользователя и добавляем в историю операций об удалении.
5. При создании сегмента с "процентом" ставим фоновую задачу через паттерн outbox. Фоновая задача потенциально может
   выполняться очень долго (если пользователей много), и есть вероятность возникновения по ходу выполнения ошибки. В
   таком случае задача прервется. Такое решение было выбрано из-за нехватки времени на более надежную реализацию.
6. Функция "попадания" пользователя в процент. Исходим из того, что идентификатор пользователя – UUID, который можем
   привести к целому числу. Глядя на остаток деления числа на 100 понимаем, попали ли в этот процент (реализация в
   пакете sorting_hat). На большом количестве пользователей должно быть нормальное распределение.
7. Если пользователь создается после того, как "процентный" сегмент был создан, то он в него не попадет. Не хватило
   времени на реализацию.
8. В задании указано, что формирование отчета происходит по одному пользователю. При этом в примере отчета указаны
   разные пользователи. Следую тексту задания, поэтому отчет формируется по одному пользователю.

## Запуск приложения и зависимостей

1. Склонировать репозиторий.

```shell
git clone https://github.com/frutonanny/slug-service.git
```

2. Перейти в репозиторий.

```shell
cd slug-service
```

3. Выполнить команду запуска приложения и его зависимостей.

```shell
docker-compose -f deployments/docker-compose.yaml --profile dev up --build --detach
```

## Основной сценарий тестирования приложения

1. Добавляем сегмент "AVITO_VOICE_MESSAGES".

```shell
curl -X POST --location "http://localhost:8081/v1/createSlug" \
    -H "Content-Type: application/json" \
    -d "{
          \"name\": \"AVITO_VOICE_MESSAGES\"
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

2. Добавляем сегмент "AVITO_PERFORMANCE_VAS".

```shell
curl -X POST --location "http://localhost:8081/v1/createSlug" \
    -H "Content-Type: application/json" \
    -d "{
          \"name\": \"AVITO_PERFORMANCE_VAS\"
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

3. Добавляем пользователю сегменты.

```shell
curl -X POST --location "http://localhost:8081/v1/modifyUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"53aa35c8-e659-44b2-882f-f6056e443c99\",
          \"add\": [
            {
              \"name\": \"AVITO_VOICE_MESSAGES\"
            },
            {
              \"name\": \"AVITO_PERFORMANCE_VAS\"
            }
          ]
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

4. Запрашиваем у пользователя сегменты.

```shell
curl -X POST --location "http://localhost:8081/v1/getUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"53aa35c8-e659-44b2-882f-f6056e443c99\"
        }"
```

В ответе ожидаем получить два сегмента "AVITO_VOICE_MESSAGES" и "AVITO_PERFORMANCE_VAS".

```json
{
  "data": {
    "slugs": [
      "AVITO_VOICE_MESSAGES",
      "AVITO_PERFORMANCE_VAS"
    ]
  }
}

```

5. Удаляем у пользователя сегмент.

```shell
curl -X POST --location "http://localhost:8081/v1/modifyUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"53aa35c8-e659-44b2-882f-f6056e443c99\",
          \"delete\": [
            \"AVITO_PERFORMANCE_VAS\"
          ]
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

6. Запрашиваем у пользователя сегменты.

```shell
curl -X POST --location "http://localhost:8081/v1/getUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"53aa35c8-e659-44b2-882f-f6056e443c99\"
        }"
```

В ответе ожидаем получить один оставшийся сегмент.

```json
{
  "data": {
    "slugs": [
      "AVITO_VOICE_MESSAGES"
    ]
  }
}
```

7. Удаляем сегмент.

```shell
curl -X POST --location "http://localhost:8081/v1/deleteSlug" \
    -H "Content-Type: application/json" \
    -d "{
          \"name\": \"AVITO_VOICE_MESSAGES\"
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

8. Запрашиваем у пользователя сегменты.

```shell
curl -X POST --location "http://localhost:8081/v1/getUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"53aa35c8-e659-44b2-882f-f6056e443c99\"
        }"
```

В ответе ожидаем получить пустой список сегментов.

```json
{
  "data": {
    "slugs": []
  }
}
```

9. Получаем отчет по пользователю.

```shell
curl -X POST --location "http://localhost:8081/v1/getReport" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"53aa35c8-e659-44b2-882f-f6056e443c99\",
          \"period\": \"2023-09\"
        }"
```

В ответе ожидаем получить ссылку на отчет.

```json
{
  "data": {
    "url": "localhost:9000/reports/report-2023-08-cb06954f-c14f-4ee6-9f44-f2b92791d55d.csv"
  }
}
```

Ссылку можно вставить в браузер, отчет должен скачаться.
Отчет должен содержать 4 события.

```csv
53aa35c8-e659-44b2-882f-f6056e443c99,AVITO_VOICE_MESSAGES,add slug,2023-08-31T20:04:16+03:00
53aa35c8-e659-44b2-882f-f6056e443c99,AVITO_PERFORMANCE_VAS,add slug,2023-08-31T20:04:16+03:00
53aa35c8-e659-44b2-882f-f6056e443c99,AVITO_PERFORMANCE_VAS,delete slug,2023-08-31T20:05:00+03:00
53aa35c8-e659-44b2-882f-f6056e443c99,AVITO_VOICE_MESSAGES,delete slug,2023-08-31T20:05:41+03:00
```

## Сценарий проверки "процентных" сегментов

1. Очистим базу данных, удалив и создав заново все зависимости.

```shell
docker-compose -f deployments/docker-compose.yaml down
docker-compose -f deployments/docker-compose.yaml --profile dev up --build --detach
```

2. Создаем пользователя с определенным ID, который попадет в нужный процент.

"030c48d3-12a3-4ea0-a334-d9a141c7a1b8" -> 51136723.
51136723 % 100 < 30.

```shell
curl -X POST --location "http://localhost:8081/v1/modifyUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"030c48d3-12a3-4ea0-a334-d9a141c7a1b8\"
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

3. Создаем пользователя с определенным ID, который не попадет в нужный процент.

"9305ea2f-0c7c-4a08-a1fa-918a496e1701" -> 2466638383.
2466638383 % 100 > 30.

```shell
curl -X POST --location "http://localhost:8081/v1/modifyUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"9305ea2f-0c7c-4a08-a1fa-918a496e1701\"
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

4. Добавляем сегмент, в который должны войти 30% пользователей. В том числе пользователь, созданный выше.

```shell
curl -X POST --location "http://localhost:8081/v1/createSlug" \
    -H "Content-Type: application/json" \
    -d "{
          \"name\": \"AVITO_DISCOUNT_30\",
          \"options\": {
            \"percent\": 30
          }
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

5. Запрашиваем у попавшего пользователя сегменты.

```shell
curl -X POST --location "http://localhost:8081/v1/getUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"030c48d3-12a3-4ea0-a334-d9a141c7a1b8\"
        }"
```

В ответе ожидаем получить процентный сегмент.

```json
{
  "data": {
    "slugs": [
      "AVITO_DISCOUNT_30"
    ]
  }
}
```

6. Запрашиваем у не попавшего пользователя сегменты.

```shell
curl -X POST --location "http://localhost:8081/v1/getUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"9305ea2f-0c7c-4a08-a1fa-918a496e1701\"
        }"
```

В ответе ожидаем получить пустые сегменты.

```json
{
  "data": {
    "slugs": []
  }
}
```

## Сценарий проверки сегментов с TTL

1. Очистим базу данных, удалив и создав заново все зависимости.

```shell
docker-compose -f deployments/docker-compose.yaml down
docker-compose -f deployments/docker-compose.yaml --profile dev up --build --detach
```

2. Добавляем сегмент, который будем у пользователя делать просроченным.

```shell
curl -X POST --location "http://localhost:8081/v1/createSlug" \
    -H "Content-Type: application/json" \
    -d "{
          \"name\": \"AVITO_PERFORMANCE_VAS_07_2023\"
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

3. Добавляем сегмент, который будем у пользователя делать непросроченным.

```shell
curl -X POST --location "http://localhost:8081/v1/createSlug" \
    -H "Content-Type: application/json" \
    -d "{
          \"name\": \"AVITO_PERFORMANCE_VAS_07_2033\"
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

4. Добавляем пользователя с просроченным и непросроченным сегментом.

```shell
curl -X POST --location "http://localhost:8081/v1/modifyUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"53aa35c8-e659-44b2-882f-f6056e443c99\",
          \"add\": [
            {
              \"name\": \"AVITO_PERFORMANCE_VAS_07_2023\",
              \"ttl\": \"2023-08-01T00:00:00Z\"
            },
            {
              \"name\": \"AVITO_PERFORMANCE_VAS_07_2033\",
              \"ttl\": \"2033-08-01T00:00:00Z\"
            }
          ]
        }"
```

В ответе ничего не ожидаем получить. Только HTTP 200.

```json
{}
```

5. Запрашиваем у пользователя сегменты.

```shell
curl -X POST --location "http://localhost:8081/v1/getUserSlugs" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": \"53aa35c8-e659-44b2-882f-f6056e443c99\"
        }"
```

В ответе ожидаем получить только непросроченный сегмент.

```json
{
  "data": {
    "slugs": [
      "AVITO_PERFORMANCE_VAS_07_2033"
    ]
  }
}
```