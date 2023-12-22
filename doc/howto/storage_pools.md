(howto-storage-pools)=
# ストレージプールを管理するには

{ref}`storage-pools` を作成、設定、表示、リサイズするための手順については以下のセクションを参照してください。

(storage-create-pool)=
## ストレージプールを作成する

Incus は初期化中にストレージプールを作成します。
同じドライバーあるいは別のドライバーを使用して、後からさらにストレージプールを追加できます。

ストレージプールを作成するには以下のコマンドを使用します:

    incus storage create <pool_name> <driver> [configuration_options...]

別途指定しない場合は、 Incus は実用的なデフォルトのサイズ（空きディスクスペースの 20%、しかし最低 5GiB で最大 30GiB）でループベースのストレージをセットアップします。

それぞれのドライバーで利用可能な設定オプションの一覧は {ref}`storage-drivers` ドキュメントを参照してください。

### 例

それぞれのストレージドライバーでストレージプールを作成する例は以下を参照してください。

`````{tabs}

````{group-tab} ディレクトリ

`pool1` という名前のディレクトリプールを作成する:

    incus storage create pool1 dir

`/data/incus` という既存のディレクトリを使って `pool2` を作成する:

    incus storage create pool2 dir source=/data/incus
````
````{group-tab} Btrfs

`pool1` という名前のループバックプールを作成する:

    incus storage create pool1 btrfs

`/some/path` にある既存の Btrfs ファイルシステムを使って `pool2` を作成する:

    incus storage create pool2 btrfs source=/some/path

`/dev/sdX` 上に `pool3` という名前のプールを作成する:

    incus storage create pool3 btrfs source=/dev/sdX
````
````{group-tab} LVM

`pool1` という名前のループバックのプールを作成する（LVM ボリュームグループ名も `pool1` になります）:

    incus storage create pool1 lvm

`my-pool` という既存の LVM ボリュームグループを使って `pool2` を作成する:

    incus storage create pool2 lvm source=my-pool

`my-vg` というボリュームグループ内の `my-pool` という既存の LVM thin-pool を使って `pool3` を作成する:

    incus storage create pool3 lvm source=my-vg lvm.thinpool_name=my-pool

`/dev/sdX` 上に `pool4` という名前のプールを作成する（LVM ボリュームグループ名も `pool4` になります）:

    incus storage create pool4 lvm source=/dev/sdX

`/dev/sdX` 上に `my-pool` というLVM ボリュームグループ名で `pool5` という名前のプールを作成する:

    incus storage create pool5 lvm source=/dev/sdX lvm.vg_name=my-pool
````
````{group-tab} ZFS

`pool1` という名前のループバックプールを作成する（ZFS zpool 名も `pool1` になります）:

    incus storage create pool1 zfs

`pool2` という名前のループバックプールを `my-tank` という ZFS zpool 名で作成する:

    incus storage create pool2 zfs zfs.pool_name=my-tank

`my-tank` という既存の ZFS zpool を使用して `pool3` を作成する:

    incus storage create pool3 zfs source=my-tank

`my-tank/slice` という既存の ZFS データセットを使用して `pool4` を作成する:

    incus storage create pool4 zfs source=my-tank/slice

`/dev/sdX` 上に `pool5` という名前のプールを作成する（ZFS zpool 名も `pool5` になります）:

    incus storage create pool5 zfs source=my-tank/zvol volume.zfs.block_mode=yes

`/dev/sdX` 上に `pool6` という名前のプールを作成する（ZFS zpool 名も `pool6` になります）:

    incus storage create pool6 zfs source=/dev/sdX

`/dev/sdX` 上に `my-tank` という ZFS zpool 名で `pool7` という名前のプールを作成する:

    incus storage create pool7 zfs source=/dev/sdX zfs.pool_name=my-tank
````
````{group-tab} Ceph RBD

デフォルトの Ceph クラスタ（名前は `ceph`）内に `pool1` という名前の OSD ストレージプールを作成する:

    incus storage create pool1 ceph

`my-cluster` という Ceph クラスタ内に `pool2` という名前の OSD ストレージプールを作成する:

    incus storage create pool2 ceph ceph.cluster_name=my-cluster

デフォルトの Ceph クラスタ内に `my-osd` という on-disk 名で `pool3` という名前の OSD ストレージプールを作成する:

    incus storage create pool3 ceph ceph.osd.pool_name=my-osd

`my-already-existing-osd` という既存の OSD ストレージプールを使って `pool4` を作成する:

    incus storage create pool4 ceph source=my-already-existing-osd

`ecpool` という既存の OSD erasure-coded ストレージプールと `rpl-pool` という OSD リプリケーテッドプールを使って `pool5` を作成する:

    incus storage create pool5 ceph source=rpl-pool ceph.osd.data_pool_name=ecpool
````
````{group-tab} CephFS

```{note}
CephFS ファイルシステムは 2 つの OSD ストレージプールからなります。そのうち 1 つは実際のデータ、もう 1 つはファイルメタデータに使用されます。
```

既存の CephFS ファイルシステム `my-filesystem` を使って `pool1` を作成する:

    incus storage create pool1 cephfs source=my-filesystem

`my-filesystem` ファイルシステムからサブディレクトリ `my-directory` を使って `pool2` を作成する:

    incus storage create pool2 cephfs source=my-filesystem/my-directory

`my-data` という名前のデータプールと`my-metadata`という名前のメタデータプールを持つ CephFS ファイルシステム `my-filesystem` を使用して `pool3` を作成する:

    incus storage create pool3 cephfs source=my-filesystem cephfs.create_missing=true cephfs.data_pool=my-data cephfs.meta_pool=my-metadata

````
````{group-tab} Ceph Object

```{note}
Ceph Object ドライバを使用する場合、事前に稼働中の Ceph Object Gateway [`radosgw`](https://docs.ceph.com/en/latest/radosgw/) の URL を用意しておく必要があります。
```

既存の Ceph Object Gateway `https://www.example.com/radosgw` を使用して `pool1` を作成する:

    incus storage create pool1 cephobject cephobject.radosgw.endpoint=https://www.example.com/radosgw

````
`````

(storage-pools-cluster)=
### クラスタ内にストレージプールを作成する

Incus クラスタを稼働していてストレージプールを追加したい場合、それぞれのクラスタメンバー内にストレージを別々に作る必要があります。
この理由は、設定、たとえばストレージのロケーションやプールのサイズがクラスタメンバー間で異なるかもしれないからです。

このため、 `--target=<cluster_member>` フラグを指定してストレージプールをペンディング状態でまず作成し、メンバーごとに適切な設定を行う必要があります。
すべてのメンバーで同じストレージプール名を使用しているか確認してください。
次に `--target` フラグなしでストレージプールを作成し、実際にセットアップします。

たとえば、以下の一連のコマンドは 3 つのクラスタメンバー上で異なるロケーションと異なるサイズで `my-pool` という名前のストレージプールをセットアップします:

```{terminal}
:input: incus storage create my-pool zfs source=/dev/sdX size=10GiB --target=vm01
Storage pool my-pool pending on member vm01
:input: incus storage create my-pool zfs source=/dev/sdX size=15GiB --target=vm02
Storage pool my-pool pending on member vm02
:input: incus storage create my-pool zfs source=/dev/sdY size=10GiB --target=vm03
Storage pool my-pool pending on member vm03
:input: incus storage create my-pool zfs
Storage pool my-pool created
```

{ref}`cluster-config-storage`も参照してください。

```{note}
ほとんどのストレージドライバでは、ストレージプールは各クラスタメンバー上にローカルに存在します。
これは 1 つのメンバー上のストレージプール内にストレージボリュームを作成しても、別のクラスタメンバー上では利用可能にはならないことを意味します。

この挙動は Ceph ベースのストレージプール（`ceph`、 `cephfs`、 `cephobject`）では異なります。これらではストレージプールは 1 つの中央のロケーション上に存在し、全てのクラスタメンバーが同じストレージボリュームを持つ同じストレージプールにアクセスします。
```

## ストレージプールを設定する

各ストレージドライバーで利用可能な設定オプションについては {ref}`storage-drivers` ドキュメントを参照してください。

（`source` のような）ストレージプールの一般的なキーはトップレベルです。
ドライバー固有のキーはドライバー名で名前空間が分けられています。

ストレージプールに設定オプションを設定するには以下のコマンドを使用します:

    incus storage set <pool_name> <key> <value>

たとえば、 `dir` ストレージプールでストレージプールのマイグレーション中に圧縮をオフにするには以下のコマンドを使用します:

    incus storage set my-dir-pool rsync.compression false

ストレージプールの設定を編集するには以下のコマンドを使用します:

    incus storage edit <pool_name>

## ストレージプールを表示する

すべての利用可能なストレージプールの一覧を表示し設定を確認できます。

以下のコマンドですべての利用可能なストレージプールを一覧表示できます:

    incus storage list

出力結果の表には（訳注: Incus の）初期化時に作成した（通常 `default` や `local` と呼ばれる）ストレージプールとあなたが追加したあらゆるストレージプールが含まれます。

特定のプールに関する詳細情報を表示するには、以下のコマンドを使用します:

    incus storage show <pool_name>

特定のプールに関する使用量を表示するには、以下のコマンドを使用します:

    incus storage info <pool_name>

(storage-resize-pool)=
## ストレージプールをリサイズする

ストレージがもっと必要な場合、`size` 設定キーを変更することでストレージプールのサイズを拡大できます:

    incus storage set <pool_name> size=<new_size>

これはループファイルをバックエンドとし Incus で管理されているストレージプールでのみ機能します。
プールは拡張（サイズを増やす）のみが可能で、縮小はできません。
