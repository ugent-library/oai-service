// Code generated by ent, DO NOT EDIT.

package item

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/ugent-library/oai-service/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.Item {
	return predicate.Item(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Item {
	return predicate.Item(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Item {
	return predicate.Item(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Item {
	return predicate.Item(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.Item {
	return predicate.Item(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.Item {
	return predicate.Item(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Item {
	return predicate.Item(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Item {
	return predicate.Item(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Item {
	return predicate.Item(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.Item {
	return predicate.Item(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.Item {
	return predicate.Item(sql.FieldContainsFold(FieldID, id))
}

// HasRecords applies the HasEdge predicate on the "records" edge.
func HasRecords() predicate.Item {
	return predicate.Item(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, RecordsTable, RecordsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRecordsWith applies the HasEdge predicate on the "records" edge with a given conditions (other predicates).
func HasRecordsWith(preds ...predicate.Record) predicate.Item {
	return predicate.Item(func(s *sql.Selector) {
		step := newRecordsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasSets applies the HasEdge predicate on the "sets" edge.
func HasSets() predicate.Item {
	return predicate.Item(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, SetsTable, SetsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSetsWith applies the HasEdge predicate on the "sets" edge with a given conditions (other predicates).
func HasSetsWith(preds ...predicate.Set) predicate.Item {
	return predicate.Item(func(s *sql.Selector) {
		step := newSetsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Item) predicate.Item {
	return predicate.Item(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Item) predicate.Item {
	return predicate.Item(func(s *sql.Selector) {
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
func Not(p predicate.Item) predicate.Item {
	return predicate.Item(func(s *sql.Selector) {
		p(s.Not())
	})
}