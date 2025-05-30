(network-sriov)=
# SR-IOV ネットワーク

<!-- Include start SR-IOV intro -->
{abbr}`SR-IOV (Single root I/O virtualization)` は仮想環境内で単一のネットワークポートを複数の仮想ネットワークインターフェースのように見せるように出来るハードウェア標準です。
<!-- Include end SR-IOV intro -->

`sriov` ネットワークタイプは親のインターフェースに接続する際に使用するプリセットを指定できるようにします。
この場合接続先の設定詳細を一切知ること無くインスタンス NIC に単に `network` オプションを設定できます。

(network-sriov-options)=
## 設定オプション

`sriov` ネットワークでは現在以下の設定キーNamespace がサポートされています。

- `user`（key/value の自由形式のユーザーメタデータ）

```{note}
{{note_ip_addresses_CIDR}}
```

`sriov` ネットワークタイプには以下の設定オプションがあります。

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_sriov-common start -->
    :end-before: <!-- config group network_sriov-common end -->
```
