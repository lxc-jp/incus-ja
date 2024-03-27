(network-integrations)=
# ネットワーク統合を設定するには

```{note}
ネットワーク統合は現状 {ref}`network-ovn` でのみ利用できます。
```

ネットワーク統合はローカル環境の Incus と Incus や他のプラットフォームにホストされたリモートのネットワークを接続するのに使えます。

## OVN 相互接続

現時点ではサポートされるネットワーク統合のタイプは OVN のみです。
これは OVN 相互接続ゲートウェイを複数の環境にまたがって OVN ネットワークをピアリングするために使えます。

これを動作させるには以下のような構成で稼働中の OVN 相互接続が必要です:

- OVN 相互接続 `NorthBound` と `SouthBound` データベース
- アベイラビリティーゾーン名 (`name` プロパティー)が適切に設定された 2 つ以上の OVN クラスター
- すべての OVN クラスターは `ovn-ic` デーモンが稼働している必要あり
- OVN クラスターが相互接続にルートを広告し学習するように設定されている
- 少なくとも 1 つのサーバーが OVN 相互接続ゲートウェイとして設定されている

さらなる詳細は [upstream のドキュメント](https://docs.ovn.org/en/latest/tutorials/ovn-interconnection.html) にあります。

## ネットワーク統合を作成

ネットワーク統合は `incus network integration create` で作成できます。
統合は Incus 環境にグローバルで、ネットワークやプロジェクトには紐づいていません。

OVN 統合の例は以下のようなものになるでしょう:

```
incus network integration create ovn-region ovn
incus network integration set ovn-region ovn.northbound_connection tcp:[192.0.2.12]:6645,tcp:[192.0.3.13]:6645,tcp:[192.0.3.14]:6645
incus network integration set ovn-region ovn.southbound_connection tcp:[192.0.2.12]:6646,tcp:[192.0.3.13]:6646,tcp:[192.0.3.14]:6646
```

## ネットワーク統合を使用

ネットワーク統合を使用するには、ピアリングする必要があります。

これは `incus network peer create` で行えます。例えば:

```
incus network peer create default region ovn-region --type=remote
```

## 設定オプション

以下の設定オプションがすべてのネットワーク統合で利用できます:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_integration-common start -->
    :end-before: <!-- config group network_integration-common end -->
```

### OVN 設定オプション

これらのオプションは OVN ネットワーク統合に固有です:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_integration-ovn start -->
    :end-before: <!-- config group network_integration-ovn end -->
```
