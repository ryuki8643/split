# プログラム解説

## 対応したオプション

- `-l`: 列分割のための列の数
- `-b`: バイト指定の分割のための文字列
- `-n`: ファイル個数分割のための文字列
- `-a`: ファイル名の桁数
- `-d`: ファイル名数字化
- 入力ファイル名
- prefix: 対応したイレギュラーな入力

## エラーハンドリング

- ファイル操作関連のエラー
- 対応したオプション以外の入力に対するエラー
- 各オプションで0以下の値が入力されたときのエラー
- オプション以外の入力が1つまたは2つ以外の時のエラー
- `b` オプションで100、100K、1000KBなどのフォーマットに合わない値が入力されたときのエラー（YBとZBについては未対応）
- `n` オプションで10、2/3、r/3、l/3、r/2/3、l/1/3等のフォーマットに合わない値が入力されたときのエラー
- 出力されるファイル数が `a` オプションで定められた範囲のファイル数を超えた時のエラー
- `l`, `n`, `b` のうち2つ以上のオプションが選択されたときのエラー

## パフォーマンスに関する工夫

- ファイルを読むときに行単位での分割なら `bufio.NewScanner` を使って1行ずつ、バイト単位の分割なら1Kずつ読み込み書き込み、都度バッファを開放することでメモリを節約して巨大ファイルに対応
- 行数を数える時も `Scanner` を使用
- `n` オプションの `CHUNK` の読み込みを最初正規表現で試みたが、`/` で split する方が早くてコードが書きやすいと判断して修正
- `n` オプションの `r` が最初についた時のラウンドロビンの書き込みについては、最初はファイルを書き込むたびに `os.Open` していたが、遅かったのと時折パニックが発生したため、一度 `Open` した後に `*os.File` を配列または変数として保存する方式に変更し、テスト時間が1秒以内に改善
- ファイル名の決定の計算量はアルファベット、数値、どちらでも桁数のオーダーに依存

## テストコードのポイント

- ファイル分割用のテストでは、ファイルの存在、容量または列の数、ファイル数、すべてのファイルの内容を足し合わせると元のファイルに戻るかを確認
- `flag` のテスト時、VS Code 上の UI で `run test` などをすると、他のオプションが付いてしまう問題を回避するため、`os.Args` に直接オプションを挿入しテスト可能にした
- `n` オプションのテストで `writer` に `&bytes.Buffer{}` を使い、標準出力もテスト
- テストカバレッジは 
