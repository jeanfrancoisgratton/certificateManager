# nxtools — Arch Linux (pacman) packaging
#
# Recipes:
#   make build     -> makepkg (build the package into PKGDEST)
#   make release   -> refresh the pacman DB (repo-add) then upload the package
#                     and database files to the Nexus archLocal repo (with reindex)
#   make clean     -> remove the makepkg working directories
#   make info      -> print the values extracted from PKGBUILD
#
# Package metadata (name / version / pkgrel / arch / binary) is read straight
# from PKGBUILD so there is a single source of truth.

SHELL := /bin/bash

PKGBUILD  := PKGBUILD

NAME      := $(shell . ./$(PKGBUILD); echo $$pkgname)
VERSION   := $(shell . ./$(PKGBUILD); echo $$pkgver)
PKGREL    := $(shell . ./$(PKGBUILD); echo $$pkgrel)
ARCH      := $(shell . ./$(PKGBUILD); echo $${arch[0]})
BINARY    := $(shell . ./$(PKGBUILD); echo $$_binaryname)

NEXUSREPO := archLocal
PKGDEST   := /data
PKGEXT    := pkg.tar.zst

# Arch database (repo-add) refresh: the pacman DB lives in the Nexus repo as
# <REPONAME>.db{,.tar.zst} / <REPONAME>.files{,.tar.zst}. Before uploading a new
# package we pull the current DB, add the package with repo-add, then re-upload.
NXTOOLS   := /opt/bin/nxtools
REPONAME  := famillegratton.net
NEXUSURL  := https://nexus.famillegratton.net:1808/archlinux
STAGING   := $(PKGDEST)/_staging

# makepkg drops the package at: $PKGDEST/<name>-<version>-<pkgrel>-<arch>.<PKGEXT>
PKG       := $(PKGDEST)/$(NAME)-$(VERSION)-$(PKGREL)-$(ARCH).$(PKGEXT)

.PHONY: build release clean info

build:
	rm -rf src pkg
	PKGDEST=$(PKGDEST) makepkg --force --syncdeps --noconfirm

release: build
	# --- refresh the Arch repo database with the freshly built package ---
	rm -rf $(STAGING); mkdir -p $(STAGING)
	cp $(PKG) $(STAGING)/
	cd $(STAGING) && \
	  $(NXTOOLS) download $(NEXUSURL)/$(REPONAME).db.tar.zst && \
	  $(NXTOOLS) download $(NEXUSURL)/$(REPONAME).files.tar.zst && \
	  repo-add $(REPONAME).db.tar.zst $(notdir $(PKG)) && \
	  rm -f $(REPONAME).db $(REPONAME).files && \
	  cp $(REPONAME).db.tar.zst $(REPONAME).db && \
	  cp $(REPONAME).files.tar.zst $(REPONAME).files
	# --- upload the package (with reindex), then the refreshed database files ---
	$(NXTOOLS) upload $(NEXUSREPO) $(PKG)
	cd $(STAGING) && for f in files files.tar.zst db db.tar.zst; do \
	  $(NXTOOLS) upload $(NEXUSREPO) $(REPONAME).$$f; \
	done
	rm -rf $(STAGING)
	${MAKE} clean

clean:
	rm -rf src pkg

info:
	@echo "NAME    = $(NAME)"
	@echo "VERSION = $(VERSION)"
	@echo "PKGREL  = $(PKGREL)"
	@echo "ARCH    = $(ARCH)"
	@echo "BINARY  = $(BINARY)"
	@echo "PKG     = $(PKG)"
