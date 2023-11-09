-include .env


export OPENAI_KEY ?= ""

assistant:
	export OPENAI_KEY=$(OPENAI_KEY) && \
	go run cmd/assistant/main.go

conversation:
	export OPENAI_KEY=$(OPENAI_KEY) && \
	go run cmd/conversation/main.go