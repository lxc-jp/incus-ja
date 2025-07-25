# API 拡張

以下の変更は 1.0 API が確定した後、Incus API に追加されました。

それらの変更はすべて後方互換であり、 `GET /1.0` の `api_extensions` を
見ることでクライアントツールにより検出可能です。

## `storage_zfs_remove_snapshots`

`storage.zfs_remove_snapshots` というデーモン設定キーが導入されました。

値の型は Boolean でデフォルトは `false` です。 `true` にセットすると、スナップショットを
復元しようとするときに必要なスナップショットをすべて削除するように Incus に
指示します。

ZFS でスナップショットの復元が出来るのは最新のスナップショットに限られるので、
この対応が必要になります。

## `container_host_shutdown_timeout`

`boot.host_shutdown_timeout` というコンテナ設定キーが導入されました。

値の型は integer でコンテナを停止しようとした後 kill するまでどれだけ
待つかを Incus に指示します。

この値は Incus デーモンのクリーンなシャットダウンのときにのみ使用されます。
デフォルトは 30s です。

## `container_stop_priority`

`boot.stop.priority` というコンテナ設定キーが導入されました。

値の型は integer でシャットダウン時のコンテナの優先度を指示します。

コンテナは優先度レベルの高いものからシャットダウンを開始します。

同じ優先度のコンテナは並列にシャットダウンします。デフォルトは 0 です。

## `container_syscall_filtering`

コンテナ設定キーに関するいくつかの新しい syscall が導入されました。

* `security.syscalls.blacklist_default` <!-- wokeignore:rule=blacklist -->
* `security.syscalls.blacklist_compat` <!-- wokeignore:rule=blacklist -->
* `security.syscalls.blacklist` <!-- wokeignore:rule=blacklist -->
* `security.syscalls.whitelist` <!-- wokeignore:rule=whitelist -->

使い方は [インスタンスの設定](instance-config) を参照してください。

## `auth_pki`

これは PKI 認証モードのサポートを指示します。

このモードではクライアントとサーバーは同じ PKI によって発行された証明書を使わなければなりません。

詳細は [セキュリティ](security.md) を参照してください。

## `container_last_used_at`

`GET /1.0/containers/<name>` エンドポイントに `last_used_at` フィールドが追加されました。

これはコンテナが開始した最新の時刻のタイムスタンプです。

コンテナが作成されたが開始はされていない場合は `last_used_at` フィールドは
`1970-01-01T00:00:00Z` になります。

## `etag`

関連性のあるすべてのエンドポイントに ETag ヘッダのサポートが追加されました。

この変更により GET のレスポンスに次の HTTP ヘッダが追加されます。

* ETag（ユーザーが変更可能なコンテンツの SHA-256）

また PUT リクエストに次の HTTP ヘッダのサポートが追加されます。

* If-Match（前回の GET で得られた ETag の値を指定）

これにより GET で Incus のオブジェクトを取得して PUT で変更する際に、
レースコンディションになったり、途中で別のクライアントがオブジェクトを
変更していた（訳注: のを上書きしてしまう）というリスク無しに PUT で
変更できるようになります。

## `patch`

HTTP の PATCH メソッドのサポートを追加します。

PUT の代わりに PATCH を使うとオブジェクトの部分的な変更が出来ます。

## `usb_devices`

USB ホットプラグのサポートを追加します。

## `https_allowed_credentials`

Incus API をすべてのウェブブラウザで（SPA 経由で）使用するには、 XHR の度に
認証情報を送る必要があります。それぞれの XHR リクエストで
[`withCredentials=true`](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest/withCredentials)
とセットします。

Firefox や Safari などいくつかのブラウザは
`Access-Control-Allow-Credentials: true` ヘッダがないレスポンスを受け入れる
ことができません。サーバーがこのヘッダ付きのレスポンスを返すことを保証するには
`core.https_allowed_credentials=true` と設定してください。

## `image_compression_algorithm`

この変更はイメージを作成する時（`POST /1.0/images`）に `compression_algorithm`
というプロパティのサポートを追加します。

このプロパティを設定するとサーバーのデフォルト値 (`images.compression_algorithm`) をオーバーライドします。

## `directory_manipulation`

Incus API 経由でディレクトリーを作成したり一覧したりでき、ファイルタイプを X-Incus-type ヘッダに付与するようになります。
現状はファイルタイプは `file` か `directory` のいずれかです。

## `container_cpu_time`

この拡張により実行中のコンテナの CPU 時間を取得できます。

## `storage_zfs_use_refquota`

この拡張により新しいサーバープロパティ `storage.zfs_use_refquota` が追加されます。
これはコンテナにサイズ制限を設定する際に `quota` の代わりに `refquota` を設定する
ように Incus に指示します。また Incus はディスク使用量を調べる際に `used` の代わりに
`usedbydataset` を使うようになります。

これはスナップショットによるディスク消費をコンテナのディスク利用の一部と
みなすかどうかを実質的に切り替えることになります。

## `storage_lvm_mount_options`

この拡張は `storage.lvm_mount_options` という新しいデーモン設定オプションを
追加します。デフォルト値は `discard` で、このオプションにより LVM LV で使用する
ファイルシステムの追加マウントオプションをユーザーが指定できるようになります。

## `network`

Incus のネットワーク管理 API 。

次のものを含みます。

* `/1.0/networks` エントリーに `managed` プロパティを追加
* ネットワーク設定オプションのすべて（詳細は [ネットワーク設定](networks.md) を参照）
* `POST /1.0/networks`（詳細は [RESTful API](rest-api.md) を参照）
* `PUT /1.0/networks/<entry>`（詳細は [RESTful API](rest-api.md) を参照）
* `PATCH /1.0/networks/<entry>`（詳細は [RESTful API](rest-api.md) を参照）
* `DELETE /1.0/networks/<entry>`（詳細は [RESTful API](rest-api.md) を参照）
* `nic` タイプのデバイスの `ipv4.address` プロパティ（`nictype` が `bridged` の場合）
* `nic` タイプのデバイスの `ipv6.address` プロパティ（`nictype` が `bridged` の場合）
* `nic` タイプのデバイスの `security.mac_filtering` プロパティ（`nictype` が `bridged` の場合）

## `profile_usedby`

プロファイルを使用しているコンテナをプロファイルエントリーの一覧の `used_by` フィールド
として新たに追加します。

## `container_push`

コンテナが push モードで作成される時、クライアントは作成元と作成先のサーバー間の
プロキシとして機能します。作成先のサーバーが NAT やファイアウォールの後ろにいて
作成元のサーバーと直接通信できず pull モードで作成できないときにこれは便利です。

## `container_exec_recording`

新しい Boolean 型の `record-output` を導入します。これは `/1.0/containers/<name>/exec`
のパラメーターでこれを `true` に設定し `wait-for-websocket` を `false` に設定すると
標準出力と標準エラー出力をディスクに保存し logs インターフェース経由で利用可能にします。

記録された出力の URL はコマンドが実行完了したら操作のメタデータに含まれます。

出力は他のログファイルと同様に、通常は 48 時間後に期限切れになります。

## `certificate_update`

REST API に次のものを追加します:

* 証明書の GET に ETag ヘッダ
* 証明書エントリーの PUT
* 証明書エントリーの PATCH

## `container_exec_signal_handling`

クライアントに送られたシグナルをコンテナ内で実行中のプロセスにフォワーディング
するサポートを `/1.0/containers/<name>/exec` に追加します。現状では SIGTERM と
SIGHUP がフォワードされます。フォワード出来るシグナルは今後さらに追加される
かもしれません。

## `gpu_devices`

コンテナに GPU を追加できるようにします。

## `container_image_properties`

設定キー空間に新しく `image` を導入します。これは読み取り専用で、親のイメージのプロパティを
含みます。

## `migration_progress`

転送の進捗が操作の一部として送信側と受信側の両方に公開されます。これは操作のメタデータの
`fs_progress` 属性として現れます。

## `id_map`

`security.idmap.isolated`、`security.idmap.isolated`、`security.idmap.size`、`raw.id_map` のフィールドを設定できるようにします。

## `network_firewall_filtering`

`ipv4.firewall` と `ipv6.firewall` という 2 つのキーを追加します。
`false` に設置すると `iptables` の FORWARDING ルールの生成をしないように
なります。 NAT ルールは対応する `ipv4.nat` や `ipv6.nat` キーが `true` に
設定されている限り引き続き追加されます。

ブリッジに対して `dnsmasq` が有効な場合、 `dnsmasq` が機能する（DHCP/DNS）
ために必要なルールは常に適用されます。

## `network_routes`

`ipv4.routes` と `ipv6.routes` を導入します。これらは Incus ブリッジに
追加のサブネットをルーティングできるようにします。

## `storage`

Incus のストレージ管理 API 。

これは次のものを含みます。

* `GET /1.0/storage-pools`
* `POST /1.0/storage-pools`（詳細は [RESTful API](rest-api.md) を参照）

* `GET /1.0/storage-pools/<name>`（詳細は [RESTful API](rest-api.md) を参照）
* `POST /1.0/storage-pools/<name>`（詳細は [RESTful API](rest-api.md) を参照）
* `PUT /1.0/storage-pools/<name>`（詳細は [RESTful API](rest-api.md) を参照）
* `PATCH /1.0/storage-pools/<name>`（詳細は [RESTful API](rest-api.md) を参照）
* `DELETE /1.0/storage-pools/<name>`（詳細は [RESTful API](rest-api.md) を参照）

* `GET /1.0/storage-pools/<name>/volumes`（詳細は [RESTful API](rest-api.md) を参照）

* `GET /1.0/storage-pools/<name>/volumes/<volume_type>`（詳細は [RESTful API](rest-api.md) を参照）
* `POST /1.0/storage-pools/<name>/volumes/<volume_type>`（詳細は [RESTful API](rest-api.md) を参照）

* `GET /1.0/storage-pools/<pool>/volumes/<volume_type>/<name>`（詳細は [RESTful API](rest-api.md) を参照）
* `POST /1.0/storage-pools/<pool>/volumes/<volume_type>/<name>`（詳細は [RESTful API](rest-api.md) を参照）
* `PUT /1.0/storage-pools/<pool>/volumes/<volume_type>/<name>`（詳細は [RESTful API](rest-api.md) を参照）
* `PATCH /1.0/storage-pools/<pool>/volumes/<volume_type>/<name>`（詳細は [RESTful API](rest-api.md) を参照）
* `DELETE /1.0/storage-pools/<pool>/volumes/<volume_type>/<name>`（詳細は [RESTful API](rest-api.md) を参照）

* すべてのストレージ設定オプション（詳細は [ストレージの設定](storage.md) を参照）

## `file_delete`

`/1.0/containers/<name>/files` の DELETE メソッドを実装します。

## `file_append`

`X-Incus-write` ヘッダを実装しました。値は `overwrite` か `append` のいずれかです。

## `network_dhcp_expiry`

`ipv4.dhcp.expiry` と `ipv6.dhcp.expiry` を導入します。 DHCP のリース期限を設定
できるようにします。

## `storage_lvm_vg_rename`

`storage.lvm.vg_name` を設定することでボリュームグループをリネームできるようにします。

## `storage_lvm_thinpool_rename`

`storage.thinpool_name` を設定することで thin pool をリネームできるようにします。

## `network_vlan`

`macvlan` ネットワークデバイスに `vlan` プロパティを新たに追加します。

これを設定すると、指定した VLAN にアタッチするように Incus に指示します。
Incus はホスト上でその VLAN を持つ既存のインターフェースを探します。
もし見つからない場合は Incus がインターフェースを作成して macvlan の親として
使用します。

## `image_create_aliases`

`POST /1.0/images` に `aliases` フィールドを新たに追加します。イメージの
作成／インポート時にエイリアスを設定できるようになります。

## `container_stateless_copy`

`POST /1.0/containers/<name>` に `live` という属性を新たに導入します。
`false` に設定すると、実行状態を転送しようとしないように Incus に伝えます。

## `container_only_migration`

`container_only` という Boolean 型の属性を導入します。 `true` に設定すると
コンテナだけがコピーや移動されるようになります。

## `storage_zfs_clone_copy`

ZFS ストレージプールに `storage_zfs_clone_copy` という Boolean 型のプロパティを導入します。
`false` に設定すると、コンテナのコピーは `zfs send` と receive 経由で行われる
ようになります。これにより作成先のコンテナは作成元のコンテナに依存しないように
なり、 ZFS プールに依存するスナップショットを維持する必要がなくなります。
しかし、これは影響するプールのストレージの使用状況が以前より非効率的になる
という結果を伴います。
このプロパティのデフォルト値は `true` です。つまり明示的に `false` に設定
しない限り、空間効率の良いスナップショットが使われます。

## `unix_device_rename`

`path` を設定することによりコンテナ内部で `unix-block`/`unix-char` デバイスをリネーム
できるようにし、ホスト上のデバイスを指定する `source` 属性が追加されます。
`path` を設定せずに `source` を設定すると、 `path` は `source` と同じものとして
扱います。 `source` や `major`/`minor` を設定せずに `path` を設定すると
`source` は `path` と同じものとして扱います。ですので、最低どちらか 1 つは
設定しなければなりません。

## `storage_rsync_bwlimit`

ストレージエンティティを転送するために `rsync` が起動される場合に
`rsync.bwlimit` を設定すると使用できるソケット I/O の量に上限を
設定します。

## `network_vxlan_interface`

ネットワークに `tunnel.NAME.interface` オプションを新たに導入します。

このキーは VXLAN トンネルにホストのどのネットワークインターフェースを使うかを
制御します。

## `storage_btrfs_mount_options`

Btrfs ストレージプールに `btrfs.mount_options` プロパティを導入します。

このキーは Btrfs ストレージプールに使われるマウントオプションを制御します。

## `entity_description`

これはエンティティにコンテナ、スナップショット、ストレージプール、ボリュームの
ような説明を追加します。

## `image_force_refresh`

既存のイメージを強制的にリフレッシュできます。

## `storage_lvm_lv_resizing`

これはコンテナの root ディスクデバイス内に `size` プロパティを設定することで
論理ボリュームをリサイズできるようにします。

## `id_map_base`

これは `security.idmap.base` を新しく導入します。これにより分離されたコンテナ
に map auto-selection するプロセスをスキップし、ホストのどの UID/GID をベース
として使うかをユーザーが指定できるようにします。

## `file_symlinks`

これは file API 経由でシンボリックリンクを転送するサポートを追加します。
X-Incus-type に `symlink` を指定できるようになり、リクエストの内容はターゲットの
パスを指定します。

## `container_push_target`

`POST /1.0/containers/<name>` に `target` フィールドを新たに追加します。
これはマイグレーション中に作成元の Incus ホストが作成先に接続するために
利用可能です。

## `network_vlan_physical`

`physical` ネットワークデバイスで `vlan` プロパティが使用できるようにします。

設定すると、 `parent` インターフェース上で指定された VLAN にアタッチするように
Incus に指示します。 Incus はホスト上でその `parent` と VLAN を既存のインターフェース
で探します。
見つからない場合は作成します。
その後コンテナにこのインターフェースを直接アタッチします。

## `storage_images_delete`

これは指定したストレージプールからイメージのストレージボリュームを
ストレージ API で削除できるようにします。

## `container_edit_metadata`

これはコンテナの `metadata.yaml` と関連するテンプレートを
`/1.0/containers/<name>/metadata` 配下の URL にアクセスすることにより
API で編集できるようにします。コンテナからイメージを発行する前にコンテナを
編集できるようになります。

## `container_snapshot_stateful_migration`

これは stateful なコンテナのスナップショットを新しいコンテナにマイグレート
できるようにします。

## `storage_driver_ceph`

これは Ceph ストレージドライバーを追加します。

## `storage_ceph_user_name`

これは Ceph ユーザーを指定できるようにします。

## `instance_types`

これはコンテナの作成リクエストに `instance_type` フィールドを追加します。
値は Incus のリソース制限に展開されます。

## `storage_volatile_initial_source`

これはストレージプール作成中に Incus に渡された実際の作成元を記録します。

## `storage_ceph_force_osd_reuse`

これは Ceph ストレージドライバーに `ceph.osd.force_reuse` プロパティを
導入します。 `true` に設定すると Incus は別の Incus インスタンスで既に使用中の
OSD ストレージプールを再利用するようになります。

## `storage_block_filesystem_btrfs`

これは `ext4` と `xfs` に加えて Btrfs をストレージボリュームファイルシステムとして
サポートするようになります。

## `resources`

これは Incus が利用可能なシステムリソースを Incus デーモンに問い合わせできるようにします。

## `kernel_limits`

これは `nofile` でコンテナがオープンできるファイルの最大数といったプロセスの
リミットを設定できるようにします。形式は `limits.kernel.[リミット名]` です。

## `storage_api_volume_rename`

これはカスタムストレージボリュームをリネームできるようにします。

## `network_sriov`

これは SR-IOV を有効にしたネットワークデバイスのサポートを追加します。

## `console`

これはコンテナのコンソールデバイスとコンソールログを利用可能にします。

## `restrict_dev_incus`

`security.guestapi` コンテナ設定キーを新たに導入します。このキーは `/dev/incus`
インターフェースがコンテナで利用可能になるかを制御します。
`false` に設定すると、コンテナが Incus デーモンと連携するのを実質無効に
することになります。

## `migration_pre_copy`

これはライブマイグレーション中に最適化されたメモリー転送をできるようにします。

## `infiniband`

これは InfiniBand ネットワークデバイスを使用できるようにします。

## `dev_incus_events`

これは `/dev/incus` ソケットに Websocket API を追加します。

`/dev/incus` ソケット上で `/1.0/events` に接続すると、 Websocket 上で
イベントのストリームを受け取れるようになります。

## `proxy`

これはコンテナに `proxy` という新しいデバイスタイプを追加します。
これによりホストとコンテナ間で接続をフォワーディングできるようになります。

## `network_dhcp_gateway`

代替のゲートウェイを設定するための `ipv4.dhcp.gateway` ネットワーク設定キーを
新たに追加します。

## `file_get_symlink`

これは file API を使ってシンボリックリンクを取得できるようにします。

## `network_leases`

`/1.0/networks/NAME/leases` API エンドポイントを追加します。 Incus が管理する
DHCP サーバーが稼働するブリッジ上のリースデータベースに問い合わせできるように
なります。

## `unix_device_hotplug`

これは Unix デバイスに `required` プロパティのサポートを追加します。

## `storage_api_local_volume_handling`

これはカスタムストレージボリュームを同じあるいは異なるストレージプール間で
コピーしたり移動したりできるようにします。

## `operation_description`

すべての操作に `description` フィールドを追加します。

## `clustering`

Incus のクラスタリング API 。

これは次の新しいエンドポイントを含みます（詳細は [RESTful API](rest-api.md) を参照）:

* `GET /1.0/cluster`
* `UPDATE /1.0/cluster`

* `GET /1.0/cluster/members`

* `GET /1.0/cluster/members/<name>`
* `POST /1.0/cluster/members/<name>`
* `DELETE /1.0/cluster/members/<name>`

次の既存のエンドポイントは以下のように変更されます:

* `POST /1.0/containers` 新しい `target` クエリパラメーターを受け付けるようになります。
* `POST /1.0/storage-pools` 新しい `target` クエリパラメーターを受け付けるようになります
* `GET /1.0/storage-pool/<name>` 新しい `target` クエリパラメーターを受け付けるようになります
* `POST /1.0/storage-pool/<pool>/volumes/<type>` 新しい `target` クエリパラメーターを受け付けるようになります
* `GET /1.0/storage-pool/<pool>/volumes/<type>/<name>` 新しい `target` クエリパラメーターを受け付けるようになります
* `POST /1.0/storage-pool/<pool>/volumes/<type>/<name>` 新しい `target` クエリパラメーターを受け付けるようになります
* `PUT /1.0/storage-pool/<pool>/volumes/<type>/<name>` 新しい `target` クエリパラメーターを受け付けるようになります
* `PATCH /1.0/storage-pool/<pool>/volumes/<type>/<name>` 新しい `target` クエリパラメーターを受け付けるようになります
* `DELETE /1.0/storage-pool/<pool>/volumes/<type>/<name>` 新しい `target` クエリパラメーターを受け付けるようになります
* `POST /1.0/networks` 新しい `target` クエリパラメーターを受け付けるようになります
* `GET /1.0/networks/<name>` 新しい `target` クエリパラメーターを受け付けるようになります

## `event_lifecycle`

これはイベント API に `lifecycle` メッセージ種別を新たに追加します。

## `storage_api_remote_volume_handling`

これはリモート間でカスタムストレージボリュームをコピーや移動できるようにします。

## `nvidia_runtime`

コンテナに `nvidia_runtime` という設定オプションを追加します。これを `true` に
設定すると NVIDIA ランタイムと CUDA ライブラリーがコンテナに渡されます。

## `container_mount_propagation`

これはディスクデバイスタイプに `propagation` オプションを新たに追加します。
これによりカーネルのマウントプロパゲーションの設定ができるようになります。

## `container_backup`

コンテナのバックアップサポートを追加します。

これは次のエンドポイントを新たに追加します（詳細は [RESTful API](rest-api.md) を参照）:

* `GET /1.0/containers/<name>/backups`
* `POST /1.0/containers/<name>/backups`

* `GET /1.0/containers/<name>/backups/<name>`
* `POST /1.0/containers/<name>/backups/<name>`
* `DELETE /1.0/containers/<name>/backups/<name>`

* `GET /1.0/containers/<name>/backups/<name>/export`

次の既存のエンドポイントは以下のように変更されます:

* `POST /1.0/containers` 新たな作成元の種別 `backup` を受け付けるようになります

## `dev_incus_images`

コンテナに `security.guestapi.images` 設定オプションを追加します。これに
より `/dev/incus` 上で `/1.0/images/FINGERPRINT/export` API が利用可能に
なります。 nested Incus を動かすコンテナがホストから生のイメージを
取得するためにこれは利用できます。

## `container_local_cross_pool_handling`

これは同じ Incus インスタンス上のストレージプール間でコンテナをコピー・移動
できるようにします。

## `proxy_unix`

proxy デバイスで Unix ソケットと abstract Unix ソケットの両方のサポートを
追加します。これらは `unix:/path/to/unix.sock`（通常のソケット）あるいは
`unix:@/tmp/unix.sock`（abstract ソケット）のようにアドレスを指定して
利用可能です。

現状サポートされている接続は次のとおりです:

* `TCP <-> TCP`
* `UNIX <-> UNIX`
* `TCP <-> UNIX`
* `UNIX <-> TCP`

## `proxy_udp`

proxy デバイスで UDP のサポートを追加します。

現状サポートされている接続は次のとおりです:

* `TCP <-> TCP`
* `UNIX <-> UNIX`
* `TCP <-> UNIX`
* `UNIX <-> TCP`
* `UDP <-> UDP`
* `TCP <-> UDP`
* `UNIX <-> UDP`

## `clustering_join`

これにより `GET /1.0/cluster` がノードに参加する際にどのようなストレージプールと
ネットワークを作成する必要があるかについての情報を返します。また、それらを作成する
際にどのノード固有の設定キーを使う必要があるかについての情報も返します。
同様に `PUT /1.0/cluster` エンドポイントも同じ形式でストレージプールとネットワークに
ついての情報を受け付け、クラスタに参加する前にこれらが自動的に作成されるようになります。

## `proxy_tcp_udp_multi_port_handling`

複数のポートにトラフィックをフォワーディングできるようにします。フォワーディングは
ポートの範囲が転送元と転送先で同じ（たとえば `1.2.3.4 0-1000 -> 5.6.7.8 1000-2000`）
場合か転送元で範囲を指定し転送先で単一のポートを指定する
（たとえば `1.2.3.4 0-1000 -> 5.6.7.8 1000`）場合に可能です。

## `network_state`

ネットワークの状態を取得できるようになります。

これは次のエンドポイントを新たに追加します（詳細は [RESTful API](rest-api.md) を参照）:

* `GET /1.0/networks/<name>/state`

## `proxy_unix_dac_properties`

これは抽象的 Unix ソケットではない Unix ソケットに GID、UID、パーミションのプロパティを追加します。

## `container_protection_delete`

`security.protection.delete` フィールドを設定できるようにします。 `true` に設定すると
コンテナが削除されるのを防ぎます。スナップショットはこの設定により影響を受けません。

## `proxy_priv_drop`

proxy デバイスに `security.uid` と `security.gid` を追加します。これは root 権限を
落とし（訳注: 非 root 権限で動作させるという意味です）、 Unix ソケットに接続する
際に用いられる UID/GID も変更します。

## `pprof_http`

これはデバッグ用の HTTP サーバーを起動するために、新たに `core.debug_address`
オプションを追加します。

このサーバーは現在`pprof` API を含んでおり、従来の`cpu-profile`, `memory-profile`
と`print-goroutines`デバッグオプションを置き換えるものです。

## `proxy_haproxy_protocol`

proxy デバイスに `proxy_protocol` キーを追加します。これは HAProxy PROXY プロトコルヘッダ
の使用を制御します。

## `network_hwaddr`

ブリッジの MAC アドレスを制御する `bridge.hwaddr` キーを追加します。

## `proxy_nat`

これは最適化された UDP/TCP プロキシを追加します。設定上可能であれば
プロキシ処理は proxy デバイスの代わりに `iptables` 経由で行われるように
なります。

## `network_nat_order`

Incus ブリッジに `ipv4.nat.order` と `ipv6.nat.order` 設定キーを導入します。
これらのキーは Incus のルールをチェイン内の既存のルールの前に置くか後に置くかを
制御します。

## `container_full`

これは `GET /1.0/containers` に `recursion=2` という新しいモードを導入します。
これにより状態、スナップショットとバックアップの構造を含むコンテナのすべての構造を
取得できるようになります。

この結果 [`incus list`](incus_list.md) は必要なすべての情報を 1 つのクエリで取得できるように
なります。

## `backup_compression`

これは新たに `backups.compression_algorithm` 設定キーを導入します。
これによりバックアップの圧縮の設定が可能になります。

## `nvidia_runtime_config`

これは `nvidia.runtime` と `libnvidia-container` ライブラリーを使用する際に追加の
いくつかの設定キーを導入します。これらのキーは NVIDIA container の対応する
環境変数にほぼそのまま置き換えられます:

* `nvidia.driver.capabilities` => `NVIDIA_DRIVER_CAPABILITIES`
* `nvidia.require.cuda` => `NVIDIA_REQUIRE_CUDA`
* `nvidia.require.driver` => `NVIDIA_REQUIRE_DRIVER`

## `storage_api_volume_snapshots`

ストレージボリュームスナップショットのサポートを追加します。これらは
コンテナスナップショットのように振る舞いますが、ボリュームに対してのみ
作成できます。

これにより次の新しいエンドポイントが追加されます（詳細は [RESTful API](rest-api.md) を参照）:

* `GET /1.0/storage-pools/<pool>/volumes/<type>/<name>/snapshots`
* `POST /1.0/storage-pools/<pool>/volumes/<type>/<name>/snapshots`

* `GET /1.0/storage-pools/<pool>/volumes/<type>/<volume>/snapshots/<name>`
* `PUT /1.0/storage-pools/<pool>/volumes/<type>/<volume>/snapshots/<name>`
* `POST /1.0/storage-pools/<pool>/volumes/<type>/<volume>/snapshots/<name>`
* `DELETE /1.0/storage-pools/<pool>/volumes/<type>/<volume>/snapshots/<name>`

## `storage_unmapped`

ストレージボリュームに新たに `security.unmapped` という Boolean 設定を導入します。

`true` に設定するとボリューム上の現在のマップをフラッシュし、以降の
idmap のトラッキングとボリューム上のリマッピングを防ぎます。

これは隔離されたコンテナ間でデータを共有するために使用できます。
この際コンテナを書き込みアクセスを要求するコンテナにアタッチした
後にデータを共有します。

## `projects`

新たに project API を追加します。プロジェクトの作成、更新、削除ができます。

現時点では、プロジェクトは、コンテナ、プロファイル、イメージを保持できます。そして、プロジェクトを切り替えることで、独立した Incus リソースのビューを見せられます。

## `network_vxlan_ttl`

新たにネットワークの設定に `tunnel.NAME.ttl` が指定できるようになります。これにより、VXLAN トンネルの TTL を増加させることができます。

## `container_incremental_copy`

新たにコンテナのインクリメンタルコピーができるようになります。`--refresh` オプションを指定してコンテナをコピーすると、見つからないファイルや、更新されたファイルのみを
コピーします。コンテナが存在しない場合は、通常のコピーを実行します。

## `usb_optional_vendorid`

名前が暗示しているように、コンテナにアタッチされた USB デバイスの
`vendorid` フィールドが省略可能になります。これによりすべての USB デバイスが
コンテナに渡されます (GPU に対してなされたのと同様)。

## `snapshot_scheduling`

これはスナップショットのスケジューリングのサポートを追加します。これにより
3 つの新しい設定キーが導入されます。 `snapshots.schedule`、`snapshots.schedule.stopped`、
そして `snapshots.pattern` です。スナップショットは最短で 1 分間隔で自動的に
作成されます。

## `snapshots_schedule_aliases`

スナップショットのスケジュールはスケジュールエイリアスのカンマ区切りリストで設定できます。
インスタンスには `<@hourly> <@daily> <@midnight> <@weekly> <@monthly> <@annually> <@yearly> <@startup>`、
ストレージボリュームには `<@hourly> <@daily> <@midnight> <@weekly> <@monthly> <@annually> <@yearly>` のエイリアスが利用できます。

## `container_copy_project`

コピー元のコンテナの dict に `project` フィールドを導入します。これにより
プロジェクト間でコンテナをコピーあるいは移動できるようになります。

## `clustering_server_address`

これはサーバーのネットワークアドレスを REST API のクライアントネットワーク
アドレスと異なる値に設定することのサポートを追加します。クライアントは
新しい `cluster.https_address` 設定キーを初期のサーバーのアドレスを指定するために
に設定できます。新しいサーバーが参加する際、クライアントは参加するサーバーの
`core.https_address` 設定キーを参加するサーバーがリッスンすべきアドレスに設定でき、
`PUT /1.0/cluster` API の `server_address` キーを参加するサーバーが
クラスタリングトラフィックに使用すべきアドレスに設定できます（`server_address`
の値は自動的に参加するサーバーの `cluster.https_address` 設定キーに
コピーされます）。

## `clustering_image_replication`

クラスタ内のノードをまたいだイメージのレプリケーションを可能にします。
新しい `cluster.images_minimal_replica` 設定キーが導入され、イメージの
レプリケーションに対するノードの最小数を指定するのに使用できます。

## `container_protection_shift`

`security.protection.shift` の設定を可能にします。これによりコンテナの
ファイルシステム上で UID/GID をシフト (再マッピング) させることを防ぎます。

## `snapshot_expiry`

これはスナップショットの有効期限のサポートを追加します。タスクは 1 分おきに実行されます。
`snapshots.expiry` 設定オプションは、`1M 2H 3d 4w 5m 6y` （それぞれ 1 分、2 時間、3 日、4 週間、5 ヶ月、6 年）といった形式を取ります。
この指定ではすべての部分を使う必要はありません。

作成されるスナップショットには、指定した式に基づいて有効期限が設定されます。
`expires_at` で定義される有効期限は、API や [`incus config edit`](incus_config_edit.md) コマンドを使って手動で編集できます。
有効な有効期限が設定されたスナップショットはタスク実行時に削除されます。
有効期限は `expires_at` に空文字列や `0001-01-01T00:00:00Z`（zero time）を設定することで無効化できます。
`snapshots.expiry` が設定されていない場合はこれがデフォルトです。

これは次のような新しいエンドポイントを追加します（詳しくは [RESTful API](rest-api.md) をご覧ください）:

* `PUT /1.0/containers/<name>/snapshots/<name>`

## `snapshot_expiry_creation`

コンテナ作成に `expires_at` を追加し、作成時にスナップショットの有効期限を上書きできます。

## `network_leases_location`

ネットワークのリースリストに `Location` フィールドを導入します。
これは、特定のリースがどのノードに存在するかを問い合わせるときに使います。

## `resources_cpu_socket`

ソケットの情報が入れ替わる場合に備えて CPU リソースにソケットフィールドを追加します。

## `resources_gpu`

サーバーリソースに新規に GPU 構造を追加し、システム上で利用可能なすべての GPU を一覧表示します。

## `resources_numa`

すべての CPU と GPU に対する NUMA ノードを表示します。

## `kernel_features`

サーバーの環境からオプショナルなカーネル機能の使用可否状態を取得します。

## `id_map_current`

内部的な `volatile.idmap.current` キーを新規に導入します。これはコンテナに
対する現在のマッピングを追跡するのに使われます。

実質的には以下が利用可能になります:

* `volatile.last_state.idmap` => ディスク上の idmap
* `volatile.idmap.current` => 現在のカーネルマップ
* `volatile.idmap.next` => 次のディスク上の idmap

これはディスク上の map が変更されていないがカーネルマップは変更されている
（例: `idmapped mounts`）ような環境を実装するために必要です。

## `event_location`

API イベントの世代の場所を公開します。

## `storage_api_remote_volume_snapshots`

ストレージボリュームをそれらのスナップショットを含んで移行できます。

## `network_nat_address`

これは Incus ブリッジに `ipv4.nat.address` と `ipv6.nat.address` 設定キーを導入します。
これらのキーはブリッジからの送信トラフィックに使うソースアドレスを制御します。

## `container_nic_routes`

これは `nic` タイプのデバイスに `ipv4.routes` と `ipv6.routes` プロパティを導入します。
ホストからコンテナへの NIC への静的ルートが追加できます。

## `cluster_internal_copy`

これは通常の `POST /1.0/containers` を実行することでクラスタノード間で
コンテナをコピーすることを可能にします。この際 Incus はマイグレーションが
必要かどうかを内部的に判定します。

## `seccomp_notify`

カーネルが `seccomp` ベースの syscall インターセプトをサポートする場合に
登録された syscall が実行されたことをコンテナから Incus に通知することが
できます。 Incus はそれを受けて様々なアクションをトリガーするかを決定します。

## `lxc_features`

これは `GET /1.0` ルート経由で [`incus info`](incus_info.md) コマンドの出力に `lxc_features`
セクションを導入します。配下の LXC ライブラリーに存在するキー・フィーチャーに
対するチェックの結果を出力します。

## `container_nic_ipvlan`

これは `nic` デバイスに `ipvlan` のタイプを導入します。

## `network_vlan_sriov`

これは SR-IOV デバイスに VLAN（`vlan`）と MAC フィルタリング（`security.mac_filtering`）のサポートを導入します。

## `storage_cephfs`

ストレージプールドライバーとして CephFS のサポートを追加します。これは
カスタムボリュームとしての利用のみが可能になり、イメージとコンテナは
CephFS ではなく Ceph（RBD）上に構築する必要があります。

## `container_nic_ipfilter`

これは `bridged` の NIC デバイスに対してコンテナの IP フィルタリング
（`security.ipv4_filtering` と `security.ipv6_filtering`）を導入します。

## `resources_v2`

`/1.0/resources` のリソース API を見直しました。主な変更は以下の通りです:

* CPU
   * ソケット、コア、スレッドのトラッキングのレポートを修正しました
   * コア毎の NUMA ノードのトラッキング
   * ソケット毎のベースとターボの周波数のトラッキング
   * コア毎の現在の周波数のトラッキング
   * CPU のキャッシュ情報の追加
   * CPU アーキテクチャをエクスポート
   * スレッドのオンライン／オフライン状態を表示
* メモリー
   * HugePages のトラッキングを追加
   * NUMA ノード毎でもメモリー消費を追跡
* GPU
   * DRM 情報を別の構造体に分離
   * DRM 構造体内にデバイスの名前とノードを公開
   * NVIDIA 構造体内にデバイスの名前とノードを公開
   * SR-IOV VF のトラッキングを追加

## `container_exec_user_group_cwd`

`POST /1.0/containers/NAME/exec` の実行時に `User`、`Group` と `Cwd` を指定するサポートを追加します。

## `container_syscall_intercept`

`security.syscalls.intercept.*` 設定キーを追加します。これはどのシステムコールを Incus がインターセプトし昇格された権限で処理するかを制御します。

## `container_disk_shift`

`disk` デバイスに `shift` プロパティを追加します。これは `idmapped mounts` のオーバーレイの使用を制御します。

## `storage_shifted`

ストレージボリュームに新しく `security.shifted` という Boolean の設定を導入します。

これを `true` に設定すると複数の隔離されたコンテナが、それらすべてがファイルシステムに
書き込み可能にしたまま、同じストレージボリュームにアタッチするのを許可します。

これは `idmapped mounts` をオーバーレイファイルシステムとして使用します。

## `resources_infiniband`

リソース API の一部として InfiniBand キャラクタデバイス（`issm`、`umad`、`uverb`）の情報を公開します。

## `daemon_storage`

これは `storage.images_volume` と `storage.backups_volume` という 2 つの新しい設定項目を導入します。これらは既存のプール上のストレージボリュームがデーモン全体のイメージとバックアップを保管するのに使えるようにします。

## `instances`

これはインスタンスの概念を導入します。現状ではインスタンスの唯一の種別は `container` です。

## `image_types`

これはイメージに新しく Type フィールドのサポートを導入します。 Type フィールドはイメージがどういう種別かを示します。

## `resources_disk_sata`

ディスクリソース API の構造体を次の項目を含むように拡張します。

* SATA デバイス（種別）の適切な検出
* デバイスパス
* ドライブの RPM
* ブロックサイズ
* ファームウェアバージョン
* シリアルナンバー

## `clustering_roles`

これはクラスタのエントリーに `roles` という新しい属性を追加し、クラスタ内のメンバーが提供する role の一覧を公開します。

## `images_expiry`

イメージの有効期限を設定できます。

## `resources_network_firmware`

ネットワークカードのエントリーに `FirmwareVersion` フィールドを追加します。

## `backup_compression_algorithm`

バックアップを作成する（`POST /1.0/containers/<name>/backups`）際に `compression_algorithm` プロパティのサポートを追加します。

このプロパティを設定するとデフォルト値（`backups.compression_algorithm`）をオーバーライドすることができます。

## `ceph_data_pool_name`

Ceph RBD を使ってストレージプールを作成する際にオプショナルな引数（`ceph.osd.data_pool_name`）のサポートを追加します。
この引数が指定されると、プールはメタデータは `pool_name` で指定されたプールに保持しつつ実際のデータは `data_pool_name` で指定されたプールに保管するようになります。

## `container_syscall_intercept_mount`

`security.syscalls.intercept.mount`、`security.syscalls.intercept.mount.allowed`、`security.syscalls.intercept.mount.shift` 設定キーを追加します。
これらは `mount` システムコールを Incus にインターセプトさせるかどうか、昇格されたパーミションでどのように処理させるかを制御します。

## `compression_squashfs`

イメージやバックアップを SquashFS ファイルシステムの形式でインポート／エクスポートするサポートを追加します。

## `container_raw_mount`

ディスクデバイスに raw mount オプションを渡すサポートを追加します。

## `container_nic_routed`

`routed` `nic` デバイスタイプを導入します。

## `container_syscall_intercept_mount_fuse`

`security.syscalls.intercept.mount.fuse` キーを追加します。これはファイルシステムのマウントを fuse 実装にリダイレクトするのに使えます。
このためにはたとえば `security.syscalls.intercept.mount.fuse=ext4=fuse2fs` のように設定します。

## `container_disk_ceph`

既存の Ceph RBD もしくは CephFS を直接 Incus コンテナに接続できます。

## `virtual-machines`

仮想マシンサポートが追加されます。

## `image_profiles`

新しいコンテナを起動するときに、イメージに適用するプロファイルのリストが指定できます。

## `clustering_architecture`

クラスタメンバーに `architecture` 属性を追加します。
この属性はクラスタメンバーのアーキテクチャを示します。

## `resources_disk_id`

リソース API のディスクのエントリーに `device_id` フィールドを追加します。

## `storage_lvm_stripes`

通常のボリュームと thin pool ボリューム上で LVM ストライプを使う機能を追加します。

## `vm_boot_priority`

ブートの順序を制御するため NIC とディスクデバイスに `boot.priority` プロパティを追加します。

## `unix_hotplug_devices`

Unix のキャラクタデバイスとブロックデバイスのホットプラグのサポートを追加します。

## `api_filtering`

インスタンスとイメージに対する GET リクエストの結果をフィルタリングする機能を追加します。

## `instance_nic_network`

NIC デバイスの `network` プロパティのサポートを追加し、管理されたネットワークへ NIC をリンクできるようにします。
これによりネットワーク設定の一部を引き継ぎ、 IP 設定のより良い検証を行うことができます。

## `clustering_sizing`

データベースの投票者とスタンバイに対してカスタムの値を指定するサポートです。
`cluster.max_voters` と `cluster.max_standby` という新しい設定キーが導入され、データベースの投票者とスタンバイの理想的な数を指定できます。

## `firewall_driver`

`ServerEnvironment` 構造体にファイアーウォールのドライバーが使用されていることを示す `Firewall` プロパティを追加します。

## `storage_lvm_vg_force_reuse`

既存の空でないボリュームグループからストレージボリュームを作成する機能を追加します。
このオプションの使用には注意が必要です。
というのは、同じボリュームグループ内に Incus 以外で作成されたボリュームとボリューム名が衝突しないことを Incus が保証できないからです。
このことはもし名前の衝突が起きたときは Incus 以外で作成されたボリュームを Incus が削除してしまう可能性があることを意味します。

## `container_syscall_intercept_hugetlbfs`

mount システムコール・インターセプションが有効にされ `hugetlbfs` が許可されたファイルシステムとして指定された場合、 Incus は別の `hugetlbfs` インスタンスを UID と GID をコンテナの root の UID と GID に設定するマウントオプションを指定してコンテナにマウントします。
これによりコンテナ内のプロセスが huge page を確実に利用できるようにします。

## `limits_hugepages`

コンテナが使用できる huge page の数を `hugetlb` cgroup を使って制限できるようにします。
この機能を使用するには `hugetlb` cgroup が利用可能になっている必要があります。
注意: `hugetlbfs` ファイルシステムの mount システムコールをインターセプトするときは、ホストの huge page のリソースをコンテナが使い切ってしまわないように huge page を制限することを推奨します。

## `container_nic_routed_gateway`

この拡張は `ipv4.gateway` と `ipv6.gateway` という NIC の設定キーを追加します。
指定可能な値は `auto` か `none` のいずれかです。
値を指定しない場合のデフォルト値は `auto` です。
`auto` に設定した場合は、デフォルトゲートウェイがコンテナ内部に追加され、ホスト側のインターフェースにも同じゲートウェイアドレスが追加されるという現在の挙動と同じになります。
`none` に設定すると、デフォルトゲートウェイもアドレスもホスト側のインターフェースには追加されません。
これにより複数のルートを持つ NIC デバイスをコンテナに追加できます。

## `projects_restrictions`

この拡張はプロジェクトに `restricted` という設定キーを追加します。
これによりプロジェクト内でセキュリティーセンシティブな機能を使うのを防ぐことができます。

## `custom_volume_snapshot_expiry`

この拡張はカスタムボリュームのスナップショットに有効期限を設定できるようにします。
有効期限は `snapshots.expiry` 設定キーにより個別に設定することも出来ますし、親のカスタムボリュームに設定してそこから作成されたすべてのスナップショットに自動的にその有効期限を適用することも出来ます。

## `volume_snapshot_scheduling`

この拡張はカスタムボリュームのスナップショットにスケジュール機能を追加します。
`snapshots.schedule` と `snapshots.pattern` という 2 つの設定キーが新たに追加されます。
スナップショットは最短で 1 分毎に作成可能です。

## `trust_ca_certificates`

この拡張により提供された CA（`server.ca`）によって信頼されたクライアント証明書のチェックが可能になります。
`core.trust_ca_certificates` を `true` に設定すると有効にできます。
有効な場合、クライアント証明書のチェックを行い、チェックが OK であれば信頼されたパスワードの要求はスキップします。
ただし、提供された CRL（`ca.crl`）に接続してきたクライアント証明書が含まれる場合は例外です。
この場合は、パスワードが求められます。

## `snapshot_disk_usage`

この拡張はスナップショットのディスク使用量を示す `/1.0/instances/<name>/snapshots/<snapshot>` の出力に `size` フィールドを新たに追加します。

## `clustering_edit_roles`

この拡張はクラスタメンバーに書き込み可能なエンドポイントを追加し、ロールの編集を可能にします。

## `container_nic_routed_host_address`

この拡張は NIC の設定キーに `ipv4.host_address` と `ipv6.host_address` を追加し、ホスト側の veth インターフェースの IP アドレスを制御できるようにします。
これは同時に複数の routed NIC を使用し、予測可能な next-hop のアドレスを使用したい場合に有用です。

さらにこの拡張は `ipv4.gateway` と `ipv6.gateway` の NIC 設定キーの振る舞いを変更します。
auto に設定するとコンテナはデフォルトゲートウェイをそれぞれ `ipv4.host_address` と `ipv6.host_address` で指定した値にします。

デフォルト値は次の通りです。

`ipv4.host_address`: `169.254.0.1`
`ipv6.host_address`: `fe80::1` 

これは以前のデフォルトの挙動と後方互換性があります。

## `container_nic_ipvlan_gateway`

この拡張は `ipv4.gateway` と `ipv6.gateway` の NIC 設定キーを追加し `auto` か `none` の値を指定できます。
指定しない場合のデフォルト値は `auto` です。
この場合は従来同様の挙動になりコンテナ内部に追加されるデフォルトゲートウェイと同じアドレスがホスト側のインターフェースにも追加されます。
`none` に設定された場合、ホスト側のインターフェースにはデフォルトゲートウェイもアドレスも追加されません。
これによりコンテナに IPVLAN の NIC デバイスを複数追加することができます。

## `resources_usb_pci`

この拡張は `/1.0/resources` の出力に USB と PC デバイスを追加します。

## `resources_cpu_threads_numa`

この拡張は `numa_node` フィールドをコアごとではなくスレッドごとに記録するように変更します。
これは一部のハードウェアでスレッドを異なる NUMA ドメインに入れる場合があるようなのでそれに対応するためのものです。

## `resources_cpu_core_die`

それぞれのコアごとに `die_id` 情報を公開します。

## `api_os`

この拡張は `/1.0` 内に `os` と `os_version` の 2 つのフィールドを追加します。

これらの値はシステム上の OS-release のデータから取得されます。

## `container_nic_routed_host_table`

この拡張は `ipv4.host_table` と `ipv6.host_table` という NIC の設定キーを導入します。
これで指定した ID のカスタムポリシーのルーティングテーブルにインスタンスの IP のための静的ルートを追加できます。

## `container_nic_ipvlan_host_table`

この拡張は `ipv4.host_table` と `ipv6.host_table` という NIC の設定キーを導入します。
これで指定した ID のカスタムポリシーのルーティングテーブルにインスタンスの IP のための静的ルートを追加できます。

## `container_nic_ipvlan_mode`

この拡張は `mode` という NIC の設定キーを導入します。
これにより `ipvlan` モードを `l2` か `l3s` のいずれかに切り替えられます。
指定しない場合、デフォルトは `l3s` （従来の挙動）です。

`l2` モードでは `ipv4.address` と `ipv6.address` キーは CIDR か単一アドレスの形式を受け付けます。
単一アドレスの形式を使う場合、デフォルトのサブネットのサイズは IPv4 では /24 、 IPv6 では /64 となります。

`l2` モードでは `ipv4.gateway` と `ipv6.gateway` キーは単一の IP アドレスのみを受け付けます。

## `resources_system`

この拡張は `/1.0/resources` の出力にシステム情報を追加します。

## `images_push_relay`

この拡張はイメージのコピーに push と relay モードを追加します。
また以下の新しいエンドポイントも追加します:

* `POST 1.0/images/<fingerprint>/export`

## `network_dns_search`

この拡張はネットワークに `dns.search` という設定オプションを追加します。

## `container_nic_routed_limits`

この拡張は routed NIC に `limits.ingress`, `limits.egress`, `limits.max` を追加します。

## `instance_nic_bridged_vlan`

この拡張は `bridged` NIC に `vlan` と `vlan.tagged` の設定を追加します。

`vlan` には参加するタグなし VLAN を指定し、 `vlan.tagged` は参加するタグ VLAN のカンマ区切りリストを指定します。

## `network_state_bond_bridge`

この拡張は `/1.0/networks/NAME/state` API に bridge と bond のセクションを追加します。

これらはそれぞれの特定のタイプに関連する追加の状態の情報を含みます。

Bond:

* Mode
* Transmit hash
* Up delay
* Down delay
* MII frequency
* MII state
* Lower devices

Bridge:

* ID
* Forward delay
* STP mode
* Default VLAN
* VLAN filtering
* Upper devices

## `resources_cpu_isolated`

この拡張は CPU スレッドに `Isolated` プロパティを追加します。
これはスレッドが物理的には `Online` ですがタスクを受け付けないように設定しているかを示します。

## `usedby_consistency`

この拡張により、可能な時は `UsedBy` が適切な `?project=` と `?target=` に対して一貫性があるようになるはずです。

`UsedBy` を持つ 5 つのエンティティーは以下の通りです:

* プロファイル
* プロジェクト
* ネットワーク
* ストレージプール
* ストレージボリューム

## `custom_block_volumes`

この拡張によりカスタムブロックボリュームを作成しインスタンスにアタッチできるようになります。
カスタムストレージボリュームの作成時に `--type` フラグが新規追加され、 `fs` と `block` の値を受け付けます。

## `clustering_failure_domains`

この拡張は `PUT /1.0/cluster/<node>` API に `failure_domain` フィールドを追加します。
これはノードの failure domain を設定するのに使えます。

## `container_syscall_filtering_allow_deny_syntax`

いくつかのシステムコールに関連したコンテナの設定キーが更新されました。

* `security.syscalls.deny_default`
* `security.syscalls.deny_compat`
* `security.syscalls.deny`
* `security.syscalls.allow`

## `resources_gpu_mdev`

`/1.0/resources` の利用可能な媒介デバイス (mediated device) のプロファイルとデバイスを公開します。

## `console_vga_type`

この拡張は `/1.0/console` エンドポイントが `?type=` 引数を取るように拡張します。
これは `console`（デフォルト）か `vga`（この拡張で追加される新しいタイプ）を指定可能です。

`/1.0/<instance name>/console?type=vga` に `POST` する際はメタデータフィールド内の操作の結果ウェブソケットにより返されるデータはターゲットの仮想マシンの SPICE Unix ソケットにアタッチされた双方向のプロキシになります。

## `projects_limits_disk`

利用可能なプロジェクトの設定キーに `limits.disk` を追加します。
これが設定されるとプロジェクト内でインスタンスボリューム、カスタムボリューム、イメージボリュームが使用できるディスクスペースの合計の量を制限できます。

## `network_type_macvlan`

ネットワークタイプ `macvlan` のサポートを追加し、このネットワークタイプに `parent` 設定キーを追加します。
これは NIC デバイスインターフェースを作る際にどの親インターフェースを使用するべきかを指定します。

さらに `macvlan` の NIC に `network` 設定キーを追加します。
これは NIC デバイスの設定の基礎として使う network を指定します。

## `network_type_sriov`

ネットワークタイプ `sriov` のサポートを追加し、このネットワークタイプに `parent` 設定キーを追加します。
これは NIC デバイスインターフェースを作る際にどの親インターフェースを使用するべきかを指定します。

さらに `sriov` の NIC に `network` 設定キーを追加します。
これは NIC デバイスの設定の基礎として使う network を指定します。

## `container_syscall_intercept_bpf_devices`

この拡張はコンテナ内で `bpf` のシステムコールをインターセプトする機能を提供します。具体的には device cgroup の `bpf` のプログラムを管理できるようにします。

## `network_type_ovn`

ネットワークタイプ `ovn` のサポートを追加し、 `bridge` タイプのネットワークを `parent` として設定できるようにします。

`ovn` という新しい NIC のデバイスタイプを追加します。これにより `network` 設定キーにどの `ovn` のタイプのネットワークに接続すべきかを指定できます。

さらにすべての `ovn` ネットワークと NIC デバイスに適用される 2 つのグローバルの設定キーを追加します:

* `network.ovn.integration_bridge` - 使用する OVS 統合ブリッジ
* `network.ovn.northbound_connection` - OVN northbound データベース接続文字列

## `projects_networks`

プロジェクトに `features.networks` 設定キーを追加し、プロジェクトがネットワークを保持できるようにします。

## `projects_networks_restricted_uplinks`

プロジェクトに `restricted.networks.uplinks` 設定キーを追加し、プロジェクト内で作られたネットワークがそのアップリンクのネットワークとしてどのネットワークが使えるかを（カンマ区切りリストで）指定します。

## `custom_volume_backup`

カスタムボリュームのバックアップサポートを追加します。

この拡張は以下の新しい API エンドポイント （詳細は [RESTful API](rest-api.md) を参照）を含みます:

* `GET /1.0/storage-pools/<pool>/<type>/<volume>/backups`
* `POST /1.0/storage-pools/<pool>/<type>/<volume>/backups`

* `GET /1.0/storage-pools/<pool>/<type>/<volume>/backups/<name>`
* `POST /1.0/storage-pools/<pool>/<type>/<volume>/backups/<name>`
* `DELETE /1.0/storage-pools/<pool>/<type>/<volume>/backups/<name>`

* `GET /1.0/storage-pools/<pool>/<type>/<volume>/backups/<name>/export`

以下の既存のエンドポイントが変更されます:

* `POST /1.0/storage-pools/<pool>/<type>/<volume>` が新しいソースタイプとして `backup` を受け付けます

## `backup_override_name`

`InstanceBackupArgs` に `Name` フィールドを追加し、バックアップをリストアする際に別のインスタンス名を指定できるようにします。

`StoragePoolVolumeBackupArgs` に `Name` と `PoolName` フィールドを追加し、カスタムボリュームのバックアップをリストアする際に別のボリューム名を指定できるようにします。

## `storage_rsync_compression`

ストレージプールに `rsync.compression` 設定キーを追加します。
このキーはストレージプールをマイグレートする際に `rsync` での圧縮を無効にするために使うことができます。

## `network_type_physical`

新たに `physical` というネットワークタイプのサポートを追加し、 `ovn` ネットワークのアップリンクとして使用できるようにします。

`physical` ネットワークの `parent` で指定するインターフェースは `ovn` ネットワークのゲートウェイに接続されます。

## `network_ovn_external_subnets`

`ovn` ネットワークがアップリンクネットワークの外部のサブネットを使用できるようにします。

`physical` ネットワークに `ipv4.routes` と `ipv6.routes` の設定を追加します。
これは子供の OVN ネットワークで `ipv4.routes.external` と `ipv6.routes.external` の設定で使用可能な外部のルートを指定します。

プロジェクトに `restricted.networks.subnets` 設定を追加します。
これはプロジェクト内の OVN ネットワークで使用可能な外部のサブネットを指定します（未設定の場合はアップリンクネットワークで定義されるすべてのルートが使用可能です）。

## `network_ovn_nat`

`ovn` ネットワークに `ipv4.nat` と `ipv6.nat` の設定を追加します。

これらの設定（訳注: `ipv4.nat` や `ipv6.nat`）を未設定でネットワークを作成する際、（訳注: `ipv4.address` や `ipv6.address` が未設定あるいは `auto` の場合に）対応するアドレス （訳注: `ipv4.nat` であれば `ipv4.address`、`ipv6.nat` であれば `ipv6.address`）がサブネット用に生成される場合は適切な NAT が生成され、`ipv4.nat` や `ipv6.nat` は `true` に設定されます。

この設定がない場合は値は `false` として扱われます。

## `network_ovn_external_routes_remove`

`ovn` ネットワークから `ipv4.routes.external` と `ipv6.routes.external` の設定を削除します。

ネットワークと NIC レベルの両方で指定するのではなく、 `ovn` NIC タイプ上で等価な設定を使えます。

## `tpm_device_type`

`tpm` デバイスタイプを導入します。

## `storage_zfs_clone_copy_rebase`

`zfs.clone_copy` に `rebase` という値を導入します。
この設定で Incus は先祖の系列上の `image` データセットを追跡し、その最上位に対して send/receive を実行します。

## `gpu_mdev`

これは仮想 CPU のサポートを追加します。
GPU デバイスに `mdev` 設定キーを追加し、`i915-GVTg_V5_4` のようなサポートされる `mdev` のタイプを指定します。

## `resources_pci_iommu`

これはリソース API の PCI エントリーに `IOMMUGroup` フィールドを追加します。

## `resources_network_usb`

リソース API のネットワークカードエントリーに `usb_address` フィールドを追加します。

## `resources_disk_address`

リソース API のディスクエントリーに `usb_address` と `pci_address` フィールドを追加します。

## `network_physical_ovn_ingress_mode`

`physical` ネットワークに `ovn.ingress_mode` 設定を追加します。

OVN NIC ネットワークの外部 IP アドレスがアップリンクネットワークにどのように広告されるかの方法を設定します。

`l2proxy`（proxy ARP/NDP）か `routed` のいずれかを指定します。

## `network_ovn_dhcp`

`ovn` ネットワークに `ipv4.dhcp` と `ipv6.dhcp` の設定を追加します。

DHCP（と IPv6 の RA）を無効にできます。デフォルトはオンです。

## `network_physical_routes_anycast`

`physical` ネットワークに `ipv4.routes.anycast` と `ipv6.routes.anycast` の Boolean の設定を追加します。デフォルトは `false` です。

`ovn.ingress_mode=routed` と共に使うと physical ネットワークをアップリンクとして使う OVN ネットワークでサブネット／ルートのオーバーラップ検出を緩和できます。

## `projects_limits_instances`

`limits.instances` を利用可能なプロジェクトの設定キーに追加します。
設定するとプロジェクト内で使われるインスタンス（VM とコンテナ）の合計数を制限します。

## `network_state_vlan`

これは `/1.0/networks/NAME/state` API に `vlan` セクションを追加します。

これらは VLAN インターフェースに関連する追加の状態の情報を含みます。

* `lower_device`
* `vid`

## `instance_nic_bridged_port_isolation`

これは `bridged` NIC に `security.port_isolation` のフィールドを追加します。

## `instance_bulk_state_change`

一括状態変更（詳細は [REST API](rest-api.md) を参照）のために次のエンドポイントを追加します:

* `PUT /1.0/instances`

## `network_gvrp`

これはオプショナルな `gvrp` プロパティを `macvlan` と `physical` ネットワークに追加し、
さらに `ipvlan`、`macvlan`、`routed`、`physical` NIC デバイスにも追加します。

設定された場合は、これは VLAN が GARP VLAN Registration Protocol を使って登録すべきかどうかを指定します。
デフォルトは `false` です。

## `instance_pool_move`

これは `POST /1.0/instances/NAME` API に `pool` フィールドを追加し、プール間でインスタンスのルートディスクを簡単に移動できるようにします。

## `gpu_sriov`

これは SR-IOV を有効にした GPU のサポートを追加します。
これにより `sriov` という GPU タイプのプロパティが追加されます。

## `pci_device_type`

これは `pci` デバイスタイプを追加します。

## `storage_volume_state`

`/1.0/storage-pools/POOL/volumes/VOLUME/state` API エンドポイントを新規追加しボリュームの使用量を取得できるようにします。

## `network_acl`

これは `/1.0/network-acls` の API エンドポイントプリフィクス以下の API にネットワークの ACL のコンセプトを追加します。

## `migration_stateful`

`migration.stateful` という設定キーを追加します。

## `disk_state_quota`

これは `disk` デバイスに `size.state` というデバイス設定キーを追加します。

## `storage_ceph_features`

ストレージプールに `ceph.rbd.features` 設定キーを追加し、新規ボリュームに使用する RBD の機能を制御します。

## `projects_compression`

`backups.compression_algorithm` と `images.compression_algorithm` 設定キーを追加します。
これらによりプロジェクトごとのバックアップとイメージの圧縮の設定が出来るようになります。

## `projects_images_remote_cache_expiry`

プロジェクトに `images.remote_cache_expiry` 設定キーを追加します。
これを設定するとキャッシュされたリモートのイメージが指定の日数使われない場合は削除されるようになります。

## `certificate_project`

API 内の証明書に `restricted` と `projects` プロパティを追加します。
`projects` は証明書がアクセスしたプロジェクト名の一覧を保持します。

## `network_ovn_acl`

OVN ネットワークと OVN NIC に `security.acls` プロパティを追加します。
これにより ネットワークに ACL をかけられるようになります。

## `projects_images_auto_update`

`images.auto_update_cached` と `images.auto_update_interval` 設定キーを追加します。
これらによりプロジェクト内のイメージの自動更新を設定できるようになります。

## `projects_restricted_cluster_target`

プロジェクトに `restricted.cluster.target` 設定キーを追加します。
これによりどのクラスタメンバーにワークロードを配置するかやメンバー間のワークロードを移動する能力を指定する --target オプションをユーザーに使わせないように出来ます。

## `images_default_architecture`

`images.default_architecture` をグローバルの設定キーとプロジェクトごとの設定キーとして追加します。
これはイメージリクエストの一部として明示的に指定しなかった場合にどのアーキテクチャを使用するかを Incus に指定します。

## `network_ovn_acl_defaults`

OVN ネットワークと NIC に `security.acls.default.{in,e}gress.action` と `security.acls.default.{in,e}gress.logged` 設定キーを追加します。
これは削除された ACL の `default.action` と `default.logged` キーの代わりになるものです。

## `gpu_mig`

これは NVIDIA MIG のサポートを追加します。
`mig` GPU type と関連する設定キーを追加します。

## `project_usage`

プロジェクトに現在のリソース割り当ての情報を取得する API エンドポイントを追加します。
API の `GET /1.0/projects/<name>/state` で利用できます。

## `network_bridge_acl`

`bridge` ネットワークに `security.acls` 設定キーを追加し、ネットワーク ACL を適用できるようにします。

さらにマッチしなかったトラフィックに対するデフォルトの振る舞いを指定する `security.acls.default.{in,e}gress.action` と `security.acls.default.{in,e}gress.logged` 設定キーを追加します。

## `warnings`

Incus の警告 API です。

この拡張は次のエンドポイントを含みます（詳細は [Restful API](rest-api.md) 参照）:

* `GET /1.0/warnings`

* `GET /1.0/warnings/<uuid>`
* `PUT /1.0/warnings/<uuid>`
* `DELETE /1.0/warnings/<uuid>`

## `projects_restricted_backups_and_snapshots`

プロジェクトに `restricted.backups` と `restricted.snapshots` 設定キーを追加し、ユーザーがバックアップやスナップショットを作成できないようにします。

## `clustering_join_token`

トラスト・パスワードを使わずに新しいクラスタメンバーを追加する際に使用する参加トークンをリクエストするための `POST /1.0/cluster/members` API エンドポイントを追加します。

## `clustering_description`

クラスタメンバーに編集可能な説明を追加します。

## `server_trusted_proxy`

`core.https_trusted_proxy` のサポートを追加します。 この設定は、Incus が HAProxy スタイルの connection ヘッダーをパースし、そのような（HAProxy などのリバースプロキシサーバーが Incus の前面に存在するような）接続の場合でヘッダーが存在する場合は、プロキシサーバーが（ヘッダーで）提供するリクエストの（実際のクライアントの）ソースアドレスへ（Incus が）ソースアドレスを書き換え（て、Incus の管理するクラスタにリクエストを送出し）ます。（Incus のログにもオリジナルのアドレスを記録します）

## `clustering_update_cert`

クラスタ全体に適用されるクラスタ証明書を更新するための `PUT /1.0/cluster/certificate` エンドポイントを追加します。

## `storage_api_project`

これはプロジェクト間でカスタムストレージボリュームをコピー／移動できるようにします。

## `server_instance_driver_operational`

これは `/1.0` エンドポイントの `driver` の出力をサーバー上で実際にサポートされ利用可能であるドライバーのみを含めるように修正します（Incus に含まれるがサーバー上では利用不可なドライバーも含めるのとは違って）。

## `server_supported_storage_drivers`

これはサーバーの環境情報にサポートされているストレージドライバーの情報を追加します。

## `event_lifecycle_requestor_address`

lifecycle requestor に address のフィールドを追加します。

## `resources_gpu_usb`

リソース API 内の `ResourcesGPUCard`（GPU エントリ）に `USBAddress`（`usb_address`）を追加します。

## `clustering_evacuation`

クラスタメンバーを退避と復元するための `POST /1.0/cluster/members/<name>/state` エンドポイントを追加します。
また設定キー `cluster.evacuate` と `volatile.evacuate.origin` も追加します。
これらはそれぞれ退避の方法（`auto`、`stop`または`migrate`）と移動したインスタンスのオリジンを設定します。

## `network_ovn_nat_address`

これは Incus の `ovn` ネットワークに `ipv4.nat.address` と `ipv6.nat.address` 設定キーを追加します。
これらのキーで OVN 仮想ネットワークからの外向きトラフィックのソースアドレスを制御します。
これらのキーは OVN ネットワークのアップリンクネットワークが `ovn.ingress_mode=routed` という設定を持つ場合にのみ指定可能です。

## `network_bgp`

これは Incus を BGP ルーターとして振る舞わせルートを `bridge` と `ovn` ネットワークに広告するようにします。

以下のグローバル設定が追加されます:

* `core.bgp_address`
* `core.bgp_asn`
* `core.bgp_routerid`

以下のネットワーク設定キーが追加されます（`bridge` と `physical`）:

* `bgp.peers.<name>.address`
* `bgp.peers.<name>.asn`
* `bgp.peers.<name>.password`

`nexthop` 設定キー（NIC type が `bridged` の場合）:

* `bgp.ipv4.nexthop`
* `bgp.ipv6.nexthop`

そして下記の NIC 特有な設定が追加されます（NIC type が `bridged` の場合）:

* `ipv4.routes.external`
* `ipv6.routes.external`

## `network_forward`

これはネットワークアドレスのフォワード機能を追加します。
`bridge` と `ovn` ネットワークで外部 IP アドレスを定義して対応するネットワーク内の内部 IP アドレス（複数指定可能）にフォワード出来ます。

## `custom_volume_refresh`

ボリュームマイグレーションに refresh オプションのサポートを追加します。

## `network_counters_errors_dropped`

これはネットワークカウンターに受信エラー数、送信エラー数とインバウンドとアウトバウンドのドロップしたパケット数を追加します。

## `metrics`

これは Incus にメトリクスを追加します。実行中のインスタンスのメトリクスを OpenMetrics 形式で返します。

この拡張は次のエンドポイントを含みます:

* `GET /1.0/metrics`

## `image_source_project`

`POST /1.0/images` に `project` フィールドを追加し、イメージコピー時にコピー元プロジェクトを設定できるようにします。

## `clustering_config`

クラスタメンバーに `config` プロパティを追加し、キー・バリュー・ペアを設定可能にします。

## `network_peer`

ネットワークピアリングを追加し、 OVN ネットワーク間のトラフィックが OVN サブシステムの外に出ずに通信できるようにします。

## `linux_sysctl`

`linux.sysctl.*` 設定キーを追加し、ユーザーが一コンテナ内の一部のカーネルパラメーターを変更できるようにします。

## `network_dns`

組み込みの DNS サーバーとゾーン API を追加し、 Incus インスタンスに DNS レコードを提供します。

以下のサーバー設定キーが追加されます:

* `core.dns_address`

以下のネットワーク設定キーが追加されます:

* `dns.zone.forward`
* `dns.zone.reverse.ipv4`
* `dns.zone.reverse.ipv6`

以下のプロジェクト設定キーが追加されます:

* `restricted.networks.zones`

DNS ゾーンを管理するために下記の REST API が追加されます:

* `/1.0/network-zones`（GET、POST）
* `/1.0/network-zones/<name>`（GET、PUT、PATCH、DELETE）

## `ovn_nic_acceleration`

OVN NIC に `acceleration` 設定キーを追加し、ハードウェアオフロードを有効にするのに使用できます。
設定値は `none` または `sriov` です。

## `certificate_self_renewal`

これはクライアント自身の信頼証明書の更新のサポートを追加します。

## `instance_project_move`

これは `POST /1.0/instances/NAME` API に `project` フィールドを追加し、インスタンスをプロジェクト間で簡単に移動できるようにします。

## `storage_volume_project_move`

これはストレージボリュームのプロジェクト間での移動のサポートを追加します。

## `cloud_init`

これは以下のキーを含む `project` 設定キー名前空間を追加します:

* `cloud-init.vendor-data`
* `cloud-init.user-data`
* `cloud-init.network-config`

これはまた `devlxd` にインスタンスのデバイスを表示する `/1.0/devices` エンドポイントを追加します。

## `network_dns_nat`

これはネットワークゾーン（DNS）に `network.nat` を設定オプションとして追加します。

デフォルトではすべてのインスタンスの NIC のレコードを生成するという現状の挙動になりますが、
`false` に設定すると外部から到達可能なアドレスのレコードのみを生成するよう Incus に指示します。

## `database_leader`

クラスタリーダーに設定される `database-leader` ロールを追加します。

## `instance_all_projects`

すべてのプロジェクトのインスタンス表示のサポートを追加します。

## `clustering_groups`

クラスタメンバーのグループ化のサポートを追加します。

これは以下の新しいエンドポイントを追加します:

* `/1.0/cluster/groups`（GET、POST）
* `/1.0/cluster/groups/<name>`（GET、POST、PUT、PATCH、DELETE）

 以下のプロジェクトの制限が追加されます:

* `restricted.cluster.groups`

## `ceph_rbd_du`

Ceph ストレージブールに `ceph.rbd.du` という Boolean の設定を追加します。
実行に時間がかかるかもしれない `rbd du` の呼び出しの使用を無効化できます。

## `instance_get_full`

これは `GET /1.0/instances/{name}` に `recursion=1` のモードを追加します。
これは状態、スナップショット、バックアップの構造体を含むすべてのインスタンスの構造体が取得できます。

## `qemu_metrics`

これは `security.agent.metrics` という Boolean 値を追加します。デフォルト値は `true` です。
`false` に設定するとメトリクスや他の状態の取得のために `incus-agent` に接続することはせず、 QEMU からの統計情報に頼ります。

## `gpu_mig_uuid`

NVIDIA `470+` ドライバー (例. `MIG-74c6a31a-fde5-5c61-973b-70e12346c202`) で使用される MIG UUID 形式のサポートを追加します。
`MIG-` の接頭辞は省略できます。

この拡張が古い `mig.gi` と `mig.ci` パラメーターに取って代わります。これらは古いドライバーとの互換性のため残されますが、
同時には設定できません。

## `event_project`

イベントの API にイベントが属するプロジェクトを公開します。

## `clustering_evacuation_live`

`cluster.evacuate` への設定値 `live-migrate` を追加します。
これはクラスタ退避の際にインスタンスのライブマイグレーションを強制します。

## `instance_allow_inconsistent_copy`

`POST /1.0/instances` のインスタンスソースに `allow_inconsistent` フィールドを追加します。
`true` の場合、 `rsync` はコピーからインスタンスを生成するときに `Partial transfer due to vanished source files`（code 24）エラーを無視します。

## `network_state_ovn`

これにより、`/1.0/networks/NAME/state` API に `ovn` セクションが追加されます。これには OVN ネットワークに関連する追加の状態情報が含まれます:

* chassis（シャーシ）

## `storage_volume_api_filtering`

ストレージボリュームの GET リクエストの結果をフィルタリングする機能を追加します。

## `image_restrictions`

この拡張機能は、イメージのプロパティに、イメージの制限やホストの要件を追加します。これらの要件は
インスタンスとホストシステムとの互換性を決定するのに役立ちます。

## `storage_zfs_export`

`zfs.export` を設定することで、プールのアンマウント時に zpool のエクスポートを無効にする機能を導入しました。

## `network_dns_records`

network zones (DNS) API を拡張し、カスタムレコードの作成と管理機能を追加します。

これにより、以下が追加されます:

* `GET /1.0/network-zones/ZONE/records`
* `POST /1.0/network-zones/ZONE/records`
* `GET /1.0/network-zones/ZONE/records/RECORD`
* `PUT /1.0/network-zones/ZONE/records/RECORD`
* `PATCH /1.0/network-zones/ZONE/records/RECORD`
* `DELETE /1.0/network-zones/ZONE/records/RECORD`

## `network_zones_all_projects`

`GET /1.0/network-zones` APIで`all-projects`でパラメータによりすべてプロジェクトのネットワークゾーンを一覧表示できるようにします。

## `storage_zfs_reserve_space`

`quota`/`refquota` に加えて、ZFS プロパティの `reservation`/`refreservation` を設定する機能を追加します。

## `network_acl_log`

ACL ファイアウォールのログを取得するための API `GET /1.0/networks-acls/NAME/log` を追加します。

## `storage_zfs_blocksize`

ZFS ストレージボリュームに新しい `zfs.blocksize` プロパティを導入し、ボリュームのブロックサイズを設定できるようになります。

## `metrics_cpu_seconds`

Incus が使用する CPU 時間をミリ秒ではなく秒単位で出力するように修正されたかどうかを検出するために使用されます。

## `instance_snapshot_never`

`snapshots.schedule`に`@never`オプションを追加し、継承を無効にすることができます。

## `certificate_token`

トラストストアに、トラストパスワードに代わる安全な手段として、トークンベースの証明書を追加します。

これは `POST /1.0/certificates` に `token` フィールドを追加します。

## `instance_nic_routed_neighbor_probe`

これは `routed` NIC が親のネットワークが利用可能かを調べるために IP 近傍探索するのを無効化できるようにします。

`ipv4.neighbor_probe` と `ipv6.neighbor_probe` の NIC 設定を追加します。未指定の場合のデフォルト値は `true` です。

## `event_hub`

これは `event-hub` というクラスタメンバーの役割と `ServerEventMode` 環境フィールドを追加します。

## `agent_nic_config`

これを `true` に設定すると、仮想マシンの起動時に `incus-agent` がインスタンスの NIC デバイスの名前と MTU を変更するための NIC 設定を適用します。

## `projects_restricted_intercept`

`restricted.container.intercept` という設定キーを追加し通常は安全なシステムコールのインターセプションオプションを可能にします。

## `metrics_authentication`

`core.metrics_authentication` というサーバー設定オプションを追加し `/1.0/metrics` のエンドポイントをクライアント認証無しでアクセスすることを可能にします。

## `images_target_project`

コピー元とは異なるプロジェクトにイメージをコピーできるようにします。

## `images_all_projects`

`GET /1.0/images`API に`all-projects` パラメータを追加することで全てのプロジェクトのイメージを一覧表示できるようにします。

## `cluster_migration_inconsistent_copy`

`POST /1.0/instances/<name>` に `allow_inconsistent` フィールドを追加します。 `true` に設定するとクラスタメンバー間で不整合なコピーを許します。

## `cluster_ovn_chassis`

`ovn-chassis` というクラスタロールを追加します。これはクラスタメンバーが OVN シャーシとしてどう振る舞うかを指定できるようにします。

## `container_syscall_intercept_sched_setscheduler`

`security.syscalls.intercept.sched_setscheduler` を追加し、コンテナ内の高度なプロセス優先度管理を可能にします。

## `storage_lvm_thinpool_metadata_size`

`storage.thinpool_metadata_size` により thin pool のメタデータボリュームサイズを指定できるようにします。

指定しない場合のデフォルトは LVM が適切な thin pool のメタデータボリュームサイズを選択します。

## `storage_volume_state_total`

これは `GET /1.0/storage-pools/{name}/volumes/{type}/{volume}/state` API に `total` フィールドを追加します。

## `instance_file_head`

`/1.0/instances/NAME/file` に HEAD を実装します。

## `instances_nic_host_name`

`instances.nic.host_name` サーバー設定キーを追加します。これは `random` か `mac` を指定できます。
指定しない場合のデフォルト値は `random` です。
`random` に設定するとランダムなホストインターフェース名を使用します。
`mac` に設定すると `lxd1122334455` の形式で名前を生成します。

## `image_copy_profile`

イメージをコピーする際にプロファイルの組を修正できるようにします。

## `container_syscall_intercept_sysinfo`

`security.syscalls.intercept.sysinfo` を追加し `sysinfo` システムコールで cgroup ベースのリソース使用状況を追加できるようにします。

## `clustering_evacuation_mode`

退避リクエストに `mode` フィールドを追加します。
これにより従来 `cluster.evacuate` で設定されていた退避モードをオーバーライドできます。

## `resources_pci_vpd`

PCI リソースエントリに VPS 構造体を追加します。
この構造体には完全な製品名と追加の設定キーバリューペアを含むベンダー提供のデータが含まれます。

## `qemu_raw_conf`

生成された qemu.conf の指定したセクションをオーバライドするための `raw.qemu.conf` 設定キーを追加します。

## `storage_cephfs_fscache`

CephFS プール上の `fscache`/`cachefilesd` をサポートするための `cephfs.fscache` 設定オプションを追加します。

## `network_load_balancer`

これはネットワークのロードバランサー機能を追加します。
`ovn` ネットワークで外部 IP アドレス上にポートを定義し、ポートから対応するネットワーク内部の単一または複数の内部 IP にトラフィックをフォワードできます。

## `vsock_api`

これは双方向の `vsock` インターフェースを導入し、 `incus-agent` と Incus サーバーがよりよく連携できるようにします。

## `instance_ready_state`

インスタンスに新しく `Ready` 状態を追加します。これは `/dev/incus` を使って設定できます。

## `network_bgp_holdtime`

特定のピアの BGP ホールドタイムを制御するために `bgp.peers.<name>.holdtime` キーを追加します。

## `storage_volumes_all_projects`

すべてのプロジェクトのストレージボリュームを一覧表示できるようにします。

## `metrics_memory_oom_total`

`/1.0/metrics` API に `incus_memory_OOM_kills_total` メトリックを追加します。
メモリーキラー（`OOM`）が発動された回数を報告します。

## `storage_buckets`

storage bucket API を追加します。ストレージプールのために S3 オブジェクトストレージのバケットの管理をできるようにします。

## `storage_buckets_create_credentials`

これはバケット作成時に管理者の初期クレデンシャルを返すようにストレージバケット API を更新します。

## `metrics_cpu_effective_total`

これは `lxd_cpu_effective_total` メトリックを `/1.0/metrics` API に追加します。
有効な CPU の総数を返します。

## `projects_networks_restricted_access`

プロジェクト内でアクセスできるネットワークを（カンマ区切りリストで）示す `restricted.networks.access` プロジェクト設定キーを追加します。
指定しない場合は、すべてのネットワークがアクセスできます（後述の `restricted.devices.nic` 設定でも許可されている場合）。 

これはまたプロジェクトの `restricted.devices.nic` 設定で制御されるネットワークアクセスにも変更をもたらします:

* `restricted.devices.nic` が `managed` に設定される場合（未設定時のデフォルト）, マネージドネットワークのみがアクセスできます。
* `restricted.devices.nic` が `allow` に設定される場合、すべてのネットワークがアクセスできます（`restricted.networks.access` 設定に依存）。
* `restricted.devices.nic` が `block` に設定される場合、どのネットワークにもアクセスできません。

## `storage_buckets_local`

これは `core.storage_buckets_address` グローバル設定を指定することでローカルストレージプール上のストレージバケットを使用できるようにします。

## `loki`

これはライフサイクルとロギングのイベントを Loki サーバーに送れるようにします。

以下のグローバル設定キーを追加します:

* `loki.api.ca_cert`: イベントを Loki サーバーに送る際に使用する CA 証明書。
* `loki.api.url`: Loki サーバーの URL。
* `loki.auth.username` と `loki.auth.password`: Loki が BASIC 二症を有効にしたリバースプロキシの背後にいる場合に使用。
* `loki.labels`: Loki イベントのラベルに使用されるカンマ区切りリストの値。
* `loki.loglevel`: Loki サーバーに送るイベントの最低のログレベル。
* `loki.types`: Loki サーバーに送られるイベントのタイプ（`lifecycle` および/または `logging`）。

## `acme`

ACME サポートを追加します。これにより [Let's Encrypt](https://letsencrypt.org/) や他の ACME サービスを使って証明書を発行できます。

以下のグローバルの設定キーを追加します:

* `acme.domain`: 証明書を発行するドメイン。
* `acme.email`: ACME サービスのアカウントに使用する email アドレス。
* `acme.ca_url`: ACME サービスのディレクトリー URL、デフォルトは `https://acme-v02.api.letsencrypt.org/directory`。

また以下のエンドポイントを追加します。これは HTTP-01 チャレンジで必要です:

* `/.well-known/acme-challenge/<token>`

## `internal_metrics`

これはメトリクスのリストに内部メトリクスを追加します。
以下を含みます。

* 実行したオペレーションの総数
* アクティブな警告の総数
* デーモンの uptime（秒数）
* Go のメモリー統計
* goroutine の数

## `cluster_join_token_expiry`

クラスタジョイントークンに有効期限を追加します。デフォルトは 3 時間ですが、`cluster.join_token_expiry` 設定キーで変更できます。

## `remote_token_expiry`

リモートの追加ジョイントークンに有効期限を追加します。
`core.remote_token_expiry` 設定キーで変更できます。デフォルトは無期限です。

## `storage_volumes_created_at`

この変更によりストレージボリュームとそのスナップショットの作成日時を保管するようになります。

これは `StorageVolume` と `StorageVolumeSnapshot` API タイプに `CreatedAt` フィールドを追加します。

## `cpu_hotplug`

これは VM に CPU ホットプラグを追加します。
CPU ピンニング使用時はホットプラグは無効になります。CPU ピンニングには NUMA デバイスのホットプラグも必要ですが、これはできないためです。

## `projects_networks_zones`

プロジェクトに`features.networks.zones`の機能を追加します。
これはネットワークゾーンが作成されたときにどのプロジェクトに関連付けされるかを変更します。
以前はネットワークゾーンは`features.networks`の値に従っており、つまりネットワークと同じプロジェクト内に作られていました。

この拡張により`features.networks`から切り離され、複数のプロジェクトがデフォルトプロジェクト（たとえば `features.networks=false`と設定）内のネットワークを共有できるようになり、共有されたネットワークに対してプロジェクト指向の（プロジェクト内のインスタンスのアドレスのみを含むような）「ビュー」を提供するプロジェクト固有の DNS ゾーンを持てるようになりました。

これはまたネットワークの`dns.zone.forward`にも変更を及ぼし、複数のゾーンに共有ネットワークを関連付けるために（プロジェクト毎に最大 1 つの）DNS ゾーン名のカンマ区切りリストを指定できるようになりました。

`dns.zone.reverse.*`設定は変更無しで、引き続き単一の DNS ゾーンのみが設定できます。
しかし結果として生成されるゾーンの内容は、正引きゾーンの 1 つを経由してネットワークを参照するすべてのプロジェクトからのアドレスをカバーするような`PTR`レコードを含むようになりました。

`features.networks=true`の設定を持つ既存のプロジェクトは自動的に`features.networks.zones=true`が設定されますが、新規のプロジェクトは明示的に設定する必要があります。

## `instance_nic_txqueuelength`

NIC デバイスの`txqueuelen`パラメーターを制御する`txqueuelen`キーを追加します。

## `cluster_member_state`

`GET /1.0/cluster/members/<member>/state` API エンドポイントとそれに関連する`ClusterMemberState` API レスポンスタイプを追加します。

## `instances_placement_scriptlet`

Starlark スクリプトレットを Incus に提供し、クラスタ内の新規インスタンスの配置を制御するカスタムロジックを使えるようにします。

Starlark スクリプトレットは新しいグローバル設定オプション`instances.placement.scriptlet`により Incus に提供されます。

## `storage_pool_source_wipe`

ストレージプールに`source.wipe`ブール値を追加し、Incus は要求されたディスクのパーティションヘッダーを消去する必要があることを示します。これにより、既存のファイルシステムがあることによる潜在的な失敗を回避できるようになります。

## `zfs_block_mode`

これにより、ZFS ブロック`filesystem`ボリュームを使用して、ZFS の上に異なるファイルシステムを使用することができるようになります。

これにより、ZFS ストレージプールに以下の新しい設定オプションが追加されます:

* `volume.zfs.block_mode`
* `volume.block.mount_options`
* `volume.block.filesystem`

## `instance_generation_id`

インスタンスの generation ID のサポートが追加されます。VM またはコンテナの generation ID は、インスタンスの時間的な世代が後ろに移動するたびに変更されます。現時点では、 generation ID は VM タイプのインスタンスを通じてのみ公開されています。これにより、VM ゲスト OS は、既に発生した可能性のある状態の複製を回避するために必要な状態を再初期化できます:

* `volatile.uuid.generation`

## `disk_io_cache`

これは、ディスクデバイスに新しい`io.cache`プロパティを導入し、VM のキャッシング動作を上書きするために使用できます。

## `amd_sev`

AMD SEV（Secure Encrypted Virtualization）のサポートを追加します。ゲスト VM のメモリーを暗号化するのに使用できます。

これは SEV 暗号化に以下の新しい設定オプションを追加します:

* `security.sev` :（bool）この VM で SEV を有効化するか
* `security.sev.policy.es` :（bool）この VM で SEV-ES を有効化するか
* `security.sev.session.dh` :（string）ゲストオーナーの`base64`エンコードされた Diffie-Hellman キー
* `security.sev.session.data` :（string）ゲストオーナーの`base64`エンコードされた session blob

## `storage_pool_loop_resize`

これはループファイルをバックエンドとするストレージプールを、プールの`size`設定を変更することでサイズを拡大できるようにします。

## `migration_vm_live`

これは VM で QEMU から QEMU へのライブマイグレーションを（Ceph を含む）共有されたストレージと共有されないストレージプールの両方で実行するサポートを追加します。

これはさらにマイグレーションの `CRIUType` `protobuf` フィールドに `3` という `CRIUType_VM_QEMU` の値を追加します。

## `ovn_nic_nesting`

これは同じインスタンスの `ovn` NIC を別の `ovn` NIC 内にネストするサポートを追加します。
これは OVN ロジカルスイッチポートを VLAN タグを使用する別の OVN NIC の内側にトンネルできるようにします。

この機能は`nested`プロパティを使って親の NIC 名を指定し、`vlan`プロパティでトンネルに使用する VLAN ID を指定することで設定できます。

## `oidc`

OpenID Connect（OIDC）認証のサポートを追加します。

これは以下の設定キーを追加します:

* `oidc.issuer`
* `oidc.client.id`
* `oidc.audience`

## `network_ovn_l3only`

これは `ovn` ネットワークを "layer 3 only" モードに設定する能力を追加します。
このモードは `ipv4.l3only` と `ipv6.l3only` 設定オプションを使うことで IPv4 または IPv6 のレベルで有効化できます。

このモードを有効にすると、ネットワークは以下のように変更されます:

* 仮想ルータの内部ポートアドレスが単一ホストのネットマスクで設定されます（例：IPv4 では /32 あるいは IPv6 では /128）。
* アクティブなインスタンスの NIC アドレスへの静的ルートが仮想ルータに追加されます。
* アクティブでないアドレス向けのパケットがアップリンクのネットワークからエスケープされるのを防ぐために、内部のサブネット全体への破棄ルートが仮想ルータに追加されます。
* 255.255.255.255 のネットマスクがインスタンス設定で使用されるように DHCPv4 サーバーが設定されます。

## `ovn_nic_acceleration_vdpa`

これは `ovn_nic_acceleration` API 拡張をアップデートします。OVN NIC の`acceleration`設定キーが Virtual Data Path Acceleration（VDPA）をサポートするための `vdpa` という値を受け付けられるようになります。

## `cluster_healing`

これにより、オフラインのクラスタメンバーを自動的に退避させるクラスタヒーリングが追加されます。

次の新しい設定キーが追加されます:

* `cluster.healing_threshold`

この設定キーは整数を取り、0（デフォルト）に設定することで無効化できます。設定された場合、その値はオフラインのクラスタメンバーが退避させられる閾値を表します。もし値が`cluster.offline_threshold`よりも低い場合、その値が代わりに使用されます。

オフラインのクラスタメンバーが退避させられると、リモートバックアップされたインスタンスのみが移行されます。ローカルインスタンスは、クラスタメンバーがオフラインになった際にそれらを移行する方法がないため、無視されます。

## `instances_state_total`

この拡張は、インスタンスの状態 API の一部である`InstanceStateDisk`と`InstanceStateMemory`に新しい`total`フィールドを追加します。

## `auth_user`

メイン API エンドポイントに現在のユーザーの詳細情報を追加します。

以下の項目を追加します:

* `auth_user_name`
* `auth_user_method`

## `security_csm`

`CSM`（Compatibility Support Module）の使用を制御する`security.csm`設定キーを追加し、レガシーなオペレーティングシステムを Incus VM 内で稼働できるようにします。

## `instances_rebuild`

この拡張はインスタンスを同じイメージ、別のイメージ、あるいは空のイメージで再構築できるようにします。
`POST /1.0/instances/<name>/rebuild?project=<project>` API エンドポイントと[`incus rebuild`](incus_rebuild.md) CLI コマンドを新しく追加します。

## `numa_cpu_placement`

これは CPU の組を指定の NUMA ノードの組内に配置できるようにします。

以下の設定キーを追加します:

* `limits.cpu.nodes` :（string）（`limits.cpu`の動的な値により選ばれた）CPU を配置する NUMA ノード ID または NUMA ノード ID の範囲のカンマ区切りリスト。

## `custom_volume_iso`

これは ISO イメージをカスタムストレージボリュームとしてインポートできるようにします。

[`incus storage volume import`](incus_storage_volume_import.md)に`--type`フラグを追加します。

## `network_allocations`

これは Incus インストール環境のネットワーク割り当てを一覧表示できるようにします。

[`incus network list-allocations`](incus_network_list-allocations.md)コマンドと`--project <PROJECT> | --all-projects`フラグを使って、使用中の IP アドレス、（インスタンスの）ハードウェアドレス、リソース URI、NAT を使用しているかどうかを各`instance`、`network`、`network forward`、`network load-balancer`について一覧表示できます。

## `zfs_delegate`

これは ZFS ストレージドライバーのボリュームに `zfs.delegate` という Boolean 設定を追加します。
有効にすると対応したシステム（ZFS 2.2 以上が必要）が使用されている場合に、ZFS dataset はコンテナに移譲され、`zfs` コマンドラインツールを使って使用できます。

## `storage_api_remote_volume_snapshot_copy`

これはストレージボリュームスナップショットをリモートとローカル間でコピーできるようにします。

## `operations_get_query_all_projects`

これは `/1.0/operations` と `/1.0/operations?recursion=1` の両方の GET API 呼び出しに `all-projects` クエリパラメーターを追加します。
このパラメーターはプロジェクト名のフィルターをバイパスできるようにします。

## `metadata_configuration`

`GET /1.0/metadata/configuration` API エンドポイントを追加し、生成されたメタデータ設定を JSON 形式で取得できるようにします。
JSON の構造は ```"configs" > `ENTITY` > `ENTITY_SECTION` > "keys" > [<CONFIG_OPTION_0>, <CONFIG_OPTION_1>, ...]``` という形式です。
どの設定オプションが含まれるかは {doc}`configuration options </config-options>` の一覧を参照してください。

## `syslog_socket`

syslog 形式のログメッセージを受信できる syslog ソケットを追加します。
これはイベント API と `incus monitor` で参照でき、Loki にフォワードできます。
この機能を有効にするには `core.syslog_socket` を `true` に設定してください。

## `event_lifecycle_name_and_project`

これは `lifecycle` イベントに `Name` と `Project` フィールドを追加します。

## `instances_nic_limits_priority`

これは NIC 単位に `limits.priority` というオプションを追加します。
これは deprecated になった（cgroup1 のみで使える）`limits.network.priority` インスタンス設定と異なり、cgroup1 と cgroup2 の両方で使えます。

## `disk_initial_volume_configuration`

この API 拡張はインスタンスのルートデバイスに初期ボリューム設定を指定できるよいうにします。
初期ボリューム設定は `initial.` という接頭辞がつけられ、インスタンス初期化時に `--device` フラグを使ってプロファイル経由または直接指定できます。

これらの設定はインスタンスの作成時にのみ適用され、その後の修正時には既存のデバイスに影響しないことに注意してください。

## `operation_wait`

この API 拡張はサーバーに `/1.0/operations/{id}/wait` エンドポイントが存在することを示します。
これはクライアントに `/1.0/events` エンドポイント経由でイベントの完了を待つ代わりにこのエンドポイントを使えることを示します。
`/1.0/events` endpoint.

## `image_restriction_privileged`

この拡張は `requirements.privileged` というイメージの制限を追加します。
`false` に設定すると特権コンテナでイメージが使用できなくなります。

## `cluster_internal_custom_volume_copy`

この拡張は単一の API 呼び出しでクラスタ内でカスタムストレージボリュームをコピーや移動できるようにします。
`POST /1.0/storage-pools/<pool>/custom?target=<target>` を呼ぶとリクエストの `source` 部分で指定されたカスタムボリュームをコピーします。
`POST /1.0/storage-pools/<pool>/custom/<volume>?target=<target>` をリクエストの `source` 部分で指定されたかステムボリュームを `source` から `target` に移動します。

## `disk_io_bus`
これはディスクデバイスに `io.bus` プロパティを追加します。ディスクがアタッチするバスを上書きするのに使えます。

## `storage_cephfs_create_missing`
これは`cephfs`ストレージプールを追加する際に Incus にストレージプールに必要なエンティティが存在しない場合に作成するように指示するのに使える`cephfs.create_missing`、`cephfs.osd_pg_num`、`cephfs.meta_pool`、`cephfs.osd_pool`構成キーを追加します。

## `instance_move_config`

このAPI拡張はインスタンスをプロジェクト間やストレージプール間で移動する際に`--profile`、`--no-profile`、`--device`、`--config`フラグを使えるようにします。

## `ovn_ssl_config`
これは OVN データベースにアクセスするために SSL CA とクライアントキーのペアを提供する新しい設定キーを追加します。
新しい設定キーは`network.ovn.ca_cert`、`network.ovn.client_cert`、`network.ovn.client_key`です。

## `certificate_description`

証明書に`description`フィールドを追加します。

## `disk_io_bus_virtio_blk`

`disk`デバイスに`virtio-blk`という値を追加し、アタッチしたディスクを`virtio-blk`バスに接続できるようにします。

## `loki_config_instance`

Loki イベントの`instance`フィールドをカスタマイズするための`loki.instance`設定キーを追加します。
通常すでに`location`フィールドでカバーされるイベントを送信する個々のシステム名ではなくクラスターの名前を公開するのに使えます。

## `instance_create_start`

`POST /1.0/instances` APIに`start`フィールドを追加します。
`true`に設定すると作成時にインスタンスが自動的に起動します。

このシナリオでは、作成と起動は単一のバックグラウンドのオペレーションの一部になります。

## `clustering_evacuation_stop_options`

`cluster.evacuate` オプションに新しいオプションを追加します:

* `stateful-stop` はインスタンスのレジュームやリストアの際に状態をディスクに保管します。
* `force-stop` はインスタンスのシャットダウンを待たずに即座に停止します。

## `boot_host_shutdown_action`

これは`boot.host_shutdown_action`インスタンス設定キーを追加します。
システムシャットダウンのデフォルトの`stop`の挙動をオーバーライドするのに使えます。
`stop`、`stateful-stop`、`force-stop`の値をサポートします。

## `agent_config_drive`

`agent:config` ディスク `source` を追加します。これでエージェントとその設定を含む ISO を VM ゲストに公開できます。

## `network_state_ovn_lr`

`GET /1.0/networks/NAME/state` API の一部である `NetworkStateOVN` 構造体に `LogicalRouter` フィールドを追加します。

OVN 論理ルーター名を取得するのに使えます。

## `image_template_permissions`

イメージメタデータテンプレートのエントリに `uid`、`gid`、`mode` フィールドを追加します。

## `storage_bucket_backup`

ストレージバケットのバックアップサポートを追加します。

以下の新しいエンドポイント（詳細は[RESTful API](rest-api.md)参照）を含みます:

* `GET /1.0/storage-pools/<pool>/buckets/<bucket>/backups`
* `POST /1.0/storage-pools/<pool>/buckets/<bucket>/backups`

* `GET /1.0/storage-pools/<pool>/buckets/<bucket>/backups/<name>`
* `POST /1.0/storage-pools/<pool>/buckets/<bucket>/backups/<name>`
* `DELETE /1.0/storage-pools/<pool>/buckets/<bucket>/backups/<name>`

* `GET /1.0/storage-pools/<pool>/buckets/<bucket>/backups/<name>/export`

## `storage_lvm_cluster`

`lvmlockd` 経由で LVM の共有された VG を使う `lvmcluster` ストレージドライバを追加します。

これを使うと、複数のサーバーが同じバッキングデバイスを参照する限りは、単一の共有された LVM プールをそれらのサーバーで使用できます。

## `shared_custom_block_volumes`

カスタムブロックボリュームに設定キー `security.shared` を追加します。
未設定か `false` の場合、カスタムブロックボリュームは複数のインスタンスにアタッチできません。
この機能はカスタムブロックボリュームを一度に複数のインスタンスにアタッチした際に発生する可能性のあるデータ消失を防ぐために追加されました。

## `auth_tls_jwt`

TLS クライアント証明書を直接使う代わりに署名された `JSON Web Token` (`JWT`) を使えるようにします。

このシナリオでは、自身の TLS クライアント証明書から `JWT` を導出し、`Authorization` HTTP ヘッダーの `bearer` トークンとして提供します。

`JWT` は証明書のフィンガープリントを `Subject` に設定し、クライアントの秘密鍵で署名される必要があります。

## `oidc_claim`

`oidc.claim` サーバー設定キーを追加します。これでどの OpenID Connect クレームをユーザー名として使うかを指定できます。

## `device_usb_serial`

デバイスタイプ `usb` に `serial` 設定キーを追加します。
この機能が追加されると、同一の `vendorid` と `productid` の複数のデバイスを区別できるようになります。

## `numa_cpu_balanced`

`limits.cpu.nodes` の新しい値として `balanced` を追加します。

`balanced` に設定すると、Incus はインスタンスの起動時に一番暇な NUMA ノードを選び、システム上の NUMA ノードで負荷が分散できるようにします。

## `image_restriction_nesting`

`requirements.nesting` というイメージ制限を追加します。 `true` はイメージがネストなしには実行できないことを意味します。

## `network_integrations`

ネットワーク統合の概念と OVN Interconnection の初期サポートを追加します。

新 API:

* `/1.0/network-integrations` (GET, POST)
* `/1.0/network-integrations/NAME` (GET, PUT, PATCH, DELETE, POST)

それぞれの統合は以下の項目で構成されます:

* name (名前)
* description (説明)
* type (種別、現在は `ovn` のみ)
* configuration (設定)
   * `ovn.northbound_connection` (OVN Interconnection データベースの接続文字列)
   * `ovn.ca_cert` (省略可能、OVN Interconnection データベースの SSL CA 証明書)
   * `ovn.client_cert` (省略可能、OVN Interconnection データベースに接続するための SSL クライアント証明書)
   * `ovn.client_key` (省略可能、OVN Interconnection データベースに接続するための SSL クライアント鍵)
   * `ovn.transit.pattern` (transit スイッチ名を制定するための Pongo2 テンプレート)

これらの統合はいくつかの新しいフィールドを通してネットワークピアにアタッチします:

* `type` (現状の挙動は `local`、 統合は `remote`)
* `target_integration` (統合への参照)

## `instance_memory_swap_bytes`

`limits.memory.swap` を全体のリミットをバイト数で指定できるようにします。

## `network_bridge_external_create`

`bridge.external_interfaces` を `interface/parent/vlan` のシンタックスを使って親のインターフェースを作れるようにします。

## `storage_zfs_vdev`

ストレージの`source`設定を拡張して`mirror`、`raidz1`、`raidz2` ZFS `vdev`のサポートを追加します。

## `container_migration_stateful`

`migration.stateful`設定キーを追加します。

これはBooleanのフラグでコンテナが起動、終了、スナップショットの機能中にステートフルモードであるときは常にtrueに設定されます。

これによりコンテナを別のシステムにコピーする際にCRIUエラーを起きにくくします。

## `profiles_all_projects`

`GET /1.0/profiles` APIで`all-projects`でパラメータによりすべてプロジェクトのプロファイルを一覧表示できるようにします。

## `instances_scriptlet_get_instances`

オプショナルのProjectやLocationフィルターを指定してインスタンススクリプトレットがインスタンスの一覧を取得できるようにします。

## `instances_scriptlet_get_cluster_members`

オプショナルでクラスタグループを指定してインスタンススクリプトレットがクラスタメンバーの一覧を取得できるようにします。

## `instances_scriptlet_get_project`

プロジェクト名を指定してインスタンススクリプトレットがプロジェクトを取得できるようにします。

## `network_acl_stateless`

ネットワークのACL内でステートレスなルールのサポートを追加します。

## `instance_state_started_at`

インスタンスステートAPIに`started_at`タイムスタンプを追加します。

## `networks_all_projects`

`GET /1.0/networks` APIで`all-projects`でパラメータによりすべてプロジェクトのネットワークを一覧表示できるようにします。

## `network_acls_all_projects`

`GET /1.0/network-acls` APIで`all-projects`でパラメータによりすべてプロジェクトのネットワークACLを一覧表示できるようにします。

## `storage_buckets_all_projects`

`GET /1.0/storage-pools/POOL/buckets` APIで`all-projects`でパラメータによりすべてプロジェクトのストレージバケットを一覧表示できるようにします。

## `resources_load`

リソースAPIにLoadセクションを追加します。

## `instance_access`

`GET /1.0/instances/NAME/access`に新しいAPIエンドポイントを追加します。これはインスタンスを操作できる人とその人たちが持つロールを表示します。

## `project_access`

`GET /1.0/projects/NAME/access`に新しいAPIエンドポイントを追加します。これはプロジェクトを操作できる人とその人たちが持つロールを表示します。

## `projects_force_delete`

`DELETE /1.0/projects`が`?force=true`を受け付けるようにします。指定するとプロジェクトをプロジェクト内部に含まれるあらゆるものとまとめて消せるようになります。

## `resources_cpu_flags`

これはリソースAPIにCPUの機能をチェックするためのCPUフラグ／拡張を追加します。

## `disk_io_bus_cache_filesystem`

これはファイルシステムに支えられたディスクで`io.bus`と`io.cache`を使えるようにします。

## `instance_oci`

OCIコンテナを稼働させる初期サポートを追加します。

## `clustering_groups_config`

クラスタリンググループにスタンダードなキーバリューの `config` オプションを追加します。
これによりそれらのグループにいつくかの制限や設定をすることができます。

## `instances_lxcfs_per_instance`

`instances.lxcfs.per_instance` というサーバー設定キーを追加します。
これによりシステム上でグローバルにではなくインスタンスごとにLXCFSを動かすかどうかを制御できます。

## `clustering_groups_vm_cpu_definition`

クラスターグループの設定を通して仮想マシンのCPU定義を制御するためのいくつかの設定オプションを追加します。

追加される設定キーは:

* `instances.vm.cpu.ARCHITECTURE.baseline`
* `instances.vm.cpu.ARCHITECTURE.flag`

## `disk_volume_subpath`

`source=volume/path`のシンタクスを用いてファイルシステムのカスタムボリュームのサブパスにアクセスできるようにします。

## `projects_limits_disk_pool`

プールごとのプロジェクトのディスク制限を追加します。
プロジェクトの制限に`limits.disk.pool.NAME`という設定オプションを追加します。

## `network_ovn_isolated`

OVNネットワークのアップリンクネットワークとして`none`を使えるようにします。
これによりネットワークを隔離できます。

## `qemu_raw_qmp`

これは仮想マシンの起動時のさまざまなステージでQMPコマンドを直接発行する設定オプションを追加します:

* `raw.qemu.qmp.early`
* `raw.qemu.qmp.pre-start`
* `raw.qemu.qmp.post-start`

## `network_load_balancer_health_check`

ロードバランサーのバックエンドにヘルスチェックを実行できるようにします。

以下の新しい設定オプションが追加されます:

* `healthcheck`
* `healthcheck.interval`
* `healthcheck.timeout`
* `healthcheck.failure_count`
* `healthcheck.success_count`

## `oidc_scopes`

`oidc.scopes`サーバー設定キーを追加します。これはアイデンティティープロバイダーからのリクエストにOIDCスコープのカンマ区切りリストを指定できます。

## `network_integrations_peer_name`

`ovn.transit.pattern`で`peerName`をテンプレート変数として使えるようにします。

## `qemu_scriptlet`

起動時のさまざまなステージで`raw.qemu.scriptlet`設定キーを使ってスクリプトを実行できるようにします。

## `instance_auto_restart`

`boot.autorestart`設定キーを追加します。`true`に設定するとインスタンスが異常終了したときに1分間に最大10回まで自動的に再起動されるようになります。

## `storage_lvm_metadatasize`

LVMストレージプールに`lvm.metadata_size`オプションを追加します。LVM物理ボリュームを作成する際のデフォルトのメタデータサイズをオーバーライドできるようにします。

## `ovn_nic_promiscuous`

OVN NICに`security.promiscuous`設定オプションを追加します。

## `ovn_nic_ip_address_none`

OVN NICの`ipv4.address`と`ipv6.address`に`none`という値を追加します。

## `instances_state_os_info`

この拡張はインスタンスの状態APIに`InstanceStateOSInfo 構造体へのポインターを追加します。

## `network_load_balancer_state`

`/1.0/networks/NAME/load-balancers/IP/state` APIエンドポイントを追加します。
これは（設定された際は）ロードバランサーのヘルスチェック情報を返します。

## `instance_nic_macvlan_mode`

これはmacvlanネットワークインタフェースに`mode`キーを追加し、Macvlanモードを設定できるようにします。

## `storage_lvm_cluster_create`

`source`を共有されたブロックデバイスに設定することで新しいLVMクラスタプールを作成できるようにします。

## `network_ovn_external_interfaces`

OVNネットワークに`bridge.external_interfaces`のサポートを追加します。

## `instances_scriptlet_get_instances_count`

これにより、インスタンス・スクリプトレットで、オプショナルなProjectまたはLocationのフィルターを指定して絞ったインスタンス数と起動中のインスタンス数を取得できます。

## `cluster_rebalance`

これはクラスターの負荷を再分散するために自動ライブマイグレーションを追加します。

この一環として、以下の設定オプションが追加されます:

* `cluster.rebalance.batch`
* `cluster.rebalance.cooldown`
* `cluster.rebalance.interval`
* `cluster.rebalance.threshold`

## `custom_volume_refresh_exclude_older_snapshots`

これは最新のターゲット（コピー先）のスナップショットより古いソース（コピー元）のスナップショットを除外するサポートを追加します。

## `storage_initial_owner`

カスタムボリュームの初期オーナーを設定できるようにします。

以下の設定オプションが追加されました:

* `initial.gid`
* `initial.mode`
* `initial.uid`

## `storage_live_migration`

ストレージプール間での仮想マシンのライブマイグレーションのサポートを追加します。

## `instance_console_screenshot`

仮想マシンの現在のVGAコンソールのスナップショットを取れるようにします。

## `image_import_alias`

イメージをアップロードする際にエイリアスを設定するために `X-Incus-aliases` HTTP ヘッダーを追加します。

## `authorization_scriptlet`

Incusクラスターでの認可を管理するために、`authorization.scriptlet`という新しい設定キーでスクリプトレットを定義できるようにします。

## `console_force`

アクティブなセッションが既にあっても、コンソールへの接続を強制するサポートを追加します。
インスタンスコンソールへの接続に`--force`という新しいフラグを追加します。

## `network_ovn_state_addresses`

これはOVNネットワークの状態を保持する構造体にアップリンクで使われるIPv4とIPv6アドレス用のフィールドを追加します。

## `qemu_scriptlet_config`

QEMUスクリプトレットの機能を拡張し、VM起動前にQEMUの設定を変更し、インスタンスからスクリプトレットで情報を渡せるようにします。

## `network_bridge_acl_devices`

デバイスがブリッジネットワークにアタッチされたときにACLを設定できるようにします。

## `instance_debug_memory`

`/1.0/instances/NAME/debug/memory`にメモリダンプAPIを追加します。

## `init_preseed_storage_volumes`

このAPI拡張はpreseed initでストレージボリュームを設定できるようにします。

## `init_preseed_profile_project`

このAPI拡張はpreseed initでプロファイルの定義の一部としてプロジェクトを指定できるようにします。

## `instance_nic_routed_host_address`

ルートを追加するVRFを指定するサポートを追加。

## `instance_smbios11`

新しいカテゴリの設定オプション、`smbios11.XYZ`が追加されました。
これは`SMBIOS Type 11`をサポートするシステム上で`SMBIOS Type 11`を通してキー・バリュー・ペアを渡せるようにします。

## `api_filtering_extended`

これはすべてのAPIのコレクションにAPIのフィルタリングのメカニズムを拡張します。

## `acme_dns01`

証明書生成のためのIncusのACMEサポートに`DNS-01`チャレンジのサポートを追加します。

## `security_iommu`

IOMMUエミュレーションを有効にするかどうか制御するための`security.iommu`設定キーを導入します。
これはLinuxでは`virtio_iommu`、WindowsではエミュレートされたIntel IOMMU経由で実行されます。

## `network_ipv4_dhcp_routes`

ブリッジネットワークとOVNネットワークに`ipv4.dhcp.routes`設定オプションを導入します。
これによりDHCPサーバで広告されるCIDRネットワークとゲートウェイアドレスの組を指定できるようになります。

## `network_state_ovn_ls`

`NetworkStateOVN`構造体に`LogicalSwitch`フィールドを追加します。これは`GET /1.0/networks/NAME/state` APIの一部になります。

これはOVN論理スイッチ名を取得するのに使われます。

## `network_dns_nameservers`

ブリッジネットワークとOVNネットワークに`dns.nameservers`設定オプションを導入します。
これによりDHCPサーバによりルーターアドバイズメント経由で広告されるIPv4とIPv6 DNSサーバアドレスを指定できるようになります。

## `acme_http01_port`

`HTTP-01`のバリデーションに使用する別のHTTPポートを制御するための`acme.http.port`を追加します。

## `network_ovn_ipv4_dhcp_expiry`

OVNネットワークに`ipv4.dhcp.expiry`を導入します。

## `instance_state_cpu_time`

インスタンスの状態APIの`CPU`の下に`allocated_time`フィールドを追加します。

## `network_io_bus`

互換のあるネットワークデバイスに`io.bus`プロパティを追加し、`virtio`（デフォルト）と`usb`から選べるようにします。

## `disk_io_bus_usb`

`disk`デバイスの`io.bus`上に`usb`という値を追加します。

## `storage_driver_linstor`

LINSTORストレージドライバを追加します。

## `instance_oci_entrypoint`

コンテナにOCIエントリポイントを設定するための新しい設定オプションの組を導入します。

* `oci.entrypoint`
* `oci.cwd`
* `oci.uid`
* `oci.gid`

これらはOCIイメージからの値を使って作成時に初期化されます。

## `network_address_set`

これは、APIエンドポイントプリフィクス`/1.0/network-address-sets`の下のAPIにネットワークアドレスセットの概念を追加します。

## `server_logging`

これはサーバーの`logging`設定キーのセットを新たに実装します。これにより複数のロギングターゲットが使えます。
以前の`loki`設定キーはこの一部として移行中です。

## `network_forward_snat`

ネットワークフォワードに`snat`設定オプションを追加します。これによりDNATがマッチするSNATが適用されるようになります。
その結果ターゲットからの新しい接続はネットワークフォワードアドレスから来たように見えるようになります。

OVNではこのために十分柔軟なSNATをサポートしていないため、この機能はブリッジネットワークに限られます。

## `memory_hotplug`

これはVMにメモリーホットプラグを追加し、実行時に再起動せずにメモリを追加できるようにします。

## `instance_nic_routed_host_tables`

これはルーテッドモードを使う`nic`デバイス上でホストルーティングテーブルを設定できるようにします。

## `instance_publish_split`

これは既存のインスタンスから分割した形式のイメージを作成できるようにします。

## `init_preseed_certificates`

このAPI拡張はプレシードの初期化の際に証明書を設定できるようにします。

## `custom_volume_sftp`

これはカスタムストレージボリュームにSFTP APIを追加します。

## `network_ovn_external_nic_address`

これはネットワークフォワード経由でインスタンスにアドレスが利用できる限りにおいてカスタムの外部IPv4またはIPv6アドレスを設定できるようにします。

## `network_physical_gateway_hwaddr`

OVNを使用する際にIPv4とIPv6ゲートウェイのMACアドレスを設定できるようにします。

## `backup_s3_upload`

インスタンスやボリュームバックアップをS3互換エンドポイントに直ちにアップロードできるようにします。

## `snapshot_manual_expiry`

インスタンスとストレージボリュームの両方に`snapshots.expiry.manual`設定キーを追加します。
スケジュールにしたがって作成されたスナップショットではなく、ユーザーが直接作成したスナップショットのデフォルトの有効期限をオーバーライドできます。

## `resources_cpu_address_sizes`

リソースAPIでCPUのアドレスサイズをトラッキングできるようにします。
これの主な利用法はクラスター内で仮想マシン内にホットプラグできる最大のメモリ量を計算することです。
