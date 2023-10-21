# ネットワークを作成するには

マネージドネットワークを作成し設定するには、[`incus network`](incus_network.md) コマンドとそのサブコマンドを使用します。
どのコマンドでも `--help` を追加すると使用方法と利用可能なフラグについてより詳細な情報を表示できます。

(network-types)=
## ネットワークタイプ

以下のネットワークタイプが利用できます:

```{list-table}
   :header-rows: 1

* - ネットワークタイプ
  - ドキュメント
  - 設定オプション
* - `bridge`
  - {ref}`network-bridge`
  - {ref}`network-bridge-options`
* - `ovn`
  - {ref}`network-ovn`
  - {ref}`network-ovn-options`
* - `macvlan`
  - {ref}`network-macvlan`
  - {ref}`network-macvlan-options`
* - `sriov`
  - {ref}`network-sriov`
  - {ref}`network-sriov-options`
* - `physical`
  - {ref}`network-physical`
  - {ref}`network-physical-options`

```

## ネットワークを作成する

ネットワークを作成するには以下のコマンドを実行します:

```bash
incus network create <name> --type=<network_type> [configuration_options...]
```

利用可能なネットワークタイプ一覧と設定オプションへのリンクは {ref}`network-types` を参照してください。

`--type` 引数を指定しない場合、デフォルトのタイプ `bridge` が使用されます。

(network-create-cluster)=
### クラスタ内にネットワークを作成する

Incus クラスタを実行していてネットワークを作成したい場合、各クラスタメンバーに別々にネットワークを作成する必要があります。
この理由はネットワーク設定は、たとえば親ネットワークインターフェースの名前のように、クラスタメンバー間で異なるかもしれないからです。

このため、まず `--target=<cluster_member>` フラグとメンバー用の適切な設定を指定して保留中のネットワークを作成する必要があります。
すべてのメンバーで同じネットワーク名を使うようにしてください。
次に実際にセットアップするために `--target` フラグなしでネットワークを作成してください。

たとえば、以下の一連のコマンドで 3 つのクラスタメンバー上に `UPLINK` という名前の物理ネットワークをセットアップします:

```{terminal}
:input: incus network create UPLINK --type=physical parent=br0 --target=vm01

Network UPLINK pending on member vm01
:input: incus network create UPLINK --type=physical parent=br0 --target=vm02
Network UPLINK pending on member vm02
:input: incus network create UPLINK --type=physical parent=br0 --target=vm03
Network UPLINK pending on member vm03
:input: incus network create UPLINK --type=physical
Network UPLINK created
```

{ref}`cluster-config-networks`も参照してください。

(network-attach)=
## インスタンスにネットワークをアタッチする

マネージドネットワークを作成後、それをインスタンスに{ref}`NIC デバイス <devices-nic>`としてアタッチできます。

そのためには、以下のコマンドを使います:

    incus network attach <network_name> <instance_name> [<device_name>] [<interface_name>]

デバイス名とインターフェース名は省略可能ですが、少なくともデバイス名は指定することをお勧めします。
指定しない場合、Incus はネットワーク名をデバイス名として使用しますが、紛らわしく問題を起こすかもしれません。
たとえば、Incus イメージは`eth0`インターフェースに IP 自動設定を行いますが、インターフェースの名前が違うと機能しません。

たとえば、`my-network`というネットワークを`my-instance`というインタンスに`eth0`デバイスとしてアタッチするには、以下のコマンドを入力します:

    incus network attach my-network my-instance eth0

### NICデバイスを追加する

[`incus network attach`](incus_network_attach.md) コマンドはインスタンスに NIC デバイスを追加するショートカットです。
別の方法として、通常通りネットワーク設定で NIC デバイスを追加できます:

    incus config device add <instance_name> <device_name> nic network=<network_name>

この方法を使う場合、必要に応じてネットワークのデフォルト設定をオーバーライドするように追加の設定をコマンドに追加できます。
すべての利用可能なデバイスオプションについては{ref}`NIC デバイス <devices-nic>`を参照してください。
