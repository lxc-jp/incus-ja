(network-ovn-peers)=
# ピアルーティング関係を作成するには

デフォルトでは、 2 つの OVN ネットワーク間のトラフィックはアップリンクのネットワークを経由します。
しかし、パケットが OVN サブシステムから出てホストのネットワークスタック（そして、場合によっては外部ネットワーク）を通過し対象ネットサークの OVN サブシステムに戻る必要があるため、この経路は非効率です。
ホストのネットワークの構成によっては、利用できる帯域幅が制限される場合があります（ホストの外部ネットワークより OVN のオーバーレイネットワークが広帯域幅のネットワークにある場合）。

このため、 Incus では 2 つの OVN ネットワーク間でルーティング関係を作成できます。
この方法を使うと 2 つのネットワーク間での通信がアップリンクのネットワーク経由ではなく OVN サブシステム内で完結できます。

さらに、ネットワーク統合を使うと、別のクラスター上で動いている場合であっても 2 つの OVN ネットワークをピアリングできます。

## ネットワーク間にルーティング関係を作成する

2 つのネットワーク間にルーティング関係を作成するには、両方のネットワークにネットワークピアを作成する必要があります。
関係は双方向でなくてはなりません。
1 のネットワークだけセットアップした場合、ルーティング関係はペンディング状態になり、アクディブにはなりません。

ピアのルーティング関係を作成する際は、対象のネットワークとの関係を特定するピアの名前を指定します。
名前は自由に選ぶことができ、後で関係を編集または削除する際に使用します。

同じプロジェクト内のネットワーク間でピアのルーティング関係を作成するには次のコマンドを使います:

    incus network peer create <network1> <peering_name> <network2> [configuration_options]
    incus network peer create <network2> <peering_name> <network1> [configuration_options]

別のプロジェクトの OVN ネットワーク間でピアのルーティング関係を作成することもできます:

    incus network peer create <network1> <peering_name> <project2/network2> [configuration_options] --project=<project1>
    incus network peer create <network2> <peering_name> <project1/network1> [configuration_options] --project=<project2>

ネットワーク統合を使ったリモートのピアリングは以下のように作成します:

    incus network peer create <network1> <peering_name> <integration name> [configuration_options] --type=remote

```{important}
プロジェクトまたはネットワークの名前が間違っている場合、コマンドは対応するプロジェクトやネットワークが存在しないというエラーは出さず、ルーティング関係はペンディング状態のままになります。
これは他のプロジェクトのユーザがプロジェクトやネットワークが存在するか調べられないようにするための（訳注: セキュリティ上の）仕様です。
```

### アのプロパティ

ピアのルーティング関係には以下のプロパティがあります。

プロパティ           | 型         | 必須 | 説明
:--                  | :--        | :--  | :--
`name`               | string     | yes  | ローカルネットワーク上のネットワークピアの名前
`description`        | string     | no   | ネットワークピアの説明
`config`             | string set | no   | 設定のキーバリューペアー（`user.*` のカスタムキーのみサポート）
`target_integration` | string     | no   | 統合の名前（リモートピアの作成時に必須）
`target_project`     | string     | yes  | 対象のネットワークがどのプロジェクト内に存在するか（作成時に必須）
`target_network`     | string     | yes  | どのネットワークとピアを作成するか（作成時に必須）
`status`             | string     | --   | 作成中か作成完了（対象のネットワークと相互にピアリングした状態）かを示すステータス

## ルーティング関係の一覧を表示する

ネットワークのネットワークピアすべての一覧を表示するには、次のコマンドを実行します:

    incus network peer list <network>

## ルーティング関係を編集する

ネットワークピアを編集するには次のコマンドを実行します:

    incus network peer edit <network> <peering_name>

このコマンドは YAML 形式のネットワークピアを編集モードで開きます。
