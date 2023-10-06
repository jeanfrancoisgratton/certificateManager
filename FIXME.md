- [x] chgdir to a relative path (as defined in the environment json file) will end up with NOENT as we never "chgdired" to `CertificateRootDir` in the first place.
- [x] when loading or saving a cert file, ensure that we're either in `conf/` before loading/saving file
- [x] defaultEnv is being overwritten at every use of the software; we need to change that file to `sampleEnv.json`, and ignore NOENT files
- [x] serial# in cert is not incremented when creating a certificate from scratch
- [x] CRITICAL : we are not checking for duplicate certificates
- [x] Comments verbosity flag is not acted up in `cm cert verify`

<br><br><br>