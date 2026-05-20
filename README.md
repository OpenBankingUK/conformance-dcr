![Open Banking Logo](https://github.com/OpenBankingUK/conformance-suite/blob/develop/docs/static_files/OBIE_logotype_blue_RGB.PNG)

The **Dynamic Client Registration Conformance Tool** is an Open Source test tool provided
by [Open Banking](https://www.openbanking.org.uk/).

The [Dynamic Client Registration](https://openbankinguk.github.io/dcr-docs-pub/v3.4/dynamic-client-registration.html) APIs allow TPPs to register one or more clients with ASPSPs in a manner that offers very low friction and removes hurdles and barriers to entry. The goal of the DCR Conformance Tool is to allow implementers of DCR to test an interface against the DCR standard.

The supporting documentation assumes technical understanding of the Open Banking ecosystem and DCR. An introduction to the concepts is available via the [Open Banking Website](https://www.openbanking.org.uk/).

## Quickstart

See the guide at [QUICK-START.md](https://github.com/OpenBankingUK/conformance-dcr/blob/develop/QUICK-START.md).

### Specification

* Dynamic Client Registration Specifications: <https://openbankinguk.github.io/dcr-docs-pub/>

# Release v1.4.0 (11th May 2026)

[See full list changes](https://github.com/OpenBankingUK/conformance-dcr/blob/develop/releases/v1.4.0.md) (v1.4.0.md)

# Development

## Requirements

* Go 1.26

## Build and Run

### Docker

```
docker run --rm -it -v [CONFIG FILE]:/config.json openbanking/conformance-dcr:v1.4.0 -config-path=/config.json
```

### From source:

```sh
gh repo clone OpenBankingUK/conformance-dcr && cd conformance-dcr && make build && ./dcr -config-path configs/config.json
```
