(run-commands)=
# インスタンス内でコマンドを実行するには

Incus では、ネットワークを経由してインスタンスにアクセスする必要なしに、Incus クライアントを使用してインスタンス内でコマンドを実行できます。

コンテナでは、これは常に機能し、Incus によって直接処理されます。
仮想マシンでは、これが機能するには、仮想マシン内で `incus-agent` プロセスが稼働している必要があります。

インスタンス内部でコマンドを実行するには [`incus exec`](incus_exec.md) コマンドを使います。
シェルコマンド（たとえば `/bin/bash`）を実行することで、インスタンスにシェルアクセスできます。

## インスタンス内部でコマンドを実行する

ホストマシンの端末から単一のコマンドを実行するには、[`incus exec`](incus_exec.md) コマンドを使います:

    incus exec <instance_name> -- <command>

たとえば、コンテナ上のパッケージリストを更新するには以下のコマンドを入力します:

    incus exec ubuntu-container -- apt-get update

### 実行モード

Incus はコマンドをインタラクティブにも非インタラクティブにも実行できます。

インタラクティブモードでは、入力（stdin）と出力（stdout, stderr）を扱うために疑似端末装置（PTS）が使用されます。
これは、ターミナル・エミュレータに接続されている場合（スクリプトから実行されていない場合）、CLI によって自動的に選択されます。
インタラクティブ・モードを強制するには、`--force-interactive`か`--mode interactive`をコマンドに追加します。

非インタラクティブ・モードでは、代わりにパイプが (stdin、stdout、stderr のそれぞれに 1 つずつ) 割り当てられます。
これにより、多くのスクリプトで必要とされるように、コマンドを実行しながら、stdin、stdout、stderr を別々に適切に取得することができます。
非インタラクティブ・モードを強制するには、`--force-noninteractive`か`--mode non-interactive`をコマンドに追加します。

### ユーザー、グループ、作業ディレクトリー

Incus はインスタンス内のデータを読まない、あるいはその中にあるものを信用しないというポリシーを持っています。
これは、Incus がユーザーやグループの解決を処理するために、`/etc/passwd`、`/etc/group`や`/etc/nsswitch.conf`のようなものを解析しないことを意味しています。

結果として、Incus はユーザーのホームディレクトリーがどこにあるか、あるいはどのような補助的なグループがあるかを知りません。

デフォルトでは、Incus は root（UID 0）、デフォルトのグループ（GID 0）としてコマンドを実行し、作業ディレクトリーは`/root`に設定されています。
ユーザー、グループ、作業ディレクトリーは以下のフラグによって上書きできます。

- `--user` - コマンドを実行するユーザー ID
- `--group` - コマンドを実行するグループ ID
- `--cwd` - コマンドを実行するディレクトリー

### 環境

以下の 2 つの方法で exec セッションに環境変数を渡せます。

インスタンスオプションとして環境変数を渡す
: インスタンス内で`ENVVAR`環境変数を`VALUE`に設定するには、`environment.ENVVAR`インスタンスオプション を設定します（{config:option}`instance-miscellaneous:environment.*`参照）:

      lxc config set <instance_name> environment.ENVVAR=VALUE

exec コマンドに環境変数を渡す
: exec コマンドに環境変数を渡すには、`--env`フラグを使います。
  たとえば以下のようにします:

      incus exec <instance_name> --env ENVVAR=VALUE -- <command>

さらに、Incus は (上記のいずれかの方法で渡されない限り) 以下のデフォルト値を設定します:

```{list-table}
   :header-rows: 1

* - 変数名
  - 条件
  - 値
* - `PATH`
  - \-
  - Concatenation of:
    - `/usr/local/sbin`
    - `/usr/local/bin`
    - `/usr/sbin`
    - `/usr/bin`
    - `/sbin`
    - `/bin`
    - `/snap` (if applicable)
    - `/etc/NIXOS` (if applicable)
* - `LANG`
  - \-
  - `C.UTF-8`
* - `HOME`
  - running as root (UID 0)
  - `/root`
* - `USER`
  - running as root (UID 0)
  - `root`
```

## インスタンスにシェルアクセスする

インスタンス内で直接コマンドを実行したい場合、インスタンス内でシェルコマンドを実行します。
たとえば、以下のコマンド（インスタンス内に`/bin/bash`コマンドが存在する想定）を入力します。

    lxc exec <instance_name> -- /bin/bash

デフォルトでは`root`ユーザーでログインします。
別のユーザーでログインしたい場合は、以下のコマンドを入力します:

    lxc exec <instance_name> -- su --login <user_name>

```{note}
インスタンス内で稼働しているオペレーティングシステムによっては、先にユーザを作成する必要があるかもしれません。
```

インスタンスシェルを終了するには、`exit`を入力するか`Ctrl`+`d`を押します。
