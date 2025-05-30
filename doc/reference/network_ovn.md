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
