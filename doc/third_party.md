# サードパーティーツールと統合
以下はネイティブまたはプラグイン経由で Incus をサポートする一般的な運用ツールの一覧です。

## Terraform / OpenTofu
[Terraform](https://www.terraform.io)と[OpenTofu](https://opentofu.org)は infrastructure as code のツールでインフラストラクチャー自体を作成することにフォーカスしたツールです。
Incus にkなしては、これはプロジェクト、プロファイル、ネットワーク、ストレージボリューム、そしてインスタンスを作成できることを意味します。

たいていの場合、一旦インスタンスとインスタンスに必要なその他すべてが用意されたら、その後の環境構築は Ansible を使うでしょう。

Incus との統合は[Incus 専用のプロバイダー](https://github.com/lxc/terraform-provider-incus)により実現されます。

## Ansible
[Ansible](https://www.ansible.com) は infrastructure as code のツールで、特にソフトウェアのプロビジョニングと構成管理にフォーカスしています。
ほとんどの作業をソフトウェアをデプロイする対象のシステムにまず接続してから行います。

そのために、Ansible は SSH やその他さまざまなプロトコルで接続ができ、そのうちの 1 つが [Incus](https://docs.ansible.com/ansible/latest/collections/community/general/incus_connection.html) です。

これにより、最初に SSH をセットアップすることなしに Incus のインスタンス内にソフトウェアを簡単にデプロイできます。

## Packer
[Packer](https://www.packer.io) はカスタム OS イメージを幅広い様々なプラットフォーム用に生成するツールです。

Packer が Incus のイメージを直接生成できるようにする [プラグイン](https://developer.hashicorp.com/packer/integrations/bketelsen/incus) が存在します。

## Distrobuilder
[Distrobuilder](https://github.com/lxc/distrobuilder) は公式の LXC と Incus イメージを生成することで良く知られているイメージ生成ツールです。
イメージの YAML 定義を読み込んで LXC コンテナーのイメージと Incus のコンテナーと VM のイメージを生成します。

Distrobuilder は既存のイメージを再パッケージするというよりも、クリーンなイメージをスクラッチから生成することにフォーカスしています。

## GARM
[GARM](https://github.com/cloudbase/garm) GitHub の自己ホストランナーを稼働できる Github Actions Runner Manager です。

[Incus](https://github.com/cloudbase/garm-provider-incus)を含む、さまざまな種類のランナーのプロバイダーを提供しています。
