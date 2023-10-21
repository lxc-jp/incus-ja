(instances-routed-nic-vm)=
# 仮想マシンに routed NIC を追加するには

インスタンスに {ref}`routed NIC デバイス <nic-routed>` を追加する際、link-local ゲートウェイ IP をデフォルトルートとして使用するようにインスタンスを設定する必要があります。
コンテナでは、これは自動的に設定されます。
仮想マシンでは、ゲートウェイは手動あるいは `cloud-init` のような仕組みで設定する必要があります。

`cloud-init` でゲートウェイを設定するには、まずインスタンスを初期化します:

    incus init images:ubuntu/22.04 jammy --vm

次に routed NIC デバイスを追加します:

    incus config device add jammy eth0 nic nictype=routed parent=my-parent-network ipv4.address=192.0.2.2 ipv6.address=2001:db8::2

このコマンドでは、`my-parent-network` が親ネットワークで、IPv4 と IPv6 アドレスは親のサブネット内です。

次に `cloud-init.network-config` 設定キーを使ってインスタンスに `netplan` 設定を追加します:

    cat <<EOF | incus config set jammy cloud-init.network-config -
    network:
      version: 2
      ethernets:
        enp5s0:
          routes:
          - to: default
            via: 169.254.0.1
            on-link: true
          - to: default
            via: fe80::1
            on-link: true
          addresses:
          - 192.0.2.2/32
          - 2001:db8::2/128
    EOF

この `netplan` 設定は必要な {ref}`スタティックな link-local next-hop アドレス <nic-routed>`（`169.254.0.1` と `fe80::1`）を追加します。
これらのルートはそれぞれ `on-link` を `true` に設定します。するとルートがインターフェースに直接接続されるよう指定されます。
また routed NIC デバイス内でのアドレスも追加します。
`netplan` の詳細は [ドキュメント](https://netplan.readthedocs.io/en/latest/) を参照してください。

```{note}
この `netplan` 設定はネームサーバーを含んでいません。
インスタンス内で DNS を使うには、有効な DNS の IP アドレスを設定する必要があります。
ホストに `incusbr0` ネットワークがあれば、ネームサーバーは代わりにその IP を指定できます。
```

これでネットワークを開始できます:

    incus start jammy

```{note}
インスタンスを輝度する前に、 proxy ARP/NDP を有効にするように {ref}`親のネットワークを設定した <nic-routed>` ことを確認してください。
```
