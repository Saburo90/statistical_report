package components

import (
	"errors"
	"fmt"
	"gitee.com/NotOnlyBooks/statistical_report/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
	"go.uber.org/zap"
)

var (
	DBConf  *conf.DBConfig
	DBGroup *xorm.EngineGroup
)

type defaultDBLoggerImpl struct {
	c *conf.DBConfig
}

func NewDBSession() *xorm.Session {
	return DBGroup.NewSession()
}

func SetupDatabase(dbConf *conf.DBConfig) error {
	DBConf = dbConf

	var (
		err    error
		master *xorm.Engine
		slaves []*xorm.Engine
	)

	for role, dsn := range dbConf.Connects {
		if role == "master" && master == nil {
			if master, err = xorm.NewMySQL(xorm.MYSQL_DRIVER, dsn); err != nil {
				return err
			}
		} else if role == "slave" {
			if slave, err := xorm.NewMySQL(xorm.MYSQL_DRIVER, dsn); err != nil {
				return err
			} else {
				slaves = append(slaves, slave)
			}
		}
	}

	if master == nil {
		return errors.New("master database can not be configured")
	}

	if len(slaves) == 0 {
		slaves = append(slaves, master)
	}

	if DBGroup, err = xorm.NewEngineGroup(master, slaves); err != nil {
		return err
	}

	if DBConf.Debug {
		DBGroup.ShowExecTime(true)
		DBGroup.ShowSQL(true)
		DBGroup.SetLogLevel(core.LOG_DEBUG)
	} else {
		DBGroup.ShowSQL(false)
		DBGroup.SetLogLevel(core.LOG_WARNING)
	}

	DBGroup.SetLogger(&defaultDBLoggerImpl{c: dbConf})

	return nil
}

func (log *defaultDBLoggerImpl) Debug(v ...interface{}) {
	zap.L().Debug(fmt.Sprint(v...))
}

func (log *defaultDBLoggerImpl) Debugf(format string, v ...interface{}) {
	zap.L().Debug(fmt.Sprintf(format, v...))
}

func (log *defaultDBLoggerImpl) Error(v ...interface{}) {
	zap.L().Error(fmt.Sprint(v...))
}

func (log *defaultDBLoggerImpl) Errorf(format string, v ...interface{}) {
	zap.L().Error(fmt.Sprintf(format, v...))
}

func (log *defaultDBLoggerImpl) Info(v ...interface{}) {
	zap.L().Info(fmt.Sprint(v...))
}

func (log *defaultDBLoggerImpl) Infof(format string, v ...interface{}) {
	zap.L().Info(fmt.Sprintf(format, v...))
}

func (log *defaultDBLoggerImpl) Warn(v ...interface{}) {
	zap.L().Warn(fmt.Sprint(v...))
}

func (log *defaultDBLoggerImpl) Warnf(format string, v ...interface{}) {
	zap.L().Warn(fmt.Sprintf(format, v...))
}

func (log *defaultDBLoggerImpl) Level() core.LogLevel {
	if log.c.Debug {
		return core.LOG_DEBUG
	}

	return core.LOG_WARNING
}

func (log *defaultDBLoggerImpl) SetLevel(l core.LogLevel) { DBGroup.SetLogLevel(l) }

func (log *defaultDBLoggerImpl) ShowSQL(show ...bool) {}

func (log *defaultDBLoggerImpl) IsShowSQL() bool {
	return log.c.Debug
}
