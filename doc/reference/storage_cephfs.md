(storage-cephfs)=
# CephFS - `cephfs`

% Include content from [storage_ceph.md](storage_ceph.md)
```{include} storage_ceph.md
    :start-after: <!-- Include start Ceph intro -->
    :end-before: <!-- Include end Ceph intro -->
```

{abbr}`CephFS (Ceph File System)` は堅牢でフル機能の POSIX 互換の分散ファイルシステムを提供する Ceph のファイルシステムコンポーネントです。
内部的には ファイルを Ceph オブジェクトにマップし、ファイルのメタデータ（たとえば、ファイルの所有権、ディレクトリーパス、アクセス権限）を別のデータプールに保管します。

## 用語

% Include content from [storage_ceph.md](storage_ceph.md)
```{include} storage_ceph.md
    :start-after: <!-- Include start Ceph terminology -->
    :end-before: <!-- Include end Ceph terminology -->
```

*CephFS ファイルシステム* は 2 つの OSD ストレージプールから構成され、ひとつは実際のデータ、もうひとつはファイルメタデータに使用されます。

## Incus の `cephfs` ドライバー

```{note}
`cephfs` ドライバはコンテントタイプ `filesystem` のカスタムストレージボリュームにのみ使用できます。

他のストレージボリュームには {ref}`Ceph <storage-ceph>` ドライバを使用してください。
そのドライバはコンテントタイプ `filesystem` のカスタムストレージボリュームにも使用できますが、 Ceph RBD イメージを使って実装しています。

使用したい CephFS ファイルシステムを事前に作成しておいて [`source`](storage-cephfs-pool-config) に指定するか、ファイルシステムと（[`cephfs.data_pool`](storage-cephfs-pool-config) と [`cephfs.meta_pool`](storage-cephfs-pool-config) で指定される名前で）データとメタデータ OSD プールを自動的に作成する[`cephfs.create_missing`](storage-cephfs-pool-config) オプションを指定します。
```

% Include content from [storage_ceph.md](storage_ceph.md)
```{include} storage_ceph.md
    :start-after: <!-- Include start Ceph driver cluster -->
    :end-before: <!-- Include end Ceph driver cluster -->
```

使用したい CephFS ファイルシステムは事前に作成する必要があり [`source`](storage-cephfs-pool-config) オプションで指定する必要があります。

% Include content from [storage_ceph.md](storage_ceph.md)
```{include} storage_ceph.md
    :start-after: <!-- Include start Ceph driver remote -->
    :end-before: <!-- Include end Ceph driver remote -->
```

% Include content from [storage_ceph.md](storage_ceph.md)
```{include} storage_ceph.md
    :start-after: <!-- Include start Ceph driver control -->
    :end-before: <!-- Include end Ceph driver control -->
```

Incus の `cephfs` ドライバーはサーバー側でスナップショットが有効な場合はスナップショットをサポートします。

## 設定オプション

`cephfs` ドライバーを使うストレージプールとこれらのプール内のストレージボリュームには以下の設定オプションが利用できます。

(storage-cephfs-pool-config)=
## ストレージプール設定

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_cephfs-common start -->
    :end-before: <!-- config group storage_cephfs-common end -->
```

{{volume_configuration}}

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group storage_volume_cephfs-common start -->
    :end-before: <!-- config group storage_volume_cephfs-common end -->
```

[^*]: {{snapshot_pattern_detail}}
