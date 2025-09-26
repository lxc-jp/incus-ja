(storage-linstor)=
# LINSTOR - `linstor`

[LINSTOR](https://linbit.com/linstor/)はオープンソースのソフトウェアデファインドストレージのソリューションで、{abbr}`DRBD (Distributed Replicated Block Device)`でレプリケートされたストレージボリュームを管理するのによく使われます。シンプルな運用にフォーカスしながらも高可用性でかつ高性能なボリュームを提供します。

LINSTORはそれ自身では配下のストレージは管理せず、代わりにZFSやLVMなどの他のコンポーネントを使ってブロックデバイスをプロビジョニングします。これらのブロックデバイスは次に[DRBD](https://linbit.com/drbd/)でレプリケーションされて、フォールトトレランスとノードのストレージの能力にかかわらず任意のクラスターノードでボリュームをマウントできるようにします。ボリュームはDRBDカーネルモジュールを使ってレプリケートされるので、レプリケーションのデータパスは完全にカーネル空間内にとどまり、ユーザー空間で実装されるソリューションと比べてオーバーヘッドを低減できます。

## 用語

LINSTORクラスターは2つのメインコンポーネントから構成されます:*コントローラー*と*サテライト*です。LINSTORコントローラーはデータベースを管理しクラスターの状態と設定を追跡し続けます。一方、サテライトはストレージを提供しクラスター間でボリュームをマウントできるようにします。クライアントはコントローラーとのみ対話し、コントローラーはユーザーのリクエストを満たすために複数のサテライトにまたがる操作を調整する責任を持ちます。

LINSTORは内部の概念にいくぶんオブジェクト指向なアプローチをとっています。これは概念が階層的な性質を持ち、下位レベルのオブジェクトが上位レベルのオブジェクトからプロパティを継承することにも見て取れます。

LINSTORは*ストレージプール*という概念を持ちます。これはLINSTORで使用できる物理ストレージを使ってボリュームを作成します。ストレージプールは（LVMやZFSなどの）バックエンドドライバー、ストレージプールをもつクラスターノードとストレージプール自身あるいはバックエンドストレージに適用できるプロパティを定義します。

LINSTORでは*リソース*はインスタンスで使用できるストレージユニットの表現です。リソースはほとんどの場合DRBDでレプリケートされたブロックデバイスです。その場合そのデバイスの1つのレプリカを表しています。リソースは*リソース定義*にグルーピングでき、これはこれのすべての子のリソースで継承される共通のプロパティを定義します。同様に*リソースグループ*はそれらの子の定義に適用される共通のプロパティを定義します。またリソースグループは指定のリソース定義にいくつのレプリカを作成し、度のストレージプールを使い、レプリカを異なるアベイラビリティゾーンに分散させるか、などの配置ルールも定義します。LINSTORとやりとりする通常の方法は希望のプロパティを持つリソースグループを定義し、次にそこからリソースを*スポーニング*させます。

## Incusの`linstor`ドライバー

```{note}
LINSTORはサテライトノード間でボリュームを移動することとマウントすることしかできません。このため、Incusのすべてのクラスタメンバーがボリュームに確実にアクセスできるようにするため、すべてのIncusノードはLINSTORサテライトノードでもあるようにしなければなりません。言い換えると、`incus`サービスが稼働する各ノードでは`linstor-satellite`サービスも動かすべきです。

しかし、これはIncusノードが必ずストレージも提供しなければならないというわけではないことに注意してください。Incusノードに「ディスクレス」なサテライトをデプロイすることでストレージ用のノードと計算用のノードを分離してLINSTORを使うこともできます。ディスクレスなノードはストレージは提供しませんが、DRBDデバイスをマウントしネットワーク越しにIOを実行できます。
```

他のストレージドライバーとは異なり、このドライバーはストレージシステムをセットアップせず、あなたがLINSTORクラスターを既にインストール済みであることを前提としています。ドライバーは{config:option}`server-miscellaneous:storage.linstor.controller_connection`オプションがIncusで使われるLINSTORコントローラーのエンドポイントに設定されることを要求します。

このトライバーはリモートとローカル両方のストレージを提供できるという点でも他のドライバーと挙動が異なります。ボリュームのディスクありのレプリカがノード上で利用可能な場合、遅延を減らすため読み書きはローカルで実行されます（ですが書き込みはレプリカに同期的にレプリケーションされる必要があるため、ネットワーク遅延の影響はあります）。同時に、ディスクレスレプリカはすべてのIOをネットワーク越しに行い、物理ストレージのあるなしにかかわらず任意のノード上でボリュームをマウントし使えるようにします。これらのハイブリッドな能力によりLINSTORは必要な際はボリュームをクラスターノード間で移動する柔軟性を持ちながらも低遅延なストレージを提供できます。

Incusの`linstor`ドライバーはリソースグループを使ってリソースを管理とスポーンします。次の表はIncusとLINSTORの概念のマッピングを表しています:

| Incusの概念      | LINSTORの概念    |
| :---             | :---             |
| ストレージプール | リソースグループ |
| ボリューム       | リソース定義     |
| スナップショット | スナップショット |

IncusはLINSTORリソースグループの完全な制御を持っていることを前提とします。
このため、IncusのLINSTORリソースグループ内にIncusで所有されてないエンティティは決して作ってはいけません。作るとIncusが削除してしまうかもしれないからです。

リソースを管理する際、IncusはLINSTORサテライトノードがどのIncusノードに対応するかを決定できる必要があります。デフォルトではIncusはノード名がLINSTORのノード名に一致する（例えば`incus cluster list`と`linstor node list`が同じノード名を表示する）ことを想定しています。Incusがスタンドアロンのサーバーとして動いている（クラスターではない）場合、ホスト名がノード名として使われます。IncusとLINSTORでノード名が一致しない場合、各Incusノードで{config:option}`server-miscellaneous:storage.linstor.satellite.name`を適切なLINSTORサテライトノード名に設定できます。

### 制限

`linstor`ドライバーは以下の制限があります:

インスタンス間でのカスタムボリュームの共有
: {ref}`content type <storage-content-types>` `filesystem`があるカスタムストレージボリュームは通常異なるクラスタメンバー上の複数のインスタンスで共有できます。
  しかし、LINSTORドライバーはcontent type `filesystem`を持つボリュームをDRBDでレプリケートされたデバイス上にファイルシステムを持つことで「シミュレート」しているため、カスタムボリュームは一度に1つのインスタンスにしか割り当てられません。

Incusインストール環境間でのリソースグループの共有
: 複数のIncusインストール環境間で同じLINSTORリソースグループを共有するのはサポートされていません。

より古いスナップショットの復元
: LINSTORは最新のスナップショット以外の復元はサポートしていません。
  しかし、古いスナップショットから新しいインスタンスを作ることはできます。
  この方法により特定のスナップショットがあなたの必要とするものを含んでいるかを確認することができます。
  正しいスナップショットを特定したら、{ref}`新しいスナップショットを削除 <storage-edit-snapshots>`してあなたが欲しいスナップショットが最新のスナップショットにしてから復元できます。

  別の方法として、復元の際にIncusがより新しいスナップショットを自動で破棄するように設定することもできます。
  そうするには、[`linstor.remove_snapshots`](storage-linstor-vol-config)設定オプションをボリューム（あるいは対応する`volume.linstor.remove_snapshots`設定をプール内のすべてのボリュームのストレージプールに）設定します。

## 設定オプション

以下の設定オプションが`linstor`ドライバーを使うストレージプールとそれらのプール内のストレージボリュームで利用できます。

(storage-linstor-pool-config)=
### ストレージプール設定

| キー                                  | 型     | デフォルト値    | 設定                                                                                                                                                                                                                                                 |
| :---                                  | :---   | :---            | :---                                                                                                                                                                                                                                                 |
| `linstor.resource_group.name`         | string | `incus`         | ストレージプールで使用されるLINSTORリソースグループ名                                                                                                                                                                                                |
| `linstor.resource_group.place_count`  | int    | 2               | リソースグループ内のリソースのために作成されるべきディスクフルレプリカの数。すでにボリュームがあるプールでこのオプションの値を増やすと、LINSTORはすべての既存のリソースが新しい値に合致するように新しくディスクフルレプリカを作成することになります |
| `linstor.resource_group.storage_pool` | string | -               | サテライトノード上にリソースが配置されるストレージプール名                                                                                                                                                                                           |
| `linstor.volume.prefix`               | string | `incus-volume-` | LINSTORが管理するボリュームの内部名に使われる接頭辞。ストレージプール作成後は変更不可                                                                                                                                                                |
| `drbd.on_no_quorum`                   | string | -               | クオラムが失われた際に使用されるDRBDポリシー（リソースグループに適用される）                                                                                                                                                                         |
| `drbd.auto_diskful`                   | string | -               | ノード上のストレージが利用可能な場合にプライマリのディスクレスリソースがディスクフルに変換されるまでの期間を表す文字列（リソースグループに適用される）                                                                                               |
| `drbd.auto_add_quorum_tiebreaker`     | bool   | `true`          | LINSTORが必要に応じて自動的にディスクレスリソースを作ってクオラムのタイブレーカーとして振る舞わせることを許可するかどうか（リソースグループに適用される）                                                                                            |

{{volume_configuration}}

(storage-linstor-vol-config)=
### ストレージボリューム設定

| キー                              | 型     | 条件                                                    | デフォルト値                                     | 説明                                                                                                                                                      |
| :---                              | :---   | :---                                                    | :---                                             | :---                                                                                                                                                      |
| `block.filesystem`                | string | content type `filesystem`をもつブロックベースボリューム | `volume.block.filesystem`と同じ                  | {{block_filesystem}}                                                                                                                                      |
| `block.mount_options`             | string | content type `filesystem`をもつブロックベースボリューム | `volume.block.mount_options`と同じ               | ブロックベースのファイルシステムボリュームのマウントオプション                                                                                            |
| `initial.gid`                     | int    | content type `filesystem`をもつカスタムボリューム       | `volume.initial.uid`と同じか`0`                  | インスタンス内のボリューム所有者のGID                                                                                                                     |
| `initial.mode`                    | int    | content type `filesystem`をもつカスタムボリューム       | `volume.initial.mode`と同じか`711`               | インスタンス内のボリュームのモード                                                                                                                        |
| `initial.uid`                     | int    | content type `filesystem`をもつカスタムボリューム       | `volume.initial.gid`と同じか`0`                  | インスタンス内のボリューム所有者のUID                                                                                                                     |
| `security.shared`                 | bool   | カスタムブロックボリューム                              | `volume.security.shared`と同じか`false`          | 複数のインスタンス間でのボリュームの共有を有効にするか                                                                                                    |
| `security.shifted`                | bool   | カスタムボリューム                                      | `volume.security.shifted`と同じか`false`         | {{enable_ID_shifting}}                                                                                                                                    |
| `security.unmapped`               | bool   | カスタムボリューム                                      | `volume.security.unmapped`と同じか`false`        | ボリュームにIPマッピングを無効化するか                                                                                                                    |
| `size`                            | string |                                                         | `volume.size`と同じ                              | ストレージボリュームのサイズ／クォータ                                                                                                                    |
| `snapshots.expiry`                | string | カスタムボリューム                                      | `volume.snapshots.expiry`                        | {{snapshot_expiry_format}}                                                                                                                                |
| `snapshots.expiry.manual`         | string | カスタムボリューム                                      | `volume.snapshots.expiry.manual` と同じ          | {{snapshot_expiry_format}}                                                                                                                                |
| `snapshots.pattern`               | string | カスタムボリューム                                      | `volume.snapshots.pattern`と同じか`snap%d`       | {{snapshot_pattern_format}} [^*]                                                                                                                          |
| `snapshots.schedule`              | string | カスタムボリューム                                      | `volume.snapshots.schedule`と同じ                | {{snapshot_schedule_format}}                                                                                                                              |
| `drbd.on_no_quorum`               | string |                                                         | -                                                | クオラムが失われた際に使用されるDRBDポリシー（リソースグループに適用される）                                                                              |
| `drbd.auto_diskful`               | string |                                                         | -                                                | ノード上のストレージが利用可能な場合にプライマリのディスクレスリソースがディスクフルに変換されるまでの期間を表す文字列（リソースグループに適用される）    |
| `drbd.auto_add_quorum_tiebreaker` | bool   |                                                         | `true`                                           | LINSTORが必要に応じて自動的にディスクレスリソースを作ってクオラムのタイブレーカーとして振る舞わせることを許可するかどうか（リソースグループに適用される） |
| `linstor.remove_snapshots`        | bool   |                                                         | `volume.linstor.remove_snapshots`と同じか`false` | 必要に応じてスナップショットを削除するか                                                                                                                  |

[^*]: {{snapshot_pattern_detail}}

```{toctree}
:maxdepth: 1
:hidden:

LINSTORのセットアップ </howto/storage_linstor_setup>
ドライバーの内部構造 <storage_linstor_internals>
```
