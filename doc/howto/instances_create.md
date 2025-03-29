(instances-create)=
# インスタンスを作成するには

インスタンスを作成するには、[`incus init`](incus_create.md)か[`incus launch`](incus_launch.md)をコマンドを使用できます。
[`incus init`](incus_create.md)はインスタンスを作成するだけですが、[`incus launch`](incus_launch.md)は作成して起動します。

## 使い方

コンテナを作成するには以下のコマンドを入力します:

    incus launch|init <image_server>:<image_name> <instance_name> [flags]

イメージ
: イメージは必要最小限のオペレーティングシステム（たとえば、Linux ディストリビューション）と Incus 関連の情報を含みます。
  さまざまなオペレーティングシステムのイメージがビルトインのリモートイメージサーバーで利用できます。
  詳細は {ref}`images` を参照してください。

  イメージがローカルにない場合、イメージサーバーとイメージの名前を指定（たとえば、公式の 22.04 Ubuntu イメージなら `images:ubuntu/22.04`）する必要があります。

インスタンス名
: インスタンス名は Incus の運用環境（そしてクラスタ内）でユニークである必要があります。
  追加の要件については {ref}`instance-properties` を参照してください。

フラグ
: フラグの完全なリストについては [`incus launch --help`](incus_launch.md) か [`incus init --help`](incus_create.md) を参照してください。
  よく使うフラグは以下のとおりです。

  - `--config` は新しいインスタンスの設定オプションを指定します
  - `--device` はプロファイルを通して提供されるデバイスの {ref}`デバイスオプション <devices>` を上書き、あるいは {ref}`ルートディスクデバイスの初期設定 <devices-disk-initial-config>` を指定します。
  - `--profile` は新しいインスタンスに使用する {ref}`プロファイル <profiles>` を指定します
  - `--network` や `--storage` は新しいインスタンスに指定のネットワークやストレージプールを使用させます
  - `--target` は指定のクラスタメンバー上にインスタンスを作成します
  - `--vm` はコンテナではなく仮想マシンを作成します

## 設定ファイルを渡す

インスタンス設定をフラグとして指定する代わりに、YAML ファイルでコマンドに渡すことができます。

たとえば、`config.yaml` の設定でコンテナを起動するには、以下のコマンドを入力します:

    incus launch images:ubuntu/22.04 ubuntu-config < config.yaml

    incus launch images:ubuntu/22.04 ubuntu-config < config.yaml

```{tip}
YAML ファイルの必要な文法を見るには既存のインスタンス設定の中身を確認（[`incus config show <instance_name> --expanded`](incus_config_show.md)）してください。
```

##  例

以下の例では [`incus launch`](incus_launch.md) を使用しますが、同じように [`incus init`](incus_create.md) も使用できます。

### システムコンテナを起動する

`images` サーバーの Ubuntu 22.04 のイメージで `ubuntu-container` というインスタンス名でシステムコンテナを起動するには、以下のコマンドを入力します:

    incus launch images:ubuntu/22.04 ubuntu-container

### アプリケーションコンテナを起動する

アプリケーション（OCI）コンテナを起動するには、まずイメージレジストリを追加する必要があります:

    incus remote add oci-docker https://docker.io --protocol=oci

次にレジストリに含まれるイメージの1つからコンテナを起動します:

    incus launch oci-docker:hello-world --ephemeral --console

### 仮想マシンを起動する

`images` サーバーの Ubuntu 22.04 のイメージで `ubuntu-vm` というインスタンス名で仮想マシンを起動するには、以下のコマンドを入力します:

    incus launch images:ubuntu/22.04 ubuntu-vm --vm

より大きいディスクサイズで起動する場合は:

    incus launch images:ubuntu/22.04 ubuntu-vm-big --vm --device root,size=30GiB

### コンテナを指定の設定で起動する

コンテナを起動しリソースを 1 つの vCPU と 192MiB の RAM に限定するには、以下のコマンドを入力します:

    incus launch images:ubuntu/22.04 ubuntu-limited --config limits.cpu=1 --config limits.memory=192MiB

### 指定のクラスタメンバー上で仮想マシンを起動する

クラスタメンバー `server2` 上で仮想マシンを起動するには、以下のコマンドを入力します:

    incus launch images:ubuntu/22.04 ubuntu-container --vm --target server2

### 指定のインスタンスタイプでコンテナを起動する

Incus ではクラウドのシンプルなインスタンスタイプが使えます。これは、インスタンスの作成時に指定できる文字列で表されます。

以下の 3 つの指定方法があります:

- `<instance type>`
- `<cloud>:<instance type>`
- `c<CPU>-m<RAM in GiB>`

たとえば、次の 3 つのインスタンスタイプは同じです:

- `t2.micro`
- `aws:t2.micro`
- `c1-m1`

コマンドラインでは、インスタンスタイプは次のように指定します:

    incus launch images:ubuntu/22.04 my-instance --type t2.micro

使えるクラウドとインスタンスタイプのリストは [`https://github.com/dustinkirkland/instance-type`](https://github.com/dustinkirkland/instance-type) で確認できます。

### ISOからブートする仮想マシンを起動する

```{note}
Windows仮想マシンを作成する際は、`image.os`プロパティを`Windows`から始まる値に確実に設定してください。
そうすることでIncusが仮想マシン内でWindowsが稼働することを想定し挙動を適切に調整します。

これは特に以下のことをもたらします:
 - いくつかの非サポートの仮想デバイスを無効化します
 - {abbr}`RTC (Real Time Clock)`クロックをUTCではなくシステムローカルタイムに基づかせます
 - Intel IOMMUコントローラへ切り替えるようにIOMMUをハンドリングします
```

ISO からブートする仮想マシンを起動するには、まず仮想マシンを作成する必要があります。
ISO イメージから仮想マシンを作成しインストールしたいとしましょう。
このシナリオでは、まず以下のコマンドで空の仮想マシンを作成します:

    incus init iso-vm --empty --vm

```{note}
インストールされているオペレーティングシステムの必要に応じて、より多くの CPU、メモリやストレージを仮想マシンに割り当てたいかもしれません。

例えば、2 CPU、4 GiB メモリと 50 GiB のストレージなら、以下のようにします:

    incus init iso-vm --empty --vm -c limits.cpu=2 -c limits.memory=4GiB -d root,size=50GiB
```

次のステップは ISO イメージをインポートし、後で仮想マシンにストレージボリュームとしてアタッチできるようにします:

    incus storage volume import <pool> <path-to-image.iso> iso-volume --type=iso

最後に、以下のコマンドでカスタム ISO ボリュームを仮想マシンにアタッチする必要があります:

    incus config device add iso-vm iso-volume disk pool=<pool> source=iso-volume boot.priority=10

`boot.priority` 設定キーは仮想マシンの起動順が確実に ISO が最初になるようにします。
仮想マシンを起動し、コンソールに接続してメニューを操作できるようにします:

    incus start iso-vm --console

シリアルコンソールでの操作が完了したら、`ctrl+a-q` を使ってコンソールから切断する必要があります。そして以下のコマンドで VGA コンソールに接続します:

    incus console iso-vm --type=vga

これでインストーラが見えるようになります。インストールが終わったら、カスタム ISO ボリュームを切り離す必要があります:

    incus storage volume detach <pool> iso-volume iso-vm

これで仮想マシンはリブートでき、リブートするとディスクから起動します。

### 仮想マシンインスタンスにIncus Agentをインストール

直接のコマンド実行(`incus exec`)、ファイル転送(`incus file`)、そして詳細な利用状況のメトリクス(`incus info`)のような機能を仮想マシンで使うために、Incusではエージェントソフトウェアが提供されています。

[images](https://images.linuxcontainers.org)の仮想マシンイメージは起動時にこのエージェントをロードするように事前に設定されています。

他の仮想マシンでは、手動でエージェントをインストールすることもできます。

```{note}
Incus Agentは現状ではLinuxの仮想マシンでのみ利用可能です。
```

Incusはマウント名`config`でリモートの`9p`ファイルシステムとして提供します。
エージェントをインストールするには、仮想マシンにアクセスし、以下のコマンドを実行する必要があります:

    mount -t 9p config /mnt
    cd /mnt
    ./install.sh

最初の行はリモートのファイルシステムをマウントポイント`/mnt`にマウントします。
次のコマンドはインストールスクリプト`install.sh`を実行しIncus Agentをインストールと起動します。
このタスクは一度だけ実行すればよいです。
