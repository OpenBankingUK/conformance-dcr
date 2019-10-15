![Open Banking Logo](https://bitbucket.org/openbankingteam/conformance-suite/raw/99b76db5f60bb4d790d6f32bffae29cbe95a3661/docs/static_files/OBIE_logotype_blue_RGB.PNG)

The **Dynamic Client Registration Conformance Tool** is an Open Source test tool provided by [Open Banking](https://www.openbanking.org.uk/). 

The [Dynamic Client Registration](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2) APIs allow TPPs to register one or more clients with a 
ASPSPs in a manner that offers very low friction and removes hurdles and barriers to entry. The goal of the DCR Conformance Tool is to allow implementers of DCR to test a interface against the DCR standard.

The supporting documentation assumes technical understanding of the Open Banking ecosystem and DCR. An introduction to the concepts is available via the [Open Banking Website](https://www.openbanking.org.uk/).

## Running 

### Prerequisites

In order to run a container you'll need docker installed.

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)

### Quickstart

Create a configuration file from `config.json.sample`.

Pull and run the latest (stable) tagged Docker image:

    > docker run --rm -v /path/to/you/local/config.json:/config.json -it "openbanking/conformance-dcr:latest" -config-path=/config.json  

## Development

### Requirements

- Go 1.12

### Running

To run against Ozone:

```bash
make build; ./dcr -config-path configs/config.json
```