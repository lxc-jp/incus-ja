(import-machines-to-instances)=
# 物理または仮想マシンを Incus インスタンスにインポートするには


Incus は既存のディスクやイメージに基づく Incus インスタンスを作成するツール（`incus-migrate`）を提供しています。

このツールは Linux マシン上で実行できます。
まず Incus サーバーに接続して空のインスタンスを作成します。このインスタンスはマイグレーション中またはマイグレーション後に設定を変更できます。
次にこのツールはあなたが用意したディスクまたはイメージからインスタンスにデータをコピーします。

```{note}
マイグレーションプロセスの最中に新しいインスタンスを設定したい場合は、マイグレーションプロセスを開始する前にあなたのインスタンスで使用したいエンティティをセットアップしてください。

デフォルトでは、新しいインスタンスは `default` プロファイルで指定されたエンティティを使用します。
設定をカスタマイズするために異なるプロファイル（あるいはプロファイルのリスト）を設定できます。
詳細は {ref}`profiles` を参照してください。
また、使用される {ref}`instance-options`、{ref}`storage pool <storage-pools>`、{ref}`ストレージボリューム <storage-volumes>`のサイズ、{ref}`network <networking>` をオーバーライドできます。

あるいは、マイグレーションの完了後にインスタンス設定を変更することもできます。
```

このツールはコンテナと仮想マシンの両方を作成できます:

* コンテナを作成する際は、コンテナのルートファイルシステムを含むディスクまたはパーティションを用意する必要があります。
  たとえば、これはあなたがツールを実行しているマシンまたはコンテナの `/` ルートディスクかもしれません。
* 仮想マシンを作成する際は、起動可能なディスク、パーティション、またはイメージを用意する必要があります。
  これは単にファイルシステムを用意するだけでは不十分であり、実行中のコンテナから仮想マシンを作成することはできないことを意味します。
  また使用中の物理マシンから仮想マシンを作成することもできません。これはマイグレーションツールがコピーしようとするディスクを使用するからです。
  代わりに、起動可能なディスク、起動可能なパーティション、または現在使用中でないディスクを用意してください。

   ````{tip}
   （Q35/`virtio-scsi` ありの QEMU/KVM 以外の）外部のハイパーバイザーから Windows VM を変換したい場合、 Windows に `virtio-win` ドライバーをインストールする必要があります。さもないと VM は起動しません。

   <details>
   <summary>Windows VM にどのようにひすようなドライバーを統合するかを見るには展開</summary>
   ホスト上で必要なツールをインストール:

   1. `virt-v2v` バージョン 2.3.4 以上（これが `--block-driver` オプションに対応する最小バージョン）をインストール。
   1. `virtio-win` パッケージをインストール、あるいは [`virtio-win.iso`](https://fedorapeople.org/groups/virt/virtio-win/direct-downloads/stable-virtio/virtio-win.iso) イメージをダウンロードして `/usr/share/virtio-win` フォルダーに配置。
   1. さらに [`rhsrvany`](https://github.com/rwmjones/rhsrvany) にインストールする必要がある場合もある。

   これで `virt-v2v` を使って外部のハイパーバイザーのイメージを Incus の `raw` イメージに変換して必要なドライバーを含められます:

   ```
   # 例 1。vmdk ディスクイメージを incus-migrate に適した raw イメージに変換する
   sudo virt-v2v --block-driver virtio-scsi -o local -of raw -os ./os -i vmx ./test-vm.vmx
   # 例 2。QEMU/KVM qcow2 イメージを変換し virtio-scsi ドライバーを統合する
   sudo virt-v2v --block-driver virtio-scsi -o local -of raw -os ./os -if qcow2 -i disk test-vm-disk.qcow2
   ```

   結果のイメージは `os` ディレクトリー内に生成され、次のステップで `incus-migrate` で使えます。
   </details>
   ````

既存のマシンを Incus インスタンスにマイグレートするには以下の手順を実行してください:

1. 最新の [Incus release](https://github.com/lxc/incus/releases) の **Assets** セクションから `bin.linux.incus-migrate` ツール（[`bin.linux.incus-migrate.aarch64`](https://github.com/lxc/incus/releases/latest/download/bin.linux.incus-migrate.aarch64)または[`bin.linux.incus-migrate.x86_64`](https://github.com/lxc/incus/releases/latest/download/bin.linux.incus-migrate.x86_64)）をダウンロードしてください。
1. ツールをインスタンスを作成したいマシン上に配置して
   （通常 `chmod u+x bin.linux.incus-migrate` を実行して）実行可能にしてください。
1. マシンに `rsync` がインストールされているか確認してください。
   インストールされていない場合は（たとえば、`sudo apt install rsync` で）インストールしてください。
1. ツールを実行します:

       sudo ./bin.linux.incus-migrate

   ツールはマイグレーションに必要な情報を入力するようプロンプトを出します。

   ```{tip}
   ツールをインタラクティブに実行する代わりの方法として、設定をパラメータでコマンドに指定することもできます。
   詳細は `./bin.linux.incus-migrate --help` を参照してください。
   ```

   1. Incus サーバーの URL を、 IP アドレスまたは DNS 名で指定してください。

      ```{note}
      Incus サーバーは {ref}`ネットワークに公開 <server-expose>` する必要があります。
      ローカルの Incus サーバーにインポートしたい場合も、それをネットワークに公開する必要があります。
      その後、ローカルサーバーにアクセスするには IP アドレスとして `127.0.0.1` を指定できます。
      ```

   1. 証明書のフィンガープリントを確認してください。
   1. 認証の方法を選択してください（{ref}`authentication` 参照）。

      たとえば、証明書トークンを選ぶ場合、 Incus サーバーにログオンしてマイグレーションツールを実行中のマシン用のトークンを [`incus config trust add`](incus_config_trust_add.md) で作成してください。
      次に生成されたトークンを、ツールを認証するのに使用してください。
   1. コンテナと仮想マシンのどちらを作成するか選択してください。
      {ref}`containers-and-vms` を参照してください。
   1. 作成するインスタンスの名前を指定してください。
   1. ルートファイルシステム（コンテナの場合）、起動可能なディスク、パーティションまたはイメージファイル（仮想マシンの場合）のパスを指定します。
   1. コンテナの場合、必要に応じてファイルシステムのマウントを追加します。
   1. 仮想マシンの場合、セキュアブートがサポートされているかを指定します。
   1. 任意で、新しいインスタンスを設定します。
      {ref}`プロファイル <profiles>`を指定するか、{ref}`設定オプション <instance-options>`や{ref}`ストレージ <storage>`を変更したり{ref}`ネットワーク <networking>`を設定する設定オプションを直接指定できます。

      あるいは、マイグレーション後に新しいインスタンスを設定することもできます。
   1. マイグレーションの設定が完了したら、マイグレーションプロセスを開始します。

   <details>
   <summary>コンテナにインポートする出力例を見るには展開</summary>

   ```{terminal}
   :input: sudo ./bin.linux.incus-migrate

   Please provide Incus server URL: https://192.0.2.7:8443
   Certificate fingerprint: xxxxxxxxxxxxxxxxx
   ok (y/n)? y

   1) Use a certificate token
   2) Use an existing TLS authentication certificate
   3) Generate a temporary TLS authentication certificate
   Please pick an authentication mechanism above: 1
   Please provide the certificate token: xxxxxxxxxxxxxxxx

   Remote Incus server:
     Hostname: bar
     Version: 5.4

   Would you like to create a container (1) or virtual-machine (2)?: 1
   Name of the new instance: foo
   Please provide the path to a root filesystem: /
   Do you want to add additional filesystem mounts? [default=no]:

   Instance to be created:
     Name: foo
     Project: default
     Type: container
     Source: /

   Additional overrides can be applied at this stage:
   1) Begin the migration with the above configuration
   2) Override profile list
   3) Set additional configuration options
   4) Change instance storage pool or volume size
   5) Change instance network

   Please pick one of the options above [default=1]: 3
   Please specify config keys and values (key=value ...): limits.cpu=2

   Instance to be created:
     Name: foo
     Project: default
     Type: container
     Source: /
     Config:
       limits.cpu: "2"

   Additional overrides can be applied at this stage:
   1) Begin the migration with the above configuration
   2) Override profile list
   3) Set additional configuration options
   4) Change instance storage pool or volume size
   5) Change instance network

   Please pick one of the options above [default=1]: 4
   Please provide the storage pool to use: default
   Do you want to change the storage size? [default=no]: yes
   Please specify the storage size: 20GiB

   Instance to be created:
     Name: foo
     Project: default
     Type: container
     Source: /
     Storage pool: default
     Storage pool size: 20GiB
     Config:
       limits.cpu: "2"

   Additional overrides can be applied at this stage:
   1) Begin the migration with the above configuration
   2) Override profile list
   3) Set additional configuration options
   4) Change instance storage pool or volume size
   5) Change instance network

   Please pick one of the options above [default=1]: 5
   Please specify the network to use for the instance: incusbr0

   Instance to be created:
     Name: foo
     Project: default
     Type: container
     Source: /
     Storage pool: default
     Storage pool size: 20GiB
     Network name: incusbr0
     Config:
       limits.cpu: "2"

   Additional overrides can be applied at this stage:
   1) Begin the migration with the above configuration
   2) Override profile list
   3) Set additional configuration options
   4) Change instance storage pool or volume size
   5) Change instance network

   Please pick one of the options above [default=1]: 1
   Instance foo successfully created
   ```

   </details>
   <details>
   <summary>仮想マシンにインポートする出力例を見るには展開</summary>

   ```{terminal}
   :input: sudo ./bin.linux.incus-migrate

   Please provide Incus server URL: https://192.0.2.7:8443
   Certificate fingerprint: xxxxxxxxxxxxxxxxx
   ok (y/n)? y

   1) Use a certificate token
   2) Use an existing TLS authentication certificate
   3) Generate a temporary TLS authentication certificate
   Please pick an authentication mechanism above: 1
   Please provide the certificate token: xxxxxxxxxxxxxxxx

   Remote Incus server:
     Hostname: bar
     Version: 5.4

   Would you like to create a container (1) or virtual-machine (2)?: 2
   Name of the new instance: foo
   Please provide the path to a root filesystem: ./virtual-machine.img
   Does the VM support UEFI Secure Boot? [default=no]: no

   Instance to be created:
     Name: foo
     Project: default
     Type: virtual-machine
     Source: ./virtual-machine.img
     Config:
       security.secureboot: "false"

   Additional overrides can be applied at this stage:
   1) Begin the migration with the above configuration
   2) Override profile list
   3) Set additional configuration options
   4) Change instance storage pool or volume size
   5) Change instance network

   Please pick one of the options above [default=1]: 3
   Please specify config keys and values (key=value ...): limits.cpu=2

   Instance to be created:
     Name: foo
     Project: default
     Type: virtual-machine
     Source: ./virtual-machine.img
     Config:
       limits.cpu: "2"
       security.secureboot: "false"

   Additional overrides can be applied at this stage:
   1) Begin the migration with the above configuration
   2) Override profile list
   3) Set additional configuration options
   4) Change instance storage pool or volume size
   5) Change instance network

   Please pick one of the options above [default=1]: 4
   Please provide the storage pool to use: default
   Do you want to change the storage size? [default=no]: yes
   Please specify the storage size: 20GiB

   Instance to be created:
     Name: foo
     Project: default
     Type: virtual-machine
     Source: ./virtual-machine.img
     Storage pool: default
     Storage pool size: 20GiB
     Config:
       limits.cpu: "2"
       security.secureboot: "false"

   Additional overrides can be applied at this stage:
   1) Begin the migration with the above configuration
   2) Override profile list
   3) Set additional configuration options
   4) Change instance storage pool or volume size
   5) Change instance network

   Please pick one of the options above [default=1]: 5
   Please specify the network to use for the instance: incusbr0

   Instance to be created:
     Name: foo
     Project: default
     Type: virtual-machine
     Source: ./virtual-machine.img
     Storage pool: default
     Storage pool size: 20GiB
     Network name: incusbr0
     Config:
       limits.cpu: "2"
       security.secureboot: "false"

   Additional overrides can be applied at this stage:
   1) Begin the migration with the above configuration
   2) Override profile list
   3) Set additional configuration options
   4) Change instance storage pool or volume size
   5) Change instance network

   Please pick one of the options above [default=1]: 1
   Instance foo successfully created
   ```

   </details>
1. マイグレーションが完了したら、新しいインスタンスをチェックし、設定を新しい環境にあわせて更新してください。
   通常は、少なくともストレージ設定（`/etc/fstab`）とネットワーク設定を更新する必要があります。
