(storage-zfs)=
# ZFS - `zfs`

{abbr}`ZFS (Zettabyte file system)` は物理ボリューム管理とファイルシステムを兼ね備えています。
ZFS のインストールは一連のストレージデバイスに広がることができ非常にスケーラブルで、ディスクを追加してストレージプールの空き容量を即座に拡大できます。

ZFS はブロックベースのファイルシステムで、あらゆる操作を検証、確認、訂正するのにチェックサムを使用することでデータ破壊から守ります。
十分な速度で動作するためには、この機構には強力な環境と大量の RAM が必要です。

さらに、 ZFS はスナップショット、レプリケーション、RAID 管理、コピー・オン・ライトのクローン、圧縮、その他の機能を提供します。

ZFS を使用するにはマシンに `zfsutils-linux` をインストールしていることを確認してください。

## 用語

ZFS は物理ストレージデバイスに基づいた論理ユニットを作成します。
これらの論理ユニットは *ZFS pools* または *zpools* と呼ばれます。
さらにそれぞれの zpool は複数の *`データセット`* に分割されます。
これらの`データセット`は以下の異なるタイプがあります。

- *ZFS ファイルシステム* はパーティションまたはマウントされたファイルシステムとして扱えます。
- *ZFS ボリューム* はブロックデバイスを表します。
- *ZFS スナップショット* は ZFS ファイルシステムまたは ZFS ボリュームの特定の状態をキャプチャーします。
  ZFS スナップショットは読み取り専用です。
- *ZFS クローン* は ZFS スナップショットの書き込み可能なコピーです。

## Incus の `zfs` ドライバー

Incus の `zfs` ドライバーは `ZFS ファイルシステム`と ZFS ボリュームをイメージとカスタムストレージボリュームに使用し、ZFS スナップショットとクローンをイメージからのインスタンス作成とインスタンスとカスタムボリュームスナップショットに使用します。
デフォルトでは Incus は ZFS プール作成時に圧縮を有効にします。

Incus は ZFS プールと`データセット`を完全制御できると想定します。
このため、ZFS プールまたは`データセット`内に、Incus が所有しないファイルシステムエンティティは、決して置くべきではありません。Incus が消してしまうかもしれないからです。

ZFS のコピー・オン・ライトが動作する仕組みのため、親の ZFS ファイルシステムはすべての子供がいなくなるまで削除できません。
その結果、Incus は削除されたがまだ参照されているすべてのオブジェクトを自動的にリネームします。
それらのオブジェクトはすべての参照がいなくなりオブジェクトが安全に削除できるようになるまでランダムな `deleted/` パスに保管されます。
この方法はスナップショットの復元に予期しない結果をもたらすかもしれないことに注意してください。
下記の {ref}`storage-zfs-limitations` を参照してください。

ZFS 0.8 以降の上で新しく作成されたすべてのプールで Incus はトリミングサポートを自動的に有効化します。
これはコントローラーによるより良いブロックの再利用を可能にすることで SSD の寿命を伸ばし、ループバックの ZFS プールを使用する際にはルートファイルシステム上の容量を解放できるようにもします。
ZFS の 0.8 より前のバージョンを稼働していてトリミングを有効にしたい場合は、少なくともバージョン 0.8 にアップグレードしてください。
そして以下のコマンドを実行し、今後作成されるプールにトリミングが自動的に有効化され、現在未使用のスペースのすべてがトリムされることを確認してください:

    zpool upgrade ZPOOL-NAME
    zpool set autotrim=on ZPOOL-NAME
    zpool trim ZPOOL-NAME

(storage-zfs-limitations)=
### 制限

`zfs` ドライバーには以下の制限があります。

プールの一部を委譲する
: ZFS はプールの一部をコンテナユーザーに委譲することをサポートしていません。
  ZFS のアップストリームではこの機能を提供すべくアクティブに作業中です。

古いスナップショットからの復元
: ZFS は最新ではないスナップショットからの復元をサポートしていません。
  ですが、古いスナップショットから新しいインスタンスを作成することはできます。
  この方法は特定のスナップショットが必要なものを含んでいるかを確認することを可能にします。
  正しいスナップショットを決定したら {ref}`指定より新しいスナップショットを削除 <storage-edit-snapshots>` して必要なスナップショットが最新になるようにして復元できるようにします。

  別の方法として、復元中により新しいスナップショットを自動的に破棄するように Incus を設定することもできます。
  そのためにはボリュームの [`zfs.remove_snapshots`](storage-zfs-vol-config)（あるいはプール内のすべてのボリュームのストレージプールの対応する `volume.zfs.remove_snapshots` 設定）を設定します。

  しかし、 [`zfs.clone_copy`](storage-zfs-pool-config) が `true` に設定される場合は、インスタンスのコピーは ZFS のスナップショットも使用することに注意してください。
  この場合は、スナップショットのすべての子孫を削除すること無しに、インスタンスを最後のコピーの前に取られたスナップショットに復元できません。
  この選択肢が選べない場合、欲しいスナップショットを新しいインスタンスにコピーしてから古いインスタンスを削除することはできます。
  しかし、インスタンスが持っていたであろう他のすべてのスナップショットは失うことになります。

I/O クォータを観測する
: I/O クォータは ZFS ファイルシステムに大きな影響は与えません。
  これは ZFS は (SPL を使用した) Solaris モジュールの移植でありネイティブな Linux ファイルシステムではないためで、 I/O の制限はネイティブ Linux ファイルシステムに適用されるからです。

ZFS の機能サポート
: idmap の使用や ZFS データセットの委任などの一部の機能には、ZFS 2.2 以上が必要なため、まだ広く利用できません。

### クォータ

ZFS は `quota` と `refquota` という 2 種類の異なるクォータのプロパティを提供します。
`quota` はスナップショットとクローンを含むデータセットの合計サイズを制限します。
`refquota` はスナップショットとクローンは含まずデータセット内のデータのサイズだけを制限します。

デフォルトでは、ストレージボリュームにクォータを設定する際は Incus は `quota` プロパティを使用します。
代わりに `refquota` プロパティを使用したい場合はボリュームの [`zfs.use_refquota`](storage-zfs-vol-config) 設定（あるいはプール内のすべてのボリュームのストレージプールの対応する `volume.zfs.use_refquota` 設定）を設定します。

また [`zfs.use_reserve_space`](storage-zfs-vol-config) (または `volume.zfs.use_reserve_space`) 設定を to use ZFS の `reservation` または `refreservation` を `quota` または `refquota` と使用するために設定することもできます。

## 設定オプション

`zfs` ドライバーを使うストレージプールとこれらのプール内のストレージボリュームには以下の設定オプションが利用できます。

(storage-zfs-pool-config)=
## ストレージプール設定

キー             | 型     | デフォルト値                                               | 説明
:--              | :---   | :------                                                    | :----------
`size`           | string | 自動（空きディスクスペースの 20%, >= 5 GiB and <= 30 GiB） | ループベースのプールを作成する際のストレージプールのサイズ（バイト単位、接尾辞のサポートあり、増やすとストレージプールのサイズを拡大）
`source`         | string | -                                                          | 既存のブロックデバイスかループファイルか ZFS データセット/プールのパス。複数のブロックデバイスは`,`で区切ります。ブロックデバイスをリストする際は`vdev`タイプを接頭辞として指定することもできます。`vdev`タイプを指定するには`vdev`タイプとブロックデバイスの間に`=`記号を書きます（例、`mirror=/dev/sda,/dev/sdb`）。`stripe`、`mirror`、`raidz1`、`raidz2`のみが`vdev`タイプとして使用できます。
`source.wipe`    | bool   | `false`                                                    | ストレージプールを作成する前に`source`で指定されたブロックデバイスの中身を消去する
`zfs.clone_copy` | string | `true`                                                     | Boolean の文字列を指定した場合は ZFS のフル `データセット`コピーの代わりに軽量なクローンを使うかどうかを制御し、 `rebase` という文字列を指定した場合は初期イメージをベースにコピーします。
`zfs.export`     | bool   | `true`                                                     | アンマウントの実行中にzpoolのエクスポートを無効にする
`zfs.pool_name`  | string | プールの名前                                               | zpool 名

{{volume_configuration}}

(storage-zfs-vol-config)=
## ストレージボリューム設定

```{rst-class} break-col-4 min-width-4-8
```

キー                      | 型     | 条件                                             | デフォルト値                                   | 説明
:--                       | :---   | :--------                                        | :------                                        | :----------
`block.filesystem`        | string | `zfs.block_mode`が有効                           | `volume.block.filesystem` と同じ               | {{block_filesystem}}
`block.mount_options`     | string | `zfs.block_mode`が有効                           | `volume.block.mount_options` と同じ            | block-backedなファイルシステムボリュームのマウントオプション
`initial.gid`             | int    | コンテントタイプ`filesystem`のカスタムボリューム | `volume.initial.uid`と同じか`0`                | インスタンス内のボリュームの所有者のGID
`initial.mode`            | int    | コンテントタイプ`filesystem`のカスタムボリューム | `volume.initial.mode`と同じか`711`             | インスタンス内のボリュームのモード
`initial.uid`             | int    | コンテントタイプ`filesystem`のカスタムボリューム | `volume.initial.gid`と同じか`0`                | インスタンス内のボリュームの所有者のUID
`security.shared`         | bool   | カスタムブロックボリューム                       | `volume.security.shared` と同じか `false`      | 複数のインスタンスでのボリュームの共有を有効にする
`security.shifted`        | bool   | カスタムボリューム                               | `volume.security.shifted` と同じか `false`     | {{enable_ID_shifting}}
`security.unmapped`       | bool   | カスタムボリューム                               | `volume.security.unmapped` と同じか `false`    | ボリュームの ID マッピングを無効にする
`size`                    | string |                                                  | `volume.size` と同じ                           | ストレージボリュームのサイズ/クォータ
`snapshots.expiry`        | string | カスタムボリューム                               | `volume.snapshots.expiry` と同じ               | {{snapshot_expiry_format}}
`snapshots.expiry.manual` | string | カスタムボリューム                               | `volume.snapshots.expiry.manual` と同じ        | {{snapshot_expiry_format}}
`snapshots.pattern`       | string | カスタムボリューム                               | `volume.snapshots.pattern` と同じか `snap%d`   | {{snapshot_pattern_format}} [^*]
`snapshots.schedule`      | string | カスタムボリューム                               | `snapshots.schedule` と同じ                    | {{snapshot_schedule_format}}
`zfs.blocksize`           | string |                                                  | `volume.zfs.blocksize` と同じ                  | ZFSブロックのサイズを512バイト～16MiBの範囲で指定します（2の累乗でなければなりません）。ブロックボリュームでは、より大きな値が設定されていても、最大値の128KiBが使用されます。
`zfs.block_mode`          | bool   |                                                  | `volume.zfs.block_mode` と同じ                 | `dataset` よりもフォーマットした `zvol` を使うかどうか
`zfs.remove_snapshots`    | bool   |                                                  | `volume.zfs.remove_snapshots` と同じか `false` | 必要に応じてスナップショットを削除するかどうか
`zfs.use_refquota`        | bool   |                                                  | `volume.zfs.use_refquota` と同じか `false`     | 領域の `quota` の代わりに `refquota` を使うかどうか
`zfs.reserve_space`       | bool   |                                                  | `volume.zfs.reserve_space` と同じか `false`    | `qouta`/`refquota` に加えて `reservation`/`refreservation` も使用するかどうか

[^*]: {{snapshot_pattern_detail}}

### ストレージバケット設定

ローカルのストレージプールドライバーでストレージバケットを有効にし、 S3 プロトコル経由でアプリケーションがバケットにアクセスできるようにするには{config:option}`server-core:core.storage_buckets_address`サーバー設定を調整する必要があります。

キー   | 型     | 条件             | デフォルト値         | 説明
:--    | :---   | :--------        | :------              | :----------
`size` | string | 適切なドライバー | `volume.size` と同じ | ストレージバケットのサイズ/クォータ
