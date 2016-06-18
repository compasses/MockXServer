default:build
build:
	@bash getversion.sh
	@go build -o MockXServer/MockXServer.run
