# Binary Services

This directory contains binaries to run services used in unit and integration tests. The [test script](../test.sh) controls the starting and stopping of these services.

* NSQ version is 1.2.0
* Redis version is 5.0.7

The NSQ binaries were downloaded from their respective download pages. Redis was built from scratch on OSX Catalina 10.15.2 and on Ubuntu 14.04.6 LTS (GNU/Linux 4.4.0-148-generic x86_64) with gcc version gcc (Ubuntu 4.8.4-2ubuntu1~14.04.4) 4.8.4.

## Test Services

bin/linux/cc-test-reporter is the CodeClimate test reporter. Travis uses it to send coverage reports to CodeClimate. See https://docs.codeclimate.com/docs/configuring-test-coverage for more info.

## Mac OSX Notes

Mac OS may silently refuse to run the services the first time you try. To get
around that, you'll have to right-click each binary and choose *Open* when Mac presents the warning message about the app coming from an unknown developer.

You can then kill the service with Control-C, or by closing the terminal window. After that first run, the [test script](../test.sh) should be able to run the services without your intervention.
