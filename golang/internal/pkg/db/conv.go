package db

import "database/sql"

func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func FromNullString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func ToNullInt64(i *int) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: int64(*i), Valid: true}
	}
	return sql.NullInt64{Valid: false}
}

func FromNullInt64(ni sql.NullInt64) *int {
	if ni.Valid {
		v := int(ni.Int64)
		return &v
	}
	return nil
}

func ToNullBool(b *bool) sql.NullBool {
	if b != nil {
		return sql.NullBool{Bool: *b, Valid: true}
	}
	return sql.NullBool{Valid: false}
}

func FromNullBool(nb sql.NullBool) *bool {
	if nb.Valid {
		return &nb.Bool
	}
	return nil
}

func ToNullFloat64(f *float64) sql.NullFloat64 {
	if f != nil {
		return sql.NullFloat64{Float64: *f, Valid: true}
	}
	return sql.NullFloat64{Valid: false}
}

func FromNullFloat64(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

func PatchField[T any](old T, newVal *T) T {
	if newVal == nil {
		return old
	}
	return *newVal
}
