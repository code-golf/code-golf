# Code Golf

This is the repository behind https://code.golf

## Quickstart

1. Install dependencies:
* [Docker](https://docs.docker.com/engine/install/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [Go](https://golang.org/doc/install)
* [make](https://www.gnu.org/software/make/) - most likely already on your system
* [mkcert](https://github.com/FiloSottile/mkcert#installation)
* [npm](https://www.npmjs.com/get-npm)

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
go run utils/update_sql_from_api.go
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
$ make test
```

2. Run the e2e tests:
```
$ make e2e
```

## Style

URL slugs are consistently abbreviated (e.g. cheeovs, langs, stats) but page
titles aren't (e.g. Achievements, Languages, Statistics).

Paginated URLs use a trailing number but only on pages after the first (e.g.
/rankings/medals/all, /rankings/medals/all/2, etc.).
