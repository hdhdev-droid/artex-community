package server

// Messages shared by the migrated assessment and defense workflows. Evidence
// excerpts and historical snapshots remain verbatim.
func init() {
	for ko, translations := range map[string][4]string{
		"근거 해시가 없습니다: %s":        {"Evidence hash is missing: %s", "缺少证据哈希：%s", "Falta el hash de las pruebas: %s", "証拠ハッシュがありません: %s"},
		"근거 무결성을 확인할 수 없습니다: %s": {"Cannot verify evidence integrity: %s", "无法验证证据完整性：%s", "No se puede verificar la integridad de las pruebas: %s", "証拠の整合性を確認できません: %s"},
		"메인 에이전트":                {"Main agent", "主代理", "Agente principal", "メインエージェント"},
		"오류: ":                   {"error: ", "错误：", "error: ", "エラー: "},

		"Closed-loop 입력이 올바르지 않습니다": {
			"Invalid closed-loop input", "闭环输入无效", "Entrada de evaluación iterativa inválida", "Closed-loop 入力が無効です",
		},
		"방어 규칙 평가 입력이 올바르지 않습니다": {
			"Invalid defense rule evaluation input", "防御规则评估输入无效", "Entrada de evaluación de reglas defensivas inválida", "防御規則評価の入力が無効です",
		},
		"보안 로그 연결 설정이 올바르지 않습니다": {
			"Invalid security log connection settings", "安全日志连接设置无效", "Configuración de registros de seguridad inválida", "セキュリティログ接続設定が無効です",
		},
		"Splunk 관리 API 연결과 인증을 확인했습니다. 실제 로그·경보 조회는 진단 실행별로 검증합니다": {
			"Splunk management API connectivity and authentication verified. Actual log and alert queries are verified per assessment.",
			"已验证Splunk管理API连接和认证。实际日志及告警查询按每次诊断验证。",
			"Conexión y autenticación de la API de Splunk verificadas. Las consultas de registros y alertas se verifican por evaluación.",
			"Splunk管理APIの接続と認証を確認しました。実際のログ・アラート照会は診断ごとに検証します。",
		},
		"저장된 인증 정보를 읽지 못했습니다. 다시 입력해 주세요": {
			"Could not read saved credentials. Enter them again.", "无法读取保存的认证信息，请重新输入。",
			"No se pudieron leer las credenciales guardadas. Introdúzcalas de nuevo.", "保存した認証情報を読み取れません。再入力してください。",
		},
		"저장된 인증 정보를 복호화하지 못했습니다. 다시 입력해 주세요": {
			"Could not decrypt saved credentials. Enter them again.", "无法解密保存的认证信息，请重新输入。",
			"No se pudieron descifrar las credenciales guardadas. Introdúzcalas de nuevo.", "保存した認証情報を復号できません。再入力してください。",
		},
		"연결 인증 정보 형식이 올바르지 않습니다": {
			"Invalid connection credential format.", "连接认证信息格式无效。",
			"Formato de credenciales de conexión inválido.", "接続認証情報の形式が無効です。",
		},
		"연결 ID를 확인해 주세요": {
			"Check the connection ID.", "请检查连接ID。", "Compruebe el ID de conexión.", "接続IDを確認してください。",
		},
		"연결 인증 정보를 입력해 주세요": {
			"Enter connection credentials.", "请输入连接认证信息。", "Introduzca las credenciales de conexión.", "接続認証情報を入力してください。",
		},
		"연결 주소를 변경할 때 인증 정보를 다시 입력해 주세요": {
			"Re-enter credentials when changing the connection address.", "更改连接地址时请重新输入认证信息。",
			"Vuelva a introducir credenciales al cambiar la dirección.", "接続先を変更する場合は認証情報を再入力してください。",
		},
		"취약점 근거·로그 연결 또는 작업 상태가 변경됐습니다. 새로고침 후 다시 조사해 주세요": {
			"Finding evidence, log connection or task state changed. Refresh before investigating again.",
			"漏洞证据、日志连接或任务状态已更改，请刷新后重新调查。",
			"Cambió la evidencia, conexión o estado de tarea. Actualice antes de investigar de nuevo.",
			"検出事項の根拠・ログ接続・タスク状態が変更されました。更新してから再調査してください。",
		},
		"침해 조사 기록을 처리하지 못했습니다": {
			"Could not process the investigation record.", "无法处理调查记录。",
			"No se pudo procesar el registro de investigación.", "侵害調査記録を処理できませんでした。",
		},
		"인증·인가로 제한된 응답": {
			"Response restricted by authentication or authorization",
			"响应受身份认证或授权限制",
			"Respuesta restringida por autenticación o autorización",
			"認証・認可によって制限された応答",
		},
		"해당 경로의 저장된 응답 %d건이 모두 401/403입니다. 보안 취약점으로 확정하지 않고 인증 조건이 갖춰진 별도 검증이 필요한 것으로 분류했습니다. AI 호출 없이 실제 응답 코드로 분류했습니다.": {
			"All %d stored responses for this path are 401/403. This indicates a need for a separate authenticated check, not a confirmed vulnerability. Classification uses response codes without an AI call.",
			"此路径保存的%d个响应均为401/403。需另行进行已认证的验证，不能据此确认漏洞。分类仅依据响应码，未调用AI。",
			"Las %d respuestas guardadas de esta ruta son 401/403. Requieren una comprobación autenticada independiente; no confirman una vulnerabilidad. La clasificación utiliza códigos de respuesta sin llamar a la IA.",
			"このパスの保存済み応答%d件はすべて401/403です。脆弱性の確定ではなく、認証条件を満たした別途検証が必要です。AIを呼ばず応答コードで分類しました。",
		},
		"인증 후 행동과 계정 간 인가를 검증하지 않았습니다.": {
			"Behavior after authentication and authorization across accounts were not verified.",
			"未验证认证后的行为或跨账户授权。",
			"No se verificaron el comportamiento autenticado ni la autorización entre cuentas.",
			"認証後の動作とアカウント間の認可は検証していません。",
		},
		"미문서화 상태코드는 별도의 명세 품질 문제일 수 있습니다.": {
			"Undocumented status codes may be a separate specification quality issue.",
			"未记录的状态码可能属于独立的规范质量问题。",
			"Los códigos no documentados pueden indicar un problema de calidad de la especificación.",
			"未記載のステータスコードは別の仕様品質の問題である可能性があります。",
		},
		"올바른 JSON 객체와 필드를 입력해 주세요": {
			"Provide a valid JSON object with supported fields.", "请输入有效的JSON对象及支持的字段。",
			"Introduzca un objeto JSON válido con campos admitidos.", "有効なJSONオブジェクトと対応するフィールドを入力してください。",
		},
		"올바른 취약점 ID를 입력해 주세요": {
			"Provide a valid finding ID.", "请输入有效的漏洞ID。",
			"Introduzca un ID de hallazgo válido.", "有効な検出事項IDを入力してください。",
		},
		"방어 검증 실행을 찾을 수 없습니다": {
			"Defense validation execution not found.", "找不到防御验证执行。",
			"No se encontró la ejecución de validación defensiva.", "防御検証の実行が見つかりません。",
		},
		"보안 로그 연결을 처리하지 못했습니다": {
			"Could not process the security log connection.", "无法处理安全日志连接。",
			"No se pudo procesar la conexión de registros de seguridad.", "セキュリティログ接続を処理できませんでした。",
		},
		"연결 설정이 변경됐습니다. 새로고침 후 다시 저장해 주세요": {
			"Connection settings changed. Refresh before saving again.", "连接设置已更改，请刷新后重新保存。",
			"La conexión cambió. Actualice antes de guardar de nuevo.", "接続設定が変更されました。更新してから再度保存してください。",
		},
		"Closed-loop 평가를 처리하지 못했습니다": {
			"Could not process the closed-loop evaluation.", "无法处理闭环评估。",
			"No se pudo procesar la evaluación iterativa.", "Closed-loop評価を処理できませんでした。",
		},
		"인증 정보가 포함된 근거 파일은 제공할 수 없습니다": {
			"Evidence containing service credentials cannot be returned.", "无法返回包含服务认证信息的证据文件。",
			"No se puede devolver evidencia con credenciales del servicio.", "サービス認証情報を含む根拠ファイルは提供できません。",
		},
		"진단 도구가 올바른 JSON 응답을 반환하지 않았습니다": {
			"The assessment tool returned an invalid JSON response.", "诊断工具返回了无效的JSON响应。",
			"La herramienta devolvió una respuesta JSON inválida.", "診断ツールが無効なJSON応答を返しました。",
		},
	} {
		errCatalog[ko] = map[string]string{"en": translations[0], "zh": translations[1], "es": translations[2], "ja": translations[3]}
	}
}
