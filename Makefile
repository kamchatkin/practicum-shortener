PORT=3000

run:
	SERVER_ADDRESS=:3000 DATABASE_DSN=postgres://goadv:goadv@goadv-db:5432/goadv go run ./cmd/shortener

get:
	curl -X GET -i 'http://localhost:${PORT}/qwerty'

getz:
	curl -X GET -i 'http://localhost:${PORT}/qwerty' \
		-H "Accept-Encoding: gzip" --compressed

post:
	curl -X POST -i 'http://localhost:${PORT}/' \
		-d 'https://ya.ru/'

postz:
	curl -X POST -i 'http://localhost:${PORT}/' \
		-d 'https://ya.ru/' \
		-H "Accept-Encoding: gzip" --compressed

api:
	curl -X POST -i 'http://localhost:${PORT}/api/shorten' \
 		-d '{"url":"https://ya.ru/"}' \
 		-H "Content-Type: application/json"

apiz:
	curl -X POST -i 'http://localhost:${PORT}/api/shorten' \
		-d '{"url":"https://ya.ru/"}' \
		-H "Content-Type: application/json" \
		-H "Accept-Encoding: gzip" \
		--compressed

c:
	docker exec -ti goadv bash

tern:

