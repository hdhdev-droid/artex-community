# 로컬 Docker 개발 환경

Windows, Linux, macOS에서 같은 Compose 구성을 사용한다.

| 호스트 | 준비 |
| --- | --- |
| Windows | Docker Desktop, Linux 컨테이너 모드 |
| Linux | Docker Engine과 Compose 플러그인, 현재 사용자의 Docker 실행 권한 |
| macOS (Intel / Apple Silicon) | 해당 Mac 아키텍처용 Docker Desktop 실행 |

프로젝트 루트에서 Windows는 다음 명령으로 시작한다.

```powershell
.\scripts\Start-Dev.ps1
```

Linux와 macOS는 다음 명령으로 시작한다. 실행 권한 변경 없이 사용할 수 있다.

```sh
sh scripts/start-dev.sh
```

두 시작 스크립트는 없는 로컬 인증 파일만 생성하고 `docker-compose.dev.yml`을 실행한다.
기존 인증 파일은 덮어쓰지 않는다. Go 백엔드 컴파일, 프런트엔드 의존성 설치와
Next 개발 서버, 진단 도구 이미지 준비가 포함된다. 테스트나 프런트엔드 배포 빌드는 실행하지 않는다.
첫 기동은 다운로드·컴파일로 시간이 걸리며 컨테이너 실행 상태만으로 앱 준비 완료를 뜻하지 않는다.

Linux/macOS 스크립트는 Docker 연결을 먼저 확인하고, 새 인증 디렉터리·파일을 소유자만
접근하도록 생성한다. 기존 인증 파일의 내용이나 권한은 변경하지 않는다.
호스트에 Go·Node·Python을 설치할 필요는 없다.

Apple Silicon Mac에서 진단 서비스는 `linux/amd64`로 실행한다.
Docker Desktop의 x86_64 에뮬레이션이 필요하며 네이티브 실행보다 느릴 수 있다.
나머지 서비스는 이미지가 제공하는 호스트 아키텍처를 사용한다.
컨테이너는 모두 Linux 환경이며, Windows/macOS 호스트 실행은 네이티브 바이너리 실행을 뜻하지 않는다.
Windows Docker Desktop에서 기본 기동을 확인했으며 Linux·Intel Mac·Apple Silicon Mac의
실제 기동 검증은 아직 수행하지 않았다.

| 구성 | 접속/역할 |
| --- | --- |
| 화면 | http://localhost:5173 |
| 백엔드 | http://localhost:8787/api/health |
| 기록 프록시 | localhost:8788 |
| PostgreSQL 16 | localhost:5432, DB/user `artex` |
| 진단 서비스 | Docker 내부 `offensive-runner:8090`, 호스트에 공개하지 않음 |

2026-10-23 현재 최신 소스로 재기동한 개발 버전은 `1.0.0-dev`다. `v1.0.0` 정식 발행과 구분한다.
`web` 서비스는 `NEXT_PUBLIC_MOCK=0`을 명시하여 실제 API·로컬 DB만 사용한다.
이전에 사용한 `localhost:5183` 샘플 서버는 종료했으며, 개발 테스트는 `localhost:5173`에서 진행한다.
작업을 생성하면 로컬 DB에 저장되고 재시작해도 유지된다. 데이터 초기화·리셋은 하지 않는다.

호스트 포트는 루프백에만 바인딩한다. 기존 Redis는 연결하거나 변경하지 않는다.
앱에 Redis 의존성이 없으므로 새 Redis를 추가하지 않는다.

## 최초 화면 설정

`http://localhost:5173/setup`에서 관리자 비밀번호를 직접 설정한다.
설정 토큰은 로컬 `.secrets/dev/setup-token` 파일에서 확인한다.
DB 비밀번호는 `.secrets/dev/postgres-password`, 서비스 간 토큰은
`.secrets/dev/offensive-token`에 보관한다. 파일 내용은 공유하거나 커밋하지 않는다.
JWT 서명 파일은 Docker `runtime` 볼륨에서 앱이 생성하며 기존 설치에서 가져오지 않는다.

관리자 설정이 끝나면 `/onboarding` 시작 안내로 이동한다. 기존 설치에서는 작업이 없는
목록의 **시작 안내** 또는 **시스템 설정 → 시작 안내**로 열 수 있다. 언어 선택은 최초
설정 화면에서는 브라우저 표시에만 적용되며, 로그인 후 시작 안내에서 **언어 저장**을
눌러 설치 전체의 화면·생성 결과 언어에 적용한다.

AI 작업, API·소스 진단, 방어 검증·침해 조사 중 사용 목적을 고르면 필요한 설정과
선택 항목을 안내한다. 설정값 존재와 실제 연결 확인을 구분한다. **연결 확인** 버튼은
실제 외부 요청, 모델 사용 비용 또는 MCP 프로세스 실행을 수반할 수 있다.
샘플 데이터·진단 작업은 자동으로 생성하거나 실행하지 않는다.

**LLM → OpenRouter**에서 프로파일 이름과 모델 ID를 입력하고 OAuth 연결을 시작할 수 있다.
공급자 인증을 마치면 기존 테이블에 새 프로파일을 저장한다. 기존 활성 모델은 바꾸지 않는다.
프로파일을 직접 선택해 연결 확인 후 활성화한다. 사용량은 OpenRouter 크레딧에서 차감된다.
키 교환은 서버에서 수행하며 콜백 화면으로 키 원문을 반환하지 않는다.

ChatGPT 구독 연결은 다음 순서로 진행한다. 현재 영어 화면의 버튼명도 함께 표기한다.

1. 같은 컴퓨터의 브라우저에서 `http://localhost:5173/system/llm`을 열고 ARTEX에 로그인한다.
2. **ChatGPT** 영역의 **계정 또는 워크스페이스 추가 / Add an account or workspace**를 선택하고 **ChatGPT로 계속 / Continue with ChatGPT**를 누른다.
3. OpenAI 로그인 화면에서 구독을 사용하는 계정·워크스페이스로 로그인하고 요청 권한을 확인하여 동의한다. ARTEX로 돌아와 **ChatGPT 연결 완료 / ChatGPT connected**가 표시되면 **모델 설정 / Model settings**을 누른다.
4. 연결된 계정이 선택되어 있는지 확인한 뒤 **사용 가능한 모델 가져오기 / Load available models**를 누른다. 반환된 목록에서 모델을 선택하고 프로파일 이름(예: `ChatGPT Local`)을 입력한다.
5. **프로파일 저장 / Save profile**을 누른다. 아래 **모델 설정 / Model profiles** 목록에 새 카드가 생긴다. 저장만으로 기본 모델은 변경되지 않는다.
6. 새 카드의 **활성으로 설정 / Set active**를 누른다. 버튼이 **활성됨 / Active**으로 바뀌고 별표가 노란색이면 기본 프로파일 활성화가 완료된 것이다.
7. 실제 요청을 시험하려면 카드 본문을 눌러 편집 창을 열고 **연결 테스트 / Test connection**을 누른다. 성공 메시지와 모델 응답을 확인한다. 이 단계는 실제 모델 요청이며 ChatGPT 플랜 사용량을 소비할 수 있다.

인증과 연결 주소는 ChatGPT 영역에서 관리한다. API 키나 주소를 별도로 입력할 필요가 없다.
기존 작업·에이전트가 특정 프로파일을 별도로 지정했다면 기본 모델 활성화와 별개로 해당 선택도 확인한다.
모델 목록이 비어 있거나 권한 오류가 표시되면 연결한 계정·워크스페이스와 플랜 사용 권한을 확인한다.
모델명은 임의로 입력하지 않고 [해당 계정이 반환한 모델 목록](https://developers.openai.com/siwc/token-sharing-open-source/models-and-inference)에서 선택한다.
허용된 모델 요청에는 연결된 ChatGPT 플랜 또는 크레딧이 사용되며, 사용량 관리 링크에서
계정과 앱 한도를 확인한다. 구독 종류에 따라 사용할 수 있는 모델·권한이 다르다.

ChatGPT 인증은 ARTEX와 브라우저가 같은 컴퓨터에 있는 로컬 설치용이다.
앱 접속은 `http://localhost:5173`, 인증 콜백만 `http://127.0.0.1:8789/auth/callback`이다.
콜백 리스너는 로그인 시작 때만 열고 만료된 요청을 정리한 뒤 닫는다. 8789 포트가 이미
사용 중이면 로그인을 시작하지 않고 오류를 표시한다. 현재 구현은 포트 8789를 고정 사용한다.
Docker의 `ARTEX_CHATGPT_CALLBACK_LISTEN=0.0.0.0:8789`는 컨테이너 내부 수신 설정이며,
호스트에는 `127.0.0.1:8789:8789`로만 게시한다. 원격 브라우저·VM 콜백 전달은 구현하지 않았다.

토큰은 실행 파일 옆 `chatgpt-oauth/credentials.enc`에 암호화해 보관하며, 기존 프로파일에는
계정 참조만 저장한다. 탐색 가능한 `data/` 워크스페이스나 브라우저 저장소에는 토큰을 넣지 않는다.
로컬 개발 Compose의 `runtime:/app` 볼륨이 이 파일과 기존 JWT 서명 파일을 함께 보존한다.
다른 배포 방식에서 컨테이너·실행 디렉터리를 교체할 때는 두 런타임 파일의 보존이 필요하다.
운영 Compose는 `auth:/app/auth` 바인드와 `ARTEX_AUTH_DIR`로 인증 파일을 함께 보존한다.
기존 운영 컨테이너는 삭제·교체 전에 [이전 절차](deployment-auth.md)를 실행해야 한다.
암호화 키가 바뀌면 기존 파일을 임의로 초기화하지 않고 저장소 오류를 알린다.
같은 설치를 이전할 때만 인증 디렉터리를 함께 옮기며, 여러 활성 인스턴스에 복제해 공유하지 않는다.
실제 키 값은 문서에 기록하지 않는다.

**계정 연결 해제**는 새 요청의 인증을 차단하고 원격 토큰 폐기를 시도한 뒤 로컬 토큰을 제거한다.
원격 폐기 확인에 실패하면 별도 안내를 표시한다. 프로파일·등록 ID는 보존되며 다시 연결할 수 있다.
이미 외부로 전송된 요청을 되돌리거나 기존 작업을 자동으로 취소하지는 않는다.

2026-10-23 후속 승인으로 온보딩·OpenRouter·ChatGPT 코드를 로컬에서 빌드·테스트하고
개발 컨테이너를 재생성했다. 현재 백엔드는 새 OAuth API와 콜백 포트 구성을 제공한다.
이후 같은 날 최신 로컬 코드로 backend/runner 이미지를 다시 빌드하고 backend/web/runner를
재생성했으며, PostgreSQL과 기존 볼륨은 유지했다. 실제 OAuth 로그인과 모델 호출은 사용자가 진행한다.
환경·포트 변경 적용에는 컨테이너 재생성이 필요하며 `restart`만으로는 포트가 추가되지 않는다.
정확한 검증 범위는 [개발 배포 완료 기록](complete/2026-10-23-local-onboarding-verification.md)을 참고한다.

## 데이터와 변경 반영

- Compose 프로젝트는 `artex-community-dev`다. 기존 운영용 Compose와 합쳐 실행하지 않는다.
- PostgreSQL, 실행 데이터, 런타임 인증 파일은 각각 별도 이름의 Docker 볼륨에 보관한다.
- 새 빈 DB에 앱이 기존 기본 스키마와 내장 항목을 자동 적용한다. 별도 init/reset 명령은 사용하지 않는다.
- 기존 프로젝트의 `data/`, `config.json`, `.env`를 앱 설정/데이터로 사용하지 않는다.
- 2026-10-20 후속 반영에서 `offensive_candidate_reviews`와 방어·조사·Closed-loop 저장 테이블을 추가했다.
  사용자 승인 후 로컬 개발 DB를 백업하고 앱의 기존 추가형 스키마 적용 경로로 반영했다.
- 프런트엔드 소스는 개발 서버에 연결되어 수정이 반영된다. `.next`와 `node_modules`는 Docker 볼륨에 둔다.
- Go 코드 변경 후 아래 명령으로 백엔드를 다시 컴파일하고 기동한다.

```powershell
docker compose -f docker-compose.dev.yml restart backend
```

진단 Python 코드나 규칙 변경 후:

```powershell
docker compose -f docker-compose.dev.yml up -d --build offensive-runner
```

상태 확인 및 데이터 보존 중지:

```powershell
docker compose -f docker-compose.dev.yml ps
docker compose -f docker-compose.dev.yml stop
```

볼륨 삭제 옵션은 사용하지 않는다. 재시작은 Windows에서 `Start-Dev.ps1`,
Linux/macOS에서 `sh scripts/start-dev.sh`로 진행한다.
프런트엔드 패키지를 바꾼 경우 `web`을 재시작하면 lockfile 기준으로 의존성을 다시 설치한다.

## 후속 기능과 검증

후보 검토 저장·이어하기, 방어 검증, 로그 재생, 독립 조사 사건, 조사 에이전트,
Splunk 연결, Closed-loop 평가·실행 이력과 보관·복원 연결을 반영했다.
화면은 `/offensive/assessments`, `/function/purple`, `/defensive`,
`/function/purple/sources`, `/closed-loop`에서 사용한다. 기존 대시보드는 유지한다.

서버·DB·에이전트·진단 서비스 테스트, 프런트엔드 테스트·타입 검사·배포 빌드를 진행했다.
정확한 검증 범위와 제한은 [후속 반영 완료 기록](complete/2026-10-20-security-followup.md)을 참고한다.
실제 외부 모델·Splunk·공격 대상 호출 및 AWS 구축은 수행하지 않았다.
AWS 실습 기능은 후속 요청으로 제거했다. 해당 화면·API·브리지 설정은 제공하지 않는다.

격리된 백엔드 테스트:

```powershell
.\scripts\Test-Dev.ps1
```

이 명령은 휘발성 `postgres-test` 컨테이너만 새로 만들며, DB·에이전트·서버 테스트의
DB를 각각 분리한다. 앱의 `postgres` 서비스나 데이터 볼륨은 변경하지 않는다.
테스트 DB의 덤프·복원과 테이블별 행 내용 지문 대조도 함께 실행한다.
실제 데이터 이전 절차는 [이전 체크리스트](migration-checklist.md)를 참고한다.
상세 테스트 로그는 백엔드 컨테이너의 `/tmp/artex-backend-check.log`에 있으며
설정 토큰이 포함될 수 있어 공유하지 않는다. 콘솔에는 결과 요약만 표시한다.

프런트엔드 배포 빌드는 개발 서버의 캐시와 분리된 컨테이너에서 실행한다.

```powershell
docker compose -f docker-compose.dev.yml --profile test run --rm --no-deps web-check
```

모델 인증과 Splunk 접속은 사용자가 설정한 연결을 사용한다.
화면과 API·모델 출력은 서버의 공통 `output_language`를 기준으로 한다.
화면 언어는 저장 성공 뒤 반영하며, 로그인된 화면에 들어오거나 브라우저로 돌아오면
서버 설정과 동기화한다. 기본값은 영어이며 기존 설치의 저장된 언어 설정을 우선한다.
원본 증거·사용자 입력·모델용 도구 설명·과거 보고서 원문은 자동으로 다시 쓰지 않는다.

`Dockerfile.dev`로 Go·Python·Node·nmap·Playwright CLI/MCP·Chromium을 포함한
개발용 백엔드 이미지를 빌드한다. 도구 설치에는 패키지 저장소와 브라우저 배포 서버 접근이 필요하다.
실제 대상 진단은 별도 승인과 설정이 필요하다.

JWT와 ChatGPT 인증 파일은 개발용 `runtime` 볼륨의 `/app`에 보존된다.
백엔드 재기동·재생성 시 이 볼륨을 유지해야 하며 샘플 데이터는 자동으로 넣지 않는다.
운영 Compose는 별도 `auth` 디렉터리를 사용한다. 기존 운영 컨테이너는
[인증 정보 이전 절차](deployment-auth.md)를 먼저 수행한다.

Windows/macOS Docker 마운트에서도 파일 변경이 감지되도록 개발 프런트는
webpack과 Watchpack polling을 사용한다. 배포 빌드는 기존 Next 빌드 방식을 유지한다.
