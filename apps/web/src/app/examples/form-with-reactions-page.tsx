import { useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { RenderPage } from "@/renderer/render.tsx";
import type { RenderPageDocument } from "@/renderer/render";

// D-EXPR example: reactions reuse the frozen $context expression engine
// (evaluateExpression) to flip field visibility / disabled state. The page
// document below is rendered through the minimal Renderer (D-COMP).
const REACTIVE_FORM_DOCUMENT: RenderPageDocument = {
  meta: {
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation", "form.controls.extended"],
  },
  body: {
    type: "section",
    children: [
      {
        type: "form",
        id: "reactive-form",
        props: {
          fields: [
            { id: "name", label: "Name (input)", type: "input" },
            {
              id: "kind",
              label: "Kind (select)",
              type: "select",
              options: [
                { value: "standard", label: "Standard" },
                { value: "priority", label: "Priority" },
              ],
            },
            { id: "approval", label: "Approval (switch)", type: "switch" },
            { id: "auditNote", label: "Audit note (textarea)", type: "textarea" },
          ],
          reactions: [
            // Visible only for admins.
            {
              id: "admin-approval",
              when: '$context.user.roles contains "admin"',
              apply: [{ fieldId: "approval", visible: true }],
            },
            {
              id: "viewer-approval",
              when: '$context.user.roles contains "viewer"',
              apply: [{ fieldId: "approval", visible: false }],
            },
            // Disabled unless the audit feature is on.
            {
              id: "audit-lock",
              when: "$context.features.audit != true",
              apply: [{ fieldId: "auditNote", disabled: true }],
            },
          ],
        },
      },
    ],
  },
};

export function FormWithReactionsPage() {
  const [isAdmin, setIsAdmin] = useState(true);
  const [auditOn, setAuditOn] = useState(true);

  const context = useMemo(
    () => ({
      user: { roles: isAdmin ? ["admin"] : ["viewer"] },
      features: { audit: auditOn },
    }),
    [isAdmin, auditOn],
  );

  return (
    <section className="space-y-6" aria-labelledby="form-with-reactions-title">
      <div className="space-y-2">
        <h1 id="form-with-reactions-title" className="text-3xl font-semibold tracking-tight">
          Form with reactions
        </h1>
        <p className="max-w-2xl text-sm leading-6 text-muted-foreground">
          D-EXPR reactions reuse the frozen $context expression engine
          (<code className="font-mono">evaluateExpression</code>). The role /
          feature toggles change the snapshot and the renderer re-applies
          field visibility and disabled state.
        </p>
      </div>

      <div className="flex flex-wrap items-center gap-3 border border-border bg-card p-4">
        <span className="text-xs font-medium uppercase tracking-[0.12em] text-muted-foreground">
          Context snapshot
        </span>
        <Button
          type="button"
          variant={isAdmin ? "default" : "outline"}
          size="sm"
          onClick={() => setIsAdmin(true)}
        >
          Admin
        </Button>
        <Button
          type="button"
          variant={!isAdmin ? "default" : "outline"}
          size="sm"
          onClick={() => setIsAdmin(false)}
        >
          Viewer
        </Button>
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={auditOn}
            onChange={(event) => setAuditOn(event.target.checked)}
          />
          audit feature
        </label>
        <code className="ml-auto rounded bg-muted px-2 py-1 text-xs">
          {JSON.stringify(context)}
        </code>
      </div>

      <RenderPage document={REACTIVE_FORM_DOCUMENT} context={context} />
    </section>
  );
}
