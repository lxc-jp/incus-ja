(storage-lvm)=
# LVM - `lvm`

{abbr}`LVM (Logical Volume Manager)` はファイルシステムというよりストレージマネージメントフレームワークです。
これは物理ストレージデバイスを管理するのに使用され、複数のストレージボリュームを作成し、配下の物理ストレージデバイスを使用し仮想化できるようにします。

この過程で物理ストレージをオーバーコミットすることが可能で、すべての利用可能なストレージが同時に使用されるわけではないシナリオに対して柔軟性を提供できることに注意してください。

LVM を使用するにはマシン上に `lvm2` がインストールされていることを確認してください。


LVM は複数の物理ストレージデバイスを組み合わせて *ボリュームグループ* にすることができます。
その後このボリュームグループから異なるタイプの *論理ボリューム* を割り当てることができます。

サポートされるボリュームタイプの 1 つに *thin pool* があります。これは許可された最大サイズの合計は利用可能な物理ストレージより大きいような薄くプロビジョンされたボリュームを作成することでリソースをオーバーコミットすることを可能にします。
別のタイプは *ボリュームスナップショット* でこれは論理ボリュームの特定の状態をキャプチャーします。

## Incus の `lvm` ドライバー

Incus の `lvm` ドライバーはイメージに論理ボリュームを、インスタンスとスナップショットにボリュームスナップショットを使用します。

Incus はボリュームグループを完全制御できると想定しています。
このため、 Incus が所有しないファイルシステムエンティティは Incus が消してしまうかもしれないので、LVM ボリュームグループ内に置くべきではありません。
しかし、既存のボリュームグループを再利用する必要がある場合（たとえば、あなたの環境ではボリュームグループが 1 つしかない場合）、[`lvm.vg.force_reuse`](storage-lvm-pool-config) を設定することでこれは可能です。

デフォルトでは LVM ストレージプールは LVM thin pool を使用しその中にすべての Incus ストレージエンティティ（イメージ、インスタンス、カスタムボリューム）の論理ボリュームを作成します。
この挙動はプール作成時に [`lvm.use_thinpool`](storage-lvm-pool-config) を `false` に設定することで変更できます。
この場合、Incus はスナップショットでないすべてのストレージエンティティに "通常の" 論理ボリュームを使用します。
これは深刻なパフォーマンスの低下とディスクの空き容量の低下を `lvm` ドライバーに必然的にもたらすことに注意してください（スピードとストレージ使用量の両面で `dir` ドライバーに近くなります）。
この理由は thin pool でない論理ボリュームがスナップショットのスナップショットをサポートしないため、ほとんどのストレージ操作が `rsync` の使用にフォールバックするためです。
さらに、 thin でないスナップショットは作成時に最大のサイズのストレージを予約しなければならないため、 thin スナップショットよりもはるかに大容量のストレージを使用するからです。
このため、このオプションはどうしても必要なユースケースの場合にのみ選択すべきです。

インスタンスの入れ替わりが激しい環境（たとえば、継続的インテグレーション）では、Incus の操作が遅くなるのを回避するため `/etc/lvm/lvm.conf` 内のバックアップの `retain_min` と `retain_days` 設定を調整すべきです。

(storage-lvmcluster)=
## Incus の `lvmcluster` ドライバー

クラスタ内で使用できる第 2 の `lvmcluster` が利用可能です。

これは `lvmlockd` と `sanlock` デーモンに依存し、共有されたディスクや一組のディスクについての分散ロックを提供します。

これは LVM ストレージプールのバッキングとして `FiberChannel LUN`、 `NVMEoF/NVMEoTCP` ディスク、`iSCSI` ドライブのようなリモートの共有されたブロックデバイスを使用可能にします。

```{note}
Thin provisioning はクラスタ LVM とは互換性がないので、ディスク使用量が増えます。
```

これを Incus で使うには、以下が必要です:

- 全てのクラスタメンバーに共有されたブロックデバイスを用意する
- `lvm`、`lvmlockd`、`sanlock` に必要なパッケージをインストールする
- `/etc/lvm/lvm.conf` 内で `use_lvmlockd = 1` と設定し `lvmlockd` を有効にする
- `/etc/lvm/lvmlocal.conf` 内に (クラスタ内で) ユニークな `host_id` の値を設定する
- `lvmlockd` と `sanlock` デーモンの両方を稼働する

## 設定オプション

`lvm` ドライバーを使うストレージプールとこれらのプール内のストレージボリュームには以下の設定オプションが利用できます。

(storage-lvm-pool-config)=
## ストレージプール設定

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_lvm-common start -->
    :end-before: <!-- config group storage_lvm-common end -->
```

{{volume_configuration}}

(storage-lvm-vol-config)=
## ストレージボリューム設定

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_volume_lvm-common start -->
    :end-before: <!-- config group storage_volume_lvm-common end -->
```

[^*]: {{snapshot_pattern_detail}}

ローカルのストレージプールドライバーでストレージバケットを有効にし、 S3 プロトコル経由でアプリケーションがバケットにアクセスできるようにするには{config:option}`server-core:core.storage_buckets_address`サーバー設定を調整する必要があります。

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_bucket_lvm-common start -->
    :end-before: <!-- config group storage_bucket_lvm-common end -->
```
