(network-ovn)=
# OVN ネットワーク

<!-- Include start OVN intro -->
{abbr}`OVN (Open Virtual Network)`は仮想ネットワーク抽象化をサポートするソフトウェアで定義されたネットワークシステムです。
あなた自身のプライベートクラウドを構築するのに使用できます。
詳細は[`www.ovn.org`](https://www.ovn.org/)をご参照ください。
<!-- Include end OVN intro -->

`ovn`ネットワークタイプは OVN{abbr}`SDN (software-defined networking)`を使って論理的なネットワークの作成を可能にします。
この種のネットワークは複数の個別のネットワーク内で同じ論理ネットワークのサブネットを使うような検証環境やマルチテナントの環境で便利です。

Incus の OVN ネットワークはより広いネットワークへの外向きのアクセスを可能にするため既存の管理された{ref}`network-bridge`や{ref}`network-physical`に接続できます。
デフォルトでは、OVN 論理ネットワークからのすべての接続はアップリンクのネットワークによって割り当てられた IP に NAT されます。

OVN ネットワークをセットアップする基本的な手順については{ref}`network-ovn-setup`をご参照ください。

% Include content from [network_bridge.md](network_bridge.md)
```{include} network_bridge.md
    :start-after: <!-- Include start MAC identifier note -->
    :end-before: <!-- Include end MAC identifier note -->
```

(network-ovn-options)=
## 設定オプション

`ovn`ネットワークタイプでは現在以下の設定キーNamespace がサポートされています:

- `bridge` （L2 インターフェースの設定）
- `dns` （DNS サーバーと名前解決の設定）
- `ipv4` （L3 IPv4 設定）
- `ipv6` （L3 IPv6 設定）
- `security` （ネットワーク ACL 設定）
- `user` （key/value の自由形式のユーザーメタデータ）

```{note}
{{note_ip_addresses_CIDR}}
```

`ovn` ネットワークタイプには以下の設定オプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_ovn-common start -->
    :end-before: <!-- config group network_ovn-common end -->
```

```{note}
`bridge.external_interfaces`オプションは不足しているVLANインターフェースを作成するための拡張フォーマットをサポートします。
拡張フォーマットは`<interfaceName>/<parentInterfaceName>/<vlanId>`です。
外部インターフェースが拡張フォーマットでリストに追加されると、ネットワークの作成時にシステムがインターフェースを自動で作成し、ネットワークの終了時に削除します。システムは`<interfaceName>`が存在しないかを確認します。インターフェース名が別の親やVLAN IDで使用中の場合、あるいはインターフェースの作成に失敗する場合、システムはエラーメッセージを表示して作成前の状態に戻します。
```

(network-ovn-child)=
## 子ネットワーク

OVNネットワークでは同じプロジェクト内の別のOVNネットワークを指す`parent`を指定して作成できます。
自身で論理ルーターを作成する代わりに、子は自身の論理スイッチとサブネットを親の論理ルーターに取り付けます。

    incus network create net1 --type=ovn network=UPLINK ipv4.address=192.0.2.1/24
    incus network create net2 --type=ovn parent=net1 ipv4.address=198.51.100.1/24
    incus launch images:debian/13 c1 --network net2

これにより複数の内部サブネットを単一の論理ルーターでルーティングしアップリンクを共有できます。

子ネットワークは自身のスイッチ、サブネット、DHCP、DNSレコード、ACLとインスタンスポートを維持します。
自身のアップリンクは持たず、親のルーターの外部ポートを経由して外部に通信し、親とは独立してNATを有効にでき、あるサブネットは変換しつつ別のサブネットは同じルーターでネイティブにルーティングできます。
子は`ipv4.nat.address`と`ipv6.nat.address`を設定し、共有のルーターの外部アドレスではなく自身のアドレスに変換することもできます。

子ネットワークには以下の内容が当てはまります：

- 論理ルーターを共有するインスタンスは、すべてがそれでルーティングされるため、デフォルトで互いに通信可能です。
  これを制限するには{ref}`network-acls`を使ってください。
- ACLルールでは同じルーターの別のネットワークのトラフィックは`@internal`ではなく`@external`にマッチます。これは`@internal`はルールが適用されるネットワークのアドレスのみをカバーするからです。
- アップリンク、外部ポート、シャーシグループとあらゆるネットワークピアーは親に属します。
  子は`network`、`parent`（ネットワークは1レベルの深さでしかネストできません）、`bridge.hwaddr`、`bridge.external_interfaces`、`bridge.multicast_relay`、あらゆる`tunnel.*`オプションを設定できず、ピアリングに参加できません。
  親のピアリングはこのサブネットもルーティングし、そのピアリングを参照するACLルールはそのトラフィックにマッチします。
- 子のサブネットは親のサブネットや同じ親の他の子のサブネットと重なってはいけません。
- `parent`はネットワークの作成時にのみ設定できます。
- 親のネットワークは子がいる間はリネームや削除できません。

(network-ovn-features)=
## サポートされている機能

`ovn`ネットワークタイプでは以下の機能がサポートされています:

- {ref}`network-acls`
- {ref}`network-forwards`
- {ref}`network-integrations`
- {ref}`network-zones`
- {ref}`network-ovn-peers`
- {ref}`network-load-balancers`

```{toctree}
:maxdepth: 1
:hidden:

OVNのセットアップ </howto/network_ovn_setup>
ルーティング関係を作成 </howto/network_ovn_peers>
ネットワークロードバランサーを設定 </howto/network_load_balancers>
```
