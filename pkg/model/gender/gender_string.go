// Code generated by "stringer -type=Gender"; DO NOT EDIT.

package gender

import "fmt"

const _Gender_name = "MasculineFeminineNeuterUnknown"

var _Gender_index = [...]uint8{0, 9, 17, 23, 30}

func (i Gender) String() string {
	if i < 0 || i >= Gender(len(_Gender_index)-1) {
		return fmt.Sprintf("Gender(%d)", i)
	}
	return _Gender_name[_Gender_index[i]:_Gender_index[i+1]]
}