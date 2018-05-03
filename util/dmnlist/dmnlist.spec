# =============================================================================
%define		name	dmnlist
%define		version	1.0.0
%define		release	1
%define		summary	Retrieve list of Decision Model and Notation objects.
%define		author	John Scherff <jscherff@24hourfit.com>
%define		gopath	%{_builddir}/go
%define		package github.com/jscherff/dmnsdk
%define		utildir util
# =============================================================================

Name:		%{name}
Version:	%{version}
Release:	%{release}%{?dist}
Summary:	%{summary}

License:	ASL 2.0
URL:		https://www.24hourfitness.com
Vendor:		24 Hour Fitness, Inc.
Prefix:		%{_bindir}
Packager: 	%{packager}
BuildRoot:	%{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)
Distribution:	el

BuildRequires:    golang >= 1.10.0, git >= 1.8.0

%description
The %{name} utility Retrieve list of Decision Model and Notation objects
and saves or displays the result in CSV format.

%prep

%build

  export GOPATH=%{gopath}
  export GIT_DIR=%{gopath}/src/%{package}/.git

  go get %{package}
  go build -ldflags='-X main.version=%{version}-%{release}' %{package}/%{utildir}/%{name}

%install

  test %{buildroot} != / && rm -rf %{buildroot}/*

  mkdir -p %{buildroot}%{_bindir}
  install -s -m 755 %{_builddir}/%{name} %{buildroot}%{_bindir}/

%clean

  test %{buildroot} != / && rm -rf %{buildroot}/*
  test %{_builddir} != / && rm -rf %{_builddir}/*

%files

  %defattr(-,root,root)
  %{_bindir}/*

%changelog
* Thu May 3 2018 - jscherff@24hourfit.com
- Initial build
