{{- $short := (shortname .Type.Name) -}}
// {{ .Name }} returns the {{ .RefType.Name }} associated with the {{ .Type.Name }}'s {{ .Field.Name }} ({{ .Field.Col.ColumnName }}).
//
// Generated from foreign key '{{ .ForeignKey.ForeignKeyName }}'.
func ({{ $short }} *{{ .Type.Name }}) Load{{ .Name }}(db XODB)  error {
	var err error
    {{ $short }}.{{ .RefType.Name }}, err={{ .RefType.Name }}By{{ .RefField.Name }}(db, {{ convext $short .Field .RefField }})
    return  err
}

