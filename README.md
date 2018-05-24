# doicheck

Check DOI values for validity via doi.org REST API:

* https://www.doi.org/factsheets/DOIProxy.html#rest-api

An example API response:

```
{
  "responseCode": 1,
  "handle": "10.1103/PhysRevLett.118.140402",
  "values": [
    {
      "index": 1,
      "type": "URL",
      "data": {
        "format": "string",
        "value": "http://link.aps.org/doi/10.1103/PhysRevLett.118.140402"
      },
      "ttl": 86400,
      "timestamp": "2017-04-06T02:10:03Z"
    },
    {
      "index": 700050,
      "type": "700050",
      "data": {
        "format": "string",
        "value": "20170405220855"
      },
      "ttl": 86400,
      "timestamp": "2017-04-06T02:10:03Z"
    },
    {
      "index": 100,
      "type": "HS_ADMIN",
      "data": {
        "format": "admin",
        "value": {
          "handle": "0.na/10.1103",
          "index": 200,
          "permissions": "111111110010"
        }
      },
      "ttl": 86400,
      "timestamp": "2017-04-06T02:10:03Z"
    }
  ]
}
```

```
$ doicheck 10.1103/PhysRevLett.118.140402
```

Check all entries in a file:

```
$ doicheck -f list.csv
```

Caches all responses in a local sqlite3 database under
`~/.config/doicheck/doi.db` - where the blob and the timestamp is recorded.

By default `doicheck` will query the local database first with some default
expiration time. To force a database update, use `-force`.

## Implementation

Test [dacap](https://github.com/ubleipzig/dacap) first.

Create a struct:

```
$ curl -s doi.org/api/handles/10.1103/PhysRevLett.118.140402 | jq . | JSONGen
type _ struct {
    Handle       string `json:"handle"`
    ResponseCode int64  `json:"responseCode"`
    Values       []struct {
        Data      interface{} `json:"data"`
        Index     int64       `json:"index"`
        Timestamp string      `json:"timestamp"`
        Ttl       int64       `json:"ttl"`
        Type      string      `json:"type"`
    } `json:"values"`
}
```
