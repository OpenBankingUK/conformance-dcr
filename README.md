![Open Banking Logo](https://github.com/OpenBankingUK/conformance-suite/blob/develop/docs/static_files/OBIE_logotype_blue_RGB.PNG)

The **Dynamic Client Registration Conformance Tool** is an Open Source test tool provided
by [Open Banking](https://www.openbanking.org.uk/).

The [Dynamic Client Registration](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2)
APIs allow TPPs to register one or more clients with a ASPSPs in a manner that offers very low friction and removes
hurdles and barriers to entry. The goal of the DCR Conformance Tool is to allow implementers of DCR to test an interface
against the DCR standard.

The supporting documentation assumes technical understanding of the Open Banking ecosystem and DCR. An introduction to
the concepts is available via the [Open Banking Website](https://www.openbanking.org.uk/).

## Quickstart

See the guide at [QUICK-START.md](https://github.com/OpenBankingUK/conformance-dcr/blob/develop/QUICK-START.md).

### Specification

* Dynamic Client Registration Specifications: <https://openbankinguk.github.io/dcr-docs-pub/>

# Release v1.3.1 (6th December 2021)

[See full list changes](https://github.com/OpenBankingUK/conformance-dcr/blob/develop/releases/v1.3.1.md) (v1.3.1.md)

# Development

## Requirements

* Go 1.17

## Build and Run

To run against Ozone:

```sh
git clone git@bitbucket.org:openbankingteam/conformance-dcr.git && cd conformance-dcr && make build && ./dcr -config-path configs/config.json
```
