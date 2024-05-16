# ossmark

## usage

### auth
Binary **ossmark** show move to */usr/local/bin*
```bash
mv ossmark /usr/local/bin/ossmark
```

The default configuration file for Osmark is **/etc/ossmark.json** .

Here are some necessary configuration fields.
```json
{
    "access_key_id": "<your ak>",
    "access_key_secret": "<your sk>",
    "bucket_name": "<your bucket name>",
}
```

### sync
Work dir for sync tool.
```json
{
    "work_dir": "~/.ossmark",
}
```

Add `--sync` flag for command tool, content check mode is
- `time` default, base on modified time.
- `local` directly using local file cover oss.
- `remote` directly using oss cover local file.

```bash
./ossmark --conf <dir>/ossmark.json --sync [mode]
```

### article server
Optional listen port for server.
```json
{
    "listen_port": 9991,
}
```

In order for **ossmark** to run continuously, a file called *ossmark.service* is provided, and you can execute the following commands
```bash
mv ossmark.service /lib/systemd/system/ossmark.service
systemctl daemon-reload
systemctl enable ossmark
systemctl start ossmark
```

## TODO
- Provide completed installation package for **ossmark**.