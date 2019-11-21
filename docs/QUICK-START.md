# `QUICK-START`

## Run Stable Version

```sh
$ curl --proto '=https' --tlsv1.2 -sSf https://bitbucket.org/openbankingteam/conformance-dcr/raw/develop/scripts/docker-run-stable.sh | bash
...
```

## Run Latest Version

```sh
$ curl --proto '=https' --tlsv1.2 -sSf https://bitbucket.org/openbankingteam/conformance-dcr/raw/develop/scripts/docker-run-latest.sh | bash
...
```

## Print (Latest) Help Menu

```sh
$ DCR_VERSION='latest'; docker --log-level=debug run --rm -it openbanking/conformance-dcr:"${DCR_VERSION}"
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

```sh
$ DCR_VERSION='latest'; docker --log-level=debug run --rm -it openbanking/conformance-dcr:"${DCR_VERSION}" -version
...
```