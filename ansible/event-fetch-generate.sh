#!/bin/sh

set -e

echo "Fetching event inventory from Ambrose's infrastructure"
ssh saguaro /usr/local/bin/wireguard-negotiator dump -i wg1 | go run event-generate.go opentech > event
