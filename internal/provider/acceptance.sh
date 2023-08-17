set -eu
set -x

cleanup_docker() {
	docker stop mock-api
}
failed() {
	cleanup_docker
	exit 1
}

command -v docker >/dev/null 2>&1 || { echo >&2 "docker command not installed or in PATH"; exit 1; }
command -v go >/dev/null 2>&1 || { echo >&2 "go command not installed or in PATH"; exit 1; }
command -v make >/dev/null 2>&1 || { echo >&2 "make command not installed or in PATH"; exit 1; }
command -v terraform >/dev/null 2>&1 || test -n "${TF_ACC_TERRAFORM_PATH:-}" || { echo >&2 "terraform command not installed or in PATH, TF_ACC_TERRAFORM_PATH not set"; exit 1; }

docker pull sadokf/searchstax-mock-api

export SEARCHSTAX_HOST=http://localhost:3000/api/rest/v2
export SEARCHSTAX_USERNAME=testUSERNAME
export SEARCHSTAX_PASSWORD=testPWD


docker run -d \
	-p 127.0.0.1:3000:3000 \
	--rm --name mock-api  sadokf/searchstax-mock-api || failed
GO111MODULE=on make testacc TEST=./internal/provider || failed
cleanup_docker

