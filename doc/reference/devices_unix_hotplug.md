(devices-unix-hotplug)=
# タイプ: `unix-hotplug`

```{note}
`unix-hotplug`デバイスタイプはコンテナでサポートされます。
ホットプラグをサポートします。
```

Unix ホットプラグデバイスは、指定した Unix デバイスをインスタンス内の（`/dev`以下の）デバイスとして出現させます。
デバイスがホストシステム上にある場合は、デバイスから読み取りやデバイスへ書き込みができます。

実装はホスト上で稼働する`systemd-udev`に依存します。

## デバイスオプション

`unix-hotplug`デバイスには以下のデバイスオプションがあります:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-unix-hotplug start -->
    :end-before: <!-- config group devices-unix-hotplug end -->
```
