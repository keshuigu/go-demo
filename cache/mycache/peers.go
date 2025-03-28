package mycache

import pb "cache/mycache/mycachepb"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 查找缓存值
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
