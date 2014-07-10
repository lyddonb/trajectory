source_js := $(wildcard web/app/*.js)
build_js := $(source_js:%.js=%.min.js)
app_bundle := web/js/app.js

buildredis:
	docker build -t lyddonb/redis redis

runredis:
	docker run --name redis -d -p 6379:6379 lyddonb/redis

killredis:
	docker stop redis
	docker rm redis

buildtrajectory:
	docker build -t lyddonb/trajectory .

runtrajectory:
	docker run --name trajectory -i -t --link redis:redis -d -p 1301:1301 -p 3001:3001 lyddonb/trajectory /bin/bash

killtrajectory:
	docker stop trajectory
	docker rm trajectory

rundebug:
	docker run -v /vagrant:/opt/go/src/github.com/lyddonb/trajectory --link redis:db -i -t lyddonb/trajectory /bin/bash 

deps:
	go get github.com/garyburd/redigo/redis

watch:
	#jsx --watch web/app/ web/js/
	watchify web/app/app.js -d -o web/js/app.js -v
	#watchify -o web/js/app.js  -v -d .

buildjs:
	browserify web/app/app.js -o web/js/app.js

buildcss:
	lessc web/less/*.less > web/css/main.css

buildall: installjsdeps buildjs buildcss
	apt-get install go build

installjsdeps:
	npm install .
	npm install -g browserify
	npm install -g less

dockerrun: 
	./trajectory --redis-port=$(REDIS_PORT_6379_TCP_PORT)

#buildjs: $(app_bundle)

#%.min.js: %.js
	#cat $^ >$@

#$(app_bundle): $(build_js)
	#uglifyjs -o $@ $<
	#echo >> $@
	#rm -f $<

.PHONY: buildjs
