(installing)=
# Incusをインストールするには

Incusをインストールする最も簡単な方法は{ref}`利用可能なパッケージの1つをインストール <installing-from-package>`ですが、{ref}`ソースからIncusをインストール <installing_from_source>`も可能です。

Incusをインストールしたら、システム上に`incus-admin`グループが存在することを確認してください。
このグループのユーザーがIncusを操作できます。
手順は{ref}`installing-manage-access`を参照してください。

## リリースを選択する

% Include content from [support.md](support.md)
```{include} support.md
    :start-after: <!-- Include start release -->
    :end-before: <!-- Include end release -->
```

本番環境にはLTSを推奨します。通常のバグフィクスとセキュリティアップデートの恩恵を受けられるからです。
しかし、長期リリースには新しい機能はやどんな種類の挙動の変更も追加されません。

LXDの最新の機能と毎月の更新を得るには、代わりに機能リリースを使ってください。

(installing-from-package)=
## Incusをパッケージからインストールする

IncusデーモンはLinuxでのみ稼働します。
クライアントツール（[`incus`](incus.md)）はほとんどのプラットフォームで利用できます。

### Linux

LinuxでIncusをインストールする最も簡単な方法は{ref}`installing-zabbly-package`です。これはDebianとUbuntuで利用できます。

(installing-zabbly-package)=
#### ZabblyのDebianとUbuntuパッケージをインストールする
現時点ではIncusをインストールする最も簡単な方法は[Zabbly](https://zabbly.com)で提供されるDebianまたはUbuntuのパッケージを使うことです。
最新の安定版リリースと（テストされていない）デイリービルドの2つのリポジトリがあります。

インストール手順は[`https://github.com/zabbly/incus`](https://github.com/zabbly/incus)にあります。

他のインストール方法については{ref}`installing`を参照してください。

1. あなたのユーザーにIncusを制御する許可を与えます。

   上記のパッケージに含まれるIncusへのアクセスは2つのグループで制御されます。

   - `incus`は基本的なユーザーアクセスを許可します。設定はできずすべてのアクションはユーザーごとのプロジェクトに限定されます。
   - `incus-admin`はIncusの完全なコントロールを許可します。

   すべてのコマンドをrootで実行することなくIncusを制御するには、あなた自身を`incus-admin`グループに追加してください。

       sudo adduser YOUR-USERNAME incus-admin
       newgrp incus-admin

   `newgrp`の手順はあなたの端末セッションを再起動しないままでIncusを利用する場合に必要です（訳注：端末を起動し直す場合は不要です）。

1. Incusを初期化します。

       incus admin init --minimal

   この手順はフォルトのオプションで最小セットアップの構成を作成します。
   初期化オプションをチューニングしたい場合、詳細は{ref}`initialize`を参照してください。

### 他のOS

```{important}
他のOS用のビルドはクライアントのみを含み、サーバは含みません。
```

````{tabs}

```{group-tab} macOS

IncusはmacOSのIncusクライアントのビルドを[Homebrew](https://brew.sh/)で公開しています（訳注：2023-10-07時点で公開されていないようです）。

機能リリースのIncusをインストールするには、以下のようにします。

    brew install incus
```

```{group-tab} Windows

Windows版のIncusクライアントは[Chocolatey](https://community.chocolatey.org/packages/incus)パッケージとして提供されています（訳注：2023-10-07時点では提供されていないようです）。
インストールするためには以下のようにします。

1. [インストール手順](https://docs.chocolatey.org/en-us/choco/setup)に従ってChocolateyをインストールします。
1. Incusクライアントをインストールします。

        choco install incus
```

````

[GitHub](https://github.com/lxc/incus/actions)にもIncusクライアントのネイティブビルドがあります。
特定のビルドをダウンロードするには以下のようにします。

1. GitHubアカウントにログインします。
1. 興味のあるブランチやタグ(たとえば、最新のリリースタグあるいは`main`)でフィルタリングします。
1. 最新のビルドを選択し、適切なアーティファクトをダウンロードします。

(installing_from_source)=
## Incusをソースからインストールする

Incusをソースコードからビルドとインストールしたい場合、以下の手順に従ってください。

Incusの開発には`liblxc`の最新バージョン（4.0.0以上が必要）を使用することをお勧めします。
さらにIncusが動作するためにはGolang 1.18以上が必要です。
Ubuntuでは次のようにインストールできます。

```bash
sudo apt update
sudo apt install acl attr autoconf automake dnsmasq-base git golang libacl1-dev libcap-dev liblxc1 liblxc-dev libsqlite3-dev libtool libudev-dev liblz4-dev libuv1-dev make pkg-config rsync squashfs-tools tar tcl xz-utils ebtables
```

デフォルトのストレージドライバである`dir`ドライバに加えて、Incusではいくつかのストレージドライバが使えます。
これらのツールをインストールすると、initramfsへの追加が行われ、ホストのブートが少しだけ遅くなるかもしれませんが、特定のドライバを使いたい場合には必要です。

```bash
sudo apt install lvm2 thin-provisioning-tools
sudo apt install btrfs-progs
```

テストスイートを実行するには、次のパッケージも必要です。

```bash
sudo apt install busybox-static curl gettext jq sqlite3 socat bind9-dnsutils
```

### ソースから最新版をビルドする

この方法はIncusの最新版をビルドしたい開発者やLinuxディストリビューションで提供されないIncusの特定のリリースをビルドするためのものです。
Linuxディストリビューションへ統合するためのソースからのビルドはここでは説明しません。
それは将来、別のドキュメントで取り扱うかもしれません。

```bash
git clone https://github.com/lxc/incus
cd incus
```

これでIncusの現在の開発ツリーをダウンロードしてソースツリー内に移動します。
その後下記の手順にしたがって実際にIncusをビルド、インストールしてください。

### ソースからリリース版をビルドする

Incusのリリースtarballは完全な依存ツリーと`libraft`とIncusデータベースのセットアップに使用する`libcowsql`のローカルコピーをバンドルしています。

```bash
tar zxvf incus-0.1.tar.gz
cd incus-0.1
```

これでリリースtarballを展開し、ソースツリー内に移動します。
その後下記の手順にしたがって実際にIncusをビルド、インストールしてください。

### ビルドを開始する

実際のビルドはMakefileの2回の別々の実行により行われます。
一つは`make deps`でこれはIncusに必要とされるライブラリをビルドします。
もう一つは`make`でIncus自体をビルドします。
`make deps`の最後に`make`の実行に必要な環境変数を設定するための手順が表示されます。
新しいバージョンのIncusがリリースされたらこれらの環境変数の設定は変わるかもしれませんので、`make deps`の最後に表示された手順を使うようにしてください。
下記の手順（例示のために表示します）はあなたがビルドするIncusのバージョンのものとは一致しないかもしれません。

ビルドには最低2GiBのRAMを搭載することを推奨します。

```{terminal}
:input: make deps

...
make[1]: Leaving directory '/root/go/deps/cowsql'
# environment

Please set the following in your environment (possibly ~/.bashrc)
#  export CGO_CFLAGS="${CGO_CFLAGS} -I$(go env GOPATH)/deps/cowsql/include/ -I$(go env GOPATH)/deps/raft/include/"
#  export CGO_LDFLAGS="${CGO_LDFLAGS} -L$(go env GOPATH)/deps/cowsql/.libs/ -L$(go env GOPATH)/deps/raft/.libs/"
#  export LD_LIBRARY_PATH="$(go env GOPATH)/deps/cowsql/.libs/:$(go env GOPATH)/deps/raft/.libs/:${LD_LIBRARY_PATH}"
#  export CGO_LDFLAGS_ALLOW="(-Wl,-wrap,pthread_create)|(-Wl,-z,now)"
:input: make
```

### ソースからのビルド結果のインストール

ビルドが完了したら、ソースツリーを維持したまま、あなたのお使いのシェルのパスに`$(go env GOPATH)/bin`を追加し、`LD_LIBRARY_PATH`環境変数を`make deps`で表示された値に設定します。これは`~/.bashrc`ファイルの場合は以下のようになります。

```bash
export PATH="${PATH}:$(go env GOPATH)/bin"
export LD_LIBRARY_PATH="$(go env GOPATH)/deps/cowsql/.libs/:$(go env GOPATH)/deps/raft/.libs/:${LD_LIBRARY_PATH}"
```

これで`incusd`と`incus`コマンドの実行ファイルが利用可能になりIncusをセットアップするのに使用できます。
`LD_LIBRARY_PATH`環境変数のおかげで実行ファイルは`$(go env GOPATH)/deps`にビルドされた依存ライブラリを自動的に見つけて使用します。

### マシンセットアップ

LXDが非特権コンテナを作成できるように、rootユーザーに対するsub{u,g}idの設定が必要です。

```bash
echo "root:1000000:1000000000" | sudo tee -a /etc/subuid /etc/subgid
```

これでデーモンを実行できます(`sudo`グループに属する全員がIncusとやりとりできるように `--group sudo` を指定します。別に指定したいグループを作ることもできます)。

```bash
sudo -E PATH=${PATH} LD_LIBRARY_PATH=${LD_LIBRARY_PATH} $(go env GOPATH)/bin/incus --group sudo
```

```{note}
`newuidmap/newgidmap`ツールがシステムに存在し、`/etc/subuid`、`/etc/subgid`が存在する場合は、rootユーザーに少なくとも10MのUID/GIDの連続した範囲を許可するように設定する必要があります。
```

(installing-manage-access)=
## Incusへのアクセスを管理する

Incusのアクセス制御はグループのメンバーシップに基づいています。
rootユーザーと`incus-admin`グループのすべてのメンバーはローカルデーモンとやりとりできます。
詳細は{ref}`security-daemon-access`を参照してください。

お使いのシステムに`incus-admin`グループが存在しない場合は、作成してIncusデーモンを再起動してください。
このグループに追加されたメンバーはIncusの完全な制御ができます。

グループのメンバーシップは通常ログイン時にのみ適用されますので、セッションを開き直すか、Incusとやりとりするシェル上で`newgrp incus-admin`コマンドを実行する必要があります。

````{important}
% Include content from [../README.md](../README.md)
```{include} ../README.md
    :start-after: <!-- Include start security note -->
    :end-before: <!-- Include end security note -->
```
````

(installing-upgrade)=
## Incusをアップグレードする

Incusを新しいバージョンにアップグレードした後、Incusはデータベースを新しいスキーマにアップデートする必要があるかもしれません。
このアップデートはIncusのアップグレードの後のデーモン起動時に自動的に実行されます。
アップデート前のデータベースのバックアップはアクティブなデータベースと同じ場所（`/var/lib/incus/database`）に保存されます。

```{important}
スキーマのアップデート後は、古いバージョンのIncusはデータベースを無効とみなすかもしれません。
これはつまりIncusをダウングレードしてもあなたのIncusの環境は利用不可能と言われるかもしれないということです。

このようなダウングレードが必要な場合は、ダウングレードを行う前にデータベースのバックアップをリストアしてください。
```
