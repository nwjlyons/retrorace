CURRENT_DIR = $(shell pwd)

e:
	cd elm; elm format --yes *.elm
	cd elm; elm make *.elm --output=../static/retrorace.js

format:
	gofmt -w $(CURRENT_DIR)

deploy:
	gofmt -w $(CURRENT_DIR)
	GOOS=linux GOARCH=amd64 go build
	ssh neillyons.io "supervisorctl stop retrorace"
	rsync -P retrorace neillyons.io:/srv/www/retrorace.neillyons.io
	rsync -P -r templates neillyons.io:/srv/www/retrorace.neillyons.io
	rsync --include="*.css" --include="*.js" -P -r static neillyons.io:/srv/www/retrorace.neillyons.io
	ssh neillyons.io "supervisorctl start retrorace"
