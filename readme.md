
# HDM

HDM is a software that do Hard Disk Management.

It handle physical and logical disks lifecycle throughout servers.

Disks States: 
- `list`: List disks Info (serial, size, labels, form factor, ...) for plugged and unplugged disks
- `index`: Index files location, size and date 
- `search`: Search for files, even on unplugged disks
- `location`: Get disk physical location 
- `forget` : Remove all stored info about a disk and his files

Add/Remove disks:
- `prepare` : Prepare a new disks (partition, encrypt, format, mount) (requires no partitions)
- `add` : Add plugged-in disks for usage (mdadm, luksOpen, mount, restart)
- `remove`: Pre or post disk unplug actions (stop, kill, umount, luksClose, spin-down, sleep)
- `erase`: Erase entire disk (requires no partitions)

Disk health :
- `health`: Report disks health 
- `test`: Test disk healthiness
- `repair`: Repair pending blocks

Disk saving/sync:
- `backupable`: Check directories can be backup (target size, target plugged)
- `backup`: backup directories
- `backups`: list backups
- `restore`: restore a file from backup, also used by `repair`

Global lifecycle:
- `check`: visit all checks commands to ensure everything is ok
- `agent`: Run a agent that inotify and self handle all possible commands and ask for human intervention


## Current requirements

- uniq label for each partitions of all disks if using labels as id -> identify disks
- can ssh to servers with ssh agent -> run hdm outside of servers
- can ssh from server to servers -> run rsync cross servers
- can run sudo on servers without password -> run as root
- lsblk >= 2.33 detect disks
- smartctl >= 7.0 -> disk health
- hdparm -> put disk in sleep
- rsync -> all backups commands
- udev -> list disks by path

## Install

HDM is a single binary file, just download and extract it to any bin directory

## Usage


## TODO

- find non backed-up paths
- get disk name for a file in any filesystem (links, overlays)
- sync directories across servers
- remove without selector should find mounted removed devices
- prepare by location
- periodic set to readonly
- save last backup time so we know we should do it again
- put sas disk in sleep mode
- get disk name from location
- refuse to prepare a new disk if label is already sued by another device
- report conflicts in labels
- agent: inotify any file change: run backup, re-index
- agent: trigger notification to operator: disk needs to be plugged, size of directory cannot be backuped, disk has failure





