package mapper

import (
	domain "tictactoe/internal/domain/model"
	record "tictactoe/internal/infrastructure/storage/model"
)

func ToRulesRecords(rules *domain.Rules) record.RulesRecord {
	return record.RulesRecord{
		UUID:        rules.UUID,
		BoardWidth:  rules.BoardWidth,
		BoardHeight: rules.BoardHeight,
		WinLength:   rules.WinLength,
	}
}

func ToRulesDomain(record record.RulesRecord) *domain.Rules {
	return &domain.Rules{
		UUID:        record.UUID,
		BoardWidth:  record.BoardWidth,
		BoardHeight: record.BoardHeight,
		WinLength:   record.WinLength,
	}
}
