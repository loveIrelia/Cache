package go_cache

//传入key选择相应的节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//从对应group查找缓存值（Http客户端）
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
