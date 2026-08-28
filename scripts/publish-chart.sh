#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
CHART_DIR="${CHART_DIR:-${REPO_ROOT}/deploy/charts/alertmanager-webhook-adapter}"
PACKAGE_DIR="${PACKAGE_DIR:-/tmp/helm-packages}"
OCI_REGISTRY="${OCI_REGISTRY:-registry-1.docker.io/bougoucharts}"
CHARTS_REPO="${CHARTS_REPO:-https://github.com/bougou/charts.git}"
CHARTS_REPO_BRANCH="${CHARTS_REPO_BRANCH:-main}"

log() {
  printf '[publish-chart] %s\n' "$*"
}

die() {
  printf '[publish-chart] ERROR: %s\n' "$*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"
}

chart_metadata() {
  require_cmd helm
  CHART_NAME="$(helm show chart "${CHART_DIR}" | awk '/^name:/ {print $2}')"
  CHART_VERSION="$(helm show chart "${CHART_DIR}" | awk '/^version:/ {print $2}')"
  CHART_PACKAGE="${CHART_NAME}-${CHART_VERSION}.tgz"
}

lint_chart() {
  require_cmd helm
  log "lint chart at ${CHART_DIR}"
  helm lint "${CHART_DIR}"
  log "render chart template"
  helm template test "${CHART_DIR}" --namespace infra >/dev/null
}

package_chart() {
  chart_metadata
  mkdir -p "${PACKAGE_DIR}"
  log "package ${CHART_NAME} ${CHART_VERSION}"
  rm -f "${PACKAGE_DIR}/${CHART_PACKAGE}"
  helm package "${CHART_DIR}" -d "${PACKAGE_DIR}"
  CHART_PACKAGE_PATH="${PACKAGE_DIR}/${CHART_PACKAGE}"
  [[ -f "${CHART_PACKAGE_PATH}" ]] || die "package not found: ${CHART_PACKAGE_PATH}"
  log "created ${CHART_PACKAGE_PATH}"
}

push_oci() {
  chart_metadata
  [[ -n "${CHART_PACKAGE_PATH:-}" ]] || CHART_PACKAGE_PATH="${PACKAGE_DIR}/${CHART_PACKAGE}"
  [[ -f "${CHART_PACKAGE_PATH}" ]] || die "chart package not found, run package first"

  : "${DOCKERHUB_CHARTS_USER:?DOCKERHUB_CHARTS_USER is required (Docker Hub bougoucharts account)}"
  : "${DOCKERHUB_CHARTS_PASS:?DOCKERHUB_CHARTS_PASS is required (Docker Hub bougoucharts account PAT)}"

  log "login to registry-1.docker.io as ${DOCKERHUB_CHARTS_USER} (bougoucharts)"
  helm registry login registry-1.docker.io \
    --username "${DOCKERHUB_CHARTS_USER}" \
    --password "${DOCKERHUB_CHARTS_PASS}"

  log "push ${CHART_PACKAGE_PATH} to oci://${OCI_REGISTRY}"
  helm push "${CHART_PACKAGE_PATH}" "oci://${OCI_REGISTRY}"
}

update_charts_index() {
  chart_metadata
  [[ -n "${CHART_PACKAGE_PATH:-}" ]] || CHART_PACKAGE_PATH="${PACKAGE_DIR}/${CHART_PACKAGE}"
  [[ -f "${CHART_PACKAGE_PATH}" ]] || die "chart package not found, run package first"
  : "${CHARTS_REPO_TOKEN:?CHARTS_REPO_TOKEN is required}"

  require_cmd git
  require_cmd yq

  local workdir
  workdir="$(mktemp -d)"
  trap 'rm -rf "${workdir}"; trap - RETURN' RETURN

  local clone_url="${CHARTS_REPO}"
  if [[ "${clone_url}" == https://github.com/* ]]; then
    clone_url="https://x-access-token:${CHARTS_REPO_TOKEN}@${clone_url#https://}"
  fi

  log "clone charts repo into ${workdir}/charts-repo"
  git clone --depth 1 --branch "${CHARTS_REPO_BRANCH}" "${clone_url}" "${workdir}/charts-repo"

  log "copy ${CHART_PACKAGE} to charts repo"
  cp "${CHART_PACKAGE_PATH}" "${workdir}/charts-repo/charts/"

  log "regenerate index.yaml"
  (
    cd "${workdir}/charts-repo"
    ./scripts/make-index.sh
    ./scripts/generate-charts-doc.sh
  )

  log "commit and push charts repo"
  (
    cd "${workdir}/charts-repo"
    git config user.name "${GIT_AUTHOR_NAME:-github-actions[bot]}"
    git config user.email "${GIT_AUTHOR_EMAIL:-github-actions[bot]@users.noreply.github.com}"
    git add index.yaml README.md "charts/${CHART_PACKAGE}"
    if git diff --cached --quiet; then
      log "charts repo already up to date"
    else
      git commit -m "alertmanager-webhook-adapter ${CHART_VERSION} chart"
      git push origin "${CHARTS_REPO_BRANCH}"
    fi
  )

  rm -rf "${workdir}"
  trap - RETURN
}

publish_all() {
  lint_chart
  package_chart
  push_oci
  if [[ "${SKIP_INDEX_UPDATE:-}" != "true" ]]; then
    update_charts_index
  else
    log "skip charts repo index update (SKIP_INDEX_UPDATE=true)"
  fi
  log "done: ${CHART_NAME} ${CHART_VERSION}"
}

usage() {
  cat <<'EOF'
Usage: publish-chart.sh <command>

Commands:
  lint            Run helm lint and helm template
  package         Create a .tgz package
  push-oci        Push packaged chart to OCI registry
  update-index    Update bougou/charts index.yaml and README
  publish         lint, package, push OCI, and update charts repo index

Environment variables:
  CHART_DIR           Path to chart source directory
  PACKAGE_DIR         Directory for generated .tgz files
  OCI_REGISTRY        OCI registry reference (default: registry-1.docker.io/bougoucharts)
  DOCKERHUB_CHARTS_USER Docker Hub username for the bougoucharts account
  DOCKERHUB_CHARTS_PASS Docker Hub PAT for the bougoucharts account
  CHARTS_REPO_TOKEN   GitHub PAT with write access to bougou/charts
  CHARTS_REPO         Charts repo URL (default: https://github.com/bougou/charts.git)
  CHARTS_REPO_BRANCH  Charts repo branch (default: main)
  SKIP_INDEX_UPDATE   Set to "true" to skip bougou/charts index update
EOF
}

main() {
  local command="${1:-}"
  case "${command}" in
    lint)
      lint_chart
      ;;
    package)
      package_chart
      ;;
    push-oci)
      package_chart
      push_oci
      ;;
    update-index)
      package_chart
      update_charts_index
      ;;
    publish)
      publish_all
      ;;
    -h|--help|help|"")
      usage
      [[ -n "${command}" ]] || exit 0
      ;;
    *)
      die "unknown command: ${command}"
      ;;
  esac
}

main "$@"
