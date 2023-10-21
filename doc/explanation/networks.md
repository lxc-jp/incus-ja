(networks)=
# ネットワークについて

あなたのインスタンスをインターネットに接続するにはいろいろな方法があります。
最も簡単な方法は Incus の初期化時にネットワークブリッジを作ってすべてのインスタンスでこのブリッジを使うことですが、 Incus はネットワークに関するさまざまな高度な設定をサポートします。

## ネットワークデバイス

インスタンスへの直接のネットワークアクセスを許可するには、 {abbr}`NIC (Network Interface Controller)` とも呼ばれるネットワークデバイスを最低 1 つ割り当てる必要があります
ネットワークデバイスは以下のどれかの方法で設定できます:

- Incus の初期化中にセットアップしたデフォルトのネットワークブリッジを使用する。
  デフォルトの設定を表示するにはデフォルトのプロファイルを確認します:

        incus profile show default

  この方法はインスタンスのネットワークを指定しない場合に使用します。
- 既存のネットワークインターフェースをインスタンスにネットワークデバイスとして追加して使用する。
  このネットワークインターフェースは Incus の制御外です。
  そのため、ネットワークインターフェースを使用するために必要なすべての情報を Incus に指定する必要があります。

  以下のようなコマンドを使用します:

        incus config device add <instance_name> <device_name> nic nictype=<nic_type> ...

  指定可能な NIC タイプの一覧とそれらの設定プロパティについては [タイプ: `nic`](devices-nic) を参照してください。

  たとえば、既存の Linux ブリッジ (`br0`) を追加するには以下のコマンドを使えます:

        incus config device add <instance_name> eth0 nic nictype=bridged parent=br0
- {doc}`マネージドネットワークを作成 </howto/network_create>` し、それをインスタンスにネットワークデバイスとして追加する。
  この方法では Incus は設定されるネットワークについてのすべての必要な情報を持っていますので、デバイスとしてインスタンスに直接アタッチできます:

        incus network attach <network_name> <instance_name> <device_name>

  詳細は{ref}`network-attach`を参照してください。

(managed-networks)=
## マネージドネットワーク

Incus でマネージドネットワークは `incus network [create|edit|set]` コマンドで作成と設定をします。

ネットワークタイプによって、 Incus はネットワークを完全に制御するか、単に外部のネットワークインターフェースを管理するかのどちらかになります。

すべての {ref}`NIC タイプ <devices-nic>` がネットワークタイプとしてサポートされているわけではないことに注意してください。
Incus はいくつかのタイプのみマネージドネットワークとしてセットアップできます。

### 完全に制御されるネットワーク

完全に制御されるネットワークではネットワークインターフェースを作成し、たとえば IP を管理する機能を含むほとんどの機能を提供します。

Incus は以下のネットワークタイプをサポートします:

{ref}`network-bridge`
: % Include content from [../reference/network_bridge.md](../reference/network_bridge.md)
  ```{include} ../reference/network_bridge.md
      :start-after: <!-- Include start bridge intro -->
      :end-before: <!-- Include end bridge intro -->
  ```

  Incus の文脈では、 `bridge` ネットワークタイプは、ブリッジを共用するインスタンスを同一の L2 ネットワークセグメントに接続するような L2 ブリッジを作成します。
  これによりインスタンス間のトラフィックを通すことができます。
  ブリッジはさらにローカルの DHCP と DNS を提供することもできます。

  これがデフォルトのネットワークタイプです。

{ref}`network-ovn`
: % Include content from [../reference/network_ovn.md](../reference/network_ovn.md)
  ```{include} ../reference/network_ovn.md
      :start-after: <!-- Include start OVN intro -->
      :end-before: <!-- Include end OVN intro -->
  ```

  Incus の文脈では、 `ovn` ネットワークタイプは論理ネットワークを作成します。
  セットアップするには OVN ツールをインストールし設定する必要があります。
  さらに、OVN にネットワーク接続を提供するアップリンクのネットワークを作成する必要があります。
  アップリンクのネットワークとして、外部ネットワークタイプの 1 つかマネージドな Incus ブリッジを使う必要があります。

  ```{tip}
  他のネットワークタイプと違って、 OVN ネットワークは {ref}`プロジェクト <projects>` 内に作成・管理できます。
  これは、制限されたプロジェクトであっても、非管理者ユーザとして自身の OVN ネットワークを作成できることを意味します。
  ```

### 外部ネットワーク

% Include content from [../reference/network_external.md](../reference/network_external.md)
```{include} ../reference/network_external.md
    :start-after: <!-- Include start external intro -->
    :end-before: <!-- Include end external intro -->
```

{ref}`network-macvlan`
: % Include content from [../reference/network_macvlan.md](../reference/network_macvlan.md)
  ```{include} ../reference/network_macvlan.md
      :start-after: <!-- Include start macvlan intro -->
      :end-before: <!-- Include end macvlan intro -->
  ```

  Incus の文脈では、 `macvlan` ネットワークタイプは親の macvlan インターフェースへインスタンスを接続する際に使用するプリセット設定を提供します。

{ref}`network-sriov`
: % Include content from [../reference/network_sriov.md](../reference/network_sriov.md)
  ```{include} ../reference/network_sriov.md
      :start-after: <!-- Include start SR-IOV intro -->
      :end-before: <!-- Include end SR-IOV intro -->
  ```

  Incus の文脈では、 `sriov` ネットワークタイプは親の SR-IOV インターフェースへインスタンスを接続する際に使用するプリセット設定を提供します。

{ref}`network-physical`
: % Include content from [../reference/network_physical.md](../reference/network_physical.md)
  ```{include} ../reference/network_physical.md
      :start-after: <!-- Include start physical intro -->
      :end-before: <!-- Include end physical intro -->
  ```

  OVN ネットワークを親インターフェースに接続する際のプリセット設定を提供します。

## お勧めの設定

一般に、マネージドネットワークは設定が容易で設定を繰り返すこと無く複数のインスタンスで同じネットワークを再利用できるので、マネージドネットワークが使用できる場合はこれを使用すべきです。

どのネットワークタイプを使用すべきかはあなたの固有の使い方によります。
完全に制御されたネットワークを選ぶと、ネットワークデバイスを使用するのに比べてより多くの機能を提供します。

一般的なお勧めとしては:

- Incus を単一のシステム上かパブリッククラウドで動かしている場合は、 {ref}`network-bridge` を使用してください。
- あなた自身のプライベートクラウドで Incus を動かしている場合は、 {ref}`network-ovn` を使用してください。

  ```{note}
  OVN は適切な運用には共有された L2 のアップリンクネットワークが必要です。
  このため、パブリッククラウドで Incus を動かしている場合は通常 OVN は使用できません。
  ```

- インスタンス NIC をマネージドネットワークに接続するためには、可能であれば `parent` プロパティより `network` プロパティを使用してください。
  こうすることで、 NIC はネットワークの設定を引き継ぎ、 `nictype` を指定する必要がなくなります。
