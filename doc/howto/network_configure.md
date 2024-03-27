(network-configure)=
# ネットワークを設定する

既存のネットワークを設定するには、 `incus network set` と `incus network unset` コマンド（単一の設定項目を設定する場合）または `incus network edit` コマンド（設定全体を編集する場合）のどちらかを使います。
特定のクラスタメンバーの設定を変更するには、 `--target` フラグを追加してください。

たとえば、以下のコマンドは物理ネットワークの DNS サーバーを設定します:

```bash
incus network set UPLINK dns.nameservers=8.8.8.8
```

利用可能な設定オプションはネットワークタイプによって異なります。
各ネットワークタイプの設定オプションへのリンクは {ref}`network-types` を参照してください。

高度なネットワーク機能を設定するためには別のコマンドがあります。
以下のドキュメントを参照してください:

- {doc}`/howto/network_acls`
- {doc}`/howto/network_forwards`
- {doc}`/howto/network_integrations`
- {doc}`/howto/network_load_balancers`
- {doc}`/howto/network_zones`
- {doc}`/howto/network_ovn_peers`（OVN のみ）
