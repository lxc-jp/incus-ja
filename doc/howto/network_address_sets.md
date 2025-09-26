(network-address-sets)=
# ネットワークアドレスセットを使うには

```{note}
ネットワークアドレスセットは{ref}`ACL <network-acls>`で使用します。これは{ref}`network-ovn`か`nftables`を使う{ref}`ブリッジネットワーク <network-bridge-firewall>`でのみ動作します。
```

ネットワークアドレスセットはCIDRサフィックスありまたはなしのIPv4かIPv6のアドレスのリストです。これらは{ref}`ACL <network-acls-rules-properties>`のsourceまたはdestinationフィールドで使えます。

## アドレスセットのプロパティ

アドレスセットには以下のプロパティがあります:

| プロパティ    | 型          | 必須 | 説明                             |
| :---          | :---        | :--- | :---                             |
| `name`        | string      | yes  | ネットワークアドレス設定の名前   |
| `description` | string      | no   | ネットワークアドレス設定の説明   |
| `addresses`   | string list | no   | イングレス・トラフィックのルール |

## アドレスセットの設定オプション

以下の設定オプションがすべてのネットワークアドレスセットで利用できます:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group network_address_set-common start -->
    :end-before: <!-- config group network_address_set-common end -->
```

## アドレスセットの作成

以下のコマンドでアドレスセットを作成します。

```bash
incus network address-set create <name> [configuration_options...]
```

これはアドレスなしでアドレスセットを作成します。後から{ref}`アドレスを追加 <manage-addresses-in-set>`できます。

(manage-addresses-in-set)=
## アドレスを追加または削除

アドレスの追加はとても簡単です:

```bash
incus network address-set add <name> <address1> <address2>
```

セットに追加するアドレスの種類に制限はありません。IPv4とIPv6のアドレスとCIDRを混合して一度に追加できます。

アドレスを削除するには、`add`コマンドの代わりに`remove`コマンドが使えます。

```bash
incus network address-set remove <name> <address1> <address2>
```

## ACLルール内でのアドレスセットの使用

{ref}`ACL <network-acls-address-sets>`内でアドレスセットを使うには、`name`の前に`$`（コマンドラインではドルをエスケープする必要があります）を追加する必要があります。こうすることで、ACLルールの`source`または`destination`フィールド内でアドレスセットを参照できます。
