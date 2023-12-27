package entity

import (
	"sdmht/lib"
	"sdmht/lib/utils"
)

const (
	AffiliateNanhua       = 0 // 南华庄势力
	AffiliateAsgard       = 1 // 阿斯加德势力
	AffiliateTakamagahara = 2 // 高天原势力
	AffiliateNilfheim     = 3 // 尼夫尔海姆势力
	AffiliateTaiqing      = 4 // 太清殿势力
	AffiliateAmaravati    = 5 // 阿摩婆罗提势力
	AffiliateHeliopolis   = 6 // 赫利奥波利斯势力
	AffiliateOlympus      = 7 // 奥林匹斯势力
	AffiliateYuqing       = 8 // 玉清殿势力

	RarityMainDeity  = 0 // 主神
	RarityOrdinary   = 1 // 普通
	RarityDare       = 2 // 稀有
	RarityLegend     = 3 // 传说
	RarityObsidian   = 4 // 黑曜
	RarityDerivative = 5 // 衍生物
)

type Unit struct {
	BaseAttribute

	// Weapon        int64  `json:"weapon"`         // 弹幕
	// Trap          int64  `json:"trap"`           // 秘术
	Health        int32 `json:"health"`         // 生命值
	Defend        int32 `json:"defend"`         // 护盾
	Attack        int32 `json:"attack"`         // 攻击力
	Move          int32 `json:"move"`           // 移动力
	IsMoving      int32 `json:"is_moving"`      // 移动中 [0:未移动 1:移动中 -1:移动结束]
	AttackPrevent bool  `json:"attack_prevent"` // 攻击防止(圣盾)
	// CounterAttack bool   `json:"counter_attack"` // 反击
	// Rebirth       bool   `json:"rebirth"`        // 重生
	HurtInstead int32 `json:"hurt_instead"` // 坚壁(保护的单位的unitID)
	NoMove      int32 `json:"no_move"`      // 临时禁止移动
	NoAttack    int32 `json:"no_attack"`    // 临时禁止攻击
	NoCure      int32 `json:"no_cure"`      // 临时禁止治疗

	ChangeAttack []*TempAttribute `json:"change_attack"` // 临时攻击力变化 [数值]持续时间
	ChangeMove   []*TempAttribute `json:"change_move"`   // 临时移动力变化 [数值]持续时间
	Cure         []*TempAttribute `json:"cure"`          // 治疗层数 [数值]持续时间
	Hurt         []*TempAttribute `json:"hurt"`          // 创伤层数 [数值]持续时间
}

type TempAttribute struct {
	Value  int32 `json:"value"`  // 数值
	Period int32 `json:"period"` // 持续时间
}

// 单位属性
type BaseAttribute struct {
	ID           int64             `json:"id" db:"id"`
	Name         string            `json:"name" db:"name"`
	SkillName    string            `json:"-" db:"skill_name"`
	Skills       map[string]*Skill `json:"skills" db:"-"`
	Rarity       int32             `json:"rarity" db:"rarity"`       // 稀有度
	Affiliate    int32             `json:"affiliate" db:"affiliate"` // 所属势力
	BaseAttack   int32             `json:"base_attack" db:"attack"`  // 基础攻击力
	MaxDefend    int32             `json:"max_defend" db:"defend"`   // 最大护盾值
	MaxHealth    int32             `json:"max_health" db:"health"`   // 最大生命值
	MaxMove      int32             `json:"max_move" db:"move"`       // 最大移动力
	BaseNoMove   bool              `json:"base_no_move" db:"-"`      // 禁止移动
	BaseNoAttack bool              `json:"base_no_attack" db:"-"`    // 禁止攻击
	BaseNoCure   bool              `json:"base_no_cure" db:"-"`      // 禁止治疗
	BaseNoEquip  bool              `json:"base_no_equip" db:"-"`     // 禁止装备
}

type Skill struct {
	TotalUseCnt int32        `json:"total_use_cut"`
	RoundUseCnt int32        `json:"round_use_cut"`
	Desc        string       `json:"desc"`
	Handler     SkillHandler `json:"-"`
}

func (u *Unit) NextRound() {
	u.IsMoving = 0
	u.Move = u.BaseAttribute.MaxMove
	if u.NoMove > 0 {
		u.NoMove--
	}
	if u.NoAttack > 0 {
		u.NoAttack--
	}
	if u.NoCure > 0 {
		u.NoCure--
	}
	// 计算攻击力
	newChangeAttack := []*TempAttribute{}
	newAttack := u.Attack
	for _, t := range u.ChangeAttack {
		t.Period -= 1
		if t.Period > 0 {
			newChangeAttack = append(newChangeAttack, t)
			newAttack += t.Value
		}
	}
	u.ChangeAttack = newChangeAttack
	if newAttack >= 0 {
		u.Attack = newAttack
	} else {
		u.Attack = 0
	}
	// 计算移动力
	newChangeMove := []*TempAttribute{}
	newMove := u.Move
	for _, t := range u.ChangeMove {
		t.Period -= 1
		if t.Period > 0 {
			newChangeMove = append(newChangeMove, t)
			newMove += t.Value
		}
	}
	u.ChangeMove = newChangeMove
	if newMove >= 0 {
		u.Move = newMove
	} else {
		u.Move = 0
	}
	// 计算治疗值
	newCure := []*TempAttribute{}
	newHealth := u.Health
	for _, t := range u.Cure {
		t.Period -= 1
		if t.Period > 0 {
			newCure = append(newCure, t)
			newHealth += t.Value
		}
	}
	u.Cure = newCure
	if newHealth >= u.MaxHealth {
		u.Health = u.MaxHealth
	} else {
		u.Health = newHealth
	}
	// 计算创伤值
	newHurt := []*TempAttribute{}
	var hurtValue int32 = 0
	for _, t := range u.ChangeMove {
		t.Period -= 1
		if t.Period > 0 {
			newHurt = append(newHurt, t)
			hurtValue += t.Value
		}
	}
	u.Hurt = newHurt
	if u.Defend >= hurtValue {
		u.Defend -= hurtValue
	} else {
		hurtValue -= u.Defend
		u.Defend = 0
		if u.Health >= hurtValue {
			u.Health -= newHealth
		} else {
			u.Health = 0
		}
	}
}

func CheckMoveing(from, to int64) error {
	// 数列阶数
	var order int64 = 4
	// 0  1  2  3
	// 4  5  6  7
	// 8  9  10 11
	// 12 13 14 15

	if from < 0 || from >= order*order ||
		to < 0 || to >= order*order {
		return lib.NewError(lib.ErrInvalidArgument, "from or to is overrange")
	}

	if from == to {
		return lib.NewError(lib.ErrInvalidArgument, "invalid moving none")
	}

	if to == from+1 {
		if from%order == order-1 {
			return lib.NewError(lib.ErrInvalidArgument, "invalid moving right")
		} else {
			return nil
		}
	} else if to == from-1 {
		if from%order == 0 {
			return lib.NewError(lib.ErrInvalidArgument, "invalid moving left")
		} else {
			return nil
		}
	} else if to == from-order {
		if from%order == 0 {
			return lib.NewError(lib.ErrInvalidArgument, "invalid moving up")
		} else {
			return nil
		}
	} else if to == from+order {
		if from/order >= order-1 && from/order < order {
			return lib.NewError(lib.ErrInvalidArgument, "invalid moving down")
		} else {
			return nil
		}
	} else {
		return lib.NewError(lib.ErrInvalidArgument, "invalid moving over one")
	}
}

type UnitQuery struct {
	Pagination *utils.Pagination `json:"page"`

	ExcludeMainDeity  bool  `json:"exclude_main_deity"`  // 不包含主神
	FilterByRarity    int32 `json:"filter_by_rarity"`    // 稀有度过滤
	FilterByAffiliate int32 `json:"filter_by_affiliate"` // 势力过滤
}
