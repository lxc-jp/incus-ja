(devices-infiniband)=
# タイプ: `infiniband`

```{note}
`infiniband`デバイスタイプはコンテナと VM の両方でサポートされます。
ホットプラグはコンテナのみでサポートし、VM ではサポートしません。
```

Incus では、InfiniBand デバイスに対する 2 種類の異なったネットワークタイプが使えます:

- `physical`: ホストの物理デバイスをインスタンスにパススルーします。
  対象のデバイスはホスト上では見えなくなり、インスタンス内に出現します。
- `sriov`: SR-IOV が有効な物理ネットワークデバイスの仮想ファンクション（virtual function）をインスタンスにパススルーします。

  ```{note}
  InfiniBandデバイスはSR-IOVをサポートしますが、他のSR-IOVが有効なデバイスと異なり、InfiniBandはSR-IOVモードの動的なデバイスの作成をサポートしません。
  このため、対応するカーネルモジュールを設定することで仮想ファンクションの数を事前に設定する必要があります。
  ```

`physical`な`infiniband`デバイスを作成するには、以下のコマンドを使用します:

    incus config device add <instance_name> <device_name> infiniband nictype=physical parent=<device>

`sriov`の`infiniband`デバイスを作成するには、以下のコマンドを使用します:

    incus config device add <instance_name> <device_name> infiniband nictype=sriov parent=<sriov_enabled_device>

## デバイスオプション

`infiniband`デバイスには以下のデバイスオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-infiniband start -->
    :end-before: <!-- config group devices-infiniband end -->
```
