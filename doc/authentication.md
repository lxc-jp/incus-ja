(authentication)=
# リモートAPI認証

Incus デーモンとのリモート通信は、HTTPS 上の JSON を使って行われます。

リモート API にアクセスするためには、クライアントは Incus サーバーとの間で認証を行う必要があります。
以下の認証方法がサポートされています:

- {ref}`authentication-tls-certs`
- {ref}`authentication-openid`

(authentication-tls-certs)=
## TLSクライアント証明書

認証に{abbr}`TLS (Transport Layer Security)`クライアント証明書を使用する場合、クライアントとサーバーの両方が最初に起動したときにキーペアを生成します。
サーバーはそのキーペアを Incus ソケットへのすべての HTTPS 接続に使用します。
クライアントは、その証明書をクライアント証明書として、あらゆるクライアント・サーバー間の通信に使用します。

証明書を再生成させるには、単に古いものを削除します。
次の接続時には、新しい証明書が生成されます。

### 通信プロトコル

サポートしているプロトコルは TLS 1.3 以上である必要があります。

`INCUS_INSECURE_TLS`環境変数をクライアントとサーバーの両方で設定することにより Incus が TLS 1.2 を受け入れるようにすることはできます。
しかし、これはサポートされる構成ではなく、時代遅れな企業プロキシを使うために強制される場合にのみ使用すべきです。

すべての通信には完全な前方秘匿を使用し、暗号は強力な楕円曲線（ECDHE-RSA や ECDHE-ECDSA など）に限定してください。

生成される鍵は最低でも 4096 ビットの RSA、できれば 384 ビットの ECDSA が望ましいです。
署名を使用する場合は、SHA-2 署名のみを信頼すべきです。

我々はクライアントとサーバーの両方を管理しているので、壊れたプロトコルや暗号の下位互換をサポートする理由はありません。

(authentication-trusted-clients)=
### 信頼できるTLSクライアント

Incus サーバーが信頼する TLS 証明書のリストは、 [`incus config trust list`](incus_config_trust_list.md) で取得できます。

信頼できるクライアントは以下のいずれかの方法で追加できます:

- {ref}`authentication-add-certs`
- {ref}`authentication-token`

サーバーとの認証を行うワークフローは、SSH の場合と同様で、未知のサーバーへの初回接続時にプロンプトが表示されます:

1. ユーザーが [`incus remote add`](incus_remote_add.md) でサーバーを追加すると、HTTPS でサーバーに接続され、その証明書がダウンロードされ、フィンガープリントがユーザーに表示されます。
1. ユーザーは、これが本当にサーバーのフィンガープリントであることを確認するよう求められます。これは、サーバーに接続して手動で確認するか、サーバーにアクセスできる人に info コマンドを実行してフィンガープリントを比較してもらうことで確認できます。
1. サーバーはクライアントの認証を試みます:

   - クライアント証明書がサーバーのトラストストアにある場合は、接続が許可されます。
   - クライアント証明書がサーバーのトラストストアにない場合、サーバーはユーザーにトークンの入力を求めます。
     提供されたトークンが一致した場合、クライアント証明書はサーバーのトラストストアに追加され、接続が許可されます。
     そうでない場合は、接続が拒否されます。

TLS クライアントから Incus へのアクセスを{ref}`authorization-tls`で制限できます。
クライアントへの信頼を取り消すには、[`incus config trust remove <fingerprint>`](incus_config_trust_remove.md) でそのクライアント証明書をサーバーから削除します。

(authentication-tls-jwt)=
#### TLS 認証の実行に `JSON Web Token` (`JWT`) を使う

認証にクライアントの TLS 証明書を直接使う代わりに、Incus はユーザーが `bearer` トークンを導出し、それを  HTTP `Authorization` ヘッダー経由で使用する方法もサポートします。

このためには、ユーザーはクライアント証明書の完全なフィンガープリントを設定した `Subject` フィールドを持つ署名した `JWT` を生成する必要があります。それは有効な `NotBefore` と `NotAfter` フィールドを持ち、クライアントの証明書の秘密鍵で署名する必要があります。

トークンは現在のクライアント証明書から[`incus remote get-client-token`](incus_remote_get-client-token.md)で生成できます。

WebSocket接続にHTTPヘッダーを設定できないブラウザーベースのクライアントでは代わりに`Authorization`ヘッダーにトークンを設定して`/tls/login`への`POST`リクエストを通じてトークンを一度送れます。成功すると、Incusはトークンを`tls_jwt`という名前の`Secure`で`HttpOnly`なクッキーに保管します。このクッキーはトークンが期限切れになるまで後続のリクエストに付与されます。`/tls/logout`へのリクエストでクッキーをクリアします。

(authentication-add-certs)=
#### 信頼できる証明書をサーバーに追加する

信頼できるクライアントを追加するには、そのクライアント証明書をサーバーのトラストストアに直接追加するのが望ましい方法です。
これを行うには、クライアント証明書をサーバーにコピーし、[`incus config trust add-certificate <file>`](incus_config_trust_add-certificate.md) で登録します。

(authentication-ca-certs)=
#### 認証局を信頼する

すべてのクライアント証明書を追加する代わりに、CA証明書をトラストストアーに追加できます。
すると、そのCAに発行されたクライアント証明書はCAのエントリー自身が提供されたかのように信頼されます。
クライアントはCAエントリーのフィンガープリントで識別され、タイプ、制限、プロジェクトを継承します。
クライアントは自身の証明書と信頼されたCAの間の中間証明書を提示する必要があり、証明書に鍵拡張使用法が含まれる場合その用途に`clientAuth`が含まれていなければなりません。

これにより更新のたびにトラストストアーを更新することなく、内部PKIで発行した短期間のクライアント証明書を使えるようになります。
それらの証明書には失効リストは参照されませんので、発行した証明書を失効するにはCAのエントリーをトラストストアーから削除思案す。
トラストストアー内に直接存在する証明書が対応するCAエントリーよりも常に優先されます。

(authentication-token)=
#### トークンを使ったクライアント証明書の追加

トークンを使って新しいクライアントを追加することもできます。
トークンは調整可能な時間（{config:option}`server-core:core.remote_token_expiry`）を過ぎるか一度使用すると無効になります。

この方法を使用するには、クライアント名の入力を促す [`incus config trust add`](incus_config_trust_add.md) を呼び出して、各クライアント用のトークンを生成します。
その後、クライアントは、生成されたトークンをプロンプトが表示されたときに入力することで、自分の証明書をサーバーのトラストストアに追加することができます。

クライアント側から新しい信頼関係を確立できるようにするには、サーバーにトラストパスワード({config:option}`server-core:core.remote_token_expiry`)を設定する必要があります。クライアントは、プロンプト時にトラストパスワードを入力することで、自分の証明書をサーバーのトラストストアに追加することができます。

<!-- Include start NAT authentication -->

```{note}
Incus サーバーが NAT の後ろ側にいる場合、クライアント用のリモートを追加する際には外部のパブリックアドレスを指定する必要があります:

    incus remote add <name> <IP_address>

サーバーでトークンを生成する際、 Incus はクライアントがサーバーにアクセスするために使える IP アドレスのリストを含めます。
しかし、サーバーが NAT の後ろ側にいる場合、これらのアドレスはクライアントが接続できないローカルアドレスの場合があります。
その場合、手動で外部アドレスを指定する必要があります。
```

<!-- Include end NAT authentication -->

あるいは、クライアントはリモートの追加時にトークンを直接提供することもできます: [`incus remote add <name> <token>`](incus_remote_add.md).

### PKI システムの使用

{abbr}`PKI (Public key infrastructure)`の設定では、システム管理者が中央の PKI を管理し、すべての Incus クライアント用のクライアント証明書とすべての Incus デーモン用のサーバー証明書を発行します。

PKI モードを有効にするには、以下の手順を実行します:

1. すべてのマシンに{abbr}`CA（Certificate authority、認証局）`の証明書を追加します:

   - `client.ca`ファイルをクライアントの設定ディレクトリー（`~/.config/incus`）に置く。
   - `server.ca`ファイルをサーバーの設定ディレクトリー（`/var/lib/incus`）に置く。
1. CA から発行された証明書をクライアントとサーバーに配置し、自動生成された証明書を置き換える。
1. サーバーを再起動します。

このモードでは、Incus デーモンへの接続はすべて、事前に発行された CA 証明書を使って行われます。

もしサーバー証明書が CA によって署名されていなければ、接続は単に通常の認証メカニズムを通過します。
サーバー証明書が有効で CA によって署名されていれば、ユーザーに証明書を求めるプロンプトを出さずに接続を続行します。

生成された証明書は自動的には信頼されないことに注意してください。そのため、{ref}`authentication-trusted-clients`で説明している方法のいずれかで、サーバーに追加する必要があります。
このモードでは、CAで発行したクライアント証明書だけがトラストストアーに追加できます。

あるいは、{config:option}`server-core:core.trust_ca_certificates`を`true`に設定することで、CAで発行された任意の有効なクライアント証明書をサーバーが自動的に信頼することもできます。
これらの証明書はトラストストアーに追加される必要はなく、サーバーへのフルアクセスを得ます。ただし、制限付き証明書としてトラストストアーに追加されている場合は除きます。
サーバーの設定ディレクトリー内の`ca.crl`ファイルにリストされている証明書は拒否されます。

```{caution}
Incus専用のCAだけを`core.trust_ca_certificates`で有効にしてください。そのCAで発行された証明書を持つ誰でもがサーバーを管理できるからです。
```

### ローカルキーの暗号化

暗号化された秘密鍵で新しいクライアント証明書を生成するには、以下のコマンドを実行します：

```sh
incus remote generate-certificate --encrypt
```

証明書か秘密鍵が既に存在する場合は生成は失敗します。
この制限はサーバーが既に信頼している資格情報を間違って上書きするのを防ぎます。
`--encrypt`オプションは新しく生成された鍵だけを暗号化します。

コマンドを実行するとパスフレーズと確認のプロンプトが出ます。
秘密鍵はディスクに書かれる前に暗号化されます。
{ref}`クライアント証明書と秘密鍵の出力ファイル <remote-client-certificate>`を選ぶには`--tls-cert <certificate-path>`と`--tls-key <key-path>`の両方を指定してください。

To encrypt an existing client key while preserving the client identity, use:

`incus` クライアントは暗号化されたクライアントキーもサポートします。上記の方法で生成された鍵は以下のコマンドでパスワードを使って暗号化できます:

```
ssh-keygen -p -o -f .config/incus/client.key
```

```{note}
[`keepalive` mode](remote-keepalive) を有効にしないと、Incus を呼び出すたびにプロンプトが表示され煩わしいかもしれません:

    $ incus list remote-host:
    Password for client.key:
    ┌──────┬───────┬──────┬──────┬──────┬───────────┐
    │ NAME │ STATE │ IPV4 │ IPV6 │ TYPE │ SNAPSHOTS │
    └──────┴───────┴──────┴──────┴──────┴───────────┘
```

```{note}
`incus` のコマンドラインは暗号化されたキーをサポートしますが、[Ansible's connection plugin](https://docs.ansible.com/ansible/latest/collections/community/general/incus_connection.html) のようなツールはサポートしません。
```

(authentication-openid)=
## OpenID Connect認証

Incus は[OpenID Connect](https://openid.net/connect/)を使用して、{abbr}`OIDC (OpenID Connect)` アイデンティティ・プロバイダーを通じてユーザーを認証することをサポートしています。

```{note}
OpenID Connect を通じた認証がサポートされていますが、まだユーザーロールの取り扱いはありません。
設定された OIDC アイデンティティ・プロバイダーを通じて認証するすべてのユーザーは、Incus へのフルアクセスを得ます。
```

Incus を OIDC 認証を使用するように設定するには、[`oidc.*`](server-options-oidc)サーバー設定オプションを設定します。
あなたの OIDC プロバイダーは[Device Authorization Grant](https://oauth.net/2/device-flow/)タイプを有効にするように設定する必要があります。

OIDC 認証で設定された Incus サーバーを指すリモートを追加するには、[`incus remote add <remote_name> <remote_address>`](incus_remote_add.md) を実行します。
その後、ウェブブラウザで認証を求められ、Incus が使用するデバイスコードを確認する必要があります。
Incus クライアントはその後、アクセストークンとリフレッシュトークンを取得し保存し、それらを Incus とのすべてのやりとりに使用します。

Incusはプロバイダーの署名鍵を使ってJWTアクセストークンをオフラインで検証します。
不透明なアクセストークンを発行するプロバイダーもサポートされます。この場合はIncusは代わりにIDトークンを検証し設定されたクライアントIDに発行されたかを確認します。
ですのでIncusで使用されるクライアントIDはプロバイダーに登録された他のどのアプリケーションとも共用してはいけません。
`oidc.audience`を設定するとこのフォールバックは無効になります。オーディエンス（audience）が不透明なアクセストークンでは検証できないからです。

Incusは`incus.allowed_subnets`（文字列のリスト）のカスタムOICDクレームをサポートします。クレームが設定されている場合、クレーム内の一覧で指定されているCIDRサブネットの1つに含まれるIPアドレスからの接続のみを許可します。

```{important}
設定済みの OIDC アイデンティティ・プロバイダーで認証されたユーザーは Incus へのフルアクセスを得ます。
ユーザーアクセスを制限するには、{ref}`authorization`も設定する必要があります。
現状では、OIDC と互換のある唯一の認可の方法は{ref}`authorization-openfga`です。
```

(authentication-server-certificate)=
## TLS サーバー証明書

Incus は {abbr}`ACME (Automatic Certificate Management Environment)` サービス（たとえば、[Let's Encrypt](https://letsencrypt.org/)）を使ったサーバー証明書の発行をサポートします。

この機能を有効にするには、[関連するサーバー設定オプション](server-options-acme)を設定してください。

Incusは`HTTP-01`と`DNS-01`チャレンジの両方をサポートしています。2つの間で設定オプションの組は異なります。

`DNS-01`では、関連する{config:option}`server-acme:acme.provider`と{config:option}`server-acme:acme.provider.environment`の値は[`lego`のドキュメント](https://go-acme.github.io/lego/dns/index.html)に直接書かれています。legoはIncusが使用しているACMEクライアントです。

`HTTP-01`では、Incusは`lego`に一時的に`80`番ポートをリッスンさせ、HTTPチャレンジが通るようにします。
クラスターでは、チャレンジはリーダー上で実行され、他のメンバーは同じポートで一時的にリッスンはしますがチャレンジリクエストはリーダーにフォワードします。
これにより証明書はすべてのクラスターメンバーの名前をカバーできます。
もしIncusサーバーがリバースプロキシの後ろにいる場合は、リバースプロキシがHTTPのトラフィックをHTTPSにリダイレクトするようにする必要があります。

External Account Binding （EAB）が必要なACMEサービスでは{config:option}`server-acme:acme.eab.kid`と{config:option}`server-acme:acme.eab.hmac`の両方をACMEサービスで指定された値に設定してください。

## 失敗のシナリオ

以下のシナリオでは認証は失敗します。

### サーバー証明書が変更された場合

サーバー証明書は以下の場合に変更されるかも知れません:

- サーバーが完全に再インストールされたため新しい証明書に変わった。
- 接続がインターセプトされた（{abbr}`MITM (Machine in the middle)`）。

このような場合、このリモートの設定内のフィンガープリントと証明書のフィンガープリントが一致しないため、クライアントはサーバーへの接続を拒否します。

この場合サーバー管理者に連絡して証明書が実際に変更されたのかを確認するのはユーザー次第です。
実際に変更されたのであれば、証明書を新しいものに置き換えるか、リモートを削除して追加し直すことができます。

### サーバーとの信頼関係が取り消された場合

別の信頼されたクライアントまたはローカルのサーバー管理者がサーバー上で対象のクライアントの信頼エントリを削除した場合、そのクライアントに対するサーバーの信頼関係は取り消されます。

この場合、サーバーは引き続き同じ証明書を使用していますが、すべての API 呼び出しは対象のクライアントが信頼されていないことを示す 403 のステータスコードを返します。
