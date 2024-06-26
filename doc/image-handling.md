(about-images)=
# イメージについて

Incus はイメージをベースとしたワークフローを使用します。
各インスタンスはイメージをベースとしています。イメージは基礎となるオペレーティングシステム（たとえば、Linux ディストリビューション）と Incus に関連するいくつかの情報を含みます。

イメージはリモートのイメージストア（概要は{ref}`image-servers`参照）から利用可能ですが、既存のインスタンスや rootfs イメージをベースにして、独自のイメージを作成できます。

リモートサーバーからローカルのイメージストアにイメージをコピーしたり、ローカルのイメージをリモートサーバーにコピーできます。
ローカルのイメージをリモートのインスタンスを作るのに使うこともできます。

各イメージはフィンガープリント（SHA256）で識別されます。
イメージを管理しやすくするために、Incus では各イメージに 1 つ以上のエイリアスを定義できます。

## キャッシュ

リモートのイメージからインスタンスを作成する際、Incus はイメージをダウンロードしローカルにキャッシュします。
イメージはローカルのイメージストアに cached フラグをセットして保管されます。
イメージは以下のいずれかが発生するまでは非公開のイメージとしてローカルに保持されます:

- {config:option}`server-images:images.remote_cache_expiry` で指定された日数の間新しいインスタンスを作成するのにイメージが使われなかった。
- イメージの有効期限（イメージのプロパティの 1 つ。どのように変更するかの情報は{ref}`images-manage-edit`参照）に達した。

Incus はイメージから新しいインスタンスが起動される度にイメージの `last_used_at` プロパティを更新することで、イメージの利用状況を記録しています。

## 自動更新

Incus はリモートサーバーからのイメージを自動的に最新に更新します。

```{note}
エイリアスを指定して取得したイメージだけが更新されます。
フィンガープリントを指定してイメージを取得した場合は、その特定のイメージバージョンを要求したことになります。
```

自動更新が有効になるかどうかはイメージをどのようにダウンロードしたかに依存します:

- インスタンス作成時にイメージがダウンロードとキャッシュされた場合は、ダウンロード時に {config:option}`server-images:images.auto_update_cached` が `true` に設定されていれば、自動的に更新されます。
- イメージがリモートサーバーから [`incus image copy`](incus_image_copy.md) コマンドでコピーされた場合は、`--auto-update`フラグが指定されていた場合のみ自動的に更新されます。

イメージのこの挙動は [`auto_update` プロパティを編集](images-manage-edit) することで変更できます。

起動時と [`images.auto_update_interval`](server-options-images) の間隔（デフォルトでは 6 時間ごと）を過ぎるたびに、Incus デーモンは自動更新とマークされコピー元のサーバーが記録されたストア内のすべてのイメージについてより新しいバージョンがあるかをチェックします。

新しいイメージが見つかったら、イメージ・ストアにダウンロードされます。
その後古いイメージを指していたエイリアスは新しいイメージを指すように変更され、古いイメージはストアから削除されます。

インスタンスの生成が遅くならないようにするため、Incus はキャッシュされたイメージからインスタンスを作成する際に新しいバージョンが利用可能かをチェックしません。
これはイメージが次の更新期間で更新されるまでの間は、新しく作成するインスタンスにイメージの古いバージョンが使われるかもしれないことを意味します。

## 特別なイメージプロパティ

プレフィックス`requirements`で始まるイメージプロパティ（たとえば、`requirements.XYZ`）は Incus がホストシステムと当該イメージで生成されるインスタンスの互換性を判断するために使用されます。
これらの互換性がない場合には、Incus はそのインスタンスを起動しません。

以下の要件がサポートされています:

% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group image-requirements start -->
    :end-before: <!-- config group image-requirements end -->
```
