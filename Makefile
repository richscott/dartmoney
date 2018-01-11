client: webroot/js/dist/index.js
server: dartmoney-server
all: server client

dartmoney-server: **/*.go
	go build dartmoney-server.go

webroot/js/dist/index.js:
	parcel build -d webroot/js/dist webroot/js/index.js

clean:
	$(RM) -r webroot/js/dist dartmoney-server
