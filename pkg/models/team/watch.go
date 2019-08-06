package team

// WatchStore interface of watch 关注
type WatchStore interface {
	Gets(uid string) []string
	Watch(uid, target string) error
	Unwatch(uid, target string) error
}
