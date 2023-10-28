(network-bridge-firewall)=
# ファイアウォールを設定するには

Linux のファイアウォールは `netfilter` をベースにしています。
Incus は同じサブシステムを使用しているため、接続に問題を引き起こすことがありえます。

ファイアウォールを動かしている場合、 Incus が管理しているブリッジとホストの間のネットワークトラフィックを許可するように設定する必要があるかもしれません。
そうしないと、一部のネットワークの機能（DHCP、DNS と外部ネットワークへのアクセス）が期待通り動かないかもしれません。

ファイアウォール（あるいは他のアプリケーション）に設定されたルールと Incus が追加するファイアウォールのルールが衝突するケースがあります。
たとえば、ファイアウォールが Incus デーモンより後に起動した場合ファイアウォールが Incus のルールを削除するかもしれず、そうするとインスタンスへのネットワーク接続を妨げるかもしれません。

## `xtables` 対 `nftables`

`netfilter` にルールを追加するには `xtables`（IPv4 には `iptables` と IPv6 には `ip6tables`）と `nftables` という異なるユーザースペースのコマンドがあります。

`xtables` は順序ありのルールのリストを提供しますが、そのため複数のシステムがルールの追加や削除を行うと問題が起きるかもしれません。
`nftables` は分離されたルールを別々の Namespace に追加することができますので、異なるアプリケーションからのルールを分離するのに役立ちます。
しかし、パケットが 1 つの Namespace でブロックされる場合、他の Namespace がそれを許可することはできません。
そのため、 1 つの Namespace が他の Namespace のルールへ影響することは依然としてあり、ファイアウォールのアプリケーションが Incus のネットワーク機能に影響することがありえます。

システムで `nftables` を利用可能な場合、 Incus はそれを検知して `nftables` モードにスイッチします。
このモードでは Incus は自身の `nftables` の Namespace を用いてルールを `nftables` に追加します。

## Incus のファイアウォールを使用する

デフォルトでは Incus が管理するブリッジはフル機能を使えるようにするためファイアウォールにルールを追加します。
システムで他のファイアウォールを使用していない場合は Incus にファイアウォールのルールを管理させることができます。

これを有効または無効にするには `ipv4.firewall` または `ipv6.firewall` {ref}`設定オプション <network-bridge-options>` を使用してください。

## 別のファイアウォールを使用する

別のアプリケーションが追加するファイアウォールのルールは Incus が追加するファイアウォールルールと干渉するかもしれません。
このため、別のファイアウォールを使用する場合は Incus のファイアウォールルールを無効にするべきです。
また Incus のインスタンスがホスト上で Incus が動かしている DHCP と DNS サーバーにアクセスできるようにするため、
インスタンスと Incus ブリッジ間のネットワークトラフィックを許可するように設定しなければなりません。

Incus のファイアウォールルールをどのように無効化し、 `firewalld` と UFW をどのように適切に設定するかは以下を参照してください。

### Incus のファイアウォールルールを無効化する

指定のネットワークブリッジ（たとえば、`incusbr0`）に Incus がファイアウォールルールを設定しないようにするためには以下のコマンドを実行してください:

    incus network set <network_bridge> ipv6.firewall false
    incus network set <network_bridge> ipv4.firewall false

### `firewalld` で信頼されたゾーンにブリッジを追加する

`firewalld` で Incus ブリッジへとブリッジからのトラフィックを許可するには、ブリッジインターフェースを `trusted` ゾーンに追加してください。
（再起動後も設定が残るように）永続的にこれを行うには以下のコマンドを実行してください:

    sudo firewall-cmd --zone=trusted --change-interface=<network_bridge> --permanent
    sudo firewall-cmd --reload

たとえば:

    sudo firewall-cmd --zone=trusted --change-interface=incusbr0 --permanent
    sudo firewall-cmd --reload

```{warning}
<!-- Include start warning -->
上に示したコマンドはシンプルな例です。
あなたの使い方に応じて、より高度なルールが必要な場合があり、その場合上の例をそのまま実行するとうっかりセキュリティリスクを引き起こす可能性があります。
<!-- Include end warning -->
```

### UFW でブリッジにルールを追加する

UFW で認識不能なトラフィックをすべてドロップするルールを入れていると、 Incus ブリッジへとブリッジからのトラフィックをブロックしてしまいます。
この場合ブリッジへとブリッジからのトラフィックを許可し、さらにブリッジへフォワードされるトラフィックを許可するルールを追加する必要があります。

そのためには次のコマンドを実行します:

    sudo ufw allow in on <network_bridge>
    sudo ufw route allow in on <network_bridge>
    sudo ufw route allow out on <network_bridge>

たとえば:

    sudo ufw allow in on incusbr0
    sudo ufw route allow in on incusbr0
    sudo ufw route allow out on incusbr0

````{warning}
% Repeat warning from above
```{include} network_bridge_firewalld.md
    :start-after: <!-- Include start warning -->
    :end-before: <!-- Include end warning -->
```

以下はより制限の強いファイアウォールの例で、ゲストからホストへのアクセスは DHCP と DNS のみに限定し、外向きの通信は全て許可します:

```
# ゲストが Incus ホストから IP を取得できるようにする
sudo ufw allow in on incusbr0 to any port 67 proto udp
sudo ufw allow in on incusbr0 to any port 547 proto udp

# ゲストが Incus ホストからホスト名を解決できるようにする
sudo ufw allow in on incusbr0 to any port 53

# ゲストが外向きの通信にアクセスできるようにする
CIDR4="$(incus network get incusbr0 ipv4.address | sed 's|\.[0-9]\+/|.0/|')"
CIDR6="$(incus network get incusbr0 ipv6.address | sed 's|:[0-9]\+/|:/|')"
sudo ufw route allow in on incusbr0 from "${CIDR4}"
sudo ufw route allow in on incusbr0 from "${CIDR6}"
```
````

(network-incus-docker)=
## Incus と Docker の接続の問題を回避する

同じホストで Incus と Docker を動かすと接続の問題を引き起こします。
この問題のよくある理由は Docker はグローバルの FOWARD のポリシーを `drop` に設定するので、それが Incus がトラフィックをフォワードすることを妨げインスタンスのネットワーク接続を失わせるということです。
詳細は [Docker on a router](https://docs.docker.com/network/iptables/#docker-on-a-router) を参照してください。

この問題を回避するためのさまざまな方法があります:

Docker をアンインストールする
: このような問題を防ぐ最も簡単な方法は、Incus を実行しているシステムから Docker をアンインストールしてシステムを再起動することです。
  代わりに、Incus のコンテナや仮想マシンの中で Docker を実行できます。

IPv4 の転送を有効にする
: Docker をアンインストールすることができない場合、Docker サービスが開始する前に IPv4 転送を有効にすることで、Docker がグローバル FORWARD ポリシーを変更するのを防ぐことができます。
  Incus ブリッジネットワークは通常、この設定を有効にします。
  しかし、Incus が Docker の後に起動すると、Docker は既にグローバル FORWARD ポリシーを変更している可能性があります。

  ```{warning}
  IPv4の転送を有効にすると、Dockerのコンテナポートがローカルネットワーク上の任意のマシンからアクセス可能になる可能性があります。
  環境によりますが、これは望ましくない場合があります。
  詳細については、[ローカルネットワークのコンテナアクセス問題](https://github.com/moby/moby/issues/14041)を参照してください。
  ```

  Docker が開始する前に IPv4 転送を有効にするためには、次の`sysctl`設定が有効になっていることを確認します:

      net.ipv4.conf.all.forwarding=1

  ```{important}
  この設定はホストの再起動時にも保持されるようにする必要があります。

  これを行う一つの方法は、次のコマンドを使用して`/etc/sysctl.d/`ディレクトリにファイルを追加することです:

      echo "net.ipv4.conf.all.forwarding=1" > /etc/sysctl.d/99-forwarding.conf
      systemctl restart systemd-sysctl

  ```

外向きネットワークトラフィックフローを許可する
: Docker のコンテナポートがローカルネットワーク上の任意のマシンからアクセス可能になる可能性を避けたい場合、Docker が提供するより複雑なソリューションを適用できます。

  次のコマンドを使用して、Incus 管理ブリッジインターフェースからの外向きネットワークトラフィックフローを明示的に許可します:

      iptables -I DOCKER-USER -i <network_bridge> -j ACCEPT
      iptables -I DOCKER-USER -o <network_bridge> -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT

  たとえば、Incus の管理ブリッジが`incusbr0`と呼ばれている場合、次のコマンドを使用して外向きトラフィックのフローを許可できます:

      iptables -I DOCKER-USER -i incusbr0 -j ACCEPT
      iptables -I DOCKER-USER -o incusbr0 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT

  ```{important}
  これらのファイアウォールルールは、ホストの再起動時にも保持されるようにする必要があります。
  これを行う方法は Linux ディストリビューションによります。
  ```
