%define debug_package   %{nil}
%define _build_id_links none
%define _name certificateManager
%define _prefix /opt
%define _version 1.22.00
%define _rel 0
%define _arch x86_64
%define _binaryname cm

Name:       certificateManager
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
cd %{_sourcedir}/%{_name}-%{_version}/src
PATH=$PATH:/opt/go/bin go build -o %{_sourcedir}/%{_binaryname} .
strip %{_sourcedir}/%{_binaryname}

%clean
rm -rf $RPM_BUILD_ROOT

%pre
exit 0

%install
install -Dpm 0755 %{_sourcedir}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post

%preun

%postun

%files
%defattr(-,root,root,-)
%{_bindir}/%{_binaryname}


%changelog
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

