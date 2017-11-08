NAME      := git-archivist
VERSION   := 0.2.2
TYPE      := beta
COMMIT    := $(shell git rev-parse HEAD)
IMAGE     := quay.io/samsung_cnct/git-archivist
TAG       ?= rc

gox:
	@go get github.com/mitchellh/gox
 
dep:
	@go get github.com/golang/dep
	@go install github.com/golang/dep/cmd/dep
	@dep ensure 	

build:
	@go build -ldflags "-X main.MajorMinorPatch=$(VERSION) \
                        -X main.ReleaseType=$(TYPE) \
                        -X main.GitCommitSha=$(COMMIT)"

install:
	@go install -ldflags "-X main.MajorMinorPatch=$(VERSION) \
                          -X main.ReleaseType=$(TYPE) \
                          -X main.GitCommitSha=$(COMMIT)"

container: gox
	@gox -ldflags "-X main.MajorMinorPatch=$(VERSION) \
                   -X main.ReleaseType=$(TYPE) \
                   -X main.GitCommitSha=$(COMMIT) \
                   -w" \
	     -osarch="linux/amd64" \
	     -output "build/{{.OS}}_{{.Arch}}/$(NAME)"
	
	docker build --rm --pull --tag $(IMAGE):$(TAG) .

tag: container
	docker tag $(IMAGE):$(TAG) $(IMAGE):$(COMMIT)

push: tag
	docker push $(IMAGE):$(COMMIT)
	docker push $(IMAGE):$(TAG)

cross-compile: gox dep
	@rm -rf build/
	@gox -ldflags "-X main.MajorMinorPatch=$(VERSION) \
                   -X main.ReleaseType=$(TYPE) \
                   -X main.GitCommitSha=$(COMMIT) \
                   -w" \
	     -osarch="linux/386" \
	     -osarch="linux/amd64" \
	     -osarch="darwin/amd64" \
	     -output "build/{{.OS}}_{{.Arch}}/$(NAME)"

dist: cross-compile
	$(eval FILES := $(shell ls build))
	@rm -rf dist && mkdir dist
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		(cd $(shell pwd)/dist && shasum -a 512 $$f.tar.gz > $$f.sha512); \
		echo $$f; \
	done

release: dist push
	@latest_tag=$$(git describe --tags `git rev-list --tags --max-count=1`); \
	comparison="$$latest_tag..HEAD"; \
	if [ -z "$$latest_tag" ]; then comparison=""; fi; \
	changelog=$$(git log $$comparison --oneline --no-merges --reverse); \
	github-release samsung-cnct/$(NAME) $(VERSION) "$$(git rev-parse --abbrev-ref HEAD)" "**Changelog**<br/>$$changelog" 'dist/*'; \
	git pull

.PHONY: dep build install container tag push cross-compile dist release  