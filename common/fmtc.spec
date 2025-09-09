################################################################################

%define debug_package  %{nil}

################################################################################

Summary:        Simple utility for rendering fmtc formatted data
Name:           fmtc
Version:        1.1.0
Release:        0%{?dist}
Group:          Applications/System
License:        Apache License, Version 2.0
URL:            https://kaos.sh/fmtc

Source0:        https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  golang >= 1.24

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
Simple utility for rendering fmtc formatted data.

################################################################################

%prep
%setup -q

%build
if [[ ! -d "%{name}/vendor" ]] ; then
  echo -e "----\nThis package requires vendored dependencies\n----"
  exit 1
elif [[ -f "%{name}/%{name}" ]] ; then
  echo -e "----\nSources must not contain precompiled binaries\n----"
  exit 1
fi


pushd %{name}
  %{__make} %{?_smp_mflags} all
  cp LICENSE ..
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -dm 755 %{buildroot}%{_mandir}/man1

install -pm 755 %{name}/%{name} %{buildroot}%{_bindir}/

./%{name}/%{name} --generate-man > %{buildroot}%{_mandir}/man1/%{name}.1

%clean
rm -rf %{buildroot}

%post
if [[ -d %{_sysconfdir}/bash_completion.d ]] ; then
  %{name} --completion=bash 1> %{_sysconfdir}/bash_completion.d/%{name} 2>/dev/null
fi

if [[ -d %{_datarootdir}/fish/vendor_completions.d ]] ; then
  %{name} --completion=fish 1> %{_datarootdir}/fish/vendor_completions.d/%{name}.fish 2>/dev/null
fi

if [[ -d %{_datadir}/zsh/site-functions ]] ; then
  %{name} --completion=zsh 1> %{_datadir}/zsh/site-functions/_%{name} 2>/dev/null
fi

%postun
if [[ $1 == 0 ]] ; then
  if [[ -f %{_sysconfdir}/bash_completion.d/%{name} ]] ; then
    rm -f %{_sysconfdir}/bash_completion.d/%{name} &>/dev/null || :
  fi

  if [[ -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish ]] ; then
    rm -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish &>/dev/null || :
  fi

  if [[ -f %{_datadir}/zsh/site-functions/_%{name} ]] ; then
    rm -f %{_datadir}/zsh/site-functions/_%{name} &>/dev/null || :
  fi
fi

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_mandir}/man1/%{name}.1.*
%{_bindir}/%{name}

################################################################################

%changelog
* Tue Sep 09 2025 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- Added '-e/--eval' option to eval escape sequences
- Dependencies update

* Tue May 06 2025 Anton Novojilov <andy@essentialkaos.com> - 1.0.2-0
- Dependencies update

* Thu Feb 06 2025 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Code refactoring
- Dependencies update

* Tue Sep 24 2024 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Migrated to v13 of ek package
- Dependencies update

* Mon Jun 17 2024 Anton Novojilov <andy@essentialkaos.com> - 0.1.2-0
- Dependencies update

* Fri Mar 22 2024 Anton Novojilov <andy@essentialkaos.com> - 0.1.1-0
- Improved support information gathering
- Dependencies update

* Wed Nov 29 2023 Anton Novojilov <andy@essentialkaos.com> - 0.1.0-0
- Initial build
