(server-migrate-lxd)=
# LXD からの移行

Incus は `lxd-to-incus` という名前のツールがあり、既存の LXD 環境を Incus に変換するのに使えます。

これが正しく動作するためには、Incus の最新の安定版をインストールし、しかし初期化はしないようにしてください。
代わりに、`incus info` と `lxc info` の両方が正常に動くことを確認し、次に `lxd-to-incus` を動かしてあなたのデータを移行してください。

この手順はデータベース全体と全てのストレージを LXD から Incus に移行し、移行後は同一のセットアップになります。

```{terminal}
:input: lxd-to-incus
:user: root
=> Looking for source server
==> Detected: snap package
=> Looking for target server
=> Connecting to source server
=> Connecting to the target server
=> Checking server versions
==> Source version: 5.19
==> Target version: 0.1
=> Validating version compatibility
=> Checking that the source server isn't empty
=> Checking that the target server is empty
=> Validating source server configuration

The migration is now ready to proceed.
At this point, the source server and all its instances will be stopped.
Instances will come back online once the migration is complete.

Proceed with the migration? [default=no]: yes
=> Stopping the source server
=> Stopping the target server
=> Wiping the target server
=> Migrating the data
=> Migrating database
=> Cleaning up target paths
=> Starting the target server
=> Checking the target server
Uninstall the LXD package? [default=no]: yes
=> Uninstalling the source server
```

```{terminal}
:input: incus list
:user: root
To start your first container, try: incus launch images:ubuntu/22.04
Or for a virtual machine: incus launch images:ubuntu/22.04 --vm

+------+---------+-----------------------+-----------------------------------------------+-----------+-----------+
| NAME |  STATE  |         IPV4          |                     IPV6                      |   TYPE    | SNAPSHOTS |
+------+---------+-----------------------+-----------------------------------------------+-----------+-----------+
| u1   | RUNNING | 10.204.220.101 (eth0) | fd42:1eb6:f1d8:4e2a:216:3eff:fe65:940d (eth0) | CONTAINER | 0         |
+------+---------+-----------------------+-----------------------------------------------+-----------+-----------+
```

また、ツールは Incus に非互換な設定を検索し、もし存在すればデータを移行する前に中止します。

```{warning}
移行中は全てのインスタンスが停止されます。
移行プロセスが一旦開始されると、簡単には後戻りできませんので、適切なダウンタイムを見込んで計画するようにしてください。
```
