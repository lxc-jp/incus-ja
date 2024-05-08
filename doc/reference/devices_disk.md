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

VM `agent`
: エージェントの実行ファイル、設定ファイル、インストールスクリプトを含む `agent` 設定の ISO を生成できます。
  これは `9p` が非サポートでエージェントをロードする別の方法が必要な環境で必要です。

  このソースタイプは VM でのみ利用可能です。

  そのようなデバイスを追加するには、以下のコマンドを使用します:

      incus config device add <instance_name> <device_name> disk source=agent:config

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

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-disk start -->
    :end-before: <!-- config group devices-disk end -->
```
