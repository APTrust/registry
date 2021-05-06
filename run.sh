#!/bin/bash
#
# run.sh
#
# Run registry unit tests or run the registry app on port 8080.
# This starts NSQ and Redis as well, since some tests and app
# features require them.
#

# ----------------------------------------------------------------------
#
# Make sure we got a valid command line arg
#
# ----------------------------------------------------------------------
if [[ $1 != "tests" && $1 != "server" ]]; then
    echo "Usage:"
    echo "./run.sh tests     - to run unit tests"
    echo "./run.sh server  - to run application on port 8080"
    exit 0
fi

# ----------------------------------------------------------------------
#
# Set the path to our local bin dir based on OS type
#
# ----------------------------------------------------------------------
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    DIR="./bin/linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    DIR="./bin/osx"
fi

# ----------------------------------------------------------------------
#
# Start the services: NSQ and Redis
#
# ----------------------------------------------------------------------
echo "Starting NSQ"
eval "$DIR/nsqd --data-path=$TMPDIR > /dev/null 2>&1 &"
NSQ_PID=$!

echo "Starting Redis"
eval "$DIR/redis-server --save "" --appendonly no > /dev/null 2>&1 &"
REDIS_PID=$!


if [[ $TRAVIS == "true" ]]; then
    APT_ENV=travis
elif [[ -z "${APT_ENV}" ]]; then
    APT_ENV=test
fi
echo "APT_ENV is $APT_ENV"

# ----------------------------------------------------------------------
#
# Run tests or server, based on command line arg.
#
# ----------------------------------------------------------------------
if [[ $1 == "tests" ]]; then
    echo "Running registry tests..."
    APT_ENV=$APT_ENV go test -p 1 ./...
elif [[ $1 == "server" ]]; then
    echo "Starting registry app..."
    APT_ENV=$APT_ENV go run registry.go
fi

EXIT_CODE=$?


# ----------------------------------------------------------------------
#
# Shut everything down after tests complete, or after user
# hits Control-C to stop server.
#
# ----------------------------------------------------------------------

echo "Killing NSQ pid $NSQ_PID"
kill $NSQ_PID

echo "Killing Redis pid $REDIS_PID"
kill $REDIS_PID

#rm *.dat

echo "Finished with code $EXIT_CODE"
exit $EXIT_CODE
