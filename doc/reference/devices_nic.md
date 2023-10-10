(devices-nic)=
# タイプ: `nic`

```{note}
`nic`デバイスタイプはコンテナと VM の両方でサポートされます。

（`ipvlan` NIC タイプを除いて）NIC はコンテナと VM の両方でホットプラグをサポートします。
```

ネットワークデバイス（*ネットワークインタフェースコントローラー*や*NIC*とも呼びます）はネットワークへの接続を提供します。
Incus はさまざまな異なるタイプのネットワークデバイス（*NICタイプ*）をサポートします。

## `nictype` 対 `network`

インスタンスにネットワークデバイスを追加する際には、追加したいデバイスのタイプを選択するのに 2 つの方法があります。`nictype`プロパティを指定するか`network`プロパティを使うかです。

これらの 2 つのデバイスオプションは相互排他であり、デバイスを作成時にどちらか 1 つのみ指定可能です。
しかし、`network`オプションを指定する際には、`nictype`オプションはネットワークタイプから自動的に導出されることに注意してください。

`nictype`
: nictype`デバイスオプションを使用する際は、Incus に管理されていないネットワークインターフェースを指定できます。
  このため、Incus がネットワークインターフェースを使用するために必要なすべての情報を指定する必要があります。

  この方法を使用する際は、`nictype`オプションはデバイス作成時に指定する必要があり、作成後は変更できません。

`network`
: `network`デバイスオプションを使用する際は、NIC は既存の{ref}`管理されたネットワーク <managed-networks>`にリンクされます。
  この場合、Incus はネットワークについて必要な情報をすべて持っているので、デバイス追加時にはネットワーク名を指定するだけでよいです。

  この方法を使用する際は、`nictype`オプションは Incus が自動的に導出します。
  値は読み取り専用で変更できません。

  ネットワークから継承される他のデバイスオプションは NIC 固有のデバイスオプションの「管理」カラムで「yes」と記載されています。
  `network`の方法を使う場合、NIC のこれらのオプションを直接カスタマイズはできません。

詳細な情報は{ref}`networks`を参照してください。

## 利用可能なNIC

次の NIC は`nictype`か`network`オプションを使って追加できます:

- [`bridged`](nic-bridged): ホスト上に存在する既存のブリッジを使い、ホストのブリッジをインスタンスに接続する仮想デバイスペアを作成します。
- [`macvlan`](nic-macvlan): 既存のネットワークデバイスをベースに MAC アドレスが異なる新しいネットワークデバイスを作成します。
- [`sriov`](nic-sriov): SR-IOV が有効な物理ネットワークデバイスの仮想ファンクション（virtual function）をインスタンスにパススルーします。
- [`physical`](nic-physical): ホストの物理デバイスをインスタンスにパススルーします。
  対象のデバイスはホスト上では見えなくなり、インスタンス内に出現します。

次の NIC は`network`オプションでのみ追加できます:

- [`ovn`](nic-ovn): 既存の OVN ネットワークを使用し、インスタンスが接続する仮想デバイスペアを作成します。

次の NIC は`nictype`オプションでのみ追加できます:

- [`ipvlan`](nic-ipvlan): 既存のネットワークデバイスをベースに MAC アドレスは同じですが IP アドレスが異なる新しいネットワークデバイスを作成します。
- [`p2p`](nic-p2p): 仮想デバイスペアを作成し、片方をインスタンス内に置き、残りの片方をホスト上に残します。
- [`routed`](nic-routed): 仮想デバイスペアを作成し、ホストからインスタンスに繋いで静的ルートとプロキシ ARP/NDP エントリーを作成します。これにより指定された親インターフェースのネットワークにインスタンスが参加できるようになります。

利用可能なデバイスオプションは NIC タイプによって異なり、以下のセクションの表に一覧表示されます。

(nic-bridged)=
### `nictype`: `bridged`

```{note}
この NIC タイプは`nictype`オプションか`network`オプションで選択できます（管理された`bridge`ネットワークの情報については{ref}`network-bridge`参照）。
```

`bridged` NIC はホストの既存のブリッジを使用し、ホストのブリッジをインスタンスに接続するための仮想デバイスのペアを作成します。

#### デバイスオプション

`bridged` タイプの NIC デバイスには以下のデバイスオプションがあります:

キー                      | 型      | デフォルト値       | 管理 | 説明
:--                       | :--     | :--                | :--  | :--
`boot.priority`           | integer | -                  | no   | VMのブート優先度（高いほうが先にブート）
`host_name`               | string  | ランダムに割り当て | no   | ホスト内でのインターフェースの名前
`hwaddr`                  | string  | ランダムに割り当て | no   | 新しいインターフェースのMACアドレス
`ipv4.address`            | string  | -                  | no   | DHCPでインスタンスに割り当てるIPv4アドレス（`security.ipv4_filtering`設定時にすべてのIPv4トラフィックを制限するには`none`と設定可能）
`ipv4.routes`             | string  | -                  | no   | ホスト上でNICに追加するIPv4静的ルートのカンマ区切りリスト
`ipv4.routes.external`    | string  | -                  | no   | NICにルーティングしアップリンクのネットワーク（BGP）で公開するIPv4静的ルートのカンマ区切りリスト
`ipv6.address`            | string  | -                  | no   | DHCPでインスタンスに割り当てるIPv6アドレス（`security.ipv6_filtering`設定時にすべてのIPv6トラフィックを制限するには`none`と設定可能）
`ipv6.routes`             | string  | -                  | no   | ホスト上でNICに追加するIPv6静的ルートのカンマ区切りリスト
`ipv6.routes.external`    | string  | -                  | no   | NICにルーティングしアップリンクのネットワーク（BGP）で公開するIPv6静的ルートのカンマ区切りリスト
`limits.egress`           | string  | -                  | no   | 外向きトラフィックのI/O制限値（さまざまな単位が使用可能、{ref}`instances-limit-units`参照）
`limits.ingress`          | string  | -                  | no   | 内向きトラフィックのI/O制限値（さまざまな単位が使用可能、{ref}`instances-limit-units`参照）
`limits.max`              | string  | -                  | no   | 内向きと外向きの両方のトラフィックI/O制限値（`limits.ingress`と`limits.egress`の両方を設定するのと同じ）
`limits.priority`         | integer | -                  | no   | 外向きトラフィックへの `skb->priority` の値（32-bit 符号なし整数）、カーネルでネットワークパケットに優先度をつけるために kernel queuing discipline （qdisc） によって使用される（この値の効果は特定の qdisc 実装、たとえば、`SKBPRIO` または `QFQ` によって異なる。この値を設定する前に kernel qdisc ドキュメントを参照のこと）
`mtu`                     | integer | 親の MTU           | yes  | 新しいインターフェースのMTU
`name`                    | string  | カーネルが割り当て | no   | インスタンス内でのインターフェースの名前
`network`                 | string  | -                  | no   | （`nictype`を直接設定する代わりに）デバイスをリンクする先の管理されたネットワーク
`parent`                  | string  | -                  | yes  | ホストデバイスの名前（`nictype`を直接設定する場合は必須）
`queue.tx.length`         | integer | -                  | no   | NICの送信キューの長さ
`security.ipv4_filtering` | bool    | `false`            | no   | インスタンスが他のインスタンスのIPv4アドレスになりすますのを防ぐ（これを設定すると`mac_filtering`も有効になります）
`security.ipv6_filtering` | bool    | `false`            | no   | インスタンスが他のインスタンスのIPv6アドレスになりすますのを防ぐ（これを設定すると`mac_filtering`も有効になります）
`security.mac_filtering`  | bool    | `false`            | no   | インスタンスが他のインスタンスのMACアドレスになりすますのを防ぐ
`security.port_isolation` | bool    | `false`            | no   | NICがポート隔離を有効にしたネットワーク内の他のNICと通信するのを防ぐ
`vlan`                    | integer | -                  | no   | タグなしのトラフィックに使用するVLAN ID（デフォルトのVLANからポートを削除するには`none`を指定）
`vlan.tagged`             | integer | -                  | no   | タグありのトラフィックに参加するVLAN IDまたはVLANの範囲のカンマ区切りリスト

(nic-macvlan)=
### `nictype`: `macvlan`

```{note}
この NIC タイプは`nictype`オプションか`network`オプションで選択できます（管理された`macvlan`ネットワークの情報については{ref}`network-macvlan`参照）。
```

`macvlan` NIC は既存の NIC をベースにしますが、MAC アドレスが異なる新しいネットワークデバイスをセットアップします。

`macvlan` NIC を使う場合、Incus ホストとインスタンス間の通信はできません。
ホストとインスタンスの両方がゲートウェイと通信できますが、それらが直接通信はできません。

#### デバイスオプション

`macvlan`タイプの NIC デバイスには以下のデバイスオプションがあります:

キー            | 型      | デフォルト値       | 管理 | 説明
:--             | :--     | :--                | :--  | :--
`boot.priority` | integer | -                  | no   | VMのブート優先度（高いほうが先にブート）
`gvrp`          | bool    | `false`            | no   | GARP VLAN Registration Protocolを使ってVLANを登録する
`hwaddr`        | string  | ランダムに割り当て | no   | 新しいインターフェースのMACアドレス
`mtu`           | integer | 親の MTU           | yes  | 新しいインターフェースのMTU
`name`          | string  | カーネルが割り当て | no   | インスタンス内部でのインターフェース名
`network`       | string  | -                  | no   | （`nictype`を直接設定する代わりに）デバイスをリンクする先の管理されたネットワーク
`parent`        | string  | -                  | yes  | ホストデバイスの名前（`nictype`を直接設定する場合は必須）
`vlan`          | integer | -                  | no   | アタッチ先のVLAN ID

(nic-sriov)=
### `nictype`: `sriov`

```{note}
この NIC タイプは`nictype`オプションか`network`オプションで選択できます（管理された`sriov`ネットワークの情報については{ref}`network-sriov`参照）。
```

`sriov` NIC は SR-IOV を有効にした物理ネットワークデバイスの仮想ファンクションをインスタンスにパススルーします。

SR-IOV を有効にしたネットワークデバイスは一組の仮想ファンクション（VF）をネットワークデバイスの単一の物理ファンクション(PF)に関連付けます。
PF は標準的な PCIe 関数です。
一方、VF はデータの移動に最適化された非常に軽量な PCIe 関数です。
PF のプロパティを変えるのを防ぐため、VF の構成機能は限定されています。

VF はシステムには通常の PCIe デバイスのように見えますので、通常の物理デバイスと全く同じようにインスタンスにパススルーできます。

VF の割り当て
: `sriov`インターフェースタイプは`parent`プロパティを通してシステム上の SR-IOV を有効にしたネットワークデバイスの名前を渡されることを想定しています。
  すると Incus はシステム上の任意の利用可能な VF をチェックします。

  デフォルトでは、Incus は見つけた最初の未使用な VF を割り当てます。
  有効になっているものが 1 つもないか、有効な VF がすべて使用中の場合、サポートされている VF の数を最大に上げて最初の未使用な VF を使用します。
  すべての利用可能な VF が使用中か、カーネルまたはカードが VF の数の増加をサポートしない場合は、Incus はエラーを返します。

  ```{note}
  Incus に特定の VF を使わせたい場合、`sriov` NIC の代わりに`physical` NIC を使用し、`parent`オプションを VF 名に設定してください。
  ```

#### デバイスオプション

`sriov`タイプの NIC デバイスには以下のデバイスオプションがあります:

キー                     | 型      | デフォルト値       | 管理 | 説明
:--                      | :--     | :--                | :--  | :--
`boot.priority`          | integer | -                  | no   | VMのブート優先度（高いほうが先にブート）
`hwaddr`                 | string  | ランダムに割り当て | no   | 新しいインターフェースのMACアドレス
`mtu`                    | integer | カーネルが割り当て | yes  | 新しいインターフェースのMTU
`name`                   | string  | カーネルが割り当て | no   | インスタンス内部でのインターフェース名
`network`                | string  | -                  | no   | （`nictype`を直接設定する代わりに）デバイスをリンクする先の管理されたネットワーク
`parent`                 | string  | -                  | yes  | ホストデバイスの名前（`nictype`を直接設定する場合は必須）
`security.mac_filtering` | bool    | `false`            | no   | インスタンスが他のインスタンスのMACアドレスになりすますのを防ぐ
`vlan`                   | integer | -                  | no   | アタッチ先のVLAN ID

(nic-ovn)=
### `nictype`: `ovn`

```{note}
この NIC タイプは`network`オプションでのみ選択できます（管理された`ovn`ネットワークの情報については{ref}`network-ovn`参照）。
```

`ovn` NIC は既存の OVN ネットワークを使用し、それにインスタンスが接続する仮想デバイスペアを作成します。

(devices-nic-hw-acceleration)=
SR-IOV ハードウェアアクセラレーション
: `acceleration=sriov`を使用するには、Incus ホスト内の Ethernet スイッチデバイスのドライバーモデル（`switchdev`）をサポートする互換性のある SR-IOV 物理 NIC を持っている必要があります。
  Incus は物理 NIC（PF）が`switchdev`モードに設定され、OVN 統合 OVS ブリッジに接続され、1 つ以上の仮想ファンクション（VF）がアクティブになっていることを前提とします。

  これを実現するには、基本的な前提条件となる以下のセットアップ手順に従ってください。

   1. PF と VF をセットアップする:

      1. PF 上でいくつかの VF をアクティベートし（以下の例では`enp9s0f0np0`とし、PCI アドレスは`0000:09:00.0`とします）、アンバインドします。
      1. `switchdev`モードと PF 上の`hw-tc-offload`を有効にします。
      1. VF をリバインドします。

      ```
      echo 4 > /sys/bus/pci/devices/0000:09:00.0/sriov_numvfs
      for i in $(lspci -nnn | grep "Virtual Function" | cut -d' ' -f1); do echo 0000:$i > /sys/bus/pci/drivers/mlx5_core/unbind; done
      devlink dev eswitch set pci/0000:09:00.0 mode switchdev
      ethtool -K enp9s0f0np0 hw-tc-offload on
      for i in $(lspci -nnn | grep "Virtual Function" | cut -d' ' -f1); do echo 0000:$i > /sys/bus/pci/drivers/mlx5_core/bind; done
      ```

   1. ハードウェアオフロードを有効にし、統合ブリッジ（通常`br-int`と呼ばれます）に PF NIC を追加して OVS をセットアップします:

      ```
      ovs-vsctl set open_vswitch . other_config:hw-offload=true
      systemctl restart openvswitch-switch
      ovs-vsctl add-port br-int enp9s0f0np0
      ip link set enp9s0f0np0 up
      ```


VDPA ハードウェアアクセラレーション
: `acceleration=vdpa`を使用するには互換性のある VDPA 物理 NIC が必要です。
  セットアップ手順は SR-IOV ハードウェアアクセラレーションと同様ですが、さらに`vhost_vdpa`モジュールをセットアップし、利用可能な VDPA 管理デバイスがあることを確認する必要があります:

  ```
  modprobe vhost_vdpa && vdpa mgmtdev show
  ```

#### デバイスオプション

`ovn` タイプの NIC デバイスには以下のデバイスオプションがあります:

キー                                   | 型      | デフォルト値       | 管理 | 説明
:--                                    | :--     | :--                | :--  | :--
`acceleration`                         | string  | `none`             | no   | ハードウェアオフローディングを有効にする（`none`か`sriov`か`vdpa`、{ref}`devices-nic-hw-acceleration`参照）
`boot.priority`                        | integer | -                  | no   | VMのブート優先度（高いほうが先にブート）
`host_name`                            | string  | ランダムに割り当て | no   | ホスト内部でのインターフェース名
`hwaddr`                               | string  | ランダムに割り当て | no   | 新しいインターフェースのMACアドレス
`ipv4.address`                         | string  | -                  | no   | DHCPでインスタンスに割り当てるIPv4アドレス
`ipv4.routes`                          | string  | -                  | no   | NICへルーティングするIPv4静的ルートのカンマ区切りリスト
`ipv4.routes.external`                 | string  | -                  | no   | NICへのルーティングとアップリンクネットワークでの公開に使用するIPv4静的ルートのカンマ区切りリスト
`ipv6.address`                         | string  | -                  | no   | DHCPでインスタンスに割り当てるIPv6アドレス
`ipv6.routes`                          | string  | -                  | no   | NICへルーティングするIPv6静的ルートのカンマ区切りリスト
`ipv6.routes.external`                 | string  | -                  | no   | NICへのルーティングとアップリンクネットワークでの公開に使用するIPv6静的ルートのカンマ区切りリスト
`name`                                 | string  | カーネルが割り当て | no   | インスタンス内部でのインターフェース名
`nested`                               | string  | -                  | no   | このNICをどの親NICの下にネストするか（`vlan`も参照）
`network`                              | string  | -                  | yes  | デバイスの接続先の管理されたネットワーク（必須）
`security.acls`                        | string  | -                  | no   | 適用するネットワークACLのカンマ区切りリスト
`security.acls.default.egress.action`  | string  | `reject`           | no   | どのACLルールにもマッチしない外向きトラフィックに使うアクション
`security.acls.default.egress.logged`  | bool    | `false`            | no   | どのACLルールにもマッチしない外向きトラフィックをログ出力するかどうか
`security.acls.default.ingress.action` | string  | `reject`           | no   | どのACLルールにもマッチしない内向きトラフィックに使うアクション
`security.acls.default.ingress.logged` | bool    | `false`            | no   | どのACLルールにもマッチしない内向きトラフィックをログ出力するかどうか
`vlan`                                 | integer | -                  | no   | ネストする際に使用する VLAN ID （`nested`も参照）

(nic-physical)=
### `nictype`: `physical`

```{note}
- この NIC タイプは`nictype`オプションまたは`network`オプションで選択できます（管理された`physical`ネットワークの情報については{ref}`network-physical`参照）。
- それぞれの親デバイスに対して`physical` NIC は1つだけ持つことができます。
```

`physical` NIC はホストからパススルーされるそのままの物理デバイスを提供します。
対象のデバイスはホストから消失し、インスタンス内に出現します（これは各ターゲットデバイスに`physical` NIC は 1 つだけ持つことができることを意味します）。

#### デバイスオプション

`physical`タイプの NIC デバイスには以下のデバイスオプションがあります:

キー            | 型      | デフォルト値       | 管理 | 説明
:--             | :--     | :--                | :--  | :--
`boot.priority` | integer | -                  | no   | VMのブート優先度（高いほうが先にブート）
`gvrp`          | bool    | `false`            | no   | GARP VLAN Registration Protocolを使ってVLANを登録する
`hwaddr`        | string  | ランダムに割り当て | no   | 新しいインターフェースのMACアドレス
`mtu`           | integer | 親の MTU           | no   | 新しいインターフェースのMTU
`name`          | string  | カーネルが割り当て | no   | インスタンス内部でのインターフェース名
`network`       | string  | -                  | no   | デバイスのリンク先（`nictype`を直接指定する代わりに）の管理ネットワーク
`parent`        | string  | -                  | yes  | ホストデバイスの名前（`nictype`を直接指定する場合は必須）
`vlan`          | integer | -                  | no   | アタッチ先のVLAN ID

(nic-ipvlan)=
### `nictype`: `ipvlan`

```{note}
- この NIC タイプはコンテナのみで利用でき、仮想マシンでは利用できません。
- この NIC タイプは`nictype`オプションでのみ選択できます。
- この NIC タイプはホットプラグをサポートしません。
```

`ipvlan` NIC は既存のネットワークデバイスを元に、同じ MAC アドレスですが IP アドレスは異なるような新しいネットワークデバイスをセットアップします。

`ipvlan` NIC を使う場合、Incus ホストとインスタンス間の通信はできません。
ホストとインスタンスの両方がゲートウェイと通信できますが、それらが直接通信はできません。

Incus は現状 L2 と L3S モードで IPVLAN をサポートします。
このモードでは、ゲートウェイは Incus により自動的に設定されますが、コンテナが起動する前に`ipv4.address`と`ipv6.address`の設定の 1 つあるいは両方を使うことにより IP アドレスを手動で指定する必要があります。

DNS
: ネームサーバーは自動的には設定されないので、コンテナ内部で設定する必要があります。
  このためには、以下の`sysctl`の設定をしてください:

   - IPv4 アドレスを使用する場合:

     ```
     net.ipv4.conf.<parent>.forwarding=1
     ```

   - IPv6 アドレスを使用する場合:

     ```
     net.ipv6.conf.<parent>.forwarding=1
     net.ipv6.conf.<parent>.proxy_ndp=1
     ```

#### デバイスオプション

`ipvlan`タイプの NIC デバイスには以下のデバイスオプションがあります:

キー              | 型      | デフォルト値                 | 説明
:--               | :--     | :--                          | :--
`gvrp`            | bool    | `false`                      | GARP VLAN Registration Protocolを使ってVLANを登録する
`hwaddr`          | string  | ランダムに割り当て           | 新しいインターフェースのMACアドレス
`ipv4.address`    | string  | -                            | インスタンスに追加するIPv4静的アドレスのカンマ区切りリスト（`l2`モードでは、CIDR形式か`/24`のサブネットの単一アドレスで指定可能）
`ipv4.gateway`    | string  | `auto` （`l3s`）, - （`l2`） | `l3s`モードでは、デフォルトIPv4ゲートウェイを自動的に追加するかどうか（`auto`か`none`を指定可能）。`l2`モードでは、ゲートウェイのIPv4アドレス
`ipv4.host_table` | integer | -                            | （メインのルーティングテーブルに加えて）IPv4の静的ルートを追加する先のカスタムポリシー・ルーティングテーブルID
`ipv6.address`    | string  | -                            | インスタンスに追加するIPv6静的アドレスのカンマ区切りリスト（`l2`モードでは、CIDR 形式か`/64`のサブネットの単一アドレスで指定可能）
`ipv6.gateway`    | string  | `auto` （`l3s`）, - （`l2`） | `l3s`モードでは、デフォルトIPv6ゲートウェイを自動的に追加するかどうか（`auto`か`none`を指定可能）。`l2`モードで、はゲートウェイのIPv6アドレス
`ipv6.host_table` | integer | -                            | （メインのルーティングテーブルに加えて）IPv6の静的ルートを追加する先のカスタムポリシー・ルーティングテーブルID
`mode`            | string  | `l3s`                        | IPVLANのモード（`l2`か`l3s`のいずれか）
`mtu`             | integer | 親の MTU                     | 新しいインターフェースのMTU
`name`            | string  | カーネルが割り当て           | インスタンス内部でのインターフェース名
`parent`          | string  | -                            | ホストデバイスの名前(必須)
`vlan`            | integer | -                            | アタッチ先のVLAN ID

(nic-p2p)=
### `nictype`: `p2p`

```{note}
この NIC タイプは`nictype`オプションでのみ選択できます。
```

`p2p` NIC は仮想デバイスペアを作成し、片方はインスタンス内に配置し、もう片方はホストに残します。

#### デバイスオプション

`p2p`タイプの NIC デバイスには以下のデバイスオプションがあります:

キー              | 型      | デフォルト値       | 説明
:--               | :--     | :--                | :--
`boot.priority`   | integer | -                  | VMのブート優先度 （高いほうが先にブート）
`host_name`       | string  | ランダムに割り当て | ホスト内でのインターフェースの名前
`hwaddr`          | string  | ランダムに割り当て | 新しいインターフェースのMACアドレス
`ipv4.routes`     | string  | -                  | ホスト上でNICに追加するIPv4静的ルートのカンマ区切りリスト
`ipv6.routes`     | string  | -                  | ホスト上でNICに追加するIPv6静的ルートのカンマ区切りリスト
`limits.egress`   | string  | -                  | 外向きトラフィックのI/O制限値（さまざまな単位が使用可能、{ref}`instances-limit-units`参照）
`limits.ingress`  | string  | -                  | 内向きトラフィックのI/O制限値（さまざまな単位が使用可能、{ref}`instances-limit-units`参照）
`limits.max`      | string  | -                  | 内向きと外向きの両方のトラフィックI/O制限値（`limits.ingress`と`limits.egress`の両方を設定するのと同じ）
`limits.priority` | integer | -                  | 外向きトラフィックへの `skb->priority` の値（32-bit 符号なし整数）、カーネルでネットワークパケットに優先度をつけるために kernel queuing discipline （qdisc） によって使用される（この値の効果は特定の qdisc 実装、たとえば、`SKBPRIO` または `QFQ` によって異なる。この値を設定する前に kernel qdisc ドキュメントを参照のこと）
`mtu`             | integer | カーネルが割り当て | 新しいインターフェースのMTU
`name`            | string  | カーネルが割り当て | インスタンス内部でのインターフェース名
`queue.tx.length` | integer | -                  | NIC の送信キューの長さ

(nic-routed)=
### `nictype`: `routed`

```{note}
この NIC タイプは`nictype`オプションでのみ選択できます。
```

`routed` NIC タイプはホストをインスタンスに接続する仮想デバイスペアを作成し、インスタンスが指定された親インターフェースのネットワークに参加できるように、静的ルートとプロキシ ARP/NDP エントリをセットアップします。
コンテナでは仮想イーサネットデバイスペアを使用し、VM では TAP デバイスを使用します。

この NIC タイプは運用上は IPVLAN に似ていて、ブリッジを設定することなくホストの MAC アドレスを共用して、インスタンスが外部ネットワークに参加できるようにします。
しかし、カーネルに IPVLAN サポートを必要としないことと、ホストとインスタンスが互いに通信できることが`ipvlan`とは異なります。

この NIC タイプは`netfilter`のルールを尊重し、ホストのルーティングテーブルを使ってパケットをルーティングしますので、ホストが複数のネットワークに接続している場合に役立ちます。

IP アドレス、ゲートウェイ、ルーティング
: インスタンスが起動する前に IP アドレスを（`ipv4.address`と`ipv6.address`の設定のいずれかあるいは両方を使って）手動で指定する必要があります。

  コンテナでは、NIC はホスト上に下記のリンクローカルゲートウェイ IP アドレスを設定し、それらをコンテナの NIC インターフェースのデフォルトゲートウェイに設定します:

      169.254.0.1
      fe80::1

  VM では、ゲートウェイは手動か`cloud-init`（{ref}`ハウツーガイド <instances-routed-nic-vm>`参照）のような仕組みを使って設定する必要があります。

  ```{note}
  お使いのコンテナイメージがインタフェースに対して DHCP を使うように設定されている場合、上記の自動的に追加される設定は削除される可能性が高いです。
  この場合、IP アドレスとゲートウェイを手動か`cloud-init`のような仕組みを使って設定する必要があります。
  ```

  この NIC タイプはインスタンスの IP アドレスすべてをインスタンスの`veth`インターフェースに向ける静的ルートをホスト上に設定します。

複数の IP アドレス
: それぞれの NIC デバイスに複数の IP アドレスを追加できます。

  しかし、代わりに複数の`routed` NIC インターフェースを使うほうが望ましいかもしれません。
  この場合、`ipv4.gateway`と`ipv6.gateway`の値を`none`に設定し、後続のインターフェースがデフォルトゲートウェイの衝突を避けるようにします。
  さらに、これらの後続のインターフェースに`ipv4.host_address`と`ipv6.host_address`を使って異なるホスト側のアドレスを指定することを検討してください。

親のインターフェース
: この NIC は`parent`のネットワークインターフェースのセットがあってもなくても利用できます。

: `parent`ネットワークインターフェースのセットがある場合、インスタンスの IP のプロキシ ARP/NDP エントリが親のインターフェースに追加され、インスタンスが親のインターフェースのネットワークにレイヤ 2 で参加できるようにします。
: これを有効にするには以下のネットワーク設定を `sysctl` でホストに適用する必要があります:

   - IPv4 アドレスを使用する場合:

     ```
     net.ipv4.conf.<parent>.forwarding=1
     ```

   - IPv6 アドレスを使用する場合:

     ```
     net.ipv6.conf.all.forwarding=1
     net.ipv6.conf.<parent>.forwarding=1
     net.ipv6.conf.all.proxy_ndp=1
     net.ipv6.conf.<parent>.proxy_ndp=1
     ```

#### デバイスオプション

`routed`タイプの NIC デバイスには以下のデバイスオプションがあります:

キー                  | 型      | デフォルト値       | 説明
:--                   | :--     | :--                | :--
`gvrp`                | bool    | `false`            | GARP VLAN Registration Protocolを使ってVLANを登録する
`host_name`           | string  | ランダムに割り当て | ホスト内でのインターフェース名
`hwaddr`              | string  | ランダムに割り当て | 新しいインターフェースのMACアドレス
`ipv4.address`        | string  | -                  | インスタンスに追加するIPv4静的アドレスのカンマ区切りリスト
`ipv4.gateway`        | string  | `auto`             | 自動的にIPv4デフォルトゲートウェイを追加するかどうか（`auto`か`none`を指定可能）
`ipv4.host_address`   | string  | `169.254.0.1`      | ホスト側の`veth`インターフェースに追加するIPv4アドレス
`ipv4.host_table`     | integer | -                  | （メインのルーティングテーブルに加えて）IPv4の静的ルートを追加する先のカスタムポリシー・ルーティングテーブルID
`ipv4.neighbor_probe` | bool    | `true`             | IPアドレスが利用可能か知るために親のネットワークを調べるかどうか
`ipv4.routes`         | string  | -                  | ホスト上でNICに追加するIPv4静的ルートのカンマ区切りリスト（L2 ARP/NDPプロキシを除く）
`ipv6.address`        | string  | -                  | インスタンスに追加するIPv6静的アドレスのカンマ区切りリスト
`ipv6.gateway`        | string  | `auto`             | 自動的にIPv6のデフォルトゲートウェイを追加するかどうか（`auto`か`none`を指定可能）
`ipv6.host_address`   | string  | `fe80::1`          | ホスト側の`veth`インターフェースに追加するIPv6アドレス
`ipv6.host_table`     | integer | -                  | （メインのルーティングテーブルに加えて）IPv6の静的ルートを追加する先のカスタムポリシー・ルーティングテーブルID
`ipv6.neighbor_probe` | bool    | `true`             | IPアドレスが利用可能か知るために親のネットワークを調べるかどうか
`ipv6.routes`         | string  | -                  | ホスト上でNICに追加するIPv6静的ルートのカンマ区切りリスト（L2 ARP/NDPプロキシを除く）
`limits.egress`       | string  | -                  | 内向きトラフィックに対するbit/sでのI/O制限値（さまざまな単位をサポート、{ref}`instances-limit-units`参照）
`limits.ingress`      | string  | -                  | 外向きトラフィックに対するbit/sでのI/O制限値（さまざまな単位をサポート、{ref}`instances-limit-units`参照）
`limits.max`          | string  | -                  | 内向きと外向き両方のトラフィックのI/O 制限値（`limits.ingress`と`limits.egress`の両方を設定するのと同じ）
`limits.priority`     | integer | -                  | 外向きトラフィックへの `skb->priority` の値（32-bit 符号なし整数）、カーネルでネットワークパケットに優先度をつけるために kernel queuing discipline （qdisc） によって使用される（この値の効果は特定の qdisc 実装、たとえば、`SKBPRIO` または `QFQ` によって異なる。この値を設定する前に kernel qdisc ドキュメントを参照のこと）
`mtu`                 | integer | 親の MTU           | 新しいインターフェースのMTU
`name`                | string  | カーネルが割り当て | インスタンス内でのインターフェース名
`parent`              | string  | -                  | インスタンスが参加するホストデバイス名
`queue.tx.length`     | integer | -                  | NICの送信キューの長さ
`vlan`                | integer | -                  | アタッチ先の VLAN ID

## `bridge`、`macvlan`、`ipvlan`を使った物理ネットワークへの接続

`bridged`、`macvlan`、`ipvlan`インターフェースタイプのいずれも、既存の物理ネットワークへ接続するために使用できます。

`macvlan`は、物理 NIC を効率的に分岐できます。つまり、物理 NIC からインスタンスで使える第 2 のインターフェースを取得できます。
この方法はブリッジデバイスと仮想イーサネットデバイスペアの作成を不要にしますし、通常はブリッジよりも良いパフォーマンスが得られます。

`macvlan`の欠点は、`macvlan`はインスタンス自身と外部との間で通信はできますが、親デバイスとは通信できないことです。
つまりインスタンスとホストが通信する必要がある場合は`macvlan`は使えません。

そのような場合は、`bridge`デバイスを選ぶのが良いでしょう。
`macvlan`では使えない MAC フィルタリングと I/O 制限も使えます。

`ipvlan`は`macvlan`と同様ですが、フォークされたデバイスが静的に割り当てられた IP アドレスを持ち、ネットワーク上の親の MAC アドレスを受け継ぐ点が異なります。
