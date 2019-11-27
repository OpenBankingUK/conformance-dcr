# Dynamic Client Registration Tool We Quick Start

This guide will assist you with the technical steps required to setup the Dynamic Client Registration Tool and run your first tests.

## Prerequisites

This guide assumes the following tools are installed and functioning correctly. Versions specified used when writing this guide.

* Docker (Client: 18.09.1, Server: 18.09.1 on OSX)
* Valid Certificates

Note for Windows 10 users - Docker on Windows 10 requires Hyper-V to be installed. Hyper-V is only available on Pro or Enterprise versions. Please refer to this guide for more information.

## How to run Dynamic Client Registration Tool

Create a configuration file using [/config.json.sample](/config.json.sample).

The following command will download the latest DCR Tool from docker hub and run it.

```sh
docker run --rm -it -v [CONFIG FILE]:/config.json openbanking/conformance-dcr:[TAG] -config-path=/config.json
```

## Optional - Downloading with Docker Content Trust (recommended)

Docker Content Trust *(DCT)* ensures that all content is securely received and verified. Open Banking cryptographically signs the images upon completion of a satisfactory image check, so that implementers can verify and trust certified content.

To verify the content has not been tampered with you can you the `DOCKER_CONTENT_TRUST` flaG. For example:

    DOCKER_CONTENT_TRUST=1 docker pull openbanking/conformance-dcr:TAG