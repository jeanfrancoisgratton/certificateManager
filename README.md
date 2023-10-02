<H1>certificateManager</H1>

A GO tool to generate and sign all of your SSL certificates

<H2>Overview</H2>

This tool uses GO's x509 package to:<br>
- Generate<br>
- Sign<br>
- Verify all of your SSL certificates, your custom Root Certificate Authority (Root CA) certs, and even your Java certs in the .JKS (Java Keystore) certificates.

<H3>What that tool does do</H3>
- Create your own custom PKI (Private Key Infrastructure), and handle all certificates from under it
- Create, if needed, a root CA cert<br>
- Create and sign "normal" SSL server certs<br>
- Verify those certs<br>
- Revoke certs<br>

<H3>What that tool does not do</H3>
- Sign certificates against a remote CA:<br>
  - No CRL (Certificate Revokation List) is implemented<br>
  - No CDP (Certificate Distribution Point) is implemented<br>
- Any operation against a remote CA, actually.<br><br>

Bear in mind : this software is intended to run on an internal network.<br>
For instance, if you wish to publish your own rootCA certificate, it's yours do deploy it at VeriSign, etc....
<br><br>
<H2>Concepts</H2>

<H3>Environments</H3>

In most use cases, you would see a single rootCA or rootCA + intermediateCA present in a PKI.<br>

You might wish to have this tool in a docker container and be able to manage multiple PKIs, so to facilitate this, we introduce the idea of *environments*.

An environment is a set of rules that define the directory structure that will be used by a PKI.
This allows handling of multiple PKIs as sandboxed environments.

The contents of an environment file is such:<br>
```json
{
  "CertificateRootDir": "/Users/jfgratton/.config/certificatemanager/certificates",
  "RootCAdir": "rootCA",
  "ServerCertsDir": "servers",
  "CertificatesConfigDir": "conf",
  "RemoveDuplicates": true
}
```

The file is in JSON format; every key in the file (except the last one, `RemoveDuplicates`) are string values representing a path. The first path **must** be absolute, while the others are relative to it.<br>
(btw... `RemoveDuplicates` is meaningless for now, as that key is not treated -yet- anywhere in my code)

You switch between environments with the `-e` flag. Not using this flag will assume that you use the default environment file, `$HOME/.config/certficatemanager/default.Env` , assuming of course that the file is there.
<br><br>
<H3>Certificates, root certificates and certificate config files</H3>
A certificate (extension .crt) is the actual x509 SSL file that you might wish to deploy on a server, for example.
The root certificate (or rootCA) is the certificate used to sign (validate) all other certificates within the PKI.
A certificate config file is the JSON file that this tool uses to generate the certificate file.

A typical certificate config file looks like this:

```json
{
  "Country": "CA",
  "Province": "Quebec",
  "Locality": "Blainville",
  "Organization": "myorg.net",
  "OrganizationalUnit": "myorg",
  "CommonName": "myorg.net root CA",
  "IsCA": true,
  "EmailAddresses": [
    "certs@myorg.net",
    "certs@org,net"
  ],
  "Duration": 10,
  "KeyUsage": [
    "certs sign",
    "crl sign",
    "digital signature"
  ],
  "DNSNames": [
    "myorg.net",
    "myorg.com",
    "lan.myorg.net"
  ],
  "IPAddresses": [
    "10.0.0.1",
    "127.0.0.1"
  ],
  "CertificateName": "sampleCert",
  "SerialNumber": 1,
  "Comments": [
    "To see which values to put in the KeyUsage field, see https://pkg.go.dev/crypto/x509#KeyUsage",
    "Strip off 'KeyUsage' from the const name and there you go.",
    "",
    "Please note that this field offers no functionality and is strictly here for documentation purposes"
  ]
}
```
<br>
<H2>PKI / environment directory structure</H2>

As mentioned above, an *environment* is a sandbox. Different environments represent different PKIs.
We'll use the variables from `sampleEnv.json` here to describe the structure. <br>

Here, I have an environment called `test (test.json)` :
```
[16:56:29|jfgratton@bergen:certificatemanager]: cm env ls
Number of environment files: 1
┏━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━┓
┃ Environment file ┃ File size ┃ Modification time   ┃
┣━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━┫
┃ test.json        ┃ 145       ┃ 2023/10/02 16:56:29 ┃
┗━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━┛

[16:56:38|jfgratton@bergen:certificatemanager]: cm env explain test
┏━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃ Environment file ┃ Certificate root dir ┃ CA dir ┃ Server certificates dir ┃ Certificates config dir ┃
┣━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ test.json        ┃ /test                ┃ CA     ┃ srv                     ┃ cfg                     ┃
┗━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

So you see, my environment, my sandbox, my PKI, sits under `/test/`

Now, I've cheated a bit here, I've already created some certs, to show you the directory structure:<br>
```[17:16:19|jfgratton@bergen:/test]: cm -e test cert ls
Number of certificates: 4
┏━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━┓
┃ Cert name    ┃ Common Name       ┃ File size ┃ Modification time   ┃
┣━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━┫
┃ gitea.json   ┃ git.myorg.net     ┃ 488       ┃ 2023/10/02 17:16:11 ┃
┃ haproxy.json ┃ haproxy.myorg.net ┃ 515       ┃ 2023/10/02 17:16:17 ┃
┃ nexus.json   ┃ nexus.myorg.net   ┃ 518       ┃ 2023/10/02 17:16:19 ┃
┃ testCA.json  ┃ myorg.net root CA ┃ 726       ┃ 2023/10/02 17:15:28 ┃
┗━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━┛
```
<br>**A NOTE ABOUT `cm cert ls`:**<br>
This command lists certificate **config** files, not certificate **files**. This means a config file might be present, but no valid certificate being present.<br>
If you wish to see that a certificate exists (and is valid) : `cm cert verify $PATH_TO_CERTIFICATE_FILE`



<H2>How do we use the software</H2>
<H3>Create an environment file</H3>

As mentioned earlier, at the initial run of the software, it will create a few files in `$HOME/.config/certificatemanager`:<br>
- sampleCert.json
- sampleCert-README.txt
- sampleEnv.json
- sampleEnv-README.txt

Read both `.txt` files for further explanations. For now, the software is not yet usable; you need an environment file to run.<br>
By default, if you do not pass an `-e` argument to the app, it will assume `-e defaultEnv.json` (you do not need to provide the extension, btw)<br>

The easiest way, then, to create that default Environment file is: `cm env create`.<br>
You could add a filename such as `test`. If you provide another name, this means that any execution of the app will need the `-e ENVFILE` args (ENVFILE being the filename you've selected)<br>

By default, the software runs with the `-e defaultEnv.json` flag as a default environment file (which is why you need to adapt the above file with sane values). This will create the correct directory structure this software needs to operate

<H3>Create a CA cert</H3>
The very first step in building your own PKI is to have a root CA (root certificate authority)<br>
Either you already have your own json CA file (`cfg/rootCA.json`, in this example), or you will need to create your own:<br>
<H4>You have your own config</H4>
`cm cert create rootCA` (assuming that your CA config file is actually `rootCA.json`)
<H4>Create your own config</H4>
`cm cert create` <-- ensure that you select TRUE for a root CA when prompted<br>

<H3>Create "standard" SSL certs</H3>
The process is exactly as the one above, except that this time you specify that you are not creating a CA certificate<br>

This means that you follow the steps, above, and if you create a new file, you will need to answer FALSE to the prompt where it asks you if this is a CA cert,\


<H2>Building, installing CertificateManager</H2>
I provide both the source code and Alpine (APK), Debian-based (DEB) or RedHat-based (RPM) binary packages.

<H3>Install from source</H3>
- Clone the repo<br>
- from the `src/` directory, run: `./build [-o OUTPUTDIR]`, where OUTPUTDIR is where you want the final binary to be copied. By default it uses `/opt/bin/`

*NOTE: The script assumes that the building user has sudo rights to write in `/opt/bin` and strip the binary from debugging code*

<H3>Install from binary packages</H3>
Go in the Releases link from this site, and pick your package, once downloaded, in say, `/tmp/` :

<H4>Alpine (APK)</H4>
```bash
$ apk add /tmp/$PACKAGENAME
```

<H4>Debian-based (DEB)</H4>
```bash
$ apt install /tmp/$PACKAGENAME
```

<H4>RedHat-based (RPM)</H4>
```bash
$ dnf localinstall /tmp/$PACKAGENAME
```

<H2>A note about some extra directories and files</H2>
The following directories:<br>
- __alpine/<br>
- __debian/<br>
- ./certificateManager.spec<br>
- ./rpmbuild-deps.sh<br><br>

Those dirs and files are needed for my own home setup. That setup relies on my own custom Docker containers to build the binary packages.
Those containers will be made available once I manage to strip them from my personal information, but it is an involved process, so for now don't count on them.
(teaser: it's a pity, it works so well :D)
