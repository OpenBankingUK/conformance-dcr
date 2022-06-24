# Dynamic Client Registration Tool Quick Start

This guide will assist you with the technical steps required to setup the Dynamic Client Registration Tool and run your
first tests.

## Prerequisites

This guide assumes the following tools are installed and functioning correctly.

Versions specified used when writing this guide.

* Docker (Client: 20.10.7, Server: 20.10.7 on Ubuntu 20.04 LTS)
* Valid Certificates

Note for Windows 10 users - Docker on Windows 10 requires Hyper-V to be installed. Hyper-V is only available on Pro or
Enterprise versions. Please refer to this guide for more information.

## How to run Dynamic Client Registration Tool

### Configuration

A template configuration file can be found at:
- [/config.json.sample](/config.json.sample) for the single SSA
- [/config.json.ssas.sample](/config.json.ssas.sample) for the multiple SSAs.

|Name                       |Type        |Description                                     |
|---------------------------|------------|------------------------------------------------|
|spec_version               | string     | Specification version                          |
|wellknown_endpoint         | string     | Open ID Connect `.well-known` endpoint for DCR |
|ssa                        | string     | Software Statement Assertion for client        |
|ssas                       | []string   | Software Statement Assertions for client (list)|
|kid                        | string     | Key ID - Identifies your key pair              |
|aud                        | string     | Audience - The intended audience that the client is being registered with. Typically the unique identifier of the organisation.|
|redirect_uris              | []string   | URIs used to callback to your application during registration, consent acquisition|
|issuer                     | string     | Unique identifier for the TPP/Client organisation, for example `software_id` as provided by Open Banking Directory. |
|private_key                | string     | Private key associated with client|
|transport_root_cas         | []string   | Root CAs for transport cert|
|transport_cert             | string     | Transport cert associated with client|
|transport_cert_subject_dn  | string     | Transport cert Subject DN associated with client - use when DCR implementation has strict checks and current implementation provides unexpected results |
|transport_key              | string     | Private key for transport|
|get_implemented            | bool       | HTTP GET method implemented as per DCR specification? |
|put_implemented            | bool       | HTTP PUT method implemented as per DCR specification? |
|delete_implemented         | bool       | HTTP DELETE method implemented as per DCR specification? |
|environment                | string     | Environment where this tool is running against, ex: sandbox or production|
|brand                      | string     | Brand name|


Sample json config (*Note* The json5 format with comments, see [/config.json.sample](/config.json.sample) for pure json sample).
```json5
{
"spec_version": "3.3", // or 3.2 
"wellknown_endpoint": "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration",
"ssa": "ex: eyJhbGciOiJQU...", // Signed JWT, Can be generated from the Open Banking DFI
"kid": "ex: lQA1TI94KVbS55vz2IHv4ifc8IA", //  Signing certificate key id
"aud": "aud", // Usually this is an OB Org Id of the ASPSP you are running DCR against
"redirect_uris": ["https://redirect-as-defined-in-the-software-statement.com"], // As configured in the Software Statement
"issuer": "ex: A67kE8qMNgz0F36clmFWbg", // Software Statement Id as defined in OB Directory
"private_key": "ex: MIIEogIBAAKCAQEAj1chaA0Hx9...", // Private key that matches the signing certificate identified by `kid` above  
"transport_root_cas": ["cert 1", "cert 2"], // Certificate chain for Transport certificate, used to validate TLS connection
"transport_cert": "ex: MIIEdTCCA12gAwIBAgIJA5N", // PEM
"transport_key": "transport key", //PEM
"transport_cert_subject_dn": "", //optional, used when standard Subject DN extraction is not returning expected string
"get_implemented": true,
"put_implemented": true,
"delete_implemented": true,
"environment": "sandbox",
"brand": "Brand/product"
}
```

**Note** that HTTP `POST` is the *only* HTTP method required by the specification, which will always be tested.

If the implementation under test supports HTTP `GET`, `PUT` or `DELETE`, they can be specified using the booleans in the
configuration as shown above.

### Run the tool

The following command will download the latest DCR Tool from docker hub and run it.

```sh
docker run --rm -it -v [CONFIG FILE]:/config.json openbanking/conformance-dcr:[TAG] -config-path=/config.json
```

Where
- `[CONFIG FILE]` defines the local, absolute path to the prepared `config.json` file.
- `[TAG]` is a tagged version of the tool
  from [DockerHub](https://hub.docker.com/r/openbanking/conformance-dcr/tags?page=1&ordering=last_updated).

## Generate DCR Compliance report

DCR Report is generated when running the tool with a `-report` flag, for security reasons you will have to download from
an embedded webserver.

```sh
docker run --rm -it -p 127.0.0.1:8080:8080 -v [CONFIG FILE]:/config.json openbanking/conformance-dcr:[TAG] -config-path=/config.json
```

Instructions will be printed how to download the report.

## Optional - Downloading with Docker Content Trust (recommended)

Docker Content Trust *(DCT)* ensures that all content is received securely and verified. Open Banking cryptographically
signs the images upon completion of a satisfactory image check, so that implementers can verify and trust certified
content.

To verify the content has not been tampered with you can you the `DOCKER_CONTENT_TRUST` flaG. For example:

    DOCKER_CONTENT_TRUST=1 docker pull openbanking/conformance-dcr:TAG