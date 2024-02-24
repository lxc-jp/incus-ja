(storage-dir)=
# ディレクトリー - `dir`

ディレクトリーストレージドライバーは基本的なバックエンドで通常のファイルとディレクトリー構造にデータを保管します。
このドライバーは素早くセットアップできディスク上のファイルを直接見ることができるので、テストには便利かもしれません。
しかし、 Incus の操作はこのドライバー用には {ref}`最適化されていません <storage-drivers-features>`。

## Incus の `dir` ドライバー

Incus の `dir` ドライバーは完全に機能し、他のドライバーと同じ機能セットを提供します。
しかし、他のドライバーよりは圧倒的に遅いです。これはインスタンス、スナップ、ショットを一瞬でコピーする代わりにイメージの解凍を行う必要があるためです。

作成時に（`source` 設定オプションを使って）別途指定されてない限り、データは `/var/lib/incus/storage-pools/` ディレクトリーに保管されます。

(storage-dir-quotas)=
### クォータ

<!-- Include start dir quotas -->
`dir` ドライバーは ext4 か XFS 上で動作しファイルシステムレベルでプロジェクトのクォータが有効な場合にストレージのクォータをサポートします。
<!-- Include end dir quotas -->

## 設定オプション

`dir` ドライバーを使うストレージプールとこれらのプール内のストレージボリュームには以下の設定オプションが利用できます。

## ストレージプール設定

キー                | 型     | デフォルト値   | 説明
:--                 | :---   | :------        | :----------
`rsync.bwlimit`     | string | `0` (no limit) | ストレージエンティティの転送に rsync を使う必要があるときにソケット I/O に指定する上限を設定
`rsync.compression` | bool   | `true`         | ストレージブールのマイグレーションの際に圧縮を使うかどうか
`source`            | string | -              | ブロックデバイスかループファイルかファイルシステムエントリのパス

{{volume_configuration}}

## ストレージボリューム設定

キー                 | 型     | 条件                       | デフォルト値                                 | 説明
:--                  | :---   | :--------                  | :------                                      | :----------
`security.shared`    | bool   | カスタムブロックボリューム | `volume.security.shared` と同じか `false`    | 複数のインスタンスでのボリュームの共有を有効にする
`security.shifted`   | bool   | カスタムボリューム         | `volume.security.shifted` と同じか `false`   | {{enable_ID_shifting}}
`security.unmapped`  | bool   | カスタムボリューム         | `volume.security.unmapped` と同じか `false`  | ボリュームの ID マッピングを無効にする
`size`               | string | 適切なドライバー           | `volume.size` と同じ                         | ストレージボリュームのサイズ/クォータ
`snapshots.expiry`   | string | カスタムボリューム         | `volume.snapshots.expiry` と同じ             | {{snapshot_expiry_format}}
`snapshots.pattern`  | string | カスタムボリューム         | `volume.snapshots.pattern` と同じか `snap%d` | {{snapshot_pattern_format}} [^*]
`snapshots.schedule` | string | カスタムボリューム         | `volume.snapshots.schedule` と同じ           | {{snapshot_schedule_format}}

[^*]: {{snapshot_pattern_detail}}

### ストレージバケット設定

ローカルのストレージプールドライバーでストレージバケットを有効にし、 S3 プロトコル経由でアプリケーションがバケットにアクセスできるようにするには{config:option}`server-core:core.storage_buckets_address`サーバー設定を調整する必要があります。

ストレージバケットは `dir` プール用の設定はありません。
他のストレージプールドライバーとは異なり、 `dir` ドライバーは `size` 設定によるバケットクォータのサポートはありません。
