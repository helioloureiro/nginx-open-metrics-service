PROJECT := nginx-openmetrics
RELEASE := 1.0
DEBVERSION := 9
ARCH := amd64

BIN_LINUX := $(PROJECT)
SRC := $(PROJECT)/main.go

CONTACT := helio@loureiro.eng.br


BUILD_OPTIONS = -modcacherw
BUILD_OPTIONS += -ldflags="-w -X 'main.Version=$$(git tag -l --sort taggerdate | tail -1)'"
BUILD_OPTIONS += -buildmode=pie
BUILD_OPTIONS += -tags netgo,osusergo
BUILD_OPTIONS += -trimpath

REGISTRY := nononono

all: $(BIN_LINUX)

init:
	go mod tidy
	go mod vendor

test: init $(SRC)
	cd $(PROJECT)
	go test -v ./...

$(BIN_LINUX): init $(SRC) test
	cd $(PROJECT)
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_LINUX) $(BUILD_OPTIONS) ./...

container:
	docker build -t $(PROJECT) .
	img_id=$(shell docker images --format=json | jq -r '. | select(.Repository == "$(PROJECT)") | .ID ') ;\
	git_tag=$(shell ./vas.sh print_tag) ;\
	docker tag $$img_id $(REGISTRY):$$git_tag
	#docker push $()


tag:
	NEWVERSION=$(shell expr $(DEBVERSION) + 1 ); \
		sed -i "/^DEBVERSION/ s/:= .*/:= $$NEWVERSION/" Makefile; \
	git commit -a -m "Uplifiting to version: $(RELEASE)-$$NEWVERSION"
	make git_tag

git_tag:
	last_commit=$(shell git log --pretty=format:"%H" | head -1); \
	git tag "$(RELEASE)-$(DEBVERSION)" $$last_commit

push:
	git push origin HEAD
	git push origin HEAD --tags

clean:
	rm -f $(PROJECT)/$(BIN_LINUX)
	rm -f *deb *dsc *build *buildinfo *changes *xz
