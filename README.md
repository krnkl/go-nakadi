[![Build Status](https://travis-ci.org/stoewer/go-nakadi.svg?branch=master)](https://travis-ci.org/stoewer/go-nakadi)
[![Coverage Status](https://coveralls.io/repos/github/stoewer/go-nakadi/badge.svg?branch=master)](https://coveralls.io/github/stoewer/go-nakadi?branch=master)
[![GoDoc](https://godoc.org/github.com/stoewer/go-nakadi?status.svg)](https://godoc.org/github.com/stoewer/go-nakadi)
---

go-nakadi
=========

Is a client for the [Nakadi](https://zalando.github.io/nakadi/manual.html) event broker written in and
for Go. the package provides convenient access to many features of Nakadi's API.

The package can be used to manage event type definitions. The EventAPI can be used to inspect existing
event types and allows further to create new event types and to update existing ones. The SubscriptionAPI
provides subscription management: existing subscriptions can be fetched from Nakadi and of course it is
also possible to create new ones. The PublishAPI of this library is used to broadcast event types of
all event type categories via Nakadi. Last but not least, the package also implements a StreamAPI, which
enables event processing on top of Nakadi's subscription based high level API.

To make the communication with Nakadi more resilient all sub APIs of this package can be configured
to retry failed requests using an exponential back-off algorithm. Please consult the
[package documentation](https://godoc.org/github.com/stoewer/go-nakadi) for further details.

Versions and stability
----------------------

This package can be considered stable and ready to use. All releases follow the rules of
[semantic versioning](http://semver.org).

Although the master branch is supposed to remain stable, there is not guarantee that braking changes will not
be merged into master when major versions are released. Therefore the repository contains version tags in
order to support vendoring tools such as [glide](https://glide.sh). The tag names follow common conventions
and have the following format `v1.0.0`.

Dependencies
------------

Build dependencies

* github.com/cenkalti/backoff
* github.com/pkg/errors

Test dependencies

* github.com/stretchr/testify
* gopkg.in/jarcoal/httpmock.v1

Run unit and integration tests
------------------------------

In oder to run the unit and integration tests all of the above dependencies must be installed. Furthermore
these tests require a [running Nakadi instance](https://zalando.github.io/nakadi/manual.html#getting-started)
on the local computer.

To run all tests invoke the following command within the `go-nakadi` root directory:

```
go test -tags=integration .
```

License
-------

This project is open source an published under the [MIT license](LICENSE).



test commit changes
