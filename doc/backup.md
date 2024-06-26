(backups)=
# Incusサーバーをバックアップする

本番環境では、常に Incus サーバーのデータをバックアップすべきです。

Incus サーバーはさまざまなエンティティを含んでおり、バックアップ戦略を選択する際には、これらのエンティティのうちどれをバックアップの対象にするかとどれぐらいの頻度で保存するかを決定する必要があります。

## 何をバックアップするか

Incus サーバーのさまざまなコンテンツはファイルシステム上に配置されるものと、それに加えて、{ref}`Incus database <database>`に記録されるものがあります。
ですので、データベースをバックアップするだけあるいはディスク上のファイルをバックアップするだけでは完全に機能するバックアップにはなりません。

Incus サーバーは以下のエンティティを含んでいます:

- インスタンス (データベースのレコードとファイルシステム)
- イメージ (データベースのレコード、イメージファイル、そしてファイルシステム)
- ネットワーク (データベースのレコードと状態ファイル)
- プロファイル (データベースのレコード)
- ストレージボリューム (データベースのレコードとファイルシステム)

これらのうちどれをバックアップする必要があるかを検討してください。
たとえば、カスタムイメージを使用していなければ、イメージはイメージサーバーに存在するのでバックアップは不要です。
`default`プロファイルしか使用していなかったり、標準の`incusbr0`ネットワークブリッジしか使用していない場合、それらは簡単に再生成できますので、バックアップする必要はないかもしれません。

## フルバックアップ

Incus サーバーのすべてのコンテンツをフルバックアップするには、`/var/lib/incus`ディレクトリーをバックアップしてください。

このディレクトリーはローカルストーレジ、Incus データベース、あなたの設定を含みます。
ただし、分離されたストレージデバイスは含みません。
つまりディレクトリーがあなたのインスタンスのデータも含むかはお使いのストレージドライバーによります。

```{important}
Incusサーバが外部ストレージ（たとえば、LVMボリュームグループ、ZFS zpool、あるいは何か他のIncus自身に直接含まれないような外部リソース）を使っている場合、それらは別途バックアップが必要です。
```

データをバックアップするには、`/var/lib/incus`の tarball を作成してください。
あなたのシステムが`/etc/subuid`と`/etc/subgid`ファイルをお使いの場合、これらのファイルもバックアップしてください。
これらをリストアするとインスタンスのファイルシステムで不要なシフトを防げます。

データをリストアするには、以下の手順を実行してください:

1. サーバー上の Incus を停止します（たとえば、`sudo systemctl stop incus.service incus.socket`で）。
1. ディレクトリー（`/var/lib/incus/`）を削除します。
1. バックアップからディレクトリーをリストアします。
1. 外部のストレージデバイスを削除しリストアします。
1. `/etc/subuid`と`/etc/subgid`ファイルがある場合はリストアします。
1. Incus を再起動します（たとえば、`sudo systemctl start incus.socket incus.service`またはマシンを再起動して）。

## 部分的なバックアップ

特定のエンティティをバックアップするだけに決めた場合、実行にはいくつかの異なる選択肢があります。
フルバックアップをしている場合であっても、追加でこれらの部分的なバックアップを検討するのが良いです。
たとえば、完全な Incus サーバーをリストアするよりも単一のインスタンスをリストアしたりプロファイルを再設定するほうが簡単で安全です。

### インスタンスとボリュームのバックアップ

インスタンスとストレージボリュームは非常に似た方法でバックアップされます（というのはインスタンスをバックアップする際は、基本的にはそのインスタンスボリュームをバックアップするからです。{ref}`storage-volume-types`参照）。

詳細な情報は{ref}`instances-backup`と{ref}`howto-storage-backup-volume`を参照してください。
以下のセクションでインスタンスとボリュームをバックアップする際の選択肢の簡単な要約を示します。

#### Incusサーバーのセカンダリバックアップ

Incus は 2 つのホスト間でインスタンスとストレージボリュームのコピーと移動を
サポートしています。
手順は{ref}`move-instances`と{ref}`howto-storage-move-volume`を参照してください。

ですので予備のサーバーがあれば、インスタンスとストレージボリュームをバックアップとして定期的にそのセカンダリサーバーにコピーできます。
必要な場合、セカンダリサーバーに切り替えたり、インスタンスやストレージボリュームをセカンダリサーバーからコピーできます。

セカンダリサーバーを純粋にストレージサーバーとして使う場合、メインの Incus サーバーほど強力である必要はありません。

#### tarballのエクスポート

`export`コマンドを使ってインスタンスとボリュームをバックアップの tarball にエクスポートできます。
デフォルトでは、これらの tarball はすべてのスナップショットを含みます。

最適化された export オプションを使用でき、すると通常はより短時間でエクスポートでき tarball のサイズも小さくなります。
しかし、バックアップの tarball をリストアする際は同じストレージドライバーを使う必要があります。

手順は{ref}`instances-backup-export`と{ref}`storage-backup-export`を参照してください。

#### スナップショット

スナップショットはインスタンスやボリュームの特定の日時での状態を保存します。
しかし、それらは同じストレージプール内に保管されますので、オリジナルのデータが削除されたり失われたりした場合はスナップショットも失われる可能性が高いです。
つまりスナップショットは非常に高速で手軽に作成とリストアができますが、安全なバックアップを構成するものではありません。

詳細な情報は{ref}`instances-snapshots`と{ref}`storage-backup-snapshots`を参照してください。

(backup-database)=
### データベースのバックアップ

{ref}`Incus database <database>`の内容をリストアする自明な方法はありませんが、それでもその内容のバックアップをとっておくと非常に便利です。
たとえば、ネットワークやプロファイルを再生成する必要がでたときに、バックアップがあれば非常に容易になります。

ローカルデータベースの内容をダンプするには以下のコマンドを使用します:

    incus admin sql local .dump > <output_file>

グローバルデータベースの内容をダンプするには以下のコマンドを使用します:

    incus admin sql global .dump > <output_file>

定期的な Incus のバックアップにこれら 2 つのコマンドを含めておくと良いでしょう。
