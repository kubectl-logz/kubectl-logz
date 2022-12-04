package internal

import (
	"testing"

	"github.com/kubectl-logz/kubectl-logz/internal/types"
	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) {
	e := parse([]byte(`[2022-12-04T11:34:26,673-0800]-[WARN ]-["pool-2-thread-1"  cid=, clu=]-[o.a.k.c.c.i.SubscriptionState]-[399]-[Consumer clientId=consumer-Intuit.asset.alias.undefined-local-1, groupId=Intuit.asset.alias.undefined-local] Resetting offset for partition sequencer-local-0 to position FetchPosition{offset=0, offsetEpoch=Optional.empty, currentLeader=LeaderAndEpoch{leader=Optional[localhost:9092 (id: 1 rack: null)], epoch=0}}.`))

	assert.True(t, e.Valid(), "valid")
	assert.False(t, e.Time.IsZero(), "non-zero time")
	assert.Equal(t, types.Level("warn"), e.Level, "warn")
	assert.Equal(t, "pool-2-thread-1", e.ThreadID, "threadId")
	assert.Equal(t, "o.a.k.c.c.i.SubscriptionState 399 Consumer clientId=consumer-Intuit.asset.alias.undefined-local-1, groupId=Intuit.asset.alias.undefined-local localhost:9092 (id: 1 rack: null) Resetting offset for partition sequencer-local-0 to position FetchPosition{offset=0, offsetEpoch=Optional.empty, currentLeader=LeaderAndEpoch{leader=Optional, epoch=0}}.", e.Msg)
}
