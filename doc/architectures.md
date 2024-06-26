(architectures)=
# アーキテクチャ

Incus は Linux カーネルと Go でサポートされるあらゆるアーキテクチャ上で稼働できます。

Incus の一部のエンティティ、たとえば、インスタンス、インスタンススナップショット、イメージはアーキテクチャに依存します。

下記のテーブルはサポートされるすべてのアーキテクチャを識別子と参照するための名前をリストアップします。
アーキテクチャ名は通常は Linux のカーネルアーキテクチャ名と揃えてあります。

ID   | カーネル名    | 注釈                             | パーソナリティ
:--- | :---          | :----                            | :------------
1    | `i686`        | 32bit Intel x86                  |
2    | `x86_64`      | 64bit Intel x86                  | `x86`
3    | `armv7l`      | 32bit ARMv7 リトルエンディアン   |
4    | `aarch64`     | 64bit ARMv8 リトルエンディアン   | `armv7l` (省略可能)
5    | `ppc`         | 32bit PowerPC ビッグエンディアン |
6    | `ppc64`       | 64bit PowerPC ビッグエンディアン | `powerpc`
7    | `ppc64le`     | 64bit PowerPC リトルエンディアン |
8    | `s390x`       | 64bit ESA/390 ビッグエンディアン |
9    | `mips`        | 32bit MIPS                       |
10   | `mips64`      | 64bit MIPS                       | `mips`
11   | `riscv32`     | 32bit RISC-V リトルエンディアン  |
12   | `riscv64`     | 64bit RISC-V リトルエンディアン  |
13   | `armv6l`      | 32bit ARMv6 リトルエンディアン   |
14   | `armv8l`      | 32bit ARMv8 リトルエンディアン   |
15   | `loongarch64` | 64bit Loongarch                  |

```{note}
Incus はカーネルアーキテクチャのみに影響し、ツールチェインで決定される特定のユーザースペースのフレーバーには影響しません。

これは Incus は ARMv7 hard-float を ARMv7 soft-float と同じとして扱い、両方を`armv7l`として参照することを意味します。
もしユーザーにとって有用であれば、正確なユーザースペースのABIがイメージとコンテナプロパティとして設定でき、簡単に問い合わせできます。
```

## 仮想マシンのサポート

Incus は以下のホストアーキテクチャーでのみ仮想マシンの動作をサポートします:

- `x86_64`
- `aarch64`
- `ppc64le`
- `s390x`

仮想マシンのファームウェアがブート可能であれば、仮想マシンのアーキテクチャは通常ホストアーキテクチャーの 32bit パーソナリティにすることもできます。
