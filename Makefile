# must be valid semver 2.0 version
ifndef $(VERSION)
LATEST_TAG=$(shell git describe --tags $(git rev-list --tags --max-count=1))
VERSION=$(LATEST_TAG:v%=%)
endif
EXECUTABLE=wurl

GO?=vgo
LDFLAGS+=-X main.Version=$(VERSION)

.PHONY: build install uninstall github-release build
.DEFAULT_GOAL := build

build:
	@echo "Building for current OS/architecture"
	$(GO) build -o ./$(EXECUTABLE) -v -ldflags="$(LDFLAGS)"

PREFIX ?= /usr/local
DESTDIR=
BIN=$(DESTDIR)$(PREFIX)/bin/

install:
	install -D ./$(EXECUTABLE) $(BIN)

uninstall:
	rm $(install_dir)/$(EXECUTABLE)

BUILD_DIR=build

temp_dir_name=$(EXECUTABLE)_$(1)_$(2)_v$(3)
build_os_arch=GOOS=$(1) GOARCH=$(2) $(GO) build -o $(3) -v -ldflags="$(LDFLAGS)"
pack_zip=zip -r -j $(1).zip $(1) && rm -rf $(1)
pack_tgz=tar -C $(1) -cpzf $(1).tar.gz ./ && rm -rf $(1)

define build_github_release
@echo "Building release package for OS $(1), arch $(2)"
$(eval temp_build_dir=$(BUILD_DIR)/$(call temp_dir_name,$(1),$(2),$(VERSION)))
@mkdir -p $(temp_build_dir)
$(eval ifeq ($(1),windows)
	temp_executable=$(temp_build_dir)/$(EXECUTABLE).exe
else
	temp_executable=$(temp_build_dir)/$(EXECUTABLE)
endif)
$(call build_os_arch,$(1),$(2),$(temp_executable))
$(eval ifeq ($(1),windows)
	pack_cmd = $(call pack_zip,$(temp_build_dir))
else
	pack_cmd = $(call pack_tgz,$(temp_build_dir))
endif)
@$(pack_cmd)
endef

github-release:
	$(call build_github_release,linux,amd64)
	$(call build_github_release,linux,386)
	$(call build_github_release,linux,arm)
	$(call build_github_release,darwin,amd64)
	$(call build_github_release,windows,amd64)
	$(call build_github_release,windows,386)

clean:
	rm -rf $(BUILD_DIR)
	$(GO) clean
