(images-remote)=
# リモートイメージを使用するには

[`incus`](incus.md) CLI コマンドはいくつかのリモートイメージサーバーを初期設定されています。
概要は{ref}`remote-image-servers`を参照してください。

## 設定されたリモートを一覧表示する

<!-- Include start list remotes -->
設定されたすべてのリモートサーバーを見るには、以下のコマンドを入力します:

    incus remote list

[simple streams形式](https://git.launchpad.net/simplestreams/tree/)を使用するリモートサーバーは純粋なイメージサーバーです。
`incus`形式を使用するサーバーは Incus サーバーであり、イメージサーバーだけとして稼働しているか、通常の Incus サーバーとして稼働するのに加えて追加のイメージを提供しているかのどちらかです。
詳細は{ref}`remote-image-server-types`を参照してください。
<!-- Include end list remotes -->

## リモート上の利用可能なイメージを一覧表示する

サーバー上のすべてのリモートイメージを一覧表示するには、以下のコマンドを入力します:

    incus image list <remote>:

結果をフィルタできます。
手順は{ref}`images-manage-filter`を参照してください。

## リモートサーバーを追加する

どのようにリモートを追加するかはサーバーが使用しているプロトコルに依存します。

### simple streamsサーバーを追加する

simple streams サーバーをリモートとして追加するには、以下のコマンドを入力します:

    incus remote add <remote_name> <URL> --protocol=simplestreams

URL は HTTPS でなければなりません。

### リモートのIncusサーバーを追加する

<!-- Include start add remotes -->
Incus サーバーをリモートして追加するには、以下のコマンドを入力します:

    incus remote add <remote_name> <IP|FQDN|URL> [flags]

認証方法によっては固有のフラグが必要です（例えば、OIDC 認証では[`incus remote add <remote_name> <IP|FQDN|URL> --auth-type=oidc`](incus_remote_add.md)を使います）。
詳細は{ref}`server-authenticate`と{ref}`authentication`を参照してください。

例えば、IP アドレスを指定してリモートを追加するには以下のコマンドを入力します:

    incus remote add my-remote 192.0.2.10

リモートサーバーのフィンガープリントを確認するプロンプトが表示され、リモートで使用している認証方法によってパスワードまたはトークンの入力を求められます。
<!-- Include end add remotes -->

## イメージを参照する

イメージを参照するには、リモートとイメージのエイリアスまたはフィンガープリントをコロンで区切って指定します。
例:

    images:ubuntu/22.04
    images:ubuntu/22.04
    local:ed7509d7e83f

(images-remote-default)=
## デフォルトのリモートを選択する

リモート名前を指定せずにイメージ名だけ指定すると、デフォルトのイメージサーバーが使用されます。

どのサーバーがデフォルトのイメージサーバーと設定されているか表示するには、以下のコマンドを入力します:

    incus remote get-default

別のリモートをデフォルトのイメージサーバーに選択するには、以下のコマンドを入力します:

    incus remote switch <remote_name>
