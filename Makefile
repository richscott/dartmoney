client: webroot/js/dist/bundle.js
server: dartmoney-server
all: server client

dartmoney-server: *.go db/*.go
	go build dartmoney-server.go

webroot/js/dist/bundle.js:
	npx webpack --config webpack.config.js

clean:
	$(RM) -r webroot/js/dist dartmoney-server
