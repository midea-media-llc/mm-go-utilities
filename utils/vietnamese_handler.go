package utils

import (
	"strings"
)

func UnsignVietnamese(s string) string {

	var strSource = strings.Split("ă,â,đ,ê,ô,ơ,ư,à,ả,ã,ạ,á,ằ,ẳ,ẵ,ặ,ắ,ầ,ẩ,ẫ,ậ,ấ,è,ẻ,ẽ,ẹ,é,ề,ể,ễ,ệ,ế,ì,ỉ,ĩ,ị,í,ò,ỏ,õ,ọ,ó,ồ,ổ,ỗ,ộ,ố,ờ,ở,ỡ,ợ,ớ,ù,ủ,ũ,ụ,ú,ừ,ử,ữ,ự,ứ,ỳ,ỷ,ỹ,ỵ,ý,Ă,Â,Đ,Ê,Ô,Ơ,Ư,À,Ả,Ã,Ạ,Á,Ằ,Ẳ,Ẵ,Ặ,Ắ,Ầ,Ẩ,Ẫ,Ậ,Ấ,È,Ẻ,Ẽ,Ẹ,É,Ề,Ể,Ễ,Ệ,Ế,Ì,Ỉ,Ĩ,Ị,Í,Ò,Ỏ,Õ,Ọ,Ó,Ồ,Ổ,Ỗ,Ộ,Ố,Ờ,Ở,Ỡ,Ợ,Ớ,Ù,Ủ,Ũ,Ụ,Ú,Ừ,Ử,Ữ,Ự,Ứ,Ứ,Ỳ,Ỷ,Ỹ,Ỵ,Ý,Ứ,Ù", ",")

	// Mang cac ky tu thay the khong dau
	var strDest = strings.Split("a,a,d,e,o,o,u,a,a,a,a,a,a,a,a,a,a,a,a,a,a,a,e,e,e,e,e,e,e,e,e,e,i,i,i,i,i,o,o,o,o,o,o,o,o,o,o,o,o,o,o,o,u,u,u,u,u,u,u,u,u,u,y,y,y,y,y,A,A,D,E,O,O,U,A,A,A,A,A,A,A,A,A,A,A,A,A,A,A,E,E,E,E,E,E,E,E,E,E,I,I,I,I,I,O,O,O,O,O,O,O,O,O,O,O,O,O,O,O,U,U,U,U,U,U,U,U,U,U,Y,Y,Y,Y,Y,Y,U,U", ",")

	var sb strings.Builder

	for _, runeValue := range s {
		idx := find(strSource, string(runeValue))

		if idx >= 0 {
			sb.WriteString(string(strDest[idx]))
		} else {
			sb.WriteString(string(runeValue))
		}
	}

	return sb.String()
}

func find(source []string, value string) int {
	for idx, item := range source {
		if item == value {
			return idx
		}
	}
	return -1
}
