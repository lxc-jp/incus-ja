(instances-manage)=
# インスタンスを管理するには

既存のインスタンスを一覧表示する際、インスタンスタイプ、状態、場所（場所を持つ場合）を表示できます。
またインスタンスをフィルターして興味があるインスタンスだけを表示できます。

````{tabs}
```{group-tab} CLI
全てのインスタンスを一覧表示するには以下のコマンドを入力します:

    incus list

表示するインスタンスをフィルターできます。例えば、インスタンスタイプ、状態、またはインスタンスが配置されているクラスタメンバーでフィルターできます:

    incus list type=container
    incus list status=running
    incus list location=server1

インスタンス名でフィルターもできます。
複数のインスタンスを一覧表示するには、名前の正規表現を使います。
例えば以下のようにします:

    incus list ubuntu.*

全てのフィルターオプションを見るには [`incus list --help`](incus_list.md) と入力します。
```

```{group-tab} API
すべてのインスタンスを一覧表示するには `/1.0/instances` エンドポイントに問い合わせます。
インスタンスのより詳細な情報を表示するには {ref}`rest-api-recursion` を使えます:

    incus query /1.0/instances?recursion=2

表示するインスタンスを名前、インスタンスタイプ、状態またはインスタンスが配置されているクラスタメンバで {ref}`フィルタ <rest-api-filtering>` できます:

    incus query /1.0/instances?filter=name+eq+ubuntu
    incus query /1.0/instances?filter=type+eq+container
    incus query /1.0/instances?filter=status+eq+running
    incus query /1.0/instances?filter=location+eq+server1

複数のインスタンスを一覧表示するには、名前の正規表現を使います。
たとえば:

    incus query /1.0/instances?filter=name+eq+ubuntu.*

詳細は [`GET /1.0/instances`](swagger:/instances/instances_get) を参照してください。
```

```{group-tab} UI
全てのインスタンスを一覧表示するには {guilabel}`Instances` へ移動します。

表示するインスタンスを、ステータス、インスタンスタイプ、使用しているプロファイルでフィルタできます。

さらに、検索テキストを入力してインスタンスを検索できます。
入力したテキストは名前、説明、ベースイメージの名前にマッチされます。
```
````

## インスタンスの詳細情報を表示する

````{tabs}
```{group-tab} CLI
インスタンスの詳細情報を表示するには以下のコマンドを入力します:

    incus info <instance_name>

インスタンスの最新のログを表示するにはコマンドに `--show-log` を追加します:

    incus info <instance_name> --show-log
```

```{group-tab} API
インスタンスの詳細情報を表示するには以下のエンドポイントに問い合わせます:

    incus query /1.0/instances/<instance_name>

詳細は [`GET /1.0/instances/{name}`](swagger:/instances/instance_get) を参照してください。
```

```{group-tab} UI
overview でインスタンスの行をクリックするとインスタンス一覧の右にインスタンスの情報のサマリが表示されます。

インスタンス名をクリックするとインスタンス詳細のページに遷移し、インスタンスについての詳細情報が表示されます。
```
````

## インスタンスを起動する

````{tabs}
```{group-tab} CLI
インスタンスを起動するには以下のコマンドを入力します:

    incus start <instance_name>

インスタンスが存在しないか既に稼働中の場合はエラーになります。

起動する際にコンソールにすぐにアタッチするには `--console` フラグを渡します。
例えば以下のようにします:

    incus start <instance_name> --console

詳細は {ref}`instances-console` を参照してください。
```

```{group-tab} API
インスタンスを起動するには、PUT リクエストを送ってインスタンス状態を変更してください:

    incus query --request PUT /1.0/instances/<instance_name>/state --data '{"action":"start"}'

<!-- Include start monitor status -->
API 呼び出しの戻り値には操作 ID が含まれます。これを使って操作の状態を問い合わせできます:

    incus query /1.0/operations/<operation_ID>

インスタンスの状態をモニターするには以下のクエリを使います:

    incus query /1.0/instances/<instance_name>/state

詳細は [`GET /1.0/instances/{name}/state`](swagger:/instances/instance_state_get) と [`PUT /1.0/instances/{name}/state`](swagger:/instances/instance_state_put) を参照してください。
<!-- Include end monitor status -->
```

```{group-tab} UI
インスタンスを起動するには、インスタンス一覧か対応するインスタンスのページに移動し、{guilabel}`Start`ボタン(▷)をクリックします。

インスタンス一覧で同時に選択し画面上部の{guilabel}`Start`ボタンをクリックすることで複数のインスタンスを起動することもできます。

インスタンス詳細ページで、{guilabel}`Console`タブを選択するとブートログがインスタンスが起動したときの情報とともに表示されます。
インスタンスが起動したら、{guilabel}`Terminal`タブを選択するとインスタンスにアクセスできます。
```
````

(instances-manage-stop)=
## インスタンスを停止する

`````{tabs}
````{group-tab} CLI
インスタンスを停止するには以下のコマンドを入力します:

    incus stop <instance_name>

インスタンスが存在しないか稼働中ではない場合はエラーになります。
````

````{group-tab} API
インスタンスを停止するには、PUT リクエストを送ってインスタンス状態を変更します:

    incus query --request PUT /1.0/instances/<instance_name>/state --data '{"action":"stop"}'

% Include content from above
```{include} ./instances_manage.md
    :start-after: <!-- Include start monitor status -->
    :end-before: <!-- Include end monitor status -->
```
````

````{group-tab} UI
インスタンスを停止するには、インスタンス一覧か対応するインスタンスのページに遷移し、 {guilabel}`Stop` ボタン (□)をクリックします。
確認のプロンプトが表示されます。

<!-- Include start skip confirmation -->
```{tip}
確認のプロンプトをスキップするには、{kbd}`Shift`キーを押しながらクリックします。
```
<!-- Include end skip confirmation -->

インスタンスを強制停止することもできます。
インスタンスの停止が長い時間かかったり、インスタンスが停止要求に応答しない場合、回転している停止ボタンをクリックし、確認プロンプトに戻ると、そこでインスタンスの強制停止を選択できます。

またインスタンス一覧で複数のインスタンスを選択し、画面上部の {guilabel}`Stop` ボタンを押すことで複数のインスタンスを停止することもできます。
````

`````

## インスタンスを削除する

インスタンスがもう不要な場合、削除できます。
削除する前にインスタンスを停止する必要があります。

`````{tabs}
```{group-tab} CLI
インスタンスを削除するには以下のコマンドを入力します:

    incus delete <instance_name>
```

```{group-tab} API
インスタンスを削除するには、インスタンスに DELETE リクエストを送ります:

    incus query --request DELETE /1.0/instances/<instance_name>

詳細は [`DELETE /1.0/instances/{name}`](swagger:/instances/instance_delete) を参照してください。
```

````{group-tab} UI
インスタンスを削除するには、インスタンスの詳細ページに遷移し {guilabel}`Delete instance` をクリックします。
確認のプロンプトが表示されます。

% Include content from above
```{include} ./instances_manage.md
    :start-after: <!-- Include start skip confirmation -->
    :end-before: <!-- Include end skip confirmation -->
```

またインスタンス一覧で複数のインスタンスを選択し、画面上部の {guilabel}`Delete` ボタンを押すことで複数のインスタンスを削除することもできます。
````
`````

```{caution}
このコマンドはインスタンスとそのスナップショットを永久的に削除します。
```

### 間違ってインスタンスを削除するのを防ぐ

間違ってインスタンスを削除するのを防ぐにはいくつかの異なる方法があります:

- 特定のインスタンスが削除されることを防ぐためには、そのインスタンスの {config:option}`instance-security:security.protection.delete` を `true` に設定します。
  手順は {ref}`instances-configure` を参照してください。
- CLI クライアントでは、[`incus delete`](incus_delete.md) コマンドを使うたびに確認のプロンプトを表示するようなエイリアスを作成します:

       incus alias add delete "delete -i"

## インスタンスを再構築する

インスタンスの root ディスクを一掃して再初期化したいがインスタンスの設定は維持したい場合、インスタンスを再構築できます。

再構築はスナップショットが 1 つも存在しないインスタンスでのみ可能です。

再構築の前にインスタンスを停止します。

````{tabs}
```{group-tab} CLI
別のイメージでインスタンスを再構築するには以下のコマンドを入力します:

    incus rebuild <image_name> <instance_name>

空のルートディスクでインスタンスを再構築するには以下のコマンドを入力します:

    incus rebuild <instance_name> --empty

`rebuild` コマンドについての詳細は、 [`incus rebuild --help`](incus_rebuild.md) を参照してください。
```

```{group-tab} API
インスタンスを別のイメージで再構築するには、インスタンスの `rebuild` エンドポイントに POST リクエストを送ります。
たとえば:

    incus query --request POST /1.0/instances/<instance_name>/rebuild --data '{"source": {"alias":"<image_alias>","server":"<server_URL>", protocol:"simplestreams"}}'

空のルートディスクでインスタンスを再構築するには、 source の type を `none` にします:

    incus query --request POST /1.0/instances/<instance_name>/rebuild --data '{"source": {"type":"none"}}'

詳細は [`POST /1.0/instances/{name}/rebuild`](swagger:/instances/instance_rebuild_post) を参照してください。
```

```{group-tab} UI
インスタンスの再構築は UI では現状ではサポートされていません。
```
````
