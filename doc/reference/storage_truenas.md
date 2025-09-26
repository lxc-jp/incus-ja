(storage-truenas)=
# TrueNAS - `truenas`
## Incus内の`truenas`ストレージドライバー

`truenas`ストレージドライバーはIncusのノードでリモートのTrueNASストレージドライバーを使って1つ以上のIncusストレージプールを持てるようにします。
ノードがクラスターの一部である場合、すべてのクラスターメンバーはストレージプールに同時にアクセスすることができ、ノード間で仮想マシン（VM）をライブマイグレーションする用途に最適になります。

ドライバーはブロックベースの方式で動作し、これはすべてのIncusのボリュームはリモートのTrueNASサーバー上のZFSのボリュームブロックとして作成されることを意味します。これらのZFSボリュームブロックデバイスはiSCSIを介してローカルのIncusノード上でアクセスされます。

既存のZFSドライバーをモデルにして、`truenas`ドライバーは標準のZFSの機能のほとんどをサポートしますが、リモートのTrueNASサーバー上で動作します。例えば、ローカルのVMのスナップショットを作成したり複製する場合は、ローカルファイルシステムを同期したあとリモートサーバー上でスナップショットと複製の操作が実行されます。複製は必要に応じてiSCSI経由で実行されます。

各ストレージプールはリモートのTrueNASホスト上のZFSデータセットに対応します。
データセットは存在しなければ自動的に作成されます。
ドライバーはリモートホスト上で利用可能なZFSの機能を使って効率的なイメージの処理、コピー操作、スナップショットの管理をネストしたZFS（ZFS-on-ZFS）を必要とせずに実行できます。

リモートのデータセットを参照するには、`source`プロパティを下記の形式で指定できます:
`[<remote host>:]<remote pool>[[/<remote dataset>]...][/]`

パスの最後が`/`で終わる場合、データセット名はIncusのストレージプール名から派生されます（例、`tank/pool1`）。

## 要件

このドライバーはTrueNAS APIを使ってリモートサーバー上でアクションを実行するために[`truenas_incus_ctl`](https://github.com/truenas/truenas_incus_ctl)ツールに依存しています。
このツールはさらに`open-iscsi`を経由してリモートのZFSボリュームのアクティベーションとデアクティベーションを管理もしています。
`truenas_incus_ctl`がインストールされていないかシステムのPATHに存在しない場合、ドライバーは無効化されます。

必要なツールをインストールするには、[`truenas\_incus\_ctl` GitHub page](https://github.com/truenas/truenas_incus_ctl)から最新版（v0.7.2+が必要）をダウンロードしてください。
さらに`open-iscsi`を以下のコマンドでインストールしてください:

    sudo apt install open-iscsi

## TrueNASホストへのログイン

APIキーを手動で作成し`truenas.api_key`プロパティでそれを指定する代わりに、`truenas_incus_ctl`ツールでリモートサーバーにログインすることもできます。

    sudo truenas_incus_ctl config login

するとTrueNASサーバーへ接続する詳細情報（認証の詳細を含む）がプロンプトで表示されます。そして設定をローカルファイルに保存します。ログインした後、iSCSIのセットアップを以下のコマンドで確認できます:

    sudo truenas_incus_ctl share iscsi setup --test

ツールが設定されたら、これを使ってリモートのデータセットを使ってストレージプールを作成できます:

    incus storage create <poolname> truenas source=[host:]<pool>[/<dataset>]/[remote-poolname]

上記のコマンドでは:

* `source`はリモートのTrueNASホスト上でストレージプールを作成する場所を指定します。
* `host`はオプションで、`truenas.host`プロパティで指定することもでき、`truenas.config`で設定することもできます。
* `remote-poolname`が指定されない場合はデフォルトとしてローカルプールの名前が使われます。

## 設定オプション

以下の設定オプションは`truenas`ドライバーを使うストレージプールとこれらのプール内のストレージボリュームで利用できます。

(storage-truenas-pool-config)=
### ストレージプール設定

| キー                     | 型      | デフォルト値 | 説明                                                                                                                                                     |
| :---                     | :---    | :---         | :---                                                                                                                                                     |
| `source`                 | string  | -            | リモートのTrueNASホスト上で使用するZFSデータセット。形式: `[<host>:]<pool>[/<dataset>][/]`。ここで`host`を省略すると、`truenas.host`で設定する必要あり。 |
| `truenas.allow_insecure` | boolean | false        | `true`に設定すると、TrueNAS APIへの安全でない（非TLS）接続を許可します。                                                                                 |
| `truenas.api_key`        | string  | -            | TrueNASホストへの認証に使うAPIキー                                                                                                                       |
| `truenas.dataset`        | string  | -            | リモートデータセット名。通常は`source`から推測されるが、上書きも可能。                                                                                   |
| `truenas.host`           | string  | -            | TrueNASシステムのホスト名またはIPアドレス。`source`内に含まれるか、設定が使われる場合は設定不要。                                                        |
| `truenas.initiator`      | string  | -            | ブロックボリュームのアタッチの際に使われるiSCSIイニシエーター名                                                                                          |
| `truenas.portal`         | string  | -            | ブロックボリュームのアタッチの際に使われるiSCSIポータルアドレス                                                                                          |

{{volume_configuration}}

(storage-truenas-vol-config)=
### ストレージボリューム設定

| キー                       | 型     | 条件                                           | デフォルト値                                     | 説明                                                                                                                                        |
| :---                       | :---   | :---                                           | :---                                             | :---                                                                                                                                        |
| `block.filesystem`         | string |                                                | `volume.block.filesystem`と同じ                  | {{block_filesystem}}                                                                                                                        |
| `block.mount_options`      | string |                                                | `volume.block.mount_options`と同じ               | ブロックベースのファイルシステムボリュームのマウントオプション                                                                              |
| `initial.gid`              | int    | content typeが`filesystem`のカスタムボリューム | `volume.initial.uid`と同じか`0`                  | インスタンス内のボリューム所有者のGID                                                                                                       |
| `initial.mode`             | int    | content typeが`filesystem`のカスタムボリューム | `volume.initial.mode`と同じか`711`               | インスタンス内のボリュームのパーミション                                                                                                    |
| `initial.uid`              | int    | content typeが`filesystem`のカスタムボリューム | `volume.initial.gid`と同じか`0`                  | インスタンス内のボリューム所有者のUID                                                                                                       |
| `security.shared`          | bool   | カスタムブロックボリューム                     | `volume.security.shared`と同じか`false`          | 複数インスタンス間でのボリュームの共有を有効にする                                                                                          |
| `security.shifted`         | bool   | カスタムボリューム                             | `volume.security.shifted`と同じか`false`         | {{enable_ID_shifting}}                                                                                                                      |
| `security.unmapped`        | bool   | カスタムボリューム                             | `volume.security.unmapped`と同じか`false`        | ボリュームのIDマッピングを無効化する                                                                                                        |
| `size`                     | string |                                                | `volume.size`と同じ                              | ストレージボリュームのサイズ／クォータ                                                                                                      |
| `snapshots.expiry`         | string | カスタムボリューム                             | `volume.snapshots.expiry`と同じ                  | {{snapshot_expiry_format}}                                                                                                                  |
| `snapshots.expiry.manual`  | string | カスタムボリューム                             | `volume.snapshots.expiry.manual`と同じ           | {{snapshot_expiry_format}}                                                                                                                  |
| `snapshots.pattern`        | string | カスタムボリューム                             | `volume.snapshots.pattern`と同じか`snap%d`       | {{snapshot_pattern_format}}                                                                                                                 |
| `snapshots.schedule`       | string | カスタムボリューム                             | `snapshots.schedule`と同じ                       | {{snapshot_schedule_format}}                                                                                                                |
| `truenas.blocksize`        | string |                                                | `volume.truenas.blocksize`と同じ                 | ZFSブロックのサイズを512バイトから16MiBの間で(2のべき乗で）指定。ブロックボリュームではより大きな値を設定しても最大値の128KiBが使われます。 |
| `truenas.remove_snapshots` | bool   |                                                | `volume.truenas.remove_snapshots`と同じか`false` | 必要に応じてスナップショットを削除する                                                                                                      |
| `truenas.use_refquota`     | bool   |                                                | `volume.truenas.use_refquota`と同じか`false`     | スペースの`quota`の代わりに`refquota`を使う                                                                                                 |
