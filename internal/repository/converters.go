package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func decimalToPgNumeric(d decimal.Decimal) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(d.String())
	return n
}

func pgNumericToDecimal(n pgtype.Numeric) decimal.Decimal {
	if !n.Valid {
		return decimal.Zero
	}

	d := decimal.NewFromBigInt(n.Int, n.Exp)
	return d
}