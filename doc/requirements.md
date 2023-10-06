# 動作環境

## Go

LXDはGo 1.18以上を必要とし、Golangのコンパイラのみでテストされています。
（訳注：以前はgccgoもサポートされていましたがGolangのみになりました）

ビルドには最低2GBのRAMを推奨します。

## 必要なカーネルバージョン

される最小のカーネルバージョンは 5.4 です。

Incusには以下の機能をサポートするカーネルが必要です。

* Namespaces （`pid`、`net`、`uts`、`ipc`と`mount`）
* Seccomp
* Native Linux AIO
  （[`io_setup(2)`](https://man7.org/linux/man-pages/man2/io_setup.2.html)など）

以下のオプションの機能はさらなるカーネルオプションを必要とします。

* Namespaces （`user`と`cgroup`）
* AppArmor （mount mediationに対するUbuntuパッチを含む）
* Control Groups （`blkio`、`cpuset`、`devices`、`memory`、`pids`と`net_prio`）
* CRIU (正確な詳細は CRIU のアップストリームを参照のこと)

さらに使用しているIncusのバージョンで必要とされるほかのカーネルの機能も必要です。

## LXC

Incusは以下のビルドオプションでビルドされたLXC 4.0.0以上を必要とします。

* `apparmor` （もしIncusのAppArmorサポートを使用するのであれば）
* `seccomp`

Ubuntuを含むさまざまなディストリビューションの最近のバージョンを動かすためには、LXCFSもインストールする必要があります。

## QEMU

仮想マシンを利用するにはQEMU 6.0以降が必要です。

## 追加のライブラリ（と開発用のヘッダ）

Incusはデータベースとして`cowsql`を使用しています。
ビルドしセットアップするためには`make deps`を実行してください。

Incusはほかにもいくつかの (たいていはパッケージ化されている)Cライブラリを使用しています。

* `libacl1`
* `libcap2`
* `libuv1`（`cowsql`で使用）
* `libsqlite3` >= 3.25.0（`cowsql`で使用）

ライブラリそのものとライブラリの開発用ヘッダ (`-dev` パッケージ)のすべてをインストールしたことを確認してください。
