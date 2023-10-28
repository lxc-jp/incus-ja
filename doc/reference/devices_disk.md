(devices-disk)=
# タイプ: `disk`

```{note}
`disk` デバイスタイプはコンテナと VM の両方でサポートされます。
コンテナと VM の両方でホットプラグをサポートします。
```

ディスクデバイスはインスタンスに追加のストレージを提供します。

コンテナにとっては、それらはインスタンス内の実質的なマウントポイントです（ホスト上の既存のファイルまたはディレクトリーのバインドマウントとしてか、あるいは、ソースがブロックデバイスの場合は通常のマウントのマウントポイント）。
仮想マシンは `9p` または `virtiofs`（使用可能な場合）を通してホスト側のマウントまたはディレクトリーを共有するか、あるいはブロックベースのディスクに対する VirtIO ディスクとして共有します。

(devices-disk-types)=
##  ディスクデバイスの種類

さまざまなソースからディスクデバイスを作成できます。  
`source` オプションに指定する値によって、追加されるディスクデバイスのタイプが決まります:

ストレージボリューム  
: 最も一般的なタイプのディスクデバイスはストレージボリュームです。  
  ストレージボリュームを追加するには、デバイスの`source`としてその名前を指定します:

      incus config device add <instance_name> <device_name> disk pool=<pool_name> source=<volume_name> [path=<path_in_instance>]

  path はファイルシステムボリュームには必要ですが、ブロックボリュームには必要ありません。

  また、[`incus storage volume attach`](incus_storage_volume_attach.md) コマンドを使用して{ref}`storage-attach-volume`することもできます。  
  どちらのコマンドも、ストレージボリュームをディスクデバイスとして追加するための同じメカニズムを使用します。

ホスト上のパス
: ホストのパス（ファイルシステムまたはブロックデバイスのいずれか）をインスタンスと共有するには、ディスクデバイスとして追加し、`source`としてホストパスを指定します:

      incus config device add <instance_name> <device_name> disk source=<path_on_host> [path=<path_in_instance>]

  path はファイルシステムボリュームには必要ですが、ブロックデバイスには必要ありません。

Ceph RBD
: Incus は、インスタンスの内部ファイルシステムを管理するために Ceph を使用できますが、既存の外部管理 Ceph RBD をインスタンスに使用したい場合は、次のコマンドで追加できます:

      incus config device add <instance_name> <device_name> disk source=ceph:<pool_name>/<volume_name> ceph.user_name=<user_name> ceph.cluster_name=<cluster_name> [path=<path_in_instance>]

  path はファイルシステムボリュームには必要ですが、ブロックデバイスには必要ありません。

CephFS
: Incus はインスタンスで内部のファイルシステムの管理に Ceph を使えますが、既存の外部で管理されている Ceph ファイルシステムをインスタンスで使用したい場合は、以下のコマンドで追加できます:

      incus config device add <instance_name> <device_name> disk source=cephfs:<fs_name>/<path> ceph.user_name=<user_name> ceph.cluster_name=<cluster_name> path=<path_in_instance>

ISO file
: 仮想マシンには ISO ファイルをディスクデバイスとして追加できます。
  ISO ファイルは VM 内部の ROM デバイスとして追加されます。

  このソースタイプは VM でのみ利用可能です。

  ISO ファイルを追加するには、そのファイルパスを`source`として指定します:

      incus config device add <instance_name> <device_name> disk source=<file_path_on_host>

VM `cloud-init`
: {config:option}`instance-cloud-init:cloud-init.vendor-data`、{config:option}`instance-cloud-init:cloud-init.user-data`から`cloud-init`設定の ISO イメージを生成し、仮想マシンにアタッチできます。

  このソースタイプは VM でのみ利用可能です。

  そのようなデバイスを追加するには、以下のコマンドを使用します:

      incus config device add <instance_name> <device_name> disk source=cloud-init:config

(devices-disk-initial-config)=
## インスタンスルートディスクデバイスの初期ボリューム設定

初期ボリューム設定は新しいインスタンスのルートディスクデバイスに個別の設定をできるようにします。
これらの設定は `initial.` という接頭辞がつき、インスタンスが作成されたときのみ適用されます。
この方法はデフォルトのストレージプールの設定とは独立に、個別の設定を持つインスタンスを作れるようにします。

たとえば、既存のプロファイルに `zfs.block_mode` の初期ボリューム設定を追加し、このプロファイルを使ってインスタンスを作成する都度適用できます:

    incus profile device set <profile_name> <device_name> initial.zfs.block_mode=true

インスタンス作成時に初期設定を直接指定もできます。たとえば:

    incus init <image> <instance_name> --device <device_name>,initial.zfs.block_mode=true

カスタムボリュームオプションに初期ボリューム設定を使ったりボリュームのサイズを設定はできないことに注意してください。

## デバイスオプション

`disk` デバイスには以下のデバイスオプションがあります:

キー                | 型      | デフォルト値  | 必須 | 説明
:--                 | :--     | :--           | :--  | :--
`boot.priority`     | integer | -             | no   | VM のブート優先度（高いほうが先にブート）
`ceph.cluster_name` | string  | `ceph`        | no   | Ceph クラスタのクラスタ名（Ceph か CephFS のソースには必須）
`ceph.user_name`    | string  | `admin`       | no   | Ceph クラスタのユーザー名（Ceph か CephFS のソースには必須）
`initial.*`         | n/a     | -             | no   | デフォルトストレージプール設定と独立して個別の設定を許可する{ref}`devices-disk-initial-config`
`io.bus`            | string  | `virtio-scsi` | no   | VMのみ: デバイスのバスを上書きする（`virtio-scsi`または`nvme`）
`io.cache`          | string  | `none`        | no   | VMのみ: デバイスのキャッシュモードを上書きする（`none`、`writeback`または`unsafe`）
`limits.max`        | string  | -             | no   | 読み取りと書き込み両方のbyte/sかIOPSによるI/O制限（`limits.read`と`limits.write`の両方を設定するのと同じ）
`limits.read`       | string  | -             | no   | byte/s（さまざまな単位が使用可能、{ref}`instances-limit-units`参照）もしくはIOPS（あとに`iops`と付けなければなりません）で指定する読み込みのI/O制限値 - {ref}`storage-configure-IO` も参照
`limits.write`      | string  | -             | no   | byte/s（さまざまな単位が使用可能、{ref}`instances-limit-units`参照）もしくはIOPS（あとに`iops`と付けなければなりません）で指定する書き込みのI/O制限値 - {ref}`storage-configure-IO` も参照
`path`              | string  | -             | yes  | ディスクをマウントするインスタンス内のパス（コンテナのみ）
`pool`              | string  | -             | no   | ディスクデバイスが属するストレージプール（Incus が管理するストレージボリュームにのみ適用可能）
`propagation`       | string  | -             | no   | バインドマウントをインスタンスとホストでどのように共有するかを管理する（`private` （デフォルト）、`shared`、`slave`、`unbindable`、 `rshared`、`rslave`、`runbindable`、`rprivate` のいずれか。完全な説明は Linux Kernel の文書 [shared subtree](https://www.kernel.org/doc/Documentation/filesystems/sharedsubtree.txt) 参照） <!-- wokeignore:rule=slave -->
`raw.mount.options` | string  | -             | no   | ファイルシステム固有のマウントオプション
`readonly`          | bool    | `false`       | no   | マウントを読み込み専用とするかどうかを制御
`recursive`         | bool    | `false`       | no   | ソースパスを再帰的にマウントするかどうかを制御
`required`          | bool    | `true`        | no   | ソースが存在しないときに失敗とするかどうかを制御
`shift`             | bool    | `false`       | no   | ソースの UID/GID をインスタンスにマッチするように変換させるためにオーバーレイの shift を設定するか（コンテナのみ）
`size`              | string  | -             | no   | byte（さまざまな単位が使用可能、 {ref}`instances-limit-units` 参照）で指定するディスクサイズ。`rootfs`（`/`） でのみサポートされます
`size.state`        | string  | -             | no   | 上の `size` と同じですが、仮想マシン内のランタイム状態を保存するために使われるファイルシステムボリュームに適用されます
`source`            | string  | -             | yes  | ファイルシステムまたはブロックデバイスのソース（詳細は{ref}`devices-disk-types`参照）
