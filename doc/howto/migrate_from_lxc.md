(migrate-from-lxc)=
# LXC から Incus にコンテナをマイグレートするには

Incus は LXC のコンテナを Incus サーバーにインポートするためのツール（`lxc-to-incus`）を提供しています。
LXC コンテナは Incus サーバーと同じマシン上に存在する必要があります。

このツールは LXC コンテナを分析し、データと設定の両方を新しい Incus コンテナにマイグレートします。

```{note}
あるいは LXC コンテナ内で `incus-migrate` ツールを使用して Incus にマイグレートすることもできます（{ref}`import-machines-to-instances` 参照）。
しかし、このツールは LXC コンテナの設定は一切マイグレートしません。
```

## ツールを取得する

あなたの Incus がインストールされた環境でツールが提供されていない場合、自分でビルドできます。
`go` （バージョン 1.18 以降）がインストールされていることを確認して以下のコマンドでツールを取得してください:

    go install github.com/lxc/incus/cmd/lxc-to-incus@latest

## LXC コンテナを用意する

一度に 1 つのコンテナをマイグレートすることもできますし、同時にあなたのすべての LXC コンテナをマイグレートすることもできます。

```{note}
マイグレートされたコンテナは元のコンテナと同じ名前を使用します。
Incus にインスタンス名としてすでに存在する名前を持つコンテナをマイグレートすることはできません。

このため、マイグレーションプロセスを開始する前に名前の衝突を引き起こす可能性のある LXC コンテナはリネームしてください。
```

マイグレーションプロセスを開始する前に、マイグレートしたいコンテナを停止してください。

## マイグレーションプロセスを開始する

コンテナをマイグレートするには `sudo lxd.lxc-to-incus [flags]` と実行してください。

たとえば、すべてのコンテナをマイグレートするには:

    sudo lxc-to-incus --all

`lxc1` コンテナだけをマイグレートするには:

    sudo lxc-to-incus --containers lxc1

2 つのコンテナ（`lxc1` と `lxc2`）をマイグレートし Incus 内の `my-storage` ストレージプールを使用するには:

    sudo lxc-to-incus --containers lxc1,lxc2 --storage my-storage

実際に実行せずにすべてのコンテナのマイグレートをテストするには:

    sudo lxc-to-incus --all --dry-run

すべてのコンテナをマイグレートするが、`rsync` の帯域幅を 5000 KB/s に限定するには:

    sudo lxc-to-incus --all --rsync-args --bwlimit=5000

すべての利用可能なフラグを確認するには `sudo lxd.lxc-to-incus --help` と実行してください。

```{note}
`linux64` アーキテクチャがサポートされない（`linux64` architecture isn't supported）というエラーが出る場合、ツールを最新版にアップデートするか LXC コンテナ内のアーキテクチャを `linux64` から `amd64` か `x86_64` に変更してください。
```

## 設定を確認する

このツールは LXC の設定と（1 つまたは複数の）コンテナの設定を分析し、可能な限りの範囲で設定をマイグレートします。
以下のような実行結果が出力されます:

```{terminal}
:input: sudo lxc-to-incus --containers lxc1

Parsing LXC configuration
Checking for unsupported LXC configuration keys
Checking for existing containers
Checking whether container has already been migrated
Validating whether incomplete AppArmor support is enabled
Validating whether mounting a minimal /dev is enabled
Validating container rootfs
Processing network configuration
Processing storage configuration
Processing environment configuration
Processing container boot configuration
Processing container apparmor configuration
Processing container seccomp configuration
Processing container SELinux configuration
Processing container capabilities configuration
Processing container architecture configuration
Creating container
Transferring container: lxc1: ...
Container 'lxc1' successfully created
```

マイグレーションプロセスが完了したら、設定を確認し、必要に応じて、マイグレートした Incus コンテナを起動する前に Incus 内の設定を更新してください。
