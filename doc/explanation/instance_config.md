(instance-config)=
# インスタンスの設定

インスタンス設定は以下の異なるカテゴリから構成されます:

インスタンスプロパティ
: インスタンスプロパティはインスタンスが作成されるときに指定されます。
  これには、たとえば、インスタンス名やアーキテクチャが含まれます。
  いくつかのプロパティは読み取り専用で作成後は変更できませんが、他のプロパティは{ref}`プロパティの値を設定する <instances-configure-properties>` または {ref}`インスタンス設定全体を編集する <instances-configure-edit>` で更新できます。

  YAML 設定内では、プロパティはトップレベルにあります。

  利用可能なインスタンスプロパティのレファレンスは {ref}`instance-properties` を参照してください。

インスタンスオプション
: インスタンスオプションはインスタンスに直接関連する設定オプションです。
  これには、たとえば、起動時のオプション、セキュリティー設定、ハードウェアのリミット、カーネルモジュール、スナップショット、そしてユーザーの鍵を含みます。
  これらのオプションはインスタンスの作成時に (`--config key=value` フラグを使って) キー/バリューペアで指定できます。
  作成後は [`incus config set`](incus_config_set.md) や [`incus config unset`](incus_config_unset.md) コマンドで変更できます。

  YAML 設定内では、オプションは `config` エントリの下に配置されます。

  利用可能なインスタンスオプションのレファレンスは {ref}`instance-options`、オプションをどのように設定するかの手順は {ref}`instances-configure-options` を参照してください。

インスタンスデバイス
: インスタンスデバイスはインスタンスにアタッチされます。
  これらは、たとえば、ネットワークインターフェース、マウントポイント、USB そして GPU デバイスが含まれます。
  通常、デバイスはインスタンスを作成した後に [`incus config device add`](incus_config_device_add.md) コマンドで追加しますが、プロファイルやインスタンスを作成するのに使用する YAML 設定ファイルに追加することもできます。

  各デバイスタイプには固有のオプションのセットがあり、*インスタンスデバイスオプション*として参照されます。

  YAML 設定内では、デバイスは `devices` エントリの下に配置されます。

  利用可能なデバイスと対応するインスタンスデバイスオプションのレファレンスについては {ref}`devices`、インスタンスデバイスをどのように追加し設定するかの手順は {ref}`instances-configure-devices` を参照してください。

```{toctree}
:maxdepth: 1
:hidden:

../reference/instance_properties.md
../reference/instance_options.md
../reference/devices.md
../reference/instance_units.md
```
