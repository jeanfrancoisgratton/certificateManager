%define debug_package   %{nil}
%define _build_id_links none
%define _name   certificatemanager
%define _prefix /opt
%define _version 0.500
%define _rel 0
%define _arch x86_64
%define _binaryname cm

Name:       certificatemanager
Version:    %{_version}
Release:    %{_rel}
Summary:    certificatemanager

Group:      SSL
License:    GPL2.0
URL:        https://github.com/jeanfrancoisgratton/certificateManager

Source0:    %{name}-%{_version}.tar.gz
BuildArchitectures: x86_64
BuildRequires: gcc
#Requires: sudo
#Obsoletes: vmman1 > 1.140

%description
RootCA and server SSL certificate manager

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
* Sat Jun 03 2023 builder <builder@famillegratton.net> 0.500-0
- new package built with tito (reinitialized tito)

* Sat Apr 22 2023 builder <builder@famillegratton.net> 0.300-0
- Completion of branch 0.300 (jean-francois@famillegratton.net)
- Started ca.EditCACertificate() (jean-francois@famillegratton.net)
- Added the capability to see the comments in the config file (jean-
  francois@famillegratton.net)
- Groundwork for branch 0.300 (jean-francois@famillegratton.net)
- Updated ROADMAP, closing branch 0.200 (jean-francois@famillegratton.net)
- Updated doc (jean-francois@famillegratton.net)
- Updated DEB packaging doc (jean-francois@famillegratton.net)
- Updated DEB packaging doc (jean-francois@famillegratton.net)
- Updated repo info (jean-francois@famillegratton.net)
- Prettifying the output, release bump (jean-francois@famillegratton.net)
- Sample config file is more explicit (jean-francois@famillegratton.net)

* Thu Apr 20 2023 builder <builder@famillegratton.net> 0.200-0
- Updated doc (jean-francois@famillegratton.net)
- Fixed KeyUsage issue, partial doc update (jean-francois@famillegratton.net)
- Doc output reformatting (jean-francois@famillegratton.net)
- Moved file around (jean-francois@famillegratton.net)
- Interim commit with revamped ca verify (jean-francois@famillegratton.net)
- Removed old samples, revamped ROADMAP.md (jean-francois@famillegratton.net)
- Completed x509.KeyUsage handling (jean-francois@famillegratton.net)
- Tag 0.200 stub (jean-francois@famillegratton.net)
- Some code uncluttering (jean-francois@famillegratton.net)
- Another round of prettifying (jean-francois@famillegratton.net)
- Simplified source file (removed clutter) (jean-francois@famillegratton.net)
- Prettified 'ca verify' output (jean-francois@famillegratton.net)
- Completed the first iteration of 'ca verify' (jean-
  francois@famillegratton.net)

* Sun Apr 16 2023 builder <builder@famillegratton.net> 0.101-0
- Version bump and fixes (jean-francois@famillegratton.net)

* Sun Apr 16 2023 builder <builder@famillegratton.net> 0.100-0
- Initial dry-run on packaging

