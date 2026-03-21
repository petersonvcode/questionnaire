#!/usr/bin/env bash
# Inline CSS and JS from dist/assets into dist/index.html

set -euo pipefail

DIST_DIR="${1:-$(dirname "$0")/dist}"
HTML_FILE="$DIST_DIR/index.html"
ASSETS_DIR="$DIST_DIR/assets"

[[ -f "$HTML_FILE" ]] || { echo "Error: $HTML_FILE not found"; exit 1; }
[[ -d "$ASSETS_DIR" ]] || { echo "Error: $ASSETS_DIR not found"; exit 1; }

# Create temp file for output
TMP_HTML=$(mktemp)
trap 'rm -f "$TMP_HTML"' EXIT

css_pat='<link[^>]*rel="stylesheet"[^>]*href="(/assets/[^"]+\.css)"'
js_pat='<script([^>]*)src="(/assets/[^"]+\.js)"([^>]*)>'

while IFS= read -r line || [[ -n "$line" ]]; do
  # Check for stylesheet link
  if [[ "$line" =~ $css_pat ]]; then
    css_file="$DIST_DIR${BASH_REMATCH[1]}"
    if [[ -f "$css_file" ]]; then
      echo '    <style>'
      # Escape </style> so it doesn't close the tag
      sed 's|</style>|<\\/style>|g' < "$css_file"
      echo '    </style>'
      rm "$css_file"
    else
      echo "$line"
    fi
  # Check for script with src
  elif [[ "$line" =~ $js_pat ]]; then
    js_file="$DIST_DIR${BASH_REMATCH[2]}"
    if [[ -f "$js_file" ]]; then
      # Preserve type="module" or other attributes if present
      attrs="${BASH_REMATCH[1]}${BASH_REMATCH[3]}"
      echo "    <script$attrs>"
      # Escape </script> so it doesn't close the tag
      sed 's|</script>|<\\/script>|g' < "$js_file"
      echo '    </script>'
      rm "$js_file"
    else
      echo "$line"
    fi
  else
    echo "$line"
  fi
done < "$HTML_FILE" > "$TMP_HTML"

mv "$TMP_HTML" "$HTML_FILE"
echo "Done: inlined assets into $HTML_FILE"
