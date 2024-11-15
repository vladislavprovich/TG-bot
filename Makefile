# name file
LINTER = golangci-lint


all: lint

# download golangci-lint'
install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# linter check code
lint:
	$(LINTER) run

# auto fix problem
lint-fix:
	$(LINTER) run --fix

# clear generate file (optional)
clean:
	rm -rf $(LINTER)

# help list
help:
	@echo "Makefile для запуску Go лінтера"
	@echo "Доступні команди:"
	@echo "  install-linter  - інсталяція golangci-lint"
	@echo "  lint            - запуск лінтера для перевірки коду"
	@echo "  lint-fix        - запуск лінтера з автоматичним виправленням помилок"
	@echo "  clean           - чистка середовища (опціонально)"

.PHONY: all lint lint-fix clean help