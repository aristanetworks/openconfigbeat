[![Build Status](https://travis-ci.com/aristanetworks/openconfigbeat.svg?token=4pFw4pt39z2iokPtps7U&branch=master)](https://travis-ci.com/aristanetworks/openconfigbeat)

# Openconfigbeat

Welcome to Openconfigbeat, the Beat for [OpenConfig](http://openconfig.net)-compatible network devices.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/aristanetworks`

Document example:

```json
{
  "@timestamp": "2016-10-06T22:38:51.770Z",
  "beat": {
    "hostname": "beater",
    "name": "openconfigbeat"
  },
  "dataset": "cairo",
  "type": "openconfigbeat",
  "update": {
    "/": {
      "openconfig-bgp:bgp": {
        "global": {
          "config": {
            "as": 55064,
            "router-id": "0.0.0.0"
          },
        },
      },
      "openconfig-vlan:vlans": {
        "vlan": [
          {
            "config": {
              "name": "default",
              "status": "ACTIVE",
              "tpid": "oc-vlan-types:TPID_0x8100",
              "vlan-id": 1
            },
            "state": {
              "status": "ACTIVE",
              "tpid": "oc-vlan-types:TPID_0x8100"
            },
            "vlan-id": 1
          }
        ]
      }
    }
  }
}
```

## Getting Started with Openconfigbeat

### Requirements

* [Golang](https://golang.org/dl/) >= 1.6
* [Glide](https://github.com/Masterminds/glide) >= 0.10.0

### Configuration
Edit `etc/beat.yml` with the appropriate configuration.

### Build

To build the binary for Openconfigbeat run the command below. This will generate a binary
in the same directory with the name openconfigbeat.

```
make
```


### Run

To run Openconfigbeat with debugging output enabled, run:

```
./openconfigbeat -e -d openconfigbeat.go
```


### Package

To be able to package Openconfigbeat the requirements are as follows:

 * [Docker Environment](https://docs.docker.com/engine/installation/) >= 1.10
 * $GOPATH/bin must be part of $PATH: `export PATH=${PATH}:${GOPATH}/bin`

To cross-compile and package Openconfigbeat for all supported platforms, run the following commands:

```
cd dev-tools/packer
make deps
make images
make
```

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/openconfigbeat.template.json, etc/openconfigbeat.asciidoc and openconfigbeat.yml

```
make update
```


### Cleanup

To clean  Openconfigbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Openconfigbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/aristanetworks
cd ${GOPATH}/github.com/aristanetworks
git clone https://github.com/aristanetworks/openconfigbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).
