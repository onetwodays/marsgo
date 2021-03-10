#goctl api plugin -plugin goctl-swagger="swagger -filename secret-im-swagger.json" -api ./signalserver.api -dir .

goctl api plugin -plugin goctl-swagger="swagger" -api signalserver.api -dir .