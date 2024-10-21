# ossmark

## usage

Binary **ossmark-\*** show move to */usr/local/bin/*
```shell
$ mv ossmark /usr/local/bin/ossmark
```

### bucket config
Command flags `--access_key_id`, `--access_key_secret`, `--bucket_name` and `--bucket_location` are used to specify the basic config of the bucket.

Alternatively, you can use `--config` to specify the configuration file path.
```json
{
    "access_key_id": "<your access_key_id>",
    "access_key_secret": "<your access_key_secret>",
    "bucket_name": "<your bucket_name>",
    "bucket_location": "<your bucket_location>"
}
```

```shell
$ ossmark-* ./ossmark-sync --access_key_id <access_key_id>  --access_key_secret <access_key_secret>  --bucket_name <bucket_name> --bucket_location <bucket_location>

$ ossmark-* --config config.json
```

### sync
Use `--work_dir` flag to specify the working directory, where the synchronized files are stored.

Use `--mode` flag to specify sync mode
- `time` default, base on modified time.
- `local` directly using local file cover oss.
- `remote` directly using oss cover local file.

```shell
$ ossmark-sync --config <dir>/config.json --mode [mode] --work_dir=./ossdata
```

Similarly, you can use configuration files to replace command-line flags.
```json
{
    "work_dir": "<your work_dir>",
    "mode": "<your mode>"
}
```

### article server
Use `--listen_port` flag to specify the server listen port.

Similarly, you can use configuration files to replace command-line flags.
```json
{
    "listen_port": 9346,
    "log_path": "<your log_path>"
}
```

In order for **ossmark-article** to run continuously, a file called *ossmark-article.service* is provided, and you can execute the following commands
```bash
mv ossmark-article.service /lib/systemd/system/ossmark-article.service
systemctl daemon-reload
systemctl enable ossmark-article
systemctl start ossmark-article
```

## TODO
- Provide completed installation package for **ossmark**.