---
myst:
  substitutions:
    type: "instance"
---

(instances-backup)=
# インスタンスをバックアップするには

インスタンスをバックアップするにはいくつかの方法があります:

- {ref}`instances-snapshots`
- {ref}`instances-backup-export`
- {ref}`instances-backup-copy`

% Include content from [storage_backup_volume.md](storage_backup_volume.md)
```{include} storage_backup_volume.md
    :start-after: <!-- Include start backup types -->
    :end-before: <!-- Include end backup types -->
```

```{note}
カスタムストレージボリュームがインスタンスにアタッチされているかもしれませんが、それらはインスタンスの一部ではありません。
ですので、インスタンスをバックアップする際カスタムストレージボリュームは保存されません。
ストレージボリュームのデータは別途バックアップする必要があります。
手順は {ref}`howto-storage-backup-volume` を参照してください。
```

(instances-snapshots)=
## インスタンスのバックアップにスナップショットを使用する

特定の日時のインスタンスをスナップショットを作成することで保存できます。スナップショットを使えばインスタンスを以前の状態に簡単に復元できます。

インスタンススナップショットはインスタンスのボリューム自身と同じストレージプールに保存されます。

% Include content from [storage_backup_volume.md](storage_backup_volume.md)
```{include} storage_backup_volume.md
    :start-after: <!-- Include start optimized snapshots -->
    :end-before: <!-- Include end optimized snapshots -->
```

### スナップショットを作成する

インスタンスのスナップショットを作成するには以下のコマンドを使います:

    incus snapshot create <instance_name> [<snapshot name>]

% Include content from [storage_backup_volume.md](storage_backup_volume.md)
```{include} storage_backup_volume.md
    :start-after: <!-- Include start create snapshot options -->
    :end-before: <!-- Include end create snapshot options -->
```

仮想マシンでは、 `--stateful` フラグを指定するとインスタンスボリュームに含まれるデータだけでなく、インスタンスの稼働状態も含めることができます。
CRIU の制限のためコンテナではこの機能は完全にはサポートされていないことに注意してください。

### スナップショットを表示、編集、削除する

インスタンスのスナップショットを表示するには以下のコマンドを使います:

    incus info <instance_name>

スナップショットを `<instance_name>/<snapshot_name>` で参照することで、インスタンスの場合と同様にスナップショットを表示または変更できます。

スナップショットの設定を表示するには、以下のコマンドを使います:

    incus config show <instance_name>/<snapshot_name>

スナップショットの有効期限を変更するには、以下のコマンドを使います:

    incus config edit <instance_name>/<snapshot_name>

```{note}
一般に、スナップショットはインスタンスの状態を保存しているため、編集できません。
唯一の例外が有効期限です。
他の設定の変更は黙って無視されます。
```

スナップショットを削除するには、以下のコマンドを使います:

    incus snapshot delete <instance_name> <snapshot_name>

### インスタンスのスナップショット作成をスケジュールする

指定した時刻（最大で 1 分ごと）に自動的にスナップショットを作成するようにインスタンスを設定できます。
そのためには、 {config:option}`instance-snapshots:snapshots.schedule` インスタンスオプションを設定してください。

たとえば、日次のスナップショットを設定するには、以下のコマンドを使います:

    incus config set <instance_name> snapshots.schedule @daily

毎日 AM 6 時にスナップショットを作成するよう設定するには、以下のコマンドを使います:

    incus config set <instance_name> snapshots.schedule "0 6 * * *"

定期的にスナップショットをスケジュールする際、自動破棄（{config:option}`instance-snapshots:snapshots.expiry`）とスナップショットの命名規則（{config:option}`instance-snapshots:snapshots.pattern`）の設定も検討してください。
また、稼働していないインスタンスのスナップショットを作成するかどうかの設定（{config:option}`instance-snapshots:snapshots.schedule.stopped`）もすると良いでしょう。

### インスタンスのスナップショットをリストアする

インスタンスを任意のスナップショットの状態に復元できます。

そのためには、以下のコマンドを使います:

    incus snapshot restore <instance_name> <snapshot_name>

スナップショットがステートフル（インスタンスの稼働状態の情報を含むことを意味します）の場合、状態をリストアするために `--stateful` を追加できます。

(instances-backup-export)=
## インスタンスのバックアップにエクスポートファイルを使用する

インスタンスの完全な内容をスタンドアロンのファイルにエクスポートし、任意の場所に保存できます。
信頼度を最大化するため、失われたり壊れたりしないように、バックアップファイルは別のファイルシステムに保存してください。

### インスタンスをエクスポートする

以下のコマンドを使ってインスタンスを圧縮ファイル（たとえば、`/path/to/my-instance.tgz`）にエクスポートします:

    incus export <instance_name> [<file_path>]

ファイルパスを指定しない場合、エクスポートファイルは作業ディレクトリーに `<instance_name>.<extension>` （たとえば、`my-container.tar.gz`）という名前で保存されます。

```{warning}
出力ファイル（`<instance_name>.<extension>` または指定した名前）がすでに存在する場合、コマンドは警告なしで既存のファイルを上書きします。
```

% Include content from [storage_backup_volume.md](storage_backup_volume.md)
```{include} storage_backup_volume.md
    :start-after: <!-- Include start export info -->
    :end-before: <!-- Include end export info -->
```

`--instance-only`
: デフォルトでは、エクスポートファイルはインスタンスのすべてのスナップショットを含みます。
  このフラグを追加すると、スナップショットを除いたインスタンスのみをエクスポートします。

### エクスポートファイルからインスタンスをリストアする

エクスポートファイル（たとえば、 `/path/to/my-backup.tgz`）を新しいインスタンスとしてインポートできます。
そのためには、以下のコマンドを使用します:

    incus import <file_path> [<instance_name>]

インスタンス名を指定しない場合、新しいインスタンスの名前はエクスポートされたインスタンスの元の名前になります。
その名前のインスタンスが指定したストレージブールにすでに（あるいはまだ）存在する場合、コマンドはエラーを返します。
その場合、バックアップをインポートする前に既存のインスタンスを削除するか、あるいはインポートの際に別のインスタンス名を指定してください。

(instances-backup-copy)=
## インスタンスをバックアップサーバーにコピーする

インスタンスをバックアップするためにセカンダリバックアップサーバーにコピーできます。

手順は {ref}`move-instances` を参照してください。
