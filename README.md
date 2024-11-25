### About
Sample REST API Service 

### Установка
1. Клонировать репозиторий
2. Создайте базу данных `music`.
3. Проверьте настройки PosgreSQL и сервера в файле `.env`. Сервер запускается на порту 8080 по умолчанию.
4. Запустите сваггер, url по умолчанию: http://localhost:8080/docs/index.html 

В терминале линукс:
```shell
git clone https://github.com/hectorbarbosa/music.git
make createdb
# build
make
# Start server (port 8080 by default)
make run
```
5. Можно вручную сделать миграцию `down` из Makefile:
```shell
make migratedown
```
6. Вспомогательный сервер с доп. информацией о песнях можно запустить из Makefile:
```shell
# build
make buildapi
# Start server (port 8080 by default)
make runapi
```
