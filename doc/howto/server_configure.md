(server-configure)=
# Incusサーバーを設定するには

Incus サーバーで利用可能なすべて設定オプションについては{ref}`server`を参照してください。

Incus サーバーがクラスタの一部の場合、一部のオプションはクラスタに適用され、また別のオプションはローカルサーバー、つまりクラスタメンバーにのみ適用されます。
{ref}`server`オプションの表で、クラスタに適用されるオプションは`global`スコープと表記され、ローカルサーバーのみに適用されるオプションは`local`スコープと表記されます。

## サーバーオプションを設定する

以下のコマンドでサーバーオプションを設定できます:

    incus config set <key> <value>

たとえば、ポート 8443 で Incus サーバーにリモートからのアクセスを許可するには、以下のコマンドを入力します:

    incus config set core.https_address :8443

クラスタ構成では、クラスタメンバーだけにサーバー設定を行うには`--target`フラグを追加してください。
たとえば、特定のクラスタメンバーでイメージの tarball を保管する場所を設定するには、以下のようなコマンドを入力してください:

    incus config set storage.images_volume my-pool/my-volume --target member02

## サーバー設定を表示する

現在のサーバー設定を表示するには、以下のコマンドを入力します:

    incus config show

クラスタ構成では、クラスタメンバーだけにサーバー設定を行うには`--target`フラグを追加してください。

## サーバー設定全体を編集する

サーバー設定全体を YAML ファイルとして編集するには、以下のコマンドを入力します:

    incus config edit

クラスタ構成では、クラスタメンバーだけにサーバー設定を行うには`--target`フラグを追加してください。
