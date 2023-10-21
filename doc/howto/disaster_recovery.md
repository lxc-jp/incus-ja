(disaster-recovery)=
# 災害時にインスタンスを復旧するには

Incus は{ref}`Incus データベース <database>`が壊れたり失われたりといったディザスタリカバリのためのツールを提供しています。

ツールはインスタンスのストレージプールをスキャンし、見つけたインスタンスをデータベースにインポートします。
失われた必要なエンティティ（通常はプロファイル、プロジェクトとネットワーク）をあなたが再作成する必要があります。

```{important}
このツールはディザスタリカバリのみに使うべきです。
適切なバックアップの代替手段としてこのツールに依存しないでください。プロファイル、ネットワークの定義、サーバ設定のようなデータは失われてしまうからです。

ツールはインタラクティブに実行する必要があり、自動化したスクリプト内では使用できません。
```

## リカバリプロセス

ツールを実行すると、データベースにまだ残っているすべてのストレージプールをスキャンし、復旧できる失われたボリュームを探します。
（ディスク上に存在するがデータベース内には存在しない）未知のストレージプールの詳細を指定することもでき、するとツールはそれらもスキャンを試みます。

（まだマウントされていなければ）指定されたストレージプールをマウントした後、ツールは Incus に関連付けられていたと思われる未知のボリュームのストレージボリュームをスキャンします。
Incus は各インスタンスのストレージボリュームに`backup.yaml`を保管していて、そこにインスタンスを復旧するために必要なすべての（インスタンス設定、アタッチしたデバイス、ストレージボリューム、プール設定も含めた）情報を保管しています。
このデータはインスタンス、ストレージボリューム、そしてストレージプールのデータベースレコードをリビルドするのに使用できます。
インスタンスを復旧する前に、ツールは`backup.yaml` ファイルの内容と（対応するスナップショットなど）ディスク上に実際に存在するものとを比較してある程度の整合性チェックを行います。
問題なければデータベースのレコードが再生成されます。

ストレージプールのデータベースレコードも作成が必要な場合、ディスカバリーフェーズにユーザーが入力した情報よりも、インスタンスの `backup.yaml` ファイルを設定のベースとして優先して使用します。
ただし、それが無い場合はユーザーが入力した情報をもとにプールのデータベースレコードを復元するようにフォールバックします。

ツールはネットワークなど失われたエンティティを再生成するためにあなたに質問します。
しかし、ツールはどのようにインスタンスが設定されていたかを知りません。
これはつまり一部の設定が`default`プロファイル経由で指定されていた場合、プロファイルに必要な設定をあなたが再度追加する必要があることを意味します。
たとえば、インスタンスで`incusbr0`ブリッジが使われていてそれを再生成するようプロンプトが出た場合、復旧されるインスタンスがそれを使うようにあなたはそれを`default`プロファイルに追加し直す必要があります。

## 例

リカバリプロセスの例を示します:

```{terminal}
:input: incus admin recover

This Incus server currently has the following storage pools:
Would you like to recover another storage pool? (yes/no) [default=no]: yes
Name of the storage pool: default
Name of the storage backend (btrfs, ceph, cephfs, cephobject, dir, lvm, zfs): zfs
Source of the storage pool (block device, volume group, dataset, path, ... as applicable): /var/lib/incus/storage-pools/default/containers
Additional storage pool configuration property (KEY=VALUE, empty when done): zfs.pool_name=default
Additional storage pool configuration property (KEY=VALUE, empty when done):
Would you like to recover another storage pool? (yes/no) [default=no]:
The recovery process will be scanning the following storage pools:
 - NEW: "default" (backend="zfs", source="/var/lib/incus/storage-pools/default/containers")
Would you like to continue with scanning for lost volumes? (yes/no) [default=yes]: yes
Scanning for unknown volumes...
The following unknown volumes have been found:
 - Container "u1" on pool "default" in project "default" (includes 0 snapshots)
 - Container "u2" on pool "default" in project "default" (includes 0 snapshots)
You are currently missing the following:
 - Network "incusbr0" in project "default"
Please create those missing entries and then hit ENTER: ^Z
[1]+  Stopped                 incus admin recover
:input: incus network create incusbr0
Network incusbr0 created
:input: fg
incus admin recover

The following unknown volumes have been found:
 - Container "u1" on pool "default" in project "default" (includes 0 snapshots)
 - Container "u2" on pool "default" in project "default" (includes 0 snapshots)
Would you like those to be recovered? (yes/no) [default=no]: yes
Starting recovery...
:input: incus list
+------+---------+------+------+-----------+-----------+
| NAME |  STATE  | IPV4 | IPV6 |   TYPE    | SNAPSHOTS |
+------+---------+------+------+-----------+-----------+
| u1   | STOPPED |      |      | CONTAINER | 0         |
+------+---------+------+------+-----------+-----------+
| u2   | STOPPED |      |      | CONTAINER | 0         |
+------+---------+------+------+-----------+-----------+
:input: incus profile device add default eth0 nic network=incusbr0 name=eth0
Device eth0 added to default
:input: incus start u1
:input: incus list
+------+---------+-------------------+---------------------------------------------+-----------+-----------+
| NAME |  STATE  |       IPV4        |                    IPV6                     |   TYPE    | SNAPSHOTS |
+------+---------+-------------------+---------------------------------------------+-----------+-----------+
| u1   | RUNNING | 192.0.2.49 (eth0) | 2001:db8:8b6:abfe:216:3eff:fe82:918e (eth0) | CONTAINER | 0         |
+------+---------+-------------------+---------------------------------------------+-----------+-----------+
| u2   | STOPPED |                   |                                             | CONTAINER | 0         |
+------+---------+-------------------+---------------------------------------------+-----------+-----------+
```
