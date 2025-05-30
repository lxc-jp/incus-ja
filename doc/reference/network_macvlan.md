(network-macvlan)=
# macvlan ネットワーク

<!-- Include start macvlan intro -->
macvlan は仮想的な {abbr}`LAN (Local Area Network)` で同じネットワークインターフェースに複数の IP アドレスを割り当てたい場合に使用できます。
基本的にはネットワークインターフェースをそれぞれの IP アドレスを持つ複数のサブインターフェースに分割することになります。
その後ランダムに生成された MAC アドレスに基づいて IP アドレスを設定できます。
<!-- Include end macvlan intro -->

`macvlan` ネットワークタイプは親のインターフェースにインスタンスを接続する際に使用するプリセットを指定できます。
この場合、接続先のネットワークについて基本的な設定詳細を一切知る必要なしに単に `network` オプションをインスタンス NIC に設定できます。

```{note}
`macvlan` NIC を使う場合、Incus ホストとインスタンス間の通信はできません。
ホストとインスタンスの両方がゲートウェイと通信できますが、それらが直接通信はできません。
```

(network-macvlan-options)=
## 設定オプション

`macvlan` ネットワークタイプでは現在以下の設定キーNamespace がサポートされています:

- `user`（key/value の自由形式のユーザーメタデータ）

```{note}
{{note_ip_addresses_CIDR}}
```

`macvlan` ネットワークタイプでは以下の設定オプションが使用できます:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_macvlan-common start -->
    :end-before: <!-- config group network_macvlan-common end -->
```
