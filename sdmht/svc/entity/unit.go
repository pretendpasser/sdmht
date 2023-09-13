package entity

import (
	"sdmht/lib/utils"
)

const (
	// MainDeity       = 1 //主神
	// AffiliatedDeity = 2 //附属神明
	Derivative = 3 //衍生物

	AffiliateNanhua       = 0 // 南华庄势力
	AffiliateAsgard       = 1 // 阿斯加德势力
	AffiliateTakamagahara = 2 // 高天原势力
	AffiliateNilfheim     = 3 // 尼夫尔海姆势力
	AffiliateTaiqing      = 4 // 太清殿势力
	AffiliateAmaravati    = 5 // 阿摩婆罗提势力
	AffiliateHeliopolis   = 6 // 赫利奥波利斯势力
	AffiliateOlympus      = 7 // 奥林匹斯势力
	AffiliateYuqing       = 8 // 玉清殿势力

	RarityMainDeity = 0 // 主神
	RarityOrdinary  = 1 // 普通
	RarityDare      = 2 // 稀有
	RarityLegend    = 3 // 传说
	RarityObsidian  = 4 // 黑曜
)

type Unit struct {
	Location int32 `json:"location"` // 位置
	BaseAttribute

	Health        uint32 `json:"health"`         // 生命值
	Defend        uint32 `json:"defend"`         // 护盾
	Move          int32  `json:"move"`           // 移动力
	AttackPrevent bool   `json:"attack_prevent"` // 攻击防止(圣盾)
	CounterAttack bool   `json:"counter_attack"` // 反击
	Rebirth       bool   `json:"rebirth"`        // 重生
	HurtInstead   bool   `json:"hurt_instead"`   // 坚壁
	NoMove        uint32 `json:"no_move"`        // 临时禁止移动
	NoAttack      uint32 `json:"no_attack"`      // 临时禁止攻击
	NoCure        uint32 `json:"no_cure"`        // 临时禁止治疗

	ChangeAttact []*TempAttribute `json:"change_attack"` // 临时攻击力变化 [数值]持续时间
	ChangeMove   []*TempAttribute `json:"change_move"`   // 临时移动力变化 [数值]持续时间
	Hurt         []*TempAttribute `json:"hurt"`          // 创伤层数 [数值]持续时间
	Cure         []*TempAttribute `json:"cure"`          // 治疗层数 [数值]持续时间
}

type TempAttribute struct {
	Value  int32  `json:"value"`  // 数值
	Period uint32 `json:"period"` // 持续时间
}

// 单位属性
type BaseAttribute struct {
	ID           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	SkillName    string `json:"skill_name" db:"skill_name"`
	Rarity       int32  `json:"rarity" db:"rarity"`       // 稀有度
	Affiliate    int32  `json:"affiliate" db:"affiliate"` // 所属势力
	Attack       uint32 `json:"attack" db:"attack"`       // 攻击力
	MaxDefend    uint32 `json:"max_defend" db:"defend"`   // 最大护盾值
	MaxHealth    uint32 `json:"max_health" db:"health"`   // 最大生命值
	MaxMove      uint32 `json:"max_move" db:"move"`       // 最大移动力
	BaseNoMove   bool   `json:"base_no_move" db:"-"`      // 禁止移动
	BaseNoAttack bool   `json:"base_no_attack" db:"-"`    // 禁止攻击
	BaseNoCure   bool   `json:"base_no_cure" db:"-"`      // 禁止治疗
	BaseNoEquip  bool   `json:"base_no_equip" db:"-"`     // 禁止装备
}

func (a *BaseAttribute) CheckMoveing(location int32) error {
	return nil
}

// func (a *BaseAttribute) Moveing(direction string) error {
// 	tempMoves := 0
// 	for _, tempMove := range a.TempMove {
// 		tempMoves += tempMove
// 	}
// 	if a.NoMove > 0 || a.Move+tempMoves <= 0 {
// 		return lib.NewError(lib.Errint32ernal, "")
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

type UnitQuery struct {
	Pagination *utils.Pagination `json:"page"`

	ExcludeMainDeity  bool  `json:"exclude_main_deity"`  // 不包含主神
	FilterByRarity    int32 `json:"filter_by_rarity"`    // 稀有度过滤
	FilterByAffiliate int32 `json:"filter_by_affiliate"` // 势力过滤
}
