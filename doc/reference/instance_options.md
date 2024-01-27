(instance-options)=
# インスタンスオプション

インスタンスオプションはインスタンスに直接関係する設定オプションです。

インスタンスオプションをどのように設定するかの手順は{ref}`instances-configure-options`を参照してください。

key/value 形式の設定は、名前空間で分けられています。
以下のオプションが利用できます:

- {ref}`instance-options-misc`
- {ref}`instance-options-boot`
- [`cloud-init` 設定](instance-options-cloud-init)
- {ref}`instance-options-limits`
- {ref}`instance-options-migration`
- {ref}`instance-options-nvidia`
- {ref}`instance-options-raw`
- {ref}`instance-options-security`
- {ref}`instance-options-snapshots`
- {ref}`instance-options-volatile`


各オプションに型が定義されていますが、すべての値は文字列として保管され、REST API で文字列としてエクスポートされる（こうすることで後方互換性を壊すことなく任意の追加の値をサポートできます）ことに注意してください。

(instance-options-misc)=
## その他のオプション

以下のセクションに一覧表示される設定オプションに加えて、以下のインスタンスオプションがサポートされます:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-miscellaneous start -->
    :end-before: <!-- config group instance-miscellaneous end -->
```

```{config:option} environment.* instance-miscellaneous
:type: "string"
:liveupdate: "yes (exec)"
:shortdesc: "インスタンスのための環境変数"

key/value の環境変数をインスタンスにエクスポートできます。
これらはその後 [`incus exec`](incus_exec.md) に設定されます。
```

(instance-options-boot)=
## ブート関連のオプション

以下のインスタンスオプションはインスタンスのブート関連の挙動を制御します:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-boot start -->
    :end-before: <!-- config group instance-boot end -->
```

(instance-options-cloud-init)=
## `cloud-init` 設定

以下のインスタンスオプションはインスタンスの[`cloud-init`](cloud-init)設定を制御します:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-cloud-init start -->
    :end-before: <!-- config group instance-cloud-init end -->
```

これらのオプションのサポートは使用するイメージに依存し、保証はされません。

`cloud-init.user-data`と`cloud-init.vendor-data`の両方を指定すると、両方のオプションの設定がマージされます。
このため、これらのオプションに設定する`cloud-init`設定が同じキーを含まないようにしてください。

(instance-options-limits)=
## リソース制限

以下のインスタンスオプションはインスタンスのリソース制限を指定します:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-resource-limits start -->
    :end-before: <!-- config group instance-resource-limits end -->
```

```{config:option} limits.kernel.* instance-resource-limits
:type: "string"
:liveupdate: "no"
:condition: "container"
:shortdesc: "インスタンスごとのカーネルリソース"

インスタンスにカーネルの制限を設定できます、たとえば、オープンできるファイル数を制限できます。
詳細は {ref}`instance-options-limits-kernel` を参照してください。
```

### PU制限

CPU 使用率を制限するための異なるオプションがあります:

- `limits.cpu`を設定して、インスタンスが見ることができ、使用することができる CPU を制限します。
  このオプションの設定方法は、{ref}`instance-options-limits-cpu`を参照してください。
- `limits.cpu.allowance`を設定して、インスタンスが利用可能な CPU にかける負荷を制限します。
  このオプションはコンテナのみで利用可能です。
  このオプションの設定方法は、{ref}`instance-options-limits-cpu-container`を参照してください。

これらのオプションは同時に設定して、インスタンスが見ることができる CPU とそれらのインスタンスの許可される使用量の両方を制限することが可能です。
しかし、`limits.cpu.allowance`を時間制限と共に使用する場合、スケジューラーに多くの制約をかけ、効率的な割り当てが難しくなる可能性があるため、`limits.cpu`の追加使用は避けるべきです。

CPU 制限は cgroup コントローラーの`cpuset`と`cpu`を組み合わせて実装しています。

(instance-options-limits-cpu)=
#### CPUピンニング

`limits.cpu`は`cpuset`コントローラーを使って、CPU を固定（ピンニング）します。
どの CPU を、またはどれぐらいの数の CPU を、インスタンスから見えるようにし、使えるようにするかを指定できます:

- どの CPU を使うかを指定するには、`limits.cpu`を CPU の組み合わせ（例:`1,2,3`）あるいは CPU の範囲（例:`0-3`）で指定できます。

  単一の CPU にピンニングするためには、CPU の個数との区別をつけるために、範囲を指定する文法（例:`1-1`）を使う必要があります。
- CPU の個数を指定した場合（例:`4`）、Incus は特定の CPU にピンニングされていないすべてのインスタンスをダイナミックに負荷分散し、マシン上の負荷を分散しようとします。
  インスタンスが起動したり停止するたびに、またシステムに CPU が追加されるたびに、インスタンスはリバランスされます。

##### 仮想マシンのCPUリミット

```{note}
Incus は`limits.cpu`オプションのライブアップデートをサポートします。
しかし、仮想マシンの場合は、対応する CPU がホットプラグされるだけです。
ゲストのオペレーティングシステムによって、新しい CPU をオンラインにするためには、インスタンスを再起動するか、なんらかの手動の操作を実行する必要があります。
```

Incus の仮想マシンはデフォルトでは 1 つの vCPU だけを割り当てられ、ホストの CPU のベンダーとタイプとマッチした CPU として現れますが、シングルコアでスレッドなしになります。

`limits.cpu`を単一の整数に設定する場合、Incus は複数の vCPU を割り当ててゲストにはフルなコアとして公開します。
これらの vCPU はホスト上の特定の物理コアにはピンニングされません。
vCPU の個数は VM の稼働中に変更できます。

`limits.cpu`を CPU ID（[`incus info --resources`](incus_info.md) で表示されます）の範囲またはカンマ区切りリストの組に設定する場合、vCPU は物理コアにピンニングされます。
このシナリオでは、Incus は CPU 設定が現実のハードウェアトポロジーとぴったり合うかチェックし、合う場合はそのトポロジーをゲスト内に複製します。
CPU ピンニングを行う場合、VM の稼働中に設定を変更することはできません。

たとえば、ピンニング設定が 8 個のスレッド、同じコアのスレッドの各ペアと 2 個の CPU に散在する偶数のコアを持つ場合、ゲストは 2 個の CPU、各 CPU に 2 個のコア、各コアに 2 個のスレッドを持ちます。
NUMA レイアウトは同様に複製され、このシナリオでは、ゲストではほとんどの場合、2 個の NUMA ノード、各 CPU ソケットに 1 個のノードを持つことになるでしょう。

複数の NUMA ノードを持つような環境では、メモリーは同様に NUMA ノードで分割され、ホスト上で適切にピンニングされ、その後ゲストに公開されます。

これらすべてにより、ゲストスケジューラはソケット、コア、スレッドを適切に判断し、メモリーを共有したり NUMA ノード間でプロセスを移動する際に NUMA トポロジーを考慮できるので、ゲスト内で非常に高パフォーマンスな操作を可能にします。

(instance-options-limits-cpu-container)=
#### 割り当てと優先度（コンテナのみ）

`limits.cpu.allowance`は、時間の制限を与えたときは CFS スケジューラのクォータを、パーセント指定をした場合は全体的な CPU シェアの仕組みを使います:

- 時間制限（たとえば、`20ms/50ms`）はハードリミットです。
  たとえば、コンテナが最大で 1 つの CPU を使用することを許可する場合は、`limits.cpu.allowance`を`100ms/100ms`のような値に設定します。この値は 1 つの CPU に相当する時間に対する相対値なので、2 つの CPU の時間を制限するには、`100ms/50ms`あるいは`200ms/100ms`のような値を使用します。
- パーセント指定を使う場合は、制限は負荷状態にある場合のみに適用されるソフトリミットです。
  設定は、同じ CPU(もしくは CPU の組)を使う他のインスタンスとの比較で、インスタンスに対するスケジューラの優先度を計算するのに使われます。
  たとえば、負荷時のコンテナの CPU 使用率を 1 つの CPU に制限するためには、`limits.cpu.allowance`を`100%`に設定します。


`limits.cpu.nodes`はインスタンスが使用する CPU を特定の NUMA ノードに限定するのに使えます。
どの NUMA ノードを使用するか指定するには、`limits.cpu.nodes`に NUMA ノード ID の組（たとえば、`0,1`）または NUMA ノードの範囲（たとえば、`0-1,2-4`）のどちらかを設定します。

`limits.cpu.priority` は、CPU の組を共有する複数のインスタンスに割り当てられた CPU の割合が同じ場合に、スケジューラの優先度スコアを計算するために使われる別の因子です。

(instance-options-limits-hugepages)=
### huge page の制限

Incus では `limits.hugepage.[size]` キーを使ってコンテナが利用できる huge page の数を制限できます。

アーキテクチャはしばしば huge page のサイズを公開しています。
利用可能な huge page サイズはアーキテクチャによって異なります。

huge page の制限は非特権コンテナ内で`hugetlbfs`ファイルシステムの`mount`システムコールをインターセプトするように Incus を設定しているときには特に有用です。
Incus が`hugetlbfs` `mount`システムコールをインターセプトすると Incus は正しい`uid`と`gid`の値を`mount`オプションに指定して`hugetblfs`ファイルシステムをコンテナにマウントします。
これにより非特権コンテナからも huge page が利用可能となります。
しかし、ホストで利用可能な huge page をコンテナが使い切ってしまうのを防ぐため、`limits.hugepages.[size]`を使ってコンテナが利用可能な huge page の数を制限することを推奨します。

huge page の制限は`hugetlb` cgroup コントローラーによって実行されます。これはこれらの制限を適用するために、ホストシステムが`hugetlb`コントローラーをレガシーあるいは cgroup の単一階層構造(訳注:cgroup v2)に公開する必要があることを意味します。

(instance-options-limits-kernel)=
### カーネルリソース制限

Incus は、インスタンスのリソース制限を設定するのに使用できる一般の名前空間キー`limits.kernel.*`を公開しています。

`limits.kernel.*`接頭辞に続いて指定されるリソースについて Incus が全く検証を行わないという意味でこれは汎用です。
Incus は対象のカーネルがサポートするすべての利用可能なリソースについて知ることはできません。
代わりに、Incus は`limits.kernel.*`接頭辞の後の対応するリソースキーとその値をカーネルに単に渡します。
カーネルが適切な検証を行います。
これによりユーザーはシステム上でサポートされる任意の制限を指定できます。

よくある制限のいくつかは以下のとおりです:

キー                       | リソース            | 説明
:--                       | :---              | :----------
`limits.kernel.as`        | `RLIMIT_AS`       | プロセスの仮想メモリーの最大サイズ
`limits.kernel.core`      | `RLIMIT_CORE`     | プロセスのコアダンプファイルの最大サイズ
`limits.kernel.cpu`       | `RLIMIT_CPU`      | プロセスが使えるCPU時間の秒単位の制限
`limits.kernel.data`      | `RLIMIT_DATA`     | プロセスのデータセグメントの最大サイズ
`limits.kernel.fsize`     | `RLIMIT_FSIZE`    | プロセスが作成できるファイルの最大サイズ
`limits.kernel.locks`     | `RLIMIT_LOCKS`    | プロセスが確立できるファイルロック数の制限
`limits.kernel.memlock`   | `RLIMIT_MEMLOCK`  | プロセスがRAM上でロックできるメモリーのバイト数の制限
`limits.kernel.nice`      | `RLIMIT_NICE`     | 引き上げることができるプロセスのnice値の最大値
`limits.kernel.nofile`    | `RLIMIT_NOFILE`   | プロセスがオープンできるファイルの最大値
`limits.kernel.nproc`     | `RLIMIT_NPROC`    | 呼び出し元プロセスのユーザーが作れるプロセスの最大数
`limits.kernel.rtprio`    | `RLIMIT_RTPRIO`   | プロセスに対して設定できるリアルタイム優先度の最大値
`limits.kernel.sigpending`| `RLIMIT_SIGPENDING` | 呼び出し元プロセスのユーザーがキューに入れられるシグナルの最大数


指定できる制限の完全なリストは `getrlimit(2)`/`setrlimit(2)`システムコールの man ページで確認できます。

`limits.kernel.*`名前空間内で制限を指定するには、`RLIMIT_`を付けずに、リソース名を小文字で指定します。
たとえば、`RLIMIT_NOFILE`は`nofile`と指定します。

制限は、コロン区切りのふたつの数字もしくは`unlimited`という文字列で指定します（たとえば、`limits.kernel.nofile=1000:2000`）。
単一の値を使って、ソフトリミットとハードリミットを同じ値に設定できます（たとえば、`limits.kernel.nofile=3000`）。

明示的に設定されないリソースは、インスタンスを起動したプロセスから継承されます。
この継承は Incus でなく、カーネルによって強制されることに注意してください。

(instance-options-migration)=
## マイグレーションオプション

以下のインスタンスオプションはインスタンスが{ref}`あるLXDサーバーから別のサーバーに移動される <move-instances>`場合の挙動を制御します:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-migration start -->
    :end-before: <!-- config group instance-migration end -->
```

(instance-options-nvidia)=
## NVIDIAとCUDAの設定

以下のインスタンスオプションはインスタンスの NVIDIA と CUDA の設定を指定します:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-nvidia start -->
    :end-before: <!-- config group instance-nvidia end -->
```

(instance-options-raw)=
## rawインスタンス設定のオーバーライド

以下のインスタンスオプションは Incus 自身が使用するバックエンド機能に直接制御できるようにします:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-raw start -->
    :end-before: <!-- config group instance-raw end -->
```

```{important}
これらの`raw.*`キーを設定すると Incus を予期せぬ形で壊してしまうかもしれません。
このため、これらのキーを設定するのは避けるほうが良いです。
```

(instance-options-qemu)=
### QEMU設定のオーバーライド

VM インスタンスに対しては、Incus は`-readconfig`コマンドラインオプションで QEMU に渡す設定ファイルを使って QEMU を設定します。
この設定ファイルは各インスタンスの起動前に生成されます。
設定ファイルは`/run/incus/<instance_name>/qemu.conf`に作られます。

デフォルトの設定はほとんどの典型的な利用ケース、VirtIO デバイスを持つモダンな UEFI ゲスト、では正常に動作します。
しかし、いくつかの状況では、生成された設定をオーバーライドする必要があります。
たとえば以下のような場合です。

- UEFI をサポートしない古いゲスト OS を実行する。
- VirtIO がゲスト OS でサポートされない場合にカスタムな仮想デバイスを指定する。
- マシンの起動前に Incus でサポートされないデバイスを追加する。
- ゲスト OS と衝突するデバイスを削除する。

設定をオーバーライドするには、`raw.qemu.conf`オプションを設定します。
これは`qemu.conf`と似たような形式ですが、いくつか拡張した形式をサポートします。
これは複数行の設定オプションですので、複数のセクションやキーを変更するのに使えます。

- 生成された設定ファイルのセクションやキーを置き換えるには、別の値を持つセクションを追加します。

  たとえば、デフォルトの`virtio-gpu-pci` GPU ドライバーをオーバーライドするには以下のセクションを使います:

  ```
  raw.qemu.conf: |-
      [device "qemu_gpu"]
      driver = "qxl-vga"
  ```

- セクションを削除するには、キー無しのセクションを指定します。
  たとえば:

  ```
  raw.qemu.conf: |-
      [device "qemu_gpu"]
  ```

- キーを削除するには、空の文字列を値として指定します。
  たとえば:

  ```
  raw.qemu.conf: |-
      [device "qemu_gpu"]
      driver = ""
  ```

- 新規のセクションを追加するには、設定ファイル内に存在しないセクション名を指定します。

QEMU で使用される設定ファイル形式は同じ名前の複数のセクションを許可します。
以下は Incus で生成される設定の抜粋です。

```
[global]
driver = "ICH9-LPC"
property = "disable_s3"
value = "1"

[global]
driver = "ICH9-LPC"
property = "disable_s4"
value = "1"
```

オーバーライドするセクションを指定するには、インデクスを指定します。
たとえば:

```
raw.qemu.conf: |-
    [global][1]
    value = "0"
```

セクションのインデクスは 0（指定しない場合のデフォルト値）から始まりますので、上の例は以下の設定を生成します:

```
[global]
driver = "ICH9-LPC"
property = "disable_s3"
value = "1"

[global]
driver = "ICH9-LPC"
property = "disable_s4"
value = "0"
```

(instance-options-security)=
## セキュリティーポリシー

以下のインスタンスオプションはインスタンスの{ref}`security`ポリシーを制御します:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-security start -->
    :end-before: <!-- config group instance-security end -->
```

(instance-options-snapshots)=
## スナップショットのスケジュールと設定

以下のインスタンスオプションは{ref}`インスタンススナップショット <instances-snapshots>`の作成と削除を制御します:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-snapshots start -->
    :end-before: <!-- config group instance-snapshots end -->
```

(instance-options-snapshots-names)=
### スナップショットの自動命名

{{snapshot_pattern_detail}}

(instance-options-volatile)=
## 揮発性の内部データ

以下の揮発性のキーはインスタンスに固有な内部データを保管するため Incus で現在内部的に使用されています:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group instance-volatile start -->
    :end-before: <!-- config group instance-volatile end --
```

```{note}
揮発性のキーはユーザは設定できません。
```
