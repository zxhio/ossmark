# ossmark

## usage

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

And some optional fields.
```json
{
    "listen_port": 9991,
}
```

In order for **ossmark** to run continuously, a file called **ossmark.service** is provided, and you can execute the following commands
```bash
mv ossmark.service /lib/systemd/system/ossmark.service
systemctl daemon-reload
systemctl enable ossmark
systemctl start ossmark
```

## TODO
- Provide completed installation package for **ossmark**.