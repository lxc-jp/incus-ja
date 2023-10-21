(network-ovn-setup)=
# Incus で OVN をセットアップするには

スタンドアロンのネットワークとしてまたは小さな Incus クラスタとして基本的な OVN ネットワークをセットアップするには以下の項を参照してください。

## スタンドアロンの OVN ネットワークをセットアップする

外向きの接続のために Incus が管理する親のブリッジネットワーク（たとえば、 `incusbr0`）に接続するスタンドアロンの OVN ネットワークを作成するには以下の手順を実行してください。

1. ローカルサーバーに OVN ツールをインストールします:

       sudo apt install ovn-host ovn-central

1. OVN の統合ブリッジを設定します:

       sudo ovs-vsctl set open_vswitch . \
          external_ids:ovn-remote=unix:/var/run/ovn/ovnsb_db.sock \
          external_ids:ovn-encap-type=geneve \
          external_ids:ovn-encap-ip=127.0.0.1

1. OVN ネットワークを作成します:

       incus network set <parent_network> ipv4.dhcp.ranges=<IP_range> ipv4.ovn.ranges=<IP_range>
       incus network create ovntest --type=ovn network=<parent_network>

1. `ovntest` ネットワークを使用するインスタンスを作成します:

       incus init images:ubuntu/22.04 c1
       incus config device override c1 eth0 network=ovntest
       incus start c1

1. [`incus list`](incus_list.md) を実行してインスタンスの情報を表示します:

   ```{terminal}
   :input: incus list
   :scroll:

   +------+---------+---------------------+----------------------------------------------+-----------+-----------+
   | NAME |  STATE  |        IPV4         |                     IPV6                     |   TYPE    | SNAPSHOTS |
   +------+---------+---------------------+----------------------------------------------+-----------+-----------+
   | c1   | RUNNING | 192.0.2.2 (eth0)    | 2001:db8:cff3:5089:216:3eff:fef0:549f (eth0) | CONTAINER | 0         |
   +------+---------+---------------------+----------------------------------------------+-----------+-----------+
   ```

## OVN 上に Incus クラスタをセットアップする

OVN ネットワークを使用する Incus クラスタをセットアップするには以下の手順を実行してください。

Incus と同様に、 OVN の分散データベースは奇数のメンバーで構成されるクラスタ上で動かす必要があります。
以下の手順は最小構成の 3 台のサーバーを使います。 3 台のサーバーでは OVN の分散データベースと OVN コントローラーの両方を動かします。
さらに Incus クラスタに OVN コントローラーのみを動かすサーバーを任意の台数追加できます。

1. OVN の分散データベースを動かしたい 3 台のマシンで次の手順を実行してください:

   1. OVN ツールをインストールします:

          sudo apt install ovn-central ovn-host

   1. マシンの起動時に OVN サービスが起動されるように自動起動を有効にします:

           systemctl enable ovn-central
           systemctl enable ovn-host

   1. OVN を停止します:

          systemctl stop ovn-central

   1. マシンの IP アドレスをメモします:

          ip -4 a

   1. `/etc/default/ovn-central` を編集します。

   1. 以下の設定をペーストします（`<server_1>`, `<server_2>` and `<server_3>` をそれぞれのマシンの IP アドレスに、 `<local>` をあなたがいるマシンの IP アドレスに置き換えてください）。

      - 最初のマシン:

        ```
        OVN_CTL_OPTS=" \
             --db-nb-addr=<local> \
             --db-nb-create-insecure-remote=yes \
             --db-sb-addr=<local> \
             --db-sb-create-insecure-remote=yes \
             --db-nb-cluster-local-addr=<local> \
             --db-sb-cluster-local-addr=<local> \
             --ovn-northd-nb-db=tcp:<server_1>:6641,tcp:<server_2>:6641,tcp:<server_3>:6641 \
             --ovn-northd-sb-db=tcp:<server_1>:6642,tcp:<server_2>:6642,tcp:<server_3>:6642"
        ```

      - 2 番目と 3 番目のマシン:

        ```
        OVN_CTL_OPTS=" \
              --db-nb-addr=<local> \
             --db-nb-cluster-remote-addr=<server_1> \
             --db-nb-create-insecure-remote=yes \
             --db-sb-addr=<local> \
             --db-sb-cluster-remote-addr=<server_1> \
             --db-sb-create-insecure-remote=yes \
             --db-nb-cluster-local-addr=<local> \
             --db-sb-cluster-local-addr=<local> \
             --ovn-northd-nb-db=tcp:<server_1>:6641,tcp:<server_2>:6641,tcp:<server_3>:6641 \
             --ovn-northd-sb-db=tcp:<server_1>:6642,tcp:<server_2>:6642,tcp:<server_3>:6642"
        ```

   1. OVN を起動します:

          systemctl start ovn-central

1. 残りのマシンでは `ovn-host` のみインストールし、自動起動を有効にしてください:

       sudo apt install ovn-host
       systemctl enable ovn-host

1. すべてのマシンで Open vSwitch（変数は上記の通りに置き換えてください）を設定します:

       sudo ovs-vsctl set open_vswitch . \
          external_ids:ovn-remote=tcp:<server_1>:6642,tcp:<server_2>:6642,tcp:<server_3>:6642 \
          external_ids:ovn-encap-type=geneve \
          external_ids:ovn-encap-ip=<local>

1. すべてのマシンで `incus admin init` を実行して Incus クラスタを作成してください。
   最初のマシンでクラスタを作成します。
   次に最初のマシンで [`incus cluster add <machine_name>`](incus_cluster_add.md) を実行してトークンを出力し、他のマシンで Incus を初期化する際にトークンを指定して他のマシンをクラスタに参加させます。
1. 最初のマシンでアップリンクネットワークを作成し設定します:

       incus network create UPLINK --type=physical parent=<uplink_interface> --target=<machine_name_1>
       incus network create UPLINK --type=physical parent=<uplink_interface> --target=<machine_name_2>
       incus network create UPLINK --type=physical parent=<uplink_interface> --target=<machine_name_3>
       incus network create UPLINK --type=physical parent=<uplink_interface> --target=<machine_name_4>
       incus network create UPLINK --type=physical \
          ipv4.ovn.ranges=<IP_range> \
          ipv6.ovn.ranges=<IP_range> \
          ipv4.gateway=<gateway> \
          ipv6.gateway=<gateway> \
          dns.nameservers=<name_server>

   必要な値を決定します。

   アップリンクネットワーク
   : アクティブな OVN シャーシがクラスタメンバー間で移動できるようにするため、ハイアベイラビリティな OVN クラスタには共有されたレイヤー 2 ネットワークが必須です（これにより OVN のルータの外部 IP が実質的に別のホストから到達可能にできます）。

     そのため管理されていないブリッジインターフェースまたは使用されていない物理インターフェースを OVN アップリンクで使用される物理ネットワークの親として指定する必要があります。
     以下の手順は手動で作成した管理されていないブリッジを使用する想定です。
     このブリッジをセットアップする手順は [ネットワークブリッジの設定](https://netplan.readthedocs.io/en/stable/examples/#how-to-configure-network-bridges) を参照してください。

   ゲートウェイ
   : `ip -4 route show default` と `ip -6 route show default` を実行してください。

   ネームサーバー
   : `resolvectl` を実行してください。

   IP の範囲
   : 割り当てられた IP を元に適切な IP の範囲を使用してください。

1. 引き続き最初のマシンで Incus を OVN DB クラスタと通信できるように設定します。
   そのためには、 `/etc/default/ovn-central` 内の `ovn-northd-nb-db` の値を確認し、以下のコマンドで Incus に指定します:

       incus config set network.ovn.northbound_connection <ovn-northd-nb-db>

1. 最後に（最初のマシンで）実際の OVN ネットワークを作成します:

       incus network create my-ovn --type=ovn

1. OVN ネットワークをテストするには、インスタンスを作成してネットワークが接続できるか確認します:

       incus launch images:ubuntu/22.04 c1 --network my-ovn
       incus launch images:ubuntu/22.04 c2 --network my-ovn
       incus launch images:ubuntu/22.04 c3 --network my-ovn
       incus launch images:ubuntu/22.04 c4 --network my-ovn
       incus list
       incus exec c4 bash
       ping <IP of c1>
       ping <nameserver>
       ping6 -n www.example.com

## OVN ログを Incus に送信

OVN コントローラーのログを Incus に送るようにするには以下の手順を実行してください。

1. syslog ソケットを有効にします:

       incus config set core.syslog_socket=true

1. `/etc/default/ovn-host` を編集用に開きます:

1. 以下の設定をペーストします:

       OVN_CTL_OPTS=" \
              --ovn-controller-log='-vsyslog:info --syslog-method=unix:/var/lib/incus/syslog.socket'"

1. OVN コントローラーを再起動します:

       systemctl restart ovn-controller.service

これで [`incus monitor`](incus_monitor.md) を使って OVN コントローラーからのネットワーク ACL トラフィックのログを見られます:

    incus monitor --type=network-acls

また Loki にログを送ることもできます。
そのためには、たとえば、{config:option}`server-loki:loki.types`設定キーに`network-acl`の値を追加してください:

    incus config set loki.types=network-acl

```{tip}
OVN `northd`、OVN north-bound `ovsdb-server`、OVN south-bound `ovsdb-server`のログもインクルードできます。
そのためには、`/etc/default/ovn-central`を編集します:

    OVN_CTL_OPTS=" \
       --ovn-northd-log='-vsyslog:info --syslog-method=unix:/var/lib/incus/syslog.socket' \
       --ovn-nb-log='-vsyslog:info --syslog-method=unix:/var/lib/incus/syslog.socket' \
       --ovn-sb-log='-vsyslog:info --syslog-method=unix:/var/lib/incus/syslog.socket'"

    sudo systemctl restart ovn-central.service
```
