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
	make test-init-site; \
	make test-init-mock-content; \
	cd ./tmp/test-site; \
	./pride build nav; \
	cd navigation; \
	../pride build nav; \

test-build-nav-root:
	cd ./tmp/test-site; \
	./pride build nav; \

test-build-nav-inner:
	cd ./tmp/test-site/navigation; \
	../pride build nav; \

test-publish:
	make test-init-site; \
	make test-init-mock-content; \
	cd ./tmp/test-site/content; \
	../pride publish ./index.md; \

