package entity

import (
	"sdmht/lib/utils"
)

const (
	OriginSquare     = 0  // 暗雾状态
	SquareExposeTime = 3  // 迷雾暴露时间
	MaxSquares       = 16 // 最大迷雾格子数

	MaxCostNum = 10 // 费用最大值
	FirstCost  = 7

	DefaultAttachCost = 3  // 攻击消耗费用
	MaxDrawCardTime   = 3  // 抽卡倒计时
	HandCardStartNum  = 3  // 起始手牌数
	HandCardMaxNum    = 10 // 手牌最大数
	DrawCardNum       = 3  // 抽牌数
)

type Scene struct {
	// 玩家ID
	PlayerID uint64 `json:"player_id"`
	// 主神单位ID
	MasterID int64 `json:"master_id"`
	// 单位
	Units map[int64]*Unit `json:"units"`
	// 迷雾 0:迷雾;+x为回到迷雾的倒计时;-x为不可开启的迷雾持续时间
	Squares []int32 `json:"squares"`
	// 单位位置
	UnitsLocation []int64 `json:"units_location"`
	// 手牌 存牌的编号
	HandCards []int64 `json:"hand_cards"`
	// 牌库 存牌的编号
	CardLibraries []int64 `json:"card_libraries"`
	// 牌库为空
	IsLibraryExpty bool `json:"is_library_empty"`
	// 附属神存活数
	RetainerAliveNum int32 `json:"retainer_alive_num"`
	// 牌库空之后的惩罚伤害
	LibraryExptyHurt int32 `json:"library_empty_hurt"`
	// 抽卡倒计时
	DrawCardCountDown int32 `json:"draw_card_count_down"`
	// 费用
	Cost int32 `json:"cost"`
	// 上一个移动的单位ID
	LastMoveUnitID int64 `json:"last_move_unit_id"`
}

func NewScene(playerID uint64, units []*Unit,
	cardLibrarys []int64, unitsLocation []int64) *Scene {
	cardLibraries := utils.SliceRandom(cardLibrarys).([]int64)
	scene := &Scene{
		PlayerID:          playerID,
		Units:             make(map[int64]*Unit),
		Squares:           make([]int32, MaxSquares),
		UnitsLocation:     unitsLocation,
		HandCards:         cardLibraries[:HandCardStartNum],
		CardLibraries:     cardLibraries[HandCardStartNum:],
		IsLibraryExpty:    false,
		RetainerAliveNum:  2,
		LibraryExptyHurt:  0,
		DrawCardCountDown: MaxDrawCardTime,
		Cost:              MaxCostNum,
	}
	for _, unit := range units {
		unit.Attack = unit.BaseAttack
		unit.Health = unit.MaxHealth
		unit.Move = unit.MaxMove
		scene.Units[unit.ID] = unit
		if unit.Rarity == RarityMainDeity {
			scene.MasterID = unit.ID
		}
	}

	return scene
}

func (s *Scene) NextRound() {
	s.Cost = MaxCostNum
	s.LastMoveUnitID = 0
	for _, unit := range s.Units {
		// unit.Skills
		unit.NextRound()
		if unit.Health == 0 {
			if unit.Rarity == 0 {
				return
			} else if unit.Rarity >= 1 && unit.Rarity <= 4 {
				s.RetainerAliveNum--
			}
			for i := range s.UnitsLocation {
				if s.UnitsLocation[i] == unit.ID {
					s.UnitsLocation[i] = 0
					break
				}
			}
			delete(s.Units, unit.ID)
		}
	}

	if s.IsLibraryExpty {
		s.LibraryExptyHurt += 2
	}

	for i, square := range s.Squares {
		if square < 0 {
			s.Squares[i]++
		} else if square > 0 {
			if s.RetainerAliveNum > 0 {
				s.Squares[i]--
			}
		}
	}

	s.WantToDrawCard()
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
func (s *Scene) RandomChangeSquare(num int, toExpose bool) {
	if num > MaxSquares {
		return
	}

	aliveMap := make(map[int64]*struct{})
	for _, unitID := range s.UnitsLocation {
		if unitID == 0 {
			continue
		}
		aliveMap[unitID] = &struct{}{}
	}

	exposed, unexposed := []int{}, []int{}               // 已暴露的，未暴露的
	unexposedAlive, unexposedNoAlive := []int{}, []int{} // 未暴露的存在单位的迷雾，未暴露的不存在单位的迷雾
	for i, square := range s.Squares {
		if square == 0 {
			unexposed = append(unexposed, i)
			if aliveMap[int64(i)] != nil {
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
