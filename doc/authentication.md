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

(authentication-add-certs)=
#### 信頼できる証明書をサーバーに追加する

信頼できるクライアントを追加するには、そのクライアント証明書をサーバーのトラストストアに直接追加するのが望ましい方法です。
これを行うには、クライアント証明書をサーバーにコピーし、[`incus config trust add-certificate <file>`](incus_config_trust_add-certificate.md) で登録します。

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

### ローカルキーの暗号化

`incus` クライアントは暗号化されたクライアントキーもサポートします。上記の方法で生成された鍵は以下のコマンドでパスワードを使って暗号化できます:

```
ssh-keygen -p -o -f .config/incus/client.key
```

```{note}
[`keepalive` mode](remote-keepalive) を有効にしないと、Incus を呼び出すたびにプロンプトが表示され煩わしいかもしれません:

    $ incus list remote-host:
    Password for client.key:
    +------+-------+------+------+------+-----------+
    | NAME | STATE | IPV4 | IPV6 | TYPE | SNAPSHOTS |
    +------+-------+------+------+------+-----------+
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

```{important}
設定済みの OIDC アイデンティティ・プロバイダーで認証されたユーザーは Incus へのフルアクセスを得ます。
ユーザーアクセスを制限するには、{ref}`authorization`も設定する必要があります。
現状では、OIDC と互換のある唯一の認可の方法は{ref}`authorization-openfga`です。
```

(authentication-server-certificate)=
## TLS サーバー証明書

Incus は {abbr}`ACME (Automatic Certificate Management Environment)` サービス（たとえば、[Let's Encrypt](https://letsencrypt.org/)）を使ったサーバー証明書の発行をサポートします。

この機能を有効にするには、以下のサーバー設定をしてください:

- {config:option}`server-acme:acme.domain`: 証明書を発行するドメイン。
- {config:option}`server-acme:acme.email`: ACME サービスのアカウントに使用する email アドレス。
- {config:option}`server-acme:acme.agree_tos`: ACME サービスの利用規約に同意するためには `true` に設定する必要あり。
- {config:option}`server-acme:acme.ca_url`: ACME サービスのディレクトリー URL。デフォルトでは Incus は "Let's Encrypt" を使用。

この機能を利用するには、 Incus は 80 番ポートを開放する必要があります。
これは [HAProxy](http://www.haproxy.org/) のようなリバースプロキシを使用することで実現できます。

以下は `incus.example.net` をドメインとして使用する HAProxy の最小限の設定です。
証明書が発行された後、 Incus は`https://incus.example.net/` でアクセスできます。

```
# Global configuration
global
  log /dev/log local0
  chroot /var/lib/haproxy
  stats socket /run/haproxy/admin.sock mode 660 level admin
  stats timeout 30s
  user haproxy
  group haproxy
  daemon
  ssl-default-bind-options ssl-min-ver TLSv1.2
  tune.ssl.default-dh-param 2048
  maxconn 100000

# Default settings
defaults
  mode tcp
  timeout connect 5s
  timeout client 30s
  timeout client-fin 30s
  timeout server 120s
  timeout tunnel 6h
  timeout http-request 5s
  maxconn 80000

# Default backend - Return HTTP 301 (TLS upgrade)
backend http-301
  mode http
  redirect scheme https code 301

# Default backend - Return HTTP 403
backend http-403
  mode http
  http-request deny deny_status 403

# HTTP dispatcher
frontend http-dispatcher
  bind :80
  mode http

  # Backend selection
  tcp-request inspect-delay 5s

  # Dispatch
  default_backend http-403
  use_backend http-301 if { hdr(host) -i incus.example.net }

# SNI dispatcher
frontend sni-dispatcher
  bind :443
  mode tcp

  # Backend selection
  tcp-request inspect-delay 5s

  # require TLS
  tcp-request content reject unless { req.ssl_hello_type 1 }

  # Dispatch
  default_backend http-403
  use_backend incus-nodes if { req.ssl_sni -i incus.example.net }

# Incus nodes
backend incus-nodes
  mode tcp

  option tcp-check

  # Multiple servers should be listed when running a cluster
  server incus-node01 1.2.3.4:8443 check
  server incus-node02 1.2.3.5:8443 check
  server incus-node03 1.2.3.6:8443 check
```

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
