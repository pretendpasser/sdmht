package entity

const (
	OpEventAttack   string = "attack"       // from(unitid)    to(position)
	OpEventMove     string = "move"         // from(unitid)    to(position)
	OpEventCard     string = "card"         // from(cardid)    to(unitid)
	OpEventSkill    string = "active_skill" // from(skillname) to(empty)
	OpEventEndRound string = "end_round"    // from(playerid)  to(empty)

	EventHurt string = "hurt" // from(unitid)    to(hurt num)
)

type SyncOperator struct {
	MatchID  uint64
	PlayerID int32
	Event    string
	From     string
	To       string
}

type SkillChecking struct {
	SelfEvent         bool  // 当前为己方事件
	NextRound         bool  // 回合切换
	BeAttack          bool  // 被攻击
	WillBeDeath       bool  // 将死亡
	HealthChange      int32 // 生命值变化
	HandCardNumChange int32 // 手牌数变化
	DeadUnit          int32 // 死亡单位ID
}
