// Code generated by ent, DO NOT EDIT.

package set

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/ugent-library/oai-service/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int64) predicate.Set {
	return predicate.Set(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int64) predicate.Set {
	return predicate.Set(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int64) predicate.Set {
	return predicate.Set(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int64) predicate.Set {
	return predicate.Set(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int64) predicate.Set {
	return predicate.Set(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int64) predicate.Set {
	return predicate.Set(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int64) predicate.Set {
	return predicate.Set(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int64) predicate.Set {
	return predicate.Set(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int64) predicate.Set {
	return predicate.Set(sql.FieldLTE(FieldID, id))
}

// SetSpec applies equality check predicate on the "set_spec" field. It's identical to SetSpecEQ.
func SetSpec(v string) predicate.Set {
	return predicate.Set(sql.FieldEQ(FieldSetSpec, v))
}

// SetName applies equality check predicate on the "set_name" field. It's identical to SetNameEQ.
func SetName(v string) predicate.Set {
	return predicate.Set(sql.FieldEQ(FieldSetName, v))
}

// SetDescription applies equality check predicate on the "set_description" field. It's identical to SetDescriptionEQ.
func SetDescription(v string) predicate.Set {
	return predicate.Set(sql.FieldEQ(FieldSetDescription, v))
}

// SetSpecEQ applies the EQ predicate on the "set_spec" field.
func SetSpecEQ(v string) predicate.Set {
	return predicate.Set(sql.FieldEQ(FieldSetSpec, v))
}

// SetSpecNEQ applies the NEQ predicate on the "set_spec" field.
func SetSpecNEQ(v string) predicate.Set {
	return predicate.Set(sql.FieldNEQ(FieldSetSpec, v))
}

// SetSpecIn applies the In predicate on the "set_spec" field.
func SetSpecIn(vs ...string) predicate.Set {
	return predicate.Set(sql.FieldIn(FieldSetSpec, vs...))
}

// SetSpecNotIn applies the NotIn predicate on the "set_spec" field.
func SetSpecNotIn(vs ...string) predicate.Set {
	return predicate.Set(sql.FieldNotIn(FieldSetSpec, vs...))
}

// SetSpecGT applies the GT predicate on the "set_spec" field.
func SetSpecGT(v string) predicate.Set {
	return predicate.Set(sql.FieldGT(FieldSetSpec, v))
}

// SetSpecGTE applies the GTE predicate on the "set_spec" field.
func SetSpecGTE(v string) predicate.Set {
	return predicate.Set(sql.FieldGTE(FieldSetSpec, v))
}

// SetSpecLT applies the LT predicate on the "set_spec" field.
func SetSpecLT(v string) predicate.Set {
	return predicate.Set(sql.FieldLT(FieldSetSpec, v))
}

// SetSpecLTE applies the LTE predicate on the "set_spec" field.
func SetSpecLTE(v string) predicate.Set {
	return predicate.Set(sql.FieldLTE(FieldSetSpec, v))
}

// SetSpecContains applies the Contains predicate on the "set_spec" field.
func SetSpecContains(v string) predicate.Set {
	return predicate.Set(sql.FieldContains(FieldSetSpec, v))
}

// SetSpecHasPrefix applies the HasPrefix predicate on the "set_spec" field.
func SetSpecHasPrefix(v string) predicate.Set {
	return predicate.Set(sql.FieldHasPrefix(FieldSetSpec, v))
}

// SetSpecHasSuffix applies the HasSuffix predicate on the "set_spec" field.
func SetSpecHasSuffix(v string) predicate.Set {
	return predicate.Set(sql.FieldHasSuffix(FieldSetSpec, v))
}

// SetSpecEqualFold applies the EqualFold predicate on the "set_spec" field.
func SetSpecEqualFold(v string) predicate.Set {
	return predicate.Set(sql.FieldEqualFold(FieldSetSpec, v))
}

// SetSpecContainsFold applies the ContainsFold predicate on the "set_spec" field.
func SetSpecContainsFold(v string) predicate.Set {
	return predicate.Set(sql.FieldContainsFold(FieldSetSpec, v))
}

// SetNameEQ applies the EQ predicate on the "set_name" field.
func SetNameEQ(v string) predicate.Set {
	return predicate.Set(sql.FieldEQ(FieldSetName, v))
}

// SetNameNEQ applies the NEQ predicate on the "set_name" field.
func SetNameNEQ(v string) predicate.Set {
	return predicate.Set(sql.FieldNEQ(FieldSetName, v))
}

// SetNameIn applies the In predicate on the "set_name" field.
func SetNameIn(vs ...string) predicate.Set {
	return predicate.Set(sql.FieldIn(FieldSetName, vs...))
}

// SetNameNotIn applies the NotIn predicate on the "set_name" field.
func SetNameNotIn(vs ...string) predicate.Set {
	return predicate.Set(sql.FieldNotIn(FieldSetName, vs...))
}

// SetNameGT applies the GT predicate on the "set_name" field.
func SetNameGT(v string) predicate.Set {
	return predicate.Set(sql.FieldGT(FieldSetName, v))
}

// SetNameGTE applies the GTE predicate on the "set_name" field.
func SetNameGTE(v string) predicate.Set {
	return predicate.Set(sql.FieldGTE(FieldSetName, v))
}

// SetNameLT applies the LT predicate on the "set_name" field.
func SetNameLT(v string) predicate.Set {
	return predicate.Set(sql.FieldLT(FieldSetName, v))
}

// SetNameLTE applies the LTE predicate on the "set_name" field.
func SetNameLTE(v string) predicate.Set {
	return predicate.Set(sql.FieldLTE(FieldSetName, v))
}

// SetNameContains applies the Contains predicate on the "set_name" field.
func SetNameContains(v string) predicate.Set {
	return predicate.Set(sql.FieldContains(FieldSetName, v))
}

// SetNameHasPrefix applies the HasPrefix predicate on the "set_name" field.
func SetNameHasPrefix(v string) predicate.Set {
	return predicate.Set(sql.FieldHasPrefix(FieldSetName, v))
}

// SetNameHasSuffix applies the HasSuffix predicate on the "set_name" field.
func SetNameHasSuffix(v string) predicate.Set {
	return predicate.Set(sql.FieldHasSuffix(FieldSetName, v))
}

// SetNameEqualFold applies the EqualFold predicate on the "set_name" field.
func SetNameEqualFold(v string) predicate.Set {
	return predicate.Set(sql.FieldEqualFold(FieldSetName, v))
}

// SetNameContainsFold applies the ContainsFold predicate on the "set_name" field.
func SetNameContainsFold(v string) predicate.Set {
	return predicate.Set(sql.FieldContainsFold(FieldSetName, v))
}

// SetDescriptionEQ applies the EQ predicate on the "set_description" field.
func SetDescriptionEQ(v string) predicate.Set {
	return predicate.Set(sql.FieldEQ(FieldSetDescription, v))
}

// SetDescriptionNEQ applies the NEQ predicate on the "set_description" field.
func SetDescriptionNEQ(v string) predicate.Set {
	return predicate.Set(sql.FieldNEQ(FieldSetDescription, v))
}

// SetDescriptionIn applies the In predicate on the "set_description" field.
func SetDescriptionIn(vs ...string) predicate.Set {
	return predicate.Set(sql.FieldIn(FieldSetDescription, vs...))
}

// SetDescriptionNotIn applies the NotIn predicate on the "set_description" field.
func SetDescriptionNotIn(vs ...string) predicate.Set {
	return predicate.Set(sql.FieldNotIn(FieldSetDescription, vs...))
}

// SetDescriptionGT applies the GT predicate on the "set_description" field.
func SetDescriptionGT(v string) predicate.Set {
	return predicate.Set(sql.FieldGT(FieldSetDescription, v))
}

// SetDescriptionGTE applies the GTE predicate on the "set_description" field.
func SetDescriptionGTE(v string) predicate.Set {
	return predicate.Set(sql.FieldGTE(FieldSetDescription, v))
}

// SetDescriptionLT applies the LT predicate on the "set_description" field.
func SetDescriptionLT(v string) predicate.Set {
	return predicate.Set(sql.FieldLT(FieldSetDescription, v))
}

// SetDescriptionLTE applies the LTE predicate on the "set_description" field.
func SetDescriptionLTE(v string) predicate.Set {
	return predicate.Set(sql.FieldLTE(FieldSetDescription, v))
}

// SetDescriptionContains applies the Contains predicate on the "set_description" field.
func SetDescriptionContains(v string) predicate.Set {
	return predicate.Set(sql.FieldContains(FieldSetDescription, v))
}

// SetDescriptionHasPrefix applies the HasPrefix predicate on the "set_description" field.
func SetDescriptionHasPrefix(v string) predicate.Set {
	return predicate.Set(sql.FieldHasPrefix(FieldSetDescription, v))
}

// SetDescriptionHasSuffix applies the HasSuffix predicate on the "set_description" field.
func SetDescriptionHasSuffix(v string) predicate.Set {
	return predicate.Set(sql.FieldHasSuffix(FieldSetDescription, v))
}

// SetDescriptionIsNil applies the IsNil predicate on the "set_description" field.
func SetDescriptionIsNil() predicate.Set {
	return predicate.Set(sql.FieldIsNull(FieldSetDescription))
}

// SetDescriptionNotNil applies the NotNil predicate on the "set_description" field.
func SetDescriptionNotNil() predicate.Set {
	return predicate.Set(sql.FieldNotNull(FieldSetDescription))
}

// SetDescriptionEqualFold applies the EqualFold predicate on the "set_description" field.
func SetDescriptionEqualFold(v string) predicate.Set {
	return predicate.Set(sql.FieldEqualFold(FieldSetDescription, v))
}

// SetDescriptionContainsFold applies the ContainsFold predicate on the "set_description" field.
func SetDescriptionContainsFold(v string) predicate.Set {
	return predicate.Set(sql.FieldContainsFold(FieldSetDescription, v))
}

// HasRecords applies the HasEdge predicate on the "records" edge.
func HasRecords() predicate.Set {
	return predicate.Set(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, RecordsTable, RecordsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRecordsWith applies the HasEdge predicate on the "records" edge with a given conditions (other predicates).
func HasRecordsWith(preds ...predicate.Record) predicate.Set {
	return predicate.Set(func(s *sql.Selector) {
		step := newRecordsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Set) predicate.Set {
	return predicate.Set(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Set) predicate.Set {
	return predicate.Set(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Set) predicate.Set {
	return predicate.Set(func(s *sql.Selector) {
		p(s.Not())
	})
}
