
fetch_reno will fetch reno of openstack project
-----------------------------------------------

build
------
1. install golang
2. go build fetchreno.go


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

./fetchreno nova 1 #will show update notes of recent 1 day

$ ./fetchreno nova 5
2016/04/13 21:55:57 Time period:  2016-04-08 21:55:57.803281469 +0800 CST 2016-04-13 21:55:57.803281469 +0800 CST
2016/04/13 21:55:59 Last commit length is: 2
2016/04/13 21:55:59 commit : a03fc2325ccf7bc2df3a5f5d926495a0fde8d66f
2016/04/13 21:55:59 commit : f738483e843fc27379b85c5401859ccc854adc5e
------------------------------------[updates releasenotes for nova]---------------------------------
---------[from 2016-04-08 21:55:57.803281469 +0800 CST to 2016-04-13 21:56:02.422824237 +0800 CST ]---------
releasenotes/notes/swap-volume-policy-9464e97aba12d1e0.yaml
---
upgrade:
  - The default policy for updating volume attachments, commonly referred to as
    swap volume, has been changed from ``rule:admin_or_owner`` to
    ``rule:admin_api``. This is because it is called from the volume service
    when migrating volumes, which is an admin-only operation by default, and
    requires calling an admin-only API in the volume service upon completion.
    So by default it would not work for non-admins.

releasenotes/notes/swap-volume-policy-9464e97aba12d1e0.yaml
---
upgrade:
  - The default policy for updating volume attachments, commonly referred to as
    swap volume, has been changed from ``rule:admin_or_owner`` to
    ``rule:admin_api``. This is because it is called from the volume service
    when migrating volumes, which is an admin-only operation by default, and
    requires calling an admin-only API in the volume service upon completion.
    So by default it would not work for non-admins.


notes
-----
.data/$repo will save last commit SHA of releatesnotes/notes of a $reop
by default, the first time it will query recent 10 day's commits on releatenotes/notes
of a $repo


know issues
-----------
some time git hub will return 403 because we don't use authory mode. server
refused us.
