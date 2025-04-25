(devices-gpu)=
# タイプ: `gpu`

GPU デバイスは、指定の GPU デバイスをインスタンス内に出現させます。

```{note}
コンテナでは、`gpu` デバイスは同時に複数の GPU にマッチングさせることができます。
VM では、各デバイスは1つの GPU にしかマッチできません。
```

以下のタイプの GPU が `gputype` デバイスオプションを使って追加できます:

- [`physical`](#gpu-physical)（コンテナと VM）: GPU 全体をインスタンスにパススルーします。 
  `gputype` が指定されない場合これがデフォルトです。
- [`mdev`](#gpu-mdev)（VM のみ）: 仮想 GPU を作成しインスタンスにパススルーします。
- [`mig`](#gpu-mig)（コンテナのみ）: MIG（Multi-Instance GPU）を作成しインスタンスにパススルーします。
- [`sriov`](#gpu-sriov)（VM のみ）: SR-IOV を有効にした GPU の仮想ファンクション（virtual function）をインスタンスに与えます。

利用可能なデバイスオプションは GPU タイプごとに異なり、以下のセクションの表に一覧表示されます。

(gpu-physical)=
## `gputype`: `physical`

```{note}
`physical` GPU タイプはコンテナと VM の両方でサポートされます。
ホットプラグはコンテナのみでサポートし、VM ではサポートしません。
```

`physical` GPU デバイスは GPU 全体をインスタンスにパススルーします。

### デバイスオプション

`physical` タイプのデバイスには以下のデバイスオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-gpu_physical start -->
    :end-before: <!-- config group devices-gpu_physical end -->
```

(gpu-mdev)=
## `gputype`: `mdev`

```{note}
`mdev` GPU タイプは VM でのみサポートされます。
ホットプラグはサポートしていません。
```

`mdev` GPU デバイスは仮想 GPU を作成しインスタンスにパススルーします。
利用可能な`mdev`プロファイルの一覧は [`incus info --resources`](incus_info.md) を実行すると確認できます。

### デバイスオプション

`mdev` タイプのデバイスには以下のデバイスオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-gpu_mdev start -->
    :end-before: <!-- config group devices-gpu_mdev end -->
```

(gpu-mig)=
## `gputype`: `mig`

```{note}
`mig` GPU タイプはコンテナでのみサポートされます。
ホットプラグはサポートしていません。
```

`mig` GPU デバイスは MIG コンピュートインスタンスを作成しインスタンスにパススルーします。
現状これは NVIDIA MIG を事前に作成しておく必要があります。

### デバイスオプション

`mig` タイプのデバイスには以下のデバイスオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-gpu_mig start -->
    :end-before: <!-- config group devices-gpu_mig end -->
```

`mig.uuid`（NVIDIA drivers 470+）か、`mig.ci`と`mig.gi`（古い NVIDIA ドライバー）の両方を設定する必要があります。

(gpu-sriov)=
## `gputype`: `sriov`

```{note}
`sriov` GPU タイプは VM でのみサポートされます。
ホットプラグはサポートしていません。
```

`sriov` GPU デバイスは SR-IOV が有効な GPU の仮想ファンクション（virtual function）をインスタンスにパススルーします。

### デバイスオプション

`sriov`タイプのデバイスには以下のデバイスオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-gpu_sriov start -->
    :end-before: <!-- config group devices-gpu_sriov end -->
```
