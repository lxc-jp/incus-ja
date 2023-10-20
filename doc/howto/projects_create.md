(projects-create)=
# プロジェクトを作成し設定するには

プロジェクトは作成時または後で設定することができます。
ただし、プロジェクトにインスタンスが含まれている場合、有効になっている機能を変更することはできません。

## プロジェクトを作成する

プロジェクトを作成するには、[`incus project create`](incus_project_create.md) コマンドを使用します。

`--config`フラグを使用して設定オプションを指定できます。
利用可能な設定オプションについては、{ref}`ref-projects`を参照してください。

たとえば、インスタンスを分離し、デフォルトプロジェクトのイメージとプロファイルにアクセスを許可する`my-project`というプロジェクトを作成するには、次のコマンドを入力します:

    incus project create my-project --config features.images=false --config features.profiles=false

セキュリティーに関する機能（たとえば、コンテナのネスト）へのアクセスをブロックし、バックアップを許可する`my-restricted-project`というプロジェクトを作成するには、次のコマンドを入力します:

    incus project create my-restricted-project --config restricted=true --config restricted.backups=allow

```{tip}
設定オプションを指定せずにプロジェクトを作成する場合、{config:option}`project-features:features.profiles`は`true`に設定されます。これはプロジェクト内でプロファイルは隔離されることを意味します。

その結果、新しいプロジェクトは`default`プロジェクトの`default`プロファイルへのアクセスは持たず、そのため（ルートディスクのような）インスタンス作成に必要な設定が不足します。
これを修正するためには、[`incus profile device add`](incus_profile_device_add.md)コマンドを使用してプロジェクトの`default`プロファイルにルートディスクデバイスを追加してください。
```

(projects-configure)=
## プロジェクトの設定
プロジェクトを設定するには、特定の設定オプションを設定するか、プロジェクト全体を編集できます。

いくつかの設定オプションは、インスタンスが含まれていないプロジェクトに対してのみ設定できます。

### 特定の設定オプションを設定する

特定の設定オプションを設定するには、[`incus project set`](incus_project_set.md) コマンドを使用します。

たとえば、`my-project`で作成できるコンテナの数を 5 つに制限するには、次のコマンドを入力します:

    incus project set my-project limits.containers=5

特定の設定オプションを解除するには、[`incus project unset`](incus_project_unset.md) コマンドを使用します。

```{note}
設定オプションを解除すると、デフォルト値に設定されます。
このデフォルト値は、プロジェクトが作成されたときに設定される初期値と異なる場合があります。
```

### プロジェクトを編集する

プロジェクトの設定全体を編集するには、[`incus project edit`](incus_project_edit.md) コマンドを使用します。

たとえば:

    incus project edit my-project
