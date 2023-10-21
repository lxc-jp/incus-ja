# Incus

Incusは次世代のシステムコンテナおよび仮想マシンマネージャーです。

<!-- Include start Incus intro -->

コンテナや仮想マシンの中で動作する完全なLinuxシステムに統一されたユーザーエクスペリエンスを提供します。
Incusは数多くの Linuxディストリビューションのイメージ（公式のUbuntuイメージとコミュニティにより提供されるイメージ）を提供しており、非常にパワフルでありながら、それでいてシンプルなREST APIを中心に構築されています。
Incusは単一のマシン上の単一のインスタンスからデータセンターのフルラック内のクラスタまでスケールし、開発とプロダクションの両方のワークロードに適しています。

Incusを使えば小さなプライベートクラウドのように感じられるシステムを簡単にセットアップできます。
あなたのマシン資源を最適に利用しながら、あらゆるワークロードを効率よく実行できます。

さまざまな環境をコンテナ化したい場合や仮想マシンを稼働させたい場合、あるいは一般にあなたのインフラを費用効率よく稼働および管理したい場合にはIncusを使うのを検討するのがお勧めです。

[`https://linuxcontainers.org/incus/try-it/`](https://linuxcontainers.org/incus/try-it/)にてオンラインでIncusを試せます。

<!-- Include end Incus intro -->

## Canonical LXDのフォーク
Incusは、[Cumulonimbus incus](https://en.wikipedia.org/wiki/Cumulonimbus_incus)と鉄床雲にちなんで名づけられましたが、CanonicalのLXDのコミュニティによるフォークです。

このフォークはLinux ContainersコミュニティからLXDプロジェクトを[Canonicalが奪取](https://linuxcontainers.org/lxd/)したことに対する回答として作られました。

このフォークの主な目的は、みなさんの貢献が歓迎され、単一の営利団体がプロジェクトを管理することがないような、真のコミュニティプロジェクトを再び提供することです。

このフォークはLXD 5.16のリリースの時点で行われ、5.16を含むそれ以前のリリースのLXDからのアップグレードを可能にしています。
この時点以降2つのプロジェクトは分化していく可能性が高いので、それ以降のLXDリリースからのアップグレードはうまく動かないかもしれません。

Incusは今後もLXDの変更を監視し関連性のある変更は取り込む予定ですが、UbuntuやCanonical製品に特化した変更や機能は取り込まれない可能性が高いです。

## 使い始めるには

インストール手順と最初のステップはIncusドキュメントの[Getting started](https://linuxcontainers.org/incus/docs/main/getting_started/)（TODO：リンク先確認）を参照してください。

- リリースのアナウンス：[`https://discuss.linuxcontainers.org/c/news/`](https://discuss.linuxcontainers.org/c/news/)
- リリースtarballs：[`https://github.com/lxc/incus/releases/`](https://github.com/lxc/incus/releases/)
- ドキュメント： [`https://incus-ja.readthedocs.io/ja/latest/`](https://incus-ja.readthedocs.io/ja/latest/)（原文： [`https://linuxcontainers.org/incus/docs/main/`](https://linuxcontainers.org/incus/docs/main/)）

## ステータス

タイプ              | サービス              | ステータス
---                 | ---                   | ---
テスト              | GitHub                | [![Build Status](https://github.com/lxc/incus/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/lxc/incus/actions?query=event%3Apush+branch%3Amain)
Goドキュメント      | Godoc                 | [![GoDoc](https://godoc.org/github.com/lxc/incus/client?status.svg)](https://godoc.org/github.com/lxc/incus/client)
静的解析            | GoReport              | [![Go Report Card](https://goreportcard.com/badge/github.com/lxc/incus)](https://goreportcard.com/report/github.com/lxc/incus)

## セキュリティ

<!-- Include start security -->

Incusのインストールが安全であることを保証するために、以下の点を考慮してください。

- OSを最新に保ち、利用可能なすべてのセキュリティパッチをインストールしてください。
- サポートされているIncusのバージョンのみを使用してください。
- IncusデーモンとリモートAPIへのアクセスを制限してください。
- 必要とされない限り、特権コンテナを使わないでください。特権的なコンテナを使う場合は、適切なセキュリティ対策をしてください。詳細は[LXCセキュリティページ](https://linuxcontainers.org/ja/lxc/security/)を参照してください。
- ネットワークインタフェースを安全に設定してください。
<!-- Include end security -->

詳しい情報は[Security](https://github.com/lxc-jp/incus-ja/blob/main/doc/explanation/security.md)を参照してください。

**重要：**
<!-- Include start security note -->
UNIXソケットを介したIncusへのローカルアクセスは、常にIncusへのフルアクセスを許可します。
これは、任意のインスタンス上のセキュリティ機能を変更できる能力に加えて、任意のインスタンスにファイルシステムパスやデバイスをアタッチする能力を含みます。

したがって、あなたのシステムへのルートアクセスを信頼できるユーザーにのみ、このようなアクセスを与えるべきです。
<!-- Include end security note -->
<!-- Include start support -->

## サポートとコミュニティ

Incusコミュニティと交流するために以下のチャンネルが用意されています。

### バグレポート

バグレポートや機能要求は以下の場所で受け付けています。[`https://github.com/lxc/incus/issues/new`](https://github.com/lxc/incus/issues/new)

### コミュニティによるサポート

コミュニテイによるサポートは[`https://discuss.linuxcontainers.org`](https://discuss.linuxcontainers.org)で取り扱います。

### 商用サポート

[Debian or Ubuntu packages](https://github.com/zabbly/incus)の利用に関する商用サポートは現在[Zabbly](https://zabbly.com)から利用できます。

## ドキュメント

公式ドキュメントは[`https://github.com/lxc-jp/incus-ja/tree/main/doc`](https://github.com/lxc-jp/incus-ja/tree/main/doc)（原文：[`https://github.com/lxc/incus/tree/main/doc`](https://github.com/lxc/incus/tree/main/doc)）にあります。

<!-- Include end support -->

## コントリビュート

修正や新機能の提供は大歓迎です。まずは、[コントリビュートガイド](CONTRIBUTING.md)をお読みください！
