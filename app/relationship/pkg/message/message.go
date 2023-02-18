package message

const (
	TopicRelationFollow       = "relation_follow"
	TopicRelationCacheRebuild = "relation_cache_rebuild"
	TopicRelationSyncCount    = "relation_sync_count"
)

const (
	TagCreate     = "create"
	TagDelete     = "delete"
	TagSync       = "sync"
	TagListMissed = "list_missed"
)

const (
	ListBatchSize    = 20
	RebuildBatchSize = 100
)
