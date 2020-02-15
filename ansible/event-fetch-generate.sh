#!/bin/sh

set -e

echo "Fetching event inventory"

ssh saguaro /usr/local/bin/wireguard-negotiator dump -i wg1 | go run event-generate.go opentech > event
