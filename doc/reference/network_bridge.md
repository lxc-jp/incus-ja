(network-bridge)=
# ブリッジネットワーク

Incus でのネットワークの設定タイプの 1 つとして、Incus はネットワークブリッジの作成と管理をサポートしています。
<!-- Include start bridge intro -->
ネットワークブリッジはインスタンス NIC が接続できる仮想的な L2 イーサネットスイッチを作成し、インスタンスが他のインスタンスやホストと通信できるようにします。
Incus のブリッジは下層のネイティブな Linux のブリッジと Open vSwitch を利用できます。
<!-- Include end bridge intro -->

`bridge`ネットワークはそれを利用する複数のインスタンスを接続する L2 ブリッジを作成しそれらのインスタンスを単一の L2 ネットワークセグメントにします。
Incus で作成されたブリッジは"managed"です。
つまり、ブリッジインターフェース自体を作成するのに加えて、Incus さらに DHCP、IPv6 ルート広告と DNS サービスを提供するローカルの`dnsmasq`プロセスをセットアップします。
デフォルトではブリッジに対して NAT も行います。

Incus ブリッジネットワークでファイアウォールを設定するための手順については{ref}`network-bridge-firewall`を参照してください。

<!-- Include start MAC identifier note -->

```{note}
静的な DHCP 割当は MAC アドレスを DHCP 識別子として使用するクライアントに依存します。
この方法はインスタンスをコピーする際に衝突するリースを回避し、静的に割り当てられたリースが正しく動くようにします。
```

<!-- Include end MAC identifier note -->

## IPv6プリフィクスサイズ

ブリッジネットワークで IPv6 を使用する場合、64 のプリフィクスサイズを使用するべきです。

より大きなサブネット（つまり 64 より小さいプリフィクスを使用する）も正常に動くはずですが、通常それらは{abbr}`SLAAC (Stateless Address Auto-configuration)`には役立ちません。

より小さなサブネットも（IPv6 の割当にはステートフル DHCPv6 を使用する場合）理論上は可能ですが、`dnsmasq`に適切にサポートされていないので問題が起きるかもしれません。より小さなサブネットを作らなければならない場合は、静的割当を使うか別のルータ広告デーモンを使用してください。

(network-bridge-options)=
## 設定オプション

`bridge`ネットワークタイプでは現在以下の設定キーNamespace がサポートされています:

- `bgp` (BGP ピア設定)
- `bridge` (L2 インターフェースの設定)
- `dns` (DNS サーバーと名前解決の設定)
- `ipv4` (L3 IPv4 設定)
- `ipv6` (L3 IPv6 設定)
- `security` (ネットワーク ACL 設定)
- `raw` (raw の設定のファイルの内容)
- `tunnel` (ホスト間のトンネリングの設定)
- `user` (key/value の自由形式のユーザーメタデータ)

```{note}
{{note_ip_addresses_CIDR}}
```

`bridge`ネットワークタイプには以下の設定オプションがあります:

キー                                   | 型      | 条件               | デフォルト                                                      | 説明
:--                                    | :--     | :--                | :--                                                             | :--
`bgp.peers.NAME.address`               | string  | BGPサーバー        | -                                                               | ピアのアドレス（IPv4かIPv6）
`bgp.peers.NAME.asn`                   | integer | BGPサーバー        | -                                                               | ピアのAS番号
`bgp.peers.NAME.password`              | string  | BGPサーバー        | - （パスワード無し）                                            | ピアのセッションパスワード（省略可能）
`bgp.peers.NAME.holdtime`              | integer | BGPサーバー        | `180`                                                           | ピアセッションホールドタイム（秒で指定、省略可能）
`bgp.ipv4.nexthop`                     | string  | BGPサーバー        | ローカルアドレス                                                | 広告されたプリフィクスのnext-hopをオーバーライド
`bgp.ipv6.nexthop`                     | string  | BGPサーバー        | ローカルアドレス                                                | 広告されたプリフィクスのnext-hopをオーバーライド
`bridge.driver`                        | string  | -                  | `native`                                                        | ブリッジのドライバー: `native`か`openvswitch`
`bridge.external_interfaces`           | string  | -                  | -                                                               | ブリッジに含める未設定のネットワークインターフェースのカンマ区切りリスト
`bridge.hwaddr`                        | string  | -                  | -                                                               | ブリッジのMACアドレス
`bridge.mtu`                           | integer | -                  | `1500`                                                          | ブリッジのMTU（tunnel使用時はデフォルト値は変わる）
`dns.domain`                           | string  | -                  | `incus`                                                         | DHCPのクライアントに広告しDNSの名前解決に使用するドメイン
`dns.mode`                             | string  | -                  | `managed`                                                       | DNSの登録モード: `none`はDNSレコード無し、`managed`はIncusが静的レコードを生成、`dynamic`はクライアントがレコードを生成
`dns.search`                           | string  | -                  | -                                                               | 完全なドメインサーチのカンマ区切りリスト、デフォルトは`dns.domain`の値
`dns.zone.forward`                     | string  | -                  | `managed`                                                       | 正引きDNSレコード用のDNSゾーン名のカンマ区切りリスト
`dns.zone.reverse.ipv4`                | string  | -                  | `managed`                                                       | IPv4逆引きDNSレコード用のDNSゾーン名
`dns.zone.reverse.ipv6`                | string  | -                  | `managed`                                                       | IPv6逆引きDNSレコード用のDNSゾーン名
`ipv4.address`                         | string  | 標準モード         | - （作成時の初期値: `auto`）                                    | ブリッジのIPv4アドレス（CIDR形式）（IPv4をオフにするには`none`、新しいランダムな未使用のサブネットを生成するには`auto`を指定）
`ipv4.dhcp`                            | bool    | IPv4 アドレス      | `true`                                                          | DHCPを使ってアドレスを割り当てるかどうか
`ipv4.dhcp.expiry`                     | string  | IPv4 DHCP          | `1h`                                                            | DHCPリースの有効期限
`ipv4.dhcp.gateway`                    | string  | IPv4 DHCP          | IPv4 アドレス                                                   | サブネットのゲートウェイのアドレス
`ipv4.dhcp.ranges`                     | string  | IPv4 DHCP          | すべてのアドレス                                                | DHCPに使用するIPv4の範囲（開始-終了の形式）のカンマ区切りリスト
`ipv4.firewall`                        | bool    | IPv4 アドレス      | `true`                                                          | このネットワークに対するファイアウォールのフィルタリングルールを生成するかどうか
`ipv4.nat`                             | bool    | IPv4 アドレス      | `false`（`ipv4.address`が`auto`の場合の作成時の初期値: `true`） | NATにするかどうか
`ipv4.nat.address`                     | string  | IPv4 アドレス      | -                                                               | ブリッジからの送信時に使うソースアドレス
`ipv4.nat.order`                       | string  | IPv4 アドレス      | `before`                                                        | 必要なNATのルールを既存のルールの前に追加するか後に追加するか
`ipv4.ovn.ranges`                      | string  | -                  | -                                                               | 子供のOVNネットワークルーターに使用するIPv4アドレスの範囲（開始-終了の形式）のカンマ区切りリスト
`ipv4.routes`                          | string  | IPv4 アドレス      | -                                                               | ブリッジへルーティングする追加のIPv4 CIDRサブネットのカンマ区切りリスト
`ipv4.routing`                         | bool    | IPv4 アドレス      | `true`                                                          | ブリッジの内外にトラフィックをルーティングするかどうか
`ipv6.address`                         | string  | 標準モード         | - （作成時の初期値: `auto`）                                    | ブリッジのIPv6アドレス（CIDR形式）（IPv6をオフにするには`none`、新しいランダムな未使用のサブネットを生成するには`auto`を指定）
`ipv6.dhcp`                            | bool    | IPv6 アドレス      | `true`                                                          | DHCP上で追加のネットワーク設定を提供するかどうか
`ipv6.dhcp.expiry`                     | string  | IPv6 DHCP          | `1h`                                                            | DHCPリースの有効期限
`ipv6.dhcp.ranges`                     | string  | IPv6 stateful DHCP | すべてのアドレス                                                | DHCPに使用するIPv6の範囲（開始-終了の形式）のカンマ区切りリスト
`ipv6.dhcp.stateful`                   | bool    | IPv6 DHCP          | `false`                                                         | DHCPを使ってアドレスを割り当てるかどうか
`ipv6.firewall`                        | bool    | IPv6 アドレス      | `true`                                                          | このネットワークに対するファイアウォールのフィルタリングルールを生成するかどうか
`ipv6.nat`                             | bool    | IPv6 アドレス      | `false`（`ipv6.address`が`auto`の場合の作成時の初期値: `true`） | NATにするかどうか
`ipv6.nat.address`                     | string  | IPv6 アドレス      | -                                                               | ブリッジからの送信時に使うソースアドレス
`ipv6.nat.order`                       | string  | IPv6 アドレス      | `before`                                                        | 必要なNATのルールを既存のルールの前に追加するか後に追加するか
`ipv6.ovn.ranges`                      | string  | -                  | -                                                               | 子供のOVNネットワークルーターに使用するIPv6アドレスの範囲（開始-終了の形式）のカンマ区切りリスト
`ipv6.routes`                          | string  | IPv6 アドレス      | -                                                               | ブリッジへルーティングする追加のIPv4 CIDRサブネットのカンマ区切りリスト
`ipv6.routing`                         | bool    | IPv6 アドレス      | `true`                                                          | ブリッジの内外にトラフィックをルーティングするかどうか
`raw.dnsmasq`                          | string  | -                  | -                                                               | 設定に追加する`dnsmasq`の設定ファイル
`security.acls`                        | string  | -                  | -                                                               | このネットワークに接続されたNICに適用するカンマ区切りのネットワークACL（{ref}`network-acls-bridge-limitations`参照）
`security.acls.default.egress.action`  | string  | `security.acls`    | `reject`                                                        | どのACLルールにもマッチしない外向きトラフィックに使うアクション
`security.acls.default.egress.logged`  | bool    | `security.acls`    | `false`                                                         | どのACLルールにもマッチしない外向きトラフィックをログ出力するかどうか
`security.acls.default.ingress.action` | string  | `security.acls`    | `reject`                                                        | どのACLルールにもマッチしない内向きトラフィックに使うアクション
`security.acls.default.ingress.logged` | bool    | `security.acls`    | `false`                                                         | どのACLルールにもマッチしない内向きトラフィックをログ出力するかどうか
`tunnel.NAME.group`                    | string  | `vxlan`            | `239.0.0.1`                                                     | `vxlan`のマルチキャスト設定（localとremoteが未設定の場合に使われます）
`tunnel.NAME.id`                       | integer | `vxlan`            | `0`                                                             | `vxlan`トンネルに使用するトンネルID
`tunnel.NAME.interface`                | string  | `vxlan`            | -                                                               | トンネルに使用するホスト・インターフェース
`tunnel.NAME.local`                    | string  | `gre` or `vxlan`   | -                                                               | トンネルに使用するローカルアドレス（マルチキャスト`vxlan`の場合は不要）
`tunnel.NAME.port`                     | integer | `vxlan`            | `0`                                                             | `vxlan`トンネルに使用するポート
`tunnel.NAME.protocol`                 | string  | 標準モード         | -                                                               | トンネリングのプロトコル: `vxlan`か`gre`
`tunnel.NAME.remote`                   | string  | `gre` or `vxlan`   | -                                                               | トンネルに使用するリモートアドレス（マルチキャスト`vxlan`の場合は不要）
`tunnel.NAME.ttl`                      | integer | `vxlan`            | `1`                                                             | マルチキャストルーティングトポロジーに使用する固有の TTL
`user.*`                               | string  | -                  | -                                                               | ユーザー指定の自由形式のキー／バリューペア

```{note}
`bridge.external_interfaces` オプションは拡張された形式をサポートし、存在しない VLAN インターフェースを作成できるようにします。

拡張された形式は `<interfaceName>/<parentInterfaceName>/<vlanId>` です。
拡張された形式でリストに外部インターフェースが追加される際、システムはネットワークの作成に応じてインターフェースを自動で作成し、ネットワークが切断される際はその後インターフェースを削除します。
システムは <interfaceName> が既に存在しないか確認します。
もしインターフェース名が別の親や VLAN ID で使用中の場合、あるいはインターフェースの作成が成功しない場合、システムはエラーメッセージを表示し、元の状態に戻します。
```

(network-bridge-features)=
## サポートされている機能

`bridge`ネットワークタイプでは以下の機能がサポートされています:

- {ref}`network-acls`
- {ref}`network-forwards`
- {ref}`network-zones`
- {ref}`network-bgp`
- [`systemd-resolved`と統合するには](network-bridge-resolved)

```{toctree}
:maxdepth: 1
:hidden:

resolvedとの統合 </howto/network_bridge_resolved>
ファイアウォールの設定 </howto/network_bridge_firewalld>
```
