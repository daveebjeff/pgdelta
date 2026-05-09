package schema

// Schema represents a full PostgreSQL schema snapshot.
type Schema struct {
	Tables           []Table
	Columns          []Column
	Indexes          []Index
	Sequences        []Sequence
	Views            []View
	MaterializedViews []MaterializedView
	Functions        []Function
	Triggers         []Trigger
	Enums            []Enum
	Extensions       []Extension
	Policies         []Policy
	Constraints      []Constraint
	ForeignKeys      []ForeignKey
	Roles            []Role
}

// SchemaDiff holds all differences between two schema snapshots.
type SchemaDiff struct {
	TableDiff           TableDiff
	ColumnDiff          ColumnDiff
	IndexDiff           IndexDiff
	SequenceDiff        SequenceDiff
	ViewDiff            ViewDiff
	MaterializedViewDiff MaterializedViewDiff
	FunctionDiff        FunctionDiff
	TriggerDiff         TriggerDiff
	EnumDiff            EnumDiff
	ExtensionDiff       ExtensionDiff
	PolicyDiff          PolicyDiff
	ConstraintDiff      ConstraintDiff
	ForeignKeyDiff      ForeignKeyDiff
	RoleDiff            RoleDiff
}

// DiffSchemas computes the full diff between two schema snapshots.
func DiffSchemas(old, new Schema) SchemaDiff {
	return SchemaDiff{
		TableDiff:            DiffTables(old.Tables, new.Tables),
		ColumnDiff:           DiffColumns(old.Columns, new.Columns),
		IndexDiff:            DiffIndexes(old.Indexes, new.Indexes),
		SequenceDiff:         DiffSequences(old.Sequences, new.Sequences),
		ViewDiff:             DiffViews(old.Views, new.Views),
		MaterializedViewDiff: DiffMaterializedViews(old.MaterializedViews, new.MaterializedViews),
		FunctionDiff:         DiffFunctions(old.Functions, new.Functions),
		TriggerDiff:          DiffTriggers(old.Triggers, new.Triggers),
		EnumDiff:             DiffEnums(old.Enums, new.Enums),
		ExtensionDiff:        DiffExtensions(old.Extensions, new.Extensions),
		PolicyDiff:           DiffPolicies(old.Policies, new.Policies),
		ConstraintDiff:       DiffConstraints(old.Constraints, new.Constraints),
		ForeignKeyDiff:       DiffForeignKeys(old.ForeignKeys, new.ForeignKeys),
		RoleDiff:             DiffRoles(old.Roles, new.Roles),
	}
}

// IsEmpty returns true if the SchemaDiff contains no changes.
func (d SchemaDiff) IsEmpty() bool {
	return len(d.TableDiff.Added) == 0 && len(d.TableDiff.Removed) == 0 &&
		len(d.ColumnDiff.Added) == 0 && len(d.ColumnDiff.Removed) == 0 &&
		len(d.IndexDiff.Added) == 0 && len(d.IndexDiff.Removed) == 0 && len(d.IndexDiff.Changed) == 0 &&
		len(d.SequenceDiff.Added) == 0 && len(d.SequenceDiff.Removed) == 0 && len(d.SequenceDiff.Changed) == 0 &&
		len(d.ViewDiff.Added) == 0 && len(d.ViewDiff.Removed) == 0 && len(d.ViewDiff.Changed) == 0 &&
		len(d.FunctionDiff.Added) == 0 && len(d.FunctionDiff.Removed) == 0 && len(d.FunctionDiff.Changed) == 0 &&
		len(d.EnumDiff.Added) == 0 && len(d.EnumDiff.Removed) == 0 && len(d.EnumDiff.Changed) == 0 &&
		len(d.ExtensionDiff.Added) == 0 && len(d.ExtensionDiff.Removed) == 0 && len(d.ExtensionDiff.Changed) == 0 &&
		len(d.RoleDiff.Added) == 0 && len(d.RoleDiff.Removed) == 0 && len(d.RoleDiff.Changed) == 0
}
