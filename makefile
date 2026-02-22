build:
	go build -o hackpad .
# @ is used to prevent the command from being printed in the terminal
# when we run make run, it will first build the project and then execute the binary
run: build
	@./hackpad 