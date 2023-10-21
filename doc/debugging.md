# Incusをデバッグするには

インスタンスの問題をデバッグする際の情報については、{ref}`instances-troubleshoot`を参照してください。

## `incus` と `incusd` のデバッグ

`incus` と `incusd` のコードをトラブルシューティングするのに役立ついくつかの異なる方法を説明します。

### `incus --debug`

クライアントのどのコマンドにも `--debug` フラグを追加することで内部についての追加情報を出力することができます。
もし有用な情報がない場合はログ出力の呼び出しで追加することができます:

    logger.Debugf("Hello: %s", "Debug")

### `incus monitor`

このコマンドはメッセージがリモートのサーバーに現れるのをモニターします。

## ローカルソケット経由でのREST API

サーバーサイドで Incus とやりとりするのに最も簡単な方法はローカルソケットを経由することです。
以下のコマンドは`GET /1.0`にアクセスし、[jq](https://stedolan.github.io/jq/tutorial/)ユーティリティを使って
JSON を人間が読みやすいように整形します:

```bash
curl --unix-socket /var/lib/incus/unix.socket incus/1.0 | jq .
```

利用可能な API については[RESTful API](rest-api.md)をご参照ください。

## HTTPS経由でのREST API

[Incus への HTTPS 接続](security.md)には、有効なクライアント証明書が必要で、最初の [`incus remote add`](incus_remote_add.md) で生成されます。
この証明書は、認証と暗号化のための接続ツールに渡す必要があります。

必要に応じて、`openssl`を使って証明書（`~/.config/incus/client.crt`）を調べることができます:

```bash
openssl x509 -text -noout -in client.crt
```

表示される行の中に以下のようなものがあるはずです:

    Certificate purposes:
    SSL client : Yes

### コマンドラインツールを使う

```bash
wget --no-check-certificate --certificate=$HOME/.config/incus/client.crt --private-key=$HOME/.config/incus/client.key -qO - https://127.0.0.1:8443/1.0
```

### ブラウザを使う

いくつかのブラウザ拡張はウェブのリクエストを作成、修正、リプレイするための便利なインターフェースを提供しています。
Incus サーバーに対して認証するには`incus`のクライアント証明書をインポート可能な形式に変換しブラウザにインポートしてください。

たとえば Windows で利用可能な形式の`client.pfx`を生成するには以下のようにします:

```bash
openssl pkcs12 -clcerts -inkey client.key -in client.crt -export -out client.pfx
```

上記のコマンドを実行し、（訳注：変換後の証明書をインポートしてから）ブラウザで[`https://127.0.0.1:8443/1.0`](https://127.0.0.1:8443/1.0)を開けば期待通り動くはずです。

## Incusデータベースをデバッグ

グローバル{ref}`データベース <database>`のファイルは Incus のデータディレクトリー（`/var/lib/incus/database/global`）の`./database/global`サブディレクトリーの下に格納されます。

クラスタの各メンバーもそのメンバー固有の何らかのデータを保持する必要があるため、Incus は通常の SQLite のデータベース(「ローカル」データベース)も使用します。
これは`./database/local.db`に置かれます。

アップグレードの前にはグローバルデータベースのディレクトリーとローカルデータベースのファイルのバックアップが作成され、 `.bak` のサフィックス付きでタグ付けされます。
アップグレード前の状態に戻す必要がある場合は、このバックアップを使うことができます。

### データベースのデータとスキーマをダンプする

データベースのデータまたはスキーマの SQL テキスト形式でのダンプを取得したい場合は、`incus admin sql <local|global> [.dump|.schema]` コマンドを使ってください。
これにより`sqlite3`コマンドラインツールの`.dump`または`.schema`ディレクティブと同じ出力を生成できます。

### コンソールからカスタムクエリを実行する

ローカルまたはグローバルデータベースに SQL クエリ（例 `SELECT`, `INSERT`, `UPDATE`）を実行する必要がある場合、`incus admin sql`コマンドを使うことができます（詳細は`incus admin sql --help`を実行してください）。

ただ、これが必要になるのは壊れたアップデートかバグからリカバーするときだけでしょう。
その場合、まず Incus チームに相談してみてください（[GitHubのイシュー](https://github.com/lxc/incus/issues/new)を作成するか[フォーラム](https://discuss.linuxcontainers.org/)に投稿）。

### Incusデーモン起動時にカスタムクエリを実行する

SQL のデータマイグレーションのバグあるいは関連する問題のためにアップグレード後に Incus デーモンが起動に失敗する場合、壊れたアップデートを修復するクエリを含んだ`.sql`ファイルを作成することで、その状況からリカバーできる可能性があります。

ローカルデータベースに対して修復を実行するには、修復に必要なクエリを含む`./database/patch.local.sql`というファイルを作成してください。
同様にグローバルデータベースの修復には`./database/patch.global.sql`というファイルを作成してください。

これらのファイルはデーモンの起動シーケンスの非常に早い段階で読み込まれ、クエリが成功したときは削除されます（クエリは SQL トランザクション内で実行されるので、クエリが失敗したときにデータベースの状態が変更されることはありません）。

上記のとおり、まず Incus チームに相談してみてください。

### クラスタデータベースをディスクに同期

クラスタデータベースの内容をディスクにフラッシュしたいなら、`incus admin sql global .sync`コマンドを使ってください。これは通常の SQLite 形式のデータベースのファイルを`./database/global/db.bin`に書き込みます。
その後`sqlite3`コマンドラインツールを使って中身を見ることができます。
