package main

import (
	"github.com/pborman/uuid"
)

func transformSection(t term) section {
	return section{
		UUID:          uuid.NewMD5(uuid.UUID{}, []byte(t.ID)).String(),
		CanonicalName: t.CanonicalName,
		TmeIdentifier: t.ID,
		Type:          "Section",
	}
}
