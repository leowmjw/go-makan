include main.mk

run: ## Build via go install mechanism
	@go install app/cmd/go-makan && go-makan

make:  ## Get latest main,mk template ;)
	@curl https://raw.githubusercontent.com/sagikazarmark/makefiles/master/go-binary/main.mk > main.mk

