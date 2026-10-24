# 상류 추가 변경·로컬 사용 가능 여부 확인 — 2026-10-23

## 상류 확인

사용자의 추가 반영 여부 요청에 따라 `git ls-remote https://github.com/jiwoochris/artex-global.git HEAD refs/heads/main 'refs/tags/*'`로 원격 참조만 조회했다.

- 원격 HEAD/main: `fdfb7a17b7da5e2de716c0e98475ec77a41c64c6`.
- community에 기록된 직접 베이스와 동일하므로 **원격 main 기준 추가 반영 대상 커밋은 0개**다.
- 원격 태그: `v0.3.15-global.1`, `v0.3.15-ko.1`. global 태그가 가리키는 커밋은 `d6964594a0221cd921b85056d604bce9d8ce54df`로 로컬 확인과 일치했다.
- `v0.3.15-global.1` 이후 9개 커밋은 이미 베이스에 포함된다. 열린 PR·미병합 브랜치·원본 ARTEX 저장소의 신규 변경까지 검토한 결과는 아니다.
- fetch/merge/cherry-pick, 원격 변경 및 베이스 메타데이터 갱신은 하지 않았다.

## 로컬 기동 상태와 복구

backend, postgres, offensive-runner는 계속 실행 중이었으나 web 컨테이너는 종료 상태였다. 최근 로그에서 Next 설정 변경 후 `/version/VERSION` 파일을 읽지 못해 종료된 것을 확인했다. 현재 소스에는 메타데이터 부재 시 처리와 Compose 읽기 전용 마운트가 이미 반영되어 있었다.

중지된 web만 기존 이미지·설치 의존성으로 다시 기동했다. `work/local-readiness/web-offline.yml`은 `npm ci`를 실행하지 않고 `npm run dev -- --webpack --hostname 0.0.0.0 --port 5173`만 실행하는 임시 Compose 덮어쓰기 파일이다.

```powershell
docker compose -f docker-compose.dev.yml -f work/local-readiness/web-offline.yml up -d --no-deps --no-build --pull never web
```

백엔드·DB·진단 서비스는 재시작하지 않았다. 패키지 설치·이미지 pull·배포 빌드·테스트는 하지 않았다. 개발 서버의 화면 제공 과정에서 Next 개발 컴파일은 발생했다.

## 확인 결과

- `http://localhost:8787/api/health`: HTTP 200, `ok: true`, `service: artex`, `version: dev`.
- `http://localhost:5173/`: 비인증 HTTP 요청은 로그인으로 이동하여 HTTP 200 응답.
- Chrome에 남아 있는 인증 세션으로 `http://localhost:5173/function/tasks` 실제 작업 목록을 확인했다. `Version · dev`, 현재 작업 **0건** 표시, 수집한 브라우저 오류 로그 없음. 인증 값은 읽거나 출력하지 않았다.
- `http://localhost:5183/`: 기존 정적 샘플 화면 HTTP 200. 실제 DB 연동 환경과 별도다.
- postgres healthy, backend/offensive-runner는 기존 프로세스 유지. 진단 실행 자체의 성공을 확인했다는 의미는 아니다.
- 증거 화면: `work/local-readiness/tasks.png`. 작업 폴더는 Git 추적 대상이 아니다.

## 사용 범위와 남은 사항

로컬 개발 환경은 `http://localhost:5173`에서 바로 접속할 수 있다. 실제 작업을 만들면 로컬 DB에 저장되는 환경이다. 이번 확인에서는 작업 생성·설정 저장·진단 실행·DB 변경을 하지 않았다.

실행 중인 백엔드는 이전 `dev` 바이너리다. 최신 v1.0.0 버전/베이스 응답 및 이후 백엔드 수정의 실제 적용은 별도 재빌드·재기동이 필요하다. 현재 확인을 최신 백엔드 코드 전체의 검증으로 간주하지 않는다. 외부 LLM/Splunk/실제 대상 진단은 실행하지 않았다. 첫 릴리스의 GitHub Actions 실행·배포 검증도 아직 남아 있다.
