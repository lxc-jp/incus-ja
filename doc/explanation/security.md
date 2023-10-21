(exp-security)=
# セキュリティーについて

% Include content from [../../README.md](../../README.md)
```{include} ../../README.md
    :start-after: <!-- Include start security -->
    :end-before: <!-- Include end security -->
```

詳細な情報は以下のセクションを参照してください。

セキュリティー上の問題を発見した場合、その問題の報告方法については [Incus のセキュリティーポリシー](https://github.com/lxc-jp/incus-ja/blob/main/SECURITY.md)（原文: [Incus security policy](https://github.com/lxc/incus/blob/main/SECURITY.md)）を参照してください。

## サポートされているバージョン

サポートされていないバージョンの Incus は実運用環境では絶対に使用しないでください。

% Include content from [../../SECURITY.md](../../SECURITY.md)
```{include} ../../SECURITY.md
    :start-after: <!-- Include start supported versions -->
    :end-before: <!-- Include end supported versions -->
```

(security-daemon-access)=
## Incus デーモンへのアクセス

Incus は Unix ソケットを介してローカルにアクセスできるデーモンで、設定されていれば{abbr}`TLS(Transport Layer Security)`ソケットを介してリモートにアクセスすることもできます。
ソケットにアクセスできる人は、ホストデバイスやファイルシステムをアタッチしたり、すべてのインスタンスのセキュリティー機能をいじったりするなど、Incus を完全に制御することができます。

したがって、デーモンへのアクセスを信頼できるユーザーに制限するようにしてください。

### Incus デーモンへのローカルアクセス

Incus デーモンは root で動作し、ローカル通信用の Unix ソケットを提供します。
Incus のアクセス制御は、グループメンバーシップに基づいて行われます。
root ユーザーと `incus-admin` グループのすべてのメンバーがローカルデーモンと対話できます。

````{important}
% Include content from [../../README.md](../../README.md)
```{include} ../../README.md
    :start-after: <!-- Include start security note -->
    :end-before: <!-- Include end security note -->
```
````

(security_remote_access)=
### リモート API へのアクセス

デフォルトでは、デーモンへのアクセスはローカルでのみ可能です。
`core.https_address`という設定オプションを設定することで、同じ API を{abbr}`TLS (Transport Layer Security)`ソケットでネットワーク上に公開することができます。
手順は {ref}`server-expose` を参照してください。
リモートクライアントは、Incus に接続して、公開用にマークされたイメージにアクセスできます。

リモートクライアントが API にアクセスできるように、信頼できるクライアントとして認証する方法がいくつかあります。
詳細は{ref}`authentication`を参照してください。

本番環境では、`core.https_address`に、（ホスト上の任意のアドレスではなく）サーバーが利用可能な単一のアドレスを設定する必要があります。
さらに、許可されたホスト/サブネットからのみ Incus ポートへのアクセスを許可するファイアウォールルールを設定する必要があります。

(container-security)=
## コンテナのセキュリティー

Incus コンテナはセキュリティーのために幅広い機能を使うことができます。

デフォルトでは、コンテナは *非特権*（*unprivileged*）であり、ユーザーNamespace 内で動作することを意味し、コンテナ内のユーザーの能力を、コンテナが所有するデバイスに対する制限された権限を持つホスト上の通常のユーザーに制限します。

コンテナ間のデータ共有が必要ない場合は、{config:option}`instance-security:security.idmap.isolated`を有効にすることで、各コンテナに対して重複しない UID/GID マップを使用し、他のコンテナに対する潜在的な{abbr}`DoS (Denial of Service、サービス拒否)`攻撃を防ぐことができます。

Incus はまた、*特権*（*privileged*）コンテナを実行することができます。
しかし、これは（訳注:コンテナ内だけで）安全に root 権限を使えるわけではなく、そのようなコンテナの中でルートアクセスを持つユーザーは、閉じ込められた状態から逃れる方法を見つけるだけでなく、ホストを DoS することができてしまう点に注意してください。

コンテナのセキュリティーと私たちが使っているカーネルの機能についてのより詳細な情報は
[LXCセキュリティページ](https://linuxcontainers.org/ja/lxc/security/)にあります。

### コンテナ名の漏洩

デフォルトの設定ではシステム上のすべての cgroup と、さらに転じて、すべての実行中のコンテナを一覧表示することが簡単にできてしまいます。

コンテナを開始する前に `/sys/kernel/slab` と `/proc/sched_debug` へのアクセスをブロックすることでコンテナ名の漏洩を防げます。
このためには以下のコマンドを実行してください:

    chmod 400 /proc/sched_debug
    chmod 700 /sys/kernel/slab/

## ネットワークセキュリティー

ネットワークインターフェースは必ず安全に設定してください。
どのような点を考慮すべきかは、使用するネットワークモードによって異なります。

### ブリッジ型NICのセキュリティー

Incus のデフォルトのネットワークモードは、各インスタンスが接続する「管理された」プライベートネットワークのブリッジを提供することです。
このモードでは、ホスト上に`incusbr0`というインターフェースがあり、それがインスタンスのブリッジとして機能します。

ホストは、管理されたブリッジごとに`dnsmasq`のインスタンスを実行し、IP アドレスの割り当てと、権威 DNS および再帰 DNS サービスの提供を担当します。

DHCPv4 を使用しているインスタンスには、IPv4 アドレスが割り当てられ、インスタンス名の DNS レコードが作成されます。
これにより、インスタンスが DHCP リクエストに偽のホスト名情報を提供して、DNS レコードを偽装することができなくなります。

`dnsmasq`サービスは、IPv6 のルータ広告機能も提供します。
つまり、インスタンスは SLAAC を使って自分の IPv6 アドレスを自動設定するので、`dnsmasq`による割り当ては行われません。
しかし、DHCPv4 を使用しているインスタンスは、SLAAC IPv6 アドレスに相当する AAAA の DNS レコードも取得します。
これは、インスタンスが IPv6 アドレスを生成する際に、IPv6 プライバシー拡張を使用していないことを前提としています。

このデフォルト構成では、DNS 名を偽装することはできませんが、インスタンスはイーサネットブリッジに接続されており、希望するレイヤー2 トラフィックを送信することができます。これは、信頼されていないインスタンスがブリッジ上で MAC または IP の偽装を効果的に行うことができることを意味します。

デフォルトの設定では、ブリッジに接続されたインスタンスがブリッジに（潜在的に悪意のある）IPv6 ルータ広告を送信することで、Incus ホストの IPv6 ルーティングテーブルを修正することも可能です。
これは、`incusbr0`インターフェースが`/proc/sys/net/ipv6/conf/incusbr0/accept_ra`を`2`に設定して作成されているためで、`forwarding`が有効であるにもかかわらず、Incus ホストがルーター広告を受け入れることを意味しています（詳細は[`/proc/sys/net/ipv4/*` Variables](https://www.kernel.org/doc/Documentation/networking/ip-sysctl.txt)を参照してください）。

しかし、Incus はいくつかのブリッジ型{abbr}`NIC(Network interface controller)`セキュリティー機能を提供しており、インスタンスがネットワーク上に送信することを許可されるトラフィックの種類を制御するために使用することができます。
これらの NIC 設定は、インスタンスが使用しているプロファイルに追加する必要がありますが、以下のように個々のインスタンスに追加することもできます。

ブリッジ型 NIC には、以下のようなセキュリティー機能があります:

キー                      | タイプ | デフォルト | 必須 | 説明
:--                       | :--    | :--        | :--  | :--
`security.mac_filtering`  | bool   | `false`    | no   | インスタンスが他のインスタンスの MAC アドレスを詐称することを防ぐ。
`security.ipv4_filtering` | bool   | `false`    | no   | インスタンスが他のインスタンスの IPv4 アドレスになりすますことを防ぎます(`mac_filtering` を有効にします)。
`security.ipv6_filtering` | bool   | `false`    | no   | インスタンスが他のインスタンスの IPv6 アドレスになりすますことを防ぎます(`mac_filtering` を有効にします)。

プロファイルで設定されたデフォルトのブリッジ型 NIC の設定は、インスタンスごとに以下の方法で上書きすることができます:

```
incus config device override <instance> <NIC> security.mac_filtering=true
```

これらの機能を併用することで、ブリッジに接続されているインスタンスが MAC アドレスや IP アドレスを詐称することを防ぐことができます。
これらのオプションは、ホスト上で利用可能なものに応じて、`xtables`（`iptables`、`ip6tables`、`ebtables`）または`nftables`を使用して実装されます。

これらのオプションは、ネストされたコンテナが異なる MAC アドレスを持つ親ネットワークを使用すること（ブリッジされた NIC や`macvlan` NIC を使用すること）を効果的に防止することができるのは注目に値します。

IP フィルタリング機能は、スプーフィングされた IP を含む ARP および NDP アドバタイジングをブロックし、スプーフィングされたソースアドレスを含むすべてのパケットをブロックします。

`security.ipv4_filtering`または`security.ipv6_filtering`が有効で、（`ipvX.address=none`またはブリッジで DHCP サービスが有効になっていないため）インスタンスに IP アドレスが割り当てられない場合、そのプロトコルのすべての IP トラフィックがインスタンスからブロックされます。

`security.ipv6_filtering` が有効な場合、IPv6 のルータ広告がインスタンスからブロックされます。

`security.ipv4_filtering`または`security.ipv6_filtering`が有効な場合、ARP、IPv4 または IPv6 ではないイーサネットフレームはすべてドロップされます。
これにより、スタックされた VLAN `Q-in-Q`（802.1ad）フレームが IP フィルタリングをバイパスすることを防ぎます。

### ルート化されたNICのセキュリティー

"routed" と呼ばれる別のネットワークモードがあります。
このモードでは、コンテナとホストの間に仮想イーサネットデバイペアを提供します。
このネットワークモードでは、Incus ホストがルータとして機能し、コンテナの IP 宛のトラフィックをコンテナの`veth`インターフェースに誘導するスタティックルートがホストに追加されます。

デフォルトでは、コンテナからのルータ広告が Incus ホスト上の IPv6 ルーティングテーブルを変更するのを防ぐために、ホスト上に作成された`veth`インターフェースは、その`accept_ra`設定が無効になっています。
それに加えて、コンテナが持っていることをホストが知らない IP に対するソースアドレスの偽装を防ぐために、ホスト上の`rp_filter`が`1`に設定されています。
