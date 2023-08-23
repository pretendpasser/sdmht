package repo

const (
	MainDeity       = 1 //主神
	AffiliatedDeity = 2 //附属神明
	Derivative      = 3 //衍生物
)

// 单位属性
type BaseAttribute struct {
	Location      int         // 位置
	Type          uint        // 类型
	Attack        int         // 攻击力
	TempAttact    map[int]int // 临时攻击力变化 [数值]回合数
	Defend        uint        // 护盾
	MaxDefend     uint        // 最大护盾值
	Health        uint        // 生命值
	MaxHealth     uint        // 最大生命值
	Move          int         // 移动力
	TempMove      map[int]int // 临时移动力变化 [数值]回合数
	AttackPrevent bool        // 攻击防止(圣盾)
	Hurt          map[int]int //创伤层数 [数值]回合数
	Cure          map[int]int //治疗层数 [数值]回合数
	NoMove        uint        // 禁止移动
	NoAttack      uint        // 禁止攻击
	NoCure        uint        // 禁止治疗
}

func (a *BaseAttribute) CheckMoveing(location int) error {
	return nil
}

// func (a *BaseAttribute) Moveing(direction string) error {
// 	tempMoves := 0
// 	for _, tempMove := range a.TempMove {
// 		tempMoves += tempMove
// 	}
// 	if a.NoMove > 0 || a.Move+tempMoves <= 0 {
// 		return lib.NewError(lib.ErrInternal, "")
// 	}
// 	switch direction {
// 	case "up":
// 		if 0 <= a.Location && a.Location < 4 {
// 			return lib.NewError(lib.ErrInvalidArgument, "can not move up")
// 		}
// 		a.Location -= 4
// 	case "down":
// 		if 12 <= a.Location && a.Location < 15 {
// 			return lib.NewError(lib.ErrInvalidArgument, "can not move down")
// 		}
// 		a.Location += 4
// 	case "left":
// 		if a.Location%4 == 0 {
// 			return lib.NewError(lib.ErrInvalidArgument, "can not move left")
// 		}
// 		a.Location -= 1
// 	case "right":
// 		if a.Location%4 == 3 {
// 			return lib.NewError(lib.ErrInvalidArgument, "can not move right")
// 		}
// 		a.Location += 1
// 	}
// 	return nil
// }
