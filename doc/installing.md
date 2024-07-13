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

いくつかの Linux ディストリビューションでは、メインレポジトリまたはサードパーティレポジトリでパッケージが利用できます。

````{tabs}

```{group-tab} Alpine
Incus と全ての依存ソフトウェアは Alpine Linux のメインレポジトリ内で `incus` として利用できます。

Incus を以下のコマンドでインストールします:

    apk add incus incus-client

仮想マシンを動かす場合は、さらに以下を実行します:

    apk add incus-vm

次にサービスを有効化と起動します:

    rc-update add incusd
    rc-service incusd start

パッケージングの問題は[こちら](https://gitlab.alpinelinux.org/alpine/aports/-/issues)に報告してください。
```

```{group-tab} Arch Linux
Incus と全ての依存ソフトウェアは Arch Linux のメインレポジトリ内で `incus` として利用できます。

Incus を以下のコマンドでインストールします:

    pacman -S incus

パッケージに問題があれば[こちら](https://gitlab.archlinux.org/archlinux/packaging/packages/incus)に報告してください。
```

```{group-tab} Debian
Debian ユーザーには現在 3 つの選択肢があります。

1. ネイティブの `incus` パッケージ

    ネイティブの `incus` パッケージは現在 Debian の testing と unstable レポジトリ内で利用できます。
    このパッケージは次回の Debian 13 (`trixie`) リリース内に含まれる予定です。

    それらのシステムでは、単に`apt install incus`と実行すれば Incus がインストールされます。
    仮想マシンを動かすには、さらに`apt install qemu-system`を実行します。
    LXDからマイグレートする場合は、`lxd-to-incus`コマンドを取得するため`apt install incus-tools`も実行します。

1. ネイティブの `incus` のバックポートされたパッケージ

   ネイティブの `incus` のバックポートされたパッケージが現在 Debian 12 (`bookworm`) ユーザーに提供されています。

    Debian 12 のシステムでは、単に `apt install incus/bookworm-backports` と実行すれば Incus がインストールされます。
    仮想マシンを動かすには、さらに`apt install qemu-system`を実行します。
    LXDからマイグレートする場合は、`lxd-to-incus`コマンドを取得するため`apt install incus-tools`も実行します。

   ****NOTE:**** バックポートされたパッケージのユーザーは Debian Bug Tracker にはバグ報告しないでください。代わりに [Incus のフォーラム](https://discuss.linuxcontainers.org) に報告するか Debian のパッケージ作成者に直接報告してください。

1. Zabbly パッケージレポジトリ

    [Zabbly](https://zabbly.com) は Debian の安定版リリース (11 と 12) 用の最新でありサポート対象である Incus のパッケージを提供します。
    これらのパッケージは Incus の全ての機能を使用するために必要なすべてを含んでいます。

    最新のインストール手順はこちらを参照してください: [`https://github.com/zabbly/incus`](https://github.com/zabbly/incus)
```

```{group-tab} Docker
Zabblyのパッケージレポジトリをベースにした、IncusのDocker/Podmanのイメージが利用手順付きで[`ghcr.io/cmspam/incus-docker`](https://ghcr.io/cmspam/incus-docker)で提供されています。
```

```{group-tab} Fedora
Incus の RPM パッケージとその依存パッケージはまだ公式の Fedora リポジトリでは利用できませんが、[`ganto/lxc4`](https://copr.fedorainfracloud.org/coprs/ganto/lxc4/) の Community Project (COPR) リポジトリでは利用できます。

`dnf` の COPR プラグインをインストールし、次に COPR レポジトリを有効にしてください:

    dnf install 'dnf-command(copr)'
    dnf copr enable ganto/lxc4

Incus を以下のコマンドでインストールします:

    dnf install incus

追加のセットアップ手順については[Getting started with Incus on Fedora](https://github.com/ganto/copr-lxc4/wiki/Getting-Started-with-Incus-on-Fedora)を参照してください。

これは Incus や Fedora の公式なプロジェクトでないことに注意してください。
パッケージの問題は[こちら](https://github.com/ganto/copr-lxc4/issues)に報告してください。
```

```{group-tab} Gentoo
Incus の全ての依存ソフトウェアは Gentoo のメインレポジトリ内に [`app-containers/incus`](https://packages.gentoo.org/packages/app-containers/incus) として利用できます。

Incus は以下のコマンドでインストールできます:

    emerge -av app-containers/incus

仮想マシンを動かす場合は、さらに以下を実行します:

    emerge -av app-emulation/qemu

重要: Incus のアップストリームと Gentoo のレポジトリに LTS と機能リリースが利用できるときに、どちらをインストールするかは後で説明します。

Incus に関連して 2 つのグループが作成されます:
`incus` は（コンテナを起動する）基本的なユーザーアクセスで、`incus-admin` は `incus admin` の制御用です。 あなたのセットアップとユースケースに応じて、通常使用するユーザーをどちらか、あるいは両方に追加してください。

インストールの後、Incus を設定できます。ですがデフォルトのままでも動くので、これは必須ではありません。

- **`openrc`**: `/etc/conf.d/incus` を編集します
- **`systemd`**: `systemctl edit --full incus.service`

`/etc/subuid` と `/etc/subgid` をセットアップします:

    echo "root:1000000:1000000000" | tee -a /etc/subuid /etc/subgid

詳細は: {ref}`User Namespace 用の ID のマッピング <userns-idmap>`

デーモンを起動します:

- **`openrc`**: `rc-service incus start`
- **`systemd`**: `systemctl start incus`

続きは [Gentoo Wiki](https://wiki.gentoo.org/wiki/Incus) 参照。
```

```{group-tab} NixOS
Incus とその依存ソフトウェアは NixOS でパッケージされていて NixOS のオプションで設定できます。利用可能なオプション一式については [`virtualisation.incus`](https://search.nixos.org/options?query=virtualisation.incus) を参照してください。

NixOS 設定に以下を加えるとサービスを有効化し開始できます。

    virtualisation.incus.enable = true;

Incus の初期化は手動で `incus admin init` を使ってもできますし、 NixOS 設定のプリシードオプションでもできます。プリシードの例は NixOS のドキュメントを参照してください。

    virtualisation.incus.preseed = {};

最後に、ユーザーを `incus-admin` グループに追加して、非ルートユーザーに Incus ソケットへのアクセス権を追加できます。それには NixOS 設定に以下を追加します:

    users.users.YOUR_USERNAME.extraGroups = ["incus-admin"];

NixOS 固有の問題については、パッケージレポジトリ内で[イシューを起票](https://github.com/NixOS/nixpkgs/issues/new/choose)してください。
```

```{group-tab} Ubuntu
Ubuntu ユーザーには現在 2 つの選択肢があります。

1. ネイティブの `incus` パッケージ

    ネイティブの `incus` パッケージは現在 Ubuntu 24.04 LTS 以降で利用できます。
    それらのシステムでは、単に`apt install incus`と実行すれば Incus がインストールされます。
    仮想マシンを動かすには、さらに`apt install qemu-system`を実行します。
    LXDからマイグレートする場合は、`lxd-to-incus`コマンドを取得するため`apt install incus-tools`も実行します。

1. Zabbly パッケージレポジトリ

    [Zabbly](https://zabbly.com) は Ubuntu の LTS リリース (20.04 と 22.04) 用の最新でありサポート対象である Incus のパッケージを提供します。
    これらのパッケージは Incus の全ての機能を使用するために必要なすべてを含んでいます。

    最新のインストール手順はこちらを参照してください: [`https://github.com/zabbly/incus`](https://github.com/zabbly/incus)
```

```{group-tab} Void Linux
Incus と全ての依存ソフトウェアは Void Linux のレポジトリ内で `incus` として利用できます。

Incus を以下のコマンドでインストールします:

    xbps-install incus incus-client

次に以下のコマンドでサービスを有効化し起動します:

    ln -s /etc/sv/incus /var/service
    ln -s /etc/sv/incus-user /var/service
    sv up incus
    sv up incus-user


パッケージに問題があれば[こちら](https://github.com/void-linux/void-packages/issues)に報告してください。
```

````

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

Windows版のIncusクライアントは[Chocolatey](https://community.chocolatey.org/packages/incus)と[Winget](https://github.com/microsoft/winget-cli)のパッケージとして提供されています。
Chocolatey または Winget を使ってインストールするには、以下の手順に従ってください:

**Chocolatey**

1. [インストール手順](https://docs.chocolatey.org/en-us/choco/setup)に従ってChocolateyをインストールします。
1. Incusクライアントをインストールします:

        choco install incus

**Winget**

1. [インストール手順](https://learn.microsoft.com/en-us/windows/package-manager/winget/#install-winget)に従って Winget をインストールします。
1. Incusクライアントをインストールします:

        winget install LinuxContainers.Incus
```

````

[GitHub](https://github.com/lxc/incus/actions)にも Incus クライアントのネイティブビルドがあります:

- Linux 用 Incus クライアント: [`bin.linux.incus.aarch64`](https://github.com/lxc/incus/releases/latest/download/bin.linux.incus.aarch64)、[`bin.linux.incus.x86_64`](https://github.com/lxc/incus/releases/latest/download/bin.linux.incus.x86_64)
- Windows 用 Incus クライアント: [`bin.windows.incus.aarch64.exe`](https://github.com/lxc/incus/releases/latest/download/bin.windows.incus.aarch64.exe)、[`bin.windows.incus.x86_64.exe`](https://github.com/lxc/incus/releases/latest/download/bin.windows.incus.x86_64.exe)
- macOS 用 Incus クライアント: [`bin.macos.incus.aarch64`](https://github.com/lxc/incus/releases/latest/download/bin.macos.incus.aarch64)、[`bin.macos.incus.x86_64`](https://github.com/lxc/incus/releases/latest/download/bin.macos.incus.x86_64)

(installing_from_source)=
## Incusをソースからインストールする

Incus をソースコードからビルドとインストールしたい場合、以下の手順に従ってください。

Incus の開発には`liblxc`の最新バージョン（5.0.0 以上が必要）を使用することをお勧めします。
さらに Incus が動作するためには最近の Go 言語（{ref}`requirements-go`参照）が動作することが必要です。

````{tabs}

```{group-tab} Alpine Linux
以下のコマンドで Alpine Linux 上で Incus をビルドするのに必要な開発リソースを取得できます:

    apk add acl-dev autoconf automake eudev-dev gettext-dev go intltool libcap-dev libtool libuv-dev linux-headers lz4-dev tcl-dev sqlite-dev lxc-dev make xz

Incus の必要な機能をすべて使えるようにするには、さらにパッケージをインストールする必要があります。
[Alpine Linux レポジトリの LXD パッケージの定義](https://gitlab.alpinelinux.org/alpine/infra/aports/-/blob/master/community/lxd/APKBUILD) から特有の関数を使う必要のあるパッケージのリストを参照できます。<!-- wokeignore:rule=master -->
また [Alpine Linux パッケージコンテンツフィルター](https://pkgs.alpinelinux.org/contents) から実行ファイル名でパッケージを見つけることができます。

メインの依存ソフトウェアをインストールします:

    apk add acl attr ca-certificates cgmanager dbus dnsmasq lxc libintl iproute2 iptables netcat-openbsd rsync squashfs-tools shadow-uidmap tar xz

仮想マシンを動かすのに必要な追加の依存ソフトウェアをインストールします:

    apk add qemu-system-x86_64 qemu-chardev-spice qemu-hw-usb-redirect qemu-hw-display-virtio-vga qemu-img qemu-ui-spice-core ovmf sgdisk util-linux-misc virtiofsd

リリース tarball あるいは git レポジトリからソースを準備した後、ビルド中の既知のバグを回避するため、以下の手順に従う必要があります:


****重要:**** システムに `/usr/local/include` が存在しない場合、ビルドエラーが出るかもしれません。
また、[`gettext` の問題](https://github.com/gosexy/gettext/issues/1)のため、以下の追加の環境変数を設定する必要があるかもしれません:

    export CGO_LDFLAGS="$CGO_LDFLAGS -L/usr/lib -lintl"
    export CGO_CPPFLAGS="-I/usr/include"
```

```{group-tab} Debian と Ubuntu
ビルドと実行時の依存ソフトウェアをインストールします:

    sudo apt update
    sudo apt install acl attr autoconf automake dnsmasq-base git golang-go libacl1-dev libcap-dev liblxc1 lxc-dev libsqlite3-dev libtool libudev-dev liblz4-dev libuv1-dev make pkg-config rsync squashfs-tools tar tcl xz-utils ebtables

****NOTE:**** DebianやUbuntuの`golang-go`のバージョンはIncusをビルドするのに必要なバージョンより古いかもしれません（{ref}`requirements-go`参照）。
そのような場合は、[upstreamから](https://go.dev/doc/install)新しいGoをインストールする必要があるかもしれません。

デフォルトのストレージドライバーである`dir`ドライバーに加えて、Incus ではいくつかのストレージドライバーが使えます。
これらのツールをインストールすると、initramfs への追加が行われ、ホストのブートが少しだけ遅くなるかもしれませんが、特定のドライバーを使いたい場合には必要です。

    sudo apt install btrfs-progs
    sudo apt install ceph-common
    sudo apt install lvm2 thin-provisioning-tools
    sudo apt install zfsutils-linux

テストスイートを実行するには、次のパッケージも必要です。

    sudo apt install busybox-static curl gettext jq sqlite3 socat bind9-dnsutils

****重要:**** `liblxc-dev` パッケージを使って `go-lxc` モジュールのビルド時にコンパイルエラーが出た場合、`liblxc` のビルド時に `INC_DEVEL` の値に `0` を指定したか確認してください。確認するためには、`/usr/include/lxc/version.h` を見てください。
もし `INC_DEVEL` の値が `1` なら、`0` に置き換えると問題を回避できます。これは Ubuntu 22.04/22.10 のパッケージのバグです。Ubuntu 23.04/23.10 ではこの問題はありません。

```

```{group-tab} OpenSUSE
以下のコマンドで OpenSUSE Tumbleweed システム上で Incus をビルドするのに必要な開発リソースを取得できます:

    sudo zypper install autoconf automake git go libacl-devel libcap-devel liblxc1 liblxc-devel sqlite3-devel libtool libudev-devel liblz4-devel libuv-devel make pkg-config tcl

さらに、通常の運用方法であれば、以下のコマンドも必要になるでしょう:

    sudo zypper install dnsmasq squashfs xz rsync tar attr acl qemu qemu-img qemu-spice qemu-hw-display-virtio-gpu-pci iptables ebtables nftables

```


````

```{note}
ARM64のCPUではUEFIで仮想マシンを扱うためにOVMFではなくAAVMFをインストールする必要があります。
一部のディストリビューションではこのために別のパッケージのインストールが必要です。
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
tar zxvf incus-6.0.0.tar.gz
cd incus-6.0.0
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
