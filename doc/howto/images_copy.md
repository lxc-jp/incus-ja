(images-copy)=
# イメージをコピーやインポートするには

イメージをイメージストアに追加するには、他のサーバーからコピーすることもできますし、ファイル（ローカルのファイルまたはウェブサーバー上のファイル）からインポートすることもできます。

## リモートからイメージをコピーする

あるサーバーから別のサーバーにイメージをコピーするには、以下のコマンドを入力します:

    incus image copy [<source_remote>:]<image> <target_remote>:

```{note}
イメージをローカルのイメージストアにコピーするには、コピー先のリモートに `local:` と指定します。
```

すべての利用可能なフラグの一覧は [`incus image copy --help`](incus_image_copy.md) を参照してください。
最も重要なものは以下のとおりです:

`--alias`
: イメージのコピーにエイリアスを割り当てる。

`--copy-aliases`
: コピー元のイメージが持つエイリアスをコピーする。

`--auto-update`
: 元のイメージが更新されたらコピーも更新する。

`--vm`
: エイリアスからコピーする際、仮想マシンを作成するのに使えるイメージをコピーする。

## ファイルからイメージをインポートする

要求される {ref}`image-format` を使用するイメージファイルを持っていれば、イメージストアにインポートできます。

そのようなイメージファイルを取得する方法はいくつかあります:

- 既存のイメージをエクスポートする（{ref}`images-manage-export`参照）
- `distrobuilder`でイメージを生成する（{ref}`images-create-build`参照）
- {ref}`remote image server <remote-image-servers>`からイメージファイルをダウンロードする（イメージをファイルにダウンロードしてインポートするより、{ref}`リモートのイメージを使用する <images-remote>`ほうが通常は簡単なことに注意してください）

### ローカルファイルシステムからインポートする

ローカルファイルシステムからイメージをインポートするには、[`incus image import`](incus_image_import.md) コマンドを使用します。
このコマンドは{ref}`統合イメージ <image-format-unified>`（圧縮されたファイルまたはディレクトリー）と{ref}`分離イメージ <image-format-split>`（2 つのファイル）の両方をサポートします。

1 つのファイルまたはディレクトリーから統合イメージをインポートするには、以下のコマンドを入力します:

    incus image import <image_file_or_directory_path> [<target_remote>:]

分離イメージをインポートするには、以下のコマンドを入力します:

    incus image import <metadata_tarball_path> <rootfs_tarball_path> [<target_remote>:]

どちらの場合も、`--alias`フラグでエイリアスを割り当てられます。
利用可能なすべてのフラグは [`incus image import --help`](incus_image_import.md) を参照してください。

### リモートウェブサーバーからファイルをインポートする

URL を指定してリモートウェブサーバーからイメージファイルをインポートできます。
この方法はイメージをユーザーに配布するためだけに Incus サーバーを稼働させる代わりに使用できます。
必要なのはカスタムヘッダ（{ref}`images-copy-http-headers`参照）をサポートする基本的なウェブサーバーだけです。

イメージファイルは統合イメージ（{ref}`image-format-unified`参照）として提供される必要があります。

リモートウェブサーバーからイメージをインポートするには、以下のコマンドを入力します:

    incus image import <URL>

`--alias`フラグでローカルのイメージにエイリアスを割り当てられます。

(images-copy-http-headers)=
#### カスタムHTTPヘッダ

Incus では以下のカスタム HTTP ヘッダをウェブサーバーで設定する必要があります:

`Incus-Image-Hash`
: ダウンロードされるイメージの SHA256 ハッシュ値。

`Incus-Image-URL`
: イメージをダウンロードする URL。

Incus はサーバーに問い合わせる際に以下のヘッダを設定します:

`Incus-Server-Architectures`
: クライアントがサポートするアーキテクチャのカンマ区切りリスト。

`Incus-Server-Version`
: 使用している Incus のバージョン。
