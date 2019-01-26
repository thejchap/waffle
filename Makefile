all: build-docker run-docker

build-docker:
	docker build -t thejchap/waffle:dev .

run-docker:
	docker run -p 3000:3000 thejchap/waffle:dev

test:
	go test ./...

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint
	./node_modules/eslint/bin/eslint.js .
