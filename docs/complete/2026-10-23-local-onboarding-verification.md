# 로컬 개발 배포 및 온보딩 검증

- 일자: 2026-10-23
- 결과: 로컬 빌드·회귀 테스트, 백엔드/웹 컨테이너 재생성, 기존 설치의 온보딩 및 5개 언어 화면 확인 완료.
- 승인 범위: 사용자 요청으로 로컬 빌드·테스트·컨테이너 재생성과 로컬 DB 정상 기동 처리를 수행했다. 실제 OAuth 로그인은 사용자가 직접 진행한다.

## 바로 사용하기

- 앱: http://localhost:5173
- 시작 안내: http://localhost:5173/onboarding
- 모델 연결: http://localhost:5173/system/llm
- 실행 버전: `1.0.0-dev` (v1.0.0 기반 로컬 개발 빌드).
- 기존 관리자 로그인과 세션이 유지된다. 샘플·데모·작업·모델 프로파일은 생성하지 않았다.
- 모델 연결 화면을 열어 두었다. ChatGPT는 `Continue with ChatGPT`, OpenRouter는 모델 ID를 입력한 뒤 `Continue with OpenRouter`로 사용자가 로그인할 수 있다. 연결 후 모델 프로파일을 저장하고 직접 선택·활성화해야 한다.
- 일반 접속은 localhost, 승인된 ChatGPT 콜백만 `http://127.0.0.1:8789/auth/callback`이다. 콜백 리스너는 로그인 시작 시 열리므로 시작 전 8789에 페이지가 없는 것은 정상이다.

## 배포 및 데이터 보존

- 기존 `artex-community-dev` 프로젝트의 backend/web만 `--no-deps --force-recreate`로 재생성했다. PostgreSQL과 offensive-runner는 기존 컨테이너를 유지했다.
- `pgdata`, `runtime`, `data` 이름 있는 볼륨과 기존 secret 파일을 그대로 사용했다. `down -v`, 초기화, reset, 기존 데이터 삭제는 수행하지 않았다.
- backend의 기존 개발 시작 스크립트가 최신 소스를 빌드하여 실행했다. 웹은 기존 설치된 node_modules를 사용하는 `work/local-readiness/web-offline.yml` override를 유지했다.
- backend 포트 8787~8789와 웹 5173은 호스트 loopback에만 게시된다. 일반 확인 URL은 localhost를 사용했다.
- `localhost:8787/api/health`와 웹 프록시 `localhost:5173/api/health` 모두 `ok=true`, `version=1.0.0-dev`였다. `/api/auth/status`는 `initialized=true`였다.
- 재생성 전후 PostgreSQL의 65개 public 테이블을 읽기 전용 트랜잭션으로 행 수·내용 해시 비교했다. 재기동 직후 62개 테이블이 완전히 같았다. `agents` 8행, `mcp_servers` 2행은 행 수가 유지되었고 정상 기동의 내장 설정 갱신/수정 시각 처리에 해당하는 변경이 있었다. `server_logs`는 기동 로그가 추가되었다.
- 언어 전환은 기존 화면의 저장 API로 검증했다. 검증 뒤 원래 표시 언어인 영어로 복원했다. 원래 기본 언어로 동작하던 설치에 `output_language=en` 설정 행이 저장되어 settings는 26→27행이 되었다. 이는 스키마 변경이 아니다.
- 최종 비교에서 위 4개 테이블 외 61개 테이블의 행 내용이 동일했다. 작업·취약점·증거·대화·프로파일 등 업무 데이터는 그대로이며 새 샘플을 넣지 않았다.
- 비교 파일에는 테이블명·행 수·해시만 보관했다. 실제 인증 키·토큰·DB 비밀번호는 출력하지 않았다.

## 실행한 검증

| 구분 | 결과 |
| --- | --- |
| Go 실행 파일 빌드 | `go build -mod=readonly -o /tmp/artex-candidate ./cmd/artex` 통과; 재생성 후 개발 시작 스크립트 빌드 및 기동 성공 |
| 프런트엔드 프로덕션 빌드 | 별도 web-check 컨테이너에서 `npm run build` 성공, TypeScript 검사와 39개 정적 페이지 생성 성공 |
| 프런트엔드 테스트 | `npm test -- --run`: 9개 파일, 33개 테스트 통과 |
| 번역 키 검사 | `npm run check:i18n`: en/ko/zh/es/ja의 4,591개 키 일치 |
| 백엔드 회귀 | `scripts/Test-Dev.ps1`: defenseeval, investigation, securitysource, closedloop, offensivereview, report, notify, internal/l10n, internal/chatgpt, llmpool, db, agent, server의 13개 패키지 통과 |
| DB 격리·복원 | 별도 tmpfs postgres-test에서 패키지별 테스트 DB 사용; dump/restore의 스키마 및 행 해시 일치; 검증 후 테스트 컨테이너 정지 |
| 추가 OAuth 경계 검사 | Origin 제한, 다른 세션의 OpenRouter 완료 거부, 콜백 host·만료·재사용 거부, 코드가 리다이렉트에 노출되지 않는지 확인 |
| ChatGPT 저장·전송 검사 | 암호문 변조 거부, 등록 보존, 동시 refresh 8개 요청 직렬화·회전 저장, invalid_grant 시 토큰 제거, 고정 endpoint/store/stream 제한, EOF·incomplete·실패 처리 확인 |
| ID 토큰 검사 | 테스트용 RSA/JWKS로 유효 토큰, nonce/audience/issuer/subject/expiry/issued-at/authorized-party 오류 거부 확인 |
| 동시성 검사 | `go test -race -count=1 ./internal/chatgpt` 통과 |
| API 인증 보호 | 인증 없는 onboarding MCP/ChatGPT 계정 조회는 각각 HTTP 401 |

OAuth 테스트는 임시 디렉터리·가짜 전송 계층·httptest를 사용했다. 실제 계정·외부 인증 서버·추론 서버와 통신하지 않았다. Vite가 향후 설정 로더 변경에 대한 안내를 출력했으나 현재 테스트 실패는 없었다.

## 실제 화면 확인

| 언어 | 온보딩 | 모델 / OpenRouter / ChatGPT |
| --- | --- | --- |
| English | 제목·설정 상태·언어·버튼 확인 | 안내·연결 버튼·빈 프로파일 목록 확인 |
| 한국어 | 제목·설정 상태·언어·버튼 확인 | 안내·연결 버튼·빈 프로파일 목록 확인 |
| 中文 | API·소스 진단 목적, runner 설정됨·미검증 표시 확인 | 안내·연결 버튼·빈 프로파일 목록 확인 |
| Español | 방어·조사 목적, Splunk 선택 안내 확인 | 긴 문구 줄바꿈·연결 버튼 확인 |
| 日本語 | 방어·조사 안내와 상태 확인 | 안내·연결 버튼·빈 프로파일 목록 확인 |

- 언어 선택→저장→번역 반영, 온보딩→모델 설정→돌아가기, 나중에 설정→기존 작업 목록→시작 안내 재진입을 확인했다.
- 미설정 모델/MCP/Splunk의 시험 버튼 비활성화, 목적별 필수/선택 표시, 실제 연결 시험이 자동 실행되지 않는 동선을 확인했다.
- OpenRouter 모델 ID가 비어 있으면 연결 버튼 비활성화, 임시 입력 시 활성화됨을 확인하고 입력은 비웠다. 로그인 버튼은 누르지 않았다.
- 브라우저의 기존 폭(약 722px)에서 실제 화면 캡처를 확인했다. 별도 모바일/다른 해상도 전체 검증은 수행하지 않았다.
- 재생성 전 04:30 UTC의 HMR 번역 누락 로그 1건이 남아 있었으나, 04:37:52 UTC 웹 재생성 후 5개 언어 검증 동안 새 console error/warn은 없었다.
- 로컬 화면 증빙 10개: `work/local-onboarding-verification/onboarding-{en,ko,zh,es,ja}.jpg`, `llm-{en,ko,zh,es,ja}.jpg` (ignored 작업 폴더).

## 이번 변경 파일

| 파일 | 변경 내용 |
| --- | --- |
| `scripts/check-dev-backend.sh` | 기본 로컬 회귀 범위에 internal/chatgpt, llmpool 포함 |
| `internal/chatgpt/transport_test.go` | 외부 통신 없는 전송 제한·스트림 종료·사용자 조치 오류 회귀 테스트 |
| `internal/chatgpt/store_test.go` | 암호문 변조, 등록 보존, 토큰 회전·동시성, invalid_grant 회귀 테스트 |
| `internal/chatgpt/identity_test.go` | 테스트 RSA/JWKS 기반 ID 토큰 검증 회귀 테스트 |
| `server/llm_oauth_test.go` | Origin, 세션 결속, 콜백 host/만료/중복/코드 노출 방지 회귀 테스트 |
| `docs/complete/2026-10-23-local-onboarding-verification.md` | 이번 실행 결과·데이터 비교·화면·미검증 범위 기록 |
| 기존 2026-10-23 온보딩·OAuth complete/todo 문서 | 과거 미검증 기록을 보존하고 이번 검증 기록 링크 추가 |

애플리케이션 구현 코드는 이번 검증에서 추가 수정하지 않았다. 발견된 빌드·테스트 실패나 재생성 후 재현되는 화면 오류는 없었다.

## 사용자가 확인할 범위

- 실제 ChatGPT/OpenRouter 계정 로그인·동의·콜백·프로파일 저장/활성화, 모델 목록 및 실제 모델 요청 성공 여부는 아직 검증하지 않았다.
- 실제 계정의 토큰 갱신·원격 폐기·구독 한도 오류, 외부 MCP/Splunk 연결은 미검증이다. 로컬 모의 테스트 통과와 실계정 성공은 구분한다.
- 기존 설치를 유지했으므로 초기 관리자 생성 화면을 실제 DB로 다시 실행하지 않았다. 인증/최초 설정 보호의 기존 테스트는 격리 DB에서 통과했다.
- 테스트 결과는 이번 로컬 환경과 수행한 범위에 한정된다.

## 후속 요청: 최신 컨테이너 갱신 및 ChatGPT 사용 순서

- 같은 날 사용자가 로컬 테스트용 컨테이너 갱신과 로그인→모델 선택→활성화 순서를 요청했다.
- `docker compose -f docker-compose.dev.yml build backend offensive-runner` 성공. 의존성 버전을 임의로 올리지 않고 현재 작업 폴더와 Dockerfile 기준으로 빌드했다.
- 기존 offline 웹 override를 유지한 채 `up -d --no-deps --force-recreate backend web offensive-runner`로 3개 서비스를 재생성했다. PostgreSQL 컨테이너와 pgdata/runtime/data 볼륨은 유지했다.
- 갱신 후 웹 프록시 health는 `ok=true`, `1.0.0-dev`; 브라우저 새로고침 후 `Continue with ChatGPT` 활성 상태와 빈 모델 프로파일 목록을 확인했다.
- 갱신 직전/직후 65개 테이블을 읽기 전용으로 비교했다. 62개 테이블 내용이 동일하며, 변경은 정상 기동 시 갱신되는 agents/mcp_servers 및 server_logs에 한정되었다. 언어 변경이나 사용자 데이터 추가·삭제는 수행하지 않았다.
- 이번 후속 실행에서는 회귀 테스트·프런트 프로덕션 빌드를 반복하지 않았다. 위의 통과 결과는 앞선 검증 실행 결과다. 실제 OAuth 로그인·모델 조회·추론은 여전히 사용자가 진행한다.
- `docs/local-development.md`: 현재 버튼명으로 로그인, 모델 가져오기, 저장, `Set active`, `Test connection` 절차를 구체화하고 이전 “미빌드·미기동” 안내를 최신 실행 상태로 고쳤다.
- 이 완료 파일에 갱신 결과를 추가했다. 애플리케이션 구현 코드는 변경하지 않았다.
- 화면 증빙: `work/local-onboarding-verification/llm-refreshed.jpg`.
