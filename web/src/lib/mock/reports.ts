import type { Finding, Severity, Task } from "../types";

// Authored Korean sample documents, not UI translations or real assessment results.
// Keep the underlying task/finding/evidence fixtures unchanged.
const samples: Record<string, { title: string; summary: string; actions: string[] }> = {
  "t-acme-web": {
    title: "외부 웹·내부 연계 진단",
    summary:
      "외부 웹 서비스의 입력 검증·관리자 인증 문제와 내부 접근 경로를 함께 보여 주는 종합 보고서 샘플입니다. 발견 상태가 pending인 항목은 검토 대기이며, 실제 검증 완료로 해석하지 않습니다.",
    actions: [
      "우선: 관리자 접근 경로와 고위험 실행 경로를 제한하고 노출된 인증 정보를 교체합니다.",
      "단기: 입력 검증과 파라미터화 쿼리를 적용하고 불필요한 소스·디버그 정보를 제거합니다.",
      "후속: 내부망 분리와 계정 권한을 점검한 뒤, 동일 범위로 재검증합니다.",
    ],
  },
  "t-acme-api": {
    title: "API 인증·인가 진단",
    summary:
      "객체 소유권 확인, 토큰 검증, 오류 응답의 정보 노출을 검토하는 API 보고서 샘플입니다. 아래에는 이 작업에 직접 속한 발견만 포함하며, 참조한 다른 작업의 발견을 중복 집계하지 않습니다.",
    actions: [
      "우선: 주문·사용자 조회마다 인증된 주체와 객체 소유권을 서버에서 대조합니다.",
      "단기: 토큰 서명·유효기간·발급자·대상을 검증하고 상세 내부 오류를 외부에 반환하지 않습니다.",
      "재검증: 서로 다른 권한의 테스트 계정으로 허용·거부 사례를 비교합니다.",
    ],
  },
  "t-shop-pay": {
    title: "결제·주문 흐름 중간 점검",
    summary:
      "일시정지된 작업의 중간 보고서 샘플입니다. 발견이 적거나 검사가 중단되었다는 이유로 결제·주문 흐름 전체를 안전하다고 판정하지 않습니다.",
    actions: [
      "우선: 식별된 구성 요소의 버전과 실제 영향을 별도로 확인합니다.",
      "재개 전: 결제 금액 검증, 주문 소유권, 중복 요청 처리의 점검 범위를 확정합니다.",
      "후속: 승인된 테스트 데이터로 재검증한 뒤 잔여 위험을 갱신합니다.",
    ],
  },
  "t-vpn-edge": {
    title: "외부 경계 서비스 점검",
    summary:
      "완료된 작업에 등록된 발견이 없는 경우의 보고서 샘플입니다. 0건은 현재 데모 목록에 발견이 없다는 의미이며, 모든 경계 서비스의 안전성이나 점검 완전성을 보증하지 않습니다.",
    actions: [
      "검토: 허용한 범위와 실제 수집 자료가 일치하는지 확인합니다.",
      "후속: 서비스 변경이나 접근 정책 변경 시 동일 기준으로 다시 점검합니다.",
    ],
  },
};

const severityNames: Record<Severity, string> = {
  critical: "치명적",
  high: "높음",
  medium: "중간",
  low: "낮음",
};
const cell = (value: string) => value.replace(/\|/g, "\\|").replace(/[\r\n]+/g, " ");
const verbatim = (value: string) => {
  const fence = "`".repeat(Math.max(3, ...Array.from(value.matchAll(/`+/g), (match) => match[0].length + 1)));
  return `${fence}text\n${value}\n${fence}`;
};

export function sampleTaskReport(task: Task | undefined, findings: Finding[]): string {
  const sample = task && Object.hasOwn(samples, task.id) ? samples[task.id] : undefined;
  if (!task || !sample) return "";
  const items = findings.filter((finding) => finding.task_id === task.id);
  const counts: Record<Severity, number> = { critical: 0, high: 0, medium: 0, low: 0 };
  for (const finding of items) counts[finding.severity]++;
  const table = items.length
    ? [
        "| ID | 위험도 | 상태 원문 | 유형 |",
        "| --- | --- | --- | --- |",
        ...items.map(
          (finding) =>
            `| ${cell(finding.id)} | ${severityNames[finding.severity]} | ${cell(finding.status)} | ${cell(finding.vulnclass)} |`,
        ),
      ].join("\n")
    : "이 작업에 등록된 발견이 없습니다. 미검증 영역은 별도 확인이 필요합니다.";
  const details = items
    .map((finding) =>
      [
        `### ${finding.id} · ${finding.vulnclass}`,
        "#### 발견 요약 원문",
        finding.summary,
        "#### 증거 원문",
        finding.evidence ? verbatim(finding.evidence) : "등록된 증거 없음",
      ].join("\n\n"),
    )
    .join("\n\n");
  return [
    `# [샘플] ARTEX ${sample.title} 보고서`,
    "> 화면 확인용 가상 데이터입니다. 실제 진단·외부 접속·개선 완료를 의미하지 않습니다. 한국어 샘플 본문과 인용한 증거 원문은 화면 언어를 바꿔도 유지됩니다.",
    "## 1. 개요",
    sample.summary,
    `- 작업 ID: ${task.id}\n- 작업 상태 원문: ${task.status}\n- 생성 시각 원문: ${task.created_at}`,
    "## 2. 작업 범위와 목표 원문",
    `**작업 설명**\n\n${task.description}\n\n**목표**\n\n${task.goal}`,
    "## 3. 발견 현황",
    `총 ${items.length}건 · 치명적 ${counts.critical} · 높음 ${counts.high} · 중간 ${counts.medium} · 낮음 ${counts.low}`,
    "발견 건수는 현재 데모 발견 목록에서 계산합니다. 위험도와 검증 상태는 별도이며, 아래 상태 원문을 함께 확인하세요.",
    table,
    "## 4. 상세 발견 및 증거",
    details || "인용할 발견·증거가 없습니다.",
    "## 5. 개선 및 재검증 계획",
    sample.actions.map((action) => `- ${action}`).join("\n"),
    "## 6. 한계와 인수 기준",
    "- 권고 조치는 계획이며 수행 완료로 표시하지 않습니다.\n- 실제 증거 수집, 수정 적용, 재검증 결과와 담당자 검토가 있어야 조치를 종료할 수 있습니다.\n- 기존 개별 취약점 보고서·사용자 입력·증거는 덮어쓰지 않습니다.",
  ].join("\n\n");
}
