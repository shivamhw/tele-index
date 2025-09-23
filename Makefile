clean: 
	- rm tele-index
build:
	go build -o tele-index
run: clean build
	./tele-index