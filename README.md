# dwld-downloader

Сервис скачиваний

Тезисы:
1. Скачивает видео с разных ресурсов при помощи yt-dlp
2. Запускает разные варианты конфигураций загрузчика по необходимости (стейджи)
3. Принимает по gRPC задачи на вставку ссылки в очередь, ее удаление, а так же предоставление статистики
4. Задачи помещаются в базу sqlite3 (Большие БД тут ни к чему)
5. Перемещает неудавшиеся скачивания по списку загрузчиков (стейджей)
6. Удавшиеся загрузки передает по ftp на указанный сервер в конфиге. При отсуствии конфига или выключенном флаге - не передает, оставляет по пути указанном в конфиге 
7. Свое состояние хранит в memcahed (полноценный redis тут ни к чему)
8. Так же держит http-echo сервер для обеспечения проверки healt-check
9. При успешной передаче по ftp - исходный файл удаляется

ROADMAP:

[+] Задать архитектуру проекта (Бот + скачивальщик) <br>
[+] Задать структуру проекта <br>
[+] Предположить первичный контракт (по ходу разработки и поддержки может модифицироваться) <br>
[+] Реализовать считывание конфигов <br>
[+] Реализовать healt-check <br>
[+] Реализовать первичную инициализацию, запуск и выключение (gracefull-shutdown) серверов (grpc, http) <br>
[+] Реализовать подключение к sqlite <br>
[+] Реализовать подключение к memcached <br>
[+] Накидать контракты слоя юзкейса <br>
[+] Реализовать контракты в grpc <br>
[+] Реализовать CRUD в sqlite <br>
[+] Реализовать CRUD в memcahed <br>

<br>

[-] Реализовать usecases <br>
[-] Добавить функцию "авторизации" скачивальщиков в сервисе предоставлющим интерфейс пользователя (bff) при их запске<br>
[-] Реализовать стейджовый скачивальщик <br>
[-] Реализовать передачу файлов по ftp <br>

<br>

[-] Написать подробный гайд по запуску этой всей истории <br>

<br>

[-] Приступить к реализации бота  <br>
