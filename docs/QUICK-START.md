# `QUICK-START`

## Print Help Menu

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

## Print Version Information

```sh
$ DCR_VERSION='latest'; docker --log-level=debug run --rm -it openbanking/conformance-dcr:"${DCR_VERSION}" -version
...
```