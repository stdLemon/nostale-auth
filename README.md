# nostale-auth
Library for obtaining code used in login packet (NoS0577)

## How to use
1. Intercept POST request sent to */api/v1/auth/iovation* (it can be done with fiddler)
2. Save value of blackbox field to file
3. Generate identity file `pkg/blackbox/gfclient_poc/create_identity.js blackbox.txt > identity.json`
4. Fill up timing range in *identity.json*
    1. Save few blackboxes from iovation request
    2. Find out what is the minimum and maximum timing for *collectionDurationMs* field

5. Create and fill *account.json* basing on template `mv account_template.json account.json`
6. Example code can be found in [gfclient_test.go](https://github.com/stdLemon/nostale-auth/blob/main/pkg/gfclient/gfclient_test.go)

## Timings example
Assuming that three blackboxes with timings listed below were captured

### 1
```json 
"collectionDurationMs": 548
```

### 2
```json 
"collectionDurationMs": 553
```

### 3
```json 
"collectionDurationMs": 558
```

### Field timing in identity.json should look like this
```json
"timing": {"min": 548, "max": 558}
```

## Compatibility note
`identity.json` generated for older blackbox schema versions is not compatible with blackbox v12.
Regenerate `identity.json` from a fresh captured blackbox before running integration test.
