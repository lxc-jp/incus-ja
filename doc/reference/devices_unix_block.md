(devices-unix-block)=
# タイプ: `unix-block`

```{note}
`unix-block`デバイスタイプはコンテナでサポートされます。
ホットプラグをサポートします。
```

Unix ブロックデバイスは、指定したブロックデバイスをインスタンス内の（`/dev`以下の）デバイスとして出現させます。
そのデバイスから読み取りやデバイスへ書き込みができます。

## デバイスオプション

`unix-block`デバイスには以下のデバイスオプションがあります:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-unix-char-block start -->
    :end-before: <!-- config group devices-unix-char-block end -->
```

(devices-unix-block-hotplugging)=
## ホットプラグ

<!-- Include start Hotplugging -->

ホットプラグは`required=false`を設定しデバイスの`source`オプションを指定した場合に有効になります。

この場合、デバイスはホスト上で出現したときに、コンテナの起動後であっても、自動的にコンテナにパススルーされます。
ホストシステムからデバイスが消えると、コンテナからも消えます。
