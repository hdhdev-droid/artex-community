# [샘플] ARTEX 외부 웹·내부 연계 진단 보고서

> 화면 확인용 가상 데이터입니다. 실제 진단·외부 접속·개선 완료를 의미하지 않습니다. 한국어 샘플 본문과 인용한 증거 원문은 화면 언어를 바꿔도 유지됩니다.

## 1. 개요

외부 웹 서비스의 입력 검증·관리자 인증 문제와 내부 접근 경로를 함께 보여 주는 종합 보고서 샘플입니다. 발견 상태가 pending인 항목은 검토 대기이며, 실제 검증 완료로 해석하지 않습니다.

- 작업 ID: t-acme-web
- 작업 상태 원문: running
- 생성 시각 원문: 2026-07-24T09:12:00Z

## 2. 작업 범위와 목표 원문

**작업 설명**

External penetration test of Acme's official site and admin back office (acme.com)

**목표**

Gain administrator access to the acme.com back office and determine whether sensitive user data can be read.

## 3. 발견 현황

총 13건 · 치명적 2 · 높음 6 · 중간 4 · 낮음 1

발견 건수는 현재 데모 발견 목록에서 계산합니다. 위험도와 검증 상태는 별도이며, 아래 상태 원문을 함께 확인하세요.

| ID | 위험도 | 상태 원문 | 유형 |
| --- | --- | --- | --- |
| f-1 | 높음 | pending | SQL Injection |
| f-4 | 중간 | pending | Reflected XSS |
| f-5 | 중간 | pending | Exposed .git |
| f-6 | 높음 | pending | Default Credentials |
| f-7 | 낮음 | pending | Open Redirect |
| f-8 | 중간 | pending | Missing Rate Limit |
| f-11 | 높음 | pending | Hardcoded Credentials |
| f-12 | 치명적 | confirmed | Deserialization RCE |
| f-13 | 높음 | pending | Privilege Escalation |
| f-14 | 중간 | pending | Network Segmentation |
| f-15 | 높음 | pending | Unauthenticated RCE |
| f-16 | 높음 | pending | Kerberoasting |
| f-17 | 치명적 | confirmed | Domain Compromise |

## 4. 상세 발견 및 증거

### f-1 · SQL Injection

#### 발견 요약 원문

Error-based SQL injection exists in the q parameter of www.acme.com/search

#### 증거 원문

```text
GET /search?q=1' AND 1=CONVERT(int,@@version)-- returns an MSSQL version error message, allowing the database structure to be read.
```

### f-4 · Reflected XSS

#### 발견 요약 원문

The search page does not escape the q parameter, causing reflected XSS

#### 증거 원문

```text
q=<script>alert(document.domain)</script> is reflected as-is in the result title.
```

### f-5 · Exposed .git

#### 발견 요약 원문

www.acme.com exposes a .git directory, allowing source code recovery

#### 증거 원문

```text
GET /.git/HEAD → 200; recovered backend source code and a database connection string comment using git-dumper.
```

### f-6 · Default Credentials

#### 발견 요약 원문

admin.acme.com admin back office uses the default password admin/admin123

#### 증거 원문

```text
Logged in successfully and entered the admin back office, where users and orders can be managed.
```

### f-7 · Open Redirect

#### 발견 요약 원문

The next parameter after login allows redirecting to an arbitrary address

#### 증거 원문

```text
Logging in with /login?next=https://evil.example results in a 302 redirect to the external site.
```

### f-8 · Missing Rate Limit

#### 발견 요약 원문

The login endpoint has no rate limiting, allowing brute-force attacks

#### 증거 원문

```text
Sending 1,000 requests per minute does not trigger a lockout, and there is no CAPTCHA.
```

### f-11 · Hardcoded Credentials

#### 발견 요약 원문

Hardcoded database credential sa/Acme@2021 found in leaked source code

#### 증거 원문

```text
Recovered config.php by restoring www.acme.com/.git with git-dumper; it contained a plaintext DB password, later confirmed reusable on the internal network.
```

### f-12 · Deserialization RCE

#### 발견 요약 원문

shop.acme.com Fastjson 1.2.24 deserialization remote code execution

#### 증거 원문

```text
Sending a JNDI payload (@type:JdbcRowSetImpl) via POST /api/import triggers a dnslog callback and drops a reverse shell (www-data@dmz-web01).
```

### f-13 · Privilege Escalation

#### 발견 요약 원문

Local privilege escalation to root from the DMZ foothold via sudo NOPASSWD misconfiguration

#### 증거 원문

```text
sudo -l shows (ALL) NOPASSWD: /usr/bin/python3 → sudo python3 -c 'os.setuid(0);os.system("/bin/sh")' immediately obtains root.
```

### f-14 · Network Segmentation

#### 발견 요약 원문

A DMZ host has dual NICs directly connected to the internal network, defeating segmentation

#### 증거 원문

```text
dmz-web01's second NIC is on 10.10.10.0/24, giving direct reach from the DMZ to the domain controller/file server/Jenkins and bypassing network segmentation.
```

### f-15 · Unauthenticated RCE

#### 발견 요약 원문

Unauthenticated Script Console remote code execution on internal Jenkins

#### 증거 원문

```text
Executed Groovy without logging in via GET /script: 'whoami'.execute().text → SYSTEM; extracted the domain account acme\svc_deploy from the credential store.
```

### f-16 · Kerberoasting

#### 발견 요약 원문

Domain service account svc_sql is Kerberoastable and has a weak password

#### 증거 원문

```text
Requested svc_sql's TGS via GetUserSPNs and cracked Sql@2020 offline with hashcat -m 13100. This account belongs to the SQL admin group.
```

### f-17 · Domain Compromise

#### 발견 요약 원문

Domain controller DC01 compromised (Domain Admin): internal objective achieved

#### 증거 원문

```text
svc_deploy belongs to Domain Admins; logged into DC01 via psexec.py to obtain SYSTEM, then used secretsdump to extract the entire domain's NTLM hashes (including krbtgt), completing the objective.
```

## 5. 개선 및 재검증 계획

- 우선: 관리자 접근 경로와 고위험 실행 경로를 제한하고 노출된 인증 정보를 교체합니다.
- 단기: 입력 검증과 파라미터화 쿼리를 적용하고 불필요한 소스·디버그 정보를 제거합니다.
- 후속: 내부망 분리와 계정 권한을 점검한 뒤, 동일 범위로 재검증합니다.

## 6. 한계와 인수 기준

- 권고 조치는 계획이며 수행 완료로 표시하지 않습니다.
- 실제 증거 수집, 수정 적용, 재검증 결과와 담당자 검토가 있어야 조치를 종료할 수 있습니다.
- 기존 개별 취약점 보고서·사용자 입력·증거는 덮어쓰지 않습니다.
