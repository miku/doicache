# doicache

Keep a local cache of DOI API responses.

```shell
$ doicache 10.1103/PhysRevLett.118.140402
http://link.aps.org/doi/10.1103/PhysRevLett.118.140402
```

## Usage

```shell
Usage of doicache:
  -db string
        leveldb directory (default "/tmp/.doicache/default")
  -dk
        dump keys
  -dkv
        dump keys and redirects
  -ttl duration
        entry expiration (default 5760h0m0s)
  -verbose
        be verbose
  -version
        show version
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

Example:

```
$ doicache < fixtures/10 | column -t
OK    10.2307/2546078                          https://www.jstor.org/stable/2546078?origin=crossref
OK    10.9783/9780812207729.91                 http://www.degruyter.com/view/books/9780812207729/9780812207729.91/9780812207729.91.xml
OK    10.1590/S0100-40422009000900046          http://www.scielo.br/scielo.php?script=sci_arttext&pid=S0100-40422009000900046&lng=pt&nrm=iso&tlng=pt
OK    10.1097/00043764-199710000-00015         https://insights.ovid.com/crossref?an=00043764-199710000-00015
H404  10.1016/jpaid.2003.07.001                NOTAVAILABLE
OK    10.1093/acrefore/9780199381135.013.205   http://classics.oxfordre.com/view/10.1093/acrefore/9780199381135.001.0001/acrefore-9780199381135-e-205
OK    10.1037/h0050516                         http://content.apa.org/journals/ccp/17/3/232b
OK    10.1016/j.avb.2016.06.006                http://linkinghub.elsevier.com/retrieve/pii/S1359178916300684
OK    10.4028/www.scientific.net/amm.29-32.61  https://www.scientific.net/AMM.29-32.61
OK    10.1136/bmj.2.1493.309                   http://www.bmj.com/cgi/doi/10.1136/bmj.2.1493.309
```

Status codes:

* OK
* H404 (invalid DOI)
* EURL (invalid URL)

## Limitation

Via LevelDB, only one process can access the cache at a time.

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
