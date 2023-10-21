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

キー        | 型     | デフォルト値 | 説明
:--         | :--    | :--          | :--
`gid`       | int    | `0`          | インスタンス（コンテナのみ）内のデバイス所有者のGID
`id`        | string | -            | GPUデバイスのDRMカードID
`mode`      | int    | `0660`       | インスタンス（コンテナのみ）内のデバイスのモード
`pci`       | string | -            | GPUデバイスのPCIアドレス
`productid` | string | -            | GPUデバイスのプロダクトID
`uid`       | int    | `0`          | インスタンス（コンテナのみ）内のデバイス所有者のUID
`vendorid`  | string | -            | GPUデバイスのベンダーID

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

キー        | 型     | デフォルト値 | 説明
:--         | :--    | :--          | :--
`id`        | string | -            | GPUデバイスのDRMカードID
`mdev`      | string | -            | 使用する`mdev`プロファイル（必須 - 例:`i915-GVTg_V5_4`）
`pci`       | string | -            | GPUデバイスのPCIアドレス
`productid` | string | -            | GPUデバイスのプロダクトID
`vendorid`  | string | -            | GPUデバイスのベンダーID

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

キー        | 型     | デフォルト値 | 説明
:--         | :--    | :--          | :--
`id`        | string | -            | GPUデバイスのDRMカードID
`mig.ci`    | int    | -            | 既存のMIGコンピュートインスタンスID
`mig.gi`    | int    | -            | 既存のMIG GPUインスタンスID
`mig.uuid`  | string | -            | 既存のMIGデバイスUUID（`MIG-`接頭辞は省略可）
`pci`       | string | -            | GPUデバイスのPCIアドレス
`productid` | string | -            | GPUデバイスのプロダクトID
`vendorid`  | string | -            | GPUデバイスのベンダーID

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

キー        | 型     | デフォルト値 | 説明
:--         | :--    | :--          | :--
`id`        | string | -            | GPUデバイスのDRMカードID
`pci`       | string | -            | 親GPUデバイスのPCIアドレス
`productid` | string | -            | 親GPUデバイスのプロダクトID
`vendorid`  | string | -            | 親GPUデバイスのベンダーID
