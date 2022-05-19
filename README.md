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

