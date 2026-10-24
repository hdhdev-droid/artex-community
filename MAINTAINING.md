# 상류 동기화와 번역 드리프트 방지 (메인테이너 안내)

한국어 · [English](MAINTAINING.en.md)

<a id="community-upstream"></a>

## artex-community의 직접 상류와 업데이트 기준

이 저장소 `hdhdev-droid/artex-community`는 `jiwoochris/artex-global`에서 포크했습니다.
직접 업데이트를 추적할 상류는 [artex-global](https://github.com/jiwoochris/artex-global)이며,
원본 프로젝트 계보는 `Autumn-27/ARTEX → artex-global → artex-community`입니다.

- 확인일: **2026-10-20**.
- 최초 기록한 상류 비교 기준: **`fdfb7a17b7da5e2de716c0e98475ec77a41c64c6`**.
- 기록 당시 두 로컬 프로젝트의 `HEAD`가 위 커밋으로 같음을 확인했습니다. 원격의 최신 상태를 조회했다는 의미는 아닙니다.
- 커뮤니티 저장소(`origin`): `https://github.com/hdhdev-droid/artex-community`.
- 직접 상류: `https://github.com/jiwoochris/artex-global`.

**이 절이 커뮤니티판의 상류 추적 기준입니다.** 아래 기존 절은 글로벌판에서 이어받은
원본 ARTEX 동기화·현지화 안내입니다. 아래의 `d003372`, `upstream` 리모트,
글로벌판 CI·릴리스 저장소 예시를 커뮤니티판의 직접 상류 기준으로 사용하지 않습니다.

### 개선 사항·신규 버전 반영 절차

1. 작업 트리와 기존 리모트를 확인합니다. 직접 상류용 이름은 `global-upstream`으로 구분하고,
   이미 존재하면 주소가 맞는지 확인합니다. 없을 때만 아래 명령으로 등록합니다.

   ```bash
   git remote add global-upstream https://github.com/jiwoochris/artex-global.git
   ```

2. 업데이트 작업 시 상류를 조회하고, 위에 기록한 기준 커밋 이후의 변경과 릴리스 안내를 확인합니다.

   ```bash
   git fetch global-upstream
   git log --oneline fdfb7a17b7da5e2de716c0e98475ec77a41c64c6..global-upstream/main
   git diff --stat fdfb7a17b7da5e2de716c0e98475ec77a41c64c6 global-upstream/main
   ```

3. 보안 수정·버그 수정·기능 개선·의존성·DB 스키마·설정 변경을 구분하여 별도 작업 브랜치에서 반영합니다.
   디렉터리를 통째로 덮어쓰지 않고 커뮤니티 고유 기능, 인증 강화, 다국어 지원을 보존합니다.
   DB 작업은 조회만 하며, 스키마 적용·초기화·리셋을 실행하지 않습니다.
4. 빌드·테스트는 자동으로 실행하지 않고, 필요한 검증 항목을 사용자에게 안내하여 요청 시 진행합니다.
   상류의 CI 결과와 커뮤니티판의 검증 결과를 구분합니다.
5. 반영한 상류 커밋·버전, 제외하거나 보류한 변경, 실제 검증 여부를 당일
   `docs/complete/YYYY-MM-DD-작업명.md`에 기록합니다. 연속된 상류 변경을 모두 반영한 경우에만
   위 비교 기준과 명령 예시의 커밋을 함께 갱신합니다. 일부 커밋만 선별했다면 기준은 유지하고 적용 목록을 기록합니다.

이 문서화는 자동 갱신이나 정기 모니터링을 설정하지 않습니다.

### community 배포 버전 관리

상류 코드 검토 대상은 `jiwoochris/artex-global`로 유지하며, **설치·자동 업데이트 대상은 `hdhdev-droid/artex-community`의 정식 릴리스만** 사용합니다.
업데이트 원본은 `selfupdate/github.go`의 `Repo` 상수이고 사용자 설정으로 변경하지 않습니다.

- 현재 community 기준 버전은 `version/VERSION`의 **`1.0.0`**입니다. 발행 태그는 반드시 이 파일과 같은 `vMAJOR.MINOR.PATCH`여야 합니다. 수정은 PATCH, 호환 기능 추가는 MINOR, 호환성이 깨지는 변경은 MAJOR를 올립니다. `v1.0.0`은 발행 준비된 기준이며 아직 태그를 발행하지 않았습니다.
- `version/VERSION`과 `version/upstream.json`을 빌드·화면·서버가 공유합니다. 태그와 VERSION이 다르면 CI가 중단됩니다. `web/package.json`의 버전은 프런트엔드 패키지 메타데이터이며 앱 업데이트 판정에 사용하지 않습니다.
- `.github/workflows/release.yml`은 community 저장소에서만 발행하며, `v*` 태그 중 세 자리 정식 버전 형식이 아닌 태그(예: `v1.0.0-rc.1`)는 거부합니다. 프런트엔드·5개 플랫폼 패키지·Docker 이미지 빌드 후 ZIP과 `SHA256SUMS`를 GitHub Release에 게시합니다.
- 태그에서 `v`를 뺀 값이 실행 파일의 `main.version`에 들어갑니다. CI에서 Linux 실행 파일의 `-version` 출력과 태그 일치 여부를 확인합니다. `-version`은 DB·업데이트 부트스트랩·서버 시작 전에 종료합니다.
- 이미지 주소는 `ghcr.io/hdhdev-droid/artex-community`입니다. Docker의 `ARTEX_TAG=vX.Y.Z`로 배포 버전을 고정할 수 있으며 기본값은 `latest`입니다. 최초 community 릴리스/이미지가 실제 발행되기 전에는 내려받을 수 없습니다.
- `build.sh`와 `build-cloud.sh`는 공통 `scripts/build-version.sh`를 사용합니다. 기본 개발 버전은 `1.0.0-dev`이며, CI가 명시하는 정식 `ARTEX_BUILD_VERSION`은 VERSION 파일과 일치해야 합니다. 개발 산출물에 정식 버전을 임의 지정하지 않습니다. 직접 `go build`한 바이너리도 VERSION에 `-dev`를 붙입니다.
- 기존 설치에 새 업데이트 저장소 설정을 반영하려면 먼저 이번 community 코드를 빌드해 수동 배포해야 합니다. 이전 실행 파일은 여전히 이전 업데이트 주소를 갖고 있습니다.

발행 전에는 기존 community 릴리스와 태그를 확인해 더 높은 버전을 정하고, 검증한 커밋의 완료 기록·변경 사항·DB 호환성을 확인합니다. 새 태그를 `origin`에 push하면 발행 워크플로가 실행되므로 승인된 배포 시에만 수행합니다. 이미 발행한 태그/파일을 덮어쓰지 않고 수정 버전을 새로 발행합니다. 상류 태그를 community 릴리스로 일괄 push하지 않습니다.

앱은 community의 최신 정식 릴리스와 현재 실행 파일의 숫자 버전을 비교합니다. 더 높은 버전과 해당 OS/CPU 패키지·체크섬이 있어야 설치할 수 있습니다. 프로그램 교체는 DB 롤백을 제공하지 않으며 Docker 컨테이너 내부 업데이트도 이미지 자체를 변경하지 않습니다.
샘플 화면은 `1.0.0-demo`와 community 기준 버전·상류 출처를 표시하며 외부 릴리스 조회·설치·롤백은 실행하지 않습니다. 정확한 포크 기준은 `version/upstream.json`의 `v0.3.15-global.1-9-gfdfb7a1`입니다. 원본 ARTEX `0.3.15` 기반 표시는 글로벌판 변경 기록에 근거합니다.

다음 버전은 VERSION 갱신 → 해당 `docs/releases/vX.Y.Z.md` 작성 → 변경 커밋 → 동일 태그 발행 순서로 관리합니다. 현재 작업 트리에는 미커밋 기능이 있으므로 기존 HEAD에 바로 태그를 붙이면 해당 기능이 릴리스에 포함되지 않습니다.

---

이 문서는 **메인테이너**가 원본 저장소 [Autumn-27/ARTEX](https://github.com/Autumn-27/ARTEX)(현재 접속 불가, 2026-10 확인)의
변경을 따라잡으면서 다국어 현지화를 유지하는 절차를 정리한 것입니다. 기여 범위·법적 책임·
현지화 방침은 [CONTRIBUTING.md](CONTRIBUTING.md)에, 사용자용 안내는 [README.ko.md](README.ko.md)에
있으므로, 이 문서는 그 방침을 **실제로 어떻게 집행하는지**에만 집중합니다.

현지화의 핵심 목표는 한 문장으로 요약됩니다. 원본의 **판단 성능을 그대로 보존하면서 사용자에게
보이는 산출물만 설정한 언어로 바꾸는 것**입니다. 상류가 갱신될 때마다 이 경계가 흐트러지기 쉬우므로,
아래 절차와 검사로 번역 드리프트를 막습니다.

---

## 1. 현지화 구조 한눈에 보기

이 저장소는 상류 ARTEX 를 **포크**해서 그 이력 위에 한국어 현지화 커밋을 쌓은 구조입니다.
상류 `main` 의 모든 커밋이 이 저장소의 이력에 포함되어 있고, 그 위에 현지화 커밋이 더해져
있습니다. 따라서 상류 변경을 가져오는 일은 "상류 `main` 과의 차이를 확인하고, 보존할 것과
번역할 것을 가려서 반영하는 일"이 됩니다.

산출물은 세 갈래로 나뉩니다.

- **원문을 그대로 두는 자산**(아래 2절). 번역하면 성능이나 상류 대조가 깨집니다.
- **코드에 고정된 출력 언어 강제**. `agent/prompt.go` 의 `langDirective()` 가 각 역할의 system
  프롬프트 끝에 "사용자 노출 출력은 한국어로 작성하라"는 지시를 덧붙입니다.
- **한국어로 번역하는 사용자 노출 문자열**. UI 는 `web/messages/ko.json` 에, 서버의
  사용자 응답 문구는 각 Go 파일의 명명 상수에 둡니다.

---

## 2. 원문을 보존하는 자산 (번역 금지)

다음 자산은 번역하지 않고 원문(중국어 또는 영어)을 유지합니다. 상류 변경이 이 자산에 닿으면
**번역 없이 그대로 반영**합니다.

- **에이전트 내부 추론 프롬프트(두뇌 본문).** `agent/promptcatalog.go` 와 DB 시드
  `agent_prompts` 에 있는 행동 지침 본문입니다. 원문(중국어)으로 벤치마크된 동작을 유지해야
  하므로 번역하면 판단에 드리프트가 생깁니다.
- **표시와 에이전트 입력을 겸하는 문자열.** 활동 타임라인에 보이면서 동시에 플래너·리포터의
  입력 컨텍스트로 되먹여지는 일부 문구(작업 중단 사유, 가로채기 차단 메시지, 트래픽 증거 헬퍼
  등)는 하나의 레코드가 두 용도를 겸하므로 원문을 보존합니다. 판정 근거는
  `work/DECISIONS-FOR-JIWOO.md` 의 "두뇌 경계 기록"에 사안별로 적혀 있습니다.
- **원본 중국어 문서·문자열.** 문서는 `README.zh.md`, UI 문자열은 `web/messages/zh.json` 에
  원문을 그대로 남겨 상류 변경과 대조하기 쉽게 합니다. 한국어 번역은 `web/messages/ko.json`
  에만 채웁니다.
- **명령·페이로드·코드·URL·식별자·로그 원문.** 분석에 필요한 원본이므로 번역하지 않습니다.
  Go 코드 주석도 우선순위가 가장 낮아 상류 대조가 끝나는 시점까지 원문을 둡니다.

---

## 3. 상류 추적 베이스

상류 리모트가 다음과 같이 설정되어 있어야 합니다. 없다면 추가합니다.

```bash
git remote add upstream https://github.com/Autumn-27/ARTEX
git remote -v   # upstream 이 보이는지 확인
```

현재 현지화가 반영을 마친 상류 베이스 커밋은 다음과 같습니다.

- **베이스 = `d003372`** (상류 `main`, 2026-10-03, PR #189 `fix/sse-same-origin` 병합).

이 값은 "이 커밋까지의 상류 변경은 전부 이 저장소에 녹아 있다"는 뜻입니다. 상류 변경을 새로
반영할 때마다 이 베이스를 7절의 방법으로 갱신합니다.

---

## 4. 상류 변경을 가져오는 절차

### 4.1 상류를 내려받고 차이를 확인합니다

```bash
git fetch upstream
git rev-list --count d003372..upstream/main        # 미반영 상류 커밋 수
git log --oneline d003372..upstream/main           # 미반영 커밋 목록
```

`git fetch` 는 상류의 원격 추적 브랜치만 갱신하므로 작업 트리와 `HEAD` 에는 영향을 주지
않습니다. 미반영 커밋이 0 이면 상류와 동기화된 상태이므로 더 할 일이 없습니다.

### 4.2 변경 파일을 분류합니다

미반영 커밋이 어떤 파일을 건드렸는지 보고, 2절의 보존 자산과 번역 대상으로 나눕니다.

```bash
git log --name-status --oneline d003372..upstream/main
```

분류 기준은 다음과 같습니다.

- `agent/promptcatalog.go`·`agent_prompts` 시드, 그리고 2절의 표시 겸 입력 문자열이 바뀌었다면
  → **번역 없이 그대로 반영**합니다.
- Go 백엔드 로직(`db/`·`llmrec/`·`server/` 등)이 바뀌었다면 → 로직은 그대로 반영하되,
  **새로 생긴 사용자 응답 문구**(`writeErr` 등)가 있는지 확인해서 한국어 상수로 번역합니다.
- UI(`web/src/**`)가 바뀌어 **새 화면 문자열**이 생겼다면 → 하드코딩하지 말고
  `web/messages/zh.json`(원문)과 `web/messages/ko.json`(번역)에 같은 키로 추가합니다.
- **탐지 규칙이 고정한 상류 지표**(`enrich/enrich.go` 의 프로버 User-Agent, `selfupdate/` 의
  자가 갱신 User-Agent, `guard/guard.go` 의 감사 마커, `db/db.go` 의 파괴명령 deny 목록,
  `cmd/artex/main.go` 의 기본 리슨·기록 프록시 포트)가 바뀌었다면 → `detections/` 의
  Sigma·Suricata 규칙과 ATT&CK 레이어, 그리고 `detections/indicators/artex_indicators.csv` 의
  값도 새 값으로 맞춥니다. 이 지표는 번역 대상이 아니라 **탐지의 근거**라, 상류가 값을 바꾸면
  규칙이 조용히 낡습니다. 5.4 의 지표 일치 테스트가 이 어긋남을 자동으로 잡습니다.

### 4.3 반영합니다

기능 단위로 병합하거나 선별 반영한 뒤, 4.2 에서 가려낸 새 문자열을 한국어로 번역합니다.
병합 과정에서 `ko.json`·`zh.json` 의 키가 어긋나거나 사용자 노출 자리에 원문이 새어 들어오기
쉬우므로, 반영 직후 반드시 5절의 검사를 돌립니다.

> **예시(2026-10-05 기준 미반영 커밋).** `git fetch upstream` 결과 상류 `main` 이
> `b55ceb1` 로 앞서 있고, 베이스 `d003372` 대비 커밋 2개(`86729b6` 모델 폴백 승인 토큰 계량
> 기능 + 병합 커밋 `b55ceb1`)가 미반영입니다. 이 커밋은 `db/llm_usage.go`·`llmrec/llmrec.go`·
> `server/intercept.go`·`server/server.go` 같은 Go 로직과 `web/src/app/(main)/system/intercept/page.tsx`·
> `web/src/lib/api.ts`·`web/src/lib/mock/handler.ts`·`web/src/lib/types.ts` 를 건드립니다.
> 따라서 메인테이너는 Go 로직은 그대로 반영하고, intercept 설정 페이지에 새로 생긴 화면
> 문자열만 `ko.json`·`zh.json` 키로 추출·번역하면 됩니다. (이 두 커밋은 이 문서를 쓴 시점에는
> 아직 반영하지 않았으므로 베이스는 `d003372` 로 둡니다.)

---

## 5. 번역 대칭과 드리프트 검사

상류 반영이나 번역 작업 뒤에 아래 세 가지를 확인합니다.

### 5.1 ko ↔ zh 메시지 대칭과 사용자 노출 CJK

`ko.json` 과 `zh.json` 의 키가 정확히 같고, `ko.json` 값에 중국어 한자가 남아 있지 않아야
합니다. 아래 스크립트가 세 수치를 출력합니다.

```bash
python3 - <<'PY'
import json, re
ko = json.load(open('web/messages/ko.json'))
zh = json.load(open('web/messages/zh.json'))
def flatten(d, p=''):
    out = {}
    if isinstance(d, dict):
        for k, v in d.items(): out.update(flatten(v, p + '/' + k))
    elif isinstance(d, list):
        for i, v in enumerate(d): out.update(flatten(v, p + '/' + str(i)))
    else: out[p] = d
    return out
fk, fz = flatten(ko), flatten(zh)
han = re.compile(r'[㐀-鿿]')
print('ko leaf keys :', len(fk))
print('zh leaf keys :', len(fz))
print('key symdiff  :', len(set(fk) ^ set(fz)))        # 0 이어야 함
print('ko vals w/CJK:', sum(1 for v in fk.values() if isinstance(v, str) and han.search(v)))  # 0 이어야 함
PY
```

기준값(2026-10-05): `ko leaf keys = 2950`, `zh leaf keys = 2950`, `key symdiff = 0`,
`ko vals w/CJK = 0`. 키 수는 상류 반영으로 늘 수 있지만, ko 와 zh 는 항상 같아야 하고
`key symdiff` 와 `ko vals w/CJK` 는 항상 0 이어야 합니다.

### 5.2 두뇌 자산의 원문 보존 확인

두뇌 본문은 중국어 원문을 유지하므로, 아래 검사에서 **한자 라인 수가 0 으로 떨어지면** 오히려
두뇌가 실수로 번역돼 오염됐다는 신호입니다.

```bash
python3 -c "import re; han=re.compile(r'[㐀-鿿]'); t=open('agent/promptcatalog.go').read(); print('promptcatalog.go CJK lines =', sum(1 for l in t.splitlines() if han.search(l)))"
```

기준값(2026-10-05): `promptcatalog.go CJK lines = 70`. 이 수가 크게 줄면 두뇌 본문이 번역됐는지
확인합니다.

### 5.3 빌드 산출물에 원문이 새지 않는지

UI 를 정적으로 내보낸 뒤 프리렌더 HTML 에 중국어가 보이면 번역 누락입니다.

```bash
cd web && npm ci && NEXT_EXPORT=1 npm run build   # out/ 생성
# out/**/*.html 에서 가시 텍스트의 중국어 한자가 0 인지 확인
```

### 5.4 탐지 지표가 상류 소스와 여전히 맞는지

`detections/` 의 규칙은 상류가 실제로 내보내는 문자열(프로버 User-Agent·자가 갱신
User-Agent·감사 마커·파괴명령 deny 목록)에 근거합니다. 상류 재동기화가 이 값을 바꾸면
번역 검사는 전부 통과하는데 배포된 규칙만 조용히 매칭을 멈춥니다. 아래 테스트가 각 지표가
상류 소스와 규칙 양쪽에 여전히 있는지 양방향으로 확인하므로, 재동기화 뒤에 함께 돌립니다.

```bash
detections/tests/indicators/run.sh   # Docker 로 격리 실행, RESULT: PASS 이면 일치
```

실패하면 어느 지표가 어긋났는지와 그 방향(상류 소스가 바뀌었는지, 규칙이 바뀌었는지)을
출력하므로, 4.2 의 마지막 분류 기준대로 규칙·레이어를 새 값에 맞춥니다. 이 테스트는 저장소
CI([`.github/workflows/detections.yml`](.github/workflows/detections.yml))에서도 규칙 트리나
위 상류 소스 파일이 바뀐 푸시·PR 마다 자동으로 돌아, 재동기화 드리프트를 머지 게이트에서 잡습니다.

새 지표를 추가하면서 **새 상류 소스 파일을 고정했다면**(예: `cmd/artex/main.go` 의 포트 지표를
넣을 때처럼), 그 파일을 반드시 위 워크플로의 `push`·`pull_request` `paths` 필터에도 추가합니다.
빠뜨리면 그 소스만 바꾼 PR 은 지표 테스트를 발화시키지 못해, 드리프트가 머지 게이트를 조용히
통과합니다. 이 동기화 자체도 지표 테스트가 자동으로 확인합니다(다섯 번째 검사 "CI triggers this
test when any pinned source changes"): 테스트가 읽는 모든 비 `detections/` 소스가 양쪽 `paths`
블록에 열거돼 있지 않으면 테스트가 실패하므로, 소스 고정과 CI 발화 조건이 어긋난 채로 머지되지
않습니다.

### 5.5 탐지 테스트 도구 핀을 올릴 때

탐지 테스트는 `sigma-cli`·SigmaHQ 검증기 플러그인(`pySigma-validators-sigmahq`)·Suricata 이미지를
고정 버전으로 돌립니다(각 `run.sh` 의 기본값, 환경 변수로 덮어쓰기 가능). 이 핀을 올리면 상류
소스가 아니라 **도구 쪽 드리프트**가 생길 수 있습니다. 특히 SigmaHQ 검증기는 판올림마다 새 관례
검사를 추가하므로, `detections/tests/sigma_lint/run.sh` 가 새 이슈를 빨갛게 드러낼 수 있습니다.
그때는 규칙을 새 관례에 맞추거나, 단독 규칙 세트에 맞지 않는 관례라면 그 사유를 적어
[`detections/tests/sigma_lint/validators.yml`](detections/tests/sigma_lint/validators.yml) 의 제외
목록에 추가합니다. 백엔드 플러그인이 지원을 바꾸면 `sigma_backends` 테스트가 같은 신호를 줍니다.

---

## 6. 빌드와 테스트로 마무리 검증

반영·번역 뒤에는 [CONTRIBUTING.md 의 개발 환경](CONTRIBUTING.md#개발-환경) 절차대로 백엔드와
프런트엔드를 검증합니다. 로컬에 Go 가 없으면 Docker 로 동일하게 돌릴 수 있습니다.

```bash
docker run --rm -v "$PWD":/src -w /src \
  -v artexko-gomod:/go/pkg/mod -v artexko-gocache:/root/.cache/go-build \
  golang:1.26 sh -c 'go build ./... && go vet ./... && go test ./... -count=1'
```

사용자 노출 문구를 번역할 때는 그 문구를 단언하는 회귀 테스트(`*_localized_test.go`)를 함께
두어, 나중에 상류 변경이 다시 중국어를 끌어와도 테스트가 잡게 합니다. 번역 검증은 반드시
역량 있는(프런티어급) 모델로 합니다. 저가·소형 모델은 출력이 원문으로 되돌아갈 수 있어
번역 적용 여부를 그 출력만으로 판단하면 안 됩니다.

---

## 7. 베이스 갱신 기록

상류 변경을 반영하고 검증까지 마쳤다면, 이 문서 3절의 **베이스 커밋 값을 새 상류 커밋으로
갱신**하고 그 변경을 같은 커밋 또는 뒤따르는 커밋에 포함합니다. 이렇게 해두면 다음 메인테이너가
"어디까지 반영됐는지"를 이 문서 한 곳에서 확인할 수 있습니다.

커밋 메시지는 [CONTRIBUTING.md 의 커밋 메시지 규칙](CONTRIBUTING.md#커밋-메시지)을 따릅니다.
예를 들어 상류 동기화 커밋은 다음과 같이 적습니다.

```
chore(upstream): 상류 d003372..b55ceb1 반영 (intercept 토큰 계량) + 신규 UI 문자열 번역
```

---

## 8. 검토·검증 수칙 (흔한 함정)

상류 반영·번역·문서 보강을 점검할 때 메인테이너가 반복해서 빠지는 함정 두 가지를 적어
둡니다. 둘 다 "검사 방법 자체가 틀려서 멀쩡한 것을 깨졌다고 오인하는" 경우라, 불필요한
되돌림을 막으려고 수칙으로 고정합니다.

### 8.1 저장소 CI 상태는 저장소를 지정해서 확인합니다

이 저장소는 상류 ARTEX 의 포크라서, 로컬 `git remote` 에 `origin`(jiwoochris/artex-global)과
`upstream`(Autumn-27/ARTEX)이 함께 등록되어 있습니다(3절 참조). 이 상태에서 `gh` 명령에
저장소를 지정하지 않으면, `gh` 가 **상류 저장소를 기본값으로 골라** 우리 워크플로가 없는
상류의 실행 결과를 보여 줍니다. 그러면 상류 CI 가 초록인 것을 보고 **우리 CI 가 통과했다고
착각**하거나, 우리 워크플로(`ci.yml`·`detections.yml`)를 "HTTP 404 … not found" 로 잘못
판단할 수 있습니다.

그래서 CI 를 확인할 때는 항상 저장소를 명시합니다.

```bash
gh run list -R jiwoochris/artex-global --workflow ci.yml --limit 5
gh run list -R jiwoochris/artex-global --workflow detections.yml --limit 5
```

한 번 설정해 두면 `-R` 를 생략해도 우리 저장소를 기본으로 보도록 바꿀 수 있습니다. 다만 이
설정은 **로컬 gh 설정**이라 저장소에 커밋되지 않으므로, 새 머신이나 새 체크아웃에서는 다시
지정해야 합니다.

```bash
gh repo set-default jiwoochris/artex-global
gh repo set-default --view   # jiwoochris/artex-global 가 보이는지 확인
```

### 8.2 문서의 외부 링크는 브라우저처럼 GET 으로 확인합니다

방어 가이드([`docs/defense-ko.md`](docs/defense-ko.md)·[`defense-en.md`](docs/defense-en.md))의
7절은 국내 공식 채널(boho.or.kr·fsec.or.kr·pipc.go.kr)의 링크를 싣습니다. 이 링크가 살아
있는지 확인할 때 `curl -I`(HEAD 요청)나 기본 User-Agent 로만 확인하면 **멀쩡한 링크를 깨진
것으로 오인**합니다. 국내 공공·보안 기관 사이트는 다음 세 가지 이유로 단순 확인을 거부하기
때문입니다.

- **HEAD 요청을 거부합니다.** 예를 들어 fsec.or.kr 은 `curl -I`(HEAD)에 400 을 돌려줍니다.
- **기본 `curl` User-Agent 를 차단합니다.** fsec.or.kr 과 pipc.go.kr 은 기본 UA 로 보낸
  GET 요청에도 400 을 돌려줍니다(브라우저 UA 로 보내면 200).
- **다른 주소로 리다이렉트합니다.** pipc.go.kr 은 `www.pipc.go.kr` 에서 `pipc.go.kr/np/` 로
  두 번 리다이렉트하므로, 리다이렉트를 따라가지 않으면 최종 상태를 놓칩니다.

따라서 링크 확인은 **브라우저 User-Agent 로, GET 으로, 리다이렉트를 따라가며** 합니다.

```bash
UA='Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0 Safari/537.36'
for u in https://www.boho.or.kr https://www.fsec.or.kr https://www.pipc.go.kr; do
  curl -sS -L -A "$UA" -o /dev/null -w "$u -> %{http_code} %{url_effective}\n" "$u"
done
```

최종 상태 코드가 200 이면 링크는 유효합니다. 상태 코드가 400·403 으로 나오면 링크가 깨진
것이 아니라 **확인 방법이 서버의 접근 정책에 막힌 것**은 아닌지 먼저 의심하고, HEAD·기본
UA·리다이렉트 미추적 같은 요인을 하나씩 제거해 다시 확인합니다. (2026-10-06 확인 기준으로
세 링크 모두 위 방법에서 200 이며, pipc.go.kr 은 2회 리다이렉트 뒤 200 입니다.)

이 수동 절차는 `scripts/check-external-links.py` 가 그대로 자동화합니다. 추적되는 모든 `.md`
에서 코드펜스·인라인 코드 밖의 외부 링크를 모으고(예약·플레이스홀더 호스트는 제외), 위와
같이 브라우저 UA·GET·리다이렉트 추적으로 상태를 확인하며 네트워크 오류·5xx·429 는 재시도해
일시적 깜빡임과 진짜 장애를 가릅니다. 결과를 네 가지로 나눕니다: OK(2xx·3xx) ·
RESTRICTED(401·403·405·429, 호스트는 살아 있고 확인 방법만 막힘) · ALLOWED
(`scripts/external-links-allowlist.txt` 에 적힌, 우리가 고칠 수 없는 상류 상속 죽은 링크) ·
DOWN(404·410·5xx·연결 오류, 깨졌을 가능성 높음).

- 네트워크 없이 점검 대상만 미리 보기: `python3 -I scripts/check-external-links.py --list`
- 릴리스·주기 점검(새로 깨진 링크가 있으면 비정상 종료): `python3 -I scripts/check-external-links.py --strict`

외부 링크 생존은 flaky 하므로 **머지 게이트에 넣지 않습니다**. 대신 비차단 워크플로
[`external-links`](.github/workflows/external-links.yml) 가 매주 월요일과 수동 실행으로
`--strict` 를 돌려, allowlist 에 없는 DOWN 이 새로 생기면 빨갛게 드러냅니다. 상류 원문 보존
파일이 물려받은 죽은 링크(예: `CHANGELOG.zh.md` 가 크레딧한, 사라진 기여자 계정)는 우리가
고칠 수 없으므로 allowlist 에 사유와 함께 적어 strict 점검에서 뺍니다.

---

## 9. Community 릴리스 발행 파이프라인

버전 원본은 `version/VERSION`이며 현재 정식 릴리스 대상은 **1.0.0**입니다.
상류 기준은 `version/upstream.json`, 릴리스 본문은 `docs/releases/v1.0.0.md`입니다.
로컬 검증과 실제 GitHub 발행 상태는 [완료 기록](docs/complete/2026-10-24-release-readiness.md)에 구분합니다.

### 발행 전 검증

1. 작업 변경을 검토하고 커밋합니다. 인증 파일·`.env`·`.secrets`·`auth`·`work`·빌드 산출물은 커밋하지 않습니다.
2. GitHub Actions의 **release → Run workflow**로 해당 브랜치를 검증합니다.
   수동 실행은 테스트와 5개 플랫폼 ZIP을 만들며 GitHub Release나 GHCR 이미지는 공개하지 않습니다.
3. 실제 OAuth 로그인·모델 선택·연결 테스트를 수행하고 운영 인증 백업과 이전 절차를 확인합니다.
4. 발행할 커밋에 `v1.0.0` 태그를 붙이고 push하면 자동 발행합니다.
   기존 태그를 덮어쓰지 않습니다. 다음 릴리스에서는 버전 파일·변경 기록을 함께 올립니다.

### 자동화의 검증과 산출물

- **validate:** community 저장소, 정식 버전 형식, 태그와 VERSION 일치, 릴리스 문서·상류 JSON 확인.
- **backend:** 외부 서비스 없이 OAuth·업데이트 복구·도메인 테스트, ChatGPT 경쟁 상태 검사,
  패키지별 별도 PostgreSQL 테스트, 인증 이전 도구 검사. 인증 관련 상세 로그는 공개하지 않음.
- **frontend:** 단위 테스트, 5개 언어 키 정합성, 정적 내보내기.
- **binaries:** 위 검증이 모두 성공한 뒤 Linux amd64/arm64, macOS amd64/arm64,
  Windows amd64를 빌드하고 시작 스크립트·skills·VERSION·upstream.json·인증 운영 안내를 ZIP에 포함.
  Linux amd64에서는 실행·버전 확인을 수행. 교차 컴파일만으로 다른 OS 실행까지 증명하지는 않음.
- **docker:** 정식 태그에서만 Linux amd64/arm64 이미지를
  `ghcr.io/hdhdev-droid/artex-community:<태그>`와 `latest`로 발행.
- **release:** 바이너리·이미지 작업 성공 뒤 ZIP, SHA256SUMS, VERSION, upstream.json을 GitHub Release에 첨부.

GHCR 최초 발행 후 공개 설치 대상이면 패키지의 공개 읽기 권한을 확인합니다.
운영에서는 `latest` 대신 검증한 태그를 고정하고,
[인증 보존·이전·롤백 절차](docs/deployment-auth.md)에 따라 앱만 교체합니다.
실행 파일 교체는 데이터베이스 롤백을 제공하지 않습니다.

### 로컬 패키징

Go 1.26.3, Node 22, npm, rsync, zip이 있는 환경에서:

```bash
ARTEX_BUILD_VERSION=v1.0.0 ARTEX_COMPRESS=0 ./build.sh --release
```

이미 정적 내보내기를 `server/webui/dist`에 준비했다면 `ARTEX_SKIP_FRONTEND=1`을 지정해 재사용합니다.
기본 개발 빌드는 `1.0.0-dev`로 표시하여 정식 릴리스와 구분합니다.
ZIP과 체크섬은 `dist`에 생성되며, 운영 데이터나 인증 파일은 패키징하지 않습니다.
