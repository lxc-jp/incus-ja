(move-instances)=
# サーバー間で既存の Incus インスタンスを移動するには

ある Incus サーバーから別のサーバーへインスタンスを移動するには [`incus move`](incus_move.md) コマンドを使います:

    incus move [<source_remote>:]<source_instance_name> <target_remote>:[<target_instance_name>]

```{note}
コンテナを移動する際にはまず停止する必要があります。
詳細は {ref}`live-migration-containers` を参照してください。

仮想マシンを移動する際は、{ref}`live-migration-vms` を有効にするか、まず仮想マシンを停止する必要があります。
```

あるいは、インスタンスを複製したい場合は [`incus copy`](incus_copy.md) コマンドを使えます:

    incus copy [<source_remote>:]<source_instance_name> <target_remote>:[<target_instance_name>]

どちらの場合も、移動元のリモートがデフォルトのリモートの場合は省略可能で、移動先でも同じインスタンス名を使用する場合は移動先インスタンス名は省略できます。
インスタンスを特定のクラスタメンバーに移動したい場合は、`--target` フラグを指定してください。
この場合、移動元と移動先のリモートは指定を省略してください。

ネットワークのセットアップに応じて、`--mode` フラグを追加して転送モードを選択できます:

`pull`（デフォルト）
: 移動先のサーバーに、移動元のサーバーへ接続させ該当のインスタンスをプルするように指示します。

`push`
: 移動元のサーバーに、移動先のサーバーへ接続させインスタンスをプッシュするように指示します。

`relay`
: クライアントに移動元と移動先の両方に接続させデータをクライアント経由で転送するよう指示します。

移動先のサーバー上でインスタンスを動かすように設定を調整する必要がある場合、（`--config`, `--device`, `--storage`, `--target-project` を使用して）設定を直接指定するか、（`--no-profiles` か `--profile` を使って）プロファイルを経由して指定できます。すべての利用可能なフラグについては [`incus move --help`](incus_move.md)  を参照してください。

(live-migration)=
## ライブマイグレーション

ライブマイグレーションとはインスタンスの稼働中にマイグレートするという意味です。
仮想マシンではフルにサポートされています。
コンテナでは限定的にサポートされています。

(live-migration-vms)=
### 仮想マシンのライブマイグレーション

仮想マシンは稼働したまま、つまり一切のダウンタイムなしで、別のサーバーに移動できます。

ライブマイグレーションを可能にするには、ステートフルマイグレーションのサポートを有効にする必要があります。
そのためには、以下の設定を確認してください。

* インスタンスの {config:option}`instance-migration:migration.stateful` を `true` に設定する。

(live-migration-containers)=
### コンテナのライブマイグレーション

コンテナについては [{abbr}`CRIU (Checkpoint/Restore in Userspace)`](https://criu.org/) を使用したライブマイグレーションが限定的にサポートされています。
しかし、広範囲に及ぶカーネルへの依存のため、非常にベーシックなコンテナ（ネットワークデバイスなしの非 `systemd` コンテナ）のみが安定してマイグレートできます。
ほとんどの実世界でのシナリオでは、コンテナを停止、移動してその後起動するのが良いです。

コンテナのライブマイグレーションを使用したい場合、マイグレーション元と先の両方のサーバーで CRIU を有効にする必要があります。

コンテナのメモリー転送を最適化するには {config:option}`instance-migration:migration.incremental.memory` プロパティを `true` に設定して CRIU の事前コピー機能を使用してください。
この設定では Incus はコンテナの一連のメモリーダーンプを実行するよう CRIU に指示します。
それぞれのダンプの後、 Incus はメモリーダーンプを指定されたリモートに送信します。
理想的なシナリオでは、各メモリーダーンプを前のメモリーダーンプとの差分にまで減らし、それによりすでに同期されたメモリーの割合を増やします。
同期されたメモリーの割合が {config:option}`instance-migration:migration.incremental.memory.goal` で設定した閾値と等しいか超えた場合、あるいは {config:option}`instance-migration:migration.incremental.memory.iterations` で指定された許容される繰り返し回数の最大値に達した場合、 Incus は CRIU に最終的なメモリーダーンプを実行し、転送するように要求します。
