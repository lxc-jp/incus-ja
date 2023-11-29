(installing)=
# Incusをインストールするには

Incus をインストールする最も簡単な方法は{ref}`利用可能なパッケージの1つをインストール <installing-from-package>`ですが、{ref}`ソースからIncusをインストール <installing_from_source>`も可能です。

Incus をインストールしたら、システム上に`incus-admin`グループが存在することを確認してください。
このグループのユーザーが Incus を操作できます。
手順は{ref}`installing-manage-access`を参照してください。

## リリースを選択する

% Include content from [support.md](support.md)
```{include} support.md
    :start-after: <!-- Include start release -->
    :end-before: <!-- Include end release -->
```

本番環境には LTS を推奨します。通常のバグフィクスとセキュリティーアップデートの恩恵を受けられるからです。
しかし、長期リリースには新しい機能はやどんな種類の挙動の変更も追加されません。

Incus の最新の機能と毎月の更新を得るには、代わりに機能リリースを使ってください。

(installing-from-package)=
## Incusをパッケージからインストールする

Incus デーモンは Linux でのみ稼働します。
クライアントツール（[`incus`](incus.md)）はほとんどのプラットフォームで利用できます。

### Linux

Linux で Incus をインストールする最も簡単な方法は{ref}`installing-zabbly-package`です。これは Debian と Ubuntu で利用できます。

また [GitHub](https://github.com/lxc/incus/actions) に Incus クライアントのネイティブビルドがあります:

- Linux 用の Incus クライアント: [`bin.linux.incus.aarch64`](https://github.com/lxc/incus/releases/latest/download/bin.linux.incus.aarch64)、[`bin.linux.incus.x86_64`](https://github.com/lxc/incus/releases/latest/download/bin.linux.incus.x86_64)
- Windows 用の Incus クライアント: [`bin.windows.incus.aarch64.exe`](https://github.com/lxc/incus/releases/latest/download/bin.windows.incus.aarch64.exe)、[`bin.windows.incus.x86_64.exe`](https://github.com/lxc/incus/releases/latest/download/bin.windows.incus.x86_64.exe)
- macOS 用の Incus クライアント: [`bin.macos.incus.aarch64`](https://github.com/lxc/incus/releases/latest/download/bin.macos.incus.aarch64)、[`bin.macos.incus.x86_64`](https://github.com/lxc/incus/releases/latest/download/bin.macos.incus.x86_64)

(installing-zabbly-package)=
#### ZabblyのDebianとUbuntuパッケージをインストールする
現時点では Incus をインストールする最も簡単な方法は[Zabbly](https://zabbly.com)で提供される Debian または Ubuntu のパッケージを使うことです。
最新の安定版リリースと（テストされていない）デイリービルドの 2 つのリポジトリがあります。

インストール手順は[`https://github.com/zabbly/incus`](https://github.com/zabbly/incus)にあります。

他のインストール方法については{ref}`installing`を参照してください。

1. あなたのユーザーに Incus を制御する許可を与えます。

   上記のパッケージに含まれる Incus へのアクセスは 2 つのグループで制御されます。

   - `incus`は基本的なユーザーアクセスを許可します。設定はできずすべてのアクションはユーザーごとのプロジェクトに限定されます。
   - `incus-admin`は Incus の完全なコントロールを許可します。

   すべてのコマンドを root で実行することなく Incus を制御するには、あなた自身を`incus-admin`グループに追加してください。

       sudo adduser YOUR-USERNAME incus-admin
       newgrp incus-admin

   `newgrp`の手順はあなたの端末セッションを再起動しないままで Incus を利用する場合に必要です（訳注：端末を起動し直す場合は不要です）。

1. Incus を初期化します。

       incus admin init --minimal

   この手順はフォルトのオプションで最小セットアップの構成を作成します。
   初期化オプションをチューニングしたい場合、詳細は{ref}`initialize`を参照してください。

### 他のOS

```{important}
他のOS用のビルドはクライアントのみを含み、サーバは含みません。
```

````{tabs}

```{group-tab} macOS

IncusはmacOSのIncusクライアントのビルドを[Homebrew](https://brew.sh/)で公開しています。

機能リリースのIncusをインストールするには、以下のようにします。

    brew install incus
```

```{group-tab} Windows

Windows版のIncusクライアントは[Chocolatey](https://community.chocolatey.org/packages/incus)パッケージとして提供されています。
インストールするためには以下のようにします。

1. [インストール手順](https://docs.chocolatey.org/en-us/choco/setup)に従ってChocolateyをインストールします。
1. Incusクライアントをインストールします。

        choco install incus
```

````

[GitHub](https://github.com/lxc/incus/actions)にも Incus クライアントのネイティブビルドがあります。
特定のビルドをダウンロードするには以下のようにします。

1. GitHub アカウントにログインします。
1. 興味のあるブランチやタグ(たとえば、最新のリリースタグあるいは`main`)でフィルタリングします。
1. 最新のビルドを選択し、適切なアーティファクトをダウンロードします。

(installing_from_source)=
## Incusをソースからインストールする

Incus をソースコードからビルドとインストールしたい場合、以下の手順に従ってください。

Incus の開発には`liblxc`の最新バージョン（4.0.0 以上が必要）を使用することをお勧めします。
さらに Incus が動作するためには最近の Go 言語（{ref}`requirements-go`参照）が動作することが必要です。
Ubuntu では次のようにインストールできます。

```bash
sudo apt update
sudo apt install acl attr autoconf automake dnsmasq-base git libacl1-dev libcap-dev liblxc1 liblxc-dev libsqlite3-dev libtool libudev-dev liblz4-dev libuv1-dev make pkg-config rsync squashfs-tools tar tcl xz-utils ebtables
command -v snap >/dev/null || sudo apt-get install snapd
sudo snap install --classic go
```

```{note}
`liblxc-dev` パッケージを使って `go-lxc` モジュールのビルド時にコンパイルエラーが出た場合、`liblxc` のビルド時に `INC_DEVEL` の値に `0` を指定したか確認してください。確認するためには、`/usr/include/lxc/version.h` を見てください。
もし `INC_DEVEL` の値が `1` なら、`0` に置き換えると問題を回避できます。これは Ubuntu 22.04/22.10 のパッケージのバグです。Ubuntu 23.04/23.10 ではこの問題はありません。
```

デフォルトのストレージドライバーである`dir`ドライバーに加えて、Incus ではいくつかのストレージドライバーが使えます。
これらのツールをインストールすると、initramfs への追加が行われ、ホストのブートが少しだけ遅くなるかもしれませんが、特定のドライバーを使いたい場合には必要です。

```bash
sudo apt install lvm2 thin-provisioning-tools
sudo apt install btrfs-progs
```

テストスイートを実行するには、次のパッケージも必要です。

```bash
sudo apt install busybox-static curl gettext jq sqlite3 socat bind9-dnsutils
```

### ソースから最新版をビルドする

この方法は Incus の最新版をビルドしたい開発者や Linux ディストリビューションで提供されない Incus の特定のリリースをビルドするためのものです。
Linux ディストリビューションへ統合するためのソースからのビルドはここでは説明しません。
それは将来、別のドキュメントで取り扱うかもしれません。

```bash
git clone https://github.com/lxc/incus
cd incus
```

これで Incus の現在の開発ツリーをダウンロードしてソースツリー内に移動します。
その後下記の手順にしたがって実際に Incus をビルド、インストールしてください。

### ソースからリリース版をビルドする

Incus のリリース tarball は完全な依存ツリーと`libraft`と Incus データベースのセットアップに使用する`libcowsql`のローカルコピーをバンドルしています。

```bash
tar zxvf incus-0.1.tar.gz
cd incus-0.1
```

これでリリース tarball を展開し、ソースツリー内に移動します。
その後下記の手順にしたがって実際に Incus をビルド、インストールしてください。

### ビルドを開始する

実際のビルドは Makefile の 2 回の別々の実行により行われます。
一つは`make deps`でこれは Incus に必要とされるライブラリをビルドします。
もう一つは`make`で Incus 自体をビルドします。
`make deps`の最後に`make`の実行に必要な環境変数を設定するための手順が表示されます。
新しいバージョンの Incus がリリースされたらこれらの環境変数の設定は変わるかもしれませんので、`make deps`の最後に表示された手順を使うようにしてください。
下記の手順（例示のために表示します）はあなたがビルドする Incus のバージョンのものとは一致しないかもしれません。

ビルドには最低 2GiB の RAM を搭載することを推奨します。

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

これで`incusd`と`incus`コマンドの実行ファイルが利用可能になり Incus をセットアップするのに使用できます。
`LD_LIBRARY_PATH`環境変数のおかげで実行ファイルは`$(go env GOPATH)/deps`にビルドされた依存ライブラリを自動的に見つけて使用します。

### マシンセットアップ

Incus が非特権コンテナを作成できるように、root ユーザーに対する sub{u,g}id の設定が必要です。

```bash
echo "root:1000000:1000000000" | sudo tee -a /etc/subuid /etc/subgid
```

これでデーモンを実行できます(`sudo`グループに属する全員が Incus とやりとりできるように `--group sudo` を指定します。別に指定したいグループを作ることもできます)。

```bash
sudo -E PATH=${PATH} LD_LIBRARY_PATH=${LD_LIBRARY_PATH} $(go env GOPATH)/bin/incusd --group sudo
```

```{note}
`newuidmap/newgidmap`ツールがシステムに存在し、`/etc/subuid`、`/etc/subgid`が存在する場合は、rootユーザーに少なくとも10MのUID/GIDの連続した範囲を許可するように設定する必要があります。
```

(installing-manage-access)=
## Incusへのアクセスを管理する

Incus のアクセス制御はグループのメンバーシップに基づいています。
root ユーザーと`incus-admin`グループのすべてのメンバーはローカルデーモンとやりとりできます。
詳細は{ref}`security-daemon-access`を参照してください。

お使いのシステムに`incus-admin`グループが存在しない場合は、作成して Incus デーモンを再起動してください。
このグループに追加されたメンバーは Incus の完全な制御ができます。

グループのメンバーシップは通常ログイン時にのみ適用されますので、セッションを開き直すか、Incus とやりとりするシェル上で`newgrp incus-admin`コマンドを実行する必要があります。

````{important}
% Include content from [../README.md](../README.md)
```{include} ../README.md
    :start-after: <!-- Include start security note -->
    :end-before: <!-- Include end security note -->
```
````

(installing-upgrade)=
## Incusをアップグレードする

Incus を新しいバージョンにアップグレードした後、Incus はデータベースを新しいスキーマにアップデートする必要があるかもしれません。
このアップデートは Incus のアップグレードの後のデーモン起動時に自動的に実行されます。
アップデート前のデータベースのバックアップはアクティブなデータベースと同じ場所（`/var/lib/incus/database`）に保存されます。

```{important}
スキーマのアップデート後は、古いバージョンのIncusはデータベースを無効とみなすかもしれません。
これはつまりIncusをダウングレードしてもあなたのIncusの環境は利用不可能と言われるかもしれないということです。

このようなダウングレードが必要な場合は、ダウングレードを行う前にデータベースのバックアップをリストアしてください。
```
