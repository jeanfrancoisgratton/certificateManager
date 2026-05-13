# Building an RPM
___

As of stubber v2.00.00, RPM building has been revamped in a major way

1. Reliance on `rpmbuild` instead of `tito`, which is being dropped because:
   - python takes quite a bit of room in a docker image (the `rpmbuilder` container)
   - reliance on python's setuptools prone to break at any time

2. Redhat-related build scripts are now under a common location, `$PROJECT_ROOTDIR/__redhat`
3. Build process is now `Makefile`-driven instead of `tito`-driven

## How to build the RPM

The steps are very simple; here is an example :

```bash
# cd __redhat/
# make {rpm | rpmcl }
# make upload
# make commitcl
# git push
```

- The `make rpmcl` step builds the RPM and update the specfile's %changelog section
- The `make rpm` step just builds the RPM, but does not update the Changelog
- The `make changelog` step updates the %changelog section of the specfile, but does not `git commit -am` it to the repo
- The `make upload` step uploads the RPM to a Nexus Repository Manager server, using `nxtools`, so this step can be ignored if those requirements (nxtools installed, and access to an NxRM server) are not met
- `make commitcl && git push` push the changelog updates to gi