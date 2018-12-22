{{- $short := (shortname .Type.Name) -}}
// {{ .Name }} returns the {{ .RefType.Name }} associated with the {{ .Type.Name }}'s {{ .Field.Name }} ({{ .Field.Col.ColumnName }}).
//
// Generated from foreign key '{{ .ForeignKey.ForeignKeyName }}'.
func ({{ $short }} *{{ .Type.Name }}) {{ .Name }}By{{ .Field.Name }}(db XODB) (*{{ .RefType.Name }}, error) {
	var err error
    {{ $short }}.{{ .RefType.Name }}, err={{ .RefType.Name }}By{{ .RefField.Name }}(db, {{ convext $short .Field .RefField }})
    return {{ $short }}.{{ .RefType.Name }}, err
}

