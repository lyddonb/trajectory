

buildredis:
	docker build -t lyddonb/redis redis

runredis:
	docker run --name redis -d -p 6379:6379 lyddonb/redis

build:
	docker build -t lyddonb/trajectory .

run:
	docker run --link redis:db -i -t lyddonb/trajectory /bin/bash

rundebug:
	docker run -v /vagrant:/opt/go/src/github.com/lyddonb/trajectory --link redis:db -i -t lyddonb/trajectory /bin/bash 

deps:
	go get github.com/garyburd/redigo/redis
