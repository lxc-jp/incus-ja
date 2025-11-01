(storage-btrfs)=
# Btrfs - `btrfs`

{abbr}`Btrfs (B-tree file system)` は {abbr}`COW (copy-on-write)` 原則に基づいたローカルファイルシステムです。
COW はデータが修正された後に既存のデータを上書きするのではなく別のブロックに保管され、データ破壊のリスクが低くなることを意味します。
他のファイルシステムと異なり、Btrfs はエクステントベースです。これはデータを連続したメモリー領域に保管することを意味します。

基本的なファイルシステムの機能に加えて、Btrfs は RAID、ボリューム管理、プーリング、スナップショット、チェックサム、圧縮、その他の機能を提供します。

Btrfs を使うにはマシンに `btrfs-progs` がインストールされているか確認してください。

## 用語

Btrfs ファイルシステムは*サブボリューム*を持つことができます。これはファイルシステムのメインツリーの名前をつけられたバイナリサブツリーでそれ自身の独立したファイルとディレクトリー階層を持ちます。
*Btrfs スナップショット*は特殊なタイプのサブボリュームで別のサブボリュームの特定の状態をキャプチャーします。
スナップショットは読み書き可または読み取り専用にできます。

## Incus の `btrfs` ドライバー

Incus の `btrfs` ドライバーはインスタンス、イメージ、スナップショットごとにサブボリュームを使用します。
新しいエンティティを作成する際（たとえば、新しいインスタンスを起動する）、 Btrfs スナップショットを作成します。

Btrfs はブロックデバイスの保管をネイティブにはサポートしていません。
このため、仮想マシンに Btrfs を使用する場合、 Incus は仮想マシンを格納するディスク上に巨大なファイルを作成します。
このアプローチはあまり効率的ではなく、スナップショット作成時に問題を引き起こすかもしれません。

Btrfs はネストした Incus 環境内のコンテナ内部でストレージバックエンドとして使用できます。
この場合、親のコンテナ自体は Btrfs を使う必要があります。
しかし、ネストした Incus のセットアップは親から Btrfs のクォータは引き継がないことに注意してください（以下の {ref}`storage-btrfs-quotas` 参照）。

(storage-btrfs-quotas)=
### クォータ

Btrfs は qgroups 経由でストレージクォータをサポートします。
Btrfs qgroups は階層的ですが、新しいサブボリュームは親のサブボリュームの qgroups に自動的に追加されるわけではありません。
これはユーザーが設定されたクォータから逃れることができることは自明であることを意味します。
このため、厳密なクォータが必要な場合は、別のストレージドライバーを検討すべきです（たとえば、`refquotas` ありの ZFS や LVM 上の Btrfs）。

クォータを使用する際は、 Btrfs のエクステントはイミュータブルであることを考慮に入れる必要があります。
ブロックが書かれると、それらは新しいエクステントに現れます。
古いエクステントはその上のすべてのデータが参照されなくなるか上書きされるまで残ります。
これはサブボリューム内で現在存在するファイルで使用されている合計容量がクォータより小さい場合でもクォータに達することがあり得ることを意味します。

```{note}
この問題は Btrfs 上で仮想マシンを使用する際にもっともよく発生します。これは Btrfs サブボリューム上に生のディスクイメージを使用する際のランダムな I/O の性質のためです。

このため、仮想マシンには Btrfs ストレージプールは決して使うべきではありません。

どうしても仮想マシンに Btrfs ストレージプールを使う必要がある場合、インスタンスのルートディスクの [`size.state`](devices-disk) をルートディスクのサイズの2倍に設定してください。
この設定により、ディスクイメージファイルの全てのブロックが qgroup クォータに達すること無しに上書きできるようになります。
[`btrfs.mount_options=compress-force`](storage-btrfs-pool-config) ストレージプールオプションでもこのシナリオを回避できます。圧縮を有効にすることの副作用で最大のエクステントサイズを縮小しブロックの再書き込みが2倍のストレージを消費しないようになるからです。
しかし、これはストレージプールのオプションなので、プール上の全てのボリュームに影響します。
```

## 設定オプション

`btrfs` ドライバーを使うストレージプールとこれらのプール内のストレージボリュームには以下の設定オプションが利用できます。

(storage-btrfs-pool-config)=
## ストレージプール設定

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_btrfs-common start -->
    :end-before: <!-- config group storage_btrfs-common end -->
```

{{volume_configuration}}

### ストレージボリューム設定

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_volume_btrfs-common start -->
    :end-before: <!-- config group storage_volume_btrfs-common end -->
```

[^*]: {{snapshot_pattern_detail}}

### ストレージバケット設定

ローカルのストレージプールドライバーでストレージバケットを有効にし、 S3 プロトコル経由でアプリケーションがバケットにアクセスできるようにするには{config:option}`server-core:core.storage_buckets_address`サーバー設定を調整する必要があります。

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_bucket_btrfs-common start -->
    :end-before: <!-- config group storage_bucket_btrfs-common end -->
```
