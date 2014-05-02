source_js := $(wildcard web/app/*.js)
build_js := $(source_js:%.js=%.min.js)
app_bundle := web/js/app.js

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

watch:
	#jsx --watch web/app/ web/js/
	watchify web/app/app.js -d -o web/js/app.js -v
	#watchify -o web/js/app.js  -v -d .

buildjs:
	browserify web/app/app.js -o web/js/app.js

buildcss:
	lessc web/less/*.less > web/css/main.css

#buildjs: $(app_bundle)

#%.min.js: %.js
	#cat $^ >$@

#$(app_bundle): $(build_js)
	#uglifyjs -o $@ $<
	#echo >> $@
	#rm -f $<

.PHONY: buildjs
