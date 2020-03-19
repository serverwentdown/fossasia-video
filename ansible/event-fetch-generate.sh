#!/bin/sh

set -e

echo "Fetching event inventory from Ambrose's infrastructure"
ssh -p 222 root@fd11:f055:a514::1 /usr/local/bin/wireguard-negotiator dump -i wg1 | go run event-generate.go opentech
