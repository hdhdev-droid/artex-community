# 2026-10-20 보안 기능 후속 반영

## 요청과 승인 범위

1~2단계 이후의 후보 검토 저장소, 방어 검증, 침해 조사, Splunk, Closed-loop를
artex-ko-plus 소스에서 선별 반영했다. 사용자가 추가로 승인한 범위에 따라
로컬 개발 DB에 테이블을 추가하고 빌드·테스트·동작 확인을 수행했다.
기존 외부 DB·서비스에는 접속하지 않았으며 AWS 인프라를 만들거나 실제 공격을 실행하지 않았다.

## 기능별 변경

| 기능 | 변경 파일과 내용 |
| --- | --- |
| 저장 스키마 | `db/schema.sql`: 후보 검토, 방어 계획·수동 관측·재생·실행·판정, 로그 연결, 조사·조사 실행, 독립 사건·실행, Closed-loop 평가·실행·관측용 14개 테이블과 인덱스·재검증 스냅샷 트리거 추가 |
| 기본 스킬 연결 | `db/db.go`: 기존 스킬 노출 시드에 API·소스 진단 스킬만 추가. DB 연결 풀 상한 유지 |
| 방어 검증 | `defenseeval/evaluate.go`, `parse.go`, `db/purple.go`, `purple_replay.go`, `purple_execution.go`, `server/purple.go`, `purple_replay.go`, `purple_execution.go`, `purple_upload.go`: 수동 기록, 로컬 로그 재생, 실제 재검증 연계, 로그 파일 판정, 개선 전후 비교 |
| 재검증 연결 | `db/finding_retests.go`, `server/finding_retests.go`: 일반 재검증 경로를 공유하면서 방어 실행을 트랜잭션으로 연결. 실행별 검증 ID를 문맥에 전달 |
| 로그 연결 | `db/security_sources.go`, `server/security_sources.go`, `securitysource/splunk.go`, `upload.go`, `investigation.go`, `investigation_agent.go`, `closedloop_rule.go`: Splunk 설정 저장, 로컬 JWT 파일에 묶인 인증 정보 암호화, TLS 검증, 제한된 읽기 전용 조회·파일 처리 |
| 조사 분석 | `investigation/types.go`, `parse.go`, `plan.go`, `analyze.go`, `defensive.go`: 고정 대상·기간·경로의 근거 분석 및 독립 사건 분석 |
| 조사 에이전트 | `agent/investigator.go`, `investigator_types.go`, `investigator_tools.go`, `investigator_context.go`, `investigator_provider.go`: 전용 읽기 도구·예산·근거 참조·중단 처리, 기존 모델/기록 경로 사용, 기존 출력 언어 지시 적용 |
| 모델 입력 가림 | `agent/investigator_redact.go`: 전송 전 복사본에서 `_KEY` 필드·할당, 비밀번호·인증 헤더·Bearer 값 가림. 원본 근거와 큰 정수 ID 보존 |
| 조사 저장·복구 | `db/investigations.go`, `investigation_runs.go`, `defensive_cases.go`, `defensive_case_runs.go`, 대응하는 `server/` 파일: 사건·실행·근거·취소·내보내기 및 재기동 시 중단 상태 복구 |
| Closed-loop | `closedloop/evaluate.go`, `db/closedloop.go`, `closedloop_live.go`, `server/closedloop.go`, `closedloop_live.go`: 독립 평가 데이터의 규칙 비교와 실제 재검증·로그 관측의 연결. 오프라인 결과와 실제 경보를 구분 |
| 선택형 실습 연결 | `server/closedloop_range.go`: 고정된 서버 설정으로만 접근하는 실습 브리지 API. 미설정 상태로 유지하며 외부 실습 인프라는 실행하지 않음 |
| 통합 조회 | `db/aegis.go`, `server/aegis.go`: 진단·사건·평가 수 집계 API. 기존 대시보드는 교체하지 않음 |
| 작업 보관·복원 | `db/task_archives.go`, `task_archives_restore.go`: 방어·조사 이력 보존 및 복원, 진행 중인 조사·재검증의 보관 차단, 과거 보관 파일과의 호환 |
| 서버 연결 | `server/server.go`: 기존 인증 미들웨어 안에 신규 API 등록, 조사 취소 상태와 복구 연결. 최초 설정 보호·비밀번호 변경에 따른 세션 무효화·언어 설정·판정 토큰 집계 유지 |
| 오류·자동 판정 번역 | `server/security_messages.go`, `offensive_review_evidence.go`: 공통 오류와 401/403 자동 분류 문구에 다국어 적용 |
| 신규 화면 | `web/src/app/(main)/defensive/page.tsx`, `closed-loop/page.tsx`, `function/purple/page.tsx`, `function/purple/sources/page.tsx`: 독립 조사·실행 비교·방어 현황·Splunk 설정 |
| 상세 화면 | `web/src/app/(main)/function/findings/detail/page.tsx`: 방어 검증·침해 조사 탭 추가. 기존 상세·날짜 언어·상속 읽기 전용 처리 유지 |
| 화면 패널 | `web/src/components/defense-state-badge.tsx`, `finding-defense-panel.tsx`, `finding-defense-execution-panel.tsx`, `finding-defense-file-upload.tsx`, `finding-defense-replay-panel.tsx`, `finding-investigation-panel.tsx`, `finding-investigation-agent-panel.tsx`, `closed-loop-live-panel.tsx`, `closed-loop-range-panel.tsx` |
| 프런트 API | `web/src/lib/purple-api.ts`, `security-api.ts`, `defensive-api.ts`, `investigation-api.ts`, `investigation-agent-api.ts`, `investigation-mapping.ts`, `closedloop-api.ts`, `closedloop-live-api.ts`, `closedloop-range-api.ts`: 기존 `ApiError`·인증·Blob 다운로드 함수와 연결 |
| 메뉴·언어 | `web/src/navigation/sidebar/sidebar-items.ts`, `web/messages/{ko,en,zh,es,ja}.json`: 기존 메뉴 보존, 신규 메뉴와 번역 키 연결, 주요 조작·안내 번역 |
| 격리 검증 | `docker-compose.dev.yml`, `scripts/Test-Dev.ps1`, `scripts/check-dev-backend.sh`: 휘발성 테스트 DB, 패키지별 독립 DB, 프런트엔드 별도 빌드 캐시. 기존 개발 서비스와 분리 |
| 테스트 보강 | 각 신규 패키지 및 `db/`, `server/`, `agent/`, `lab/offensive/`의 관련 테스트 이전. `agent/investigator_redact_test.go` 추가. `server/task_templates_test.go`의 불필요한 백그라운드 실행 제거. `securitysource/closedloop_rule_test.go`의 대용량 응답 fixture 안정화 및 원문 출력 제거 |
| 안내 문서 | `docs/local-development.md`, `docs/offensive-assessments.md`: 저장소 준비 상태와 검증 방법을 현재 구현에 맞게 갱신 |

## 로컬 적용

- `artex-community-dev`의 앱 DB만 적용했다.
- 적용 전 백업: `.secrets/dev/community-before-followup-20261020-170053.dump`.
  인증 정보가 포함될 수 있으므로 공유·커밋하지 않는다.
- 백엔드 재기동 시 기존 추가형 스키마 적용 경로를 사용했다.
- public 테이블 수 51 → 65, `offensive_candidate_reviews` 존재 확인.
- 관리자 계정은 이미 설정된 상태임을 상태 API로만 확인했다. 비밀번호·모델 인증 값은 조회하거나 변경하지 않았다.
- 앱·프록시·PostgreSQL·진단 서비스는 계속 로컬 Docker 컨테이너로 실행한다.
- 검증 완료 후 `postgres-test`만 정지했다. 앱 DB와 다른 개발 컨테이너는 유지했다.
- 별도로 추가되어 있던 Linux/macOS 시작 스크립트와 AMD64 진단 이미지 설정을 보존했다.

## 검증 결과

- `defenseeval`, `investigation`, `securitysource`, `closedloop`, `offensivereview`,
  `db`, `agent`, `server` 패키지 테스트 통과.
- DB·에이전트·서버 테스트는 각각 다른 임시 DB를 사용했다.
  초기 공유 DB 재실행 시 남은 작업을 서버가 복구하면서 테스트 임시 디렉터리 정리와 충돌했다.
  패키지별 DB 분리와 해당 CRUD fixture 정리 후 서버 전체 테스트가 통과했다.
- 대용량 Splunk 응답 fixture에서 미소비 요청 본문으로 인한 연결 종료가
  크기 제한 대신 수신 중단으로 관측된 사례가 있었다. 요청을 소비하도록 고친 후 관련 패키지가 통과했다.
- Python 실행기: 35개 테스트 통과. 외부 네트워크가 없는 임시 컨테이너 사용.
- 프런트엔드: 6개 파일, 18개 테스트 통과.
  신규 `web/src/i18n/security-messages.test.ts`에서 5개 언어의 신규 보안 문구를 실제 포맷팅해 검사했다.
  `closed-loop-live-panel.tsx`의 JSON 예제를 번역 문법으로 처리하던 오류를 `t.raw` 사용으로 수정했다.
  Windows 마운트의 수정 내용이 개발 서버에 즉시 반영되지 않아 재시작했으며,
  새로 연 화면에서 해당 오류가 사라지고 오프라인 비교가 정상 실행됨을 확인했다.
- TypeScript 검사 통과.
- 프런트엔드 배포 빌드 통과, 신규 화면을 포함한 정적 페이지 36개 생성.
  첫 빌드는 Next가 `next-env.d.ts`를 쓰려 해 읽기 전용 마운트에서 실패했으며,
  컨테이너 내부 소스 복사본을 쓰도록 바꾼 뒤 통과했다.
- 최종 작업 트리에서 5개 언어의 번역 키 4,097개 정합성 검사 통과. 이 검사는 번역 품질이나 모든 문구의 번역 완료를 의미하지 않는다.
- 로컬 백엔드 `/api/health`, 프런트 프록시 `/api/auth/status` 정상 응답.
- 실행 중인 개발 서버가 새 경로를 인식하지 못해 처음 404가 발생했다.
  프런트엔드 재시작 후 로그인된 브라우저에서 신규 메뉴와 방어 현황 화면이 정상 표시되는 것을 확인했다.
- API·소스 진단 화면에서 저장소 미준비 경고가 사라지고 실행 목록·입력 화면이 표시되는 것을 확인했다.
- Closed-loop 화면의 합성 예제로 오프라인 비교를 실행했다.
  평가용 4개 이벤트에서 기준 규칙 TP/FP/TN/FN=1/1/1/1, 후보 규칙=2/0/2/0이 표시됐다.
  실제 경보 출처가 없으므로 MTTD는 계산 불가로 표시됐다.
  비교 결과는 DB에 저장하지 않았으며 실제 공격·탐지 성능을 검증한 결과가 아니다.

## 검증 한계와 별도 설정

- 실제 외부 모델, 운영 Splunk, 외부 공격 대상은 호출하지 않았다.
  모델 호출·조회·취소·복구는 가짜 제공자와 테스트 서버를 이용해 검증했다.
- AWS·실습 브리지 배포와 전체 운영용 브라우저·침투 도구 설치는 하지 않았다.
  실습 탭은 브리지 미설정 상태이며 자동으로 외부 환경을 만들지 않는다.
- 한국어·중국어 신규 화면 문구를 가져오고, 다른 세 언어의 주요 조작·고지와 공통 오류를 번역했다.
  긴 도움말·일부 동적 분석 문구에는 한국어 원문이 남아 있다. 모든 신규 문구의 완전한 5개 언어 번역은 완료하지 않았다.
- 업로드한 원본 로그는 로컬 증거 보존·내보내기에 포함된다.
  가림 처리는 모델 입력용 복사본에 적용하며 모든 임의 형식의 비밀 탐지를 보증하지 않는다.
- 기존 실제 데이터 이전, 다른 설치의 암호화 인증 정보 이전, 운영 배포는 수행하지 않았다.
