# community 업데이트 원본·버전 관리 — 2026-10-23

## 요청과 반영

요청: 업데이트 대상을 `artex-community`로 맞추고 버전을 관리할 수 있도록 구성.
로컬 `origin`이 `git@github.com:hdhdev-droid/artex-community.git`인 것을 확인했다. 상류 코드 검토 대상 `jiwoochris/artex-global`과 실제 설치 대상 community 릴리스를 분리했다.

| 파일 | 변경 내용 |
| --- | --- |
| `selfupdate/github.go` | 업데이트 저장소를 `hdhdev-droid/artex-community`로 변경. 기존 HTTPS·체크섬·숫자 버전 비교·개발 버전 차단 유지. |
| `.github/workflows/release.yml` | community 저장소 제한, 정식 `vMAJOR.MINOR.PATCH` 태그 검사, Linux 바이너리 버전과 태그 일치 확인 추가. community GHCR 주소 사용. 바이너리와 이미지 작업 성공 뒤 Release 게시. |
| `docker-compose.yml`, `install.sh` | 실제 기본 이미지와 설치 실패 안내를 `ghcr.io/hdhdev-droid/artex-community`로 일치시킴. |
| `cmd/artex/main.go` | `-version` 추가. DB·업데이트 부트스트랩·서버 시작 전에 버전만 출력하고 종료. |
| `build-cloud.sh` | 기존 `build.sh`와 같이 명시 버전 또는 Git describe/dev를 실행 파일에 주입. |
| `web/src/lib/types.ts`, `web/src/lib/mock/handler.ts` | 업데이트 응답의 demo 모드 추가. 샘플에서 dev/community 원본을 반환하고 설치·롤백은 거부. |
| `web/src/app/(main)/system/settings/_components/update-card.tsx` | 업데이트 원본 표시. 샘플임을 표시하고 조회 버튼 비활성화, 가짜 OS/CPU·독립 실행 표시 및 재시작 안내 제거. |
| `web/messages/{ko,en,zh,es,ja}.json` | 원본·데모 모드·제한 안내 및 데모 작업 거부 메시지 추가. |
| `detections/tests/sigma_match/events/artex_selfupdate_egress.json` | self-update egress 예제의 저장소 경로를 실제 community API 경로와 일치시킴. 탐지 User-Agent는 유지. |
| `MAINTAINING.md`, `MAINTAINING.en.md` | 태그 기반 버전 정책, 배포 단계, 개발 빌드·기존 설치 전환·Docker 차이 기록. |
| `README.md`, `README.ko.md`, `README.es.md`, `README.ja.md` | Docker 이미지 및 설치 안내를 community 주소로 변경. |

## 버전 정책

- Git 태그 `vMAJOR.MINOR.PATCH`를 단일 릴리스 버전 기준으로 사용. 별도 VERSION 파일이나 프런트엔드 패키지 버전을 앱 버전과 중복 관리하지 않는다.
- 실제 바이너리에는 `v` 없는 버전이 들어가며 `/api/health`, 업데이트 판정, `-version`이 이를 사용한다.
- 개발/dirty 빌드는 정식 버전으로 자동 업데이트하지 않는다.
- 이번 작업에서는 첫 버전 번호를 임의 확정하거나 태그·커밋·원격 Release를 발행하지 않았다. 로컬 태그 목록은 비어 있었으며 원격 릴리스 현황은 조회하지 않았다.
- 최초 배포는 검증된 커밋에 승인된 태그를 발행한 뒤 수행한다. 기존 실행 파일의 원본 주소는 소스 수정만으로 바뀌지 않으므로 새 community 빌드를 먼저 수동 배포해야 한다.

## 확인 및 미실행 범위

- 코드·워크플로 의존 순서·리소스 주소·번역 JSON의 변경 부분을 정적으로 검토했다.
- 변경 대상 `git diff --check`에서 오류가 없었으며, 5개 언어 JSON 파싱 및 새 번역 키 존재를 확인했다. 컴파일·기능 테스트 결과를 의미하지 않는다.
- 빌드·테스트, GitHub 조회·push·발행, Docker 이미지 빌드·pull·재기동, 실제 업데이트/롤백, DB 조회·변경은 실행하지 않았다. 기존 기능과 사용자 입력·증거 데이터는 변경하지 않았다.
- 현재 `localhost:5183`은 이전에 생성한 정적 샘플 산출물을 제공한다. 이번 소스 변경은 재빌드하기 전까지 그 화면에 나타나지 않는다.
- 후속 검증은 요청 시 로컬 빌드·테스트로 진행할 수 있다. 실제 태그 발행, 5개 플랫폼 설치, GHCR 접근 권한 및 community Release 다운로드·재시작·롤백의 종단 검증은 별도 배포 단계다.
