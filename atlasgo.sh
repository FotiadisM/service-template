#!/bin/sh

# This is a little script that can be downloaded from the internet to install the Atlas CLI.
# It attempts to determine architecture and distribution to download and install the correct binary.
# It runs on most Unix shells like {a,ba,da,k,z}sh.

set -u

VERSION="${VERSION:-0.1}"
ATLAS_UPDATE_SRV="${ATLAS_UPDATE_SRV:-https://release.ariga.io/atlas}"
ATLAS_VERSION="${ATLAS_VERSION:-latest}"
ATLAS_DEBUG="${ATLAS_DEBUG:-false}"
CI="${CI:-false}"

main() {
    # Prerequisites
    need_cmd uname
    need_cmd mktemp
    need_cmd chmod
    need_cmd chown
    need_cmd mkdir

    local _skip_prompt=false
    # Consume flags
    while [ $# -gt 0 ]; do
        case "$1" in
            "-y")
                _skip_prompt=true
        esac
        shift
    done

    # Under some circumstances we skip prompting,
    # e.g. in CI systems or if there is no tty to read input from.
    if [ "$CI" = "true" ] || ! sh -c ': >/dev/tty' > /dev/null 2>&1; then
        _skip_prompt=true
    fi

    # Determine architecture and cpu type.
    get_architecture || return 1
    local _arch="$ARCH"
    local _platform="$PLATFORM"
    local _chipset="$CHIPSET"

    # Create a temporary download folder.
    local _dir
    _dir="$(ensure mktemp -d)"
    ensure mkdir -p "$_dir"

    # Build filename and download path.
    local _file="atlas-$_arch-$ATLAS_VERSION"
    local _url="$ATLAS_UPDATE_SRV/$_file"
    if [ "$ATLAS_DEBUG" = "true" ]; then
      _url="$_url?test=1"
    fi

    local _path="/usr/local/bin/atlas"
    if check_cmd "atlas"; then
        _path=$(which atlas)
    fi

    # Since this script is most likely piped into sh, there is no stdin to prompt the user on.
    # In that case explicitly connect to /dev/tty to read user input.
    if [ "$_skip_prompt" = "false" ]; then
        # We want to prompt the user. If there is a stdin, use it. If not, use /dev/tty.
        local _yn
        local _prompt="Install '$_file' to '$_path'? [y/N] "
        if [ -t 0 ]; then
            read -p "$_prompt" _yn
        else
            read -p "$_prompt" _yn < /dev/tty
        fi

        case "$_yn" in
            "y" | "yes" | "Y") ;;
            *) exit ;;
        esac
    fi

    echo "Downloading $_url"
    local _curlVersion=$(curl --version | head -n 1 | awk '{ print $2 }')
    local _ua="Atlas Installer/$VERSION ($_platform; $_chipset) cURL $_curlVersion"
    (cd "$_dir" && ensure curl -L -o "$_file" -A "$_ua" "$_url")

    case "$_arch" in
        *linux*)
            # Install the binary in path
            local _install="install -o root -g root -m 0755 $_dir/$_file $_path"
            if ! [ "$(id -u)" = 0 ]; then
              _install="sudo $_install"
            fi
            ensure_silent eval "$_install"
            ;;

        *darwin*)
            ensure chmod +x "$_dir/$_file"
            # On Mac sometimes the default path does not exist.
            local _mkdir="mkdir -p ${_path%/*}"
            local _mv="mv $_dir/$_file $_path"
            if ! [ "$(id -u)" = 0 ]; then
                _mkdir="sudo $_mkdir"
                _mv="sudo $_mv"
            fi
            ensure_silent eval "$_mkdir"
            ensure_silent eval "$_mv"
            if ! [ "$(id -u)" = 0 ]; then
                ensure_silent sudo chown root: "$_path"
            fi
            ;;
    esac

    # Run once to ensure atlas is installed correctly.
    ensure atlas version

    echo "Installation successful!"
}

get_architecture() {
    local _ostype _cputype _platform
    _ostype="$(uname -s)"
    _cputype="$(uname -m)"

    if [ "$_ostype" = Darwin ] && [ "$_cputype" = i386 ]; then
        # Darwin `uname -m` lies
        if sysctl hw.optional.x86_64 | grep -q ': 1'; then
            _cputype=x86_64
        fi
    fi

    CHIPSET="$_cputype"

    case "$_cputype" in

        xscale | arm | armv6l | armv7l | armv8l | aarch64 | arm64)
            _cputype=arm64
            ;;

        x86_64 | x86-64 | x64 | amd64)
            _cputype=amd64
            ;;

        *)
            err "unknown CPU type: $_cputype"

    esac

    case "$_ostype" in

        Linux | FreeBSD | NetBSD | DragonFly)
            _ostype=linux
            _platform=Linux
            # If the libc implementation is musl, or the glibc version is <2.31, use the musl build.
            if ldd --version 2>&1 | grep -q 'musl' || \
               [ $(version "$(ldd --version | awk '/ldd/{print $NF}')") -lt $(version "2.31") ]; then
              _cputype="$_cputype-musl"
            fi
            ;;

        Darwin)
            _ostype=darwin
            _platform=MacOS
            # We currently don't have an arm build for M1.
            # However, M1 is capable of running amd64 binaries.
            _cputype=amd64
            ;;

        *)
            err "unrecognized OS type: $_ostype"
            ;;

    esac

    ARCH="$_ostype-$_cputype"
    PLATFORM="$_platform"
}

need_cmd() {
    if ! check_cmd "$1"; then
        err "need '$1' (command not found)"
    fi
}

check_cmd() {
    command -v "$1" > /dev/null 2>&1
}

ensure() {
    if ! "$@"; then err "command failed: $*"; fi
}

ensure_silent() {
  if ! "$@"; then exit 1; fi
}

err() {
    echo "$1" >&2
    exit 1
}

version() {
  echo "$@" | awk -F. '{ printf("%d%03d%03d%03d\n", $1,$2,$3,$4); }';
}

main "$@" || exit 1
