# Code Golf

This is the repository behind the website and community https://code.golf

## Quickstart

Alternatively, see the Vagrant section below to create a virtual machine with everything pre-installed.

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

3. Install the NPM packages:

> *NOTE*: if your host OS is not the same architecture / executable format as
> your Docker environment, this can result in incorrect format binaries installed
> into `node_modules`, so you may want to skip this step.

```
$ npm install
```

4. Bring up the website:
```
$ make dev
```

5. Optionally, load information from the code.golf database.
```
$ go run utils/update_sql_from_api.go
```

6. Navigate to https://localhost

## Hacking

Some of ancillary scripts are written in Raku, to run these ensure you have a
recent install of Raku installed and use Zef to install the dependencies:
```
$ zef install --deps-only .
```

## TypeScript

The `js/` directory contains the TypeScript files which will be transpiled by
`esbuild` into JavaScript files for serving. `.tsx` files can additionally
make use of [JSX](https://www.typescriptlang.org/docs/handbook/jsx.html).

## Linting

Run `make lint` to lint the code before a pull request. This lints the TypeScript code, then the Go code.

In Visual Studio Code, the following settings are helpful for editor support for ESLint:

```
"eslint.validate": ["typescript", "typescriptreact"],
"eslint.format.enable": true,
"editor.defaultFormatter": "dbaeumer.vscode-eslint"
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

## Vagrant

Using Vagrant and VirtualBox, with one command you can create a virtual machine, install all of the other dependencies onto it (Docker, Docker Compose, Go, Raku + zef dependencies, npm, etc.), and forward port 443 (https) to it.

1. Install the dependencies:
* [Vagrant](https://www.vagrantup.com/downloads)
* [VirtualBox](https://www.virtualbox.org/)

If you have homebrew, you can install with:
```
$ brew install vagrant virtualbox
```

2. Create the virtual machine:
```
$ vagrant up
```

3. Install certificates. Ideally, you should run the following command outside of the virtual machine, the same as in the quickstart, because it will install the certificate in your browser. Alternatively, you can run the command in a virtual machine shell and then install the certificates manually.
```
$ make cert
```
or
```
$ vagrant ssh --command 'cd /vagrant/ && make cert'
```

5. Bring up the website:
```
$ vagrant ssh --command 'cd /vagrant/ && make dev'
```

5. Optionally, load information from the code.golf database.
```
$ vagrant ssh --command 'cd /vagrant/ && go run utils/update_sql_from_api.go'
```

6. Navigate to https://localhost

## Database Access

If you have [PostgreSQL](https://www.postgresql.org/download/) installed You can access the SQL database directly with the following command. If you are using Vagrant, PostgreSQL is pre-installed.
```
$ make db-dev
```

## API

Validate API definition with [vacuum](https://api.quobix.com/report?url=https://code.golf/api).

## Style

URL slugs are consistently abbreviated (e.g. cheevos, langs, stats) but page
titles aren't (e.g. Achievements, Languages, Statistics).

Paginated URLs use a trailing number but only on pages after the first (e.g.
/rankings/medals/all, /rankings/medals/all/2, etc.).

## arm64/aarch64 Architecture (Ex: Apple Silicon)

If your machine has arm64 architecture, then there are a few options for running a local instance of the site.
You can attempt to use QEMU emulation which is supported by Docker Desktop, but may be slower.
The default when you `make dev` is that the docker images will be built with arm64 architecture.

You may be able to use some languages without building them locally.
Here is a list of some of the languages that are able to run the sample code, without rebuilding individual docker containers:
AWK
Bash
Berry
Lua
Perl
PHP
Python
Raku
Ruby
sed
SQL
Wren

For other languages, you will need to build them locally. For example:
```
$ ./build-langs --no-push C
```

Otherwise, you may see the following error:
```
assertion failed [!result.is_error]: Unable to open /proc/sys/vm/mmap_min_addr
(VMAllocationTracker.cpp:281 init)
 signal: trace/breakpoint trap
 ```
