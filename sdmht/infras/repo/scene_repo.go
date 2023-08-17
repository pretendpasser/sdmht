package repo

import (
	"sdmht/lib/utils"
)

const (
	MainDeity       = 1 //主神
	AffiliatedDeity = 2 //附属神明
	Derivative      = 3 //衍生物

	OriginSquare     = 0
	SquareExposeTime = 3 // 迷雾暴露时间

	MaxDrawCardTime = 3
)

// 场景
type Scene struct {
	Square            [16]int // 0:迷雾;+x为回到迷雾的倒计时;-x为不可开启的迷雾持续时间
	HandCard          [10]int // 手牌 存牌的编号
	CardLibrary       [20]int // 牌库 存牌的编号
	DrawCardCountDown int     // 抽卡倒计时
}

// 操作全场迷雾
func (s *Scene) OperatorAllSquare(time int) {
	if time <= 0 {
		s.Square = [16]int{
			time, time, time, time,
			time, time, time, time,
			time, time, time, time,
			time, time, time, time,
		}
	} else {
		for i, v := range s.Square {
			if v != 0 {
				continue
			}
			s.Square[i] = 3
		}
	}
}

// 随机 num 数量的开雾/盖雾
func (s *Scene) RandomChangeSquare(num int, toExpose bool, alives []int) {
	if num > 16 {
		return
	}

	aliveMap := make(map[int]*struct{})
	for _, alive := range alives {
		aliveMap[alive] = &struct{}{}
	}

	exposed, unexposed := []int{}, []int{}               // 已暴露的，未暴露的
	unexposedAlive, unexposedNoAlive := []int{}, []int{} // 未暴露的存在单位的迷雾，未暴露的不存在单位的迷雾
	for i, square := range s.Square {
		if square == 0 {
			unexposed = append(unexposed, i)
			if aliveMap[i] != nil {
				unexposedAlive = append(unexposedAlive, i)
			} else {
				unexposedNoAlive = append(unexposedNoAlive, i)
			}
		} else if square > 0 {
			exposed = append(exposed, i)
		}
	}

	if toExpose {
		//  开雾数 >= 迷雾数
		if num >= len(unexposed) {
			// 迷雾中不全是单位
			if len(unexposedAlive) != 0 && len(unexposedNoAlive) != 0 {
				// 开所有无单位的迷雾
				for _, square := range unexposedNoAlive {
					s.Square[square] = SquareExposeTime
				}
				return
			}

			// 迷雾中全 有/没有 单位
			// 开全部迷雾
			for _, square := range unexposed {
				s.Square[square] = SquareExposeTime
			}
			return
		}

		// 无单位的迷雾数 >= 开雾数
		if len(unexposedNoAlive) >= num {
			utils.SliceRandom(&unexposedNoAlive)
			for _, square := range unexposedNoAlive[:num] {
				s.Square[square] = SquareExposeTime
			}
			return
		} else {
			num -= len(unexposedNoAlive)
			utils.SliceRandom(&unexposedAlive)
			for i := 0; i < num; i++ {
				unexposedNoAlive = append(unexposedNoAlive, unexposedAlive[i])
			}
			utils.SliceRandom(&unexposedNoAlive)
			for _, square := range unexposedNoAlive {
				s.Square[square] = SquareExposeTime
			}
		}
	} else {
		if len(exposed) <= num {
			for _, square := range exposed {
				s.Square[square] = OriginSquare
			}
			return
		}
		utils.SliceRandom(&exposed)
		for i := 0; i < num; i++ {
			s.Square[exposed[i]] = OriginSquare
		}
	}
}

// 获取 num 数量 已开/未开 的迷雾，不足时从另一边取
func (s *Scene) RandomGetSquare(num int, isExpose bool) []int {
	res := []int{}
	exposed, unexposed := []int{}, []int{} //已暴露的，未暴露的
	for i, square := range s.Square {
		if square <= 0 {
			unexposed = append(unexposed, i)
		} else if square > 0 {
			exposed = append(exposed, i)
		}
	}

	if isExpose {
		exposedLength := len(exposed)
		if exposedLength == num {
			res = exposed
		} else if exposedLength > num {
			utils.SliceRandom(&exposed)
			res = exposed[:num]
		} else {
			res = append(res, exposed...)
			num -= exposedLength
			utils.SliceRandom(&unexposed)
			res = append(res, unexposed[:num]...)
		}
	} else {
		unexposedLength := len(unexposed)
		if unexposedLength == num {
			res = unexposed
		} else if unexposedLength > num {
			utils.SliceRandom(&unexposed)
			res = unexposed[:num]
		} else {
			res = append(res, unexposed...)
			num -= unexposedLength
			utils.SliceRandom(&exposed)
			res = append(res, exposed[:num]...)
		}
	}
	utils.SliceRandom(&res)
	return res
}
