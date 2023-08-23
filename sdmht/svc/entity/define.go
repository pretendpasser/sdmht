package entity

// type MatchRoom struct {
// 	ID      uint64     `json:"id"`
// 	Players [2]*Player `json:"players"`
// }

type Player struct {
	ID    uint64 `json:"id"`
	Scene *Scene `json:"scene"`
}
