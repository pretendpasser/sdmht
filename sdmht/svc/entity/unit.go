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
	Location int32 // 位置
	BaseAttribute

	Health        uint32 // 生命值
	Defend        uint32 // 护盾
	Move          int32  // 移动力
	AttackPrevent bool   // 攻击防止(圣盾)
	CounterAttack bool   // 反击
	Rebirth       bool   // 重生
	HurtInstead   bool   // 坚壁
	NoMove        uint32 // 临时禁止移动
	NoAttack      uint32 // 临时禁止攻击
	NoCure        uint32 // 临时禁止治疗

	TempAttact []*TempAttribute // 临时攻击力变化 [数值]持续时间
	TempMove   []*TempAttribute // 临时移动力变化 [数值]持续时间
	Hurt       []*TempAttribute // 创伤层数 [数值]持续时间
	Cure       []*TempAttribute // 治疗层数 [数值]持续时间
}

type TempAttribute struct {
	Value  uint32 // 数值
	Period uint32 // 持续时间
}

// 单位属性
type BaseAttribute struct {
	ID           uint64 `db:"id"`
	Name         string `db:"name"`
	SkillName    string `db:"skill_name"`
	Rarity       int32  `db:"rarity"`    // 稀有度
	Affiliate    int32  `db:"affiliate"` // 所属势力
	Attack       uint32 `db:"attack"`    // 攻击力
	MaxDefend    uint32 `db:"defend"`    // 最大护盾值
	MaxHealth    uint32 `db:"health"`    // 最大生命值
	MaxMove      uint32 `db:"move"`      // 最大移动力
	BaseNoMove   bool   `db:"-"`         // 禁止移动
	BaseNoAttack bool   `db:"-"`         // 禁止攻击
	BaseNoCure   bool   `db:"-"`         // 禁止治疗
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
