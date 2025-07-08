# build_for_docker.sh
#
# This builds registry for Docker, injecting the commit ID
# and build date into the binary, so we can display them
# in the Registry footer.

OUTPUT_FILE=${1:-"./registry-exe"}

GIT_REV="-X 'github.com/APTrust/registry/common.CommitID=$(git rev-parse HEAD)'"
BUILD_DATE="-X 'github.com/APTrust/registry/common.BuildDate=$(date)'"

go build -ldflags="$GIT_REV $BUILD_DATE" -o $OUTPUT_FILE registry.go
