package helper

import "github.com/sirupsen/logrus"

type FinalizeTransactionOpt struct {
	Commit   func() error
	Rollback func() error

	Err *error

	LogTag string
}

func FinalizeTransaction(opt FinalizeTransactionOpt) {
	logTag := opt.LogTag + "[FinalizeTransaction]"

	if p := recover(); p != nil {
		_ = opt.Rollback()
		logrus.WithField("panic", p).Error("transaction rolled back due to panic")
		panic(p)
	}

	if opt.Err != nil {
		errRollback := opt.Rollback()
		if errRollback != nil {
			logrus.WithFields(logrus.Fields{
				"error":          *opt.Err,
				"error_rollback": errRollback,
			}).Errorf("%s[opt.Rollback] Error during transaction", logTag)
		}
		return
	}

	errCommit := opt.Commit()
	if errCommit != nil {
		logrus.WithFields(logrus.Fields{
			"error_commit": errCommit,
		}).Errorf("%s[opt.Commit] Error during transaction", logTag)
	}
}
