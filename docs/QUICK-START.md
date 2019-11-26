# `QUICK-START`

Create a configuration file from [/config.json.sample](/config.json.sample). See the sample configuration at [./config_ozone_sample.json](./config_ozone_sample.json), if you are having problems.

## Print (Latest) Help Menu

**NB**: `-help` output is subject to change.

```sh
$ ( \
  DCR_VERSION="latest"; \
  docker --log-level=debug run \
    --rm \
    -it \
    -v "${CONFIG_PATH}":/home/app/.config/conformance-dcr/config.json \
    openbanking/conformance-dcr:"${DCR_VERSION}" \
      -help \
)
...
Dynamic Client Registration Conformance Tool cli
Usage of dcr:
  -config-path string
    	Config file path
  -debug
    	Enable debug defaults to disabled
  -filter string
    	Filter scenarios containing value
  -report
    	Enable report output defaults to disabled
  -version
    	Print the version details of conformance-dcr
...
```

## Print (Latest) Version Information

**NB**: `-version` output is subject to change.

```sh
$ ( \
  DCR_VERSION="latest"; \
  docker --log-level=debug run \
    --rm \
    -it \
    -v "${CONFIG_PATH}":/home/app/.config/conformance-dcr/config.json \
    openbanking/conformance-dcr:"${DCR_VERSION}" \
      -version \
)
...
```

## Run Latest Version

```sh
$ ( \
  DCR_VERSION="latest"; \
  CONFIG_PATH="$(pwd)/configs/config.json"; \
  docker --log-level=debug run \
    --rm \
    -it \
    -v "${CONFIG_PATH}":/home/app/.config/conformance-dcr/config.json \
    openbanking/conformance-dcr:"${DCR_VERSION}" \
      -debug \
      -report \
      -config-path=/home/app/.config/conformance-dcr/config.json \
)
...
=== Scenario: DCR-001 - Validate OIDC Config Registration URL
	Test case: Validate Registration URL
		PASS Registration Endpoint Validate
=== Scenario: DCR-002 - Dynamically create a new software client
	Test case: Register software client
		PASS Generate signed software client claims
2019/11/26 10:45:58 getting claims from authoriser
2019/11/26 10:45:58 setting signed claims in context var: jwt_claims
		PASS Software client register
2019/11/26 10:45:58 get jwt claims from ctx var: jwt_claims
2019/11/26 10:45:58 request:
 POST /dynamic-client-registration/v3.1/register HTTP/1.1
Host: ob19-rs1.o3bank.co.uk:4501
Accept: application/json
Content-Type: application/jwt
...
```

See complete at [./logs/ozone_run_debug.log](./logs/ozone_run_debug.log).

## Run Stable Version

```sh
$ ( \
  DCR_VERSION="v1.0.0"; \
  CONFIG_PATH="$(pwd)/configs/config.json"; \
  docker --log-level=debug run \
    --rm \
    -it \
    -v "${CONFIG_PATH}":/home/app/.config/conformance-dcr/config.json \
    openbanking/conformance-dcr:"${DCR_VERSION}" \
      -debug \
      -report \
      -config-path=/home/app/.config/conformance-dcr/config.json \
)
...
```
