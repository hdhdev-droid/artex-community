# 번역 후속 반영 및 전체 소스 조사 완료 기록 — 2026-10-20

## 요청과 처리 범위

앞서 확인한 번역 누락을 반영하고, 설정 언어가 적용되지 않는 부분을 초기 데이터 생성 코드까지 조사하여 별도 목록을 작성했다. 이번 작업으로 전체 프로젝트의 모든 문구 번역이 완료된 것은 아니다. 추가 발견 항목은 [언어 조사 목록](../todo/2026-10-20-language-audit.md)에 L01–L18로 기록했다.

빌드·테스트·DB 조회/변경·초기화·재시작·외부 API 호출은 수행하지 않았다. 비밀 설정 파일과 실제 키 값은 읽지 않았다. 같은 작업 디렉터리에서 진행 중인 마이그레이션 변경은 유지했으며 아래에는 이 언어 작업에서 추가한 변경만 기록한다.

## 변경 파일과 내용

| 파일 | 변경 |
| --- | --- |
| `web/messages/en.json`, `es.json`, `ja.json` | purple/investigation/investigationAgent/Aegis/ClosedLoop/LiveLoop/RangeLoop의 남아 있던 한국어 값을 해당 언어로 번역. 이미 번역된 값 보존 |
| `web/messages/{ko,en,zh,es,ja}.json` | uiControls 47개, chat.agentFallback 1개, 기본 예제 3개, generatedMessages 74개 추가. 5개 언어 4,174개 키 |
| `web/src/lib/generated-message-patterns.json` | 기존 진단기/조사기의 고정 원문 템플릿 74개 식별 목록 추가 |
| `web/src/lib/generated-message-i18n.ts` | 알려진 제품 생성 문구만 선택 언어로 표시하는 훅. 동적 값과 알 수 없는 원문 보존 |
| `web/src/app/(main)/offensive/assessments/page.tsx` | 진단기 요약·제약·경고 등 알려진 문구를 표시 시 번역 |
| `web/src/components/finding-investigation-panel.tsx` | 분석 결과·신호·시간선·목록·이력의 알려진 분석 문구 표시 번역, 날짜 locale |
| `web/src/lib/agent-i18n.ts` | 알려진 기본 변수 예제에만 적용하는 varExample 추가. 실제 시드 동작과 다른 주석 정정 |
| `web/src/components/agent-editor.tsx` | 변수 예제 툴팁을 varExample에 연결 |
| `web/src/app/(main)/_components/sidebar/layout-controls.tsx` | 화면 설정 라벨·선택 항목·접근성 이름 번역 |
| `web/src/app/(main)/function/llm-records/page.tsx` | Worker/Profile/Model/Tokens 및 모델 필터 안내 번역 |
| `web/src/app/(main)/chat/page.tsx` | 빈 대화의 기본 에이전트 명칭 번역 |
| `web/src/components/ui/dialog.tsx`, `sheet.tsx` | 닫기 안내 번역 |
| `web/src/components/ui/sidebar.tsx` | 사이드바 제목·설명·토글 접근성 이름 번역 |
| `web/src/components/ui/breadcrumb.tsx`, `pagination.tsx`, `carousel.tsx`, `spinner.tsx` | 탐색·이전/다음·더 보기·로딩 안내 번역. 기존 호출자가 전달한 페이지 버튼 문구는 유지 |
| `web/src/app/(main)/closed-loop/page.tsx`, `defensive/page.tsx`, `function/purple/page.tsx`, `function/tasks/page.tsx`, `system/intercept/page.tsx` | locale 없는 날짜/숫자 표시에 기존 getDateLocale 적용 |
| `web/src/components/closed-loop-live-panel.tsx`, `closed-loop-range-panel.tsx` | 실행·관측·조회 구간 날짜 locale |
| `web/src/components/finding-defense-execution-panel.tsx`, `finding-defense-panel.tsx`, `finding-defense-replay-panel.tsx` | 실행·관측·재생 이력 날짜 locale |
| `web/src/components/finding-investigation-agent-panel.tsx` | 조사 기록·실행·조회·증거 날짜/시간 locale |
| `web/src/components/side-question-workspace.tsx`, `traffic-evidence-viewer.tsx` | 보조 질문 날짜/시간과 트래픽 바이트 숫자 locale |
| `server/security_sources.go` | Splunk 연결 성공 안내를 기존 trMsg로 연결 |
| `server/security_messages.go` | 연결 인증 정보·연결 ID·조사 충돌 등의 알려진 오류 및 연결 성공 안내 번역 추가 |
| `config/msg_i18n.go` | 등록된 운영 로그 템플릿이 활성 서버 출력 언어를 사용하도록 변경, 기동 전 구성 언어 fallback 유지 |
| `docs/todo/2026-10-20-language-audit.md` | 추가 발견 18개 그룹, 초기 데이터별 상태, 정상 원문 보존 대상, 조사 한계 기록 |
| `docs/todo/2026-10-20-language-source-candidates.csv` | 전체 정적 검색의 경로·행별 후보 목록 |
| `docs/todo/2026-10-20-language-scan-coverage.csv` | 후보 없는 파일을 포함한 검색 범위 기록 |
| `docs/todo/2026-10-20-language-identical-values.csv` | 영어와 같은 es/zh/ja 메시지 키 검토 목록 |

## 정적 확인

- JSON 읽기 및 중복 키 확인: 중복 없음.
- 5개 언어 키 집합 동일: 각 4,174개.
- en/es/ja/zh 메시지 값의 한국어 잔존: 0개.
- 추가한 템플릿의 동적 자리표시자를 유지했다. 기존 ICU select의 분기 문구 차이는 실제 파라미터 누락으로 취급하지 않았다.
- 변경 코드와 표시 경로를 읽어 확인했다. 빌드·타입 검사·자동 테스트·브라우저 동작 검증은 하지 않았다.
- UI 번역은 저장된 원본 진단 파일/증거/해시/사용자 입력을 변경하지 않는다.
- 실제 설치 DB의 사용자 수정 데이터와 번역 품질의 원어민 검수는 이번 정적 조사에 포함되지 않는다.
