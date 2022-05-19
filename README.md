# rest_desktop

Простой пример REST-сервиса с клиентом и сервером.

Запуск сервера (Linux):
```bash
cd rest_srv
SERVERPORT=8082 go run rest_srv
```

Запуск клиента (Windows):
```bash
cd rest_clnt
set SERVERPORT=8082
go run rest_clnt
# {"id":0}
go run rest_clnt unregister
#
go run rest_clnt unregister
# desktop with cname=THINKPAD-X230 not found
go run rest_clnt register
# {"id":1}
go run rest_clnt update
#
```

Эксперименты с curl (обработка ошибок):
```bash
curl http://localhost:8083/register/
# expect method POST at /register/, got GET
curl -X POST http://localhost:8083/register/
# mime: no media type
cat params.json
# {"cname":"EPSON-L3150","cip":"10.0.0.3","user":"printer","at":"2022-05-19T20:44:16.6639767Z"}
curl --verbose -X POST -H "Content-Type: application/json" -d @params.json http://localhost:8082/register/
# < HTTP/1.1 200 OK
# {"id":2}
curl --verbose -X DELETE -H "Content-Type: application/json" -d @params.json http://localhost:8082/unregister/
# < HTTP/1.1 200 OK
```
