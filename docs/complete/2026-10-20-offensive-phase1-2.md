# 2026-10-20 완료 기록: API·소스 진단 1~2단계 코드 반영

## 요청과 범위

artex-ko-plus의 공통 호환성·API/소스 진단·후보 검토를 artex-community에 선별 반영했다.
기존 기능 유지, 빌드·테스트 및 DB 변경 금지 조건을 적용했다.
앞서 수정한 포크 출처 문서는 보존했다.

## 변경 파일

| 파일 | 내용 |
|---|---|
| `server/server.go` | 기존 인증 미들웨어 안에 진단 라우트 등록, 후보 검토 런타임 필드 추가 |
| `server/offensive_assessments.go` | 대상 업로드, 진단 실행·조회, 산출물 다운로드, 기존 작업 기반 검토 연결 |
| `server/offensive_review_evidence.go` | 저장된 진단 근거·해시 확인, 제한된 소스 발췌, 인증 제한 자동 분류; Windows 경로 처리 및 `_KEY` 값 가림 보완 |
| `server/offensive_reviews.go` | 후보별 큐·중단·이어하기·저장, 저장소 준비 여부 조회 및 미준비 시 호출 거절 |
| `db/offensive_reviews.go` | 체크포인트·소유권·임대·복구 어댑터 소스 추가. 스키마 생성·실제 DB 작업 없음 |
| `offensivereview/model.go`, `types.go`, `language.go` | 후보 1건 단위 모델 호출, 결과·근거 참조 검증, 기존 출력 언어 적용 |
| `server/errmsg_catalog.go` | 신규 저장소·연결 안내의 번역 추가 |
| `web/src/lib/api.ts` | 기존 ApiError·인증 처리 유지, FormData 및 Blob 응답 지원 |
| `web/src/lib/offensive-api.ts` | 신규 API 타입·호출 추가, 업로드도 기존 공통 HTTP 함수 재사용 |
| `web/src/app/(main)/offensive/assessments/page.tsx` | 등록·범위 선택·진단·이력·후보·다운로드·모델 선택·묶음 검토·중단·이어하기 화면 |
| `web/src/navigation/sidebar/sidebar-items.ts` | 기존 메뉴 순서·대시보드를 유지하면서 진단 메뉴 추가 |
| `web/messages/{ko,en,zh,es,ja}.json` | 신규 화면·상태·메뉴 번역 추가 |
| `lab/offensive/runner.py`, `api_discovery.py`, `rules/*.yml` | 별도 API·소스 진단 서비스와 로컬 분석 규칙 이전 |
| `Dockerfile.offensive-runner`, `docker-compose.offensive.yml` | 선택형 진단 서비스 이미지·앱 연결 구성. 기본 배포 파일 유지 |
| `.dockerignore` | 새 구성의 인증 파일과 로컬 환경이 빌드 컨텍스트에 들어가지 않도록 제외 |
| `skills/api-assessment/SKILL.md`, `skills/source-assessment/SKILL.md` | 진단 검토용 스킬 파일 이전. DB 노출 시드 변경 없음 |
| `docs/offensive-assessments.md` | 사용 흐름, 연결 조건, 제한, 후속 검증 항목 |

## 미변경·미실행

- `db/schema.sql`, `db/db.go`, 기존 인증·실행 엔진·대시보드 및 기본 Compose 설정은 변경하지 않았다.
- 기존 DB·인증 정보·데이터를 복사하지 않았다.
- 빌드·테스트·타입 검사·린트·진단 실행·모델 호출·Docker 빌드·프로세스 기동은 수행하지 않았다.
- DB 접속·SQL 실행·스키마 적용·init·reset은 수행하지 않았다.
- 확인은 변경 파일의 코드·호출 연결·Git 변경 목록을 읽는 정적 검토에 한정했다.

## 완료와 보류 구분

1~2단계의 기능 연결 코드를 반영했다. 런타임 검증 완료나 즉시 사용 가능 상태를 의미하지 않는다.
진단 서비스는 운영자가 별도로 준비해야 하며, 기존 앱 Dockerfile은 사전 빌드된 community 바이너리가 필요하다.

기본 community DB에는 `offensive_candidate_reviews` 테이블이 없으므로
**AI 후보 검토·판정 저장·이어하기 활성화는 보류**다.
DB 변경 금지 조건을 지키기 위해 자동 생성·임시 저장 우회를 넣지 않았고, 미준비 시 API와 UI에서 차단한다.
도구 진단 결과는 별도 파일 저장소를 이용한다.

신규 화면과 모델 출력 언어는 기존 다국어 체계에 연결했지만, 가져온 도구 설명·일부 오류·자동 판정 원문은
번역이 추가로 필요하다. 방어 검증·침해 조사·Closed-loop·AWS 실습 기능은 이번에 반영하지 않았다.
