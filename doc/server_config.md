(server)=
# サーバー設定

Incus サーバーは key/value 設定オプションで設定できます。

key/value 設定は名前空間が分けられています。
以下のオプションが利用可能です:

- {ref}`server-options-core`
- {ref}`server-options-acme`
- {ref}`server-options-cluster`
- {ref}`server-options-images`
- {ref}`server-options-logging`
- {ref}`server-options-misc`
- {ref}`server-options-oidc`
- {ref}`server-options-openfga`

設定オプションをどのように設定するかについての手順は{ref}`server-configure`を参照してください。

```{note}
このページの表で`global`スコープと表記されたオプションは即時に全てのクラスタメンバーに適用されます。
`local`スコープと表記されたオプションはメンバーごとに設定する必要があります。
```

(server-options-core)=
## コア設定

以下のサーバーオプションはコアデーモンの設定を制御します:

% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group server-core start -->
    :end-before: <!-- config group server-core end -->
```

(server-options-acme)=
## ACME設定

以下のサーバーオプションは{ref}`ACME <authentication-server-certificate>`設定を制御します:

% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group server-acme start -->
    :end-before: <!-- config group server-acme end -->
```

(server-options-oidc)=
## OpenID Connect 設定

以下のサーバーオプションは{ref}`authentication-openid`による外部ユーザー認証を設定します:

% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group server-oidc start -->
    :end-before: <!-- config group server-oidc end -->
```

(server-options-openfga)=
## OpenFGA 設定

以下のサーバーオプションは {ref}`authorization-openfga` を使った外部ユーザー認可を設定します:

% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group server-openfga start -->
    :end-before: <!-- config group server-openfga end -->
```

(server-options-cluster)=
## クラスタ設定

以下のサーバーオプションは{ref}`clustering`を制御します:

% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group server-cluster start -->
    :end-before: <!-- config group server-cluster end -->
```

(server-options-images)=
## イメージ設定

以下のサーバーオプションは{ref}`images`をどう取り扱うかを設定します:

% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group server-images start -->
    :end-before: <!-- config group server-images end -->
```

(server-options-logging)=
## ログ設定

ログシステムは複数のターゲットが設定できるようになりました。ターゲットはユニークな名前（例：`loki1`、`syslog01`）により識別されます。
各ターゲットは独立して設定でき、固有のログ種別が割り当てられます。

### サポートされるターゲット

- `loki` -  Grafana Lokiサーバーにログを送信
- `syslog` - リモートのsyslogエンドポイントにログを送信

### 設定例

```
logging.loki01.target.type: loki
logging.loki01.target.address: https://loki01.int.example.net
logging.loki01.target.username: foo
logging.loki01.target.password: bar
logging.loki01.types: lifecycle,network-acl
logging.loki01.lifecycle.types: instance

logging.syslog01.target.type: syslog
logging.syslog01.target.address: syslog01.int.example.net
logging.syslog01.target.facility: security
logging.syslog01.types: logging
logging.syslog01.logging.level: warning
```

% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group server-logging start -->
    :end-before: <!-- config group server-logging end -->
```

(server-options-misc)=
## その他設定

以下のサーバーオプションは{ref}`instances`のサーバー固有設定、MAAS 統合、{ref}`OVN <network-ovn>`統合、{ref}`バックアップ <backups>`、{ref}`storage`を設定します:

% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group server-miscellaneous start -->
    :end-before: <!-- config group server-miscellaneous end -->
```

(server-options-user)=
## ユーザーオプション

追加のユーザー定義キーは`user.`ネームスペース内で利用可能です。
ユーザー定義のキーは常に`string`型で`global`スコープになります。
`user.ui.`で始まるキーはウェブUI設定オプションで使用され、認証していないユーザーにも見えることに注意してください。
