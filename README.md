# doicache

Keep a local cache of DOI API responses.

```shell
$ doicache 10.1103/PhysRevLett.118.140402
http://link.aps.org/doi/10.1103/PhysRevLett.118.140402
```

Dump all keys:

```
$ doicache -dk
10.1103/PhysRevLett.118.140402
```

Adjust expiration date:

```
$ doicache -verbose -ttl 1s 10.1103/PhysRevLett.118.140402
INFO[0000] entry expired
INFO[0000] https://doi.org/api/handles/10.1103/PhysRevLett.118.140402
INFO[0001] {"Date":"2018-05-25T01:19:02.177003048+02:00","Blob":"eyJyZ..."}
http://link.aps.org/doi/10.1103/PhysRevLett.118.140402
```

Read input from a file:

```
$ doicache < file
```

----

API docs: https://www.doi.org/factsheets/DOIProxy.html#rest-api - an example
response:

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

