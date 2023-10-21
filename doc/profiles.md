(profiles)=
# プロファイルを使用するには

プロファイルは一組の設定オプションを保持します。
プロファイルにはインスタンスオプション、デバイスとデバイスオプションを含められます。

1 つのインスタンスには任意の数のプロファイルを適用できます。
プロファイルは指定された順番に適用され、その結果最後に指定したプロファイルが特定のキーを上書きします。
どのような場合でも、インスタンス固有の設定はプロファイル由来のものを上書きします。

```{note}
プロファイルはコンテナと仮想マシンに適用できます。
ですので、どちらのタイプに有効なオプションとデバイスを含めることができます。

インスタンスタイプに適用できない設定を含むプロファイルを適用すると、この設定は無視されエラーにはなりません。
```

新しいインスタンスを起動する際にプロファイルを指定しない場合は、自動的には`default`プロファイルが適用されます。
このプロファイルはネットワークインターフェースとルートディスクを定義します。
`default`プロファイルはリネームや削除はできません。

## プロファイルを表示する

すべての利用可能なプロファイルを一覧表示するには以下のコマンドを入力します:

    incus profile list

プロファイルの内容を表示するには以下のコマンドを入力します:

    incus profile show <profile_name>

## 空のプロファイルを作成する

空のプロファイルを作成するには以下のコマンドを入力します:

    incus profile create <profile_name>

(profiles-edit)=
## プロファイルを編集する

プロファイルの特定の設定オプションを設定するか、あるいは YAML 形式でプロファイル全体を編集できます。

### プロファイルの特定の設定オプションを設定する

プロファイルのインスタンスオプションを設定するには、[`incus profile set`](incus_profile_set.md) コマンドを使います。
プロファイル名とインスタンスオプションのキーとバリューを指定します:

    incus profile set <profile_name> <option_key>=<option_value> <option_key>=<option_value> ...

プロファイルのインスタンスデバイスを追加と変更するには、[`incus profile device add`](incus_profile_device_add.md) コマンドを使います。
プロファイル名、デバイス名、デバイスタイプと({ref}`デバイスタイプ <devices>`ごとの)必要に応じてデバイスオプションを指定します:

    incus profile device add <profile_name> <device_name> <device_type> <device_option_key>=<device_option_value> <device_option_key>=<device_option_value> ...

以前にプロファイルに追加したデバイスのインスタンスデバイスオプションを設定するには、[`incus profile device set`](incus_profile_device_set.md) コマンドを使います:

    incus profile device set <profile_name> <device_name> <device_option_key>=<device_option_value> <device_option_key>=<device_option_value> ...

### プロファイル全体を編集する

個々の設定オプションを別々に設定する代わりに、YAML 形式で一度にすべてのオプションを提供できます。

既存のプロファイルまたはインスタンス設定の中身で必要なマークアップをチェックします。
たとえば、`default`プロファイルは以下のようになっているかもしれません:

    config: {}
    description: Default Incus profile
    devices:
      eth0:
        name: eth0
        network: incusbr0
        type: nic
      root:
        path: /
        pool: default
        type: disk
    name: default
    used_by:

インスタンスオプションは`config`の下に提供されます。
インスタンスデバイスとインスタンスデバイスオプションは`devices`の下に提供されます。

ターミナルの標準エディタを使ってプロファイルを編集するには、以下のコマンドを入力します:

    incus profile edit <profile_name>

別の方法として、設定を含む YAML ファイル（たとえば、`profile.yaml`）を作成して、以下のコマンドで設定をプロファイルに書き込めます:

    incus profile edit <profile_name> < profile.yaml

## インスタンスにプロファイルを適用する

インスタンスにプロファイルを適用するには以下のコマンドを入力します:

    incus profile add <instance_name> <profile_name>

```{tip}
プロファイル追加後に [`incus config show <instance_name>`](incus_config_show.md) を実行して設定を確認します。

プロファイルが`profiles`の下に一覧表示されます。
しかし、プロファイルからの設定オプションは（`--expanded`フラグを追加しない限り）`config`の下には表示されません。
この挙動の理由はこれらの設定はプロファイルからは取得されインスタンス設定から取得されるわけではないからです。

これはプロファイルを編集する場合、変更はプロファイルを使用している全てのインスタンスに自動的に適用されることを意味します。
```

インスタンスの起動時に`--profile`フラグを追加してプロファイルを指定することもできます:

    incus launch <image> <instance_name> --profile <profile> --profile <profile> ...

## インスタンスからプロファイルを削除する

インスタンスからプロファイルを削除するには以下のコマンドを入力します:

    incus profile remove <instance_name> <profile_name>
