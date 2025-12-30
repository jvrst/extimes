# extimes

Explain Unix timestamps; pipe stdin into extimes & stdout appends a human readable timestamp.


I.e. see this example:
```sh
echo "{"time": 1767135589}" | extimes

{time: 1767135589 (2025-12-30T23:59:49+01:00)}
```

