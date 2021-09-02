include main.mk

run: ## Build via go install mechanism
	@go install app/cmd/go-makan && go-makan

make:  ## Get latest main,mk template ;)
	@curl https://raw.githubusercontent.com/sagikazarmark/makefiles/master/go-binary/main.mk > main.mk

makelib: ## Get template  for lib ;)
	@curl https://raw.githubusercontent.com/sagikazarmark/makefiles/master/go-library/main.mk > main.mk

makeapp: ## Get template for app
	@curl https://raw.githubusercontent.com/sagikazarmark/makefiles/master/go-app/main.mk > main.mk

