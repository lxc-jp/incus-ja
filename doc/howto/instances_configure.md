(instances-configure)=
# インスタンスを設定するには

{ref}`instance-properties` か {ref}`instance-options` を設定するか {ref}`devices` を追加し設定することでインスタンスを設定できます。

設定方法は以下の項を参照してください。

```{note}
異なるインスタンス設定を保管し再利用するには、{ref}`プロファイル <profiles>` を使用してください。
```

(instances-configure-options)=
## インスタンスオプションを設定する

{ref}`インスタンスを作成する <instances-create>` 際にインスタンスオプションを指定できます。
あるいは、インスタンス作成後にインスタンスオプションを変更できます。

````{tabs}
```{group-tab} CLI
[`incus config set`](incus_config_set.md) コマンドを使ってインスタンスオプションを変更できます。
インスタンス名とインスタンスオプションのキーとバリューを指定します:

    incus config set <instance_name> <option_key>=<option_value> <option_key>=<option_value> ...
```

```{group-tab} API
インスタンスに PATCH リクエストを送るとインスタンスオプションを変更します。
インスタンス名とインスタンスオプションのキーとバリューを指定します:

    incus query --request PATCH /1.0/instances/<instance_name> --data '{"config": {"<option_key>":"<option_value>","<option_key>":"<option_value>"}}'

詳細は [`PATCH /1.0/instances/{name}`](swagger:/instances/instance_patch) を参照してください。
```

```{group-tab} UI
インスタンスオプションを変更するには、インスタンス詳細ページの {guilabel}`Configuration` タブに遷移して {guilabel}`Edit instance` をクリックします。

変更したい設定オプションを探して値を変更します。
{guilabel}`Save changes` をクリックして変更した設定を保存してください。

UI に表示されていないインスタンスオプションを変更するには、 {ref}`instances-configure-edit` の手順に従ってください。
```
````

利用可能なオプションの一覧とどのオプションがどのインスタンスタイプで利用可能かの情報は {ref}`instance-options` を参照してください。

例えば、コンテナのメモリーリミットを変更するには:

````{tabs}
```{group-tab} CLI
メモリーリミットを 8 GiB に設定するには、以下のコマンドを入力します:

    incus config set my-container limits.memory=8GiB
```

```{group-tab} API
メモリーリミットを 8 GiB に設定するには、以下のリクエストを送ります:

    incus query --request PATCH /1.0/instances/my-container --data '{"config": {"limits.memory":"8GiB"}}'
```

```{group-tab} UI
メモリーリミットを 8 GiB に設定するには、インスタンス詳細ページの {guilabel}`Configuration` タブに遷移して {guilabel}`Advanced > Resource limits` を選択します。
そして {guilabel}`Edit instance` をクリックします。

**Memory limit** の {guilabel}`Override` を選んで、absolute の値として 8 GiB を入力します。

![Setting the memory limit for an instance to 8 GiB](/images/UI/limits_memory_example.png)
```
````

```{note}
一部のインスタンスオプションはインスタンスが稼働中に即座に更新されます。
他のインスタンスオプションはインスタンスの再起動後に更新されます。

どのオプションがインスタンス稼働中に即座に反映されるかの情報は {ref}`instance-options` の "ライブアップデート" 列を参照してください。
```

(instances-configure-properties)=
## インスタンスプロパティを設定する

````{tabs}
```{group-tab} CLI
インスタンス作成後にインスタンスプロパティを変更するには、 `--property` フラグを指定して [`incus config set`](incus_config_set.md) コマンドを使います。
インスタンス名とインスタンスプロパティのキーとバリューを指定します:

    incus config set <instance_name> <property_key>=<property_value> <property_key>=<property_value> ... --property

同じフラグを使って、設定オプションを解除するのと全く同じようにインスタンスプロパティも設定解除できます:

    incus config unset <instance_name> <property_key> --property

指定したプロパティの値を取得もできます:

    incus config get <instance_name> <property_key> --property
```

```{group-tab} API
API でインスタンスプロパティを変更するには、インスタンスオプションの変更と同じ仕組みを使います。
唯一の違いはプロパティは設定の root レベルにありますが、オプションは `config` フィールドは以下にあることです。

ですので、インスタンスプロパティを設定するには、インスタンスに PATCH リクエストを送ります:

    incus query --request PATCH /1.0/instances/<instance_name> --data '{"<property_key>":"<property_value>","<property_key>":"property_value>"}}'

インスタンスプロパティを設定解除するには、設定解除したいプロパティを除いた完全なインスタンス設定を含む PUT リクエストをくります。

詳細は [`PATCH /1.0/instances/{name}`](swagger:/instances/instance_patch) と [`PUT /1.0/instances/{name}`](swagger:/instances/instance_put) を参照してください。
```

```{group-tab} UI
Incus UI はインスタンスオプションとインスタンスプロパティを区別しません。
ですので、 {ref}`インスタンスオプションを設定する <instances-configure-options>` と同じようにインスタンスプロパティを設定できます。
```
````

(instances-configure-devices)=
## デバイスを設定する

一般的に、デバイスはコンテナの稼働中に追加または削除できます。
仮想マシンはいくつかのデバイスタイプではホットプラグをサポートしますが、すべてではありません。

利用可能なデバイスタイプとそのオプションについては {ref}`devices` を参照してください。

```{note}
各デバイスのエントリはインスタンスごとにユニークな名前により識別します。

プロファイルに定義されたデバイスは、プロファイルがインスタンスに割り当てられる順番でインスタンスに適用されます。
インスタンス設定内に直接定義されたデバイスは最後に適用されます。
各ステージで、より以前のステージに同じ名前のデバイスがある場合は、デバイスエントリ全体が最後の定義により上書きされます。

デバイス名は最大64文字です。
```

`````{tabs}
````{group-tab} CLI
インスタンスにデバイスを追加して設定するには、 [`incus config device add`](incus_config_device_add.md) コマンドを使います。

インスタンス名、デバイス名、デバイスタイプと ({ref}`デバイスタイプ <devices>` ごとに) 必要に応じてデバイスオプションを指定します:

    incus config device add <instance_name> <device_name> <device_type> <device_option_key>=<device_option_value> <device_option_key>=<device_option_value> ...

例えば、ホストシステムの `/share/c1` 上のストレージをインスタンスのパス `/opt` に追加するには、以下のコマンドを入力します:

    incus config device add my-container disk-storage-device disk source=/share/c1 path=/opt

以前追加したデバイスのインスタンスデバイスオプションを設定するには、 [`incus config device set`](incus_config_device_set.md) コマンドを使います:

    incus config device set <instance_name> <device_name> <device_option_key>=<device_option_value> <device_option_key>=<device_option_value> ...

```{note}
デバイスオプションは {ref}`インスタンスの作成 <instances-create>` 時に `--device` フラグを使って指定することもできます。
これは {ref}`プロファイル <profiles>` を通して提供されるデバイスのデバイスオプションを上書きしたい場合に有用です。
```

デバイスを除去するには、[`incus config device remove`](incus_config_device_remove.md) コマンドを使います。
利用可能なコマンドの完全なリストは [`incus config device --help`](incus_config_device.md) を参照してください。

````

````{group-tab} API
インスタンスにデバイスを追加して設定するには、インスタンス設定を変更するのと同じ仕組みを使います。
デバイス設定は設定の `devices` フィールドの下に配置されています。

インスタンス名、デバイス名、デバイスタイプと ({ref}`デバイスタイプ <devices>` ごとに) 必要に応じてデバイスオプションを指定します:

    incus query --request PATCH /1.0/instances/<instance_name> --data '{"devices": {"<device_name>": {"type":"<device_type>","<device_option_key>":"<device_option_value>","<device_option_key>":"device_option_value>"}}}'

例えば、ホストシステムの `/share/c1` 上のストレージをインスタンスのパス `/opt` に追加するには、以下のコマンドを入力します:

    incus query --request PATCH /1.0/instances/my-container --data '{"devices": {"disk-storage-device": {"type":"disk","source":"/share/c1","path":"/opt"}}}'

詳細は [`PATCH /1.0/instances/{name}`](swagger:/instances/instance_patch) を参照してください。
````

````{group-tab} UI
現時点では UI はデバイスは限定的にサポートしています。

インスタンスにデバイスをアタッチするには、まず作成する必要があります。
次に（{ref}`インスタンスオプションを設定する <instances-configure-options>` と同じ方法で）インスタンス設定を変更し、インスタンスにデバイスをアタッチします。

```{note}
インスタンス設定に表示されるデバイスの一部は {ref}`プロファイル <profiles>` から継承されているか {ref}`プロジェクト <projects>` で定義されたものです。
これらのデバイスはインスタンスごとには編集できません。
```

デバイスの追加と変更は現時点では UI ではサポートされていません。 {ref}`instances-configure-edit` の手順に従ってください。
````

`````

## インスタンス設定を表示する

````{tabs}
```{group-tab} CLI
書き込み可能なインスタンスプロパティ、インスタンスオプション、デバイスとデバイスオプションを含むインスタンスの現在の設定を表示するには、以下のコマンドを入力します:

    incus config show <instance_name> --expanded
```

```{group-tab} API
書き込み可能なインスタンスプロパティ、インスタンスオプション、デバイスとデバイスオプションを含むインスタンスの現在の設定を取得するには、インスタンスに GET リクエストを送ります:

    incus query /1.0/instances/<instance_name>

詳細は [`GET /1.0/instances/{name}`](swagger:/instances/instance_get) を参照してください。
```

```{group-tab} UI
インスタンスの現在の設定を表示するには、 {guilabel}`Instances` に遷移し、インスタンスを選択して、 {guilabel}`Configuration` タブに切り替えます。

インスタンスプロパティ、インスタンスオプション、デバイス、デバイスオプション（または現時点では UI でサポートされていない情報）を含む完全な設定を見るには、 {guilabel}`YAML configuration` を選択します。
この画面はインスタンス設定の完全な YAML を表示します。
```
````

(instances-configure-edit)=
## インスタンス設定全体を編集する

`````{tabs}
````{group-tab} CLI
書き込み可能なインスタンスプロパティ、インスタンスオプション、デバイスとデバイスオプションを含むインスタンス設定全体を編集するには、以下のコマンドを入力します:

    incus config edit <instance_name>

```{note}
利便性のため、 [`incus config edit`](incus_config_edit.md) コマンドは読み取り専用のインスタンスプロパティを含む設定全体を表示します。
しかし、これらのプロパティは変更できません。
変更しても無視されます。
```
````

````{group-tab} API
書き込み可能なインスタンスプロパティ、インスタンスオプション、デバイスとデバイスオプションを含むインスタンス設定全体を編集するには、インスタンスに PUT リクエストを送ります:

    incus query --request PUT /1.0/instances/<instance_name> --data '<instance_configuration>'

詳細は [`PUT /1.0/instances/{name}`](swagger:/instances/instance_put) を参照してください。

```{note}
提供する設定内に読み取り専用のインスタンスプロパティの変更を含めた場合、それらは無視されます。
```
````

````{group-tab} UI
インスタンスを設定する UI フォームを使う代わりに、インスタンスの YAML 設定を編集することもできます。
UI で編集できない設定を変更する場合はこの方法が必要です。

```{important}
更新の際、変更を保存せずに YAML 設定画面から遷移しないでください。
遷移すると、変更は失われます。
```

インスタンスの YAML 設定を編集するには、インスタンス詳細ページに遷移し、 {guilabel}`Configuration` タブに切り替えて {guilabel}`YAML configuration` を選択します。
そして {guilabel}`Edit instance` をクリックします。

YAML 設定を必要に応じて編集します。
そして {guilabel}`Save changes` をクリックして変更した設定を保存します。

```{note}
利便性のため、 YAML は読み取り専用のインスタンスプロパティを含む設定全体を含みます。
しかし、これらのプロパティは変更できません。
変更しても無視されます。
```
````
`````
