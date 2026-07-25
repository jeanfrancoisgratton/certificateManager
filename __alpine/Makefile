# nxtools — Alpine (APK) packaging
#
# Recipes:
#   make build     -> abuild -r  (build the .apk and refresh the local repo index)
#   make release   -> upload the freshly built .apk to the Nexus apkLocal repo
#   make clean     -> run all abuild cleanup tasks
#   make info      -> print the values extracted from APKBUILD
#
# Package metadata (name / version / pkgrel / arch / binary) is read straight
# from APKBUILD so there is a single source of truth.

SHELL := /bin/sh

APKBUILD  := APKBUILD

NAME      := $(shell . ./$(APKBUILD); echo $$pkgname)
VERSION   := $(shell . ./$(APKBUILD); echo $$pkgver)
PKGREL    := $(shell . ./$(APKBUILD); echo $$pkgrel)
ARCH      := $(shell . ./$(APKBUILD); echo $$arch)
BINARY    := $(shell . ./$(APKBUILD); echo $$binaryname)

NEXUSREPO := apkLocal

# abuild -r drops the package at:
#   /data/packages/<name>/<arch>/<name>-<version>-r<pkgrel>.apk
APK       := /data/packages/$(NAME)/$(ARCH)/$(NAME)-$(VERSION)-r$(PKGREL).apk

.PHONY: build release clean info

# Alpine reindex is a no-op (Nexus rebuilds APKINDEX on every upload), so the
# upload deliberately omits the -r/--reindex flag.
build:
	abuild -r

release: build
	/opt/bin/nxtools upload $(NEXUSREPO) $(APK)

clean:
	abuild clean cleanpkg cleancache cleanoldpkg

info:
	@echo "NAME    = $(NAME)"
	@echo "VERSION = $(VERSION)"
	@echo "PKGREL  = $(PKGREL)"
	@echo "ARCH    = $(ARCH)"
	@echo "BINARY  = $(BINARY)"
	@echo "APK     = $(APK)"
