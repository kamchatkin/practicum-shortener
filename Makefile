run:
	go run ./cmd/shortener

get:
	curl -X GET -i 'http://localhost:8080/qwerty'

getz:
	curl -X GET -i 'http://localhost:8080/qwerty' \
		-H "Accept-Encoding: gzip" --compressed

post:
	curl -X POST -i 'http://localhost:8080/' \
		-d 'https://ya.ru/'

postz:
	curl -X POST -i 'http://localhost:8080/' \
		-d 'https://ya.ru/' \
		-H "Accept-Encoding: gzip" --compressed

api:
	curl -X POST -i 'http://localhost:8080/api/shorten' \
 		-d '{"url":"https://ya.ru/"}' \
 		-H "Content-Type: application/json"

apiz:
	curl -X POST -i 'http://localhost:8080/api/shorten' \
		-d '{"url":"https://ya.ru/"}' \
		-H "Content-Type: application/json" \
		-H "Accept-Encoding: gzip" \
		--compressed

c:
	docker exec -ti goadv bash
