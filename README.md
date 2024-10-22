# ossmark

Ossmark contains the following components internally:
- **ossmark-sync**, oss bucket sync tool
- **ossmark-article**, provides HTTP service for displaying markdown files in oss.

## usage

Before install, you should edit *ossmark/conf/ossmark.json* for config, then run `bash ossmark.sh install` .

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
$ ossmark-* --access_key_id <access_key_id> --access_key_secret <access_key_secret> --bucket_name <bucket_name> --bucket_location <bucket_location>

$ ossmark-* --config ossmark.json
```

### ossmark-sync
Use `--work_dir` flag to specify the working directory, where the synchronized files are stored.

Use `--mode` flag to specify sync mode
- `time` default, base on modified time.
- `local` directly using local file cover oss.
- `remote` directly using oss cover local file.

```shell
$ ossmark-sync --config ossmark.json --mode <mode> --work_dir=<work_dir>
```

Alternatively, you can add these fields to the configuration file instead of command-line flags.
```json
{
    "work_dir": "<your work_dir>",
    "mode": "<your mode>"
}
```

### ossmark-article
Use `--listen_port` flag to specify the server listen port.

Similarly, you can add these fields to the configuration file instead of command-line flags.
```json
{
    "listen_port": 9346,
    "log_path": "<your log_path>"
}
```

## TODO
- Add **ossmark-image** to show a thumbnail of oss.