JIS X 0208
===
This package is a library for the Japanese character code JIS X 0208.

[Wikipedia](https://ja.wikipedia.org/wiki/JIS_X_0208) の JIS X 0208 項には以下のようにあります：

> JIS X 0208は、日本語表記、地名、人名などで用いられる6,879図形文字を含む、主として情報交換用の2バイト符号化文字集合を規定する日本産業規格 (JIS) である。現行の規格名称は7ビット及び8ビットの2バイト情報交換用符号化漢字集合 (7-bit and 8-bit double byte coded KANJI sets for information interchange) である。1978年にJIS C 6226として制定され、1983年、1990年および1997年に改正された。JIS漢字コード、JIS漢字、JIS第1第2水準漢字、JIS基本漢字などの通称がある。

このパッケージでは、文字が JIS X 0208 に含まれるかどうかを判定します。

参考にした資料は以下です：

* [図書館員のコンピュータ基礎講座:JIS X 0208コード表](https://www.asahi-net.or.jp/~ax2s-kmtn/ref/jisx0208.html)

JIS X 0208 の文字集合は `unicode.RangeTable` として定義していますので、直接利用可能です。また、いくつかの関数も定義してあります。
詳細は [ドキュメント](https://pkg.go.dev/github.com/ikawaha/jisx0208) や[ブログ](https://zenn.dev/ikawaha/articles/20210116-ab1ac4a692ae8bb4d9cf)を参照ください。

