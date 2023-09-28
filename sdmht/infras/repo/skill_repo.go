package repo

import (
	"sdmht/lib/log"
	"sdmht/sdmht/svc/entity"
	"sync/atomic"
)

// 触发式
// 开局初始化技能及其触发event
// checking skill

type SkillRepo struct {
	skill     map[string]entity.SkillHandler
	skillDesc map[string]string
}

func NewSkillRepo() *SkillRepo {
	repo := &SkillRepo{}
	repo.skill, repo.skillDesc = skillInit()
	return repo
}

func (r *SkillRepo) Get(name string) string {
	return r.skillDesc[name]
}

func (r *SkillRepo) Find() map[string]string {
	return r.skillDesc
}

func (r *SkillRepo) Checking(name string) entity.SkillHandler {
	return r.skill[name]
}

func skillInit() (map[string]entity.SkillHandler, map[string]string) {
	skill := make(map[string]entity.SkillHandler)
	skillDesc := make(map[string]string)
	skill["德罗普尼尔"] = SkillDeLuoPuNier
	skillDesc["德罗普尼尔"] = SkillDescDeLuoPuNier()
	skill["机略"] = SkillJiLue
	skillDesc["机略"] = SkillDescJiLue()
	skill["悖逆"] = SkillBeiNi
	skillDesc["悖逆"] = SkillDescBeiNi()
	skill["要塞广寒"] = SkillYaoSaiGuangHan
	skillDesc["要塞广寒"] = SkillDescYaoSaiGuangHan()

	return skill, skillDesc
}

// id == 0, means get all Subsidiary Deitys, return length is less or equal than 2;
// id != 0, means get the other Subsidiary Deity, return length is less or equal than 1;
func getSubsidiaryDeityID(units []*entity.Unit, id int32) []int32 {
	unitID := []int32{}
	for _, unit := range units {
		if unit.Rarity >= 1 && unit.Rarity <= 4 {
			if unit.UnitID == id {
				continue
			}
			unitID = append(unitID, unit.UnitID)
		}
	}
	return unitID
}

func SkillDescDeLuoPuNier() string {
	return `获得两张弹幕卡【永恒之枪】（每回合至多使用1次）。`
}
func SkillDeLuoPuNier(m *entity.Match, unitID int32, _ entity.SkillChecking) {
	// unit := m.Players[m.WhoseTurn].Units[unitID]

}

func SkillDescJiLue() string {
	return `每获得一张卡时，随机解除1格迷雾（每回合至多发动4次）。`
}
func SkillJiLue(m *entity.Match, unitID int32, checking entity.SkillChecking) {
	skill, ok := m.Players[m.WhoseTurn].Units[unitID].Skills["机略"]
	if !ok {
		log.S().Errorw("机略", "system error", "not found skill name")
	}

	// 自己回合开始重置
	if checking.NextRound && checking.SelfEvent && skill.RoundUseCnt != 0 {
		skill.RoundUseCnt = 0
	}

	// continue: 自己回合 且 手牌变动
	if !checking.SelfEvent || checking.HandCardNumChange == 0 {
		log.S().Infow("机略", "SelfEvent", checking.SelfEvent,
			"HandCardNumChange", checking.HandCardNumChange)
		return
	}

	if skill.RoundUseCnt >= 4 {
		log.S().Info("机略 ", "Max Round Use")
		return
	}

	changeNum := checking.HandCardNumChange
	if changeNum+skill.RoundUseCnt > 4 {
		changeNum = 4 - skill.RoundUseCnt
	}
	atomic.AddInt32(&skill.RoundUseCnt, changeNum)
	if changeNum == 0 {
		return
	}

	log.S().Infow("机略", "HandCardNumChange", checking.HandCardNumChange,
		"useNum", changeNum, "totalRoundUse", skill.RoundUseCnt)
	m.Players[m.GetOtherPlayer()].RandomChangeSquare(int(changeNum), true)

}

func SkillDescBeiNi() string {
	return `另一附属神永远处于迷雾（仍然会受到伤害）。己方回合开始时，受此效果影响的单位获得【圣盾】，永久+1攻击力。`
}
func SkillBeiNi(m *entity.Match, unitID int32, checking entity.SkillChecking) {
	if !(checking.SelfEvent && checking.NextRound) && !checking.WillBeDeath {
		return
	}

	subsidiaryDeityIDs := getSubsidiaryDeityID(m.Players[m.WhoseTurn].Units, unitID)
	if len(subsidiaryDeityIDs) == 0 {
		log.S().Info("悖逆 ", "not get other subsidiary deity")
		return
	}
	otherUnit := m.Players[m.WhoseTurn].Units[subsidiaryDeityIDs[0]]
	if checking.WillBeDeath {
		// atomic.StoreInt32(&otherUnit.PermanentlyCover, 0)
		log.S().Info("悖逆 ", "WillBeDeath")
		return
	}
	log.S().Infow("悖逆", "next round", m.CurRound)
	otherUnit.AttackPrevent = true
	// atomic.StoreInt32(&otherUnit.PermanentlyCover, 1)
	atomic.AddUint32(&otherUnit.Attack, 1)
}

func SkillDescYaoSaiGuangHan() string {
	return `【反击】，自身失去生命值时，随机移动至迷雾区域（无迷雾时原地不动）。装填弹幕卡时，永久+1攻击力，己方随机覆盖1格迷雾（每回合至多发动3次）。`
}
func SkillYaoSaiGuangHan(m *entity.Match, unitID int32, _ entity.SkillChecking) {
	// unit := m.Players[m.WhoseTurn].Units[unitID]
	// if !unit.CounterAttack {
	// 	AttackCounter(unit)
	// }

}

// 【反击】
func AttackCounter(unit *entity.Unit) {
	// unit.CounterAttack = true
}
