# Code Golf

This is the repository behind https://code.golf

## Quickstart

1. Install mkcert:
```
$ yay mkcert
```

2. Install the local CA:
```
$ make cert
Using the local CA at "~/.local/share/mkcert" ✨
The local CA is now installed in the system trust store! ⚡️
The local CA is now installed in the Firefox and/or Chrome/Chromium trust store (requires browser restart)! 🦊


Created a new certificate valid for the following names 📜
 - "localhost"

The certificate is at "./localhost.pem" and the key at "./localhost-key.pem" ✅
```

3. Build the assets:
```
$ ./build-assets
```

4. Bring up the website:
```
$ make dev
```

5. Optionally, load information from the code.golf database.
```
pip install -r utils/requirements.txt
utils/update_sql_from_api.py
```

6. Navigate to https://localhost

## Hacking

Some of ancillary scripts are written in Raku, to run these ensure you have a
recent install of Raku installed and use Zef to install the dependencies:
```
$ zef install --deps-only .
```

## Testing

1. Run the unit tests:
```
$ go test ./...
```

2. Run the e2e tests:
```
$ prove
```
