# nostale-auth
Library for obtaining code used in login packet (NoS0577)

## How to use
1. Intercept POST request sent to */api/v1/auth/iovation* (it can be done with fiddler)
2. Save value of blackbox field to file
3. Generate identity file `gfclient_poc/create_identity.js blackbox.txt > identity.json`
4. Fill up timings in *identity.json*
    1. Save few blackboxes from iovation request
    2. Find out what are the minimum and maximum timings for each field (dP, dF, dW and dC)
    3. Timing *d* in blackbox is a sum of all other timings plus additional ms, only added value should be used in *timing.d* field

5. Create and fill *account.json* basing on template `mv account_template.json account.json`
6. Example code can be found in [gf_client_test.go](https://github.com/stdLemon/nostale-auth/blob/main/gf_client_test.go)

## Timings example
Assuming that three blackboxes with timings listed below were captured

### 1
```json 
"dP": 31,
"dF": 151,
"dW": 351,
"dC": 5,
"d": 548
```

### 2
```json 
"dP": 32,
"dF": 152,
"dW": 352,
"dC": 6,
"d": 553
```

### 3
```json 
"dP": 33,
"dF": 153,
"dW": 353,
"dC": 7,
"d": 558
```

### identity timings should look like this
```json
"dP": {"min": 31, "max": 33},
"dF": {"min": 151, "max": 153},
"dW": {"min": 351, "max": 353},
"dC": {"min": 5, "max": 7},
"d": {"min": 10, "max": 12}
```
