# 2019_1_SleeplessNights
##Микросервисы системы
###Consul
Consul выполняет роль service discovery и key-value хранилища для конфигурации системы.  
Внешние зависимости:
* docker
* jq  
Consul собирается из оффициального контейнера, чтобы его получить нужно сделать  
`docker pull consul`  
Сам consul и интерфейсы микросервисов для него конфигурируются в файле consul/config.json, этот файл имеет особую структуру, подробнее см. в документации к consul
## О конфигурации системы
Модель конфига:  
```
{  
  "dev": {
    "log_level": "trace",  
    "auth": {  
      "cookie_name": "session_token",  
      "expiration_time": 86400,  
      "secret": "-"  
    },  
    "cors": {  
      "credentials": true,  
      "domains": "https://sleepless-nights--frontend.herokuapp.com",  
      "headers": [  
        "X-Requested_With",  
        "Content-type",  
        "User-Agent",  
        "Cache-Control",  
        "Cookie",  
        "Origin",  
        "Accept-Encoding",  
        "Connection",  
        "Host",  
        "Upgrade-Insecure-Requests",  
        "User-Agent",  
        "Referer",  
        "Access-Control-Request-Method",  
        "Access-Control-Request-Headers"  
      ],  
      "methods": [  
        "GET",  
        "POST",  
        "PATCH",  
        "DELETE",  
        "OPTIONS"  
      ]  
    },  
    "db": {  
      "host": "localhost",  
      "port": 5432,  
      "user": "maxim",  
      "password": "-",  
      "dbname": "questions"  
    },  
    "game": {  
      "max_rooms": 100,  
      "room": {  
        "io_queues_len": 50,  
        "main_loop_interval": 500  
      }  
    },  
    "grpc_servers": {  
      "main": {  
        "host": "127.0.0.1",  
        "port": "8000"  
      },  
      "game": {  
        "host": "127.0.0.1",  
        "port": "8001"  
      },  
      "auth": {  
        "host": "127.0.0.1",  
        "port": "8002"  
      }  
    },  
    "http_servers": {  
      "main": {  
        "host": "127.0.0.1",  
        "port": "8080"  
      },  
      "game": {  
        "host": "127.0.0.1",  
        "port": "8081"  
      }  
    },  
    "scripts_path": "/home/maxim/Dev/Technopark/2_semester/Golang/2019_1_SleeplessNights/main_microservice/scripts",  
    "static_path": "/home/maxim/Dev/Technopark/2_semester/Golang/2019_1_SleeplessNights/main_microservice/static"  
  },
  
  "test": {
  },  
  "local": {
  },  
  "prod": {
  }  
}
```
## Об архитектуре игры
Сейчас все файлы, относящиеся к бизнес логике игры, лежат в директории TheGame. Уже в следующем модуле нам понадобится выделять эту логику в отдельный микросервис, поэтому на данном этапе очень важно, чтобы не возникало никаких внешних зависимостей.
В игре присутствуют следующие сущности:
* Game [синглтон, фасад] - объект-интерпритация всей бизнес логики игры. Если мы извне хотим сделать что-то, что относится к игре, то нам определённо нужно обратиться именно к этому объекту. Знает о всех вспомагательных сущностях сервиса, типа комнат, балансировщика и т.д.
* Room [фасад] - объекты, несущие ответственность непосредственно за игровой процесс. Ничего не знает о том, как устроено приложение, но знает всё о процессе игры. Здесь мы работаем с игроками (ВАЖНО игрок != пользователь), вопросами, игровым циклом и т.д.
* Player [посредник] - интерфейс для взаимодействия с внешним миром. На самом деле, игровой логике абсолютно всё равно, с кем она работает, поэтому в этой сущностиважно не завязываться конкретно на игрока-пользователя. Именно поэтому Player - интерфейс, а не структура. Такой подход, возможно, несколько усложняет реализацию, но зато в результате для нашей системы не будет разницы, кто есть игрок, это может быть клиент по websocket, клиент по HTTP2, локальный бот с двумя каналами - что угодно, реализующее уникальное дуплексное соединение.
* PlayerFactory [синглтон, фабричный метод (не путать с абстрактной фабрикой)] - раз наш игрок весь из себя такой абстрактный, то нужно решить вопрос с тем, откуда таких игроков брать. Я вижу две реализации - открыть интерфейс Player наружу и предложить остальным подсистемам самим готовить игроков, или сделать по методу для каждого возможного источника, откуда Player может взяться. Второй вариант семантически и для реализации гораздо проще.
* GameField - по сути просто обёртка для матрицы игрового поля. Нужна просто для разделения ответственности внутри комнаты.
* QuestionPack - обёртка для удобного представления хранимых данных, тоже ничего сложного.
* Message - структура, реализующая шаблон сообщения, которыми общаются клиент и сервер.
####В квадратных скобках перечислены используемые паттерны проектирования, если плохо представляете себе сущность, попробуйте погуглить эти слова, ну и спросите меня
