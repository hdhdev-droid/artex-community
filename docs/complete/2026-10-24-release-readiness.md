# 2026-10-24 운영 인증 보존·릴리스 준비·로컬 갱신

## 요청과 범위

실제 OAuth 로그인은 사용자가 수행한다. 기존 데이터를 유지하면서 운영 JWT·ChatGPT 인증 정보 영속화, 남은 로컬 검증, v1.0.0 릴리스 준비와 개발 컨테이너 갱신을 진행했다. 앞서 승인된 로컬 빌드·테스트·정상 기동 처리를 적용했다. 실제 서비스 연결이나 진단 실행, 운영 데이터 이전, 원격 배포·GitHub 발행은 실행하지 않았다.

## 변경 파일

| 파일 | 변경 내용 |
| --- | --- |
| `server/auth_storage.go` | `ARTEX_AUTH_DIR` 선택. 기본 실행 파일 옆 경로는 유지. 탐색 가능한 data 아래 경로·심볼릭 링크 우회·기존 파일의 일부만 옮긴 상태를 거부 |
| `server/auth.go` | JWT 읽기 실패·손상 시 새 키로 덮어쓰지 않음. OAuth 파일만 남아 있을 때 생성 금지. 최초 생성은 배타적 쓰기, 파일 동기화. 기존 data 경로 이전 오류를 무시하지 않음 |
| `server/server.go` | 선택한 인증 경로를 JWT와 ChatGPT가 함께 사용. 기존 암호화 파일을 읽지 못하면 복원 안내 후 기동 중단 |
| `docker-compose.yml`, `Dockerfile` | 운영 `./auth:/app/auth` 마운트와 인증 경로 설정. 이미지에도 `/app/auth` 영속화 지점 선언. Node 22와 개발 환경에서 검증한 Playwright 버전으로 정렬 |
| `.gitignore`, `.dockerignore` | 인증 디렉터리·중첩 JWT 파일·로컬 설정·검증 산출물이 커밋 또는 이미지 빌드 문맥에 섞이지 않도록 제외 |
| `scripts/migrate-docker-auth.py` | 중지된 이전 컨테이너에서 JWT와 암호화된 OAuth 파일을 함께 복사. 사전 점검, 동일 파일 재실행, 충돌 시 중단. 원본 삭제·내용 출력 없음 |
| `scripts/test_migrate_docker_auth.py` | 이전 도구의 사전 점검·재실행·불일치 보호 검사 |
| `server/auth_storage_test.go`, `internal/chatgpt/store_test.go` | 파일 이전 후 세션·암호화 저장소 복원, 손상·누락·경로 보호, 새 저장소 인스턴스에서 인증 복원, 오프라인 연결 해제 검사 |
| `server/chat_llm_resolve_test.go`, `core_test.go`, `mgmt_test.go`, `task_categories_test.go`, `task_control_routes_test.go`, `task_archives_test.go` | 인증 디렉터리를 작업 파일 디렉터리와 분리하는 실제 운영 조건으로 테스트 픽스처 수정 |
| `securitysource/closedloop_rule_test.go` | 큰 응답 전에 요청 본문을 소비하도록 모의 서버 수정. 간헐적 연결 리셋을 제거하며 실제 응답 중단 검사는 유지 |
| `scripts/check-dev-backend.sh`, `scripts/Test-Dev.ps1` | 업데이트 교체·롤백 검사 추가. Windows PowerShell 5에서 복원 스크립트의 `r` 문자까지 제거하던 줄바꿈 인용 오류 수정 |
| `.github/workflows/release.yml` | 백엔드·인증 이전·경쟁 상태·격리 PostgreSQL 검사 추가. 검증 성공 후 5개 플랫폼 패키징. 수동 실행은 공개 발행하지 않으며 정식 태그 push만 Release·GHCR 발행. 기본 권한 최소화 |
| `.github/workflows/ci.yml`, `build.sh` | CI 설명 갱신, 릴리스 ZIP에 인증 운영 안내 포함 |
| `MAINTAINING.md`, `docs/local-development.md`, `docs/deployment-auth.md`, `docs/releases/v1.0.0.md` | community 릴리스 기준·검증·기존 컨테이너 이전·백업·복원·롤백 안내 정리 |

## 검증 결과

- 백엔드 14개 패키지 통과: defenseeval, investigation, securitysource, closedloop, offensivereview, report, notify, internal/l10n, internal/chatgpt, llmpool, selfupdate, db, agent, server.
- DB·agent·server는 별도 휘발성 테스트 DB로 실행했다. 최초 관리자 설정 보호와 비밀번호 변경에 따른 세션 무효화 등 기존 인증 검사를 포함한다. 앱 DB에서 초기화·리셋·테스트 데이터 삽입은 하지 않았다.
- 격리 DB 덤프·복원 후 모든 테이블의 행 지문 일치. 수정한 `Test-Dev.ps1`을 Windows PowerShell 5에서 처음부터 끝까지 재실행해 통과했다.
- ChatGPT 경쟁 상태 검사(`-race`) 통과. 실제 OAuth 제공자는 호출하지 않고 모의 응답으로 갱신·오류·연결 해제를 검사했다.
- 인증 이전 Python 검사 3개가 Windows·Linux에서 통과. 가짜 인증 파일을 넣은 별도 Docker 컨테이너로 실제 `docker cp` 사전 점검·복사·재실행까지 성공했다. 실 운영 인증은 이 도구로 옮기지 않았다.
- 프런트엔드 9개 파일·33개 테스트, 5개 언어 4,591개 번역 키 정합성, TypeScript 포함 정적 릴리스 빌드(39개 페이지) 통과.
- Linux amd64/arm64, macOS amd64/arm64, Windows amd64의 **1.0.0** ZIP 및 SHA256SUMS 생성·검증. Linux amd64 도움말·버전 실행, Windows amd64 버전 실행 통과.
- 운영 Docker Linux amd64 이미지 `artex-community:1.0.0-local` 빌드·실행 확인. 인증 경로 `/app/auth`, 인증 파일 미포함, Node·Python·Playwright 실행 확인. Docker의 AUTH_DIR 관련 경고는 비밀 값이 아닌 경로 환경 변수에 대한 이름 기반 경고다.
- 운영 Compose와 릴리스 YAML 검증. 수동 실행에는 태그를 선택하더라도 발행하지 않도록 이벤트 조건도 제한했다. 실제 GitHub Actions 실행과 GHCR 발행은 하지 않았다.
- 실제 localhost 화면에서 온보딩 영어·한국어·중국어·스페인어·일본어 전환·저장 확인. 폭 390px 모바일 화면과 키보드 Tab 이동 확인. 원래 영어와 기본 화면 크기로 복원했다. LLM 화면에서 ChatGPT 연결 버튼·기존 모델 프로필을 확인했고 브라우저 오류는 없었다. OAuth 시작·모델 호출·활성 프로필 변경은 하지 않았다.

## 데이터와 로컬 실행 상태

- 주소: `http://localhost:5173/system/llm?from=onboarding`. 백엔드 health 정상, 개발 표시 **1.0.0-dev**. 정식 패키지는 **1.0.0**으로 구분한다.
- backend·web·offensive-runner를 재빌드/재생성했다. 마지막 인증 경로 보완 뒤 backend를 다시 컴파일·기동했다. PostgreSQL 앱 컨테이너와 pgdata·runtime·data 볼륨은 유지했다.
- 컨테이너 재생성 전후 JWT와 기존 ChatGPT 암호화 파일의 동일성을 내용 출력 없이 확인했다. 브라우저의 기존 로그인도 유지됐다. 기존 LLM 프로필 1개는 행 내용까지 동일했다.
- 65개 테이블 비교: 재생성 직후 62개 동일. agents 8개·mcp_servers 2개는 정상 기동 갱신, server_logs는 기동 로그가 추가됐다. 언어를 저장한 settings도 변경됐으며 27개 행과 원래 영어 설정은 유지했다. 사용자 작업·증거·모델 프로필·보안 소스는 보존했다. 샘플·데모 데이터는 넣지 않았다.
- 검증 자료와 5개 플랫폼 ZIP은 Git에서 제외된 `work/final-release-verification/`에 보관한다. 실제 인증 파일 내용·키·토큰은 기록하지 않는다.
- 마지막 재기동까지 포함한 최종 비교는 65개 중 61개 테이블이 동일했고, 차이는 위 정상 기동·언어 설정 4개 테이블에 한정됐다. 임시 인증 이전·웹 빌드 컨테이너는 정리하고 테스트 PostgreSQL은 중지했다. 앱 PostgreSQL은 계속 healthy 상태다.

## 사용자 또는 발행 단계에서 남은 확인

1. 실제 ChatGPT/OpenRouter 로그인, 사용할 계정·모델 선택, 프로필 활성화와 모델 연결 테스트. 기존 프로필을 자동 활성화하지 않았다. 실제 토큰 만료·제공자 폐기·요금 한도는 사용자 계정으로 확인해야 한다.
2. 실제 MCP·Splunk·대상 진단과 전체 작업 실행은 연결 설정 후 별도 확인. 이번 검증은 외부 서비스에 요청하지 않았다.
3. 변경 검토·커밋 후 GitHub 수동 검증, 최종 태그 push, GHCR 공개 읽기 권한과 운영 설치 확인. 현재 원격 발행이나 태그 생성은 하지 않았다.
4. macOS와 ARM64의 실제 실행 및 ARM64 운영 Docker 빌드는 해당 환경/Actions에서 확인한다. 이번 로컬 결과는 교차 컴파일과 Linux amd64 이미지 검증이다. 새 설치 화면에서 실제 관리자 비밀번호 입력은 하지 않았으며, 첫 설정 경로는 격리 DB 통합 검사로 검증했다.
