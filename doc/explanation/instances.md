(expl-instances)=
# インスタンスについて

Incus は以下のインスタンスタイプをサポートします:

システムコンテナ
: システムコンテナは共有されたカーネルを使って完全なLinuxディストリビューションを稼働します。
  これらのコンテナは完全はLinuxディストリビューションを実行し、仮想マシンと非常によく似ていますが、ホストシステムとカーネルを共有する点が異なります。

  これらは非常にオーバーヘッドが低く、非常にコンパクトにパッケージされ、概して仮想マシンとほぼ同一の体験を提供しますが、仮想マシンと違って必須のハードウェアサポートやオーバーヘッドがありません。

  システムコンテナは`liblxc` (LXC)を使って実装されています。

アプリケーションコンテナ
: アプリケーションコンテナはビルド済みのイメージを使って単一のアプリケーションを実行します。
  この種のコンテナはDockerやKubernetesなどによって普及しました。

  純粋なLinux環境を提供してその上にソフトウェアをインストールする必要なしに、インストール済みの、そしてたいていは事前に設定済みのソフトウェアが含まれます。

  Incusは任意のOCI互換のイメージレポジトリ（例えば Docker Hub）からアプリケーションコンテナイメージを取得し、利用できます。

  アプリケーションコンテナは`umoci`と`skopeo`の助けを借りて`liblxc` (LXC)を使って実装されています。

仮想マシン
: {abbr}`Virtual machines (VMs)`は完全に仮想化されたシステムです。
  仮想マシンはIncusでネイティブにサポートされていて、システムコンテナと別の選択肢を提供します。

  コンテナ内で正常に動作しないものもあります。異なるカーネルを必要とするものや独自のカーネルモジュールはコンテナではなく仮想マシン内で動かす必要があります。

  同様に、完全なPCIデバイスなど、ある種のデバイスパススルーは仮想マシンでのみ正しく動きます。

  ユーザー体験を一貫性のあるものにするために、組み込みのエージェントがIncusにより提供されており、インタラクティブなコマンド実行やファイル転送ができます。

  仮想マシンはQEMUを使って実装されています。

  ```{note}
  現状、仮想マシンはコンテナよりサポートする機能が少ないですが、将来には両方のインスタンスタイプで同じ機能セットをサポートする計画です。

  仮想マシンでどの機能が利用可能かを見るには、{ref}`instance-options` ドキュメントの条件のカラムを確認してください。
  ```

インスタンスタイプのより詳細な情報は{ref}`containers-and-vms`を参照してください。
