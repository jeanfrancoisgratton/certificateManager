SPEC      := certificatemanager.spec
NAME      := $(shell rpmspec -q --qf '%{name}\n'    $(SPEC) | head -1)
VERSION   := $(shell rpmspec -q --qf '%{version}\n' $(SPEC) | head -1)
RELEASE   := $(shell rpmspec -q --qf '%{release}\n' $(SPEC) | head -1)
SOURCES   := $(HOME)/rpmbuild/SOURCES
RPMS      := $(HOME)/rpmbuild/RPMS
DEST      := /data
REPOROOT  := $(shell git rev-parse --show-toplevel)
NEXUSREPO := dnfLocal
TAG       := v$(VERSION)

.PHONY: tag release tarball check-tags changelog rpm rpmcl upload commitcl clean

# Intended workflow: make release

tarball:
	git -C $(REPOROOT) archive --prefix=$(NAME)-$(VERSION)/ HEAD \
	  | gzip > $(SOURCES)/$(NAME)-$(VERSION).tar.gz

check-tags:
	@# 1. The condition that silently breaks changelog generation:
	@if git describe --exact-match --tags HEAD >/dev/null 2>&1; then \
	  echo "WARNING: HEAD is already tagged $$(git describe --exact-match --tags HEAD)."; \
	  echo "         updateChangelog.sh will fall back to the previous tag for the range."; \
	fi
	@# 2. Local/remote tag agreement:
	@local=$$(git rev-parse -q --verify "refs/tags/$(TAG)" 2>/dev/null); \
	 remote=$$(git ls-remote --tags origin "refs/tags/$(TAG)" 2>/dev/null | cut -f1); \
	 if   [ -n "$$local" ] && [ -z "$$remote" ]; then echo "NOTE: $(TAG) exists locally but is not pushed."; \
	 elif [ -z "$$local" ] && [ -n "$$remote" ]; then echo "NOTE: $(TAG) is on origin but not local (run: git fetch --tags)."; \
	 elif [ -n "$$local" ] && [ "$$local" != "$$remote" ]; then \
	   echo "ERROR: $(TAG) points to different commits locally ($$local) and on origin ($$remote)."; \
	   exit 1; \
	 fi

changelog: check-tags tarball
	$(REPOROOT)/__redhat/updateChangelog.sh

rpm: tarball
	rpmbuild -bb $(SPEC)
	find $(RPMS) -name '*.rpm' -exec cp {} $(DEST)/ \;

rpmcl: tarball changelog
	rpmbuild -bb $(SPEC)
	find $(RPMS) -name '*.rpm' -exec cp {} $(DEST)/ \;

upload:
	/opt/bin/nxtools upload -r $(NEXUSREPO) $(DEST)/$(NAME)-$(VERSION)-$(RELEASE).x86_64.rpm

commitcl:
	@git add ${SPEC}
	@if git diff --cached --quiet -- ${SPEC}; then \
	  echo "NOTE: no staged changes to $(SPEC), skipping changelog commit."; \
	else \
	  git commit -m "chore: update changelog for $(VERSION)-$(RELEASE)"; \
	fi

clean:
	rm -f  $(SOURCES)/$(NAME)-$(VERSION).tar.gz
	rm -rf $(HOME)/rpmbuild/BUILD/$(NAME)-$(VERSION)

tag:
	@if git rev-parse -q --verify "refs/tags/$(TAG)" >/dev/null 2>&1; then \
	  echo "NOTE: $(TAG) already exists locally, skipping tag creation (release-only bump, version unchanged)."; \
	else \
	  git tag -a $(TAG) -m '$(TAG)'; \
	fi
	@if git ls-remote --tags origin "refs/tags/$(TAG)" 2>/dev/null | grep -q .; then \
	  echo "NOTE: $(TAG) already on origin, skipping push."; \
	else \
	  git push origin $(TAG); \
	fi

release:
	${MAKE} rpmcl
	${MAKE} commitcl
	git push
	${MAKE} upload
	${MAKE} tag
