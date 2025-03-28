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

キー                                   | 型      | 条件                   | デフォルト                                                                        | 説明
:--                                    | :--     | :--                    | :--                                                                               | :--
`network`                              | string  | -                      | -                                                                                 | 外部ネットワークへのアクセスに使うアップリンクのネットワークまたは``none`で隔離されたままにする
`bridge.external_interfaces`           | string  | -                      | -                                                                                 | ブリッジに含む未設定のネットワークインタフェースのカンマ区切りリスト
`bridge.hwaddr`                        | string  | -                      | -                                                                                 | ブリッジのMACアドレス
`bridge.mtu`                           | integer | -                      | `1442`                                                                            | ブリッジのMTU（デフォルトではホストからホストへのGeneveトンネルを許可します）
`dns.nameservers`                      | string  | -                      | アップリンクDNSサーバー（アップリンクが設定されていない場合はIPv4とIPv6アドレス） | DHCPクライアントへとルーターアドバータイズメント経由で広告するDNSサーバーのIPアドレス。IPv4とIPv6の両方のアドレスがDHCPクライアントに広告され、最初のIPv6アドレスはさらにRA経由のRDNSSとしても広告されます。
`dns.domain`                           | string  | -                      | `incus`                                                                           | DHCPのクライアントに広告しDNSの名前解決に使用するドメイン
`dns.search`                           | string  | -                      | -                                                                                 | 完全なドメインサーチのカンマ区切りリスト（デフォルトは`dns.domain`の値）
`dns.zone.forward`                     | string  | -                      | -                                                                                 | 正引きDNSレコード用のDNSゾーン名のカンマ区切りリスト
`dns.zone.reverse.ipv4`                | string  | -                      | -                                                                                 | IPv4逆引きDNSレコード用のDNSゾーン名
`dns.zone.reverse.ipv6`                | string  | -                      | -                                                                                 | IPv6逆引きDNSレコード用のDNSゾーン名
`ipv4.address`                         | string  | 標準モード             | - （作成時の初期値: `auto`）                                                      | ブリッジのIPv4アドレス（CIDR形式）。IPv4をオフにするには`none`、新しいランダムな未使用のサブネットを生成するには`auto`を指定。
`ipv4.dhcp`                            | bool    | IPv4 アドレス          | `true`                                                                            | DHCPを使ってアドレスを割り当てるかどうか
`ipv4.dhcp.expiry`                     | string  | IPv4 DHCP              | `1h`                                                                              | DHCPリースをいつ期限切れにするか
`ipv4.dhcp.routes`                     | string  | IPv4 DHCP              | -                                                                                 | DHCPオプション121経由で提供する静的ルート、代替のサブネット（CIDR）とゲートウェイアドレスのカンマ区切りリスト（dnsmasqとOVNと同じ形式）
`ipv4.l3only`                          | bool    | IPv4 アドレス          | `false`                                                                           | layer 3 only モード を有効にするかどうか
`ipv4.nat`                             | bool    | IPv4 アドレス          | `false` （`ipv4.address`が`auto`の場合の作成時の初期値: `true`）                  | NATするかどうか
`ipv4.nat.address`                     | string  | IPv4 アドレス          | -                                                                                 | ネットワークからの外向きトラフィックに使用されるソースアドレス（アップリンクに`ovn.ingress_mode=routed`が必要）
`ipv6.address`                         | string  | 標準モード             | - （作成時の初期値: `auto`）                                                      | ブリッジのIPv6アドレス（CIDR形式）。IPv6をオフにするには`none`、新しいランダムな未使用のサブネットを生成するには`auto`を指定。
`ipv6.dhcp`                            | bool    | IPv6 アドレス          | `true`                                                                            | Whether to provide additional network configuration over DHCP
`ipv6.dhcp.stateful`                   | bool    | IPv6 DHCP              | `false`                                                                           | DHCPを使ってアドレスを割り当てるかどうか
`ipv6.l3only`                          | bool    | IPv6 DHCP ステートフル | `false`                                                                           | layer 3 only モード を有効にするかどうか
`ipv6.nat`                             | bool    | IPv6 アドレス          | `false` （`ipv6.address`が`auto`の場合の作成時の初期値: `true`）                  | NATするかどうか
`ipv6.nat.address`                     | string  | IPv6 アドレス          | -                                                                                 | ネットワークからの外向きトラフィックに使用されるソースアドレス（アップリンクに`ovn.ingress_mode=routed`が必要）
`security.acls`                        | string  | -                      | -                                                                                 | このネットワークに接続するNICに適用するネットワークACLのカンマ区切りリスト
`security.acls.default.egress.action`  | string  | `security.acls`        | `reject`                                                                          | どのACLルールにもマッチしない外向きトラフィックに使うアクション
`security.acls.default.egress.logged`  | bool    | `security.acls`        | `false`                                                                           | どのACLルールにもマッチしない外向きトラフィックをログ出力するかどうか
`security.acls.default.ingress.action` | string  | `security.acls`        | `reject`                                                                          | どのACLルールにもマッチしない内向きトラフィックに使うアクション
`security.acls.default.ingress.logged` | bool    | `security.acls`        | `false`                                                                           | どのACLルールにもマッチしない内向きトラフィックをログ出力するかどうか
`user.*`                               | string  | -                      | -                                                                                 | ユーザー指定の自由形式のキー／バリューペア

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
