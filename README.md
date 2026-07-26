![certificateManager](images/cm_banner2.png)

# certificateManager

A Go tool to generate, sign, verify and revoke all of your SSL/TLS certificates from your own PKI.

The binary is installed as `cm`.

## Overview

`certificateManager` uses Go's `crypto/x509` package to build and operate a self-contained
Public Key Infrastructure (PKI). With it you can:

- Create your own custom PKI and manage every certificate under it
- Create a root CA (Certificate Authority) certificate
- Create and sign "standard" SSL/TLS server certificates
- Verify certificates
- Revoke certificates
- Optionally export server certificates to Java Keystore (`.jks`) and PKCS#12 (`.p12`) formats

### What this tool does *not* do

- It does not sign certificates against a **remote** CA
- No CRL (Certificate Revocation List) is implemented
- No CDP (CRL Distribution Point) is implemented
- No operation of any kind against a remote CA

This software is intended to run on an **internal network**. If you want your root CA
to be publicly trusted, it is up to you to publish/deploy it with a public authority.

## Concepts

### Environments

In most cases a PKI has a single root CA (or a root CA plus an intermediate CA). You may,
however, want to manage several independent PKIs — for example inside a Docker container.
To make this possible, the tool introduces the idea of an **environment**.

An environment is a small JSON file that describes the directory layout of one PKI. Each
environment is a self-contained sandbox, which lets you manage several PKIs side by side.

Environment files live in `$HOME/.config/JFG/certificatemanager/`. An environment file
looks like this:

```json
{
  "CertificateRootDir": "/localrep/devops/_certificates/famillegratton",
  "RootCAdir": "/localrep/devops/_certificates/famillegratton/CA",
  "ServerCertsDir": "/localrep/devops/_certificates/famillegratton/srv",
  "CertificatesConfigDir": "/localrep/devops/_certificates/famillegratton/cfg",
  "RemoveDuplicates": true
}
```

- `CertificateRootDir` — absolute path to the root of the PKI
- `RootCAdir` — directory holding the root CA certificate, its key and the PKI database
- `ServerCertsDir` — directory holding the server certificates, keys, CSRs and Java exports
- `CertificatesConfigDir` — directory holding the certificate **config** files
- `RemoveDuplicates` — when `true`, refuses to create a certificate whose subject already
  exists (as a valid entry) in the PKI database (`index.txt`)

When you create an environment with `cm env add`, you enter a directory **name** for each
of the sub-directories and the tool stores them as absolute paths, rooted at
`CertificateRootDir`.

You select an environment with the global `-e` flag. If you omit it, the tool uses
`defaultEnv.json`. The `.json` extension is always implied, so `-e test` and
`-e test.json` are equivalent.

### Certificates, root certificates and config files

- A **certificate** (`.crt`) is the actual x509 SSL/TLS file you deploy on a server.
- The **root certificate** (root CA) is the certificate that signs (validates) every other
  certificate in the PKI.
- A **certificate config file** is the JSON file this tool reads to generate a certificate.

A typical certificate config file:

```json
{
  "Country": "CA",
  "Province": "Quebec",
  "Locality": "Blainville",
  "Organization": "famillegratton",
  "OrganizationalUnit": "Computers",
  "CommonName": "bergen.famillegratton.net",
  "IsCA": false,
  "EmailAddresses": [
    "certs@famillegratton.net",
    "jfgratton@famillegratton.net"
  ],
  "Duration": 3,
  "KeyUsage": [
    "digital signature",
    "key encipherment",
    "data encipherment",
    "key agreement"
  ],
  "DNSNames": [
    "bergen",
    "bergen.famillegratton.net"
  ],
  "CertificateName": "bergen.famillegratton.net",
  "SerialNumber": 26,
  "Comments": [
    "pc principal"
  ]
}
```

Notes:

- `Duration` is expressed in **years**. Because of the Apple 825-day rule, CA certificates
  (`IsCA: true`) are capped at **2 years**.
- For the valid `KeyUsage` strings, see the reference table under
  [Key usage values](#key-usage-values).
- `SerialNumber` and `Comments` are optional; `Comments` are purely documentation and have
  no functional effect.

## PKI / environment directory structure

An environment is a sandbox; different environments are different PKIs. Given this
environment:

```bash
[jfgratton@bergen certificatemanager]$ cm -e test env info test
```

the PKI directory tree looks like this (root CA dir = `CA`, server dir = `srv`, config
dir = `cfg`):

```text
test
├── CA
│   ├── index.txt
│   ├── index.txt.attr
│   ├── newcerts
│   │   ├── 0002.pem
│   │   ├── 0003.pem
│   │   └── 0004.pem
│   ├── serial
│   ├── testCA.crt
│   └── testCA.key
├── cfg
│   ├── gitea.json
│   ├── haproxy.json
│   ├── nexus.json
│   └── testCA.json
└── srv
    ├── certs
    │   ├── gitea.crt
    │   ├── haproxy.crt
    │   └── nexus.crt
    ├── csr
    │   ├── gitea.csr
    │   ├── haproxy.csr
    │   └── nexus.csr
    ├── java
    └── private
        ├── gitea.key
        ├── haproxy.key
        └── nexus.key
```

- `CA/` — the root CA certificate and key live here directly
  - `index.txt` — the PKI database, listing every certificate (including the root CA)
  - `index.txt.attr` — the attribute file for the database
  - `serial` — the last serial number issued
  - `newcerts/` — a copy of each issued certificate, named after its serial in hexadecimal
- `cfg/` — the certificate **config** files
- `srv/certs/` — the issued certificates (`.crt`)
- `srv/csr/` — the certificate signing requests (kept in case you want to publish your PKI)
- `srv/private/` — the certificate private keys
- `srv/java/` — server certificates exported to PKCS#12 (`.p12`) and Java Keystore (`.jks`)

> **A note about `cm cert ls`:** this command lists certificate **config** files, not the
> certificates themselves. A config file may exist without a corresponding valid
> certificate. To confirm a certificate exists and is valid, use
> `cm cert verify PATH_TO_CERT`.

## Usage

The global flag `-e ENV` selects the environment file (default: `defaultEnv.json`).

### Environments — `cm env`

| Command | Aliases | Description |
|---------|---------|-------------|
| `cm env add [FILE]` | `create` | Create an environment file (prompts for directories). Defaults to `defaultEnv.json`. |
| `cm env list [DIR]` | `ls` | List the environment files. |
| `cm env info [FILE...]` | `explain` | Show the details of one or more environment files. |
| `cm env rm [FILE]` | `remove` | Delete an environment file. Defaults to `defaultEnv.json`. |

The `.json` extension is always optional.

### Certificates — `cm cert`

| Command | Aliases | Description |
|---------|---------|-------------|
| `cm cert create [CONFIG]` | | Create a CA or standard certificate. Prompts interactively if no config file is given. |
| `cm cert verify FILE` | | Verify a certificate. The file does not need to live inside the current PKI. |
| `cm cert list` | `ls` | List the certificate config files in the environment. |
| `cm cert revoke [CONFIG]` | `rm`, `remove` | Revoke a certificate by its config file name. |

Flags:

- `cm cert create -b, --keysize N` — private key size in bits (default `4096`)
- `cm cert create -j, --java` — also export a PKCS#12 (`.p12`) and Java Keystore (`.jks`)
- `cm cert verify -v, --verbose` — display the full output
- `cm cert verify -c, --comments` — display the config file comments (if any)
- `cm cert revoke -r, --remove` — also physically remove the certificate artefacts from the PKI

### First run

On its first run the tool creates its per-user config directory,
`$HOME/.config/JFG/certificatemanager/`. Before you can do anything useful you need an
environment file describing where your PKI will live.

The easiest way is:

```bash
cm env add            # creates the default environment (defaultEnv.json)
cm env add test       # creates a named environment (used later with: cm -e test ...)
```

### Create a CA certificate

The first step in building a PKI is a root CA.

If you already have a CA config file (e.g. `cfg/rootCA.json`):

```bash
cm cert create rootCA
```

Or create one interactively — answer **TRUE** when asked whether this is a CA certificate:

```bash
cm cert create
```

### Create a standard SSL/TLS certificate

Same process as above, but answer **FALSE** to the "is this a CA certificate?" prompt (or
use a config file whose `IsCA` is `false`):

```bash
cm cert create bergen.famillegratton.net
```

### Java certificates

Passing `-j` to `cm cert create` converts the freshly-created `.crt` into a PKCS#12 (`.p12`)
file and then into a Java Keystore (`.jks`). Both are written to `srv/java/`.

> **Caveat:** the `.p12` → `.jks` conversion relies on the external `keytool` binary,
> which ships with any Java SDK/JRE. This dependency is intentionally **not** declared in the
> binary packages, because Java package names differ wildly across distributions. Make sure
> `keytool` is on your `PATH` if you use `-j`.

### Revoke a certificate

```bash
cm cert revoke haproxy.json        # mark as revoked in the PKI database
cm cert revoke -r haproxy.json     # also delete the certificate's files
```

You name the certificate **config** file (as shown by `cm cert ls`).

### Key usage values

Valid `KeyUsage` strings (case-insensitive):

`digital signature`, `content commitment`, `key encipherment`, `data encipherment`,
`key agreement`, `cert sign`, `crl sign`, `encipher only`, `decipher only`.

Two convenience catch-alls are also accepted when creating a cert interactively:
`CADEFAULTS` (sensible defaults for a root CA) and `CERTDEFAULTS` (defaults for a standard
certificate).

## Building and installing

Both the source code and binary packages (Alpine `.apk`, Debian `.deb`, RedHat `.rpm`,
Arch Linux) are provided.

### Build from source

Requires Go (see `go.version` for the version this release is built with).

```bash
git clone https://github.com/jeanfrancoisgratton/certificatemanager.git
cd certificatemanager/src
./build.sh [OUTPUT_DIR]     # defaults to /opt/bin
```

The build script names the binary `cm` on the `main`/`develop` branches, and
`cm-<branch>` on any other branch.

### Install from a binary package

Grab the package for your distribution, then:

| Distro | Command |
|--------|---------|
| Alpine | `apk add [--allow-untrusted] /tmp/PACKAGE.apk` |
| Debian/Ubuntu | `apt install /tmp/PACKAGE.deb` |
| RedHat/Fedora | `dnf localinstall /tmp/PACKAGE.rpm` |
| Arch Linux | `pacman -U /tmp/PACKAGE.pkg.tar.zst` |

(`--allow-untrusted` is only needed if you do not have the signing keys for the Alpine
package repository.)

## Packaging directories

The `__alpine/`, `__archlinux/`, `__debian/` and `__redhat/` directories hold the packaging
metadata and build scripts used to produce the binary packages. They are geared toward the
maintainer's own container-based build pipeline and are not required to build or run the
tool from source.

## Documentation

- Changelog: [`docs/CHANGELOG.md`](docs/CHANGELOG.md)
- License: [`docs/LICENSE`](docs/LICENSE)
- Sample environment and certificate config files: [`docs/samples/`](docs/samples/)
