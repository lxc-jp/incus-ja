(network-external)=
# 外部ネットワーク

<!-- Include start external intro -->
外部ネットワークは既に存在するネットワークを使用します。
そのため、 Incus がそれらを制御するには限界があるため、ネットワーク ACL、ネットワークフォワードやネットワークゾーンのような Incus の機能はサポートされません。

外部ネットワークを使用する主な目的は親インターフェースによるアップリンクのネットワークを提供することです。
この外部ネットワークはインスタンスや他のネットワークを親のインターフェースに接続する際のプリセットを指定します。

Incus は以下の外部ネットワークタイプをサポートします:
<!-- Include end external intro -->

```{toctree}
:maxdepth: 1
/reference/network_macvlan
/reference/network_sriov
/reference/network_physical
```
