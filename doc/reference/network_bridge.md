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

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_bridge-common start -->
    :end-before: <!-- config group network_bridge-common end -->
```

```{note}
`bridge.external_interfaces` オプションは拡張された形式をサポートし、存在しない VLAN インターフェースを作成できるようにします。

拡張された形式は `<interfaceName>/<parentInterfaceName>/<vlanId>` です。
拡張された形式でリストに外部インターフェースが追加される際、システムはネットワークの作成に応じてインターフェースを自動で作成し、ネットワークが切断される際はその後インターフェースを削除します。
システムは`<interfaceName>`が既に存在しないか確認します。
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
