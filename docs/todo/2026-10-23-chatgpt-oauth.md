# ChatGPT 구독 OAuth 후속 조건

> 후속 승인에 따라 로컬 빌드·테스트·컨테이너 재생성 및 5개 언어 화면 확인을 완료했다. [수행한 범위와 남은 실계정 검증](../complete/2026-10-23-local-onboarding-verification.md)을 기준으로 확인한다. 아래는 최초 계획·구현 시점의 기록이다.

- 일자: 2026-10-23
- 상태: 콜백 예외 사용자 승인, 코드 반영 / 빌드·테스트·실계정 검증 대기

## 지원 방식과 승인 범위

OpenAI의 오픈소스 앱용 Sign in with ChatGPT는 사용자 동의를 받아 해당 계정의
ChatGPT 플랜으로 허용된 Responses API 요청을 수행하는 흐름을 제공한다.
일반 OpenAI API 키나 ARTEX 관리자 로그인과는 별도의 모델 인증이다.

공식 등록 규격은 `http://127.0.0.1:<port>/auth/callback`의 HTTP loopback 콜백을 요구하며
localhost로 대체하지 않도록 명시한다. 현재 사용자 규칙 11번은 localhost 사용을 요구한다.
사용자가 “인증 콜백만 127.0.0.1 허용”으로 승인했다. 앱 접속은 localhost를 유지하고
콜백만 `http://127.0.0.1:8789/auth/callback`을 사용하도록 코드를 반영했다.
DB 변경·빌드·테스트·재기동이나 실제 계정 연결은 이번 승인에 포함되지 않는다.

## 반영 범위

1. 호스트 식별자와 계정별 등록 정보를 보호된 런타임 파일에 보관한다. DB 스키마는 추가하지 않는다.
2. PKCE S256, state, nonce를 새로 생성하고 콜백 리스너를 준비한 뒤 브라우저 인증을 시작한다.
3. 최초 등록에서 발급받은 client ID를 사용해 코드를 교환한다. 재인증은 저장된 계정별 client ID를 재사용한다.
4. 공개 JWKS로 ID 토큰 서명, issuer, audience, 만료, nonce, 계정 일치를 검증한다. 실제 발급 scope로 구독 사용 권한을 판정한다.
5. 토큰은 브라우저 저장소나 일반 프로파일 키 문자열에 복제하지 않는다. 보호된 런타임 저장소의 참조를 기존 프로파일과 연결한다.
6. 계정별 모델 목록을 명시적 동작으로 가져온다. 기존 Responses 지원에 `store=false`, `stream=true`, 갱신 토큰 처리와 완료 이벤트 검증을 연결한다.
7. 토큰 갱신은 계정별로 직렬화하고 회전된 토큰을 원자적으로 보관한다. 연결 해제 시 원격 폐기를 시도하며 실패 여부를 사용자에게 구분해 알린다.
8. UI에 사용 중인 계정, 구독 사용 안내, 사용량 관리, 재인증·연결 해제를 5개 언어로 제공한다. ARTEX 관리자 인증은 유지한다.
9. Docker에서는 호스트 브라우저가 콜백 리스너에 접근하도록 loopback 포트 게시가 필요하다. 현재 컨테이너 변경·재시작은 수행하지 않았다.

기존 Codex/ChatGPT 앱의 인증 파일은 읽거나 재사용하지 않는다. 외부 인증·추론 호출과
실제 계정 연결 검증은 수행하지 않았다.

상세 변경 파일과 실행 조건은 [구현 완료 기록](../complete/2026-10-23-chatgpt-oauth.md)에 기록했다.
이 문서의 구현 항목은 소스에 반영했으며, 아래 동작 검증은 아직 남아 있다.

## 허용 후 검증할 항목

- 빌드·타입 검사, 기존 API 키/OpenRouter 프로파일 회귀, 5개 언어 화면과 키보드·모바일 동선.
- state/nonce/PKCE, 만료·취소·중복 콜백·다른 ARTEX 세션, 다른 client ID·subject·audience·서명 거부.
- 최초 등록 실패 뒤 발급 client ID 재사용, 계정·워크스페이스 분리, 권한 없는 계정의 추론 차단.
- 갱신 토큰 회전·동시 요청·저장 실패·재인증과 해제 경쟁, 재시작 후 암호화 저장소 읽기.
- 명시적 모델 조회, 프로파일 저장·활성화·모델 시험, 도구 호출·추론 스트림의 정상 완료와 중도 종료.
- 구독 한도 오류 시 자동 재시도·다른 공급자로 자동 전환하지 않는지, 사용량 관리 안내.
- 원격 토큰 폐기 성공·실패와 로컬 해제, 기존 진행 중 요청 처리, 콜백 리스너 만료 후 종료.

실계정 인증·모델 호출은 외부 연결과 플랜 사용량을 수반하므로 현재 실행하지 않았다.

## 공식 근거

- [오픈소스 앱 개요](https://developers.openai.com/siwc/token-sharing-open-source)
- [등록·로그인 및 loopback 콜백 규격](https://developers.openai.com/siwc/token-sharing-open-source/sign-in)
- [계정·갱신·로그아웃](https://developers.openai.com/siwc/token-sharing-open-source/profiles-and-sessions)
- [모델 및 추론](https://developers.openai.com/siwc/token-sharing-open-source/models-and-inference)
- [사용자 경험 지침](https://developers.openai.com/siwc/ui-ux-guidelines)
