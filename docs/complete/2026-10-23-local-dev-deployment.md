# 최신 로컬 개발 배포·확인 — 2026-10-23

## 요청과 승인 범위

- 사용자가 최신 로컬 개발 버전 재기동·테스트와 샘플/데모 종료를 요청했다.
- 백엔드 `db.Open`이 스키마·내장 설정을 자동 적용하는 점을 설명하고, 기존 DB 조회 전용 규칙의 예외로 **로컬 DB의 정상 기동 처리 허용** 답변을 받은 뒤 재기동했다.
- DB 초기화·리셋·삭제, 실제 진단 실행, 외부 LLM/Splunk 호출은 하지 않았다.

## 반영 및 실행

- `docker-compose.dev.yml`: web에 `NEXT_PUBLIC_MOCK=0` 명시. 기존 DB를 사용하는 실제 개발 화면으로 고정했다.
- 기존 backend를 최신 소스로 재빌드·재기동했다. `GET http://localhost:8787/api/health`에서 `ok:true`, `version:1.0.0-dev` 확인.
- web은 기존 설치 의존성을 사용한 `work/local-readiness/web-offline.yml` 덮어쓰기로 재생성했다. 이미지 pull이나 npm 설치를 하지 않았다.
- PostgreSQL·offensive-runner 컨테이너 및 영속 볼륨은 그대로 유지했다.
- 명령행을 확인한 뒤 `work/report-validation/serve.py`만 종료했다. `localhost:5183` 연결 실패로 종료 확인. 샘플 소스/과거 검증 문서는 삭제하지 않았으며 실제 개발 모드에서는 사용하지 않는다.
- `docs/local-development.md`에 현재 버전·실제 접속 주소·샘플 종료를 기록했다.

## 빌드·테스트

기존 이미지·캐시를 사용하는 `--pull=never --network none` 격리 컨테이너에서 실행했다. 운영 DB 접속 정보나 비밀 파일은 검증 작업 디렉터리에 복사하지 않았다.

| 검증 | 결과 |
| --- | --- |
| `go build -mod=readonly ./...` | 종료 코드 0 |
| `go test -mod=readonly ./...` | 종료 코드 0, 22개 테스트 패키지 통과·2개 테스트 없음. DB 의존 테스트는 실제 DB 연동 성공으로 간주하지 않음 |
| 실행 파일 빌드 및 `-version` | 빌드 종료 코드 0, 별도 격리 실행 결과 `1.0.0-dev` |
| `npm test -- --run` | 9개 파일·33개 테스트 통과 |
| `npm run check:i18n` | 5개 언어 4,511개 키 일치 |
| `NEXT_PUBLIC_MOCK=0 npm run build` | 종료 코드 0, TypeScript 검사 및 프로덕션 빌드 완료 |

임시 검증 스크립트 마지막 줄의 Windows 줄바꿈 때문에 버전 명령 리다이렉션/빌드 종료 파일명 확인에 문제가 있었다. LF로 정리했으며, 이미 완료된 빌드·테스트의 종료 코드를 회수하고 빌드된 바이너리의 버전을 별도 실행해 확인했다. 앱 결함은 아니며 실패한 스크립트 종료를 전체 성공으로 간주하지 않았다. 종료된 검증 컨테이너는 정리했다.

## 실제 화면·데이터 확인

- `http://localhost:5173/function/tasks`: 실제 작업 0건, `Version · 1.0.0-dev` 표시.
- 시스템 설정: community `v1.0.0`, `hdhdev-droid/artex-community`, `v0.3.15-global.1-9-gfdfb7a1`, 정확한 베이스 커밋 및 원본 ARTEX `0.3.15` 표시.
- API·소스 진단: 로딩 완료 후 등록 폼 활성화, 실행 이력 없음. 방어 검증 화면 로딩, 조사 화면 저장 사건 0건, Closed-loop 저장된 실행 없음 확인.
- 위 탐색에서 수집한 브라우저 오류 로그 없음. 개발 서버 최초 페이지 컴파일로 메뉴 이동이 지연된 경우 최종 로딩 완료를 확인했다.
- 기존 설정 화면의 업데이트 조회는 실제 GitHub 릴리스가 없다는 응답을 표시한다. 정식 릴리스 미발행 상태이며 로컬 앱 사용을 차단하지 않는다. 새 릴리스를 발행한 것은 아니다.
- 재기동 전후 및 완료 후 테이블 열 구조: 702개, 지문 `ff6c2559c8c739c0ff51eb412ca8e33a` 동일. 작업 0건 유지. 모든 DB 행/인덱스의 완전 동일성을 검증했다는 의미는 아니다. 승인된 기동 시 내장 설정 처리 자체는 실행됐다.
- 실제 작업 생성·변경·삭제와 외부 연동·진단 실행까지 테스트한 것은 아니다. 사용자 테스트가 가능한 상태로 기동했다.

## 결과 위치

- 사용자 접속: **http://localhost:5173/function/tasks** (인증 세션이 없으면 기존 계정으로 로그인).
- 버전/베이스: **http://localhost:5173/system/settings**.
- 로그·검증 결과: `work/dev-deploy/`.
- 화면 증거: `work/dev-deploy/tasks-current.png`, `work/dev-deploy/version-settings.png`.
