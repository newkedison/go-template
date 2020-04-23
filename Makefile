first: all

all: TEMPLATE

TEMPLATE: main.go
	go build -o $@

run: TEMPLATE
	chmod +x TEMPLATE && ./TEMPLATE
