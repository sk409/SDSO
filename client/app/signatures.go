package main

type signature interface {
	name() string
	description() string
	diagnostician() diagnostician
}

type xssSignatureAngleBrackets struct {
}

func newXSSSignatureAngleBrackets() *xssSignatureAngleBrackets {
	return &xssSignatureAngleBrackets{}
}

func (x *xssSignatureAngleBrackets) name() string {
	return "XSS"
}

func (x *xssSignatureAngleBrackets) description() string {
	return "<がエスケープされていません。"
}

func (x *xssSignatureAngleBrackets) diagnostician() diagnostician {
	text := "<*******DAST************"
	return newTextMatchDiagnostician([]string{text}, []string{text})
}

type osCommandInjectionSleep struct {
}

func newOSCommandInjectionSleep() *osCommandInjectionSleep {
	return &osCommandInjectionSleep{}
}

func (o *osCommandInjectionSleep) name() string {
	return "OSコマンドインジェクション"
}

func (o *osCommandInjectionSleep) description() string {
	return "受け取ったコマンドをそのまま実行しています。"
}

func (o *osCommandInjectionSleep) diagnostician() diagnostician {
	return newTimeDiagnostician([]string{"sleep 3"}, 3)
}

type sqlInjectionSleep struct {
}

func newSQLInjectionSleep() *sqlInjectionSleep {
	return &sqlInjectionSleep{}
}

func (s *sqlInjectionSleep) name() string {
	return "SQLインジェクション"
}

func (s *sqlInjectionSleep) description() string {
	return "受け取った値をそのままクエリに使用しています。"
}

func (s *sqlInjectionSleep) diagnostician() diagnostician {
	// TODO: MySQL以外にも対応
	return newTimeDiagnostician([]string{"select sleep(3);", "'';select sleep(3);", "');select sleep(3);#"}, 3)
}
