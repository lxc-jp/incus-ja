(image-format)=
# イメージ形式

イメージはルートファイルシステムとイメージを記述するメタデータファイルを含みます。
またイメージを使用するインスタンス内部でファイルを生成するためのテンプレートも含められます。

イメージは統合イメージ（単一ファイル）か分離イメージ（2 つのファイル）としてパッケージできます。

## 中身

コンテナのイメージは以下のディレクトリー構造を持ちます:

```
metadata.yaml
rootfs/
templates/
```

仮想マシンのイメージは以下のディレクトリー構造を持ちます:

```
metadata.yaml
rootfs.img
templates/
```

どちらのインスタンスタイプでも、`templates/`ディレクトリーは省略可能です。

### メタデータ

`metadata.yaml`ファイルはイメージが Incus 内で稼働するために関連する情報を含みます。
以下の情報を含んでいます:

```yaml
architecture: x86_64
creation_date: 1424284563
properties:
  description: Ubuntu 22.04 LTS Intel 64bit
  os: Ubuntu
  release: jammy 22.04
templates:
  ...
```

`architecture`と`creation_date`フィールドは必須です。
`properties`フィールドはイメージのデフォルトプロパティのセットを含みます。
`os`, `release`, `name`, `description`フィールドはよく使われますが、必須ではありません。

`templates`フィールドは省略可能です。
テンプレートをどのように設定するかの情報は{ref}`image_format_templates`を参照してください。

### ルートファイルシステム

コンテナでは、`rootfs/`ディレクトリーがコンテナ内のルートディレクトリー（`/`）の完全なファイルシステムツリーを含みます。

仮想マシンは`rootfs/`ディレクトリーの代わりに`rootfs.img` `qcow2`ファイルを使います。
このファイルはメインのディスクデバイスになります。

(image_format_templates)=
### テンプレート（省略可能＿

インスタンス内部でファイルを動的に作成するのにテンプレートを使用できます。
そのためには、`metadata.yaml`ファイル内でテンプレートルールを設定し、`templates/`ディレクトリー内にテンプレートファイルを配置します。

一般的なルールとして、パッケージに所有されるファイルはテンプレート化は決してするべきではないです。そうでないとインスタンスの通常のオペレーションで上書きされてしまうでしょう。

#### テンプレートルール

生成すべき各ファイルに対して、`metadata.yaml`ファイル内でルールを作成します。
たとえば:

```yaml
templates:
  /etc/hosts:
    when:
      - create
      - rename
    template: hosts.tpl
    properties:
      foo: bar
  /etc/hostname:
    when:
      - start
    template: hostname.tpl
  /etc/network/interfaces:
    when:
      - create
    template: interfaces.tpl
    create_only: true
  /home/foo/setup.sh:
    when:
      - create
    template: setup.sh.tpl
    create_only: true
    uid: 1000
    gid: 1000
    mode: 755
```

`when`キーは以下の 1 つ以上を指定できます:

- `create` - 新規インスタンスがイメージから作成された時に実行
- `copy` - 既存インスタンスからインスタンスが作成されたときに実行
- `start` - インスタンスが開始する度に毎回実行

`template`キーは`templates/`ディレクトリー内のテンプレートファイルを指します。

`properties`キーでユーザー定義のテンプレートプロパティをテンプレートファイルに渡せます。

ファイルが存在しない場合にのみ Incus にファイルを作らせ、ファイルが存在する場合は上書きしてほしくない場合は、`create_only`キーをセットします。

`uid`、`gid`、`mode` キーはファイルの所有者とパーミションを制御するのに使えます。

#### テンプレートファイル

テンプレートファイルは[Pongo2](https://www.schlachter.tech/solutions/pongo2-template-engine/)形式を使います。

テンプレートファイルは常に以下のコンテキストを受け取ります。

| 変数           | 型                               | 説明
| -------------- | -------------------------------- | ------------------------------------------------------------------------------------- |
| `trigger`      | `string`                         | テンプレートをトリガーしたイベント名                                                  |
| `path`         | `string`                         | テンプレートを使用するファイルのパス                                                  |
| `instance`     | `map[string]string`              | インスタンスプロパティのキー/値マップ（名前、アーキテクチャ、特権、一時的）           |
| `config`       | `map[string]string`              | インスタンス設定のキー/値マップ                                                       |
| `devices`      | `map[string]map[string]string`   | インスタンスに割り当てられたデバイスのキー/値マップ                                   |
| `properties`   | `map[string]string`              | `metadata.yaml`で指定されたテンプレートプロパティのキー/値マップ                      |

利便性のため、以下の関数が Pongo2 テンプレートにエクスポートされます。

- `config_get("user.foo", "bar")` - `user.foo`の値か、未設定の場合は`"bar"`を返します。

## イメージのtarball

Incus は 2 種類の Incus 固有のイメージ形式、統合 tarball と分離 tarball をサポートします。

これらの tarball は圧縮されていても構いません。
Incus は tarball の広範囲の圧縮アルゴリズムをサポートします。
しかし、互換性のためには`gzip`または`xz`を使うのが良いです。

(image-format-unified)=
### 統合tarball

統合 tarball は単一の tarball（通常`*.tar.xz`）で、イメージの完全な中身を含みます。それにはメタデータ、ルートファイルシステムと省略可能なテンプレートファイルが含まれます。

これが Incus 自身がイメージを公開する際に内部的に使用している形式です。
通常こちらのほうが扱いやすいので、Incus 固有のイメージを作る際は統合形式を使うのが良いです。

この形式のイメージの識別子は tarball の SHA-256 ハッシュ値です。

(image-format-split)=
### 分離tarball

分離イメージは 2 つの tarball から構成されます。
1 つは （通常`*.tar.xz`）はメタデータと省略可能なテンプレートファイルを含む tarball で、もう 1 つは実際のインスタンスデータを含む tarball、`squashfs`、または `qcow2` イメージです。

コンテナでは、2 つめのファイルはたいていは SquashFS でフォーマットされたファイルシステムツリーですが、同じツリーの tarball でもよいです。
仮想マシンでは、2 つめのファイルは`qcow2`でフォーマットされたディスクイメージです。

tarball は外部で圧縮されていてもよい（`.tar.xz`、`.tar.gz`、…）ですが `squashfs` と `qcow2` はそれぞれのネイティブの圧縮オプションで内部を圧縮してもよいです。

この形式は既に利用可能である既存の Incus 以外の rootfs tarball から簡単にイメージをビルドできるように設計されています。
Incus と他のツールの両方で使用するイメージを作りたい場合もこの形式を使うのが良いです。

これらのイメージのイメージ識別子はメタデータとデータファイルを（この順で）結合したものの SHA-256 ハッシュ値です。
