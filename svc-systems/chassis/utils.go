package chassis

import (
	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	dberrors "github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

func findOrNull(conn *persistencemgr.ConnPool, table, key string) (string, error) {
	r, e := conn.Read(table, key)
	if e != nil {
		switch e.ErrNo() {
		case dberrors.DBKeyNotFound:
			return "", nil
		default:
			return "", e
		}
	}
	return r, nil
}
