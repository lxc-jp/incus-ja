(dev-incus)=
# インスタンス〜ホスト間の通信

ホストされているワークロード（インスタンス）とそのホストのコミュニケーションは
厳密には必要とされているわけではないですが、とても便利な機能です。

Incus ではこの機能は `/dev/incus/sock` というノードを通して実装されており、
このノードはすべての Incus のインスタンスに対して作成、セットアップされます。

このファイルはインスタンス内部のプロセスが接続できる Unix ソケットです。
マルチスレッドで動いているので複数のクライアントが同時に接続できます。

```{note}
インスタンスのソケットへのアクセスを許可するには {config:option}`instance-security:security.guestapi` を `true`（これがデフォルトです）に設定する必要があります。
```

## 実装詳細

ホストでは Incus は `/var/lib/incus/guestapi/sock` をバインドして新しいコネクションの
リッスンを開始します。

このソケットは、Incus が開始させたすべてのインスタンス内の `/dev/incus/sock` に
公開されます。

4096 を超えるインスタンスを扱うのに単一のソケットが必要です。そうでなければ、
Incus は各々のインスタンスに異なるソケットをバインドする必要があり、
ファイルディスクリプター数の上限にすぐ到達してしまいます。

## 認証

`/dev/incus/sock` への問い合わせは依頼するインスタンスに関連した情報のみを
返します。リクエストがどこから来たかを知るために、 Incus は初期のソケットの
ユーザークレデンシャルを取り出し、 Incus が管理しているインスタンスのリストと比較します。

## プロトコル

`/dev/incus/sock` のプロトコルは JSON メッセージを用いたプレーンテキストの
HTTP であり、 Incus プロトコルのローカル版に非常に似ています。

メインの Incus API とは異なり、 `/dev/incus/sock` API にはバックグラウンド処理と
認証サポートはありません。

## REST-API

### API の構造

* `/`
   * `/1.0`
      * `/1.0/config`
         * `/1.0/config/{key}`
      * `/1.0/devices`
      * `/1.0/events`
      * `/1.0/images/{fingerprint}/export`
      * `/1.0/meta-data`

### API の詳細

#### `/`

##### GET

* 説明: サポートされている API のリスト
* 出力: サポートされている API エンドポイント URL のリスト（デフォルトでは ['/1.0']`）

戻り値:

```json
[
    "/1.0"
]
```

#### `/1.0`

##### GET

* 説明: 1.0 API についての情報
* 出力: dict 形式のオブジェクト

戻り値:

```json
{
    "api_version": "1.0",
    "location": "foo.example.com",
    "instance_type": "container",
    "state": "Started",
}
```

#### PATCH

* 説明: インスタンスの状態を更新する（有効な状態は `Ready` と `Started`）
* 戻り値: 無し

 入力:

 ```json
 {
    "state": "Ready"
 }
```

#### `/1.0/config`

##### GET

* 説明: 設定キーの一覧
* 出力: 設定キー URL のリスト

設定キーの名前はインスタンスの設定の名前と一致するようにしています。
しかし、設定の namespace のすべてが `/dev/incus/sock` にエクスポート
されているわけではありません。
現在は `cloud-init.*` と `user.*` キーのみがインスタンスにアクセス可能となっています。

現時点ではインスタンスが書き込み可能な名前空間はありません。

戻り値:

```json
[
    "/1.0/config/user.a"
]
```

#### `/1.0/config/<KEY>`

##### GET

* 説明: そのキーの値
* 出力: プレーンテキストの値

戻り値:

    blah

#### `/1.0/devices`

##### GET

* 説明: インスタンスのデバイスのマップ
* 出力: dict

戻り値:

```json
{
    "eth0": {
        "name": "eth0",
        "network": "incusbr0",
        "type": "nic"
    },
    "root": {
        "path": "/",
        "pool": "default",
        "type": "disk"
    }
}
```

#### `/1.0/events`

##### GET

* 説明: この API ではプロトコルが WebSocket にアップグレードされます。
* 出力: 無し（イベントのフローが終わることがなくずっと続く）

サポートされる引数は以下の通りです:

* type: 購読する通知の種別のカンマ区切りリスト（デフォルトは all）

通知の種別には以下のものがあります:

* `config`（あらゆる `user.*` 設定キーの変更）
* `device`（あらゆるデバイスの追加、変更、削除）

この API は決して終了しません。それぞれの通知は別々の JSON の dict として送られます:

```json
{
    "timestamp": "2017-12-21T18:28:26.846603815-05:00",
    "type": "device",
    "metadata": {
        "name": "kvm",
        "action": "added",
        "config": {
            "type": "unix-char",
            "path": "/dev/kvm"
        }
    }
}
```

```json
{
    "timestamp": "2017-12-21T18:28:26.846603815-05:00",
    "type": "config",
    "metadata": {
        "key": "user.foo",
        "old_value": "",
        "value": "bar"
    }
}
```

#### `/1.0/images/<FINGERPRINT>/export`

##### GET

* 説明: 公開されたあるいはキャッシュされたイメージをホストからダウンロードする
* 出力: 生のイメージあるいはエラー
* アクセス権: `security.devlxd.images` を `true` に設定する必要があります

戻り値:

    デーモン API の /1.0/images/<FINGERPRINT>/export を参照してください。

#### `/1.0/meta-data`

##### GET

* 説明: cloud-init と互換性のあるコンテナのメタデータ
* 出力: cloud-init のメタデータ

戻り値:

    #cloud-config
    instance-id: af6a01c7-f847-4688-a2a4-37fddd744625
    local-hostname: abc
