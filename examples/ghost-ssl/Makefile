start: certs
	docker-compose up


stop:
	docker-compose down
	rm -rf certs

certs:
	mkdir -p ./certs
	mkcert -cert-file certs/blog.test.localhost -key-file certs/blog.test.localhost.key blog.test.localhost
