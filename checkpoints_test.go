package kugo

import (
	"context"
	"testing"

	"github.com/tj/assert"
)

func Test_Checkpoints(t *testing.T) {
	t.SkipNow()
	c := New(WithEndpoint("http://localhost:1442"))
	point, err := c.CheckpointBySlot(context.Background(), CheckpointBySlotInput{SlotNo: 51540727})
	assert.Nil(t, err)

	assert.Equal(t, "fe5f9af58ab0511a77524f4d2a0b930213b3bb1353e11e3d69e83129b9fbe65a", point.HeaderHash)
	assert.Equal(t, 51540722, point.SlotNo)
}
