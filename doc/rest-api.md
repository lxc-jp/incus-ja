# REST API

Incus とクライアント間のすべての通信は HTTP 上の RESTful API を使用します。
この API や（リモートの通信では）TLS あるいは（ローカルの操作では）Unix ソケットで通信します。

どのようにリモートの API にアクセスするかについての情報は{ref}`authentication`を参照してください。

```{tip}
- どのように API が使われるかの例を見るにはIncusクライアント（[`incus`](incus.md)）で`--debug`フラグを追加してコマンドを実行してください。
デバッグ情報がAPIの呼び出しと戻り値を表示します。
- 手軽にに API を呼び出せるように、Incus クライアントは [`incus query`](incus_query.md) コマンドを提供しています。
```

## API のバージョニング

サポートされている API のメジャーバージョンのリストは `GET /` を使って取得できます。

後方互換性を壊す場合は API のメジャーバージョンが上がります。

後方互換性を壊さずに追加される機能は `api_extensions` の追加という形になり、
特定の機能がサーバーでサポートされているかクライアントがチェックすることで
利用できます。

## 戻り値

次の 3 つの標準的な戻り値の型があります。

* 標準の戻り値
* バックグラウンド操作
* エラー

### 標準の戻り値

標準の同期的な操作に対しては以下のような dict が返されます:
```js
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "metadata": {}                          // リソースやアクションに固有な追加のメタデータ
}
```

HTTP ステータスコードは必ず 200 です。

### バックグラウンド操作

リクエストの結果がバックグラウンド操作になる場合、 HTTP ステータスコードは 202（Accepted）
になり、操作の URL を指す HTTP の Location ヘッダが返されます。

レスポンスボディは以下のような構造を持つ dict です:

```js
{
    "type": "async",
    "status": "OK",
    "status_code": 100,
    "operation": "/1.0/instances/<id>",                     // バックグラウンド操作の URL
    "metadata": {}                                          // 操作のメタデータ（下記参照）
}
```

操作のメタデータの構造は以下のようになります:

```js
{
    "id": "a40f5541-5e98-454f-b3b6-8a51ef5dbd3c",           // 操作の UUID
    "class": "websocket",                                   // 操作の種別（task, websocket, token のいずれか）
    "created_at": "2015-11-17T22:32:02.226176091-05:00",    // 操作の作成日時
    "updated_at": "2015-11-17T22:32:02.226176091-05:00",    // 操作の最終更新日時
    "status": "Running",                                    // 文字列表記での操作の状態
    "status_code": 103,                                     // 整数表記での操作の状態（status ではなくこちらを利用してください。訳注: 詳しくは下記のステータスコードの項を参照）
    "resources": {                                          // リソース種別（container, snapshots, images のいずれか）の dict を影響を受けるリソース
      "instances": [
        "/1.0/instances/test"
      ]
    },
    "metadata": {                                           // 対象となっている（この例では exec）操作に固有なメタデータ
      "fds": {
        "0": "2a4a97af81529f6608dca31f03a7b7e47acc0b8dc6514496eb25e325f9e4fa6a",
        "control": "5b64c661ef313b423b5317ba9cb6410e40b705806c28255f601c0ef603f079a7"
      }
    },
    "may_cancel": false,                                    //（REST で DELETE を使用して）操作がキャンセル可能かどうか
    "err": ""                                               // 操作が失敗した場合にエラー文字列が設定されます
}
```

対象の操作に対して追加のリクエストを送って情報を取り出さなくても、
何が起こっているかユーザーにとってわかりやすい形でボディは構成されています。
ボディに含まれるすべての情報はバックグラウンド操作の URL から取得する
こともできます。

### エラー

さまざまな状況によっては操作を行う前に直ぐに問題が起きる場合があり、
そういう場合には以下のような値が返されます:

```js
{
    "type": "error",
    "error": "Failure",
    "error_code": 400,
    "metadata": {}                      // エラーについてのさらなる詳細
}
```

HTTP ステータスコードは 400、401、403、404、409、412、500 のいずれかです。

## ステータスコード
Incus REST API はステータス情報を返す必要があります。それはエラーの理由だったり、
操作の現在の状態だったり、 Incus が提供する様々なリソースの状態だったりします。

デバッグをシンプルにするため、ステータスは常に文字列表記と整数表記で
重複して返されます。ステータスの整数表記の値は将来に渡って不変なので
API クライアントが個々の値に依存できます。文字列表記のステータスは
人間が API を手動で実行したときに何が起きているかをより簡単に判断
できるように用意されています。

ほとんどのケースでこれらは `status` と `status_code` と呼ばれ、前者は
ユーザーフレンドリーな文字列表記で後者は固定の数値です。

整数表記のコードは常に 3 桁の数字で以下の範囲の値となっています。

* 100 to 199: リソースの状態（started、stopped、ready、…）
* 200 to 399: 成功したアクションの結果
* 400 to 599: 失敗したアクションの結果
* 600 to 999: 将来使用するために予約されている番号の範囲

### 現在使用されているステータスコード一覧

コード | 意味
:---   | :------
100    | 操作が作成された
101    | 開始された
102    | 停止された
103    | 実行中
104    | キャンセル中
105    | ペンディング
106    | 開始中
107    | 停止中
108    | 中断中
109    | 凍結中
110    | 凍結された
111    | 解凍された
112    | エラー
113    | 準備完了
200    | 成功
400    | 失敗
401    | キャンセルされた

(rest-api-recursion)=
## 再帰

巨大な一覧のクエリを最適化するために、コレクションに対して再帰が実装されています。
コレクションに対するクエリの GET リクエストに `recursion` パラメーターを指定できます。

デフォルト値は 0 でコレクションのメンバーの URL が返されることを意味します。
1 を指定するとこれらの URL がそれが指すオブジェクト（通常は dict 形式）で
置き換えられます。

再帰はジョブへのポインタ（URL）をオブジェクトそのもので単に置き換えるように
実装されています。

(rest-api-filtering)=
## フィルタ

検索結果をある値でフィルタするために、コレクションにフィルタが実装されています。
コレクションに対する GET クエリに `filter` 引数を渡せます。

フィルタにはデフォルト値はありません。これは見つかったすべての結果が返されることを意味します。
フィルタの引数には以下のような言語を設定します。

    ?filter=field_name eq desired_field_assignment

この言語は REST API のフィルタロジックを構成するための OData の慣習に従います。
フィルタは下記の論理演算子もサポートします。
not（`not`）、equals（`eq`）、not equals（`ne`）、and（`and`）、or（`or`）
フィルタは左結合で評価されます。
空白を含む値はクォートで囲むことができます。
ネストしたフィルタもサポートされます。
たとえば設定内のフィールドに対してフィルタするには以下のように指定します:

    ?filter=config.field_name eq desired_field_assignment

device の属性についてフィルタするには以下のように指定します:

    ?filter=devices.device_name.field_name eq desired_field_assignment

以下に上記の異なるフィルタの方法を含む GET クエリをいくつか示します:

    containers?filter=name eq "my container" and status eq Running

    containers?filter=config.image.os eq ubuntu or devices.eth0.nictype eq bridged

    images?filter=Properties.os eq Centos and not UpdateSource.Protocol eq simplestreams

## 非同期操作

完了までに 1 秒以上かかるかもしれない操作はバックグラウンドで実行しなければ
なりません。そしてクライアントにはバックグラウンド操作 ID を返します。

クライアントは操作のステータス更新をポーリングするか long-poll API を使って
通知を待つことが出来ます。

## 通知

通知のために WebSocket ベースの API が利用できます。クライアントへ送られる
トラフィックを制限するためにいくつかの異なる通知種別が存在します。

リモート操作の状態をポーリングしなくて済むように、リモート操作を開始する
前に操作の通知をクライアントが常に購読しておくのがお勧めです。

## PUT と PATCH の使い分け

Incus API は既存のオブジェクトを変更するのに PUT と PATCH の両方をサポートします。

PUT はオブジェクト全体を新しい定義で置き換えます。典型的には GET で現在の
オブジェクトの状態を取得した後に PUT が呼ばれます。

レースコンディションを避けるため、 GET のレスポンスから ETag ヘッダを読み取り
PUT リクエストの If-Match ヘッダに設定するべきです。こうしておけば GET と
PUT の間にオブジェクトが他から変更されていた場合は更新が失敗するようになります。

PATCH は変更したいプロパティだけを指定することでオブジェクト内の単一の
フィールドを変更するのに用いられます。キーを削除するには通常は空の値を
設定すれば良いようになっていますが、 PATCH ではキーの削除は出来ず、代わりに
PUT を使う必要がある場合もあります。

## API 構造

Incus は API エンドポイントを記述する [Swagger](https://swagger.io/) 仕様を自動生成しています。
この API 仕様の YAML 版が [`rest-api.yaml`](https://github.com/lxc/incus/blob/main/doc/rest-api.yaml) にあります。
手軽にウェブで見る場合は {doc}`api` を参照してください。
