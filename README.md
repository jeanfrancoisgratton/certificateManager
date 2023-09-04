<H1>certificateManager</H1>

A tool to generate and sign all of your SSL certificates

<H2>Overview</H2>

This tool uses GO's x509 package to generate, sign and verify all of your SSL certificates, your custom Root Certificate Authority (Root CA) certs, and even your Java certs in the .JKS (Java Keystore) certificates.

<H3>What that tool does do</H3>
- Create a directory structure to hold all of the infrastructure needed to maintain a root CA and its dependent certs<br>
- Create, if needed, a root CA cert<br>
- Create and sign "normal" SSL server certs<br>
- Verify those certs<br>
- Revoke (in a future version) certs, with CRL<br>

<H3>What that tool does not do</H3>
- Sign certificates against a remote CA<br>
- Any operation against a remote CA, actually<br>

The work involved in dealing with all of the use-cases with remote CAs is not worth the effort, especially considering that I started that project for a specific home use-case

<H2>How does it work</H2>
At initial run, a sample file will be created in $HOME/.config/certificateManager/.<br>
This file should be copied as `defaultEnv.json` ASAP, with sane values.<br>
Read `environmentSample-README.txt` in the same directory for examples.

By default, the software runs with the `-e defaultEnv.json` flag as a default environment file (which is why you need to adapt the above file with sane values). This will create the correct directory structure this software needs to operate

You first need to either adapt `sampleEnv.json` as `defaultEnv.json`, or use `cm env create -e $ENVIRONMENTNAME` to create a new environment file

Then, you need to create your certificates, with `cm cert create $CERTIFICATENAME`. You will be prompted for the values needed to create the given certificate.
All relevant files will be stored in the directories as specified in the environment file (specified with the `-e` flag) you've chosen.

You can optionally create a Java (.jks) cert with the `-j` flag. This flag is ignored when the created certificate is a root CA.

All certificates, even CA certs, can be verified with `cm cert verify $CERTNAME`.

Not a single certificate or environment operation needs to specify an absolute path. All paths are computed from the environment file in `$HOME/.config/certificateManager/`

<H2>Building, installing CertificateManager</H2>
I provide both the source code and Alpine (APK), Debian-based (DEB) or RedHat-based (RPM) binary packages.

<H3>Install from source</H3>
- Clone the repo<br>
- from the `src/` directory, run: `./build [-o OUTPUTDIR]`, where OUTPUTDIR is where you want the final binary to be copied. By default it uses `/opt/bin/`

*NOTE: The script assumes that the building user has sudo rights to write in `/opt/bin` and strip the binary from debugging code*

<H3>Install from binary packages</H3>
Go in the Releases link from this site, and pick your package, once downloaded, in say, `/tmp/` :

<H4>Alpine (APK)</H4>
```
$ apk add /tmp/$PACKAGENAME
```

<H4>Debian-based (DEB)</H4>
```
$ apt install /tmp/$PACKAGENAME
```

<H4>RedHat-based (RPM)</H4>
```
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
