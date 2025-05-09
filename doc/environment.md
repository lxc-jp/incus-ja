# 環境変数

以下の環境変数を設定することで、Incus のクライアントとデーモンをユーザーの環境に適合させることができ、いくつかの高度な機能を有効または無効にすることができます。

## クライアントとサーバー共通の環境変数

名前                 | 説明
:---                 | :----
`INCUS_DIR`          | Incusのデータディレクトリー
`INCUS_INSECURE_TLS` | trueに設定するとクライアント<->サーバー通信とサーバー<->イメージサーバーの両方（サーバー<->サーバーとクラスタは影響を受けない）ですべてのデフォルトのGoのcipherを許可する
`PATH`               | 実行ファイルの検索対象のパスのリスト
`http_proxy`         | HTTP用のプロキシサーバーのURL
`https_proxy`        | HTTPS用のプロキシサーバーのURL
`no_proxy`           | プロキシが不要なドメイン、IPアドレスあるいはCIDRレンジのリスト

## クライアントの環境変数

名前                | 説明
:---                | :----
`EDITOR`            | 使用するテキストエディタ
`INCUS_CONF`        | LXC設定ディレクトリーのパス
`INCUS_GLOBAL_CONF` | LXCグローバル設定ディレクトリーのパス
`INCUS_PROJECT`     | 使用するプロジェクトの名前（設定されたデフォルトのプロジェクトよりも優先されます）
`INCUS_REMOTE`      | 使用するリモートの名前（設定されたデフォルトのリモートよりも優先されます）
`VISUAL`            | （`EDITOR` が設定されてないときに）使用するテキストエディタ

## サーバーの環境変数

名前                            | 説明
:---                            | :----
`INCUS_AGENT_PATH`              | `incus-agent` のビルドを含むディレクトリのパス
`INCUS_CLUSTER_UPDATE`          | クラスタアップデートの際に呼ぶスクリプト
`INCUS_DEVMONITOR_DIR`          | デバイスモニターでモニターするパス。主にテスト用
`INCUS_DOCUMENTATION`           | ウェブサーバーで配信するドキュメントのパス
`INCUS_EDK2_PATH`               | `*_CODE.fd`と`*_VARS.fd`を含むEDK2ファームウェアビルドへのパス
`INCUS_EXEC_PATH`               | （サブコマンド実行時に使用される）Incus実行ファイルのフルパス
`INCUS_IDMAPPED_MOUNTS_DISABLE` | idmapを使ったマウントを無効にする（従来のUIDシフトを試す際に有用です）
`INCUS_LXC_TEMPLATE_CONFIG`     | LXCテンプレート設定ディレクトリー
`INCUS_SECURITY_APPARMOR`       | `false`に設定するとAppArmorを無効にします
`INCUS_SKIP_INSTANCE_TYPES`     | `true`に設定するとインスタンスタイプ定義のダウンロードをスキップします
`INCUS_UI`                      | ウェブサーバーを配信する web UI のパス
`INCUS_USBIDS_PATH`             | hwdata `usb.ids` ファイルのパス
