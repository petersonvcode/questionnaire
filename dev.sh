#!/usr/bin/env bash
set -euo pipefail
set -m

cd "$(dirname "${BASH_SOURCE[0]}")"

back_pid=""
front_pid=""

cleanup() {
  trap - EXIT INT TERM
  [[ -n "$back_pid" ]] && kill -- -"$back_pid" 2>/dev/null || true
  [[ -n "$front_pid" ]] && kill -- -"$front_pid" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

need() {
  command -v "$1" >/dev/null || { echo "missing command: $1"; exit 1; }
}
need go
need make

export NVM_DIR="${NVM_DIR:-$HOME/.nvm}"
if [[ ! -s "$NVM_DIR/nvm.sh" ]]; then
  echo "nvm not found at $NVM_DIR"
  exit 1
fi
# shellcheck disable=SC1091
set +u
. "$NVM_DIR/nvm.sh"
set -u

echo "checking backend dependencies..."
(cd backend && go mod download)

echo "checking frontend dependencies..."
(
  set +u
  cd frontend
  nvm install
  nvm use
  npm install
)

mkdir -p logs
: > logs/backend.log
: > logs/frontend.log

echo "backend  http://localhost:9090  (logs/backend.log)"
echo "frontend http://localhost:3000  (logs/frontend.log)"
echo "follow logs with: tail -f logs/backend.log logs/frontend.log"

(cd backend && ENV=local make run-local) >>logs/backend.log 2>&1 &
back_pid=$!
(set +u; cd frontend && nvm use && npm run dev) >>logs/frontend.log 2>&1 &
front_pid=$!

wait
