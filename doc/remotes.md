# リモートサーバーを追加するには

リモートサーバーは Incus コマンドラインクライアント内の概念です。
デフォルトでは、コマンドラインクライアントはローカルの Incus デーモンとやりとりしますが、他のサーバーやクラスタを追加できます。

リモートサーバーの用途の 1 つはローカルサーバーでインスタンスを作成するのに使えるイメージを配布することです。
詳細は{ref}`remote-image-servers`を参照してください。

完全な Incus サーバーをお使いのクライアントにリモートサーバーとして追加することもできます。
この場合、ローカルのデーモンと同様にリモートサーバーとやりとりできます。
例えば、リモートサーバー上のインスタンスを管理したりサーバー設定を更新できます。

## 認証

Incus サーバーをリモートサーバーとして追加できるようにするには、サーバーの API が公開されている必要があります。
それはつまり、{config:option}`server-core:core.https_address`サーバー設定オプションが設定されている必要があることを意味します。

サーバーを追加する際は、{ref}`authentication`の方法で認証する必要があります。

詳細は{ref}`server-expose`を参照してください。

## 追加されたリモートを一覧表示する

% Include parts of the content from file [howto/images_remote.md](howto/images_remote.md)
```{include} howto/images_remote.md
   :start-after: <!-- Include start list remotes -->
   :end-before: <!-- Include end list remotes -->
```

## リモートのIncusサーバーを追加する

% Include parts of the content from file [howto/images_remote.md](howto/images_remote.md)
```{include} howto/images_remote.md
   :start-after: <!-- Include start add remotes -->
   :end-before: <!-- Include end add remotes -->
```

## デフォルトのリモートを選択する

Incus コマンドラインクライアントは`local`リモート、つまりローカルの Incus デーモン、に接続する用に初期設定されています。

別のリモートをデフォルトのリモートとして選択するには、以下のように入力します:

    incus remote switch <remote_name>

どのサーバーがデフォルトのリモートとして設定されているか確認するには、以下のように入力します。

    incus remote get-default

## グローバルのリモートを設定する

グローバルなシステム毎の設定としてリモートを設定できます。
これらのリモートは、設定を追加した Incus サーバーのすべてのユーザーで利用できます。

ユーザーはこれらのシステムで設定されたリモートを（例えば [`incus remote rename`](incus_remote_rename.md)または[`incus remote set-url`](incus_remote_set-url.md)を実行することで）オーバーライドできます。
その結果、リモートと対応する証明書がユーザー設定にコピーされます。

グローバルリモートを設定するには、`/etc/incus/`に置かれた`config.yml`ファイルを編集します。

リモートへの接続用の証明書は同じ場所の`servercerts`ディレクトリー(例えば、 `/etc/incus/servercerts/`)に保管する必要があります。
証明書はリモート名に対応する(例えば、`foo.crt`)必要があります。

以下の設定例を参照してください:

```
remotes:
  foo:
    addr: https://192.0.2.4:8443
    auth_type: tls
    project: default
    protocol: incus
    public: false
  bar:
    addr: https://192.0.2.5:8443
    auth_type: tls
    project: default
    protocol: incus
    public: false
```
