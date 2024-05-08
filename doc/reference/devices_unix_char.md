(devices-unix-char)=
# タイプ: `unix-char`

```{note}
`unix-char`デバイスタイプはコンテナでサポートされます。
ホットプラグをサポートします。
```

Unix キャラクタデバイスは、指定したキャラクタデバイスをインスタンス内の（`/dev`以下の）デバイスとして出現させます。
そのデバイスから読み取りやデバイスへ書き込みができます。

## デバイスオプション

`unix-char`デバイスには以下のデバイスオプションがあります:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-unix-char-block start -->
    :end-before: <!-- config group devices-unix-char-block end -->
```

(devices-unix-char-hotplugging)=
## ホットプラグ

% Include content from [devices_unix_block.md](device_unix_block.md)
```{include} devices_unix_block.md
    :start-after: Hotplugging -->
```
