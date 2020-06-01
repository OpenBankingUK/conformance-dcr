### Release v1.2.0 (13th May 2020)

New scenarios with wrong response types, added organisation details to report.

[See full list changes](releases/v1.2.0.md) 

### Release v1.1.0 (22nd April 2020)

- New scenario for updating software client
- New scenario for updating software client with wrong id
- New scenario for retrieving software client with wrong id
- Token endpoint sign method now comes from wellknown and limited to PS256
- Added missing license
- Limit `response_type` to `code` and/or `code id_token` if more present in the .wellknown endpoint

### Release v1.0.4

- Increase timeout to 10 seconds on http clients to help slower endpoints

### Release v1.0.3

- Fixed `response_types` property in register software from static to `response_types_supported` from .wellknown    

### Release v1.0.2

- Removed unused RS256 flag
- Fixed `request_object_signing_alg` claims value from `none` to first found in .wellknown  
- Fixed missing `scope` in client credentials grant call  
- Fixed wrong header token value calculation for `client_secret_basic` token endpoint auth method   

### Release v1.0.1

- Support report download via http
- Patch to fix 3rd party library bug jwt-go
- Fix content type sent on client register to application/jose
- Added debug file to report zip 
