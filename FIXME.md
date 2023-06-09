- [x] customize packaging files
- [x] key usage shows as a numerical value instead of strings in `ca verify`
- [x] completely borked the x509.KeyUsage to []strings mapping
- [x] APK package might not work as it needs libresolv.so.2, which just does not exist on Alpine. Workaround needed (CGO_ENABLED=0, maybe ?)
- [ ] ca.prompt4values() is beyond ugly... Interface{}, much ? :-p
- [ ] not that big an issue, but exportable variables that do not need to be exportable should not be marked as exportable.
- [ ] `cert create` creates a certificate, but it is not signed

<br><br><br>