# Copyright (C) 2016  Arista Networks, Inc.
# Use of this source code is governed by the Apache License 2.0
# that can be found in the LICENSE file.

FROM alpine

RUN apk add --no-cache curl jq wget \
  && export TAG=$(curl -s \
    https://api.github.com/repos/aristanetworks/openconfigbeat/releases/latest \
    | jq -r '.tag_name' ) \
    && wget https://github.com/aristanetworks/openconfigbeat/releases/download/$TAG/openconfigbeat \
    -O /usr/bin/openconfigbeat \
    && chmod +x /usr/bin/openconfigbeat

ENTRYPOINT ["/usr/bin/openconfigbeat"]
