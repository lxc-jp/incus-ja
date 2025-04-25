(devices-proxy)=
# タイプ: `proxy`

```{note}
`proxy`デバイスタイプはコンテナ（NAT と非 NAT モード）と VM（NAT モードのみ）でサポートされます。
コンテナと VM の両方でホットプラグをサポートします。
```

プロキシデバイスにより、ホストとインスタンス間のネットワーク接続を転送できます。
この方法で、ホストのアドレスの一つに到達したトラフィックをインスタンス内のアドレスに転送したり、その逆でインスタンス内にアドレスを持ちホストを通して接続することができます。

{ref}`devices-proxy-nat-mode`では、プロキシデバイスを TCP と UDP のプロキシに使用することができます。
NAT モードではない場合、Unix ソケット間のトラフィックをプロキシすることもできます（これはたとえば、コンテナからホストシステムへのグラフィカルな GUI やオーディオトラフィックを転送するのに便利です）。また、プロトコル間でもプロキシすることができます（たとえば、ホストシステム上に TCP リスナーを設置し、そのトラフィックをコンテナ内の Unix ソケットに転送することができます）。

利用できる接続タイプは次の通りです:

- `tcp <-> tcp`
- `udp <-> udp`
- `unix <-> unix`
- `tcp <-> unix`
- `unix <-> tcp`
- `udp <-> tcp`
- `tcp <-> udp`
- `udp <-> unix`
- `unix <-> udp`

`proxy`デバイスを追加するには、以下のコマンドを使用します:

    incus config device add <instance_name> <device_name> proxy listen=<type>:<addr>:<port>[-<port>][,<port>] connect=<type>:<addr>:<port> bind=<host/instance>

(devices-proxy-nat-mode)=
## NATモード

プロキシデバイスは NAT モード（`nat=true`）もサポートします。NAT モードではパケットは別の接続を通してプロキシされるのではなく NAT を使ってフォワードされます。
これはターゲットの送り先が HAProxy の PROXY プロトコル（非 NAT モードでプロキシデバイスを使う場合はこれはクライアントアドレスを渡す唯一の方法です）をサポートする必要なく、クライアントのアドレスを維持できるという利点があります。

しかし、NAT モードはインスタンスが稼働しているホストがゲートウェイの場合（たとえば `incusbr0`を使用しているケース）のみサポートされます。

NAT モードでサポートされる接続のタイプは以下の通りです:

- `tcp <-> tcp`
- `udp <-> udp`

プロキシデバイスを`nat=true`に設定する際は、以下のようにターゲットのインスタンスが NIC デバイス上に静的 IP を持つようにする必要があります。

## IPアドレスを指定する

インスタンス NIC に静的 IP を設定するには、以下のコマンドを使用します:

    incus config device set <instance_name> <nic_name> ipv4.address=<ipv4_address> ipv6.address=<ipv6_address>

静的な IPv6 アドレスを設定するためには、親のマネージドネットワークは`ipv6.dhcp.stateful`を有効にする必要があります。

IPv6 アドレスを設定する場合は以下のような角括弧の記法を使います。たとえば以下のようにします:

    connect=tcp:[2001:db8::1]:80

connect のアドレスをワイルドカード（IPv4 では 0.0.0.0、IPv6 では[::]にします）に設定することで、接続アドレスをインスタンスの IP アドレスになるように指定できます。

```{note}
listen のアドレスも非 NAT モードではワイルドカードのアドレスが使用できます。
しかし、NAT モードを使う際は Incus ホスト上の IP アドレスを指定する必要があります。
```

## デバイスオプション

`proxy` デバイスには以下のデバイスオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-proxy start -->
    :end-before: <!-- config group devices-proxy end -->
```
