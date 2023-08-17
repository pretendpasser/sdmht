package repo

import "sdmht/lib"

// 单位属性
type BaseAttribute struct {
	Location      int  // 位置
	Type          uint // 类型
	Attack        uint // 攻击力
	Health        uint // 生命值
	Move          uint // 移动力
	AttackPrevent bool // 攻击防止(圣盾)
	NoMove        uint // 禁止移动
	NoAttack      uint // 禁止攻击
}

func (a *BaseAttribute) Moveing(direction string) error {
	switch direction {
	case "up":
		if 0 <= a.Location && a.Location < 4 {
			return lib.NewError(lib.ErrInvalidArgument, "can not move up")
		}
		a.Location -= 4
	case "down":
		if 12 <= a.Location && a.Location < 15 {
			return lib.NewError(lib.ErrInvalidArgument, "can not move down")
		}
		a.Location += 4
	case "left":
		if a.Location%4 == 0 {
			return lib.NewError(lib.ErrInvalidArgument, "can not move left")
		}
		a.Location -= 1
	case "right":
		if a.Location%4 == 3 {
			return lib.NewError(lib.ErrInvalidArgument, "can not move right")
		}
		a.Location += 1
	}
	return nil
}
