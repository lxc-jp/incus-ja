(devices-tpm)=
# タイプ: `tpm`

```{note}
`tpm`デバイスタイプは、コンテナと VM の両方でサポートされています。
ただし、コンテナではホットプラグがサポートされていますが、VM ではサポートされていません。
```

TPM デバイスは、{abbr}`TPM (Trusted Platform Module)`エミュレータへのアクセスを有効にします。

TPM デバイスは、ブートプロセスを検証し、ブートチェーンのステップが改ざんされていないことを確認するために使用できます。また、暗号化キーを安全に生成および保存することもできます。

Incus は、TPM 2.0 をサポートするソフトウェア TPM を使用します。  
コンテナの主な使用例は、証明書のシールで、これによりキーがコンテナの外部に保存され、攻撃者がそれらを取得することがほぼ不可能になります。  
仮想マシンでは、TPM を証明書のシールに使用するだけでなく、ブートプロセスの検証にも使用できます。これにより、たとえば、Windows BitLocker と互換性のあるフルディスク暗号化が可能になります。

## デバイスオプション

`tpm`デバイスには以下のデバイスオプションがあります:

% Include content from [config_options.txt](../config_options.txt)
```{include} ../config_options.txt
    :start-after: <!-- config group devices-tpm start -->
    :end-before: <!-- config group devices-tpm end -->
```
