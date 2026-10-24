# ChatGPT 구독 OAuth 코드 반영

> 후속 로컬 빌드·테스트·컨테이너 재생성 및 5개 언어 화면 확인은 [개발 배포 검증 기록](2026-10-23-local-onboarding-verification.md)에 기록했다. 아래의 미수행 표기는 최초 코드 반영 시점 기준이다. 실제 OAuth 로그인은 사용자 확인 단계로 남아 있다.

- 일자: 2026-10-23
- 상태: 코드 반영·소스 검토·서식 정리 완료 / 빌드·타입 검사·테스트·실계정·화면 검증 미수행
- 승인: “인증 콜백만 127.0.0.1 허용”. 일반 앱 주소는 localhost 유지.
- 선행 작업: [온보딩·5개 언어·OpenRouter](2026-10-23-onboarding-openrouter.md).

## 반영 내용

1. 기존 관리자 인증 뒤 ChatGPT 연결을 시작한다. 최초 등록은 dynamic client, 후속 로그인은
   저장된 issued client ID를 사용한다. 호스트 UUID를 유지하고 PKCE S256·state·nonce를
   요청마다 생성한다. state는 ARTEX 로그인 세션에 결속하며 5분·1회 사용으로 제한한다.
2. 콜백 전용 주소는 `http://127.0.0.1:8789/auth/callback`이다. 로그인 버튼에서 리스너를
   준비하고, 유효한 콜백 코드는 서버에 남긴 뒤 원래 localhost 화면으로 돌려보낸다.
   인증된 완료 요청만 코드를 교환한다. 콜백의 임의 origin/client ID 변경은 허용하지 않는다.
3. 고정 OpenAI 인증 엔드포인트와 공개 JWKS를 사용한다. RS256 서명·issuer·audience·exp·iat·
   nonce·subject·authorized party를 확인한다. 검증 실패 시 기존 연결의 토큰을 덮어쓰지 않는다.
   최초 코드 교환에 실패해도 발급 client ID는 미연결 등록으로 보존하여 새 로그인에 재사용한다.
4. 토큰 응답의 실제 권한을 확인한다. 계정/워크스페이스 등록을 분리하고 계정 참조만 기존
   `llm_profiles` 저장 경로에 넣는다. 새 DB 테이블·스키마·의존 라이브러리는 추가하지 않았다.
5. 토큰은 기존 설치 JWT 서명 키에서 별도 도메인으로 파생한 AES-GCM 키로 암호화한다.
   실행 파일 옆 보호된 런타임 디렉터리에 임시 파일→동기화→rename으로 저장한다.
   Unix 파일 권한은 0600/디렉터리 0700이다. 실제 키·토큰 값은 읽거나 출력하지 않았다.
6. 만료 전 갱신, 회전된 토큰 일괄 저장, 동시 갱신 직렬화, 재인증·해제 경쟁 방지와
   invalid_grant 재연결 안내를 추가했다. refresh 요청에는 scope를 재전송하지 않는다.
7. 모델은 버튼을 누를 때 계정의 `/v1/models`에서 가져오고 표시 가능한 순서를 유지한다.
   새 프로파일은 조회한 모델로 저장하며 자동 활성화/자동 대체 대상으로 추가하지 않는다.
8. 기존 Responses SDK를 재사용한다. 구독 프로파일은 고정 `/v1/responses`·`store=false`·
   `stream=true`를 사용한다. 내부 Complete 호출도 스트림을 누적해 반환한다.
   `response.completed` 전 EOF·incomplete·failed는 실패로 처리한다. 인증/구독 한도 오류는
   SDK 및 작업 재시도/폴백에서 자동 재요청하거나 다른 유료 공급자로 넘기지 않는다.
9. 계정 해제는 discovery의 허용된 revocation endpoint로 최대 두 번 시도한다. 원격 확인
   여부와 별개로 로컬 토큰을 제거하고 등록 정보는 보존한다. 원격 확인 실패는 별도 표시한다.
   이미 진행 중인 외부 요청/작업의 취소는 포함하지 않는다.
10. 기존 LLM 설정에 계정 선택·재인증·모델 조회·프로파일 저장·해제·사용량 관리, 로그인 완료
    안내를 추가했다. 연결 관리 항목은 일반 API 키 입력으로 덮어쓸 수 없게 했다.
    채팅 모델 선택 영역에 사용하는 ChatGPT 계정/플랜을 표시한다. 5개 언어 문구를 추가했다.
    사용자 입력·모델 출력·증거 및 기존 작업 이력은 수정하지 않았다.

## 변경 파일

| 파일 | 변경 내용 |
| --- | --- |
| `internal/chatgpt/store.go` | 암호화 파일, 호스트·계정 등록, 토큰 교환/갱신, 모델 목록, 해제 |
| `internal/chatgpt/identity.go` | 공개 JWKS와 ID 토큰 검증, 폐기 엔드포인트 조회 |
| `internal/chatgpt/transport.go` | 고정 추론 주소, 매 요청 토큰 해석, 완료 이벤트 검증, 사용자 조치 필요 오류 분류 |
| `agent/chatgpt.go`, `agent/provider.go` | 기존 SDK 연결, Complete의 스트림 누적, 구독 요청의 프록시/리다이렉트/자동 재시도 제한 |
| `server/chatgpt_oauth.go` | 세션 결속 OAuth API, 제한된 콜백 리스너, 계정·모델·프로파일 연결 |
| `server/chatgpt_messages.go` | 기존 오류 카탈로그에 5개 언어 안내 등록 |
| `server/server.go` | 저장소/라우트 등록, 해제된 계정의 준비 상태, 연결 시험 오류 안내 |
| `server/server_mgmt.go`, `server/dto.go` | 프로파일 계정 메타데이터, OAuth 설정 보호, 모델 조회 분기 |
| `server/task_llm.go`, `server/engine.go`, `llmpool/pool.go` | 구독 한도·재인증 오류의 자동 재시도 및 공급자 전환 차단 |
| `web/src/app/(main)/system/llm/_components/chatgpt-connect.tsx` | 계정 연결 관리 카드 |
| `web/src/app/(main)/system/llm/chatgpt/callback/{page,layout}.tsx` | 토큰 없는 완료 화면, no-referrer |
| `web/src/app/(main)/system/llm/page.tsx` | 연결 카드, 계정 상태, 관리되는 입력 보호 |
| `web/src/app/(main)/chat/page.tsx` | 사용 중인 계정/플랜 및 사용량 링크 |
| `web/src/lib/types.ts` | 비밀 없는 프로파일 인증 메타데이터 |
| `web/messages/{en,ko,zh,es,ja}.json` | 연결 UI·오류 안내 번역 |
| `docker-compose{,.dev}.yml` | 콜백 환경 및 호스트 loopback 8789 포트 |
| `.gitignore`, `.dockerignore` | 런타임 인증 파일 배포/커밋 제외 |
| `README.ko.md`, `docs/local-development.md` | 사용 순서, 저장소, 콜백, 이후 기동 조건 |
| `docs/todo/2026-10-23-{chatgpt-oauth,first-run-onboarding}.md` | 승인/반영 상태와 남은 검증 |
| `docs/complete/2026-10-23-onboarding-openrouter.md` | 후속 구현 기록 연결 |

## 실제 수행 및 남은 확인

- 수행: 공식 OpenAI 문서와 설치된 SDK 소스 확인, 변경 호출 경로 수동 검토, Go/TSX 서식 정리.
- 미수행: DB 조회·변경·마이그레이션, 앱 빌드·타입 검사·린트·테스트, 브라우저 화면 시험,
  앱/컨테이너 기동·재기동·재생성, OAuth 로그인·모델·MCP·Splunk 호출.
- 실제 인증 파일을 생성하거나 기존 Codex/ChatGPT 인증 파일을 읽지 않았다. 공개 문서 열람 외
  외부 인증·추론 호출은 하지 않았다. 기존 실행 중 컨테이너의 서식 도구만 사용했다.
- 현재 백엔드에는 새 코드를 적용하지 않았다. 8789 포트는 Compose 파일에만 추가되어 있다.
  이후 허용된 백엔드 빌드와 컨테이너 재생성이 필요하다. 단순 restart로 포트는 바뀌지 않는다.
- 로컬 개발 runtime 볼륨은 인증 파일과 기존 서명 파일을 보존한다. 기본 배포 Compose의
  data 바인드만으로는 두 파일을 보존하지 않는다. 다른 배포 방식은 별도 런타임 보존이 필요하다.
- 향후 사용자가 연결·저장·해제 버튼을 누르면 로컬 인증 파일/기존 프로파일에 정상적인
  설정 저장이 발생한다. 이번 소스 작업에서는 해당 API를 실행하지 않았다.
- [후속 검증 목록](../todo/2026-10-23-chatgpt-oauth.md)의 실제 로그인·갱신·한도·해제 및
  5개 언어 화면 검증이 남아 있다. 소스 검토를 동작 성공으로 간주하지 않는다.

## 공식 근거

- [등록과 로그인](https://developers.openai.com/siwc/token-sharing-open-source/sign-in)
- [계정과 세션](https://developers.openai.com/siwc/token-sharing-open-source/profiles-and-sessions)
- [모델과 추론](https://developers.openai.com/siwc/token-sharing-open-source/models-and-inference)
- [사용자 경험 지침](https://developers.openai.com/siwc/ui-ux-guidelines)
