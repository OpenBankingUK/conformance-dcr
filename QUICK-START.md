# Dynamic Client Registration Tool We Quick Start

This guide will assist you with the technical steps required to setup the Dynamic Client Registration Tool and run your first tests.

## Prerequisites

This guide assumes the following tools are installed and functioning correctly. Versions specified used when writing this guide.

* Docker (Client: 18.09.1, Server: 18.09.1 on OSX)
* Valid Certificates

Note for Windows 10 users - Docker on Windows 10 requires Hyper-V to be installed. Hyper-V is only available on Pro or Enterprise versions. Please refer to this guide for more information.

## How to run Dynamic Client Registration Tool

### Configuration

A template configuration file can be found at [/config.json.sample](/config.json.sample).

|Name                       |Type        |Description                               |
|---------------------------|------------|------------------------------------------|
|wellknown_endpoint         | string     | Open ID Connect wellknown endpoint for DCR|
|ssa                        | string     | Software Statement Assertion for client   |
|kid                        | string     | Key ID - Identifies your key pair|
|aud                        | string     | Audience - The intended audience that the client is being registerd with. Typically the unique identifier of the organisation.|
|issuer                     | string     | Unique identifer for the TPP/Client organisation, for example `software_id` as provided by Open Banking Directory. |
|redirect_uris              | []string   | URIs used to callback to your application during registration, consent acquisition|
|private_key                | string     | Private key associated with client|
|transport_root_cas         | []string   | Root CAs for transport cert|
|transport_cert             | string     | Transport cert associated with client|
|transport_key              | string     | Private key for transport|
|get_implemented            | bool       | HTTP GET method implemented as per DCR specification|
|put_implemented            | bool       | HTTP PUT method implemented as per DCR specification|
|delete_implemented         | bool       | HTTP DELETE method implemented as per DCR specification|
|environment                | string     | Environment where this tool is running against, ex: sandbox|
|brand                      | string     | Brand name|

Note that HTTP POST is the only HTTP method required by the specification, which will always be tested. If the implementation under test supports HTTP GET, PUT or DELETE, they can be specified using the booleans in the configuration as shown above.

### Run the tool

The following command will download the latest DCR Tool from docker hub and run it.

```sh
docker run --rm -it -v [CONFIG FILE]:/config.json openbanking/conformance-dcr:[TAG] -config-path=/config.json
```

## Generate DCR Compliance report

DCR Report is generated when running the tool with a `-report` flag, for security reasons you will have
to download from a embedded webserver.

```sh
docker run --rm -it -p 8080:8080 -v [CONFIG FILE]:/config.json openbanking/conformance-dcr:[TAG] -config-path=/config.json
```

Instructions will be printed how to download the report.

## Optional - Downloading with Docker Content Trust (recommended)

Docker Content Trust *(DCT)* ensures that all content is securely received and verified. Open Banking cryptographically signs the images upon completion of a satisfactory image check, so that implementers can verify and trust certified content.

To verify the content has not been tampered with you can you the `DOCKER_CONTENT_TRUST` flaG. For example:

    DOCKER_CONTENT_TRUST=1 docker pull openbanking/conformance-dcr:TAG