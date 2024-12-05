#!/bin/sh
# Copy all executables from ./bin to /usr/local/bin
cp ./bin/linux/nsqadmin /usr/local/bin/
cp ./bin/linux/nsqd /usr/local/bin/
cp ./bin/linux/nsqlookupd /usr/local/bin/
cp ./bin/linux/redis-cli /usr/local/bin/
cp ./bin/linux/redis-server /usr/local/bin/
# Execute the original command
exec "$@"

