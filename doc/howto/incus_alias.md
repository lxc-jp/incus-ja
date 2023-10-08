(incus-alias)=
# コマンドエイリアスを追加するには

Incus コマンドラインクライアントでは良く使うコマンドのエイリアスを追加できます。
長いコマンドのショートカットとして、あるいは既存のコマンドに自動的にフラグを追加するために、エイリアスを使用できます。

コマンドエイリアスを管理するには、[`incus alias`](incus_alias.md)コマンドを使用します。

例えば、インスタンスを削除する際に必ず確認を求めるようにするには`incus delete`に常に`incus delete -i`を実行するようにエイリアスを作成します:

    lxc alias add delete "delete -i"

登録されたすべてののエイリアスを表示するには[`incus alias list`](incus_alias_list.md)を実行します。
すべての利用可能なサブコマンドを表示するには[`incus alias --help`](incus_alias.md)を実行してください。
