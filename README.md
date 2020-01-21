# Code Golf

[![Travis](https://travis-ci.org/code-golf/code-golf.svg)](https://travis-ci.org/code-golf/code-golf) [![Coveralls](https://coveralls.io/repos/github/code-golf/code-golf/badge.svg)](https://coveralls.io/github/code-golf/code-golf)

This is the repository behind https://code-golf.io

## Quickstart

1. Install mkcert:
```
$ yay mkcert
```

2. Install the local CA:
```
$ mkcert -install localhost
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

4. Build the languages:
```
$ ./build-langs
```

5. Bring up the website:
```
$ make dev
```

6. Navigate to https://localhost

## Testing

1. Run the unit tests:
```
$ go test ./...
```

2. Run the e2e tests:
```
$ prove
```

## Thanks

[![BrowserStack](browserstack.png)](https://www.browserstack.com)

Cross browser testing powered by [BrowserStack](https://www.browserstack.com)
