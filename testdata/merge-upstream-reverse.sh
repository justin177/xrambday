#!/bin/sh
set -eu

if [ "$#" -lt 1 ] || [ "$#" -gt 3 ]; then
  echo "usage: $0 <upstream-config.json> [merge-fragment.json] [output.json]" >&2
  exit 2
fi

script_dir="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"

base_config="$1"
merge_config="${2:-$script_dir/upstream-reverse-portal.merge.json}"
output_config="${3:-upstream.merged.json}"

jq -s '
  .[0] as $base |
  .[1] as $add |
  $base
  | .reverse = (.reverse // {})
  | .reverse.portals = ((.reverse.portals // []) + ($add.reverse.portals // []))
  | .inbounds = ((.inbounds // []) + ($add.inbounds // []))
  | .inbounds[0].settings = (.inbounds[0].settings // {})
  | .inbounds[0].settings.fallbacks = (($add.fallbacks // []) + (.inbounds[0].settings.fallbacks // []))
  | .routing = (.routing // {})
  | .routing.rules = (($add.routing.rules // []) + (.routing.rules // []))
' "$base_config" "$merge_config" > "$output_config"

echo "merged config written to $output_config"
