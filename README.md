# rest_desktop

Простой пример REST-сервиса с клиентом и сервером.

Подготовка к работе:
```bash
go get github.com/lib/pq

sudo -u postgres psql --port=5432 testdb testuser
> CREATE SEQUENCE seq_rest_srv_table
    INCREMENT 1
    START 1000
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;
OK

> CREATE TABLE rest_srv_table
(
    id INTEGER NOT NULL DEFAULT NEXTVAL('seq_rest_srv_table'::regclass), -- первичный ключ
    cname CHARACTER VARYING(100) NOT NULL, -- имя компьютера
    cip CHARACTER VARYING(22) NOT NULL, -- ip-адрес компьютера (для ipv6 надо 22 символа)
    "user" CHARACTER VARYING(100) NOT NULL, -- имя пользователя
    "at" TIMESTAMP NOT NULL,
    CONSTRAINT pk_rest_srv PRIMARY KEY (id),
    CONSTRAINT unq_rest_srv_cname UNIQUE (cname)
)
TABLESPACE pg_default;
OK
```

Запуск сервера (Linux):
```bash
go get github.com/lib/pq

cd rest_srv
export SERVERURL=10.0.0.2
export SERVERPORT=8082
export DBURL=postgres://testuser:testpassword@10.0.0.2/testdb?sslmode=disable
go run rest_srv
```

Запуск клиента (Windows):
```bash
cd rest_clnt
set SERVERURL=10.0.0.2
set SERVERPORT=8082
go run rest_clnt
# {"id":1000}
go run rest_clnt unregister
#
go run rest_clnt unregister
# desktop with cname=THINKPAD-X230 not found
go run rest_clnt register
# {"id":1001}
go run rest_clnt update
#
```

Эксперименты с curl (обработка ошибок):
```bash
curl http://10.0.0.2:8083/register/
# expect method POST at /register/, got GET
curl -X POST http://10.0.0.2:8083/register/
# mime: no media type
cat params.json
# {"cname":"EPSON-L3150","cip":"10.0.0.3","user":"printer","at":"2022-05-19T20:44:16.6639767Z"}
curl --verbose -X POST -H "Content-Type: application/json" -d @params.json http://10.0.0.2:8082/register/
# < HTTP/1.1 200 OK
# {"id":1002}
curl --verbose -X DELETE -H "Content-Type: application/json" -d @params.json http://10.0.0.2:8082/unregister/
# < HTTP/1.1 200 OK
```

## Ремарка по операциям c БД

Работа с БД сделана "на скорую руку", фактически операции `INSERT` и `UPDATE` объединены в один общий `INSERT` с дополнительным `ON CONFLICT ON CONSTRAINT unique DO UPDATE SET`, тем самым я решаю три проблемы сразу:
 - что делать, если нарушено ограничение уникальности?
 - лень было делать отдельно добавление и отдельно обновление (при желании можно посмотреть ветку `v0.1-local_store`, где методы `CreateDesktop` и `UpdateDesktop` различаются)
 - что делать, если desktop не зарегистрирован при попытке его обновить?

Но так делать нельзя, потому что не следует забывать про необходимость поддержания БД в согласованном состоянии, и поэтому необходимо реализовать операции изменения данных в БД с использованием оптимистического блокирования, т.е. в WHERE должна попасть и старое значение и новое. Таким образом `UpdateDesktop` должен делать что-то вроде:
```sql
UPDATE public.rest_srv_table
SET cname=$new_cname,cip=$new_ip,user=$new_user,at=$at
WHERE cname=$old_cname,cip=$old_ip,user=$old_user
RETURNING id;
```

Таким образом на выходе возможны три ситуации:
 - ошибка с откатом транзакции, если кто-то успел раньше нас?
 - операция изменения прошла удачно и мы узнали id
 - операция изменения не удалась, id мы не узнали, значит кто-то успел обновить информацию в БД раньше нас, а наши данные устарели... значит надо принять решения что делать дальше? повторить? уведомить пользователя? забить? и т.д. и т.п - это именно та причина, по которой все эти REST-интерфейсы прихо ложаться в работу с реляционными БД, потому что на одной стороне "ACID и вот это вот всё" а с другой curl-утилита, которой пофик на целостность данных в БД и на её согласованность ;)

## Бонус 1

Чтобы программа работала и в Windows, и Linux, все зависимости и настройки вынесены в environment variables.

## Бонус 2

Не стал делать, но...
...для работы сервера как systemd-сервис, просто добавляем его запуск в /etc/systemd/system/ с полными указаниями путей. Тут по хорошему надо бы подключить к программе `libsystemd-dev` (не разбирался как это в golang) и дёргать там software watch-dog, чтобы systemd за нас перезапустил рухнувший или зависший сервис. А что, а вдруг?!

В упаковке deb-пакета вижу только одну проблему: разруливание зависимости от `github.com/lib/pq`. Тут у меня уже второй час ночи, и... пошёл ка я спать)