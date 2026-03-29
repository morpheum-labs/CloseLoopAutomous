#!/usr/bin/env bash
# Build the arms Docker image (Go static binary; see Dockerfile).
# Optionally tag and push to Docker Hub or another registry.
#
# Usage:
#   bash docker-build.sh [tag]
#   bash docker-build.sh [--push] [tag]     # tag and push after build
#
# Default image tag (when [tag] is omitted): derived from the repo — nearest annotated/lightweight tag,
# commits since that tag, and abbreviated commit id (git describe --tags --long --always from repo root).
# Example: v0.3.3-11-ge66807f. Override with an explicit [tag] argument or ARMS_IMAGE_TAG.
#
# Environment — build:
#   ARMS_IMAGE_NAME     Image name (default: arms)
#   ARMS_IMAGE_TAG      Force image tag (skips git describe; use for CI without git metadata)
#   ARMS_DOCKERFILE     Path to Dockerfile (default: Dockerfile next to this script)
#
# Environment — push (--push):
#   DOCKER_SPACE_SORA   Registry username or full path (e.g. myuser or ghcr.io/myorg)
#   DOCKER_TOKEN_SORA   Registry password or API token (stdin to docker login)
#   SKIP_BUILD          If 1, skip build and only tag/push existing local image
#
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${ROOT}"

PUSH=0
TAG=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --push) PUSH=1; shift ;;
    *) TAG="${1}"; shift ;;
  esac
done

if [[ -z "${TAG}" ]]; then
  if [[ -n "${ARMS_IMAGE_TAG:-}" ]]; then
    TAG="${ARMS_IMAGE_TAG}"
  else
    GIT_TOPLEVEL="$(git -C "${ROOT}" rev-parse --show-toplevel 2>/dev/null || true)"
    if [[ -n "${GIT_TOPLEVEL}" ]]; then
      TAG="$(git -C "${GIT_TOPLEVEL}" describe --tags --long --always 2>/dev/null || echo latest)"
    else
      TAG="latest"
    fi
  fi
fi

IMAGE_NAME="${ARMS_IMAGE_NAME:-arms}"
DOCKERFILE="${ARMS_DOCKERFILE:-Dockerfile}"

DOCKER_REGISTRY="${DOCKER_SPACE_SORA:-}"
DOCKER_TOKEN="${DOCKER_TOKEN_SORA:-}"

if [[ ! -f "${DOCKERFILE}" ]]; then
  echo "Error: Dockerfile not found: ${DOCKERFILE}" >&2
  exit 1
fi

if [[ "${SKIP_BUILD:-0}" != "1" ]]; then
  echo "Building ${IMAGE_NAME}:${TAG}..."
  docker build \
    -f "${DOCKERFILE}" \
    -t "${IMAGE_NAME}:${TAG}" \
    "${ROOT}"
  echo "OK: ${IMAGE_NAME}:${TAG}"
fi

if [[ "${PUSH}" != "1" ]]; then
  exit 0
fi

if [[ -z "${DOCKER_REGISTRY}" || -z "${DOCKER_TOKEN}" ]]; then
  echo "Error: --push requires DOCKER_SPACE_SORA and DOCKER_TOKEN_SORA." >&2
  echo "  DOCKER_SPACE_SORA = Docker Hub user or ghcr.io/org namespace" >&2
  echo "  DOCKER_TOKEN_SORA = registry password or token" >&2
  exit 1
fi

REMOTE_IMAGE="${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}"
echo "Tagging and pushing ${REMOTE_IMAGE}..."

if [[ "${DOCKER_REGISTRY}" == *"/"* ]]; then
  REGISTRY_HOST="${DOCKER_REGISTRY%%/*}"
  echo "${DOCKER_TOKEN}" | docker login "${REGISTRY_HOST}" -u "${DOCKER_REGISTRY#*/}" --password-stdin
else
  echo "${DOCKER_TOKEN}" | docker login -u "${DOCKER_REGISTRY}" --password-stdin
fi

docker tag "${IMAGE_NAME}:${TAG}" "${REMOTE_IMAGE}"
docker push "${REMOTE_IMAGE}"

REMOTE_LATEST="${DOCKER_REGISTRY}/${IMAGE_NAME}:latest"
docker tag "${IMAGE_NAME}:${TAG}" "${REMOTE_LATEST}"
docker push "${REMOTE_LATEST}"

echo "Published ${REMOTE_IMAGE} and ${REMOTE_LATEST}"
