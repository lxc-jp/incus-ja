(server-expose)=
# Incusをネットワークに公開するには

デフォルトでは、Incus は Unix ソケットを介してローカルユーザーからのみ使用でき、ネットワーク経由でアクセスすることはできません。

Incus をネットワークに公開するには、ローカル Unix ソケット以外のアドレスをリッスンするように設定する必要があります。これを行うには、{config:option}`server-core:core.https_address` サーバー設定オプションを設定します。

たとえば、Incus サーバーをポート`8443`でアクセスできるようにするには、以下のコマンドを入力します:

    incus config set core.https_address :8443

特定の IP アドレスからのアクセスを許可するには、`ip addr`を使用して利用可能なアドレスを見つけ、それを設定します。
たとえば:

```{terminal}
:input: ip addr

1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: enp5s0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    link/ether 10:66:6a:e3:f3:3f brd ff:ff:ff:ff:ff:ff
    inet 10.68.216.12/24 metric 100 brd 10.68.216.255 scope global dynamic enp5s0
       valid_lft 3028sec preferred_lft 3028sec
    inet6 fd42:e819:7a51:5a7b:1266:6aff:fee3:f33f/64 scope global mngtmpaddr noprefixroute
       valid_lft forever preferred_lft forever
    inet6 fe80::1266:6aff:fee3:f33f/64 scope link
       valid_lft forever preferred_lft forever
3: incusbr0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default qlen 1000
    link/ether 10:66:6a:8d:f3:72 brd ff:ff:ff:ff:ff:ff
    inet 10.64.82.1/24 scope global incusbr0
       valid_lft forever preferred_lft forever
    inet6 fd42:f4ab:4399:e6eb::1/64 scope global
       valid_lft forever preferred_lft forever
:input: incus config set core.https_address 10.68.216.12
```

すべてのリモートクライアントは Incus に接続して公開利用とマークされた任意のイメージにアクセスできます。

(server-authenticate)=
## Incusサーバーでの認証

リモート API にアクセスできるようにするには、クライアントは Incus サーバーに認証しなければなりません。
いくつかの認証方法があります。詳細は{ref}`authentication`を参照してください。

お勧めの方法はクライアントの TLS 証明書をトラストトークンを使ってサーバーのトラストストアに追加することです。
トラストトークンを使ってクライアントを認証するには、以下の手順を実行します:

1. サーバーで、以下のコマンドを入力します:

       incus config trust add <client_name>

   クライアント証明書を追加するのに使用できるトークンをコマンドが生成し表示します。
1. クライアントで、以下のコマンドでサーバーを追加します:

       incus remote add <remote_name> <token>

   % Include content from [../authentication.md](../authentication.md)
```{include} ../authentication.md
    :start-after: <!-- Include start NAT authentication -->
    :end-before: <!-- Include end NAT authentication -->
```

詳細や他の認証方法については{ref}`authentication`を参照してください。
