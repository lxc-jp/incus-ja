(exp-clustering)=
# クラスタリングについて

全体のワークロードを複数のサーバーに分散するため、 Incus はクラスタリングモードで動かせます。
このシナリオでは、クラスタメンバーとそのインスタンスの設定を保持する同じ分散データベースを任意の台数の Incus サーバーで共有します。
Incus クラスタは [`incus`](incus.md) クライアントまたは REST API を使って管理できます。

(clustering-members)=
## クラスタメンバー

Incus クラスタは 1 台のブートストラップサーバーと少なくともさらに 2 台のクラスタメンバーから構成されます。
クラスタは状態を [分散データベース](../database.md) に保管します。これは Raft アルゴリズムを使用して複製される[Cowsql](https://github.com/cowsql/cowsql/) データベースです。

2 台のメンバーだけでもクラスタを作成することは出来なくはないですが、少なくとも 3 台のクラスタメンバーを強く推奨します。
このセットアップでは、クラスタは少なくとも 1 台のメンバーの消失に耐えることができ、分散状態の過半数を確立できます。

クラスタを作成する際、 Cowsql データベースは 3 番目のメンバーがクラスタにジョインするまではブートストラップサーバー上でのみ稼働します。
そして 2 番目と 3 番目のサーバーはデータベースの複製を受信します。

詳細は {ref}`cluster-form` を参照してください。

(clustering-member-roles)=
### メンバーロール

3 台のメンバーのクラスタでは、すべてのメンバーがクラスタの状態を保管する分散データベースを複製します。
クラスタのメンバーがさらに増えると、一部のメンバーだけがデータベースを複製します。
残りのメンバーはデータベースへアクセスしますが、複製はしません。

任意の時点で、選出されたリーダーが 1 つ存在し、他のメンバーの健康状態をモニターします。

データベースを複製する各メンバーは *voter* か *stand-by* のロールを持ちます。
クラスタリーダーがオフラインになると voter の 1 つが新しいリーダーに選出されます。
voter のメンバーがオフラインになると stand-by メンバーが自動的に voter に昇格します。
データベース (そしてクラスタ) は voter の過半数がオンラインである限り利用可能です。

以下のロールが Incus クラスタメンバーに割り当て可能です。
自動のロールは Incus 自身によって割り当てられユーザーによる変更は出来ません。

| ロール                  | 自動     | 説明 |
| :---                  | :--------     | :---------- |
| `database`            | yes           | 分散データベースの voter メンバー |
| `database-leader`     | yes           | 分散データベースの現在のリーダー |
| `database-standby`    | yes           | 分散データベースの stand-by（voter ではない）メンバー |
| `event-hub`           | no            | 内部 Incus イベントへの交換ポイント（hub）（最低 2 つは必要）|
| `ovn-chassis`         | no            | OVN ネットワークのアップリンクゲートウェイの候補 |


voter メンバーのデフォルトの数（{config:option}`server-cluster:cluster.max_voters`）は 3 です。
stand-by メンバーのデフォルトの数（{config:option}`server-cluster:cluster.max_standby`）は 2 です。
この設定では、クラスタを稼働したまま一度に最大で 1 つの voter メンバーの電源を切ることができます。

詳細は {ref}`cluster-manage` を参照してください。

(clustering-offline-members)=
#### オフラインメンバーと障害耐性

クラスタメンバーがダウンして設定されたオフラインの閾値を超えると、ステータスはオフラインと記録されます。
この場合、このメンバーに対する操作はできなくなり、すべてのメンバーの状態変更を必要とする操作もできなくなります。

オフラインのメンバーがオンラインに戻るとすぐに操作が再びできるようになります。

オフラインになったメンバーがリーダーそのものだった場合、他のメンバーは新しいリーダーを選出します。

サーバーを再びオンラインに復旧できないあるいはしたくない場合、[クラスタからメンバーを削除](cluster-manage-delete-members) できます。

応答しないメンバーをオフラインと判断する秒数は[`cluster.offline_threshold`](server-options-cluster)設定で調整できます。
デフォルト値は 20 秒です。
最小値は 10 秒です。

オフラインのメンバーからインスタンスを自動的に{ref}`退避 <cluster-evacuate>`するには、{config:option}`server-cluster:cluster.offline_threshold`設定をゼロでない値に設定してください。

詳細は{ref}`cluster-recover`を参照してください。

#### Failure domains

オフラインになったメンバーにロールを割り当てる際に、どのクラスタメンバーを優先するかを指示するために failure domain を使用できます。
たとえば、現在データベースロールを持つクラスタメンバーがシャットダウンした場合、 Incus はデータベースロールを同じ failure domain 内の別のクラスタメンバーがあればそれに割り当てようとします。

クラスタメンバーの failure domain を更新するには、[`incus cluster edit <member>`](incus_cluster_edit.md) コマンドを使って `failure_domain` プロパティを `default` から他の文字列に変更します。

(clustering-member-config)=
### メンバー設定

Incus クラスタメンバーは一般的に同一のシステムと想定されています。
それはクラスタにジョインするすべての Incus サーバーはブートストラップサーバーとストレージプールとネットワークについて同一の設定を持つ必要があるということです。

少し異なるディスクの順序やネットワークインターフェースの名前付けのようなことに対応するため、ストレージとネットワークに関連してメンバー固有のいくつかの設定が例外的に用意されています。

クラスタ内にそのような設定が存在する場合、追加するサーバーにはそれらの設定に対する値を提供する必要があります。
たいていの場合、これはインタラクティブな `incus admin init` コマンドで実行され、ユーザーにストレージやネットワークに関連する設定の値の入力を求めます。

通常これらの設定には以下のものが含まれます:

- ストレージプールのソースデバイスとサイズ
- ZFS プール、 LVM thin pool、または LVM ボリュームグループの名前
- ブリッジネットワークの外部インターフェースと BGP の next-hop
- 管理された `physical` または `macvlan` ネットワークの親のネットワークデバイス名

詳細は {ref}`cluster-config-storage` と {ref}`cluster-config-networks` を参照してください。

事前に質問を調べたい（スクリプトでの自動化に有用）場合、 `/1.0/cluster` API エンドポイントをクエリしてください。
これは `incus query /1.0/cluster` あるいは他の API クライアントを使って実行できます。

## イメージ

デフォルトでは、 Incus はデータベースメンバーと同じ数のクラスタメンバーにイメージを複製します。
通常これはクラスタ内で最大 3 つのコピーを持つことを意味します。

障害耐性とイメージがローカルで利用できる確率を改善するためこの数を増やすことができます。
そのためには、{config:option}`server-cluster:cluster.images_minimal_replica` 設定を変更してください。
すべてのクラスタメンバーにイメージをコピーするには`-1`という特別な値を使用できます。

(cluster-groups)=
## クラスタグループ

Incus のクラスタではクラスタグループにメンバーを追加できます。
これらのクラスタグループは、すべての利用可能なメンバーのサブセットに属するクラスタメンバー上で、インスタンスを起動するのに使用できます。
たとえば、GPU を持つすべてのメンバーからなるクラスタメンバーを作って、GPU が必要なすべてのインスタンスをこのクラスタグループ上で起動できます。

デフォルトでは、すべてのクラスタメンバーは `default` グループに属します。

詳細は {ref}`howto-cluster-groups` と {ref}`cluster-target-instance` を参照してください。

(clustering-instance-placement)=
## インスタンスの自動配置

クラスタのセットアップでは各インスタンスはクラスタメンバーの 1 つの上で稼働します。
インスタンスを起動する際、特定のクラスタメンバー、クラスタグループをターゲットにするか、あるいは Incus に自動的にどれかのクラスタメンバーに割り当てさせることもできます。

デフォルトでは、自動的な割り当てはインスタンス数が一番少ないクラスタメンバーを選択します。
複数のメンバーが同じインスタンス数の場合は、それらの 1 つがランダムで選ばれます。

しかし、この挙動を {config:option}`cluster-cluster:scheduler.instance` 設定で制御することもできます:

- クラスタメンバーの `scheduler.instance` が `all` に設定されると、以下の条件でこのクラスタメンバーが選ばれます:

  - インスタンスが `--target` を指定せずに作成され、かつクラスタメンバーのインスタンス数が最小である。
   - インスタンスがこのクラスタメンバー上で稼働するようにターゲットされた。
   - インスタンスがこのクラスタメンバーが所属するクラスタグループのメンバー上で稼働するようにターゲットされ、かつクラスタメンバーがそのクラスタグループの他のメンバーと比べてインスタンス数が最小である。

- クラスタメンバーの `scheduler.instance` が `manual` に設定されると、以下の条件でこのクラスタメンバーが選ばれる:

   - インスタンスがこのクラスタメンバー上で稼働するようにターゲットされた。

- クラスタメンバーの `scheduler.instance` が `group` に設定されると、以下の条件でこのクラスタメンバーが選ばれる:

   - インスタンスがこのクラスタメンバー上で稼働するようにターゲットされた。
   - インスタンスがこのクラスタメンバーが所属するクラスタグループのメンバー上で稼働するようにターゲットされ、かつクラスタメンバーがそのクラスタグループの他のメンバーと比べてインスタンス数が最小である。

(clustering-instance-placement-scriptlet)=
### インスタンス配置スクリプトレット

Incus では埋め込まれたスクリプト(スクリプトレット)を使って自動的なインスタンス配置を制御するカスタムロジックを使用できます。
この方法は、組み込みのインスタンス配置機能よりも柔軟性が高いです。

インスタンス配置スクリプトレットは[Starlark言語](https://github.com/bazelbuild/starlark) (Python のサブセット)で記述する必要があります。
スクリプトレットは、Incus がインスタンスをどこに配置するかを知る必要があるたびに呼び出されます。
スクリプトレットは、配置されるインスタンスに関する情報と、インスタンスをホストできる候補のクラスタメンバーに関する情報を受け取ります。
スクリプトレットからクラスタメンバー候補の状態と利用可能なハードウェアリソースについての情報を要求することもできます。

インスタンス配置スクリプトレットは`instance_placement`関数を以下のシグネチャで実装する必要があります:

   `instance_placement(request, candidate_members)`:

- `request`は、[`scriptlet.InstancePlacement`](https://pkg.go.dev/github.com/lxc/incus/shared/api/scriptlet/#InstancePlacement) の展開された表現を含むオブジェクトである。このリクエストには、`project`および`reason`フィールドが含まれています。`reason`は、`new`、`evacuation`、または`relocation`のいずれかである。
- `candidate_members`は、[`api.ClusterMember`](https://pkg.go.dev/github.com/lxc/incus/shared/api#ClusterMember) エントリを表すクラスタメンバーオブジェクトの`list`である。

たとえば:

```python
def instance_placement(request, candidate_members):
    # 情報ログ出力の例。これは Incus のログに出力されます。
    log_info("instance placement started: ", request)

    # インスタンスのリクエストに基づいてロジックを適用する例。
    if request.name == "foo":
        # エラーログ出力の例。これは Incus のログに出力されます。
        log_error("Invalid name supplied: ", request.name)

        fail("Invalid name") # エラーで終了してインスタンス配置を拒否します。

    # 提供された第1候補のサーバーにインスタンスを配置する。
    set_target(candidate_members[0].server_name)

    return # インスタンス配置を進めるために空を返す。
```

スクリプトレットは Incus に適用するためには`instances.placement.scriptlet`グローバル設定に設定する必要があります。

たとえばスクリプトレットが`instance_placement.star`というファイルに保存されている場合、Incus には以下のように適用できます:

    cat instance_placement.star | incus config set instances.placement.scriptlet=-

Incus に現在適用されているスクリプトレットを見るには`incus config get instances.placement.scriptlet`コマンドを使用してください。

スクリプトレットでは（Starlark で提供される関数に加えて）以下の関数が利用できます:

- `log_info(*messages)`: `info`レベルで Incus のログにログエントリを追加する。`messages`は 1 つ以上のメッセージの引数。
- `log_warn(*messages)`: `warn`レベルで Incus のログにログエントリを追加する。`messages`は 1 つ以上のメッセージの引数。
- `log_error(*messages)`: `error`レベルで Incus のログにログエントリを追加する。`messages`は 1 つ以上のメッセージの引数。
- `set_target(member_name)`: インスタンスが作成されるべきクラスタメンバーを設定する。`member_name`はインスタンスが作成されるべきクラスタメンバーの名前。この関数が呼ばれなければ、Incus は組み込みのインスタンス配置ロジックを使用する。
- `get_cluster_member_resources(member_name)`: クラスタメンバーのリソースについての情報を取得する。[`api.Resources`](https://pkg.go.dev/github.com/lxc/incus/shared/api#Resources)の形式でリソースについての情報を含むオブジェクトを返す。`member_name`はリソース情報を取得する対象のクラスタメンバーの名前。
- `get_cluster_member_state(member_name)`: クラスタメンバーの状態を取得する。[`api.ClusterMemberState`](https://pkg.go.dev/github.com/lxc/incus/shared/api#ClusterMemberState)の形式でクラスタメンバーの状態を含むオブジェクトを返す。`member_name`は状態を取得する対象のクラスタメンバーの名前。
- `get_instance_resources()`: インスタンスが必要とするリソースについての情報を取得する。[`scriptlet.InstanceResources`](https://pkg.go.dev/github.com/lxc/incus/shared/api/scriptlet/#InstanceResources)の形式でリソース情報を含むオブジェクトを返す。
- `get_instances(location, project)`: プロジェクトやロケーションフィルターに基づいてインスタンスの一覧を取得する。[`[]api.Instance`](https://pkg.go.dev/github.com/lxc/incus/shared/api#Instance)の形式でインスタンスの一覧を返す。
- `get_instances_count(location, project, pending)`: プロジェクトやロケーションのフィルターに基づくインスタンス数。この数にはまだデータベースレコードが存在しない現在作成中のインスタンスを含む場合があります。
- `get_cluster_members(group)`: クラスタグループに基づいてクラスタメンバーの一覧を取得する。[`[]api.ClusterMember`](https://pkg.go.dev/github.com/lxc/incus/shared/api#ClusterMember)の形式でクラスタメンバーの一覧を返す。
- `get_project(name)`: プロジェクト名に基づいてプロジェクトオブジェクトを取得する。[`api.Project`](https://pkg.go.dev/github.com/lxc/incus/shared/api#Project)の形式でプロジェクトオブジェクトを返す。

```{note}
オブジェクト内のフィールド名は対応する Go の型の JSON フィールド名と同じです。
```
