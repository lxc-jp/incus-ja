(howto-cluster-groups)=
# クラスタグループをセットアップするには

ラスタメンバーは {ref}`cluster-groups` にアサインできます。
デフォルトでは、すべてのクラスタメンバーは `default` グループに属しています。

クラスタグループを作成するには、[`incus cluster group create`](incus_cluster_group_create.md) コマンドを使用します。
たとえば:

    incus cluster group create gpu

クラスタメンバーを 1 つまたは複数のグループに割り当てるには、[`incus cluster group assign`](incus_cluster_group_assign.md) コマンドを使用します。
このコマンドは、指定したクラスタメンバーを現在所属しているすべてのクラスタグループから削除し、その後、指定したグループまたはグループに追加します。

たとえば、`server1`を`gpu`グループのみに割り当てるには、次のコマンドを使用します:

    incus cluster group assign server1 gpu

`server1`を`gpu`グループに割り当てるとともに、`default`グループにも保持させるためには、以下のコマンドを使用します:

    incus cluster group assign server1 default,gpu

他のグループからメンバーを削除せずに特定のグループにクラスタメンバーを追加するには [`incus cluster group add`](incus_cluster_group_add.md) コマンドを使います。

たとえば、`server1` を `default` グループに残したまま `gpu` に追加するには、以下のコマンドを使います:

    incus cluster group add server1 gpu

## クラスタグループメンバー上でインスタンスを起動する

クラスタグループがある場合、インスタンスを特定のメンバー上で動かすようにターゲットする代わりに、クラスタグループのいずれかのメンバー上で動かすようにターゲットできます。

```{note}
クラスタグループにインスタンスをターゲットできるようにするには {config:option}`cluster-cluster:scheduler.instance` は `all`（デフォルト）または `group` に設定する必要があります。

詳細は{ref}`clustering-instance-placement`を参照してください。
```

クラスタグループのメンバー上でインスタンスを起動するには、{ref}`cluster-target-instance` の指示に従ってください。ただし `--target` フラグではグループ名の前に `@` をつけて指定してください。
たとえば:

    incus launch images:ubuntu/22.04 c1 --target=@gpu
