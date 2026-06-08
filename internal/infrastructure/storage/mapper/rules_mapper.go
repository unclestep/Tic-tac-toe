package mapper

import (
	dmodel "tictactoe/internal/domain/model"
	dsmodel "tictactoe/internal/infrastructure/storage/model"
)

func ToRulesStorage(rules *dmodel.Rules) *dsmodel.RulesRecord {
	return &dsmodel.RulesRecord{
		UUID:        rules.UUID,
		BoardWidth:  rules.BoardWidth,
		BoardHeight: rules.BoardHeight,
		WinLength:   rules.WinLength,
	}
}

func ToRulesDomain(record *dsmodel.RulesRecord) *dmodel.Rules {
	return &dmodel.Rules{
		UUID:        record.UUID,
		BoardWidth:  record.BoardWidth,
		BoardHeight: record.BoardHeight,
		WinLength:   record.WinLength,
	}
}
