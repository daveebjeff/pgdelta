package migrate_test

import (
	"strings"
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

var baseOp = schema.Operator{
	Schema:     "public",
	Name:       "+",
	LeftType:   "integer",
	RightType:  "integer",
	ResultType: "integer",
	Procedure:  "int4pl",
}

func TestCreateOperatorSQL_Basic(t *testing.T) {
	sql := migrate.CreateOperatorSQL(baseOp)
	if !strings.Contains(sql, "CREATE OPERATOR public.+") {
		t.Errorf("expected CREATE OPERATOR public.+, got: %s", sql)
	}
	if !strings.Contains(sql, "PROCEDURE = int4pl") {
		t.Errorf("expected PROCEDURE = int4pl, got: %s", sql)
	}
	if !strings.Contains(sql, "LEFTARG = integer") {
		t.Errorf("expected LEFTARG = integer, got: %s", sql)
	}
}

func TestCreateOperatorSQL_WithCommutator(t *testing.T) {
	op := baseOp
	op.Commutator = "OPERATOR(public.+)"
	sql := migrate.CreateOperatorSQL(op)
	if !strings.Contains(sql, "COMMUTATOR = OPERATOR(public.+)") {
		t.Errorf("expected COMMUTATOR clause, got: %s", sql)
	}
}

func TestDropOperatorSQL(t *testing.T) {
	sql := migrate.DropOperatorSQL(baseOp)
	expected := "DROP OPERATOR public.+ (integer, integer);"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
}

func TestOperatorDiffSQL_AddedAndRemoved(t *testing.T) {
	diff := schema.OperatorDiff{
		Added:   []schema.Operator{baseOp},
		Removed: []schema.Operator{baseOp},
	}
	stmts := migrate.OperatorDiffSQL(diff)
	if len(stmts) != 2 {
		t.Errorf("expected 2 statements, got %d", len(stmts))
	}
}

func TestOperatorDiffSQL_Changed(t *testing.T) {
	changed := baseOp
	changed.Procedure = "custom_add"
	diff := schema.OperatorDiff{
		Changed: []schema.Operator{changed},
	}
	stmts := migrate.OperatorDiffSQL(diff)
	if len(stmts) != 2 {
		t.Errorf("expected DROP + CREATE for changed operator, got %d stmts", len(stmts))
	}
	if !strings.HasPrefix(stmts[0], "DROP") {
		t.Errorf("expected DROP first, got: %s", stmts[0])
	}
	if !strings.HasPrefix(stmts[1], "CREATE") {
		t.Errorf("expected CREATE second, got: %s", stmts[1])
	}
}
