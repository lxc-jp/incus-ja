(authorization)=
# 認可

Unix ソケットで Incus とやり取りする際、`incus-admin`グループのメンバーは Incus API へのフルアクセスを得ます。
一方、`incus` グループのみのメンバーはそのユーザーに割り当てられた単一のプロジェクトに制限されます。

ネットワーク越しに Incus とやり取りする際は（手順は{ref}`server-expose`参照）、さらに認証しユーザーアクセスを制限できます。
以下の 2 つの認可の方法がサポートされます:

- {ref}`authorization-tls`
- {ref}`authorization-openfga`

(authorization-tls)=
## TLS 認可

Incus はネイティブで {ref}`authentication-trusted-clients` を 1 つまたは複数のプロジェクトに制限することをサポートします。
クライアント証明書が制限される際、クライアントはまたグローバルな設定変更やアクセス可能なプロジェクトの設定（限度や制約）の変更も行えません。

アクセスを制限するには、[`incus config trust edit <fingerprint>`](incus_config_trust_edit.md) を使用します。
`restricted`キーを`true`に設定し、クライアントのアクセスを制限するプロジェクトのリストを指定します。
プロジェクトのリストが空の場合、クライアントはどのプロジェクトへのアクセスも許可されません。

この認可の方法はクライアントが TLS で認証する場合は、他の認可の方法が設定されているかどうかによらず、必ず使用されます。

(authorization-openfga)=
## Open Fine-Grained Authorization (OpenFGA)

Incus は [{abbr}`OpenFGA (Open Fine-Grained Authorization)`](https://openfga.dev) との統合をサポートします。
この認可の方法はきめ細かく設定ができます。
例えば、ユーザーアクセスを単一のインスタンスに制限するのに使えます。

OpenFGA を認可に使うには、あなた自身で OpenFGA サーバーを設定し稼働させる必要があります。
Incus でこの認可の方法を有効にするには、[`openfga.*`](server-options-openfga) サーバー設定オプションを設定する必要があります。
Incus は OpenFGA サーバーに接続し、{ref}`openfga-model` を書き込み、以降の全てのリクエストへの認可をこのサーバーに問い合わせます。

(openfga-model)=
### OpenFGA モデル

OpenFGA では、特定の API リソースへのアクセスはユーザーとそのリソースの関連によって決定されます。
これらの関連は [OpenFGA 認可モデル](https://openfga.dev/docs/concepts#what-is-an-authorization-model)で決まります。
Incus の OpenFGA 認可モデルは API リソースを他のリソースとの関連とユーザーやグループのそのリソースへの関連に基づいて記述します。
いくつかの便利な関連がモデルに組み込まれています:

- `server -> admin`: Incus へのフルアクセス。
- `server -> operator`: サーバー設定、証明書、ストレージプールの編集権限を除いた、Incus へのフルアクセス。
- `server -> viewer`: サーバーレベルの設定を参照できるが変種出来ない。プロジェクトやその中身は参照できない。
- `project -> manager`: 編集権限を含む、単一プロジェクトへのフルアクセス。
- `project -> operator`: 編集権限を除いた、単一プロジェクトへのフルアクセス。
- `project -> viewer`: 単一プロジェクトへの参照権限。
- `instance -> manager`: 編集権限を含む、単一インスタンスへのフルアクセス。
- `instance -> operator`: 編集権限を除いた、単一インスタンスへのフルアクセス。
- `instance -> user`: 単一インスタンスへの参照権限に加えて、 `exec`、`console`、`file` APIの権限。
- `instance -> viewer`: 単一インスタンスへの参照権限。

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

```{literalinclude} ../internal/server/auth/driver_openfga_model.openfga
---
language: none
---
```
