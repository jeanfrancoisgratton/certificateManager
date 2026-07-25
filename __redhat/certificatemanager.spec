%define debug_package   %{nil}
%define _build_id_links none
%define _name certificatemanager
%define _prefix /opt
%define _version 1.8.0
%define _rel 2
%define _arch x86_64
%define _bash_completionsdir /usr/share/bash-completion/completions
%define _zsh_completionsdir  /usr/share/zsh/site-functions
%define _binaryname cm

Name:       certificatemanager
Version:    %{_version}
Release:    %{_rel}
Summary:    Certificates and CA management tool

Group:      SSL tools
License:    GPL2.0
URL:        https://github.com/jeanfrancoisgratton/certificateManager

Source0:    %{name}-%{_version}.tar.gz
BuildArchitectures: x86_64
BuildRequires: gcc

%description
Certificates and CA management tool

%prep
%autosetup

%build
cd src
go mod download
PATH=$PATH:/opt/go/bin CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o %{_builddir}/%{name}-%{version}/%{_binaryname} .

%clean
rm -rf $RPM_BUILD_ROOT

%pre

%install
rm -rf %{buildroot}
install -Dpm 0755 %{_builddir}/%{name}-%{version}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post
# Bash completion — always install
%{_bindir}/%{_binaryname} completion bash > %{_bash_completionsdir}/{_binaryname}

# Zsh completion — only if zsh is present
if command -v zsh > /dev/null 2>&1; then
    mkdir -p %{_zsh_completionsdir}/zsh/site-functions
    %{_bindir}/%{_binaryname} completion zsh > %{_zsh_completionsdir}/_{_binaryname}
fi

%preun

%postun
if [ $1 -eq 0 ]; then
    # $1 == 0 means this is a full uninstall, not an upgrade
    rm -f %{_bash_completionsdir}/%{_binaryname}
    rm -f %{_zsh_completionsdir}/_%{_binaryname}
fi

%files
%defattr(-,root,root,-)
%{_bindir}/%{_binaryname}


%changelog
* Sat Jul 25 2026 Binary package builder <builder@famillegratton.net> 1.8.0-2
- RPMBUILDER: added missing variables in specfile
- APKBUILDER: post packaging scripts misnamed
- corrected wrong specfile name issue

* Sat Jul 25 2026 Binary package builder <builder@famillegratton.net> 1.8.0-1
- Refactored BUILDERs packaging, updated builddeps, upgraded helperFunctions and customError
- sanitized code, cleared go vet ./...
- Fixed typo
- added changelog to specfile
- added archlinux support to binary packaging, go version bump, builddeps updates, simplified error handling

* Wed May 13 2026 Binary package builder <builder@famillegratton.net> 1.70.00-0
- added archlinux support to binary packaging, go version bump, builddeps updates, simplified error handling



* Mon Aug 25 2025 Binary package builder <builder@famillegratton.net> 1.61.00-0
- GO version bump and misc (jean-francois@famillegratton.net)
- updated the outdated CHANGELOG.md, nothing more (jean-
  francois@famillegratton.net)



* Mon Jan 13 2025 APK Builder <builder@famillegratton.net> 1.59.00-0
- CN will now default to cert name if omitted (jean-
  francois@famillegratton.net)
- removed go install code (builder@famillegratton.net)

* Sat Nov 09 2024 APK Builder <builder@famillegratton.net> 1.58.00-0
- Fixed config path issue, removed sample certs, GO version bump (jean-
  francois@famillegratton.net)
- Fixed go dep (builder@famillegratton.net)

* Tue Oct 15 2024 RPM Builder <builder@famillegratton.net> 1.57.00-2
- re-initiialized tito

* Sat Jul 20 2024 DEB Builder <builder@famillegratton.net> 1.57.00-0
- new package built with tito

* Sun Jun 02 2024 RPM Builder <builder@famillegratton.net> 1.55.00-0
- Now preventing CAs to last longer than 800 days (jean-
  francois@famillegratton.net)
- Minor doc update (jean-francois@famillegratton.net)
- Updated builddeps (jean-francois@famillegratton.net)
- Forgot updating to doc files (jean-francois@famillegratton.net)
- Fixed rpmbuild script that used old specfile name
  (builder@famillegratton.net)

* Sun May 12 2024 RPM Builder <builder@famillegratton.net> 1.52.00-0
- Hopefully done with the error handling issues (jean-
  francois@famillegratton.net)
- Fixed most error handling issues in the cert subpackage; more to come (jean-
  francois@famillegratton.net)
- Software version bump (jean-francois@famillegratton.net)
- Firming up error handling (jean-francois@famillegratton.net)

* Sun May 12 2024 RPM Builder <builder@famillegratton.net> 1.51.00-0
- Fixed error handling (jean-francois@famillegratton.net)
- New release for DEB package (builder@famillegratton.net)
- New package release to lowercase its name everywhere
  (builder@famillegratton.net)

* Sun May 05 2024 RPM Builder <builder@famillegratton.net> 1.50.00-1
- 

* Sun May 05 2024 RPM Builder <builder@famillegratton.net>
- new package built with tito

* Thu Mar 28 2024 RPM Builder <builder@famillegratton.net>
- Minor output fixes, config move, go version bump (jean-
  francois@famillegratton.net)

* Fri Mar 01 2024 RPM Builder <builder@famillegratton.net>
- missing checksums file (jean-francois@famillegratton.net)
- Automatic commit of package [certificateManager] release [1.24.00-0].
  (builder@famillegratton.net)
- go mod tidy (jean-francois@famillegratton.net)
- GO version bump (jean-francois@famillegratton.net)
- updated gitignore (jean-francois@famillegratton.net)
- upgraded some modules in go.mod (jean-francois@famillegratton.net)
- Updated GO version, packages versions (jean-francois@famillegratton.net)
- More cleanup (jean-francois@famillegratton.net)
- Repo cleanup (jean-francois@famillegratton.net)

* Fri Mar 01 2024 RPM Builder <builder@famillegratton.net>
- go mod tidy (jean-francois@famillegratton.net)
- GO version bump (jean-francois@famillegratton.net)
- updated gitignore (jean-francois@famillegratton.net)
- upgraded some modules in go.mod (jean-francois@famillegratton.net)
- Updated GO version, packages versions (jean-francois@famillegratton.net)
- More cleanup (jean-francois@famillegratton.net)
- Repo cleanup (jean-francois@famillegratton.net)

* Thu Dec 07 2023 RPM Builder <builder@famillegratton.net> 1.23.00-0
- GO version bump (jean-francois@famillegratton.net)

* Thu Nov 02 2023 RPM Builder <builder@famillegratton.net> 1.22.00-0
- Fixed issue where already existing Java certs for a non-existant cert
  remained (jean-francois@famillegratton.net)
- sync zenika-> (jean-francois@famillegratton.net)
- Interim sync (jean-francois@famillegratton.net)
- Fixed duplicate path issue in cm cert rm (jean-francois@famillegratton.net)
- forgotten doc updates (builder@famillegratton.net)

* Thu Nov 02 2023 RPM Builder <builder@famillegratton.net> 1.21.00-0
- Version bump in cmd/root.go (jean-francois@famillegratton.net)
- env subcommand tweaks (jean-francois@famillegratton.net)

* Tue Oct 31 2023 RPM Builder <builder@famillegratton.net> 1.20.06-0
- more explicit error message in cert verify (jean-francois@famillegratton.net)
- GO pkgs upgrades (jean-francois@famillegratton.net)

* Mon Oct 16 2023 RPM Builder <builder@famillegratton.net> 1.20.05-0
- New version numbering scheme (jean-francois@famillegratton.net)
- Cosmetic change (jean-francois@famillegratton.net)

* Sun Oct 15 2023 RPM Builder <builder@famillegratton.net> 1.205-0
- Fixed wrong path for private keys (jean-francois@famillegratton.net)
- version bump for debian package (builder@famillegratton.net)

* Sun Oct 15 2023 RPM Builder <builder@famillegratton.net> 1.200-0
- Completed directory name simplification (jean-francois@famillegratton.net)
- Completed filepath.Join() revamp (jean-francois@famillegratton.net)
- Finished with the environment package (jean-francois@famillegratton.net)
- GO and software version bumps (jean-francois@famillegratton.net)

* Thu Oct 05 2023 RPM Builder <builder@famillegratton.net> 1.100-0
- Added Java cert in Doc, fixed duplicate certificate creation issue (jean-
  francois@famillegratton.net)
- Fixed cm cert verify (jean-francois@famillegratton.net)
- Fixed issue where comments where not displayed in cm cert verify (jean-
  francois@famillegratton.net)

* Tue Oct 03 2023 RPM Builder <builder@famillegratton.net> 1.010-1
- Fixed typo in directory name (jean-francois@famillegratton.net)

* Tue Oct 03 2023 RPM Builder <builder@famillegratton.net> 1.010-0
- Fixed issue where serial number was not incremented in the certificate (jean-
  francois@famillegratton.net)
- Build packages already take care of strip (jean-francois@famillegratton.net)
- Fixed issue with the 'strip' binary failing on arm64 arch (jean-
  francois@famillegratton.net)

* Tue Oct 03 2023 RPM Builder <builder@famillegratton.net> 1.001-0
- Completed README.md (jean-francois@famillegratton.net)
- Fixed bug where config dir was being recursively scanned (jean-
  francois@famillegratton.net)
- Version bump : minor fixes (jean-francois@famillegratton.net)
- Package name change (jean-francois@famillegratton.net)
- Sync zenika-> (jean-francois@famillegratton.net)
- Doc update (...part3) (jean-francois@famillegratton.net)
- Doc update (...final?) (jean-francois@famillegratton.net)
- Doc update part 1 (jean-francois@famillegratton.net)
- Fixed apk packaging (jean-francois@famillegratton.net)
- Fixed version in deb packaging (jean-francois@famillegratton.net)

* Sun Oct 01 2023 RPM Builder <builder@famillegratton.net> 1.000-0
- First prod-ready version

