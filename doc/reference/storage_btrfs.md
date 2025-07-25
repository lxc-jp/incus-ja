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

キー                  | 型     | デフォルト値                                               | 説明
:--                   | :---   | :--------                                                  | :----------
`btrfs.mount_options` | string | `user_subvol_rm_allowed`                                   | ブロックデバイスのマウントオプション
`size`                | string | 自動（空きディスクスペースの 20%, >= 5 GiB and <= 30 GiB） | ループベースのプールを作成する際のストレージプールのサイズ（バイト単位、接尾辞のサポートあり、増やすとストレージプールのサイズを拡大）
`source`              | string | -                                                          | 既存のブロックデバイス、ループファイル、あるいはBtrfsサブボリュームのパス
`source.wipe`         | bool   | `false`                                                    | ストレージプールを作成する前に`source`で指定されたブロックデバイスの中身を消去する

{{volume_configuration}}

### ストレージボリューム設定

キー                      | 型     | 条件                                             | デフォルト値                                 | 説明
:--                       | :---   | :--------                                        | :------                                      | :----------
`initial.gid`             | int    | コンテントタイプ`filesystem`のカスタムボリューム | `volume.initial.uid`と同じか`0`              | インスタンス内のボリュームの所有者のGID
`initial.mode`            | int    | コンテントタイプ`filesystem`のカスタムボリューム | `volume.initial.mode`と同じか`711`           | インスタンス内のボリュームのモード
`initial.uid`             | int    | コンテントタイプ`filesystem`のカスタムボリューム | `volume.initial.gid`と同じか`0`              | インスタンス内のボリュームの所有者のUID
`security.shared`         | bool   | カスタムブロックボリューム                       | `volume.security.shared` と同じか `false`    | 複数のインスタンスでのボリュームの共有を有効にする
`security.shifted`        | bool   | カスタムボリューム                               | `volume.security.shifted` と同じか `false`   | {{enable_ID_shifting}}
`security.unmapped`       | bool   | カスタムボリューム                               | `volume.security.unmapped` と同じか `false`  | ボリュームへの id マッピングを無効にする
`size`                    | string | 適切なドライバー                                 | `volume.size` と同じ                         | ストレージボリュームのサイズ/クォータ
`snapshots.expiry`        | string | カスタムボリューム                               | `volume.snapshots.expiry` と同じ             | {{snapshot_expiry_format}}
`snapshots.expiry.manual` | string | カスタムボリューム                               | `volume.snapshots.expiry.manual` と同じ      | {{snapshot_expiry_format}}
`snapshots.pattern`       | string | カスタムボリューム                               | `volume.snapshots.pattern` と同じか `snap%d` | {{snapshot_pattern_format}} [^*]
`snapshots.schedule`      | string | カスタムボリューム                               | `volume.snapshots.schedule` と同じ           | {{snapshot_schedule_format}}

[^*]: {{snapshot_pattern_detail}}

### ストレージバケット設定

ローカルのストレージプールドライバーでストレージバケットを有効にし、 S3 プロトコル経由でアプリケーションがバケットにアクセスできるようにするには{config:option}`server-core:core.storage_buckets_address`サーバー設定を調整する必要があります。

キー   | 型     | 条件             | デフォルト値         | 説明
:--    | :---   | :--------        | :------              | :----------
`size` | string | 適切なドライバー | `volume.size` と同じ | ストレージバケットのサイズ/クォータ
