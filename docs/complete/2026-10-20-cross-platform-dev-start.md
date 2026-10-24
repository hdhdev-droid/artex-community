# 2026-10-20 Linux·macOS 개발 환경 시작 지원

## 변경 파일

- `scripts/start-dev.sh`: Linux와 macOS 공통 POSIX 셸 시작 스크립트 추가.
  Docker 연결 확인 후 없는 개발 인증 파일만 생성하고 기존 개발 Compose 실행.
  시스템 난수 장치와 기본 도구를 사용하며 값은 출력하지 않는다.
  신규 인증 파일은 소유자만 읽고 쓸 수 있도록 생성하고 기존 파일은 유지한다.
- `docker-compose.dev.yml`: 진단 서비스만 `linux/amd64`로 지정.
  Apple Silicon에서는 Docker 에뮬레이션으로 스캐너 바이너리를 실행하도록 구성.
  다른 서비스와 기존 Windows 시작 방식은 유지한다.
- `docs/local-development.md`: Windows·Linux·Intel Mac·Apple Silicon Mac 준비 조건,
  실행 명령, 에뮬레이션 및 검증 범위 설명 추가.

## 확인 범위

기존 Compose·Windows 시작 스크립트와의 경로·인증 파일·실행 명령 연결을 정적으로 확인했다.
기존 `.gitattributes`가 셸 스크립트의 LF 줄바꿈을 지정하므로 변경하지 않았다.
이번 작업에서는 컨테이너 재시작·빌드·테스트·DB 접근·인증 파일 생성을 수행하지 않았다.
Linux 및 Mac 실기기 실행은 아직 검증하지 않았다.
기존의 Windows Docker Desktop 기동 확인과 구분한다.
