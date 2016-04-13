
fetch_reno will fetch reno of openstack project
-----------------------------------------------

build
------
1. install golang
2. go build main.go


configure
---------
1. add project repo name to proj.txt
$ cat proj.txt
nova
neutron
magnum
cinder
glance
ironic
keystone

run
---
./get_reno.sh

reno update results will be saved into a txt file named as today.

notes
-----
.data/$repo will save last commit SHA of releatesnotes/notes of a $reop
by default, the first time it will query recent 10 day's commits on releatenotes/notes
of a $repo


know issues
-----------
some time git hub will return 403 because we don't use authory mode. server
refused us.
