# 2026-10-20 잔여 이슈 보완 및 AWS 실습 제거

## 요청 범위

- AWS 실습 기능 삭제.
- 남은 개발 환경·언어 적용·검증 문제 보완.
- 사용자 추가 답변에 따라 외부 LLM·Splunk 호출은 제외했다.
- 실제 원본 데이터 이전은 제외하고 절차와 임시 DB 복원 검증만 준비했다.
- 기존 앱 DB·파일·인증 정보는 초기화·재설정하지 않았다. 설치 키를 읽거나 복사하지 않았다.

## 변경 파일

| 파일 | 변경 내용 |
| --- | --- |
| `server/closedloop_range.go`, `closedloop_range_test.go` | AWS 실습 브리지 API 및 해당 테스트 제거 |
| `server/closedloop.go` | 브리지 경로 등록 제거. 일반 실행 비교·오프라인 평가 유지 |
| `web/src/components/closed-loop-range-panel.tsx`, `web/src/lib/closedloop-range-api.ts` | 실습 화면과 API 클라이언트 제거 |
| `web/src/app/(main)/closed-loop/page.tsx`, `web/messages/{ko,en,zh,es,ja}.json` | AWS 탭과 RangeLoop 번역 제거. 과거 `tab=range`는 기본 탭으로 표시 |
| `Dockerfile.dev`, `docker-compose.dev.yml` | 개발 백엔드에 Go·Node·Python·nmap·Playwright CLI/MCP·Chromium 포함. Playwright 패키지 버전 고정. 기본 언어 en. 개발 프런트는 webpack/polling으로 마운트 변경 감지 |
| `web/src/i18n/server-locale.ts`, `locale-store.ts`, `web/src/app/(main)/layout.tsx` | 서버 `output_language`를 설치 공통 언어로 사용. 진입·포커스 복귀·표시 상태 변경 시 동기화. 오래된 읽기 응답이 새 저장 결과를 덮어쓰지 않도록 처리 |
| `web/src/app/(main)/system/settings/page.tsx`, 언어별 `settingsPage.language.desc` | 저장 성공 뒤 화면 언어 전환. 실패 시 기존 언어 유지, 중복 저장 방지, 공통 언어 안내 |
| `web/src/lib/agent-i18n.ts`, 언어별 `builtinAgents.reporter` | 알려진 reporter 초기 이름·설명만 번역. 사용자 수정값 보존 |
| `web/src/lib/generated-message-patterns.json`, 언어별 `generatedMessages` | 기본 실행 차단 규칙 이름 20개, 자산 기본 메모 4개, 조사 계획 문구 136개 추가 |
| `web/src/app/(main)/system/intercept/page.tsx`, `system/intercept/assets/page.tsx` | 규칙 이름과 기본 메모의 표시 번역. 편집 원문·DB 규칙·차단 동작 유지 |
| `web/src/components/finding-investigation-panel.tsx` | 계획 요약·가설·조건·필요 로그/필드·조사 지침·정상 대안·다음 단계·참고자료 제목 표시 번역 |
| `internal/l10n/template.go`, `templates.json` | 고정 문서·검증 템플릿의 다국어 변환. 사용자 데이터나 이미 완성된 문서는 입력하지 않음 |
| `report/language.go`, `report.go`, `findings.go` | 새 Markdown·CSV 보고서의 제목·항목·집계·증거 링크 라벨 번역. 사용자 제목·보고서 본문·원본 증거 유지 |
| `notify/language.go`, `markdown.go`, `html.go`, `feishu.go`, `telegram.go`, `event.go` | 새 알림 제목·요약·메일/메시지 템플릿·링크 문구 번역. 실제 발송은 하지 않음 |
| `defenseeval/evaluate.go`, `closedloop/evaluate.go`, `securitysource/splunk.go`, `server/purple.go` | 입력 검증 오류를 동적 값 삽입 전에 번역 |
| `server/errmsg_i18n.go`, `security_messages.go`, `server.go` | 알려진 래핑 오류 접두부 번역, 공통 런타임 언어를 보고서·검증 패키지에 연결 |
| `web/src/i18n/client-text.ts`, `web/src/lib/api.ts` | 인증·업로드·다운로드 클라이언트 대체 오류 문구에 현재 언어 적용 |
| `web/src/app/(main)/offensive/assessments/page.tsx` | 요약이 없는 실행 기록의 번역 호출 타입 오류 수정 |
| `web/src/app/(main)/function/tasks/detail/_tabs/sessions-tab.tsx`, `broadcast-tab.tsx`, `web/src/components/ui/chart.tsx` | 토큰·차트 숫자 locale 지정, 시간선 날짜 포맷이 초기 언어에 고정되는 문제 수정 |
| `web/src/i18n/server-locale.test.ts`, `security-messages.test.ts` | 저장 실패·늦은 응답·다른 브라우저 언어 동기화 및 신규 번역 형식 검사 |
| `report/language_test.go`, `notify/language_test.go`, `internal/l10n/template_test.go` | 원본 증거 보존, 문서·알림 번역, 포맷 인수 보존 검사 |
| `server/logt_catalog_test.go` | 서버 런타임 언어가 다른 테스트에서 남는 영향을 격리. 초기 언어와 런타임 언어 우선순위를 각각 확인 |
| `scripts/check-dev-backend.sh`, `Test-Dev.ps1`, `check-dev-restore.sh` | 문서 패키지 검사 추가. 별도 테스트 DB 덤프·복원 및 public 테이블별 행 수·내용 지문 대조. PowerShell 줄바꿈 호환 |
| `docs/local-development.md`, `docs/migration-checklist.md`, 언어 감사 문서 | 현재 환경과 실제 이전 절차·검증 한계 기록 |

## 검증 결과

- `defenseeval`, `investigation`, `securitysource`, `closedloop`, `offensivereview`,
  `report`, `notify`, `internal/l10n`, `db`, `agent`, `server` 테스트 통과.
- 최초 서버 검사에서 로그 번역 테스트가 실패했다. 런타임 언어 우선순위 변경 후
  테스트 간 전역 설정이 남는 문제였으며, 초기 상태를 격리하고 런타임 우선순위 검사를 추가한 뒤 통과했다.
- 임시 PostgreSQL 덤프·복원 후 테이블별 행 수·내용 지문 일치.
  첫 실행은 검증 자체는 통과했지만 PowerShell의 마지막 CR 문자 때문에 실패 코드가 반환됐다.
  입력 줄바꿈을 정규화한 뒤 복원 검증을 다시 실행해 정상 종료를 확인했다.
- 프런트 타입 검사 및 7개 파일의 21개 테스트 통과.
- 5개 언어 키 4,192개 정합성 및 신규 보안/생성 문구 ICU 포맷 검사 통과.
- 프런트 배포 빌드 통과, 정적 페이지 36개 생성.
- 개발 이미지 빌드 통과. Go 1.26.3, Python 3.11.2, nmap 7.93,
  Playwright 1.64.1, CLI 0.1.22, MCP 0.0.83 설치 및 명령 실행 확인.
- 로그인된 실제 로컬 브라우저에서 Closed-loop 탭이 실행 비교·오프라인 평가 2개만 표시됨을 확인.
- 설정 화면의 설치 공통 언어 안내 문구 표시 확인. 사용자 언어 설정값은 변경하지 않음.
- `/api/health` 정상. backend·web·offensive-runner 실행 중, 앱 PostgreSQL healthy.
- 검증 후 `postgres-test`만 정지했다.
- `git diff --check` 통과. Git의 LF/CRLF 안내와 Vite/webpack 의존성 분석 경고는 남지만 검사·빌드 실패는 아니다.

## 범위와 한계

- 실제 외부 연결·공격·알림 발송·운영 배포·실제 원본 데이터 이전은 사용자 선택에 따라 제외했다.
- 복원 검증은 테스트로 생성한 DB에 대한 검증이다. 실제 설치의 업로드 파일·암호화 정보까지 이전한 결과가 아니다.
- 기존 증거, 모델 생성 본문, 사용자 편집값, 모델용 도구·스킬 설명, 과거 저장 기록은 재작성하지 않는다.
- 이번 변경은 위 항목의 보완이며, 기존 전역 언어 감사 L01–L18의 모든 소스 후보를
  번역 완료로 판정한 작업은 아니다. 잔여 후보와 원문 유지 대상은 감사 문서에 구분했다.
