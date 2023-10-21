# バッケージ作成の推奨
以下は Incus のパッケージ作成者向けの推奨事項です。

以下の推奨に従うとさまざまな Linux ディストリビューションでより期待通りな経験を提供できるでしょう。

## パッケージ

通常は、少なくとも `incus` と `incus-client` パッケージに分割するのが良いでしょう。
後者はデーモンやその依存物をインストールせずに `incus` コマンドだけをインストールできるようにします。

さらに、`fuidshift`、`lxc-to-incus`、`incus-benchmark`、`incus-migrate` のような使用頻度の低いツールを `incus-tools` パッケージに分離すると便利かもしれません。

## グループ

2 つのグループを提供すると良いです:

- `incus-admin` は `unix.socket` ソケットへのアクセスを許可し、Incus への実質的に完全な制御を許可します。
- `incus` は `user.socket`  ソケットへのアクセスを許可し、制限された Incus プロジェクトを利用できるようにします。

## 初期化スクリプト

以下は `systemd` の使用を想定しています。
`systemd` を使用しないディストリビューションでは似たような命名規則に従うのがよいですが、ソケットアクティベーションのような点は一部差異があるでしょう。

- `incus.service` は `incusd` デーモンを起動・停止するメインのユニットです。
- `incus.socket` は `incus.service` ユニット用のソケットアクティベーションのユニットです。存在する場合、`incus.service` は単体では起動しないようにします。
- `incus-user.service` は `incus-user` デーモンを起動・停止するメインのユニットです。
- `incus-user.socket` は `incus-user.service` ユニット用のソケットアクティベーションのユニットです。存在する場合、`incus-user.service` は単体では起動しないようにします。
- `incus-startup.service` は `incusd activateifneeded` コマンドを使って必要であればデーモンの起動をトリガーします。さらに `incusd shutdown` を呼んでホストのシャットダウン時にインスタンスを順番にシャットダウンします。

## バイナリ

`incusd` と `incus-user` デーモンはユーザーの `PATH` の通らない場所に置くのが良いです。
`incus-agent` も同様で、デーモンの `PATH` に存在する必要がありますが、ユーザーは利用できないようにしてください。

ユーザーに利用できるようにするべきメインのバイナリは `incus` です。

これらに加えて、以下のオプショナルなバイナリも利用できるようにしてください:

- `fuidshift`（root のみに限定するべき）
- `incus-benchmark`
- `incus-migrate`
- `lxc-to-incus`
- `lxd-to-incus`（root のみに限定するべき）
