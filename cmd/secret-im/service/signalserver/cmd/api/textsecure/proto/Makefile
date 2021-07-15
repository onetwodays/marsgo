all:
	protoc --go_out=../../queue  queue_protocol.proto
	protoc --go_out=.. PubSubMessage.proto TextSecure.proto  WebSocketResources.proto
	protoc --go_out=.. Groups.proto StorageService.proto