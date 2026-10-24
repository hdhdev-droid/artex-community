# 온보딩·5개 언어·OpenRouter OAuth 코드 반영

> 후속 승인으로 로컬 빌드·테스트·컨테이너 재생성과 5개 언어 화면 확인을 완료했다. [개발 배포 검증 기록](2026-10-23-local-onboarding-verification.md)에 결과를 기록했다. 아래 미수행 표기는 최초 구현 시점 기준이며 실제 OAuth 로그인은 사용자가 진행한다.

- 일자: 2026-10-23
- 상태: 온보딩 및 OpenRouter 코드 반영, 빌드·테스트·실계정 연결 미검증
- ChatGPT 구독 OAuth: 후속으로 콜백 예외 승인 및 코드 반영. [후속 완료 기록](2026-10-23-chatgpt-oauth.md) 참고.

## 반영 내용

- 최초 설정에서 표시 언어 선택, 관리자 설정 성공 후 `/onboarding` 이동.
- 기존 작업 목록의 빈 상태 및 시스템 설정에 시작 안내 진입점 추가. 기존 사용자는 강제로 이동시키지 않는다.
- AI 작업 / API·소스 진단 / 방어 검증·침해 조사별 필수·선택 연결 안내.
- 설정됨·미검증, 연결 확인됨, 실패, 미설정, 확인 불가를 구분. 조회 실패를 빈 설정으로 취급하지 않는다.
- 모델, 진단 서비스, MCP, Splunk 확인은 기존 기능을 재사용하고 버튼을 누를 때만 수행한다.
- 재방문·브라우저 포커스 복귀·새로고침에서는 설정만 조회하고 과거 연결 확인 결과를 지운다. 시험 결과는 DB나 브라우저 영구 저장소에 보관하지 않는다.
- MCP 요약 API는 이름·ID·활성 여부만 반환한다. 환경 변수, 실행 인수, URL을 포함하는 기존 전체 설정 API는 체크리스트에서 호출하지 않는다.
- 언어와 안내 목적만 브라우저에서 임시 보관한다. 서버 언어는 로그인 후 사용자가 저장할 때만 기존 설정 API로 반영한다.
- 한국어·영어·일본어·중국어·스페인어 신규 문구 추가. 사용자 입력·증거 원문은 변경하지 않는다.
- OpenRouter OAuth PKCE S256, 5분 만료, 세션에 결속된 단일 사용 state, 고정 토큰 교환 주소와 리다이렉트 차단 추가.
- 토큰 교환은 백엔드에서 처리하고 브라우저에는 프로파일 ID만 반환한다. upstream 오류 본문이나 키 원문을 로그·오류 응답에 싣지 않는다.
- 완료 시 기존 저장 경로로 새 모델 프로파일을 생성한다. 기존 프로파일 및 활성 모델은 유지하고 자동 대체 목록에서도 제외한다. 연결 확인·활성화는 사용자가 진행한다.
- 기존 초기 관리자 설정 보호와 JWT 인증 미들웨어를 유지한다. 새 API는 관리자 인증을 요구한다.

## 수정 파일

| 파일 | 변경 |
| --- | --- |
| `web/src/app/(main)/onboarding/page.tsx` | 시작 안내, 언어 저장, 목적별 상태와 연결 확인, 기존 설정 화면 이동 |
| `web/src/components/setup-language.tsx` | 인증 전 표시 언어 선택 |
| `web/src/components/onboarding-entry.tsx` | 진입 버튼과 고정 내부 경로 복귀 버튼 |
| `web/src/app/(auth)/setup/page.tsx` | 언어 선택 삽입, 설정 성공 후 시작 안내로 이동 |
| `web/src/app/(main)/_components/main-content.tsx` | 설정 화면에서 시작 안내 복귀 동선 |
| `web/src/app/(main)/function/tasks/page.tsx` | 빈 목록에 시작 안내 진입점 |
| `web/src/app/(main)/system/settings/page.tsx` | 기존 설치용 시작 안내 진입점 |
| `web/src/lib/api.ts` | 비밀 값 없는 MCP 요약 조회 |
| `server/onboarding.go` | 인증된 MCP 메타데이터 조회 |
| `server/llm_oauth.go` | OpenRouter OAuth 시작·교환·프로파일 저장 |
| `server/server.go` | 새 인증 API 라우트와 프로세스 내 만료 상태 |
| `web/src/app/(main)/system/llm/_components/openrouter-connect.tsx` | OpenRouter 계정 연결 입력·시작 UI |
| `web/src/app/(main)/system/llm/oauth/callback/page.tsx` | 코드·state 처리, URL에서 쿼리 제거, 결과 안내 |
| `web/src/app/(main)/system/llm/oauth/callback/layout.tsx` | no-referrer 정책 |
| `web/src/app/(main)/system/llm/page.tsx` | OAuth 연결 UI 배치 |
| `web/messages/{en,ko,ja,zh,es}.json` | `onboarding`, `modelOAuth` 번역 |
| `docs/local-development.md`, `README.ko.md` | 최초 사용 및 모델 연결 안내 |
| `docs/todo/2026-10-23-first-run-onboarding.md` | 계획 상태와 완료 기록 연결 |
| `docs/todo/2026-10-23-chatgpt-oauth.md` | ChatGPT 지원 근거·콜백 예외·후속 구현 범위 |

## 수행 및 미수행

- 수행: 공식 인증 문서 확인, 기존 호출 경로와 변경 코드 정적 검토, 신규 Go/TSX 코드 서식 정리.
- 미수행: DB 조회·변경·스키마 적용, 앱 빌드, 타입 검사, 린트, 자동 테스트, 실제 화면 시험, OAuth 로그인, 모델·MCP·Splunk 연결 시험, 진단 실행.
- 앱 프로세스와 컨테이너를 재시작하지 않았다. 기존 컨테이너의 gofmt/Biome만 서식 정리에 사용했다.
- 새 백엔드 API를 사용하려면 이후 승인된 빌드·재기동이 필요하다. 현재 실행 중인 백엔드에서 새 OAuth 기능이 동작한다고 확인한 것은 아니다.
- 새 테이블·DB 스키마 변경·샘플·데모를 추가하지 않았다. 사용자가 이후 UI에서 언어 저장이나 OAuth 완료를 실행하면 기존 API/테이블에 정상 설정 저장이 발생한다.

## 남은 검증

사용자 허용 후 5개 언어·모바일·신규 설치·기존 설치 진입, 조회 실패, PKCE 만료·취소·잘못된 state·세션 불일치·중복 콜백, 실제 OpenRouter 인증 및 모델 연결을 확인해야 한다.
성공한 코드 교환 뒤 DB 저장에 실패하면 발급된 OpenRouter 키가 남을 수 있으므로 공급자에서 연결 키를 확인하는 안내를 표시한다.
실제 공급자 승인 화면·과금·계정 연결은 이번 작업에서 조작하지 않았다.

## 참고 문서

- [OpenRouter OAuth PKCE](https://openrouter.ai/docs/guides/overview/auth/oauth)
- [ChatGPT 오픈소스 앱 등록](https://developers.openai.com/siwc/token-sharing-open-source/sign-in)
