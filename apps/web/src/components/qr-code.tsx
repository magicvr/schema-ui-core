// QR code SVG renderer (W11 · M-01): renders the module matrix of a QR code
// (qrcode-generator, MIT) as inline SVG — no canvas dependency, works in
// jsdom tests and offline builds. Used by the MFA enrollment surface to show
// a scannable QR for the otpauth:// URI.
import { useMemo } from "react";

import qrcode from "qrcode-generator";

export interface QrCodeProps {
  /** Text to encode (e.g. an otpauth:// URI). Empty → renders nothing. */
  value: string;
  /** Rendered size in CSS pixels (the matrix is square). */
  size?: number;
  /** Extra class names for the wrapping figure. */
  className?: string;
  /** Accessible label. */
  label?: string;
}

/**
 * Renders the QR module matrix as SVG. Each dark module becomes a 1×1 rect in
 * a unit viewBox scaled to `size`; a white background rect keeps scanners
 * happy (quiet zone included via the 4-module default margin of the encoder).
 */
export function QrCode({ value, size = 160, className, label = "QR code" }: QrCodeProps) {
  const cells = useMemo(() => {
    if (value === "") {
      return null;
    }
    const qr = qrcode(0, "M");
    qr.addData(value);
    qr.make();
    const count = qr.getModuleCount();
    const dark: Array<{ x: number; y: number }> = [];
    for (let row = 0; row < count; row += 1) {
      for (let col = 0; col < count; col += 1) {
        if (qr.isDark(row, col)) {
          dark.push({ x: col, y: row });
        }
      }
    }
    return { count, dark };
  }, [value]);

  if (cells === null) {
    return null;
  }

  // F-005 (grok A-002): the encoder's module matrix has no quiet zone — add a
  // 4-module white margin around it so scanners can locate the finder patterns.
  const quietZone = 4;
  const extent = cells.count + quietZone * 2;
  return (
    <figure className={className}>
      <svg
        viewBox={`0 0 ${extent} ${extent}`}
        width={size}
        height={size}
        role="img"
        aria-label={label}
        shapeRendering="crispEdges"
        className="rounded-md border border-border bg-background"
        data-qr-code="true"
      >
        <rect width={extent} height={extent} fill="white" />
        {cells.dark.map((cell) => (
          <rect
            key={`${cell.x},${cell.y}`}
            x={cell.x + quietZone}
            y={cell.y + quietZone}
            width={1}
            height={1}
            fill="black"
          />
        ))}
      </svg>
    </figure>
  );
}
