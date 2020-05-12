
# fossasia-video

The experimental FOSSASIA video recording setup

Looking for the video recorder's guide? The collaborative version is [on Google Docs](https://docs.google.com/document/d/1LWKt0kXWGVuPHCLvBBH35frJFtfhQuNIZ7DSS2HNLZI/edit) and a [PDF is available in this repository](cheatsheet.pdf).

There is also the [video equipment sheet](https://docs.google.com/spreadsheets/d/1XLckJRG4ng2NoOfqy-EKca2LJdfC-Z6zfE6Q3FqKoM0/edit)

## Overview

## Installing Debian

For the recording machines, get a fresh copy of Debian, and install it with the following settings:

- Keyboard: American
- Username: opentech
- Hostname: model-increment
  - Example: x220-01
- Add GNOME Desktop
- Add OpenSSH Server
- Leave Print & System Utilities enabled

If a GNOME Desktop and SSH daemon is already installed, a reinstall is not required but recommended. 

## WireGuard

To set up the overlay WireGuard network to manage machines everywhere, a server needs to be set up. [`wireguard-negotiator`](https://github.com/serverwentdown/wireguard-negotiator) is a tool written to automate some of the key exchange, that must be run on a publicly accessible server as root:

```
# Create WireGuard interface
ip link add dev wg1 type wireguard
ip addr add fd11:f055:a514:0000::1/64 dev wg1
# Configure the interface once. Assumes the configuration file exists
# See WireGuard docs on how to write this configuraion file
wg setconf wg1 /etc/wireguard/wg1.conf
# Bring up the interface
ip link set wg1 up
wireguard-negotiator server -i wg1 -e [hostname] -l :8080 -I -B
```

On every recording machine, the client can be easily configured like so:

```
# For now, we have to manually install WireGuard
echo "deb http://deb.debian.org/debian/ unstable main" | sudo tee /etc/apt/sources.list.d/unstable.list
printf 'Package: *\nPin: release a=unstable\nPin-Priority: 90\n' | sudo tee /etc/apt/preferences.d/limit-unstable
sudo apt update
sudo apt install wireguard
# TODO: validate this section
sudo systemctl enable --now systemd-networkd
wget -O wgn http://[hostname]:8080
chmod +x wgn
sudo ./wgn request -s http://[hostname]:8080
```

## Configuring Rooms

To specify the room for a specific host, do the following:

```
echo the_room_id > ~opentech/room_id
echo teh_room_type > ~opentech/room_type
# Example for Event Hall 2-1:
echo EH2 > ~opentech/room_type
echo EH2-1 > ~opentech/room_id
```

See the [Video Equipment Google Sheet](https://docs.google.com/spreadsheets/d/1XLckJRG4ng2NoOfqy-EKca2LJdfC-Z6zfE6Q3FqKoM0/edit?usp=sharing) for the specific IDs. To get edit access, DM @serverwentdown on Telegram. 

This step is optional. This and the following steps should be done for every change in room or setup of the laptop.

## Exporting Hosts

Export all hosts on the overlay network from WireGuard configuration:

```
wireguard-negotiator dump > wireguard.list
```

This list of IPs can be exported into our Ansible hosts format as such:

```
go run event-generate.go opentech < wireguard.list > event
```

The script `event-fetch-generate.sh` does the exporting and generating of the `event` inventory.

## Running Playbooks

Before running any Playbooks, switch to SSH authentication:

```
ansible-playbook -Kf 8 -kc paramiko -i event ssh-key.yml
ansible -Kf 8 -kc paramiko -b -i event all -a reboot
```

Now, plays can be run like so:

```
ansible-playbook -Kf 4 -i [inventory] [playbook]
# Example:
ansible-playbook -Kf 4 -i event recorders.yml
ansible-playbook -Kf 4 -i event recorders-stop.yml
ansible-playbook -Kf 4 -i event reboot.yml
ansible-playbook -Kf 4 -i event upgrade.yml
```

And run commands like so:

```
ansible -Kf 8 -b -i event all -a 'apt install firmware-iwlwifi'
```

