![Open Banking Logo](https://bitbucket.org/openbankingteam/conformance-suite/raw/99b76db5f60bb4d790d6f32bffae29cbe95a3661/docs/static_files/OBIE_logotype_blue_RGB.PNG)

The **Dynamic Client Registration Conformance Tool** is an Open Source test tool provided by [Open Banking](https://www.openbanking.org.uk/).

The [Dynamic Client Registration](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2) APIs allow TPPs to register one or more clients with a
ASPSPs in a manner that offers very low friction and removes hurdles and barriers to entry. The goal of the DCR Conformance Tool is to allow implementers of DCR to test a interface against the DCR standard.

The supporting documentation assumes technical understanding of the Open Banking ecosystem and DCR. An introduction to the concepts is available via the [Open Banking Website](https://www.openbanking.org.uk/).

## Quickstart

See the guide at [QUICK-START.md](https://bitbucket.org/openbankingteam/conformance-dcr/src/develop/QUICK-START.md).

## Release Notes 
* * *

# Release v1.0.1

The release is called **v1.0.1**, it updates the tool:
- Support report download via http
- Patch to fix 3rd party library bug jwt-go
- Fix content type sent on client register to application/jose
- Added debug file to report zip 

# Development

## Requirements

* Go 1.13

## Build and Run

To run against Ozone:

```sh
git clone git@bitbucket.org:openbankingteam/conformance-dcr.git && cd conformance-dcr && make build && ./dcr -config-path configs/config.json
```
