package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

func baseEventTrig() schema.EventTrigger {
	return schema.EventTrigger{
		Name:     "audit_ddl",
		Event:    "ddl_command_end",
		FuncName: "public.audit_ddl_func",
		Enabled:  "ENABLE",
		Tags:     []string{"CREATE TABLE", "DROP TABLE"},
	}
}

func TestCreateEventTriggerSQL_WithTags(t *testing.T) {
	et := baseEventTrig()
	sql := CreateEventTriggerSQL(et)
	assert.Equal(t, "CREATE EVENT TRIGGER audit_ddl ON ddl_command_end WHEN TAG IN ('CREATE TABLE', 'DROP TABLE') EXECUTE FUNCTION public.audit_ddl_func();", sql)
}

func TestCreateEventTriggerSQL_NoTags(t *testing.T) {
	et := baseEventTrig()
	et.Tags = nil
	sql := CreateEventTriggerSQL(et)
	assert.Equal(t, "CREATE EVENT TRIGGER audit_ddl ON ddl_command_end EXECUTE FUNCTION public.audit_ddl_func();", sql)
}

func TestDropEventTriggerSQL(t *testing.T) {
	et := baseEventTrig()
	sql := DropEventTriggerSQL(et)
	assert.Equal(t, "DROP EVENT TRIGGER audit_ddl;", sql)
}

func TestAlterEventTriggerSQL_Disable(t *testing.T) {
	et := baseEventTrig()
	et.Enabled = "DISABLE"
	sql := AlterEventTriggerSQL(et)
	assert.Equal(t, "ALTER EVENT TRIGGER audit_ddl DISABLE;", sql)
}

func TestEventTriggerDiffSQL_AddedAndRemoved(t *testing.T) {
	et := baseEventTrig()
	stmts := EventTriggerDiffSQL([]schema.EventTrigger{et}, nil)
	assert.Len(t, stmts, 1)
	assert.Equal(t, "DROP EVENT TRIGGER audit_ddl;", stmts[0])

	stmts = EventTriggerDiffSQL(nil, []schema.EventTrigger{et})
	assert.Len(t, stmts, 1)
	assert.Contains(t, stmts[0], "CREATE EVENT TRIGGER")
}

func TestEventTriggerDiffSQL_Changed(t *testing.T) {
	old := baseEventTrig()
	new := baseEventTrig()
	new.FuncName = "public.new_audit_func"
	stmts := EventTriggerDiffSQL([]schema.EventTrigger{old}, []schema.EventTrigger{new})
	assert.Len(t, stmts, 2)
	assert.Equal(t, "DROP EVENT TRIGGER audit_ddl;", stmts[0])
	assert.Contains(t, stmts[1], "public.new_audit_func")
}
