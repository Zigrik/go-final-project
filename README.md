# Финальный проект курса Яндекс Практикум "Go-разработчик с нуля"
Это веб сервер, на котором реализован планировщик задач с хранением этих задач в базе данных.

Все задачи со ☆ выполнены.

По умолчанию запрос сервера http://localhost:7540/
Порт 7540
БД scheduler.db
Пароль по умолчанию не установлен. При сборке Докер-образа пароль "yandex". Токен в tests\settings.go прописан под этот пароль

**Параметы tests\settings.go:**

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6Wzk0LDIwMSwxOTQsMTk0LDE0NSw1Myw1NiwxODcsMjM3LDE3MywyMTcsMTI0LDM3LDE0OCw2MSwxOTcsMjU0LDIyMywxMzQsMjMsMjMwLDE5MCwyMTcsMTI4LDExMyw3OSw5NCwxNTQsMTU1LDQ2LDM2LDIzMF19.GRO2foVtAg8K6iZ6qX1R-Wf4v1VSjMKeiHgaL2klhI4`

**Cборка Докер-образа:**
docker build -t yandex-final .

Запуск контейнера:
docker run -d -p 7540:7540 -v $(pwd)/data:/app yandex-final

Запуск контейнера с переменными окружения:
docker run -d -p 7540:7540 -v $(pwd)/data:/app -e TODO_PORT=7540 -e TODO_DBFILE=/app/scheduler.db -e TODO_PASSWORD=yandex yandex-final

**От себя:**
-для удобства проверки старался максимально следовать инструкциям задания, handler`ы разнесены по соответствующим .go файлам
-из того что не было в задани добавил logger для отслеживания ошибок, это немного усложнило проект, но я посчитал это необходимым
-marshal/unmarshal json функции добавил в api.go
-afterNow функция максимально упрощена, т.к. ее функционал нужно было применить в нескольких функциях, а условия ее применения различались. В итоге код этой функции ушел напрямую в эти функции.