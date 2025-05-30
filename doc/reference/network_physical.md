(network-physical)=
# 物理ネットワーク

<!-- Include start physical intro -->
物理 (`physical`) ネットワークタイプは既存のネットワークに接続します。これはネットワークインターフェースまたはブリッジになることができ、OVN のためのアップリンクネットワークとしての役目を果たします。
<!-- Include end physical intro -->

このネットワークタイプは OVN ネットワークを親インターフェースに接続する際に使用するプリセットの設定を提供したり、インスタンスが物理インターフェースを NIC として使用できるようにします。この場合、インスタンス NIC は接続先の設定詳細を知ること無く単に `network` オプションを設定できるようにします。

(network-physical-options)=
## 設定オプション

物理ネットワークでは現在以下の設定キーNamespace がサポートされています:

- `bgp`（BGP ピア設定）
- `dns`（DNS サーバーと名前解決の設定）
- `ipv4`（L3 IPv4 設定）
- `ipv6`（L3 IPv6 設定）
- `ovn`（OVN 設定）
- `user`（key/value の自由形式のユーザーメタデータ）

```{note}
{{note_ip_addresses_CIDR}}
```

物理ネットワークタイプには以下の設定オプションがあります:

## BGPオプション

OVNダウンストリームネットワークのBGPピアの設定には以下のオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_physical-bgp start -->
    :end-before: <!-- config group network_physical-bgp end -->
```

## DNSオプション

物理ネットワークで使われるDNSサーバーとサーチドメインを制御するには以下のオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_physical-dns start -->
    :end-before: <!-- config group network_physical-dns end -->
```

## IPV4オプション

物理ネットワークのIPv4の設定には以下のオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_physical-ipv4 start -->
    :end-before: <!-- config group network_physical-ipv4 end -->
```

## IPV6オプション

物理ネットワークのIPv6の設定には以下のオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_physical-ipv6 start -->
    :end-before: <!-- config group network_physical-ipv6 end -->
```

## OVNオプション

物理ネットワークをOVNアップリンクとして使用する際に適用するオプションには以下があります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_physical-ovn start -->
    :end-before: <!-- config group network_physical-ovn end -->
```

## 共通オプション

これらは他の機能に関係なくすべての物理ネットワークに適用されるオプションです:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_physical-common start -->
    :end-before: <!-- config group network_physical-common end -->
```

(network-physical-features)=
## サポートされている機能

物理ネットワークタイプでは以下の機能がサポートされています:

- {ref}`network-bgp`
