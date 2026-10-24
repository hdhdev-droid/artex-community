# 언어 감사 잔여 항목 반영 — 2026-10-20

## 요청과 범위

`docs/todo/2026-10-20-language-audit.md`의 잔여 항목을 ko/en/zh/es/ja 표시 경로에 반영했다.
기존 기능과 사용자 입력·증거 원문을 보존했다. 기존 작업 디렉터리의 다른 변경은 되돌리지 않았다.
현재 사용자 규칙에 따라 빌드·타입 검사·테스트 실행·DB 작업·컨테이너 재시작은 하지 않았다.
외부 LLM·Splunk·번역 서비스 호출과 실제 데이터 이전도 하지 않았다.

## 파일별 변경

| 파일 | 변경 내용 |
| --- | --- |
| `web/src/lib/builtin-description-i18n.ts` | 기본 도구·스킬·파라미터의 표시 번역 함수. 도구는 system 여부·식별자·원문 일치, 스킬은 이름·원문 일치, 파라미터는 시스템 도구의 알려진 원문 일치 조건 적용 |
| `web/src/lib/builtin-description-sources.json` | 기본 도구 52개 전체 원문을 식별용으로 보존. 번역 요약을 모델용 설명으로 저장하지 않음 |
| `web/src/lib/builtin-parameter-sources.json` | 파라미터 설명 161개의 식별 원문. 숫자·커서·ID·DSL·공유 증거 스키마·상속 작업 제한 포함 |
| `web/src/lib/builtin-skill-sources.json` | 기본 스킬 5개의 식별 원문. 실제 SKILL.md와 모델 지침은 수정하지 않음 |
| `web/src/lib/builtin-agent-sources.json`, `agent-i18n.ts` | 알려진 기본 에이전트 이름·설명과 일치할 때만 번역. 사용자 수정 이름·설명은 화면에서도 원문 유지 |
| `web/src/components/agent-editor.tsx` | 스킬·도구 탭 제목, 스킬 설명·도구 요약·트리거 도구 설명 표시 번역 |
| `web/src/app/(main)/system/tools/page.tsx` | 목록에 번역 요약 표시·검색 지원. 상세 편집기에는 원문과 번역 요약을 분리. 파라미터도 같은 방식으로 표시. 에이전트 이름·식별자·요청 방식 라벨 번역 |
| `web/src/app/(main)/system/skills/page.tsx` | 선택한 기본 스킬 설명, 라이선스·호환성 라벨 번역. 파일 편집기 본문은 그대로 유지 |
| `web/src/app/(main)/system/agents/page.tsx`, `system/agents/detail/page.tsx` | 스킬·에이전트·식별자 제목을 기존 번역 경로에 연결 |
| `web/src/app/(main)/function/commands/page.tsx`, `function/tasks/page.tsx`, `function/tasks/detail/_tabs/broadcast-tab.tsx` | 실행 에이전트·토큰·전달 내용 제목 번역. 값과 실행 내용은 그대로 유지 |
| `web/src/app/(main)/system/intercept/page.tsx`, `web/src/lib/generated-message-patterns.json` | 규칙 편집 원문 아래 알려진 이름·거부 문구의 번역 안내 추가. 입력 상태와 저장 요청은 원문 유지 |
| `config/followup_i18n.go`, `messages_followup.json`, `msg_i18n.go` | 기존 운영 메시지 경로의 추가 사전. 65개 고정 원문을 정확히 조회하며 동적 값은 번역 후 포맷 |
| `server/assembly.go`, `orchestration.go`, `server.go` | 감사에서 지적된 시드·프롬프트 갱신·작업·도구·MCP·트래픽 운영 로그를 기존 logT 경로에 연결 |
| `guard/guard.go`, `default_rules.json` | 플랫폼 차단 프레이밍·취소·승인 거부를 현재 언어로 생성. 기본 규칙 19개의 이름과 사유가 모두 원래 값일 때만 사유 번역 |
| `internal/l10n/templates.json` | 알림 오류와 보관 경고·복원 활동·보관 ID 오류 템플릿 추가 |
| `notify/dingtalk.go`, `email.go`, `feishu.go`, `filter.go`, `http.go`, `mask.go`, `telegram.go`, `webhook.go`, `wecom.go` | 고정 검증·전송 오류를 포맷 전에 번역. 원인 오류, HTTP 응답, SMTP 응답, 사용자 값은 포맷 인자로 유지. 재시도·영구 오류 구분·주소 검증은 변경하지 않음 |
| `server/task_archives.go` | 새 정리 경고·복원 활동 제목·동적 보관 ID 오류 번역 |
| `server/security_messages.go`, `offensive_review_evidence.go`, `conversations.go` | 증거 해시·무결성 오류 템플릿과 새 대화 오류 제목 번역. 파일명·오류 상세·증거 내용은 보존 |
| `web/src/i18n/client-text.ts`, `web/src/lib/mock/handler.ts` | React 외부의 데모 오류에 현재 언어 적용. 범위 오류는 기존 scopeEditor 오류 코드를 재사용. 의존 작업 ID·범위 행 번호는 번역 후 삽입 |
| `web/src/app/(main)/system/logs/page.tsx` | 데모 모드의 알려진 예시 로그 6개만 표시 시 번역. 실제 SSE/과거 로그 원문은 그대로 표시 |
| `web/src/app/(external)/page.tsx` | 기존 리다이렉트 뒤 실행되지 않는 Coming Soon 반환문 제거 |
| `web/messages/{ko,en,zh,es,ja}.json` | 위 표시 요약·설명·라벨·오류·예시 로그의 5개 언어 메시지 추가 |
| `guard/guard_test.go`, `server/conversations_localized_test.go`, `web/src/i18n/security-messages.test.ts` | 기존 고정 언어 기대값에 언어를 명시하고 사용자 차단 사유 보존 회귀 조건·신규 메시지 범위를 추가. 테스트 실행은 하지 않음 |
| `docs/todo/2026-10-20-language-audit.md` | L01–L18의 현재 처리 상태와 역사적 조사 결과를 분리 |

## 원문 보존 방식

- 도구·파라미터 번역은 **읽기용 요약**이다. 모델이 사용하는 전체 설명, 파라미터 스키마와 기본값, 편집 필드와 저장 요청은 원문을 유지한다.
- 알려진 원문과 다르면 번역 매핑을 적용하지 않는다. 사용자 정의 도구 설명은 번역하지 않는다.
- 차단 여부·규칙 우선순위·정규식·승인 흐름은 변경하지 않는다. `Interceptor.Log`에는 기존 원본 결정을 기록하고, 새 모델 응답의 플랫폼 안내만 지역화한다.
- 사용자 제목·메모·질문·모델 응답·증거·HTTP 요청/응답·외부 오류·이미 저장된 로그는 재작성하지 않는다.
- `web/src/lib/mock/data.ts`의 시나리오 제목·본문·대화·취약점·트래픽은 실제 입력과 증거를 흉내 내는 fixture이므로 동일한 보존 기준을 적용한다.
- DB 시드 변경·실제 DB 수정·키 열람·키 복사·키 전송은 하지 않았다.

## 정적 확인

- TypeScript 파서로 소스 228개를 읽어 문법 오류가 없음을 확인했다. 타입 검사나 트랜스파일·빌드는 실행하지 않았다.
- 신규 메타데이터·데모 메시지의 5개 언어 템플릿 1,380개를 ICU 문법 파서로 확인했다.
- 메시지 키는 5개 언어 모두 4,502개로 일치하며 중복 JSON 키가 없다.
- 추가 운영 메시지 65개와 공통 템플릿 전체 159개의 행 구조·중복 원문·printf 인자 순서를 확인했다.
- 기본 거부 규칙 19개의 메시지가 추가 운영 메시지 사전에 모두 존재한다.
- 수정한 Go 파일 19개를 gofmt로 구문 파싱·서식 정리했다.
  컨테이너의 소스 마운트가 읽기 전용이어서 최초 직접 저장은 실패했다.
  마운트를 바꾸지 않고 포맷 결과를 받아 호스트 작업 폴더에 저장했다.
- 정적 확인 중 도구 편집기 JSX의 닫는 중괄호 누락을 수정한 후 소스 파싱을 다시 확인했다.
- `git diff --check`에서 공백 오류가 없었다.

## 실행 확인 범위

이번 변경에는 빌드·타입 검사·단위/통합 테스트·브라우저 언어 전환·실제 알림 발송·모델 실행 검증을 수행하지 않았다.
앞선 작업의 통과 결과를 이번 수정의 실행 검증 결과로 재사용하지 않는다.
실행 검증이 필요하면 사용자의 요청에 따라 기존 로컬 개발 환경에서 진행한다.
상류 기본 도구·스킬·파라미터 설명이 바뀌면 식별 원문과 표시 번역을 함께 갱신한다.
