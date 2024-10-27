export GO_VERSION := "1.22.3"
offergen_tester_image := "offergen-tester-image"
export OFFERGEN_TESTER_IMAGE := offergen_tester_image

src_mount := "-v " + justfile_directory() / "offergen" + ":/source"
docker_run := "docker run --rm"


build:
    docker build \
        --build-arg GO_VERSION=${GO_VERSION} \
        -t {{offergen_tester_image}} .

update-modules: build
    {{docker_run}} -it --workdir /source {{src_mount}} golang:${GO_VERSION}-alpine go mod tidy

generate: build
    {{docker_run}} {{src_mount}} --entrypoint templ {{offergen_tester_image}} generate

generate-watch: build
    {{docker_run}} -it {{src_mount}} --entrypoint templ {{offergen_tester_image}} generate -watch

fmt: build
    {{docker_run}} -it {{src_mount}} --entrypoint bash {{offergen_tester_image}} -c "gofmt -w . && templ fmt templates"

shell: build
    {{docker_run}} -p 8080:80 -it {{src_mount}} --entrypoint /bin/bash {{offergen_tester_image}}


shell-compose: build
    docker compose run --entrypoint /bin/bash offergen

test: build
    docker compose run --rm --entrypoint go offergen test ./... -v

sql:
    docker compose exec postgres psql testdb -U testuser

up: build
    docker compose up -d

logs:
    docker compose logs

down:
    docker compose down --remove-orphans --timeout=0 --volumes


build-prod: generate
    {{docker_run}} {{src_mount}} --entrypoint /bin/bash {{offergen_tester_image}} -c 'export GOOS=linux && export GOARCH=amd64 && go build -o build/'
