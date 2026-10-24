# [샘플] ARTEX API 인증·인가 진단 보고서

> 화면 확인용 가상 데이터입니다. 실제 진단·외부 접속·개선 완료를 의미하지 않습니다. 한국어 샘플 본문과 인용한 증거 원문은 화면 언어를 바꿔도 유지됩니다.

## 1. 개요

객체 소유권 확인, 토큰 검증, 오류 응답의 정보 노출을 검토하는 API 보고서 샘플입니다. 아래에는 이 작업에 직접 속한 발견만 포함하며, 참조한 다른 작업의 발견을 중복 집계하지 않습니다.

- 작업 ID: t-acme-api
- 작업 상태 원문: running
- 생성 시각 원문: 2026-07-25T14:05:00Z

## 2. 작업 범위와 목표 원문

**작업 설명**

api.acme.com authorization bypass and injection testing

**목표**

Assess authorization bypass (IDOR) and injection risk on the api.acme.com order/user endpoints.

## 3. 발견 현황

총 3건 · 치명적 0 · 높음 2 · 중간 0 · 낮음 1

발견 건수는 현재 데모 발견 목록에서 계산합니다. 위험도와 검증 상태는 별도이며, 아래 상태 원문을 함께 확인하세요.

| ID | 위험도 | 상태 원문 | 유형 |
| --- | --- | --- | --- |
| f-2 | 높음 | pending | IDOR |
| f-3 | 높음 | pending | Weak JWT |
| f-9 | 낮음 | pending | Verbose Error |

## 4. 상세 발견 및 증거

### f-2 · IDOR

#### 발견 요약 원문

Authorization bypass on api.acme.com/v1/orders?id= allows viewing other users' orders

#### 증거 원문

```text
Changing id=1001 to id=1002 returns another user's order (including shipping address and phone number), with no ownership check.
```

### f-3 · Weak JWT

#### 발견 요약 원문

The API JWT uses a weak key that can be cracked and forged offline

#### 증거 원문

```text
HS256, key 'secret', cracked with john in 5 seconds → an arbitrary sub can be forged to bypass authorization.
```

### f-9 · Verbose Error

#### 발견 요약 원문

The API returns a stack trace on 500 errors, exposing paths and framework version

#### 증거 원문

```text
Triggering a 500 error returns a Node.js stack trace, exposing absolute paths and dependency versions.
```

## 5. 개선 및 재검증 계획

- 우선: 주문·사용자 조회마다 인증된 주체와 객체 소유권을 서버에서 대조합니다.
- 단기: 토큰 서명·유효기간·발급자·대상을 검증하고 상세 내부 오류를 외부에 반환하지 않습니다.
- 재검증: 서로 다른 권한의 테스트 계정으로 허용·거부 사례를 비교합니다.

## 6. 한계와 인수 기준

- 권고 조치는 계획이며 수행 완료로 표시하지 않습니다.
- 실제 증거 수집, 수정 적용, 재검증 결과와 담당자 검토가 있어야 조치를 종료할 수 있습니다.
- 기존 개별 취약점 보고서·사용자 입력·증거는 덮어쓰지 않습니다.
