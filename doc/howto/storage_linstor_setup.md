(storage-linstor-setup)=
# IncusでLINSTORをセットアップするには

LINSTORクラスターをセットアップしてIncusのストレージプロバイダーとして設定するにはこのガイドに従ってください。

このガイドでは、Ubuntu 24.04を動かしている3つのノード（`server01`, `server02`, `server03`）をセットアップします。このすべてでIncusのインスタンスを動かしLINSTORクラスターにストレージを提供します。ストレージ専用のノードを持ちネットワーク経由でストレージを利用するだけいった他の設定もサポートされています。ストレージの設定にかかわらず、すべてのIncusノードはノードにボリュームをマウントできるようにするためLINSTORサテライトサービスを動かすのがよいです。

また注目すべきは、ここではLINSTORストレージバックエンドとしてLVM Thinを使いますが、通常のLVMとZFSもサポートされています。

1. 3つのマシン上で以下の手順を実行し、必須のLINSTORコンポーネントをインストールします:

   1. LINBIT PPAを追加します:

          sudo add-apt-repository ppa:linbit/linbit-drbd9-stack
          sudo apt update

   1. 必須パッケージをインストールします:

          sudo apt install lvm2 drbd-dkms drbd-utils linstor-satellite

   1. マシンの起動時にLINSTORサテライトサービスが常に開始されるよう自動起動を有効にします:

          sudo systemctl enable --now linstor-satellite

1. 最初のマシン（`server01`）上で以下の手順を実行し、LINSTORコントローラーをセットアップし、LINSTORクラスターを起動します:

   1. LINSTORコントローラーとクライアントパッケージをインストールします:

          sudo apt install linstor-controller linstor-client python3-setuptools

   1. マシンの起動時にLINSTORコントローラーサービスが常に開始されるよう自動起動を有効にします:

          sudo systemctl enable --now linstor-controller

   1. LINSTORクラスターにノードを追加します（`<server_1>`、`<server_2>`、`<server_3>`を対応するマシンのIPアドレスに置き換えて実行します）。この例では`server01`は（コントローラーとサテライトを）組み合わせたノードですが、他の2つはサテライトのみです。

          linstor node create server01 <server_1> --node-type combined
          linstor node create server02 <server_2> --node-type satellite
          linstor node create server03 <server_3> --node-type satellite

   1. すべてのノードがオンラインになり、ノード名がIncusクラスター内のノード名と一致していることを確認します:

      ```{terminal}
      :input: linstor node list
      :scroll:

      ╭─────────────────────────────────────────────────────────────╮
      ┊ Node     ┊ NodeType  ┊ Addresses                   ┊ State  ┊
      ╞═════════════════════════════════════════════════════════════╡
      ┊ server01 ┊ COMBINED  ┊ 10.172.117.211:3366 (PLAIN) ┊ Online ┊
      ┊ server02 ┊ SATELLITE ┊ 10.172.117.35:3366 (PLAIN)  ┊ Online ┊
      ┊ server03 ┊ SATELLITE ┊ 10.172.117.232:3366 (PLAIN) ┊ Online ┊
      ╰─────────────────────────────────────────────────────────────╯
      ```

   1. すべてのノードで必要な機能が使えるようになっていることを確認します。この例の場合、`LVMThin`と`DRBD`ですが、利用可能となっています:

      ```{terminal}
      :input: linstor node info
      :scroll:

      ╭───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
      ┊ Node     ┊ Diskless ┊ LVM ┊ LVMThin ┊ ZFS/Thin ┊ File/Thin ┊ SPDK ┊ EXOS ┊ Remote SPDK ┊ Storage Spaces ┊ Storage Spaces/Thin ┊
      ╞═══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╡
      ┊ server01 ┊ +        ┊ +   ┊ +       ┊ +        ┊ +         ┊ -    ┊ -    ┊ +           ┊ -              ┊ -                   ┊
      ┊ server02 ┊ +        ┊ +   ┊ +       ┊ +        ┊ +         ┊ -    ┊ -    ┊ +           ┊ -              ┊ -                   ┊
      ┊ server03 ┊ +        ┊ +   ┊ +       ┊ +        ┊ +         ┊ -    ┊ -    ┊ +           ┊ -              ┊ -                   ┊
      ╰───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯

      ╭───────────────────────────────────────────────────────────────────────╮
      ┊ Node     ┊ DRBD ┊ LUKS ┊ NVMe ┊ Cache ┊ BCache ┊ WriteCache ┊ Storage ┊
      ╞═══════════════════════════════════════════════════════════════════════╡
      ┊ server01 ┊ +    ┊ -    ┊ -    ┊ +     ┊ -      ┊ +          ┊ +       ┊
      ┊ server02 ┊ +    ┊ -    ┊ -    ┊ +     ┊ -      ┊ +          ┊ +       ┊
      ┊ server03 ┊ +    ┊ -    ┊ -    ┊ +     ┊ -      ┊ +          ┊ +       ┊
      ╰───────────────────────────────────────────────────────────────────────╯
      ```

   1. クラスターにストレージを寄与するそれぞれのサテライトノードでストレージプールを作成します。この例ではすべてのサテライトノードがストレージを寄与しますが、専用のストレージノードを用意するようなセットアップでは、それらのノードでだけストレージプールを作成します。`vgcreate`と`pvcreate`を使って手動でLVMボリュームグループをセットアップしLINSTORにこれらのボリュームグループを使ってストレージプールをセットアップするようにさせることもできますが、`linstor physical-storage create-device-pool`はこのセットアップを手軽に自動化できます。またプールを構成する複数のデバイスを指定することもできますが、この例では各ノードにある単一の`/dev/nvme1n1`デバイスを持っているものとします:

          linstor physical-storage create-device-pool --storage-pool nvme_pool --pool-name nvme_pool lvmthin server01 /dev/nvme1n1
          linstor physical-storage create-device-pool --storage-pool nvme_pool --pool-name nvme_pool lvmthin server02 /dev/nvme1n1
          linstor physical-storage create-device-pool --storage-pool nvme_pool --pool-name nvme_pool lvmthin server03 /dev/nvme1n1

   1. すべてのストレージプールが作成され、サイズが期待どおりであることを確認します:

      ```{terminal}
      :input: linstor storage-pool list
      :scroll:

      ╭────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
      ┊ StoragePool          ┊ Node     ┊ Driver   ┊ PoolName                    ┊ FreeCapacity ┊ TotalCapacity ┊ CanSnapshots ┊ State ┊ SharedName                    ┊
      ╞════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╡
      ┊ DfltDisklessStorPool ┊ server01 ┊ DISKLESS ┊                             ┊              ┊               ┊ False        ┊ Ok    ┊ server01;DfltDisklessStorPool ┊
      ┊ DfltDisklessStorPool ┊ server02 ┊ DISKLESS ┊                             ┊              ┊               ┊ False        ┊ Ok    ┊ server02;DfltDisklessStorPool ┊
      ┊ DfltDisklessStorPool ┊ server03 ┊ DISKLESS ┊                             ┊              ┊               ┊ False        ┊ Ok    ┊ server03;DfltDisklessStorPool ┊
      ┊ nvme_pool            ┊ server01 ┊ LVM_THIN ┊ linstor_nvme_pool/nvme_pool ┊    49.89 GiB ┊     49.89 GiB ┊ True         ┊ Ok    ┊ server01;nvme_pool            ┊
      ┊ nvme_pool            ┊ server02 ┊ LVM_THIN ┊ linstor_nvme_pool/nvme_pool ┊    49.89 GiB ┊     49.89 GiB ┊ True         ┊ Ok    ┊ server02;nvme_pool            ┊
      ┊ nvme_pool            ┊ server03 ┊ LVM_THIN ┊ linstor_nvme_pool/nvme_pool ┊    49.89 GiB ┊     49.89 GiB ┊ True         ┊ Ok    ┊ server03;nvme_pool            ┊
      ╰────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
      ```

   1. IncusがLINSTORコントローラーと通信できるように設定します（`<server_1>`をコントローラーマシンのIPアドレスに置き換えます）:

          incus config set storage.linstor.controller_connection=http://<server_1>:3370

   1. Incusにストレージプールを作成します。このIncusストレージプール上のボリュームにLINSTORが`nvme_pool`ストレージプールを確実に使うように`linstor.resource_group.storage_pool`オプションを指定します。これは複数のLINSTORストレージプールがある場合に特に有効です（例：NVMEドライブに1つとSATA HDDに別の1つ）:

          incus storage create remote linstor --target server01
          incus storage create remote linstor --target server02
          incus storage create remote linstor --target server03
          incus storage create remote linstor linstor.resource_group.storage_pool=nvme_pool

   1. LINSTORリソースグループがIncusにより作成されたことを確認します:

      ```{terminal}
      :input: linstor resource-group list
      :scroll:
      ╭──────────────────────────────────────────────────────────────────────────────────────╮
      ┊ ResourceGroup ┊ SelectFilter              ┊ VlmNrs ┊ Description                     ┊
      ╞══════════════════════════════════════════════════════════════════════════════════════╡
      ┊ DfltRscGrp    ┊ PlaceCount: 2             ┊        ┊                                 ┊
      ╞┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄╡
      ┊ remote        ┊ PlaceCount: 2             ┊        ┊ Resource group managed by Incus ┊
      ┊               ┊ StoragePool(s): nvme_pool ┊        ┊                                 ┊
      ╰──────────────────────────────────────────────────────────────────────────────────────╯
      ```

   1. ストレージをテストするため、いくつかボリュームとインスタンスを作成します:

          incus launch images:ubuntu/24.04 c1 --storage remote
          incus storage volume create remote fsvol
          incus storage volume attach remote fsvol c1 /mnt

          incus launch images:ubuntu/24.04 v1 --storage remote --vm -c migration.stateful=true
          incus storage volume create remote vol --type block size=42GiB
          incus storage volume attach remote vol v1

   1. Incusで作成されたリソースがLINSTORでどのように見えているかを確認します:

      ```{terminal}
      :input: linstor resource-definition list --show-props Aux/Incus/name Aux/Incus/type Aux/Incus/content-type
      :scroll:
      ╭─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
      ┊ ResourceName                                  ┊ Port ┊ ResourceGroup ┊ Layers       ┊ State ┊ Aux/Incus/name                                                                ┊ Aux/Incus/type   ┊ Aux/Incus/content-type ┊
      ╞═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╡
      ┊ incus-volume-1cb987892f6748299a7f894a483e4e7e ┊ 7004 ┊ remote        ┊ DRBD,STORAGE ┊ ok    ┊ incus-volume-v1                                                               ┊ virtual-machines ┊ block                  ┊
      ┊ incus-volume-5b680bf0dd6f4b39b784c1c151dd510c ┊ 7002 ┊ remote        ┊ DRBD,STORAGE ┊ ok    ┊ incus-volume-default_fsvol                                                    ┊ custom           ┊ filesystem             ┊
      ┊ incus-volume-5d7ee1b9c5224f73b3dd3c3a4ff46fed ┊ 7000 ┊ remote        ┊ DRBD,STORAGE ┊ ok    ┊ incus-volume-198e0b3f6b3685418d9c21b58445686f939596b1fccd8e295191fe515d1ab32c ┊ images           ┊ filesystem             ┊
      ┊ incus-volume-9f7ed7091da346e2b7c764348ffada54 ┊ 7001 ┊ remote        ┊ DRBD,STORAGE ┊ ok    ┊ incus-volume-c1                                                               ┊ containers       ┊ filesystem             ┊
      ┊ incus-volume-10991980d449418b9b8714b769f030d7 ┊ 7005 ┊ remote        ┊ DRBD,STORAGE ┊ ok    ┊ incus-volume-default_vol                                                      ┊ custom           ┊ block                  ┊
      ┊ incus-volume-af0e3529ad514b7b89c7a3a9b8b718ff ┊ 7003 ┊ remote        ┊ DRBD,STORAGE ┊ ok    ┊ incus-volume-dfc28af5f731668509b897ce7eb30d07c5bfe50502da4b2f19421a8a0b05137a ┊ images           ┊ block                  ┊
      ╰─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
      ```
