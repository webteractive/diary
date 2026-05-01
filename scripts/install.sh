#!/usr/bin/env sh
set -eu

REPO="${DIARY_REPO:-webteractive/diary}"
INSTALL_DIR="${DIARY_INSTALL_DIR:-/usr/local/bin}"
VERSION="${DIARY_VERSION:-latest}"

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "diary install: missing required command: $1" >&2
    exit 1
  fi
}

detect_os() {
  case "$(uname -s)" in
    Darwin) echo "darwin" ;;
    Linux) echo "linux" ;;
    MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
    *)
      echo "diary install: unsupported OS: $(uname -s)" >&2
      exit 1
      ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *)
      echo "diary install: unsupported architecture: $(uname -m)" >&2
      exit 1
      ;;
  esac
}

download() {
  url="$1"
  dest="$2"

  if command -v curl >/dev/null 2>&1; then
    if [ "$dest" = "-" ]; then
      curl -fsSL "$url" -w '%{url_effective}'
    else
      curl -fsSL "$url" -o "$dest"
    fi
    return
  fi

  if command -v wget >/dev/null 2>&1; then
    if [ "$dest" = "-" ]; then
      wget -q "$url" -O /dev/null --server-response 2>&1 | awk '/^  Location: / {print $2}' | tail -n 1
    else
      wget -q "$url" -O "$dest"
    fi
    return
  fi

  echo "diary install: missing required command: curl or wget" >&2
  exit 1
}

need tar
need mktemp

os="$(detect_os)"
arch="$(detect_arch)"

if [ "$VERSION" = "latest" ]; then
  need sed
  latest_url="$(download "https://github.com/${REPO}/releases/latest" "-")"
  version_label="$(printf '%s' "$latest_url" | sed 's#.*/tag/##')"
  if [ -z "$version_label" ] || [ "$version_label" = "$latest_url" ]; then
    echo "diary install: could not resolve latest release version" >&2
    exit 1
  fi
else
  version_label="$VERSION"
fi
base_url="https://github.com/${REPO}/releases/download/${version_label}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT INT TERM

if [ "$os" = "windows" ]; then
  need unzip
  archive="diary_${version_label}_${os}_${arch}.zip"
else
  archive="diary_${version_label}_${os}_${arch}.tar.gz"
fi

echo "diary install: downloading ${archive}"
download "${base_url}/${archive}" "${tmpdir}/${archive}"

if [ "$os" = "windows" ]; then
  unzip -q "${tmpdir}/${archive}" -d "$tmpdir"
  binary="${tmpdir}/diary.exe"
else
  tar -xzf "${tmpdir}/${archive}" -C "$tmpdir"
  binary="${tmpdir}/diary"
fi

if [ ! -f "$binary" ]; then
  echo "diary install: archive did not contain expected binary" >&2
  exit 1
fi

chmod +x "$binary"
mkdir -p "$INSTALL_DIR"

dest="${INSTALL_DIR}/diary"
if [ "$os" = "windows" ]; then
  dest="${INSTALL_DIR}/diary.exe"
fi

if [ -w "$INSTALL_DIR" ]; then
  mv "$binary" "$dest"
else
  echo "diary install: ${INSTALL_DIR} is not writable, trying sudo" >&2
  sudo mv "$binary" "$dest"
fi

echo "diary installed to ${dest}"
"$dest" --help >/dev/null
