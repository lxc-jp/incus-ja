(howto-storage-buckets)=
# ストレージバケットとキーを管理するには

{ref}`storage-buckets` を作成、設定、表示、リサイズするための手順およびストレージバケットキーを管理する方法については以下のセクションを参照してください。

## S3アドレスを設定する

S3 アドレスを設定することにより、ローカルストレージ（`dir`、`btrfs`、`lvm`、または`zfs`プール）上のストレージバケットを使用することが可能になります。これにより、S3 プロトコルを通じてバケットにアクセスできるようになります。

S3 アドレスを設定するには、{config:option}`server-core:core.storage_buckets_address` サーバー設定オプションを設定します。たとえば:

    incus config set core.storage_buckets_address :8555

## ストレージバケットを管理する

ストレージバケットは S3 プロトコルを使って公開されるオブジェクトストレージを提供します。

カスタムストレージボリュームとは異なり、ストレージバケットはインスタンスに追加されるのではなく、それらの URL を通してアプリケーションから直接アクセスされます。

詳細は {ref}`storage-buckets` を参照してください。

### ストレージバケットを作成する

ストレージプール内にストレージバケットを作成するには、以下のコマンドを使用します:

    incus storage bucket create <pool_name> <bucket_name> [configuration_options...]

それぞれのドライバーで利用可能なストレージバケット設定オプションの一覧については {ref}`storage-drivers` を参照してください。

クラスタメンバーにストレージバケットを追加するには `--target` フラグを追加してください:

    incus storage bucket create <pool_name> <bucket_name> --target=<cluster_member> [configuration_options...]

```{note}
ほとんどのストレージドライバでは、ストレージバケットはクラスタ間でリプリケートされず、作成されたメンバー上にのみ存在します。
この挙動は `cephobject` ストレージプールでは異なります。 `cephobject` ではバケットはどのクラスタメンバーからも利用できます。
```

### ストレージバケットを設定するには

各ストレージドライバーで利用可能な設定オプションについては {ref}`storage-drivers` ドキュメントを参照してください。

ストレージバケットの設定オプションを設定するには以下のコマンドを使用します:

    incus storage bucket set <pool_name> <bucket_name> <key> <value>

たとえば、バケットにクォータサイズを設定するには、以下のコマンドを使用します:

    incus storage bucket set my-pool my-bucket size 1MiB

以下のコマンドでストレージバケットの設定を編集することもできます:

    incus storage bucket edit <pool_name> <bucket_name>

ストレージバケットとそのキーを削除するには以下のコマンドを使用します:

    incus storage bucket delete <pool_name> <bucket_name>

### ストレージバケットを表示するには

ストレージプール内のすべての利用可能なストレージバケットの一覧を表示し設定を確認できます。

ストレージプール内のすべての利用可能なストレージバケットを一覧表示するには、以下のコマンドを使用します:

    incus storage bucket list <pool_name>

特定のバケットの詳細情報を表示するには、以下のコマンドを使用します:

    incus storage bucket show <pool_name> <bucket_name>

### ストレージバケットをリサイズするには

デフォルトではストレージバケットにはクォータは適用されません。

ストレージバケットクォータを設定するには、サイズを設定します:

    incus storage bucket set <pool_name> <bucket_name> size <new_size>

```{important}
- ストレージバケットの拡大は通常は正常に動作します（ストレージプールが十分なストレージを持つ場合）。
- ストレージバケットを現在の使用量より縮小することはできません。

```

## ストレージバケットキーを管理する

アプリケーションがストレージバケットにアクセスするためには *アクセスキー* と *シークレットキー* からなる S3 クレデンシャルを使う必要があります。
特定のバケットに対して複数のセットのクレデンシャルを作成できます。

それぞれのクレデンシャルのセットにはキー名を設定します。
キー名は参照のためだけに用いられ、アプリケーションがクレデンシャルを使用する際に提供する必要はありません。

それぞれのクレデンシャルのセットには *ロール* が設定されます。それはバケットにどの操作を実行できるかを指定します。

使用可能なロールは以下のとおりです:

- `admin` - バケットへのフルアクセス。
- `read-only` - バケットへの読み取り専用アクセス（一覧とファイルの取得のみ）。

バケットキー作成時にロールが指定されない場合、使用されるロールは `read-only` になります。

### ストレージバケットキーを作成する

ストレージバケットにクレデンシャルのセットを作成するには、以下のコマンドを使用します:

    incus storage bucket key create <pool_name> <bucket_name> <key_name> [configuration_options...]

ストレージバケットに特定のロールを持つクレデンシャルのセットを作成するには、以下のコマンドを使用します:

    incus storage bucket key create <pool_name> <bucket_name> <key_name> --role=admin [configuration_options...]

これらのコマンドはランダムなクレデンシャルキーのセットを生成し表示します。

### ストレージバケットキーを編集または削除するには

既存のバケットキーを編集するには以下のコマンドを使用します:

    incus storage bucket key edit <pool_name> <bucket_name> <key_name>

既存のバケットキーを削除するには以下のコマンドを使用します:

    incus storage bucket key delete <pool_name> <bucket_name> <key_name>

### ストレージバケットのキーを表示するには

既存のバケットに定義されているキーを表示するには以下のコマンドを使用します:

    incus storage bucket key list <pool_name> <bucket_name>

特定のバケットキーを表示するには以下のコマンドを使用します:

    incus storage bucket key show <pool_name> <bucket_name> <key_name>
