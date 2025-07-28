test-init-site:
	rm -r ./tmp/test-site; \
	go run main.go new site test-site; \
	go build -o pride; \
	mv test-site ./tmp; \
	mv pride ./tmp/test-site/pride \

test-init-mock-content:
	cd ./tmp/test-site; \
	./pride new content about; \
	./pride new content posts/go_and_rust; \
	./pride new content posts/js_and_python; \
	./pride new content posts/angular_and_react; \
	./pride new content articles/do_you_dev; \
	./pride new content articles/burnout_or_hungry; \
	./pride new content posts/archive/old.md; \
	./pride new content posts/archive/older.md; \

test-build-nav:
	cd ./tmp/test-site; \
	./pride build nav; \