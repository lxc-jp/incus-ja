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

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_zfs-common start -->
    :end-before: <!-- config group storage_zfs-common end -->
```

{{volume_configuration}}

(storage-zfs-vol-config)=
## ストレージボリューム設定

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_volume_zfs-common start -->
    :end-before: <!-- config group storage_volume_zfs-common end -->
```

[^*]: {{snapshot_pattern_detail}}

### ストレージバケット設定

ローカルのストレージプールドライバーでストレージバケットを有効にし、 S3 プロトコル経由でアプリケーションがバケットにアクセスできるようにするには{config:option}`server-core:core.storage_buckets_address`サーバー設定を調整する必要があります。

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_bucket_zfs-common start -->
    :end-before: <!-- config group storage_bucket_zfs-common end -->
```
