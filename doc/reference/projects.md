(ref-projects)=
# プロジェクトの設定

プロジェクトは、キー/値の設定オプションのセットを通じて設定することができます。
これらのオプションを設定する方法については、{ref}`projects-configure` を参照してください。

キー/値の設定は名前空間化されています。
次のオプションが利用可能です:

- {ref}`project-features`
- {ref}`project-limits`
- {ref}`project-restrictions`
- {ref}`project-specific-config`

(project-features)=
## プロジェクトの機能

プロジェクトの機能は、プロジェクト内でどのエンティティが隔離され、どのエンティティが`default`プロジェクトから継承されるかを定義します。

`feature.*` オプションが `true` に設定されている場合、対応するエンティティはプロジェクト内で隔離されます。

```{note}
特定のオプションを明示的に設定せずにプロジェクトを作成すると、このオプションは以下の表で与えられた初期値に設定されます。

ただし、`feature.*` オプションのいずれかを解除すると、初期値に戻るのではなく、デフォルト値に戻ります。
すべての `feature.*` オプションのデフォルト値は `false` です。
```

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group project-features start -->
    :end-before: <!-- config group project-features end -->
```

(project-limits)=
## プロジェクトの制限

プロジェクトの制限は、プロジェクトに属するコンテナや VM が使用できるリソースの上限を定義します。

`limits.*` オプションによっては、プロジェクト内で許可されるエンティティの数に制限が適用されることがあります（たとえば {config:option}`project-limits:limits.containers` や {config:option}`project-limits:limits.networks`）。また、プロジェクト内のすべてのインスタンスのリソース使用量の合計値に制限が適用されることもあります（たとえば、 {config:option}`project-limits:limits.cpu` や {config:option}`project-limits:limits.processes`）。
後者の場合、制限は通常、各インスタンスに設定されている {ref}`instance-options-limits` に適用されます（直接またはプロファイル経由で設定されている場合）、実際に使用されているリソースではありません。

たとえば、プロジェクトの {config:option}`project-limits:limits.memory` 設定を `50GB` に設定した場合、プロジェクトのインスタンスで定義されたすべての {config:option}`project-limits:limits.memory` 設定キーの個別の値の合計が 50GB 未満に保たれます。
{config:option}`project-limits:limits.memory` 設定の合計が 50GB を超えるインスタンスを作成しようとすると、エラーが発生します。

同様に、プロジェクトの {config:option}`project-limits:limits.cpu` 設定キーを `100` に設定すると、個々の {config:option}`project-limits:limits.cpu` 値の合計が 100 未満に保たれます。

プロジェクトの制限を使用する場合、以下の条件を満たす必要があります:

- `limits.*` 設定のいずれかを設定し、インスタンスに対応する設定がある場合、プロジェクト内のすべてのインスタンスに対応する設定が定義されている必要があります（直接またはプロファイル経由で設定）。
  インスタンスの設定オプションについては {ref}`instance-options-limits` を参照してください。
- {ref}`instance-options-limits-cpu` が有効になっている場合、{ref}`instance-options-limits-cpu` 設定は使用できません。
  これは、プロジェクトで {ref}`instance-options-limits-cpu` を使用するためには、プロジェクト内の各インスタンスの {ref}`instance-options-limits-cpu` 設定を CPU の数、または CPU のセットや範囲ではなく、数値に設定する必要があることを意味します。
- {config:option}`project-limits:limits.memory` 設定は、パーセンテージではなく絶対値で設定する必要があります。

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group project-limits start -->
    :end-before: <!-- config group project-limits end -->
```

(project-restrictions)=
## プロジェクトの制約

プロジェクトのインスタンスがセキュリティーに関連する機能（コンテナのネストや raw LXC 設定など）にアクセスできないようにするには、{config:option}`project-restricted:restricted` 設定オプションを `true` に設定します。
その後、さまざまな `restricted.*` オプションを使用して、通常は {config:option}`project-restricted:restricted` によってブロックされる個々の機能を選択し、プロジェクトのインスタンスで使用できるように許可できます。

たとえば、プロジェクトを制限し、すべてのセキュリティー関連機能をブロックしつつ、コンテナのネストを許可するには、次のコマンドを入力します:

    incus project set <project_name> restricted=true
    incus project set <project_name> restricted.containers.nesting=allow

セキュリティーに関連する各機能には、関連する `restricted.*` プロジェクト設定オプションがあります。
機能の使用を許可する場合は、その `restricted.*` オプションの値を変更してください。
ほとんどの `restricted.*` 設定は、`block`（デフォルト）または `allow` に設定できる二値スイッチです。
ただし、一部のオプションは、より細かい制御のために他の値をサポートしています。

```{note}
`restricted.*` オプションを有効にするには、`restricted` 設定を `true` に設定する必要があります。
`restricted` が `false` に設定されている場合、`restricted.*` オプションを変更しても効果はありません。

すべての `restricted.*` キーを `allow` に設定することは、`restricted` 自体を `false` に設定することと同等です。
```

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group project-restricted start -->
    :end-before: <!-- config group project-restricted end -->
```

(project-specific-config)=
## プロジェクト固有の設定

プロジェクトに対していくつかの {ref}`server` オプションを上書きできます。
また、プロジェクトにユーザーメタデータを追加することができます。

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group project-specific start -->
    :end-before: <!-- config group project-specific end -->
```
