# 2026-10-20 로컬 Docker 개발 환경 구축

## 요청 및 범위

남은 이전 작업을 정리하고 DB를 포함한 로컬 Docker 개발 환경을 기동했다.
이번 환경 구축에 필요한 Go 컴파일·진단 이미지 빌드·의존성 설치·개발 서버 기동을 수행했다.
기존 DB·데이터·Redis를 초기화하거나 이전하지 않았다.
테스트 스위트·린트·타입 검사·프런트엔드 배포 빌드는 수행하지 않았다.

## 변경 파일

| 파일 | 변경 내용 |
| --- | --- |
| `docker-compose.dev.yml` | 별도 `artex-community-dev` 프로젝트에 PostgreSQL 16, Go 백엔드, Next 개발 서버, 진단 서비스 구성. 호스트는 루프백에만 포트 노출. 데이터·캐시 볼륨 분리 |
| `scripts/Start-Dev.ps1` | 누락된 개발 인증 파일만 생성하고 컨테이너 빌드·기동. 기존 인증 파일 유지, 값 출력 없음 |
| `scripts/dev-backend.sh` | 컨테이너 안에서 인증 파일을 읽고 백엔드를 컴파일·기동. 별도 데이터·런타임 경로 사용 |
| `docs/local-development.md` | 접속, 최초 설정, 재시작·중지, 저장소와 남은 기능 안내 |
| `README.md`, `README.ko.md` | 로컬 개발 안내 링크 추가 |

`.secrets/dev/`에 생성한 DB 비밀번호·설정 토큰·서비스 토큰은 기존 ignore 규칙으로
Git 및 이미지 빌드 컨텍스트에서 제외된다. 기존 JWT 서명 값이나 모델 인증 값을 조회·복사하지 않았다.

## 실행 결과

- Docker Desktop 엔진 29.8.0 확인.
- 기존 `redis-server` 컨테이너 유지.
- 새 서비스 4개 모두 실행 중, PostgreSQL healthy.
- `http://localhost:5173/setup` HTTP 200.
- `http://localhost:5173/api/auth/status` HTTP 200,
  `initialized: false`, `setup_token_required: true`.
- `http://localhost:8787/api/health` HTTP 200, `ok: true`, `version: dev`.
- 백엔드에서 내부 진단 서비스 `/health` 응답 `status: ok` 확인.
- 새 DB는 앱의 기존 자동 스키마·내장 항목 적용 후 public 테이블 51개.
  읽기 전용 조회로 `offensive_candidate_reviews`가 없는 것을 확인했다.
- DB 스키마 소스 수정, 별도 SQL init/reset, 기존 DB 접속은 하지 않았다.

## 사용자가 설정할 항목과 후속 작업

- 최초 화면에서 `.secrets/dev/setup-token`을 사용해 관리자 비밀번호를 설정한다.
- 모델 연결은 비어 있다. 사용자 모델 설정 전에는 AI 엔진이 대기한다.
- AI 후보 검토 저장·중단·이어하기는 별도 테이블 적용 승인 및 구현 이후 활성화 가능하다.
- 실제 업로드·진단·다운로드·AI 검토, 기존 화면·인증 회귀는 아직 검증하지 않았다.
- 동적 문구 번역, 방어 검증·침해 조사·Splunk·Closed-loop 이전은 남아 있다.
- 개발용 Go 이미지에 운영용 전체 외부 도구를 설치한 것은 아니다.

기동 상태 확인과 전체 기능 검증을 구분한다. 기존 1~2단계 완료 기록의
빌드·기동 미실행 설명은 그 작업 당시 상태이며, 이번 기록이 이후 개발 환경 기동 결과다.
