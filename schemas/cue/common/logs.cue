package common

LogsSortOrder:     "Descending" | "Ascending"                 @cuetsy(kind="enum")
LogsDedupStrategy: "none" | "exact" | "numbers" | "signature" @cog(kind="enum",memberNames="none|exact|numbers|signature")
