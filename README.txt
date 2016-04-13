
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
./get_reno.sh <day> # to get release notes of repos in proj.txt in <day> days
reno update results will be saved into a txt file named as today.

for example ::

2016年04月13日

or

./main nova 1 #will show update notes of recent 1 day



notes
-----
.data/$repo will save last commit SHA of releatesnotes/notes of a $reop
by default, the first time it will query recent 10 day's commits on releatenotes/notes
of a $repo


know issues
-----------
some time git hub will return 403 because we don't use authory mode. server
refused us.
