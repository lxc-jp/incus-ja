(devices-usb)=
# タイプ: `usb`

```{note}
`usb`デバイスタイプはコンテナと VM の両方でサポートされます。
コンテナと VM の両方でホットプラグをサポートします。
```

USB デバイスは、指定された USB デバイスをインスタンスに出現させます。
パフォーマンスの問題のため、高スループットまたは低レイテンシを要求するデバイスの使用は避けてください。

コンテナでは、（`/dev/bus/usb`にある）`libusb`デバイスのみがインスタンスに渡されます。
この方法はユーザースペースのドライバーを持つデバイスで機能します。
専用のカーネルドライバーを必要とするデバイスは、代わりに[`unix-char`デバイス](devices-unix-char)か[`unix-hotplug`デバイス](devices-unix-hotplug)を使用してください。

仮想マシンでは、USB デバイス全体がパススルーされますので、あらゆる USB デバイスがサポートされます。
デバイスがインスタンスに渡されると、ホストからは消失します。

## デバイスオプション

`usb`デバイスには以下のデバイスオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-usb start -->
    :end-before: <!-- config group devices-usb end -->
```
