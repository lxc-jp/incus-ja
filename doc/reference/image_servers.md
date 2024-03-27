(image-servers)=
# デフォルトイメージサーバー

[`incus`](incus.md) CLI コマンドは下記のデフォルトリモートイメージサーバーが初期設定されています:

`images:`
: このサーバーはさまざまな Linux ディストリビューションの非公式イメージを提供します。
  イメージは[Linux Containers](https://linuxcontainers.org/)チームによりメンテナンスされ、コンパクトで最小限にビルドされています。

  利用可能なイメージの概要については[`images.linuxcontainers.org`](https://images.linuxcontainers.org)を参照してください。

さらなるイメージサーバーは `incus remote add` で追加できます。

(image-server-types)=
## イメージサーバータイプ

Incus は下記のタイプのリモートイメージサーバーをサポートします:

simple streams サーバー
: [simple streams形式](https://git.launchpad.net/simplestreams/tree/)を使う純粋なイメージサーバー。
  静的ファイルのみで構成されますので、これらのサーバーを動かすのに特別なソフトウェアは必要ありません。
  デフォルトの `images:` サーバーは simplestreams タイプを使っています。

公開 Incus サーバー
: イメージを配布するためだけに稼働し、このサーバー自身ではインスタンスを稼働しない Incus サーバー。

  Incus サーバーをポート 8443 で公開で利用可能にするには、{config:option}`server-core:core.https_address`設定オプションを`:8443`に設定し、認証方法をなにも設定しないようにします（詳細は{ref}`server-expose`参照）。
  そして共有したいイメージを`public`にセットします。

Incus サーバー
: ネットワーク越しに管理できる通常の Incus サーバー、イメージサーバーとしても利用可能。

  セキュリティー上の理由により、リモート API へのアクセスを制限し、アクセス制御のための認証方法を設定するほうが良いです。
  詳細な情報は{ref}`server-expose`と{ref}`authentication`を参照してください。

(image-server-tooling)=
## simplestreams サーバーを管理するためのツール
Incus には Simple streams 形式を使ってファイルシステムツリーを管理するのに使える `incus-simplestreams` というツールが含まれています。

このツールは`incus-simplestreams add`でコンテナ（`squashfs`）か仮想マシン（`qcow2`）イメージのいずれかのインポートをサポートし、`incus-simplestreams list`で利用可能なすべてのイメージをフィンガープリントとともに一覧表示し、`incus-simplestreams remove`でサーバーからイメージを削除できます。

このファイルシステムツリーは有効な証明書を持つ HTTPS をサポートする通常のウェブサーバー上に配置する必要があります。

Incus メタデータ tarball が付属しないイメージをインポートする際には、`incus-simplestreams generate-metadata` コマンドを使うといくつかの質問を元に基本的なメタデータの tarball を生成できます。
