all: 
	cd ../src
	CGO_ENABLED=0 go build . -o ./cm.exe

prepare:
	sudo /opt/bin/install_golang.sh 1.20.5 amd64
	sudo apk update

install:
	test -d $(DESTDIR)$(PREFIX)/bin || mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -pm 755 $(BUILDDIR)/cm.exe $(DESTDIR)$(PREFIX)/bin/cm
	strip $(DESTDIR)$(PREFIX)/bin/cm
