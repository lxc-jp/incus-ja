(cluster-member-config)=
# クラスタメンバーの設定

各クラスタメンバーは以下のサポートされる Namespace 内で独自のキー・バリュー設定を持てます:

- `user`（ユーザーのメタデータ用に自由形式のキー・バリュー）
- `scheduler`（メンバーが自クラスタによりどのように動的にターゲットされるかに関連するオプション）

現状サポートされるキーは以下のとおりです:

% Include content from [../config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group cluster-cluster start -->
    :end-before: <!-- config group cluster-cluster end -->
```
