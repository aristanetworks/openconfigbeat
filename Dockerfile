# Copyright (C) 2016  Arista Networks, Inc.
# Use of this source code is governed by the Apache License 2.0
# that can be found in the LICENSE file.

FROM golang:1.7.3

RUN mkdir -p /go/src/github.com/aristanetworks/openconfigbeat
WORKDIR /go/src/github.com/aristanetworks/openconfigbeat
COPY ./ .
RUN make

ENTRYPOINT ["./openconfigbeat", "-e", "-d", "openconfigbeat.go"]
