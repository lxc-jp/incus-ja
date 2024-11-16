(network-bridge-resolved)=
# `systemd-resolved` と統合するには

Incus を実行するシステムが DNS ルックアップの実行に `systemd-resolved` を使用する場合、 `resolved` に Incus が名前解決できるドメインを通知するべきです。
そうするには、 Incus ネットワークブリッジにより提供される DNS サーバーとドメインを `resolved` の設定に追加してください。

```{note}
この機能を使いたい場合、 `dns.mode` オプション ({ref}`network-bridge-options` 参照) を `managed` か `dynamic` に設定する必要があります。
```

(network-bridge-resolved-configure)=
## resolved を設定する

ネットワークブリッジを `resolved` 設定に追加するには、対応するブリッジの DNS アドレスとドメインを指定します。

DNS アドレス
: IPv4 アドレス、 IPv6 アドレス、あるいは両方を使用できます。
  アドレスはサブネットのネットマスク無しで指定する必要があります。

  ブリッジの IPv4 アドレスを取得するには以下のコマンドを使用します:

      incus network get <network_bridge> ipv4.address

  ブリッジの IPv6 アドレスを取得するには以下のコマンドを使用します:

      incus network get <network_bridge> ipv6.address

DNS ドメイン
: ブリッジの DNS ドメイン名を取得するには以下のコマンドを使用します:

      incus network get <network_bridge> dns.domain

  このオプションが設定されていない場合、デフォルトのドメイン名は `incus` です。

`resolved` を設定するには以下のコマンドを使用します:

    resolvectl dns <network_bridge> <dns_address>
    resolvectl domain <network_bridge> ~<dns_domain>

```{note}
`resolved`でDNSドメインを指定する場合、ドメイン名に `~` の接頭辞をつけてください。
`~` により `resolved` がこのドメインをルックアップするためだけに対応するネームサーバーを使うようになります。

ご利用のシェルによっては `~` が展開されるのを防ぐために DNS ドメインを引用符で囲む必要があるかもしれません。
```

DNSSEC と DNS over TLS

: `incus`のDNSサーバーはDNSSEC や DNS over TLS をサポートしません。

  resolved の設定によってはサーバーがDNSSEC や DNS over TLS をサポートしないため、設定がエラーになります。

  ブリッジでのみ両方を無効化するには、次のコマンドを実行します:

      resolvectl dnssec <network_bridge> off
      resolvectl dnsovertls <network_bridge> off

たとえば:

    resolvectl dns incusbr0 192.0.2.10
    resolvectl domain incusbr0 '~incus'
    resolvectl dnssec incusbr0 off
    resolvectl dnsovertls incusbr0 off

```{note}
別の方法として、 `systemd-resolve` コマンドを使用することもできます。
このコマンドは `systemd` の新しいリリースでは廃止予定となっていますが、後方互換性のため引き続き提供されています。

    systemd-resolve --interface <network_bridge> --set-domain ~<dns_domain> --set-dns <dns_address> --set-dnsovertls=off --set-dnssec=off
```

`resolved` の設定はブリッジが存在する限り残ります。
リブートのたびに Incus が再起動した後に上記のコマンドを実行するか、下記のように設定を永続的にする必要があります。

## `resolved` の設定を永続的にする

システムの起動時に適用され Incus がネットワークインターフェースを作成したときに有効になるように `systemd-resolved` の DNS 設定を自動化できます。

そうするには、 `/etc/systemd/system/lxd-dns-<network_bridge>.service` という名前の `systemd` ユニットファイルを以下の内容で作成してください:

```
[Unit]
Description=Incus per-link DNS configuration for <network_bridge>
BindsTo=sys-subsystem-net-devices-<network_bridge>.device
After=sys-subsystem-net-devices-<network_bridge>.device

[Service]
Type=oneshot
ExecStart=/usr/bin/resolvectl dns <network_bridge> <dns_address>
ExecStart=/usr/bin/resolvectl domain <network_bridge> ~<dns_domain>
ExecStart=/usr/bin/resolvectl dnssec <network_bridge> off
ExecStart=/usr/bin/resolvectl dnsovertls <network_bridge> off
ExecStopPost=/usr/bin/resolvectl revert <network_bridge>
RemainAfterExit=yes

[Install]
WantedBy=sys-subsystem-net-devices-<network_bridge>.device
```

ファイル名と内容で `<network_bridge>` をブリッジの名前（たとえば `incusbr0`）に置き換えてください。
さらに `<dns_address>` と `<dns_domain>` を {ref}`network-bridge-resolved-configure` に書かれているように置き換えてください。

次に以下のコマンドでサービスの自動起動を有効にし起動します:

    sudo systemctl daemon-reload
    sudo systemctl enable --now incus-dns-<network_bridge>

（Incus が既に実行中のため）対応するブリッジが既に存在する場合、以下のコマンドでサービスが起動したかを確認できます:

    sudo systemctl status incus-dns-<network_bridge>.service

以下のような出力になるはずです:

```{terminal}
:input: sudo systemctl status incus-dns-incusbr0.service

● incus-dns-incusbr0.service - Incus per-link DNS configuration for incusbr0
     Loaded: loaded (/etc/systemd/system/incus-dns-incusbr0.service; enabled; vendor preset: enabled)
     Active: inactive (dead) since Mon 2021-06-14 17:03:12 BST; 1min 2s ago
    Process: 9433 ExecStart=/usr/bin/resolvectl dns incusbr0 n.n.n.n (code=exited, status=0/SUCCESS)
    Process: 9434 ExecStart=/usr/bin/resolvectl domain incusbr0 ~incus (code=exited, status=0/SUCCESS)
   Main PID: 9434 (code=exited, status=0/SUCCESS)
```

`resolved` に設定が反映されたか確認するには、 `resolvectl status <network_bridge>` を実行します:

```{terminal}
:input: resolvectl status incusbr0

Link 6 (incusbr0)
      Current Scopes: DNS
DefaultRoute setting: no
       LLMNR setting: yes
MulticastDNS setting: no
  DNSOverTLS setting: no
      DNSSEC setting: no
    DNSSEC supported: no
  Current DNS Server: n.n.n.n
         DNS Servers: n.n.n.n
          DNS Domain: ~incus
```
