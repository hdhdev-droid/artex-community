package offensivereview

import "github.com/Autumn-27/artex/config"

func reviewLanguage() string {
	switch config.ActiveLanguage() {
	case "ko":
		return "Korean"
	case "zh":
		return "Chinese"
	case "es":
		return "Spanish"
	case "ja":
		return "Japanese"
	default:
		return "English"
	}
}

func reviewLimitation(source bool) string {
	switch config.ActiveLanguage() {
	case "ko":
		if source {
			return staticLimit
		}
		return apiLimit
	case "zh":
		if source {
			return "仅根据提供的源代码进行静态判断，未验证运行环境中的可利用性或实际复现。"
		}
		return "仅审核已保存的API请求和响应。未发起新请求或跨账号比较，响应不代表漏洞利用成功。"
	case "es":
		if source {
			return "Conclusión estática basada en el código proporcionado; no se verificó la explotación ni la reproducción en ejecución."
		}
		return "Solo se revisaron solicitudes y respuestas guardadas. No se hicieron nuevas solicitudes ni comparaciones entre cuentas; las respuestas no demuestran explotación."
	case "ja":
		if source {
			return "提示されたソースの静的判断です。実行環境での悪用可能性や再現成功は確認していません。"
		}
		return "保存済みのAPIリクエストとレスポンスのみ確認しました。新規リクエストやアカウント間比較は行っておらず、応答は悪用成功を示しません。"
	default:
		if source {
			return "Static assessment of supplied source evidence only; runtime exploitability and successful reproduction were not verified."
		}
		return "Only stored API requests and responses were reviewed. No new requests or cross-account comparisons were performed; observed responses do not establish exploitation."
	}
}
