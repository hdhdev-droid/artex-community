# 언어 설정 미적용 조사 — 2026-10-20

## 후속 반영 상태

### 로컬 검증 완료 — 2026-10-23

[로컬 빌드·테스트 및 5개 언어 화면 검증 기록](../complete/2026-10-23-language-validation.md)을 추가했다.
Go·프런트엔드 빌드, DB 없는 Go 테스트, 프런트엔드 26개 테스트 및 4,503개 번역 키 정합성 검사를 통과했다.
5개 언어 전환은 DB와 분리된 실제 데모 화면에서 확인했다. 검증 중 발견한 빌드 오류와 접근성·중국어 검색 문구도 수정했다.
실 DB 언어 저장·동기화와 외부 서비스 호출은 이번 검증에서 제외했다. 아래 소스 반영 기록은 당시 결과로 보존한다.

### 잔여 항목 반영 — 2026-10-20

사용자 요청에 따라 아래 미반영 경로를 추가 수정했다. 상세 파일·보존 기준·확인 결과는
[언어 감사 잔여 항목 완료 기록](../complete/2026-10-20-language-audit-completion.md)에 기록했다.
이 표의 완료는 **소스 반영 기준**이며, 이번 변경의 빌드·타입 검사·테스트·실행 검증 완료를 뜻하지 않는다.

| 감사 항목 | 현재 처리 |
| --- | --- |
| L01–L02 | 앞선 작업의 설치 공통 언어·저장 성공 후 전환·진입/포커스 동기화 유지 |
| L03 | reporter 표시 번역 유지. 다른 기본 에이전트도 시드 이름·설명과 정확히 일치할 때만 번역하여 사용자 수정값 보존 |
| L04 | 앞선 작업의 기본 자산 규칙 메모 표시 번역 유지 |
| L05 | 기본 실행 규칙 이름 번역 유지. 편집기에는 번역 안내를 원문 입력란과 분리해 표시. 기본 거부 문구 19개 추가 |
| L06–L07 | 앞선 작업의 조사 계획·새 보고서 템플릿 번역 유지 |
| L08 | 새 알림 본문에 이어 채널 검증·SMTP·HTTP·전송 오류의 고정 틀도 번역. 외부 응답 사유는 그대로 유지 |
| L09 | 앞선 동적 입력 오류 번역 유지. 증거 파일명·보관 ID를 끼워 넣기 전 템플릿 번역 추가 |
| L10 | 앞선 클라이언트 인증·업로드·다운로드 대체 오류 번역 유지 |
| L11 | 기본 도구 52개·파라미터 161개에 표시 요약, 기본 스킬 5개에 설명 번역 추가. 식별자·기본 원문 일치 조건으로 사용자 편집값 보호. 모델용 전체 설명·프롬프트·스키마·파일 편집 본문은 원문 유지 |
| L12–L13 | Skill/Tools/Worker/Token/Agent/Key/Method/license/compatibility/payload 표시 라벨을 번역 키로 연결. 프로토콜 필드명은 유지 |
| L14–L15 | 앞선 숫자·차트·시간선 locale 처리 유지 |
| L16 | 지적된 assembly/orchestration/server 운영 로그, 새 보관 경고·복원 활동·대화 오류 라벨 번역. 이미 저장된 활동은 재작성하지 않음 |
| L17 | 새 차단 응답의 플랫폼 정책 프레이밍과 취소·승인 거부 사유 번역. 기본 규칙 이름·메시지가 모두 일치할 때만 기본 사유 번역. 판정·정규식·승인 흐름·원본 결정 기록 유지 |
| L18 | 데모 고정 오류 49개와 동적 오류 3개, 예시 로그 6개 번역. 도달 불가능한 Coming Soon 제거. 데모 작업·취약점·대화·HTTP/증거 fixture는 입력·증거 원문으로 유지 |

위 확정 항목에 대한 이번 반영은 완료했다. 후보 CSV의 3,814행을 모두 치환한 것은 아니다.
이 목록에는 번역 사전, 주석, 모델 지침, 기술 식별자, 사용자 입력과 증거를 표현하는 fixture가 포함된다.
새 기본 설명이 원문 카탈로그와 달라지면 원문으로 표시한다. 상류 업데이트 시 해당 카탈로그와 번역을 함께 갱신해야 한다.

### 이전 반영 이력

아래 본문은 최초 조사 시점의 기록이다. 이후 [잔여 이슈 보완 기록](../complete/2026-10-20-remaining-issues.md)에서
L01–L04, L06–L07, L10, L14–L15에 해당하는 공통 언어·초기 표시·조사 계획·보고서·날짜/숫자 경로를 보완했다.
L05는 기본 규칙 이름 표시를, L08은 새 알림 렌더링을, L09는 방어/Closed-loop/Splunk의 주요 동적 입력 오류를 보완했다.
편집기 원문과 외부 오류를 전부 번역한 것은 아니다.
L11–L13, L16–L18의 나머지 기존 소스 후보와 모델용 원문은 이 보완의 완료 항목으로 계산하지 않는다.
AWS 실습 화면·API는 삭제되어 해당 기능의 번역도 제거됐다.

## 범위와 판정 기준

artex-community 작업 중인 소스를 기준으로 조사했다. 이전 검토의 번역 누락을 반영한 뒤, 프런트엔드·서버·진단기·보고서·알림·초기 데이터 생성 코드를 정적으로 검색했다.

- 소스 검색: 테스트·의존성·빌드 산출물을 제외한 Go/TS/TSX/Python/SQL 472개 파일, 155,493행.
- 별도 확인: 5개 언어 메시지 파일, 배포 설정의 기본 언어, 기본 스킬 5개 설명.
- 비밀 설정 파일, 실제 키, 업로드 파일, 실행 결과, 실제 DB 데이터는 열람하지 않았다.
- 빌드·테스트·프로세스 실행·컨테이너 재시작·DB 작업은 하지 않았다.
- 다른 작업이 진행 중인 작업 디렉터리의 스냅샷이므로 이후 변경으로 행 번호가 달라질 수 있다.

**검색 결과 전체를 미번역 오류로 세면 안 된다.** 번역 사전 원문, 프로토콜 식별자, 프롬프트, 원본 증거, 데모 데이터, 미사용 코드도 검색된다. 아래 확정 항목과 후보 목록을 구분한다. 정적 검색으로 동적 생성 문구나 모든 실제 화면 상태를 완전히 보장할 수는 없다.

## 이번에 반영한 범위

| 항목 | 결과 |
| --- | --- |
| 영어·스페인어·일본어 메시지 안의 한국어 잔존 | 이전 검토에서 언어별 840개였던 범위가 현재 0개. 작업 도중 이미 반영된 번역은 보존했다. |
| 메시지 키 | ko/en/zh/es/ja 모두 4,174개, 키 집합 동일, 중복 JSON 키 없음 |
| 화면 설정·공통 UI·LLM 기록·대화 빈 화면 | 기존 next-intl 경로로 연결 |
| 진단기 고정 문구 | 50개 원문 템플릿을 표시 시 번역. 원본 산출물은 변경하지 않음 |
| 조사 분석 고정 문구 | 24개 원문 템플릿을 결과·요약·이력 표시에서 번역 |
| 기본 에이전트 변수 예제 | 중국어 예제 3개를 알려진 초기값과 일치할 때만 번역 |
| 날짜·시간 | 이번 반영 대상 화면과 관련 패널에 기존 getDateLocale 적용 |
| Splunk 연결 확인·일부 오류 | 기존 서버 메시지 사전과 trMsg 사용 |
| 이미 등록된 운영 로그 번역 | 서버 런타임 출력 언어 사용, 초기화 전에는 구성 기본 언어 사용 |

## 추가로 확인한 미적용 항목

이 표는 **최초 추가 조사 당시의 미반영 항목**이다. 현재 처리 결과는 문서 상단 표를 기준으로 한다.

| ID / 우선순위 | 노출 경로 | 코드 근거 | 현상 / 필요한 처리 |
| --- | --- | --- | --- |
| L01 / 높음 | 언어 설정 전반 | `web/src/i18n/locale-store.ts:10`, `web/src/app/(main)/system/settings/page.tsx:124`, `server/errmsg_i18n.go:36` | UI는 브라우저 저장소, API·모델은 서버 공통 output_language를 사용한다. 다른 브라우저에서 언어를 바꾸거나 서버 저장이 실패하면 화면과 응답이 달라질 수 있다. 요청별 언어인지 설치 공통 언어인지 정책 정리가 필요하다. |
| L02 / 높음 | 신규 설치 초기 언어 | `web/src/i18n/config.ts:8`, `docker-compose.dev.yml:61` | 프런트 기본 en, 로컬 Docker의 ARTEX_LANG은 ko. 서버 output_language가 아직 없는 초기 상태에서 불일치 가능하다. |
| L03 / 높음 | 기본 보고서 작성 에이전트 | `server/orchestration.go:738`, `web/src/lib/agent-i18n.ts:17` | reporter는 builtin=false이고 현재 시스템 시드 번역 대상에는 retester만 있다. 이름·설명이 한국어로 표시된다. 사용자가 수정한 에이전트와 구분해 알려진 시드만 번역해야 한다. |
| L04 / 높음 | 기본 자산 차단 규칙 | `db/db.go:302`, `web/src/app/(main)/system/intercept/assets/page.tsx:241` | .gov/.gov.cn/.edu/.edu.cn 기본 메모 4개가 중국어이며 화면에서 rule.note를 그대로 표시한다. 사용자 메모는 보존하고 초기 메모만 표시용 번역이 필요하다. |
| L05 / 높음 | 기본 실행 차단 규칙 | `db/db.go:364`, `db/db.go:559`, `web/src/app/(main)/system/intercept/page.tsx:720` | 기본 규칙 17개와 기존 설치 보완용 규칙 2개의 이름·거부 안내가 중국어. 목록은 rule.name 원문, 편집기는 name/message 원문을 사용한다. 편집 원본과 표시 라벨을 분리해야 한다. |
| L06 / 높음 | 기본 조사 계획·독립 Defensive 계획 | `investigation/plan.go:24`, `investigation/defensive.go:7`, `web/src/components/finding-investigation-panel.tsx:95` | 계획 요약, 가설 제목·설명·조건·필요 로그·필드·조사 지침·정상 대안·일부 다음 단계가 한국어다. 이번 24개 분석 결과 번역에는 계획 문구 전체가 포함되지 않는다. |
| L07 / 높음 | Markdown 보고서·내보내기 | `report/report.go:51`, `report/findings.go:43`, `server/server.go:2288`, `server/finding_traffic_export.go:41` | 보고서 제목·섹션·집계·증거 라벨이 한국어로 생성된다. 원본 증거/기존 보고서 본문을 바꾸지 않고 새로 생성하는 문서 틀에 언어를 적용해야 한다. |
| L08 / 높음 | 알림·메일 | `notify/markdown.go:75`, `notify/html.go:46`, `notify/notify.go:45` | 심각도·상태 라벨은 이미 언어별 매핑이 있지만 요약 제목, 본문 안내, 전체 보기 링크 등은 한국어가 남는다. 채널별 템플릿과 오류를 별도로 처리해야 한다. |
| L09 / 높음 | 동적 API 오류 | `server/purple.go:46`, `server/errmsg_i18n.go:19`, `defenseeval/evaluate.go:69`, `closedloop/evaluate.go:126`, `securitysource/splunk.go:103` | trMsg는 완성 문자열의 정확한 일치만 처리한다. 숫자를 먼저 끼워 넣은 문구, `%w`로 연결한 사유, 패키지 오류·외부 오류는 사전에 원문 템플릿만 추가해도 번역되지 않는다. 오류 코드 또는 포맷 이전 번역 필요. |
| L10 / 중간 | 업로드·다운로드·인증 실패 안내 | `web/src/lib/api.ts:143`, `:470`, `:482`, `:671`, `:920`, `:1181` | 클라이언트 fallback 오류가 한국어. 응답 본문 오류와 구분해 브라우저 선택 언어의 기존 번역을 사용해야 한다. |
| L11 / 중간 | 도구·스킬 목록과 에이전트 편집기 | `server/assembly.go:220`, `server/orchestration.go:483`, `web/src/components/agent-editor.tsx:766`, `:806`, `:1270`, `web/src/app/(main)/system/tools/page.tsx:331`, `web/src/app/(main)/system/skills/page.tsx:1019` | 시드/파일에서 읽은 설명을 그대로 표시한다. 도구 설명·파라미터 설명은 모델 입력이기도 하므로 UI 표시 번역과 편집 원문을 분리해야 한다. |
| L12 / 중간 | 일부 기존 화면 라벨 | `web/src/components/agent-editor.tsx:326`, `web/src/app/(main)/function/commands/page.tsx:215`, `web/src/app/(main)/function/tasks/page.tsx:996`, `web/src/app/(main)/system/agents/detail/page.tsx:29` | Skill/Tools/Worker/Token/Agent 등 일반 명칭이 영어 고정. MCP/URL/제품명 같은 고유 식별자와 구분해 번역 필요. |
| L13 / 중간 | 기타 편집·상세 라벨 | `web/src/app/(main)/system/tools/page.tsx:747`, `:804`, `web/src/app/(main)/system/skills/page.tsx:1269`, `:1277`, `web/src/app/(main)/function/tasks/detail/_tabs/broadcast-tab.tsx:419` | Key/Method/license/compatibility/payload 등 필드 제목이 영어 고정. 스키마 키 자체는 보존하고 제목에 번역을 적용할 수 있다. |
| L14 / 중간 | 세션 토큰 수·차트 | `web/src/app/(main)/function/tasks/detail/_tabs/sessions-tab.tsx:1699`, `:1889`, `web/src/components/ui/chart.tsx:258` | 9곳의 toLocaleString()에 명시적 locale이 없어 브라우저 기본 숫자 형식을 따른다. |
| L15 / 중간 | 실행 시간선의 날짜 그룹 | `web/src/app/(main)/function/tasks/detail/_tabs/broadcast-tab.tsx:160` | dayFmt/clockFmt가 모듈 로드 시 getDateLocale로 한 번 생성된다. 새로고침 없이 언어를 바꾸면 기존 formatter가 남을 수 있다. |
| L16 / 중간 | 운영 로그·저장된 활동 메시지 | `server/orchestration.go:757`, `server/server.go:1797`, `server/task_archives.go:349` | logT/config.T를 통하지 않는 중국어 로그와 고정 한국어 활동 메시지가 있다. 이미 저장된 로그에는 언어를 바꿔도 소급 적용되지 않는다. |
| L17 / 중간 | 플랫폼 차단 상세 | `guard/guard.go:124`, `:148` | 플랫폼 차단 사유·접두 설명이 중국어로 구성된다. 실행 기록/모델 도구 응답에 남는 문구로, 규칙 이름 번역만으로 해결되지 않는다. |
| L18 / 낮음 | 데모·미사용 화면 | `web/src/lib/mock/data.ts`, `web/src/lib/mock/handler.ts:1698`, `web/src/app/(main)/system/logs/page.tsx:35`, `web/src/app/(external)/page.tsx:5` | 데모 데이터·데모 오류는 한국어/중국어가 섞여 있고 외부 페이지는 Coming Soon 고정. 실제 운영 데이터와 별도 처리 대상이다. |

## 초기 데이터 상세

DB 초기화/변경 없이 **시드 코드와 표시 경로**만 확인했다. 현재 설치 DB에 실제로 어떤 행·사용자 수정값이 있는지는 확인하지 않았다.

| 초기 데이터 | 소스 | 현재 표시 / 상태 |
| --- | --- | --- |
| 기본 agents 이름·설명 | `db/db.go:189`, `web/src/lib/agent-i18n.ts` | 기존 key 기반 번역 사용. DB의 중국어 원문 자체가 미번역 화면을 뜻하지는 않는다. |
| 기본 프롬프트 변수 라벨 | `db/db.go`, `web/src/components/agent-editor.tsx:448` | 기존 varLabel 번역 사용 |
| 기본 프롬프트 변수 예제 3개 | `web/src/lib/agent-i18n.ts:65`, `web/src/components/agent-editor.tsx:452` | 이번 반영. 원본 예제와 일치할 때만 번역 |
| retester 에이전트 | `server/finding_retests.go`, `web/src/lib/agent-i18n.ts:17` | builtin=false여도 기존 시스템 번역 대상 |
| reporter 에이전트 | `server/orchestration.go:738` | 미반영 L03 |
| 자산 차단 기본 메모 4개 | `db/db.go:302` | 미반영 L04 |
| 실행 차단 초기 규칙 17개 + 보완 규칙 2개 | `db/db.go:364`, `:559` | 미반영 L05. 설치 상태에 따라 보유 규칙은 다름 |
| 도구 description/schema 파라미터 설명 | `server/assembly.go:220`, `server/orchestration.go:483`, `db/tools.go:39` | 미반영 L11. 모델용 설명과 UI 표시가 공유됨 |
| 기본 스킬 5개 설명 | `skills/{api-recon,scopesentry,playwright-cli,api-assessment,source-assessment}/SKILL.md` | 중국어 2개·영어 1개·한국어 2개 원문 표시. 미반영 L11 |
| agent_prompts 초기 본문 | `server/assembly.go:194`, `agent/prompt.go`, `agent/promptcatalog.go` | 편집기 원문은 언어 설정으로 자동 번역되지 않음. 모델 출력 지침과 프롬프트 편집 문서의 언어는 별개 |
| 기본 MCP 설정/기술 식별자 | `db/schema.sql:860` | URL·이름·프로토콜 설정은 번역 대상이 아님 |
| 목업 작업·취약점·로그 | `web/src/lib/mock/`, `web/src/app/(main)/system/logs/page.tsx` | 미반영 L18. 실제 DB 초기 행과 구분 |

## 번역 오류로 판단하지 않은 항목

- `web/src/lib/chat-mentions.ts`의 중국어 분류는 서버와 합의된 멘션 직렬화 형식이다. 표시 이름은 별도 번역이 있다. 임의 치환하면 기존 멘션 파싱이 깨진다.
- `web/src/navigation/sidebar/sidebar-items.ts` 원문 라벨은 nav-main의 id 기반 번역을 거친다.
- `web/src/config/app-config.ts`, `web/src/i18n/config.ts`는 언어별 메타데이터 및 언어 자체의 이름이다.
- `companyScopeErrorText`의 한국어는 현재 목업 handler가 사용한다. 실제 자산 화면은 오류 코드와 번역 키를 사용하므로 운영 화면 결함으로 합산하지 않았다.
- `nav-documents.tsx`는 현재 app-sidebar에서 사용이 주석 처리된 코드다. Documents/Open/Share/Delete는 잔존 소스이지만 현재 메뉴 노출로 확정하지 않았다.
- 모델 생성 보고서, 사용자 질문·제목·메모, 업로드 원문, HTTP 응답, 외부 도구 오류, 원본 증거는 기존 언어 그대로 보존해야 한다. 새 모델 응답에는 언어 지침이 있어도 기존 저장 결과가 자동 재번역되지는 않는다.
- API 후보 검토의 고정 한계 안내와 조사 에이전트 시스템 지침에는 이미 출력 언어 경로가 있다. 모델 자체의 지침 준수까지 정적 검토로 보장하지 않는다.
- 언어 간 동일 문자열은 곧 미번역이라는 뜻이 아니다. 기술명·식별자·공유 표기가 포함된다.

## 전체 위치 목록

1. [소스 검색 후보 3,814행](2026-10-20-language-source-candidates.csv): 246개 파일의 경로·행·검출 규칙·분류. CJK 리터럴 3,118, 영어 오류 후보 649, 영어 JSX 라벨 후보 30, 영어 UI 속성 3, 기본 locale 숫자 표시 9, 스킬 설명 5. 한 위치가 여러 규칙에 걸릴 수 있으며 번역 사전과 정상 코드도 포함된다. **미번역 결함 3,814개라는 의미가 아니다.**
2. [검색 범위 472개 파일](2026-10-20-language-scan-coverage.csv): 파일별 소스 행 수와 후보 수. 검색 결과가 없는 파일도 포함.
3. [영어와 동일한 메시지 155개](2026-10-20-language-identical-values.csv): es 78, zh 39, ja 38개. 키별 검토용 목록으로 대부분 짧은 기술명·공유 표기이며 별도 판정이 필요하다.

후보 추출은 문자열·화면 라벨·오류 생성·locale 생략을 찾는 정적 규칙이다. TypeScript 전체 문법 분석이나 런타임 데이터 흐름 추적은 아니므로 정규식/템플릿/주석에 따른 오탐·누락 가능성이 있다. Python은 표준 AST로 문자열을 읽었다. 원문 내용과 비밀 값은 CSV에 복사하지 않았다.

## 후속 순서

1. L01–L02 언어 적용 범위를 정한 뒤 새 요청·새 출력의 언어를 일치시킨다.
2. L03–L06, L11은 초기값을 식별하는 표시용 매핑으로 처리한다. DB 원문이나 사용자 수정값을 덮어쓰지 않는다.
3. L07–L10, L16–L17은 새로 생성하는 보고서·알림·에러·활동 메시지에 적용한다.
4. L12–L15를 정리하고, 데모/미사용 코드는 따로 처리한다.
5. 사용자가 요청할 때 5개 언어 전환, 최초 접속, 설정 저장 실패, 내보내기, 알림, 기존 데이터 표시를 빌드·실행 환경에서 확인한다.
