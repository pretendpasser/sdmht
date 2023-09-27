package entity

import (
	"sdmht/lib/utils"
)

const (
	OriginSquare     = 0
	SquareExposeTime = 3  // 迷雾暴露时间
	MaxSquares       = 16 // 最大迷雾数

	MaxDrawCardTime  = 3  // 抽卡倒计时
	HandCardStartNum = 3  // 起始手牌数
	HandCardMaxNum   = 10 // 手牌最大数
	DrawCardNum      = 3  // 抽牌数
)

type Scene struct {
	// 0:迷雾;+x为回到迷雾的倒计时;-x为不可开启的迷雾持续时间
	Squares []int32 `json:"squares"`
	// 手牌 存牌的编号
	HandCards []int64 `json:"hand_cards"`
	// 牌库 存牌的编号
	CardLibraries []int64 `json:"card_libraries"`
	// 牌库为空
	IsLibraryExpty bool `json:"is_library_empty"`
	// 牌库空之后的惩罚伤害
	LibraryExptyHurt int32 `json:"library_empty_hurt"`
	// 抽卡倒计时
	DrawCardCountDown int32 `json:"draw_card_count_down"`
}

func NewScene(cardLibrarys []int64) *Scene {
	cardLibraries := utils.SliceRandom(cardLibrarys).([]int64)
	return &Scene{
		Squares:           make([]int32, MaxSquares),
		HandCards:         cardLibraries[:HandCardStartNum],
		CardLibraries:     cardLibraries[HandCardStartNum:],
		IsLibraryExpty:    false,
		LibraryExptyHurt:  0,
		DrawCardCountDown: MaxDrawCardTime,
	}
}

func (s *Scene) NextRound() {
	if s.IsLibraryExpty {
		s.LibraryExptyHurt += 2
	}
	s.WantToDrawCard()

	for i, square := range s.Squares {
		if square < 0 {
			s.Squares[i]++
		} else if square > 0 {
			s.Squares[i]--
		}
	}
}

func (s *Scene) WantToDrawCard() {
	s.DrawCardCountDown--
	if !s.IsLibraryExpty && s.DrawCardCountDown == 0 {
		cardLibraryLength := len(s.CardLibraries)
		if cardLibraryLength > DrawCardNum {
			s.HandCards = append(s.HandCards, s.CardLibraries[:DrawCardNum]...)
			if len(s.HandCards) > HandCardMaxNum {
				s.HandCards = s.HandCards[:HandCardMaxNum]
			}
			s.CardLibraries = s.CardLibraries[DrawCardNum:]
		} else {
			s.HandCards = append(s.HandCards, s.CardLibraries[:cardLibraryLength]...)
			if len(s.HandCards) > HandCardMaxNum {
				s.HandCards = s.HandCards[:HandCardMaxNum]
			}
			s.CardLibraries = []int64{}
			s.IsLibraryExpty = true
		}
	}

	if s.DrawCardCountDown == 0 {
		s.DrawCardCountDown = MaxDrawCardTime
	}
}

// 操作全场迷雾
func (s *Scene) OperatorAllSquare(time int32) {
	if time <= 0 {
		s.Squares = []int32{
			time, time, time, time,
			time, time, time, time,
			time, time, time, time,
			time, time, time, time,
		}
	} else {
		for i, square := range s.Squares {
			if square != 0 {
				continue
			}
			s.Squares[i] = SquareExposeTime
		}
	}
}

// 随机 num 数量的开雾/盖雾
func (s *Scene) RandomChangeSquare(num int, toExpose bool, alives []int) {
	if num > MaxSquares {
		return
	}

	aliveMap := make(map[int]*struct{})
	for _, alive := range alives {
		aliveMap[alive] = &struct{}{}
	}

	exposed, unexposed := []int{}, []int{}               // 已暴露的，未暴露的
	unexposedAlive, unexposedNoAlive := []int{}, []int{} // 未暴露的存在单位的迷雾，未暴露的不存在单位的迷雾
	for i, square := range s.Squares {
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
					s.Squares[square] = SquareExposeTime
				}
				return
			}

			// 迷雾中全 有/没有 单位
			// 开全部迷雾
			for _, square := range unexposed {
				s.Squares[square] = SquareExposeTime
			}
			return
		}

		// 无单位的迷雾数 >= 开雾数
		if len(unexposedNoAlive) >= num {
			unexposedNoAlive = utils.SliceRandom(unexposedNoAlive).([]int)
			for _, square := range unexposedNoAlive[:num] {
				s.Squares[square] = SquareExposeTime
			}
			return
		} else {
			num -= len(unexposedNoAlive)
			unexposedAlive = utils.SliceRandom(unexposedAlive).([]int)
			for i := 0; i < num; i++ {
				unexposedNoAlive = append(unexposedNoAlive, unexposedAlive[i])
			}
			unexposedNoAlive = utils.SliceRandom(unexposedNoAlive).([]int)
			for _, square := range unexposedNoAlive {
				s.Squares[square] = SquareExposeTime
			}
		}
	} else {
		if len(exposed) <= num {
			for _, square := range exposed {
				s.Squares[square] = OriginSquare
			}
			return
		}
		exposed = utils.SliceRandom(exposed).([]int)
		for i := 0; i < num; i++ {
			s.Squares[exposed[i]] = OriginSquare
		}
	}
}

// 获取 num 数量 已开/未开的迷雾，不足时从另一边取
// num 为 0 时，获取所有的  已开/未开 的迷雾
func (s *Scene) RandomGetSquare(num int, isExpose bool) []int {
	res := []int{}
	exposed, unexposed := []int{}, []int{} //已暴露的，未暴露的
	for i, square := range s.Squares {
		if square <= 0 {
			unexposed = append(unexposed, i)
		} else if square > 0 {
			exposed = append(exposed, i)
		}
	}
	exposed = utils.SliceRandom(exposed).([]int)
	unexposed = utils.SliceRandom(unexposed).([]int)

	if isExpose {
		if num == 0 {
			return exposed
		}
		exposedLength := len(exposed)
		if exposedLength == num {
			res = exposed
		} else if exposedLength > num {
			res = exposed[:num]
		} else {
			res = append(res, exposed...)
			num -= exposedLength
			res = append(res, unexposed[:num]...)
		}
	} else {
		if num == 0 {
			return unexposed
		}
		unexposedLength := len(unexposed)
		if unexposedLength == num {
			res = unexposed
		} else if unexposedLength > num {
			res = unexposed[:num]
		} else {
			res = append(res, unexposed...)
			num -= unexposedLength
			res = append(res, exposed[:num]...)
		}
	}

	return res
}
