package models

import (
	"github.com/lib/pq"
)

func UniqArray(list []ForumUser, user, forum string,) []ForumUser {
	for _, u := range list {
		if user == u.UserPK {
			return list
		}
	}
	list = append(list, ForumUser{forum, user})
	return list

}


func parseError(err error) error {
	if err!=nil {
		switch err.(*pq.Error).Code {
		case pq.ErrorCode("23505"):
			return UniqueError
		case pq.ErrorCode("23503"):
			return FKConstraintError
		default:
			return err
		}
	}
	return nil
}