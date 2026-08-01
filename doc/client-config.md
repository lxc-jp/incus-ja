(client-config)=
# CLI設定ファイル
Incusクライアントは環境データを保管し振る舞いを調整するために設定ファイルを使います。このドキュメントではファイルの構造とオプションを説明します。

```{note}
設定ファイル手で編集するとCLIが不整合な状態になってしまうかもしれないので、[`incus remote`](incus_remote.md)と[`incus alias`](incus_alias.md)コマンド経由でのみ編集するのが良いです。
```

## 設定ファイル
デフォルトでは、Incus CLIは`$HOME/.config/incus/config.yml`を使いますが、`INCUS_CONF`環境変数で指定される任意のディレクトリー内の`config.yml`ファイルを使うこともできます。

ファイルやパスが存在しない場合は、以下のデフォルトを使います：

```yaml
default-remote: local
remotes:
  images:
    protocol: simplestreams
    public: true
    addr: https://images.linuxcontainers.org
  local:
    protocol: incus
    public: false
    addr: unix://
aliases: {}
defaults:
  list_format: ""
  console_type: ""
  console_spice_command: ""
  no_color: false
```

## `remotes`
Incusのリモートの情報は[`incus remote`](incus_remote.md)コマンドで管理され、`remotes`セクションに保管されます。remotesは以下の設定キーを持つことができます。

### `addr`
リモートサーバーのアドレス。IP、FQDN、URLを受け付けます。`addr1,addr2,...`のようにカンマで区切って複数のアドレスを設定できます。

### `last_working_address`
リモートが複数のアドレスを持つ場合に、リモートの最後に動いていたアドレスを登録するためにCLIで使われます。

### `auth_type`
サーバーAPIに認証するのにクライアントが使う方法。`tls`と`oidc`を受け付けます。

### `keepalive`
`keepalive`機能が接続を開いたままにしておくタイムアウトの秒数。

### `project`
リモートを使うときにCLIで使うデフォルトプロジェクト。[`incus project switch`](incus_project_switch.md)コマンドで定義されます。`--project`フラグと`INCUS_PROJECT`環境変数で上書きされます。

### `protocol`
リモートと通信するのにCLIが使うプロトコル。リモートをIncusあるいはイメージサーバーのどちらとして使うかを決めます。以下の選択肢を受け付けます：

- `incus`: ネットワーク超しにアクセスできるプライベートIncusサーバー
- `oci`: アプリケーションコンテナイメージを提供する[Open Container Initiative (OCI)](https://opencontainers.org/)サーバー
- `public`: イメージのみを提供するネットワーク超しにアクセスするIncusサーバー
- `simplestream`: イメージを提供する[Simple Stream](https://git.launchpad.net/simplestreams/tree/)サーバー

### `credentials_helper`
OCIサービス認証を処理するヘルパーコマンドを定義します。

### `public`
リモートがイメージ専用のサーバー（public）かどうかを定義します。

## `default_remotes`
コマンドを実行するのにCLIで使うデフォルトのIncusのリモートを定義します。[`incus remote switch`](incus_remote_switch.md)コマンドで定義されます。この設定は`INCUS_REMOTE`環境変数あるいは`<remote>:`接頭辞で上書きされます。

## `aliases`
[`incus alias`](incus_alias.md)コマンドで管理されるエイリアスを保管します。各エイリアスはエイリアス名とエイリアスのコマンド（引用符で囲んだ）で構成するキー／値のペアです。

## `defaults`
特定のCLIコマンドでdefaultのCLIの振る舞いを決定します。

### `list_format`
一覧がどのように出力されるかを決定します。`--format`（`-f`）フラグと等価です。`csv`、`json`、`table`、`yaml`、`compact`、`markdown`を受け付けます。

### `console_type`
[`incus console`](incus_console.md)コマンドで使うデフォルトのコンソールタイプを決定します。指定できるオプションは以下のとおりです：

- `console`: text based console
- `vga`: graphic UI console

### `console_spice_command`
[`incus console`](incus_console.md)で使うVGAコンソールを提供するための代替のSPICEコマンドを決定します。

### `no_color`
CLIのカラー表示を無効化します。
