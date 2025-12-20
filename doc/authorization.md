(authorization)=
# 認可

Unix ソケットで Incus とやり取りする際、`incus-admin`グループのメンバーは Incus API へのフルアクセスを得ます。
一方、`incus` グループのみのメンバーはそのユーザーに割り当てられた単一のプロジェクトに制限されます。

ネットワーク越しに Incus とやり取りする際は（手順は{ref}`server-expose`参照）、さらに認証しユーザーアクセスを制限できます。
以下の 3 つの認可の方法がサポートされます:

- {ref}`authorization-tls`
- {ref}`authorization-openfga`
- {ref}`authorization-scriptlet`

(authorization-tls)=
## TLS 認可

Incus はネイティブで {ref}`authentication-trusted-clients` を 1 つまたは複数のプロジェクトに制限することをサポートします。
クライアント証明書が制限される際、クライアントはまたグローバルな設定変更やアクセス可能なプロジェクトの設定（限度や制約）の変更も行えません。

アクセスを制限するには、[`incus config trust edit <fingerprint>`](incus_config_trust_edit.md) を使用します。
`restricted`キーを`true`に設定し、クライアントのアクセスを制限するプロジェクトのリストを指定します。
プロジェクトのリストが空の場合、クライアントはどのプロジェクトへのアクセスも許可されません。

{ref}`OpenFGA authorization <authorization-openfga>`が設定されている場合でも、クライアントが TLS で認証する場合は、この認可の方法が使われます。

(authorization-openfga)=
## Open Fine-Grained Authorization (OpenFGA)

Incus は [{abbr}`OpenFGA (Open Fine-Grained Authorization)`](https://openfga.dev) との統合をサポートします。
この認可の方法はきめ細かく設定ができます。
例えば、ユーザーアクセスを単一のインスタンスに制限するのに使えます。

OpenFGA を認可に使うには、あなた自身で OpenFGA サーバーを設定し稼働させる必要があります。
Incus は OpenFGA サーバーに接続し、{ref}`openfga-model` を書き込み、以降の全てのリクエストへの認可をこのサーバーに問い合わせます。

Incusでこの認可の方法を有効にするには、[`openfga.*`](server-options-openfga)サーバー設定オプションを設定してください。
OpenFGAを有効にするには`openfga.api.url`と`openfga.api.token`の両方を設定する必要があります。`openfga.store.id`はオプショナルです。指定しない場合はIncusが新しいストアーを生成します。

(openfga-model)=
### OpenFGA モデル

OpenFGA では、特定の API リソースへのアクセスはユーザーとそのリソースの関連によって決定されます。
これらの関連は [OpenFGA 認可モデル](https://openfga.dev/docs/concepts#what-is-an-authorization-model)で決まります。
Incus の OpenFGA 認可モデルは API リソースを他のリソースとの関連とユーザーやグループのそのリソースへの関連に基づいて記述します。

```{literalinclude} ../internal/server/auth/driver_openfga_model.openfga
---
language: none
---
```

```{important}
ホストへのルート権限を信頼して与えられないユーザーに対しては以下の関連は許可すべきではありません:

- `server -> admin`
- `server -> operator`
- `server -> can_edit`
- `server -> can_create_storage_pools`
- `server -> can_create_projects`
- `server -> can_create_certificates`
- `certificate -> can_edit`
- `storage_pool -> can_edit`
- `project -> manager`

他の関連は許可しても構いません。
しかし、適切な{ref}`project-restrictions`を適用する必要があります。
```

完全な Incus の OpenFGA の認可モデル `internal/server/auth/driver_openfga_model.openfga` 内で定義されます:

(authorization-scriptlet)=
## 認可スクリプトレット

Incusはきめ細やかな認可を管理するスクリプトレットの定義をサポートします。これにより外部ツールに依存することなく詳細な認可のルールを書くことができます。

認可スクリプトレットを使うためには、`authorization.scriptlet`サーバー設定オプションに`authorize`という関数を実装するスクリプトレットを書きます。この関数は3つの引数を取ります:

- `details`：`Username`（ユーザー名あるいは証明書のフィンガープリント）、`Protocol`（認可のプロトコル）、`IsAllProjectsRequest`（リクエストがすべてのプロジェクトに対してされるかどうか）、`ProjectName`（プロジェクト名）のアトリビュートを持つオブジェクト
- `object`：ユーザーが認可をリクエストする対象のオブジェクト
- `entitlement`：ユーザーが希望する認可レベル

この関数はユーザーが対象のオブジェクトに指定の認可レベルでアクセスできるかできないかを示すBooleanの値を返す必要があります。

さらに、アクセスAPIを使ってユーザーが一覧表示できるように、2つのオプショナルな関数を定義できます:

- `get_instance_access`：2つの引数（`project_name`と`instance_name`）を取り、指定のインスタンスにアクセスできるユーザー一覧を返す
- `get_project_access`：1つの引数（`project_name`）を取り、指定のプロジェクトにアクセスできるユーザー一覧を返す
