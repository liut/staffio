package models

import (
	"testing"
)

func TestSplitName(t *testing.T) {
	var (
		items = []struct{ cn, sn, gn string }{
			{"张飞", "张", "飞"},
			{"戏志才", "戏", "志才"},
			{"蔡迎慧 ", "蔡", "迎慧"},
			{"西门春雪", "西门", "春雪"},
			{"暂时不支持", "暂时不支持", ""},
		}
	)
	for _, n := range items {
		sn, gn := SplitName(n.cn)
		if sn == n.sn && gn == n.gn {
			t.Logf("%s %s", sn, gn)
		} else {
			t.Errorf("unexpect result %s %s", sn, gn)
		}

	}
}
