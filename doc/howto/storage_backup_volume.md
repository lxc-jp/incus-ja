---
myst:
  substitutions:
    type: "volume"
---

(howto-storage-backup-volume)=
# カスタムストレージボリュームをバックアップするには

カスタムストレージボリュームをバックアップするにはいくつかの方法があります:

- {ref}`storage-backup-snapshots`
- {ref}`storage-backup-export`
- {ref}`storage-copy-volume`

<!-- Include start backup types -->
どの方法を選ぶかはあなたのユースケースとお使いのストレージドライバーによって異なります。

一般に、スナップショットは高速で（ストレージドライバーによりますが）空間効率が良いですが、{{type}}と同じストレージプールに保存されますので信頼性はあまり高くありません。
エクスポートファイルは別のディスク上に保存できますのでより信頼性が高いです。
これらは別のストレージプールに{{type}}をリストアするのに使用できます。
ネットワークで接続された別の Incus サーバーがある場合、定期的に{{type}}をこの別のサーバーにコピーするのも高い信頼性が得られます。そしてこの方法は{{type}}のスナップショットをバックアップするためにも使えます。
<!-- Include end backup types -->

```{note}
カスタムストレージボリュームがインスタンスにアタッチされているかもしれませんが、それらはインスタンスの一部ではありません。
ですので、{ref}`インスタンスをバックアップ <instances-backup>`する際カスタムストレージボリュームは保存されません。
ストレージボリュームのデータは別途バックアップする必要があります。
```

(storage-backup-snapshots)=
## ボリュームのバックアップのスナップショットを使用する

特定の日時のストレージボリュームをスナップショットを作成することで保存できます。スナップショットを使えばストレージボリュームを以前の状態に簡単に復元できます。
スナップショットはボリューム自身と同じストレージプールに保存されます。

<!-- Include start optimized snapshots -->
ほとんどのストレージドライバーはスナップショットの最適化された作成をサポートします（{ref}`storage-drivers-features`参照）。
これらのドライバーではスナップショットの作成は高速で空間効率も良いです。
`dir` ドライバーでは、スナップショットの機能は利用できますが、あまり効率がよくありません。
`lvm` ドライバーでは、スナップショットの作成は高速ですが、スナップショットのリストアが効率的なのは thin-pool モードを使っているときだけです。
<!-- Include end optimized snapshots -->

### カスタムストレージボリュームのスナップショットを作成する

カスタムストレージボリュームのスナップショットを作成するには以下のコマンド使用します:

    incus storage volume snapshot <pool_name> <volume_name> [<snapshot_name>]

<!-- Include start create snapshot options -->
既存のスナップショットを置き換えるには、スナップショット名とともに `--reuse` フラグを追加します。

デフォルトでは、 `snapshots.expiry` 設定オプションが設定されていない限り、スナップショットは永遠に保存されます。
全般的な期限が設定されていてもスナップショットを維持するには、 `--no-expiry` フラグを使用してください。
<!-- Include end create snapshot options -->

(storage-edit-snapshots)=
### スナップショットを表示、編集、削除する

ストレージボリュームのスナップショットを表示するには以下のコマンドを使います:

    incus storage volume info <pool_name> <volume_name>

スナップショットを `<volume_name>/<snapshot_name>` で参照することで、ストレージボリュームの場合と同様にスナップショットを表示または変更できます。

スナップショットの情報を表示するには、以下のコマンドを使います:

    incus storage volume show <pool_name> <volume_name>/<snapshot_name>

スナップショットを編集（たとえば、説明を追加したり有効期限を編集）には以下のコマンドを使います:

    incus storage volume edit <pool_name> <volume_name>/<snapshot_name>

スナップショットを削除するには、以下のコマンドを使います:

    incus storage volume delete <pool_name> <volume_name>/<snapshot_name>

### カスタムストレージボリュームのスナップショット作成をスケジュールする

指定した時刻に自動的にスナップショットを作成するようにカスタムストレージボリュームを設定できます。
そのためには、ストレージボリュームの `snapshots.schedule` 設定オプションを設定してください（{ref}`storage-configure-volume`参照）。

たとえば、日次のスナップショットを設定するには、以下のコマンドを使います:

    incus storage volume set <pool_name> <volume_name> snapshots.schedule @daily

毎日 AM 6 時にスナップショットを作成するよう設定するには、以下のコマンドを使います:

    incus storage volume set <pool_name> <volume_name> snapshots.schedule "0 6 * * *"

定期的にスナップショットをスケジュールする際、自動破棄（`snapshots.expiry`）とスナップショットの命名規則（`snapshots.pattern`）の設定も検討してください。
設定オプションの詳細は{ref}`storage-drivers`のドキュメントを参照してください。

### カスタムストレージボリュームのスナップショットをリストアする

カスタムストレージボリュームを任意のスナップショットの状態に復元できます。

そのためには、まずストレージボリュームを使用しているすべてのインスタンスを停止する必要があります。
その後以下のコマンドを入力します:

    incus storage volume restore <pool_name> <volume_name> <snapshot_name>

スナップショットを同じまたは別のストレージプール（リモートのストレージプールでも可）内に新しいカスタムストレージボリュームをリストアもできます。
そのためには、以下のコマンドを入力します:

    incus storage volume copy <source_pool_name>/<source_volume_name>/<source_snapshot_name> <target_pool_name>/<target_volume_name>

(storage-backup-export)=
## ボリュームのバックアップにエクスポートファイルを使用する

カスタムストレージボリュームの完全な内容をスタンドアロンのファイルにエクスポートし、任意の場所に保存できます。
信頼度を最大化するため、失われたり壊れたりしないように、バックアップファイルは別のファイルシステムに保存してください。

### カスタムストレージボリュームをエクスポートする

以下のコマンドを使ってインスタンスを圧縮ファイル（たとえば、`/path/to/my-instance.tgz`）にエクスポートします:

    incus storage volume export <pool_name> <volume_name> [<file_path>]

ファイルパスを指定しない場合、エクスポートファイルは作業ディレクトリーに `backup.tar.gz` という名前で保存されます。

```{warning}
出力ファイルがすでに存在する場合、コマンドは警告なしで既存のファイルを上書きします。
```

<!-- Include start export info -->
コマンドに以下のフラグを追加できます:

`--compression`
: デフォルトでは、出力ファイルは `gzip` 圧縮を使用します。
  別の圧縮アルゴリズム（たとえば、`bzip2`）を指定したり、`--compression=none` で圧縮しないようにできます。

`--optimized-storage`
: ストレージプールが `btrfs` か `zfs` ドライバーを使用している場合、 `--optimized-storage` フラグを指定すると個別のファイルのアーカイブではなくドライバー固有のバイナリ形式でデータを保存します。
  この場合、エスクポートファイルは同じストレージドライバーを使うプールでのみ使用できます。

  最適化されたモードでボリュームをエクスポートするほうが個別のファイルをエクスポートするより通常は高速です。
  スナップショットはメインボリュームからの差分としてエクスポートされるため、サイズが小さくなりアクセスが容易になります。
<!-- Include end export info -->

`--volume-only`
: デフォルトでは、エクスポートファイルはストレージボリュームのすべてのスナップショットを含みます。
  このフラグを追加すると、スナップショットを除いたストレージボリュームのみをエクスポートします。

### エクスポートファイルからカスタムストレージボリュームをリストアする

エクスポートファイル（たとえば、 `/path/to/my-backup.tgz`）を新しいカスタムストレージボリュームとしてインポートできます。
そのためには、以下のコマンドを使用します:

    incus storage volume import <pool_name> <file_path> [<volume_name>]

ボリューム名を指定しない場合、新しいボリュームの名前はエクスポートされたストレージボリュームの元の名前になります。
その名前のボリュームが指定したストレージブールにすでに（あるいはまだ）存在する場合、コマンドはエラーを返します。
その場合、バックアップをインポートする前に既存のボリュームを削除するか、あるいはインポートの際に別のボリューム名を指定してください。
