# 動作環境

(requirements-go)=
## Go

Incus は Go 1.23 以上を必要とし、Go 言語のコンパイラのみでテストされています。

ビルドには最低 2GB の RAM を推奨します。

## 必要なカーネルバージョン

サポートされる最小のカーネルバージョンは 5.15 です。

Incus には以下の機能をサポートするカーネルが必要です。

* Namespaces （`pid`、`net`、`uts`、`ipc`と`mount`）
* Seccomp
* Native Linux AIO
  （[`io_setup(2)`](https://man7.org/linux/man-pages/man2/io_setup.2.html)など）

以下のオプションの機能はさらなるカーネルオプションを必要とします。

* Namespaces （`user`と`cgroup`）
* AppArmor
* Control Groups （`blkio`、`cpuset`、`devices`、`memory`、`pids`）
* CRIU (正確な詳細は CRIU のアップストリームを参照のこと)

さらに使用している Incus のバージョンで必要とされるほかのカーネルの機能も必要です。

## LXC

Incus は以下のビルドオプションでビルドされた LXC 5.0.0 以上を必要とします。

* `apparmor` （もし Incus の AppArmor サポートを使用するのであれば）
* `seccomp`

コンテナ内のリソース消費を適切にレポートするために、LXCFSのインストールを強く推奨します。

## OCI

OCIコンテナを動かすには、Incusは現状では`skopeo`に依存しています。
`skopeo`はユーザーの`PATH`に存在する必要があります。

## QEMU

仮想マシンを利用するには QEMU 6.0 以降が必要です。

`virtiofsd`を使用する場合、`virtiofsd`の[Rustでのリライト](https://gitlab.com/virtio-fs/virtiofsd)のみサポートされます。

## OVS/OVN

IncusでOVNネットワークを使う際に必要なOVSとOVNの最小バージョンは以下のとおりです:

* OVS: 2.15.0
* OVN: 23.03.0

## 追加のライブラリ（と開発用のヘッダ）

Incus はデータベースとして`cowsql`を使用しています。
ビルドしセットアップするためには`make deps`を実行してください。

Incus はほかにもいくつかの (たいていはパッケージ化されている)C ライブラリを使用しています。

* `libacl1`
* `libcap2`
* `libuv1`（`cowsql`で使用）
* `libsqlite3` >= 3.25.0（`cowsql`で使用）

ライブラリそのものとライブラリの開発用ヘッダ (`-dev` パッケージ)のすべてをインストールしたことを確認してください。
