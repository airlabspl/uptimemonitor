.PHONY: watch
watch:
	composer run dev

.PHONY: test
test:
	php artisan test