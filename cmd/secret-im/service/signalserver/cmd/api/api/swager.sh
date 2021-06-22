#goctl api plugin -plugin goctl-swagger="swagger -filename secret-im-swagger.json" -api ./signalserver.api -dir .

goctl api plugin -plugin goctl-swagger="swagger -filename user.json" -api signalserver.api -dir .
docker run --rm -p 8083:8080 -e SWAGGER_JSON=/home/zh/go/marsgo/cmd/secret-im/service/signalserver/cmd/api/api/user.json  swaggerapi/swagger-ui