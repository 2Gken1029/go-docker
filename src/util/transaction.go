package util

import (
	"gorm.io/gorm"
)

// TransactAndReturnData は、指定された GORM データベースに対してトランザクションを開始し、
// 渡された関数を実行する。関数の実行結果とトランザクションのエラーを返す。
// トランザクション内で発生した panic をキャッチし、ロールバックを行い、
// エラーがなく正常に終了した場合はトランザクションをコミットする。
//
// パラメータ:
//   - db: GORM データベースへのポインタ
//   - txFunc: トランザクション内で実行される関数
//
// 戻り値:
//   - data: トランザクション内で実行された関数の結果（省略可能）
//   - err: トランザクションのエラーまたは関数のエラー
//
// 注意: この関数はトランザクションを確実に処理するために defer を使用しており、
// トランザクション内で panic が発生した場合にはロールバックが行われる。
func TransactAndReturnData(db *gorm.DB, txFunc func(*gorm.DB) (interface{}, error)) (data interface{}, err error) {
	// トランザクションを開始
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error // トランザクションの開始に失敗した場合、エラーを返す
	}

	// 関数の実行後に必ず実行されるdefer文
	defer func() {
		// panicが発生した場合、トランザクションをロールバックしてpanicを再発生させる
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// エラーが発生した場合、トランザクションをロールバック
			tx.Rollback()
		} else {
			// エラーがなく、正常終了した場合、トランザクションをコミット
			err = tx.Commit().Error
		}
	}()

	// トランザクション内で渡された関数(txFunc)を実行し、結果をdataに格納
	data, err = txFunc(tx)
	return
}
