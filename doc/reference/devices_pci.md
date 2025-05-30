(devices-pci)=
# タイプ: `pci`

```{note}
`pci`デバイスタイプは VM でサポートされます。
ホットプラグはサポートされません。
```

PCI デバイスは生の PCI デバイスをホストから仮想マシンにパススルーするために使用されます。

これらや主にサウンドカードやビデオキャプチャーカードのような特別な単一機能の PCI カードに使われることを意図しています。
理論上は、GPU やネットワークカードなどより高度な PCI デバイスも使用できますが、通常はそれらのデバイスのために Incus が提供する個別のデバイスタイプ（[`gpu`デバイス](devices-gpu)や[`nic` デバイス](devices-nic)）を使うほうがより便利です。

## デバイスオプション

`pci` デバイスには以下のデバイスオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-pci start -->
    :end-before: <!-- config group devices-pci end -->
```
