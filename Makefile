all:
	@cat Makefile | grep : | grep -v PHONY | grep -v @ | sed 's/:/ /' | awk '{print $$1}' | sort

#-------------------------------------------------------------------------------

.PHONY: install-deps
install-deps:
	glide install && gometalinter.v2 --install

.PHONY: build
build:
	go build -ldflags="-s -w" -o bin/dockerhub main.go

.PHONY: lint
lint:
	gometalinter.v2 ./main.go

.PHONY: package
package: build
	upx --brute bin/dockerhub
	mkdir -p docker-hub
	rm -Rf docker-hub
	mkdir -p docker-hub
	cp -rv bin docker-hub/
	cp -v *.png docker-hub/
	cp -v *.plist docker-hub/
	cd docker-hub/ && \
		zip -r docker-hub.zip * && \
		mv -v docker-hub.zip ../docker-hub.alfredworkflow

#-------------------------------------------------------------------------------

.PHONY: tag
tag:
	@ if [ $$(git status -s -uall | wc -l) != 1 ]; then echo 'ERROR: Git workspace must be clean.'; exit 1; fi;

	@echo "This release will be tagged as: $$(cat ./VERSION)"
	@echo "This version should match your release. If it doesn't, re-run 'make version'."
	@echo "---------------------------------------------------------------------"
	@read -p "Press any key to continue, or press Control+C to cancel. " x;

	@echo " "
	@chag update $$(cat ./VERSION)
	@echo " "

	@echo "These are the contents of the CHANGELOG for this release. Are these correct?"
	@echo "---------------------------------------------------------------------"
	@chag contents
	@echo "---------------------------------------------------------------------"
	@echo "Are these release notes correct? If not, cancel and update CHANGELOG.md."
	@read -p "Press any key to continue, or press Control+C to cancel. " x;

	@echo " "

	git add .
	git commit -a -m "Preparing the $$(cat ./VERSION) release."
	chag tag --sign

#-------------------------------------------------------------------------------

.PHONY: version
version:
	@echo "Current version: $$(cat ./VERSION)"
	@read -p "Enter new version number: " nv; \
	printf "$$nv" > ./VERSION
