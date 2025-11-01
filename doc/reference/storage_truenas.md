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

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_truenas-common start -->
    :end-before: <!-- config group storage_truenas-common end -->
```

{{volume_configuration}}

(storage-truenas-vol-config)=
### ストレージボリューム設定

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_volume_truenas-common start -->
    :end-before: <!-- config group storage_volume_truenas-common end -->
```
